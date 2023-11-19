package inclustercapioperator

import (
	"context"
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterOperator represents the operator for managing clusters.
type ClusterOperator struct {
	KubeconfigPath string
	Clientset      *kubernetes.Clientset
	DynamicClient  dynamic.Interface

	APIResource      metav1.APIResource
	GVResource       schema.GroupVersionResource
	DefaultNamespace string
}

func NewClusterOperator(kubeconfigPath string) (*ClusterOperator, error) {
	var config *rest.Config
	var err error

	// Use in-cluster config if running inside Kubernetes
	if kubeconfigPath == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating in-cluster Kubernetes client: %v", err)
		}
	} else {
		// Use kubeconfig from the specified path for local development
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("Error building kubeconfig: %v", err)
		}
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

	kubeconfig, ok := secret.Data["value"]
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

	gvResource := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1", // Replace with the actual version
		Resource: apiResource.Name,
	}

	return &ClusterOperator{
		KubeconfigPath:   "",
		Clientset:        clientset,
		DynamicClient:    dynamicClient,
		APIResource:      apiResource,
		GVResource:       gvResource,
		DefaultNamespace: "default",
	}, nil
}

func (co *ClusterOperator) GetWorkloadClusters() ([]string, error) {
	clusters, err := co.DynamicClient.Resource(co.GVResource).Namespace("default").List(context.TODO(), metav1.ListOptions{})
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

func (co *ClusterOperator) GetKubernetesSecrets() ([]string, error) {
	// List all secrets in the default namespace
	secrets, err := co.Clientset.CoreV1().Secrets(co.DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing secrets: %v", err)
	}

	var kubeconfigSecretsNames []string
	// Filter secrets based on their names
	for _, secret := range secrets.Items {
		// Check if the secret name ends with "-kubeconfig"
		if strings.HasSuffix(secret.GetName(), "-kubeconfig") {
			kubeconfigSecretsNames = append(kubeconfigSecretsNames, secret.GetName())
		}
	}

	return kubeconfigSecretsNames, nil
}

func (co *ClusterOperator) GetMgmtIngresses() ([]string, error) {
	ingresses, err := co.Clientset.NetworkingV1().Ingresses(co.DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Eror listing ingresses from all namespaces: %v\n", err)
	}

	var kubernetesIngresses []string
	for _, ingress := range ingresses.Items {
		kubernetesIngresses = append(kubernetesIngresses, ingress.GetName())
	}

	return kubernetesIngresses, nil
}
