/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"reflect"

	"k8s.io/kubernetes/pkg/api"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/conversion"
)

func addConversionFuncs() {
	// Add non-generated conversion functions
	err := api.Scheme.AddConversionFuncs(
		convert_api_PodSpec_To_v1_PodSpec,
		convert_v1_PodSpec_To_api_PodSpec,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}

// The following two PodSpec conversions functions where copied from pkg/api/conversion.go
// for the generated functions to work properly.
// This should be fixed: https://github.com/kubernetes/kubernetes/issues/12977
func convert_api_PodSpec_To_v1_PodSpec(in *api.PodSpec, out *v1.PodSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.PodSpec))(in)
	}
	if in.Volumes != nil {
		out.Volumes = make([]v1.Volume, len(in.Volumes))
		for i := range in.Volumes {
			if err := convert_api_Volume_To_v1_Volume(&in.Volumes[i], &out.Volumes[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	if in.Containers != nil {
		out.Containers = make([]v1.Container, len(in.Containers))
		for i := range in.Containers {
			if err := convert_api_Container_To_v1_Container(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	out.RestartPolicy = v1.RestartPolicy(in.RestartPolicy)
	if in.TerminationGracePeriodSeconds != nil {
		out.TerminationGracePeriodSeconds = new(int64)
		*out.TerminationGracePeriodSeconds = *in.TerminationGracePeriodSeconds
	} else {
		out.TerminationGracePeriodSeconds = nil
	}
	if in.ActiveDeadlineSeconds != nil {
		out.ActiveDeadlineSeconds = new(int64)
		*out.ActiveDeadlineSeconds = *in.ActiveDeadlineSeconds
	} else {
		out.ActiveDeadlineSeconds = nil
	}
	out.DNSPolicy = v1.DNSPolicy(in.DNSPolicy)
	if in.NodeSelector != nil {
		out.NodeSelector = make(map[string]string)
		for key, val := range in.NodeSelector {
			out.NodeSelector[key] = val
		}
	} else {
		out.NodeSelector = nil
	}
	out.ServiceAccountName = in.ServiceAccountName
	// DeprecatedServiceAccount is an alias for ServiceAccountName.
	out.DeprecatedServiceAccount = in.ServiceAccountName
	out.NodeName = in.NodeName
	out.HostNetwork = in.HostNetwork
	if in.ImagePullSecrets != nil {
		out.ImagePullSecrets = make([]v1.LocalObjectReference, len(in.ImagePullSecrets))
		for i := range in.ImagePullSecrets {
			if err := convert_api_LocalObjectReference_To_v1_LocalObjectReference(&in.ImagePullSecrets[i], &out.ImagePullSecrets[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ImagePullSecrets = nil
	}
	return nil
}

func convert_v1_PodSpec_To_api_PodSpec(in *v1.PodSpec, out *api.PodSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.PodSpec))(in)
	}
	if in.Volumes != nil {
		out.Volumes = make([]api.Volume, len(in.Volumes))
		for i := range in.Volumes {
			if err := convert_v1_Volume_To_api_Volume(&in.Volumes[i], &out.Volumes[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	if in.Containers != nil {
		out.Containers = make([]api.Container, len(in.Containers))
		for i := range in.Containers {
			if err := convert_v1_Container_To_api_Container(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	out.RestartPolicy = api.RestartPolicy(in.RestartPolicy)
	if in.TerminationGracePeriodSeconds != nil {
		out.TerminationGracePeriodSeconds = new(int64)
		*out.TerminationGracePeriodSeconds = *in.TerminationGracePeriodSeconds
	} else {
		out.TerminationGracePeriodSeconds = nil
	}
	if in.ActiveDeadlineSeconds != nil {
		out.ActiveDeadlineSeconds = new(int64)
		*out.ActiveDeadlineSeconds = *in.ActiveDeadlineSeconds
	} else {
		out.ActiveDeadlineSeconds = nil
	}
	out.DNSPolicy = api.DNSPolicy(in.DNSPolicy)
	if in.NodeSelector != nil {
		out.NodeSelector = make(map[string]string)
		for key, val := range in.NodeSelector {
			out.NodeSelector[key] = val
		}
	} else {
		out.NodeSelector = nil
	}
	// We support DeprecatedServiceAccount as an alias for ServiceAccountName.
	// If both are specified, ServiceAccountName (the new field) wins.
	out.ServiceAccountName = in.ServiceAccountName
	if in.ServiceAccountName == "" {
		out.ServiceAccountName = in.DeprecatedServiceAccount
	}
	out.NodeName = in.NodeName
	out.HostNetwork = in.HostNetwork
	if in.ImagePullSecrets != nil {
		out.ImagePullSecrets = make([]api.LocalObjectReference, len(in.ImagePullSecrets))
		for i := range in.ImagePullSecrets {
			if err := convert_v1_LocalObjectReference_To_api_LocalObjectReference(&in.ImagePullSecrets[i], &out.ImagePullSecrets[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ImagePullSecrets = nil
	}
	return nil
}
