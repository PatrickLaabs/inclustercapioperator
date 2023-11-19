# In-Cluster ClusterAPI Operator

Various Operations possible through ClusterAPI CRDs within a
go program.

## Limitations

- It is possbile to get a set of provisioned clusters on the management-cluster from the ClusterAPI API.
- Get the available kubeconfigs for your cluster. currently, its only returning the values.