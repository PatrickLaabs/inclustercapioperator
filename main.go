package inclustercapioperator

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetWorkloadClusters() ([]string, error) {
	// Use in-cluster config if running inside Kubernetes, otherwise use kubeconfig file
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernetes client: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernetes client: %v", err)
	}

	// Replace "default" and "management-prod-cluster-kubeconfig" with your Secret's namespace and name
	secret, err := clientset.CoreV1().Secrets("default").Get(context.TODO(), "management-prod-cluster-kubeconfig", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error getting kubeconfig from Secret: %v", err)
	}

	kubeconfig, ok := secret.Data["kubeconfig"]
	if !ok {
		return nil, errors.New("kubeconfig not found in Secret")
	}

	// Use the retrieved kubeconfig
	config, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Error building kubeconfig: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating dynamic client: %v", err)
	}

	apiResource := metav1.APIResource{
		Name: "clusters",
		Kind: "Cluster",
	}

	gvResouce := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1", // Replace with the actual version
		Resource: apiResource.Name,
	}

	clusters, err := dynamicClient.Resource(gvResouce).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing clusters: %v", err)
	}

	for _, cluster := range clusters.Items {
		fmt.Printf("Cluster: %s\n", cluster.GetName())
	}

	var clusterNames []string
	for _, cluster := range clusters.Items {
		clusterNames = append(clusterNames, cluster.GetName())
	}

	return clusterNames, nil
}
