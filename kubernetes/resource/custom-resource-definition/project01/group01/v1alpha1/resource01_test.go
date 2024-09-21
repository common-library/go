package v1alpha1_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/common-library/go/json"
	"github.com/common-library/go/kubernetes/resource/client"
	"github.com/common-library/go/kubernetes/resource/custom-resource-definition/project01/group01/v1alpha1"
	"github.com/emicklei/go-restful/v3"
	appsV1 "k8s.io/api/apps/v1"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

func TestResource01(t *testing.T) {
	serverHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if body, err := json.ToString(v1alpha1.Resource01); err != nil {
				t.Fatal(err)
			} else {
				fmt.Fprintf(w, body)
			}

		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			if body, err := ioutil.ReadAll(r.Body); err != nil {
				t.Fatal(err)
			} else if answer, err := json.ToString(v1alpha1.Resource01); err != nil {
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
		if resource01, err := client.Get[apiextensionsV1.CustomResourceDefinition](restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", v1alpha1.Resource01.Name); err != nil {
			t.Fatal(err)
		} else if compare, err := json.ToString(resource01); err != nil {
			t.Fatal(err)
		} else if answer, err := json.ToString(v1alpha1.Resource01); err != nil {
			t.Fatal(err)
		} else if string(compare) != answer {
			t.Log(string(compare))
			t.Log(answer)
			t.Fatal("invalid")
		}

		if err := client.Post(restClient, "customresourcedefinitions", &v1alpha1.Resource01); err != nil {
			t.Fatal(err)
		}

		if err := client.Put(restClient, "customresourcedefinitions", &v1alpha1.Resource01); err != nil {
			t.Fatal(err)
		}

		if err := client.Delete(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", v1alpha1.Resource01.Name); err != nil {
			t.Fatal(err)
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
