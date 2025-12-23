// Package client provides Kubernetes REST client utilities for resource management.
//
// This package simplifies interactions with the Kubernetes API server by providing
// type-safe generic functions for common CRUD operations on Kubernetes resources.
//
// Features:
//   - In-cluster and external client configuration
//   - Generic resource Get operations
//   - Resource Create (Post) operations
//   - Resource Update (Put) operations
//   - Resource Delete operations
//   - Support for custom resources
//
// Example:
//
//	// In-cluster
//	restClient, _ := client.GetClientForInCluster()
//
//	// Get ConfigMap
//	cm, err := client.Get[v1.ConfigMap](restClient, "", "v1", "default", "configmaps", "my-config")
package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetClientForInCluster creates a Kubernetes REST client for in-cluster use.
//
// This function is designed to be called from within a Kubernetes pod and uses
// the service account token and cluster CA certificate mounted by Kubernetes.
//
// Returns:
//   - rest.Interface: Kubernetes REST client interface
//   - error: Error if client creation fails
//
// The function performs these steps:
//  1. Loads in-cluster configuration from environment and mounted files
//  2. Creates a Kubernetes clientset using the configuration
//  3. Returns the REST client from the clientset
//
// In-cluster configuration requirements:
//   - KUBERNETES_SERVICE_HOST environment variable
//   - KUBERNETES_SERVICE_PORT environment variable
//   - Service account token at /var/run/secrets/kubernetes.io/serviceaccount/token
//   - CA certificate at /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
//
// Example:
//
//	// Inside a Kubernetes pod
//	restClient, err := client.GetClientForInCluster()
//	if err != nil {
//	    log.Fatal("Failed to create in-cluster client:", err)
//	}
//
//	// Use client to get a ConfigMap
//	configMap, err := client.Get[coreV1.ConfigMap](
//	    restClient,
//	    "",        // group (empty for core)
//	    "v1",      // version
//	    "default", // namespace
//	    "configmaps",
//	    "my-config",
//	)
//
// Note: This function will fail when run outside a Kubernetes cluster.
// Use GetClientUsingConfig for external access.
func GetClientForInCluster() (rest.Interface, error) {
	if config, err := rest.InClusterConfig(); err != nil {
		return nil, err
	} else if clientset, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return clientset.RESTClient(), nil
	}
}

// GetClientUsingConfig creates a Kubernetes REST client using a provided configuration.
//
// This function is designed for clients running outside a Kubernetes cluster,
// such as CLI tools or external applications.
//
// Parameters:
//   - config: Kubernetes REST client configuration
//
// Returns:
//   - rest.Interface: Kubernetes REST client interface
//   - error: Error if client creation fails
//
// The configuration typically includes:
//   - API server URL
//   - Authentication credentials (token, certificate, or username/password)
//   - TLS settings
//   - Timeout settings
//
// Example with kubeconfig:
//
//	import "k8s.io/client-go/tools/clientcmd"
//
//	// Load kubeconfig file
//	config, err := clientcmd.BuildConfigFromFlags("", "/path/to/kubeconfig")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create client
//	restClient, err := client.GetClientUsingConfig(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use client
//	pods, _ := client.Get[coreV1.PodList](restClient, "", "v1", "default", "pods", "")
//
// Example with manual configuration:
//
//	config := &rest.Config{
//	    Host:        "https://k8s.example.com:6443",
//	    BearerToken: "your-token-here",
//	    TLSClientConfig: rest.TLSClientConfig{
//	        Insecure: false,
//	        CAFile:   "/path/to/ca.crt",
//	    },
//	}
//
//	restClient, err := client.GetClientUsingConfig(config)
func GetClientUsingConfig(config *rest.Config) (rest.Interface, error) {
	if clientset, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return clientset.RESTClient(), nil
	}
}

// Get retrieves a Kubernetes resource by name and returns it as type T.
//
// This generic function fetches a single Kubernetes resource from the API server
// and unmarshals it into the specified Go type.
//
// Type Parameters:
//   - T: Target Go type (must match the resource being retrieved)
//
// Parameters:
//   - restClient: Kubernetes REST client interface
//   - group: API group (empty "" for core resources like Pod, ConfigMap)
//   - version: API version (e.g., "v1", "v1beta1")
//   - namespace: Namespace name (empty "" for cluster-scoped resources)
//   - resource: Resource type in plural form (e.g., "pods", "configmaps", "deployments")
//   - name: Resource name
//
// Returns:
//   - T: Retrieved resource as type T
//   - error: Error if request fails or unmarshaling fails
//
// The function constructs the appropriate API path:
//   - Core resources: /api/{version}/namespaces/{namespace}/{resource}/{name}
//   - Group resources: /apis/{group}/{version}/namespaces/{namespace}/{resource}/{name}
//   - Cluster-scoped: /api(s)/{group}/{version}/{resource}/{name}
//
// Example - Get ConfigMap:
//
//	import coreV1 "k8s.io/api/core/v1"
//
//	configMap, err := client.Get[coreV1.ConfigMap](
//	    restClient,
//	    "",           // group (empty for core)
//	    "v1",         // version
//	    "default",    // namespace
//	    "configmaps", // resource
//	    "app-config", // name
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("ConfigMap data: %v\n", configMap.Data)
//
// Example - Get Deployment:
//
//	import appsV1 "k8s.io/api/apps/v1"
//
//	deployment, err := client.Get[appsV1.Deployment](
//	    restClient,
//	    "apps",         // group
//	    "v1",           // version
//	    "production",   // namespace
//	    "deployments",  // resource
//	    "web-app",      // name
//	)
//
// Example - Get Node (cluster-scoped):
//
//	node, err := client.Get[coreV1.Node](
//	    restClient,
//	    "",      // group
//	    "v1",    // version
//	    "",      // namespace (empty for cluster-scoped)
//	    "nodes", // resource
//	    "worker-node-1",
//	)
//
// Example - Get Custom Resource:
//
//	cr, err := client.Get[MyCustomResource](
//	    restClient,
//	    "example.com",     // group
//	    "v1alpha1",        // version
//	    "default",         // namespace
//	    "mycustomresources", // resource (plural)
//	    "my-instance",     // name
//	)
func Get[T any](restClient rest.Interface, group, version, namespace, resource, name string) (T, error) {
	var t T

	if body, err := do(restClient, http.MethodGet, group, version, namespace, resource, name, nil); err != nil {
		return t, err
	} else if err := json.Unmarshal([]byte(body), &t); err != nil {
		return t, err
	} else {
		return t, nil
	}
}

// Post creates a new Kubernetes resource.
//
// This function sends a POST request to the Kubernetes API server to create
// a new resource. The resource's API group, version, and namespace are extracted
// from the object's metadata.
//
// Parameters:
//   - restClient: Kubernetes REST client interface
//   - resource: Resource type in plural form (e.g., "pods", "configmaps", "deployments")
//   - object: Kubernetes resource object to create (must implement runtime.Object)
//
// Returns:
//   - error: Error if creation fails
//
// The object parameter must:
//   - Implement runtime.Object interface
//   - Have TypeMeta set (APIVersion and Kind)
//   - Have ObjectMeta with at least Name and optionally Namespace
//   - Not have ResourceVersion set (will be assigned by API server)
//
// Example - Create ConfigMap:
//
//	import coreV1 "k8s.io/api/core/v1"
//	import metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//
//	configMap := &coreV1.ConfigMap{
//	    TypeMeta: metaV1.TypeMeta{
//	        APIVersion: "v1",
//	        Kind:       "ConfigMap",
//	    },
//	    ObjectMeta: metaV1.ObjectMeta{
//	        Name:      "app-config",
//	        Namespace: "default",
//	    },
//	    Data: map[string]string{
//	        "config.yaml": "key: value",
//	    },
//	}
//
//	err := client.Post(restClient, "configmaps", configMap)
//	if err != nil {
//	    log.Fatal("Failed to create ConfigMap:", err)
//	}
//
// Example - Create Deployment:
//
//	import appsV1 "k8s.io/api/apps/v1"
//
//	replicas := int32(3)
//	deployment := &appsV1.Deployment{
//	    TypeMeta: metaV1.TypeMeta{
//	        APIVersion: "apps/v1",
//	        Kind:       "Deployment",
//	    },
//	    ObjectMeta: metaV1.ObjectMeta{
//	        Name:      "web-app",
//	        Namespace: "production",
//	    },
//	    Spec: appsV1.DeploymentSpec{
//	        Replicas: &replicas,
//	        Selector: &metaV1.LabelSelector{
//	            MatchLabels: map[string]string{"app": "web"},
//	        },
//	        Template: coreV1.PodTemplateSpec{
//	            ObjectMeta: metaV1.ObjectMeta{
//	                Labels: map[string]string{"app": "web"},
//	            },
//	            Spec: coreV1.PodSpec{
//	                Containers: []coreV1.Container{{
//	                    Name:  "web",
//	                    Image: "nginx:1.21",
//	                }},
//	            },
//	        },
//	    },
//	}
//
//	err := client.Post(restClient, "deployments", deployment)
//
// Example - Create Custom Resource:
//
//	customResource := &MyCustomResource{
//	    TypeMeta: metaV1.TypeMeta{
//	        APIVersion: "example.com/v1alpha1",
//	        Kind:       "MyCustomResource",
//	    },
//	    ObjectMeta: metaV1.ObjectMeta{
//	        Name:      "my-instance",
//	        Namespace: "default",
//	    },
//	    Spec: MySpec{
//	        Field1: "value1",
//	    },
//	}
//
//	err := client.Post(restClient, "mycustomresources", customResource)
func Post(restClient rest.Interface, resource string, object runtime.Object) error {
	return do_for_object(restClient, http.MethodPost, resource, object)
}

// Put updates an existing Kubernetes resource.
//
// This function sends a PUT request to the Kubernetes API server to update
// an existing resource. The resource must already exist and must include
// a valid ResourceVersion for optimistic concurrency control.
//
// Parameters:
//   - restClient: Kubernetes REST client interface
//   - resource: Resource type in plural form (e.g., "pods", "configmaps", "deployments")
//   - object: Kubernetes resource object with updates (must implement runtime.Object)
//
// Returns:
//   - error: Error if update fails
//
// The object parameter must:
//   - Implement runtime.Object interface
//   - Have TypeMeta set (APIVersion and Kind)
//   - Have ObjectMeta with Name, Namespace, and ResourceVersion
//   - Match an existing resource on the server
//
// Important: ResourceVersion must match the current version on the server,
// otherwise the update will fail with a conflict error. To ensure this,
// first Get the resource, modify it, then Put it back.
//
// Example - Update ConfigMap:
//
//	import coreV1 "k8s.io/api/core/v1"
//
//	// First, get the current ConfigMap
//	configMap, err := client.Get[coreV1.ConfigMap](
//	    restClient, "", "v1", "default", "configmaps", "app-config",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Modify the ConfigMap
//	configMap.Data["new-key"] = "new-value"
//
//	// Update the resource
//	err = client.Put(restClient, "configmaps", &configMap)
//	if err != nil {
//	    log.Fatal("Failed to update ConfigMap:", err)
//	}
//
// Example - Update Deployment replicas:
//
//	import appsV1 "k8s.io/api/apps/v1"
//
//	// Get current deployment
//	deployment, err := client.Get[appsV1.Deployment](
//	    restClient, "apps", "v1", "production", "deployments", "web-app",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Scale up
//	replicas := int32(5)
//	deployment.Spec.Replicas = &replicas
//
//	// Update
//	err = client.Put(restClient, "deployments", &deployment)
//	if err != nil {
//	    log.Fatal("Failed to scale deployment:", err)
//	}
//
// Example - Update with retry on conflict:
//
//	import "k8s.io/apimachinery/pkg/api/errors"
//
//	for i := 0; i < 3; i++ {
//	    // Get latest version
//	    configMap, err := client.Get[coreV1.ConfigMap](
//	        restClient, "", "v1", "default", "configmaps", "app-config",
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Modify
//	    configMap.Data["key"] = "value"
//
//	    // Try to update
//	    err = client.Put(restClient, "configmaps", &configMap)
//	    if err == nil {
//	        break // Success
//	    }
//
//	    // Retry if conflict
//	    if !errors.IsConflict(err) {
//	        log.Fatal(err)
//	    }
//	}
func Put(restClient rest.Interface, resource string, object runtime.Object) error {
	return do_for_object(restClient, http.MethodPut, resource, object)
}

// Delete removes a Kubernetes resource by name.
//
// This function sends a DELETE request to the Kubernetes API server to remove
// the specified resource. The deletion may be immediate or graceful depending
// on the resource type and cluster configuration.
//
// Parameters:
//   - restClient: Kubernetes REST client interface
//   - group: API group (empty "" for core resources)
//   - version: API version (e.g., "v1", "v1beta1")
//   - namespace: Namespace name (empty "" for cluster-scoped resources)
//   - resource: Resource type in plural form (e.g., "pods", "configmaps")
//   - name: Resource name to delete
//
// Returns:
//   - error: Error if deletion fails
//
// Deletion behavior:
//   - Pods: Graceful shutdown with termination period
//   - Deployments: Cascading deletion of owned ReplicaSets and Pods
//   - ConfigMaps/Secrets: Immediate deletion
//   - Custom Resources: Depends on finalizers
//
// Example - Delete ConfigMap:
//
//	err := client.Delete(
//	    restClient,
//	    "",           // group
//	    "v1",         // version
//	    "default",    // namespace
//	    "configmaps", // resource
//	    "app-config", // name
//	)
//	if err != nil {
//	    log.Fatal("Failed to delete ConfigMap:", err)
//	}
//
// Example - Delete Deployment:
//
//	err := client.Delete(
//	    restClient,
//	    "apps",        // group
//	    "v1",          // version
//	    "production",  // namespace
//	    "deployments", // resource
//	    "web-app",     // name
//	)
//	if err != nil {
//	    log.Fatal("Failed to delete Deployment:", err)
//	}
//
// Example - Delete Pod:
//
//	err := client.Delete(
//	    restClient,
//	    "",        // group
//	    "v1",      // version
//	    "default", // namespace
//	    "pods",    // resource
//	    "web-app-pod-12345",
//	)
//
// Example - Delete Node (cluster-scoped):
//
//	err := client.Delete(
//	    restClient,
//	    "",      // group
//	    "v1",    // version
//	    "",      // namespace (empty for cluster-scoped)
//	    "nodes", // resource
//	    "worker-node-1",
//	)
//
// Example - Delete Custom Resource:
//
//	err := client.Delete(
//	    restClient,
//	    "example.com",        // group
//	    "v1alpha1",           // version
//	    "default",            // namespace
//	    "mycustomresources",  // resource
//	    "my-instance",        // name
//	)
//
// Example - Check if resource exists before deletion:
//
//	import "k8s.io/apimachinery/pkg/api/errors"
//
//	err := client.Delete(restClient, "", "v1", "default", "configmaps", "app-config")
//	if err != nil {
//	    if errors.IsNotFound(err) {
//	        log.Println("ConfigMap not found, already deleted")
//	    } else {
//	        log.Fatal("Failed to delete:", err)
//	    }
//	}
func Delete(restClient rest.Interface, group, version, namespace, resource, name string) error {
	_, err := do(restClient, http.MethodDelete, group, version, namespace, resource, name, nil)
	return err
}

func do_for_object(restClient rest.Interface, method, resource string, object runtime.Object) error {
	if metaObject, err := meta.Accessor(object); err != nil {
		return err
	} else {
		group := object.GetObjectKind().GroupVersionKind().Group
		version := object.GetObjectKind().GroupVersionKind().Version
		namespace := metaObject.GetNamespace()

		name := ""
		switch method {
		case http.MethodPost:
			name = ""
		default:
			name = metaObject.GetName()
		}

		_, err := do(restClient, method, group, version, namespace, resource, name, object)
		return err
	}
}

func do(restClient rest.Interface, method, group, version, namespace, resource, name string, object runtime.Object) (string, error) {
	if restClient == nil {
		return "", errors.New("restClient is nil")
	}

	var request *rest.Request
	switch method {
	case http.MethodGet:
		request = restClient.Get()
	case http.MethodPost:
		request = restClient.Post()
	case http.MethodPut:
		request = restClient.Put()
	case http.MethodDelete:
		request = restClient.Delete()
	}

	absPath := "api"
	if len(group) != 0 {
		absPath += "s/" + group
	}
	absPath += "/" + version

	request.AbsPath(absPath)

	if len(namespace) != 0 {
		request.Namespace(namespace)
	}

	if len(resource) != 0 {
		request.Resource(resource)
	}

	if len(name) != 0 {
		request.Name(name)
	}

	if object != nil {
		if body, err := json.Marshal(object); err != nil {
			return "", err
		} else {
			request.Body(body)
		}
	}

	result := request.Do(context.TODO())

	body, _ := result.Raw()

	return string(body), result.Error()
}
