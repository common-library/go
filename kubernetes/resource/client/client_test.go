package client_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/common-library/go/json"
	"github.com/common-library/go/kubernetes/resource/client"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
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
			key: value,
		},
	}

	serverHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if body, err := json.ToString(configMap); err != nil {
				t.Fatal(err)
			} else {
				fmt.Fprintf(w, body)
			}
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			if body, err := ioutil.ReadAll(r.Body); err != nil {
				t.Fatal(err)
			} else if answer, err := json.ToString(configMap); err != nil {
				t.Fatal(err)
			} else if string(body) != answer {
				t.Log(string(body))
				t.Log(answer)
				t.Fatal("invalid")
			}
		case http.MethodDelete:
		}
	}

	testFunc := func(restClient rest.Interface) {
		defer func() {
			if err := client.Delete(restClient, "", "v1", namespace, "configmaps", name); err != nil {
				t.Fatal(err)
			}
		}()

		if err := client.Post(restClient, "configmaps", &configMap); err != nil {
			t.Fatal(err)
		}

		if err := client.Put(restClient, "configmaps", &configMap); err != nil {
			t.Fatal(err)
		}

		if configMap, err := client.Get[coreV1.ConfigMap](restClient, "", "v1", namespace, "configmaps", name); err != nil {
			t.Fatal(err)
		} else if configMap.Name != name || configMap.Data[key] != value {
			t.Fatal(configMap)
		}
	}

	test(t, serverHandlerFunc, testFunc)
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

func test(t *testing.T, serverHandlerFunc http.HandlerFunc, testFunc func(rest.Interface)) {
	server := httptest.NewUnstartedServer(serverHandlerFunc)
	server.EnableHTTP2 = false
	server.StartTLS()
	defer server.Close()

	config := rest.ClientContentConfig{
		ContentType:  restful.MIME_JSON,
		GroupVersion: appsV1.SchemeGroupVersion,
		Negotiator:   runtime.NewClientNegotiator(scheme.Codecs.WithoutConversion(), appsV1.SchemeGroupVersion),
	}

	if baseURL, err := url.Parse(server.URL); err != nil {
		t.Fatal(err)
	} else if restClient, err := rest.NewRESTClient(baseURL, "", config, nil, server.Client()); err != nil {
		t.Fatal(err)
	} else {
		testFunc(restClient)
	}
}
