package v1alpha1

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource01Spec struct {
	Field01 int    `json:"field01"`
	Field02 string `json:"field02"`
}

type Resource01 struct {
	metaV1.TypeMeta `json:",inline"`

	metaV1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec Resource01Spec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

func (in *Resource01) DeepCopyObject() runtime.Object {
	out := new(Resource01)
	*out = *in

	return out
}

type Resource01List struct {
	metaV1.TypeMeta `json:",inline"`

	metaV1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []*Resource01 `json:"items" protobuf:"bytes,2,rep,name=items"`
}
