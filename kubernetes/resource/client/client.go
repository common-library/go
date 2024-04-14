// Package client provides kubernetes client implementations.
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

// GetClientForInCluster is get a client that runs within the pod.
//
// ex) restClient, err := client.GetClientForInCluster()
func GetClientForInCluster() (rest.Interface, error) {
	if config, err := rest.InClusterConfig(); err != nil {
		return nil, err
	} else if clientset, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return clientset.RESTClient(), nil
	}
}

// GetClientUsingConfig is get the client running outside the pod.
//
// ex) restClient, err := client.GetClientUsingConfig(config)
func GetClientUsingConfig(config *rest.Config) (rest.Interface, error) {
	if clientset, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return clientset.RESTClient(), nil
	}
}

// Get is get kubernetes resource.
//
// ex) configMap, err := client.Get[coreV1.ConfigMap](restClient, "", "v1", namespace, "configmaps", name)
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

// Post is to create a kubernetes resource.
//
// ex)
//
//	configMap := coreV1.ConfigMap{...}
//	err := client.Post(restClient, "configmaps", &configMap)
func Post(restClient rest.Interface, resource string, object runtime.Object) error {
	return do_for_object(restClient, http.MethodPost, resource, object)
}

// Put is to update a kubernetes resource.
//
// ex)
//
//	configMap := coreV1.ConfigMap{...}
//	err := client.Put(restClient, "configmaps", &configMap)
func Put(restClient rest.Interface, resource string, object runtime.Object) error {
	return do_for_object(restClient, http.MethodPut, resource, object)
}

// Delete is to delete a kubernetes resource.
//
// ex) err := client.Delete(restClient, "", "v1", namespace, "configmaps", name)
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
