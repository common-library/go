# Kubernetes

Utilities for Kubernetes resource management with type-safe REST client operations.

## Overview

The kubernetes package provides simplified interfaces for interacting with Kubernetes API servers through REST clients. It supports both in-cluster and external client configurations with generic type-safe operations for CRUD (Create, Read, Update, Delete) operations on Kubernetes resources.

## Features

- **In-Cluster Client** - Automatic configuration for pods running in Kubernetes
- **External Client** - Configuration support for external applications
- **Generic Get** - Type-safe resource retrieval
- **Create Operations** - Post resources to Kubernetes API
- **Update Operations** - Put updated resources
- **Delete Operations** - Remove resources by name
- **Custom Resources** - Full support for CRDs
- **Standard Resources** - Support for all Kubernetes native resources

## Installation

```bash
go get -u github.com/common-library/go/kubernetes/resource/client
go get -u k8s.io/client-go
go get -u k8s.io/api
```

## Quick Start

### In-Cluster Usage

```go
import (
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
)

// Inside a Kubernetes pod
restClient, err := client.GetClientForInCluster()
if err != nil {
    log.Fatal(err)
}

// Get a ConfigMap
configMap, err := client.Get[coreV1.ConfigMap](
    restClient,
    "",           // group (empty for core resources)
    "v1",         // version
    "default",    // namespace
    "configmaps", // resource type
    "app-config", // name
)
```

### External Usage

```go
import (
    "k8s.io/client-go/tools/clientcmd"
    "github.com/common-library/go/kubernetes/resource/client"
)

// Load kubeconfig
config, err := clientcmd.BuildConfigFromFlags("", "/path/to/kubeconfig")
if err != nil {
    log.Fatal(err)
}

// Create client
restClient, err := client.GetClientUsingConfig(config)
if err != nil {
    log.Fatal(err)
}

// Use client...
```

## API Reference

### GetClientForInCluster

```go
func GetClientForInCluster() (rest.Interface, error)
```

Creates a Kubernetes REST client for in-cluster use from within a pod.

**Returns:**
- `rest.Interface` - Kubernetes REST client
- `error` - Error if client creation fails

**Requirements:**
- Must run inside a Kubernetes pod
- Service account token must be mounted
- CA certificate must be available

### GetClientUsingConfig

```go
func GetClientUsingConfig(config *rest.Config) (rest.Interface, error)
```

Creates a Kubernetes REST client using a provided configuration.

**Parameters:**
- `config` - Kubernetes REST client configuration

**Returns:**
- `rest.Interface` - Kubernetes REST client
- `error` - Error if client creation fails

### Get

```go
func Get[T any](
    restClient rest.Interface,
    group, version, namespace, resource, name string,
) (T, error)
```

Retrieves a Kubernetes resource by name.

**Type Parameters:**
- `T` - Target Go type for the resource

**Parameters:**
- `restClient` - Kubernetes REST client
- `group` - API group (empty for core resources)
- `version` - API version (e.g., "v1")
- `namespace` - Namespace (empty for cluster-scoped)
- `resource` - Resource type in plural (e.g., "pods")
- `name` - Resource name

**Returns:**
- `T` - Retrieved resource
- `error` - Error if retrieval fails

### Post

```go
func Post(
    restClient rest.Interface,
    resource string,
    object runtime.Object,
) error
```

Creates a new Kubernetes resource.

**Parameters:**
- `restClient` - Kubernetes REST client
- `resource` - Resource type in plural
- `object` - Resource object to create

**Returns:**
- `error` - Error if creation fails

### Put

```go
func Put(
    restClient rest.Interface,
    resource string,
    object runtime.Object,
) error
```

Updates an existing Kubernetes resource.

**Parameters:**
- `restClient` - Kubernetes REST client
- `resource` - Resource type in plural
- `object` - Resource object with updates (must include ResourceVersion)

**Returns:**
- `error` - Error if update fails

### Delete

```go
func Delete(
    restClient rest.Interface,
    group, version, namespace, resource, name string,
) error
```

Deletes a Kubernetes resource by name.

**Parameters:**
- `restClient` - Kubernetes REST client
- `group` - API group (empty for core resources)
- `version` - API version
- `namespace` - Namespace (empty for cluster-scoped)
- `resource` - Resource type in plural
- `name` - Resource name to delete

**Returns:**
- `error` - Error if deletion fails

## Complete Examples

### ConfigMap Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    // Get client
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create ConfigMap
    configMap := &coreV1.ConfigMap{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "v1",
            Kind:       "ConfigMap",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "app-config",
            Namespace: "default",
        },
        Data: map[string]string{
            "database.host": "localhost",
            "database.port": "5432",
            "app.mode":      "production",
        },
    }
    
    err = client.Post(restClient, "configmaps", configMap)
    if err != nil {
        log.Fatal("Failed to create ConfigMap:", err)
    }
    fmt.Println("ConfigMap created successfully")
    
    // Get ConfigMap
    retrieved, err := client.Get[coreV1.ConfigMap](
        restClient, "", "v1", "default", "configmaps", "app-config",
    )
    if err != nil {
        log.Fatal("Failed to get ConfigMap:", err)
    }
    fmt.Printf("ConfigMap data: %v\n", retrieved.Data)
    
    // Update ConfigMap
    retrieved.Data["app.mode"] = "staging"
    err = client.Put(restClient, "configmaps", &retrieved)
    if err != nil {
        log.Fatal("Failed to update ConfigMap:", err)
    }
    fmt.Println("ConfigMap updated successfully")
    
    // Delete ConfigMap
    err = client.Delete(restClient, "", "v1", "default", "configmaps", "app-config")
    if err != nil {
        log.Fatal("Failed to delete ConfigMap:", err)
    }
    fmt.Println("ConfigMap deleted successfully")
}
```

### Deployment Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    appsV1 "k8s.io/api/apps/v1"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Deployment
    replicas := int32(3)
    deployment := &appsV1.Deployment{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "apps/v1",
            Kind:       "Deployment",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "nginx-deployment",
            Namespace: "default",
            Labels: map[string]string{
                "app": "nginx",
            },
        },
        Spec: appsV1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metaV1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": "nginx",
                },
            },
            Template: coreV1.PodTemplateSpec{
                ObjectMeta: metaV1.ObjectMeta{
                    Labels: map[string]string{
                        "app": "nginx",
                    },
                },
                Spec: coreV1.PodSpec{
                    Containers: []coreV1.Container{{
                        Name:  "nginx",
                        Image: "nginx:1.21",
                        Ports: []coreV1.ContainerPort{{
                            ContainerPort: 80,
                        }},
                    }},
                },
            },
        },
    }
    
    err = client.Post(restClient, "deployments", deployment)
    if err != nil {
        log.Fatal("Failed to create Deployment:", err)
    }
    fmt.Println("Deployment created")
    
    // Scale Deployment
    retrieved, err := client.Get[appsV1.Deployment](
        restClient, "apps", "v1", "default", "deployments", "nginx-deployment",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    newReplicas := int32(5)
    retrieved.Spec.Replicas = &newReplicas
    
    err = client.Put(restClient, "deployments", &retrieved)
    if err != nil {
        log.Fatal("Failed to scale Deployment:", err)
    }
    fmt.Println("Deployment scaled to 5 replicas")
}
```

### Secret Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Secret
    secret := &coreV1.Secret{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "v1",
            Kind:       "Secret",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "db-credentials",
            Namespace: "default",
        },
        Type: coreV1.SecretTypeOpaque,
        StringData: map[string]string{
            "username": "admin",
            "password": "super-secret-password",
        },
    }
    
    err = client.Post(restClient, "secrets", secret)
    if err != nil {
        log.Fatal("Failed to create Secret:", err)
    }
    fmt.Println("Secret created")
    
    // Get Secret
    retrieved, err := client.Get[coreV1.Secret](
        restClient, "", "v1", "default", "secrets", "db-credentials",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Secret has %d keys\n", len(retrieved.Data))
}
```

### Service Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Service
    service := &coreV1.Service{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "v1",
            Kind:       "Service",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "nginx-service",
            Namespace: "default",
        },
        Spec: coreV1.ServiceSpec{
            Selector: map[string]string{
                "app": "nginx",
            },
            Type: coreV1.ServiceTypeLoadBalancer,
            Ports: []coreV1.ServicePort{{
                Port:       80,
                TargetPort: intstr.FromInt(80),
                Protocol:   coreV1.ProtocolTCP,
            }},
        },
    }
    
    err = client.Post(restClient, "services", service)
    if err != nil {
        log.Fatal("Failed to create Service:", err)
    }
    fmt.Println("Service created")
}
```

### Namespace Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Namespace
    namespace := &coreV1.Namespace{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "v1",
            Kind:       "Namespace",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name: "production",
            Labels: map[string]string{
                "environment": "production",
            },
        },
    }
    
    err = client.Post(restClient, "namespaces", namespace)
    if err != nil {
        log.Fatal("Failed to create Namespace:", err)
    }
    fmt.Println("Namespace created")
    
    // Get Namespace (cluster-scoped resource)
    retrieved, err := client.Get[coreV1.Namespace](
        restClient,
        "",            // group
        "v1",          // version
        "",            // namespace (empty for cluster-scoped)
        "namespaces",  // resource
        "production",  // name
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Namespace: %s\n", retrieved.Name)
    fmt.Printf("Labels: %v\n", retrieved.Labels)
}
```

### Pod Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Pod
    pod := &coreV1.Pod{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "v1",
            Kind:       "Pod",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "test-pod",
            Namespace: "default",
        },
        Spec: coreV1.PodSpec{
            Containers: []coreV1.Container{{
                Name:  "nginx",
                Image: "nginx:1.21",
                Ports: []coreV1.ContainerPort{{
                    ContainerPort: 80,
                }},
            }},
        },
    }
    
    err = client.Post(restClient, "pods", pod)
    if err != nil {
        log.Fatal("Failed to create Pod:", err)
    }
    fmt.Println("Pod created")
    
    // Get Pod status
    retrieved, err := client.Get[coreV1.Pod](
        restClient, "", "v1", "default", "pods", "test-pod",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Pod phase: %s\n", retrieved.Status.Phase)
    fmt.Printf("Pod IP: %s\n", retrieved.Status.PodIP)
}
```

### Custom Resource Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/kubernetes/resource/client"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Custom Resource Definition
type MyCustomResource struct {
    metaV1.TypeMeta   `json:",inline"`
    metaV1.ObjectMeta `json:"metadata,omitempty"`
    Spec              MySpec   `json:"spec,omitempty"`
    Status            MyStatus `json:"status,omitempty"`
}

type MySpec struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

type MyStatus struct {
    Ready bool `json:"ready"`
}

func main() {
    restClient, err := client.GetClientForInCluster()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Custom Resource
    cr := &MyCustomResource{
        TypeMeta: metaV1.TypeMeta{
            APIVersion: "example.com/v1alpha1",
            Kind:       "MyCustomResource",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:      "my-instance",
            Namespace: "default",
        },
        Spec: MySpec{
            Field1: "value1",
            Field2: 42,
        },
    }
    
    err = client.Post(restClient, "mycustomresources", cr)
    if err != nil {
        log.Fatal("Failed to create Custom Resource:", err)
    }
    fmt.Println("Custom Resource created")
    
    // Get Custom Resource
    retrieved, err := client.Get[MyCustomResource](
        restClient,
        "example.com",        // group
        "v1alpha1",           // version
        "default",            // namespace
        "mycustomresources",  // resource (plural)
        "my-instance",        // name
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("CR Spec: %+v\n", retrieved.Spec)
}
```

### External Client with Kubeconfig

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/common-library/go/kubernetes/resource/client"
    coreV1 "k8s.io/api/core/v1"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Get kubeconfig path
    home, _ := os.UserHomeDir()
    kubeconfig := filepath.Join(home, ".kube", "config")
    
    // Load config
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatal("Failed to load kubeconfig:", err)
    }
    
    // Create client
    restClient, err := client.GetClientUsingConfig(config)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    
    // List pods (using empty name to get list)
    // Note: This example shows the pattern, but you'd typically use
    // the List operation from client-go for listing
    
    // Get specific pod
    pod, err := client.Get[coreV1.Pod](
        restClient, "", "v1", "default", "pods", "my-pod",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Pod: %s, Status: %s\n", pod.Name, pod.Status.Phase)
}
```

## Best Practices

### 1. Handle ResourceVersion for Updates

```go
// Good: Get, modify, then put
configMap, _ := client.Get[coreV1.ConfigMap](
    restClient, "", "v1", "default", "configmaps", "app-config",
)
configMap.Data["key"] = "value"
client.Put(restClient, "configmaps", &configMap)

// Avoid: Creating object without ResourceVersion
// Will fail with conflict error
```

### 2. Use Proper API Groups

```go
// Good: Core resources use empty group
client.Get[coreV1.Pod](restClient, "", "v1", "default", "pods", "my-pod")

// Good: Apps resources use "apps" group
client.Get[appsV1.Deployment](restClient, "apps", "v1", "default", "deployments", "my-app")

// Avoid: Wrong group
// client.Get[coreV1.Pod](restClient, "core", "v1", ...) // Wrong
```

### 3. Check Errors Properly

```go
import "k8s.io/apimachinery/pkg/api/errors"

// Good: Check specific error types
pod, err := client.Get[coreV1.Pod](restClient, "", "v1", "default", "pods", "my-pod")
if err != nil {
    if errors.IsNotFound(err) {
        log.Println("Pod not found")
    } else if errors.IsUnauthorized(err) {
        log.Println("Unauthorized")
    } else {
        log.Fatal(err)
    }
}

// Avoid: Ignoring errors
// pod, _ := client.Get[...]
```

### 4. Set Required Metadata

```go
// Good: Complete metadata
deployment := &appsV1.Deployment{
    TypeMeta: metaV1.TypeMeta{
        APIVersion: "apps/v1",
        Kind:       "Deployment",
    },
    ObjectMeta: metaV1.ObjectMeta{
        Name:      "my-app",
        Namespace: "default",
    },
    Spec: appsV1.DeploymentSpec{...},
}

// Avoid: Missing TypeMeta or ObjectMeta
```

### 5. Use Retry Logic for Conflicts

```go
// Good: Retry on conflict
for i := 0; i < 3; i++ {
    configMap, err := client.Get[coreV1.ConfigMap](...)
    if err != nil {
        log.Fatal(err)
    }
    
    configMap.Data["key"] = "value"
    
    err = client.Put(restClient, "configmaps", &configMap)
    if err == nil {
        break
    }
    
    if !errors.IsConflict(err) {
        log.Fatal(err)
    }
}
```

## Error Handling

### Common Errors

```go
import "k8s.io/apimachinery/pkg/api/errors"

pod, err := client.Get[coreV1.Pod](restClient, "", "v1", "default", "pods", "my-pod")
if err != nil {
    switch {
    case errors.IsNotFound(err):
        // Resource doesn't exist (404)
        log.Println("Pod not found")
        
    case errors.IsConflict(err):
        // Update conflict (409) - ResourceVersion mismatch
        log.Println("Conflict, retry with latest version")
        
    case errors.IsUnauthorized(err):
        // Authentication failed (401)
        log.Println("Unauthorized")
        
    case errors.IsForbidden(err):
        // Permission denied (403)
        log.Println("Forbidden")
        
    case errors.IsAlreadyExists(err):
        // Resource already exists (409)
        log.Println("Already exists")
        
    case errors.IsInvalid(err):
        // Invalid resource (422)
        log.Println("Invalid resource")
        
    default:
        log.Fatal("Unknown error:", err)
    }
}
```

## Dependencies

- `k8s.io/client-go` - Kubernetes client library
- `k8s.io/api` - Kubernetes API types
- `k8s.io/apimachinery` - Kubernetes API machinery

## Further Reading

- [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)
- [client-go Documentation](https://github.com/kubernetes/client-go)
- [Kubernetes API Reference](https://kubernetes.io/docs/reference/kubernetes-api/)
- [Custom Resource Definitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)
