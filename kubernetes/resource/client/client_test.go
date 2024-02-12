package client_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/kubernetes/resource/client"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func TestGetClientForInCluster(t *testing.T) {
	if _, err := client.GetClientForInCluster(); err != nil && err.Error() != "unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined" {
		t.Fatal(err)
	}
}

func TestGetClientUsingConfig(t *testing.T) {
	config := &rest.Config{}

	if _, err := client.GetClientUsingConfig(config); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	return

	restClient, err := client.GetClientForInCluster()
	if err != nil {
		t.Fatal(err)
	}

	name := uuid.New().String()
	namespace := "default"
	key := "key"
	value := "value"
	configMap := coreV1.ConfigMap{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	if err := client.Post(restClient, "", "v1", namespace, "configmaps", &configMap); err != nil {
		t.Fatal(err)
	} else if configMap, err := client.Get[coreV1.ConfigMap](restClient, "", "v1", namespace, "configmaps", name); err != nil {
		t.Fatal(err)
	} else if configMap.Name != name || configMap.Data[key] != value {
		t.Fatal("invalid name - ", configMap.Name, configMap.Data)
	}

	value2 := "value2"
	configMap.Data[key] = value2
	if err := client.Put(restClient, "", "v1", namespace, "configmaps", name, &configMap); err != nil {
		t.Fatal(err)
	} else if configMap, err := client.Get[coreV1.ConfigMap](restClient, "", "v1", namespace, "configmaps", name); err != nil {
		t.Fatal(err)
	} else if configMap.Name != name || configMap.Data[key] != value2 {
		t.Fatal("invalid name - ", configMap.Name, configMap.Data)
	}

	if list, err := client.Get[coreV1.ConfigMapList](restClient, "", "v1", namespace, "configmaps", ""); err != nil {
		t.Fatal(err)
	} else {
		for _, item := range list.Items {
			t.Log(item.Name)
		}
	}

	if err := client.Delete(restClient, "", "v1", namespace, "configmaps", name); err != nil {
		t.Fatal(err)
	}
}

func TestPost(t *testing.T) {
	TestGet(t)
}

func TestPut(t *testing.T) {
	TestGet(t)
}

func TestDelete(t *testing.T) {
	TestGet(t)
}
