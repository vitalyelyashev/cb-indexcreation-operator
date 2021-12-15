/*
Copyright 2021 Vitaly Elyashev.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"time"

	core "k8s.io/api/core/v1"

	cachev1alpha1 "forescout.com/m/api/v1alpha1"
	"github.com/couchbase/gocb/v2"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CouchbaseIndexReconciler reconciles a CouchbaseIndex object
type CouchbaseIndexReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	mylog = ctrl.Log.WithName("Reconcile")
)

//+kubebuilder:rbac:groups=cache.forescout.com,resources=couchbaseindices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.forescout.com,resources=couchbaseindices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.forescout.com,resources=couchbaseindices/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CouchbaseIndex object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *CouchbaseIndexReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	mylog.Info("Start reconciliation")

	couchbaseIndex := &cachev1alpha1.CouchbaseIndex{}

	err := r.Get(ctx, req.NamespacedName, couchbaseIndex)
	if err != nil {
		mylog.Error(err, "Custom Resource Not exists")
		return ctrl.Result{}, err
	}

	mylog.Info("Succesfully assigned couchbase index")

	url, username, password, err := GetSecretData(r, ctx, req, couchbaseIndex.Spec.SecretName)

	if err != nil {
		mylog.Error(err, "Cannot recieve data from secret")

		return ctrl.Result{}, err
	}

	indexData := couchbaseIndex.Spec.Index

	//Connect to Couchbase Cluster
	cluster, err := gocb.Connect(
		url,
		gocb.ClusterOptions{
			Username: username,
			Password: password,
		})
	if err != nil {
		mylog.Error(err, "Cannot connect to the cluster")

		return ctrl.Result{}, err
	}

	mylog.Info("Succesfully connected to couchbase cluster", "ConnectionUrl:", url, "Username:", username, "Password:", password)

	if !controllerutil.ContainsFinalizer(couchbaseIndex, cachev1alpha1.IndexFinalizer) {
		controllerutil.AddFinalizer(couchbaseIndex, cachev1alpha1.IndexFinalizer)
		if err := r.Update(ctx, couchbaseIndex); err != nil {
			mylog.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		} else {
			mylog.Info("Successfully added a finalizer")
		}
	}

	qi := cluster.QueryIndexes()
	isIndexMarkedtoBeDeleted := couchbaseIndex.GetDeletionTimestamp() != nil

	if isIndexMarkedtoBeDeleted {

		return ReconcileDelete(r, ctx, indexData, qi, couchbaseIndex)
	}

	mylog.Info("Index status: ", "current status", couchbaseIndex.Status.Type)

	switch conditionType := couchbaseIndex.Status.Type; conditionType {

	case cachev1alpha1.ConditionInProgress:
		mylog.Info("the status is InProgress - Check if Index already created")

		isExist, index, err := IsIndexExist(indexData.IndexName, indexData.BucketName, indexData.DeepCopy().IsPrimary, cluster)
		if err != nil {
			mylog.Error(err, "Failed to get index")

			couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
			r.Status().Update(ctx, couchbaseIndex)
			return ctrl.Result{}, err
		}

		if !isExist {
			mylog.Info("The index is not created yet - retry later")

			return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
		} else {
			UpdateIndexIfNeeded(index, indexData, qi)
			couchbaseIndex.Status.Type = cachev1alpha1.ConditionReady
			r.Status().Update(ctx, couchbaseIndex)
			return ctrl.Result{}, nil
		}
	case cachev1alpha1.ConditionReady:
	case cachev1alpha1.ConditionFailed:
		err = CreateIndex(indexData, cluster, couchbaseIndex, r, ctx, qi)
		if err != nil {
			mylog.Error(err, "Failed to create index")

			couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
			r.Status().Update(ctx, couchbaseIndex)
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	default:
		err = CreateIndex(indexData, cluster, couchbaseIndex, r, ctx, qi)
		if err != nil {
			mylog.Error(err, "Failed to create index")

			couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
			r.Status().Update(ctx, couchbaseIndex)
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}
	mylog.Info("Reconcile loop complete")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CouchbaseIndexReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.CouchbaseIndex{}).
		Complete(r)
}

func ReconcileDelete(r *CouchbaseIndexReconciler,
	ctx context.Context,
	indexData cachev1alpha1.IndexData,
	qi *gocb.QueryIndexManager,
	couchbaseIndex *cachev1alpha1.CouchbaseIndex) (ctrl.Result, error) {
	mylog.Info("Index should be deleted")
	var err error
	if indexData.IsPrimary {
		err = qi.DropPrimaryIndex(indexData.BucketName, &gocb.DropPrimaryQueryIndexOptions{})
	} else {
		err = qi.DropIndex(indexData.BucketName, indexData.IndexName, &gocb.DropQueryIndexOptions{})
	}

	if err != nil {
		mylog.Error(err, "Failed to delete a index")
	} else {
		mylog.Info("Succesfully removed index.", "IndexName:", indexData.IndexName)
		controllerutil.RemoveFinalizer(couchbaseIndex, cachev1alpha1.IndexFinalizer)
		if err := r.Update(ctx, couchbaseIndex); err != nil {
			mylog.Error(err, "unable to remove finalizer")
			return ctrl.Result{}, err
		} else {
			mylog.Info("Successfully removed a finalizer")
		}
	}
	return ctrl.Result{}, nil
}
func CreateIndex(
	indexData cachev1alpha1.IndexData,
	cluster *gocb.Cluster,
	couchbaseIndex *cachev1alpha1.CouchbaseIndex,
	r *CouchbaseIndexReconciler,
	ctx context.Context,
	qi *gocb.QueryIndexManager) (err error) {
	isExist, index, err := IsIndexExist(indexData.IndexName, indexData.BucketName, indexData.IsPrimary, cluster)
	mylog.Info("the status is Not Exists - try to create an index")
	if err != nil {
		mylog.Error(err, "Failed to get an index")
		couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
		r.Status().Update(ctx, couchbaseIndex)
		return err
	}

	if !isExist {
		//create new Index

		mylog.Info("Create index asyncronousely. Status will be changed to InProgress")
		if couchbaseIndex.Spec.Index.IsPrimary {
			if err := qi.CreatePrimaryIndex(indexData.BucketName, &gocb.CreatePrimaryQueryIndexOptions{}); err != nil {
				if errors.Is(err, gocb.ErrIndexExists) {
					mylog.Info("Primary Index already exists")
				} else {
					couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
					r.Status().Update(ctx, couchbaseIndex)
					mylog.Error(err, "Failed to create primary index")
					return err
				}
			}
		} else {
			if err := qi.CreateIndex(indexData.BucketName, indexData.IndexName, indexData.Parameters, &gocb.CreateQueryIndexOptions{}); err != nil {
				if errors.Is(err, gocb.ErrIndexExists) {
					mylog.Info("Index already exists")
				} else {
					couchbaseIndex.Status.Type = cachev1alpha1.ConditionFailed
					r.Status().Update(ctx, couchbaseIndex)
					mylog.Error(err, "Failed to create index")
					return err
				}
			}
		}

		couchbaseIndex.Status.Type = cachev1alpha1.ConditionInProgress
		r.Status().Update(ctx, couchbaseIndex)

		return nil
	} else {
		UpdateIndexIfNeeded(index, indexData, qi)
		mylog.Info("Index already exists")
		couchbaseIndex.Status.Type = cachev1alpha1.ConditionReady
		r.Status().Update(ctx, couchbaseIndex)
	}
	return nil
}
func UpdateIndexIfNeeded(index gocb.QueryIndex, indexData cachev1alpha1.IndexData, qi *gocb.QueryIndexManager) {

	//qi.CreateIndex()
	mylog.Info("The index is already exists", "Condition", index.Condition, "Key:", index.IndexKey, "Type:", index.Type)
}
func IsIndexExist(name string, bucket string, isPrimary bool, cluster *gocb.Cluster) (found bool, index gocb.QueryIndex, err error) {
	qi := cluster.QueryIndexes()
	if isPrimary {
		name = "#primary"
	}
	all, err := qi.GetAllIndexes(bucket, &gocb.GetAllQueryIndexesOptions{})
	if err != nil {
		return
	}

	for _, i := range all {
		if i.Name == name {
			found = true
			index = i
			break
		}
	}
	return

}

func GetSecretData(r *CouchbaseIndexReconciler,
	ctx context.Context,
	req ctrl.Request,
	name string) (url string, username string, password string, err error) {
	secret := &core.Secret{}
	mylog.Info("trying to get secret from the K8s ", "secret name", name, "namespace", req.Namespace)

	err = r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: req.Namespace}, secret)
	if err != nil {
		mylog.Error(err, "failed to get a secret")
		return
	}
	mylog.Info("Succesfully recieved secrets")

	url = string(secret.Data["url"][:])
	username = string(secret.Data["username"][:])
	password = string(secret.Data["password"][:])
	return
}
