# Kubernetes

Kubernetes client utilities and resource management.

## Overview

The kubernetes package provides simplified interfaces for interacting with Kubernetes clusters, including resource management, client creation, and common operations.

## Subpackages

### resource/client

Kubernetes client for resource operations (create, read, update, delete).

[📖 Documentation](resource/client/README.md)

**Features:**
- RESTful Kubernetes client operations
- Dynamic resource management
- Support for all Kubernetes API versions
- Simplified CRUD operations
- Namespace-aware operations

**Quick Example:**
```go
import "github.com/common-library/go/kubernetes/resource/client"

client := &client.Client{}
err := client.CreateClient("/path/to/kubeconfig")

// Create resource
err = client.Create("apps/v1", "deployments", "default", deploymentYaml)

// Get resource
data, err := client.Read("v1", "pods", "default", "mypod")

// Update resource
err = client.Update("apps/v1", "deployments", "default", "myapp", updatedYaml)

// Delete resource
err = client.Delete("v1", "services", "default", "myservice")
```

**Supported Operations:**
- `CreateClient` - Initialize client with kubeconfig
- `Create` - Create new resources
- `Read` - Get resource data
- `Update` - Update existing resources
- `Delete` - Remove resources
- `List` - List resources by label selector

## Use Cases

- **Custom Controllers** - Build Kubernetes operators and controllers
- **CI/CD Tools** - Automate deployment and resource management
- **Cluster Management** - Programmatic cluster administration
- **Testing** - Integration tests with Kubernetes resources
- **Resource Validation** - Validate and manipulate Kubernetes manifests

## Installation

```bash
go get -u github.com/common-library/go/kubernetes/resource/client
```

## Best Practices

1. **Use Kubeconfig** - Store credentials securely in kubeconfig files
2. **Namespace Isolation** - Always specify namespaces for resources
3. **Error Handling** - Check errors for all Kubernetes operations
4. **Resource Lifecycle** - Clean up resources in defer statements
5. **API Versioning** - Use stable API versions (v1, apps/v1) in production

## Common Patterns

### Creating Resources

```go
client := &client.Client{}
client.CreateClient("/home/user/.kube/config")

// Create deployment
deploymentYAML := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
`

err := client.Create("apps/v1", "deployments", "default", deploymentYAML)
```

### Reading Resources

```go
// Get pod details
podData, err := client.Read("v1", "pods", "default", "nginx-pod")
if err != nil {
    log.Fatal(err)
}
fmt.Println(podData)
```

### Listing Resources

```go
// List all pods with label app=nginx
pods, err := client.List("v1", "pods", "default", "app=nginx")
```

## Dependencies

- `k8s.io/client-go` - Official Kubernetes Go client
- `k8s.io/apimachinery` - Kubernetes API machinery

## Further Reading

- [resource/client Package Documentation](resource/client/README.md)
- [Kubernetes API Documentation](https://kubernetes.io/docs/reference/kubernetes-api/)
- [client-go Documentation](https://github.com/kubernetes/client-go)
