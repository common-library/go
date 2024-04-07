package v1alpha1_test

import (
	"testing"

	"github.com/common-library/go/kubernetes/resource/client"
	crd "github.com/common-library/go/kubernetes/resource/custom-resource-definition/project01/group01/v1alpha1"
	"github.com/common-library/go/kubernetes/resource/custom-resource/project01/group01/v1alpha1"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResource01(t *testing.T) {
	return

	restClient, err := client.GetClientForInCluster()
	if err != nil {
		t.Fatal(err)
	}

	resource01 := v1alpha1.Resource01{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: crd.Group + "/" + crd.Version,
			Kind:       "Resource01",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-01",
			Namespace: "test",
		},
		Spec: v1alpha1.Resource01Spec{
			Field01: 1,
			Field02: "value",
		},
	}

	if err := client.Post(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", &crd.Resource01); err != nil && apiErrors.IsAlreadyExists(err) == false {
		t.Fatal(err)
	} else if err := client.Post(restClient, crd.Group, crd.Version, resource01.Namespace, "resource01s", &resource01); err != nil && apiErrors.IsAlreadyExists(err) == false {
		t.Fatal(err)
	} else if err := client.Delete(restClient, crd.Group, crd.Version, resource01.Namespace, "resource01s", resource01.Name); err != nil {
		t.Fatal(err)
	} else if err := client.Delete(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", crd.Resource01.Name); err != nil {
		t.Fatal(err)
	}
}
