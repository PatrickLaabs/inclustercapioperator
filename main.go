package github.com/PatrickLaabs/inclustercapioperator

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetWorkloadClusters() {
	// Use in-cluster config if running inside Kubernetes, otherwise use kubeconfig file
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := "/home/patrick/Development/management-prod-cluster.kubeconfig"
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("Error building kubeconfig: %v\n", err)
			os.Exit(1)
		}
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
}
