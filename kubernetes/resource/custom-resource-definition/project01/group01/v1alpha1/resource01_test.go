package v1alpha1_test

import (
	"testing"

	"github.com/common-library/go/kubernetes/resource/client"
	"github.com/common-library/go/kubernetes/resource/custom-resource-definition/project01/group01/v1alpha1"
)

func TestResource01(t *testing.T) {
	return

	restClient, err := client.GetClientForInCluster()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Post(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", &v1alpha1.Resource01); err != nil {
		t.Fatal(err)
	} else if err := client.Delete(restClient, "apiextensions.k8s.io", "v1", "", "customresourcedefinitions", v1alpha1.Resource01.Name); err != nil {
		t.Fatal(err)
	}
}
