package v1alpha1

import (
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Resource01 = apiextensionsV1.CustomResourceDefinition{
	TypeMeta: metaV1.TypeMeta{
		APIVersion: "apiextensions.k8s.io/v1",
		Kind:       "CustomResourceDefinition",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name: "resource01s" + "." + Group,
	},
	Spec: apiextensionsV1.CustomResourceDefinitionSpec{
		Group: Group,
		Names: apiextensionsV1.CustomResourceDefinitionNames{
			Plural:   "resource01s",
			Singular: "resource01",
			Kind:     "Resource01",
		},
		Scope: apiextensionsV1.ResourceScope("Namespaced"),
		Versions: []apiextensionsV1.CustomResourceDefinitionVersion{
			apiextensionsV1.CustomResourceDefinitionVersion{
				Name:    Version,
				Served:  true,
				Storage: true,
				Schema: &apiextensionsV1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsV1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]apiextensionsV1.JSONSchemaProps{
							"spec": apiextensionsV1.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apiextensionsV1.JSONSchemaProps{
									"field01": apiextensionsV1.JSONSchemaProps{Type: "integer"},
									"field02": apiextensionsV1.JSONSchemaProps{Type: "string"},
								},
							},
						},
					},
				},
			},
		},
	},
}
