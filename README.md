# cb-indexcreation-operator
This operator allows to create Couchbase query indexes declaratively.

Description:
What you can do with the operator:
1. Create both primary a and regular indexes using yamls. There are examples in /config/samples for both use cases
2. Connection details to the cluster provided by external secret (configurable)

What you cannot do with the operator:
1. There is no update/delete index support as for now.
2. Batch creation is not supported for now. You need to define yaml per index.
3. Only user name/password authentication is supported.
4. There is no support to reconciliation with the changes out of the operator.
5. There is no security scanning on input parameters. So SQL injection is not prevented
6. Deffered indexes creation is not supported yet
7. Only kustomization templatization is supported - there is no helm charts support

Instructions
1. Clone the repo
2. Run 'make deploy' to deploy it into your cluster. The access to K8s cluster should be configured already. The operator is deployed from my DockerHub account.
3. Or, if you want to have your own image - 
		change the make file 
		'make build' 
		'make docker-build docker-push' 
		and than 'make deploy'
