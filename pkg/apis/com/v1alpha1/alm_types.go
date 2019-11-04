package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceDescriptorSpec defines the desired state of an ALM MicroService
// +k8s:openapi-gen=true
type ServiceDescriptorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	JVMOptions string `json:"JVMOptions"`
	Version    string `json:"Version"`
}

// NimrodDescriptorSpec defines the desired state of the Nimrod ALM MicroService
// +k8s:openapi-gen=true
type NimrodDescriptorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	JVMOptions string `json:"JVMOptions"`
	Version    string `json:"Version"`
	// lm-themes
	ThemesConfigMap string `json:"ThemesConfigMap"`
	// lm-locales
	LocalesConfigMap string `json:"LocalesConfigMap"`
}

// ConfiguratorDescriptorSpec defines the desired state of the Configurator ALM MicroService
// +k8s:openapi-gen=true
type ConfiguratorDescriptorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	JVMOptions string `json:"JVMOptions"`
	Run        bool   `json:"Run"`
}

// ALMSpec defines the desired state of ALM
// +k8s:openapi-gen=true
type ALMSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Secure                 bool                       `json:"secure"`
	SpringCloudConfigLabel string                     `json:"springCloudConfigLabel"`
	SpringProfilesActive   string                     `json:"springProfilesActive"`
	DeploymentType         string                     `json:"deploymentType"`
	DockerRepo             string                     `json:"dockerRepo"`
	Release                string                     `json:"release"`
	Configurator           ConfiguratorDescriptorSpec `json:"configurator"`
	Conductor              ServiceDescriptorSpec      `json:"conductor"`
	Apollo                 ServiceDescriptorSpec      `json:"apollo"`
	Galileo                ServiceDescriptorSpec      `json:"galileo"`
	Talledega              ServiceDescriptorSpec      `json:"talledega"`
	Daytona                ServiceDescriptorSpec      `json:"daytona"`
	Nimrod                 NimrodDescriptorSpec       `json:"nimrod"`
	Ishtar                 ServiceDescriptorSpec      `json:"ishtar"`
	Relay                  ServiceDescriptorSpec      `json:"relay"`
	Watchtower             ServiceDescriptorSpec      `json:"watchtower"`
	Doki                   ServiceDescriptorSpec      `json:"doki"`
	Brent                  ServiceDescriptorSpec      `json:"brent"`
}

// ALMStatus defines the observed state of ALM
// +k8s:openapi-gen=true
type ALMStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	IshtarHealthy bool `json:"ishtarHealthy"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ALM is the Schema for the alms API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ALM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ALMSpec   `json:"spec,omitempty"`
	Status ALMStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ALMList contains a list of ALM
type ALMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ALM `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ALM{}, &ALMList{})
}
