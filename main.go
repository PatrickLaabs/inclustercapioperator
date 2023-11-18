package main

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
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

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// TODO: Interact with the Cluster API to list clusters

	// Example: List all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing namespaces: %v\n", err)
		os.Exit(1)
	}

	for _, ns := range namespaces.Items {
		fmt.Printf("Namespace: %s\n", ns.Name)
	}
}
