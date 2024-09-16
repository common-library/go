package v1alpha1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/common-library/go/kubernetes/resource/client"
	crd "github.com/common-library/go/kubernetes/resource/custom-resource-definition/project01/group01/v1alpha1"
	"github.com/common-library/go/kubernetes/resource/custom-resource/project01/group01/v1alpha1"
	"github.com/emicklei/go-restful/v3"
	appsV1 "k8s.io/api/apps/v1"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

func TestResource01(t *testing.T) {
	resource01 := v1alpha1.Resource01{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: crd.Group + "/" + crd.Version,
			Kind:       "Resource01",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-01",
			Namespace: "default",
		},
		Spec: v1alpha1.Resource01Spec{
			Field01: 1,
			Field02: "value",
		},
	}

	object := resource01.DeepCopyObject().(*v1alpha1.Resource01)
	if object.Spec != resource01.Spec {
		t.Log(object)
		t.Fatal("invalid")
	}

	serverHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "{}")
		}
	}

	testFunc := func(restClient rest.Interface) {
		defer func() {
			if err := client.Delete(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", crd.Resource01.Name); err != nil {
				t.Fatal(err)
			}
		}()

		defer func() {
			if err := client.Delete(restClient, crd.Group, crd.Version, resource01.Namespace, "resource01s", resource01.Name); err != nil {
				t.Fatal(err)
			}
		}()

		if err := client.Post(restClient, "customresourcedefinitions", &crd.Resource01); err != nil {
			t.Fatal(err)
		}

		if err := client.Post(restClient, "resource01s", &resource01); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Get[apiextensionsV1.CustomResourceDefinition](restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", crd.Resource01.Name); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Get[v1alpha1.Resource01](restClient, crd.Group, crd.Version, resource01.Namespace, "resource01s", resource01.Name); err != nil {
			t.Fatal(err)
		}
	}

	test(t, serverHandlerFunc, testFunc)
}

func TestResource01List(t *testing.T) {
	serverHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{}")
	}

	testFunc := func(restClient rest.Interface) {
		if resource01List, err := client.Get[v1alpha1.Resource01List](restClient, crd.Group, crd.Version, "", "resource01s", ""); err != nil {
			t.Fatal(err)
		} else {
			for _, resource01 := range resource01List.Items {
				t.Log(resource01.Name)
			}
		}
	}

	test(t, serverHandlerFunc, testFunc)
}

func test(t *testing.T, serverHandlerFunc http.HandlerFunc, testFunc func(rest.Interface)) {
	t.Parallel()

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
