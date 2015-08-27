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

package etcd

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/daemon"
	"k8s.io/kubernetes/pkg/registry/generic"
	etcdgeneric "k8s.io/kubernetes/pkg/registry/generic/etcd"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/storage"
)

// rest implements a RESTStorage for daemons against etcd
type REST struct {
	*etcdgeneric.Etcd
}

// daemonPrefix is the location for daemons in etcd, only exposed
// for testing
var daemonPrefix = "/daemons"

// NewREST returns a RESTStorage object that will work against daemons.
func NewREST(s storage.Interface) *REST {
	store := &etcdgeneric.Etcd{
		NewFunc: func() runtime.Object { return &api.Daemon{} },

		// NewListFunc returns an object capable of storing results of an etcd list.
		NewListFunc: func() runtime.Object { return &api.DaemonList{} },
		// Produces a path that etcd understands, to the root of the resource
		// by combining the namespace in the context with the given prefix
		KeyRootFunc: func(ctx api.Context) string {
			return etcdgeneric.NamespaceKeyRootFunc(ctx, daemonPrefix)
		},
		// Produces a path that etcd understands, to the resource by combining
		// the namespace in the context with the given prefix
		KeyFunc: func(ctx api.Context, name string) (string, error) {
			return etcdgeneric.NamespaceKeyFunc(ctx, daemonPrefix, name)
		},
		// Retrieve the name field of a daemon
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.Daemon).Name, nil
		},
		// Used to match objects based on labels/fields for list and watch
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return daemon.MatchDaemon(label, field)
		},
		EndpointName: "daemons",

		// Used to validate daemon creation
		CreateStrategy: daemon.Strategy,

		// Used to validate daemon updates
		UpdateStrategy: daemon.Strategy,

		Storage: s,
	}

	return &REST{store}
}
