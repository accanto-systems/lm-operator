// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ALM) DeepCopyInto(out *ALM) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ALM.
func (in *ALM) DeepCopy() *ALM {
	if in == nil {
		return nil
	}
	out := new(ALM)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ALM) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ALMList) DeepCopyInto(out *ALMList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ALM, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ALMList.
func (in *ALMList) DeepCopy() *ALMList {
	if in == nil {
		return nil
	}
	out := new(ALMList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ALMList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ALMSpec) DeepCopyInto(out *ALMSpec) {
	*out = *in
	out.Configurator = in.Configurator
	out.Conductor = in.Conductor
	out.Apollo = in.Apollo
	out.Galileo = in.Galileo
	out.Talledega = in.Talledega
	out.Daytona = in.Daytona
	out.Nimrod = in.Nimrod
	out.Ishtar = in.Ishtar
	out.Relay = in.Relay
	out.Watchtower = in.Watchtower
	out.Doki = in.Doki
	out.Brent = in.Brent
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ALMSpec.
func (in *ALMSpec) DeepCopy() *ALMSpec {
	if in == nil {
		return nil
	}
	out := new(ALMSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ALMStatus) DeepCopyInto(out *ALMStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ALMStatus.
func (in *ALMStatus) DeepCopy() *ALMStatus {
	if in == nil {
		return nil
	}
	out := new(ALMStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfiguratorDescriptorSpec) DeepCopyInto(out *ConfiguratorDescriptorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfiguratorDescriptorSpec.
func (in *ConfiguratorDescriptorSpec) DeepCopy() *ConfiguratorDescriptorSpec {
	if in == nil {
		return nil
	}
	out := new(ConfiguratorDescriptorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NimrodDescriptorSpec) DeepCopyInto(out *NimrodDescriptorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NimrodDescriptorSpec.
func (in *NimrodDescriptorSpec) DeepCopy() *NimrodDescriptorSpec {
	if in == nil {
		return nil
	}
	out := new(NimrodDescriptorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceDescriptorSpec) DeepCopyInto(out *ServiceDescriptorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceDescriptorSpec.
func (in *ServiceDescriptorSpec) DeepCopy() *ServiceDescriptorSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceDescriptorSpec)
	in.DeepCopyInto(out)
	return out
}
