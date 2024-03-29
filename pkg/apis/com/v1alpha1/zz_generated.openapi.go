// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/com/v1alpha1.ALM":           schema_pkg_apis_com_v1alpha1_ALM(ref),
		"./pkg/apis/com/v1alpha1.ALMSpec":       schema_pkg_apis_com_v1alpha1_ALMSpec(ref),
		"./pkg/apis/com/v1alpha1.ALMStatus":     schema_pkg_apis_com_v1alpha1_ALMStatus(ref),
		"./pkg/apis/com/v1alpha1.Service":       schema_pkg_apis_com_v1alpha1_Service(ref),
		"./pkg/apis/com/v1alpha1.ServiceSpec":   schema_pkg_apis_com_v1alpha1_ServiceSpec(ref),
		"./pkg/apis/com/v1alpha1.ServiceStatus": schema_pkg_apis_com_v1alpha1_ServiceStatus(ref),
	}
}

func schema_pkg_apis_com_v1alpha1_ALM(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ALM is the Schema for the alms API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ALMSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ALMStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/com/v1alpha1.ALMSpec", "./pkg/apis/com/v1alpha1.ALMStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_com_v1alpha1_ALMSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ALMSpec defines the desired state of ALM",
				Properties: map[string]spec.Schema{
					"springCloudConfigLabel": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"springProfilesActive": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"deploymentType": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"dockerRepo": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"conductor": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"apollo": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"galileo": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"talledega": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"daytona": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"nimrod": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"ishtar": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"relay": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"watchtower": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"doki": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
				},
				Required: []string{"springCloudConfigLabel", "springProfilesActive", "deploymentType", "dockerRepo", "conductor", "apollo", "galileo", "talledega", "daytona", "nimrod", "ishtar", "relay", "watchtower", "doki"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/com/v1alpha1.ServiceSpec"},
	}
}

func schema_pkg_apis_com_v1alpha1_ALMStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ALMStatus defines the observed state of ALM",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_com_v1alpha1_Service(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Service is the Schema for the services API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/com/v1alpha1.ServiceStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/com/v1alpha1.ServiceSpec", "./pkg/apis/com/v1alpha1.ServiceStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_com_v1alpha1_ServiceSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ServiceSpec defines the desired state of an ALM MicroService",
				Properties: map[string]spec.Schema{
					"jVMOptions": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"jVMOptions", "version"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_com_v1alpha1_ServiceStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ServiceStatus defines the observed state of Service",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
