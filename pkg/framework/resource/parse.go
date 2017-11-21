/*
Copyright 2017 The Kubernetes Authors.

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

package resource

import (
	"strings"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/framework/internal/inject"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// Parser discovers resources from an API server and parses them into
// indexed data structures
type Parser struct {
	resources  openapi.Resources
	discovery  discovery.DiscoveryInterface
	rest       rest.Interface
	apiGroup   string
	apiVersion string
}

// NewParser returns a new Parser
func NewParser() Parser {
	return Parser{
		inject.FactorySingleton.GetResources(),
		inject.FactorySingleton.GetDiscovery(),
		inject.FactorySingleton.GetRest(),
		inject.FactorySingleton.GetApiGroup(),
		inject.FactorySingleton.GetApiVersion(),
	}
}

// Resources discovers and indexes resources from the API server.
// It returns a map of resource name to resources matching that name ordered
// by preference as reported by the server
func (p *Parser) Resources() (Resources, error) {
	gvs, err := p.discovery.ServerResources()
	if err != nil {
		return nil, err
	}

	resources, byGVR := p.indexResources(gvs)
	err = p.attachSubResources(gvs, resources, byGVR)
	return resources, err
}

// subResource returns a Resource name, subresource name pair and true if Resource is subresource.
// Returns a Resource name, empty string and false if Resource is not a subresource.
func (*Parser) subResource(resource *v1.APIResource) (string, string, bool) {
	parts := strings.Split(resource.Name, "/")
	if len(parts) > 1 {
		return parts[0], parts[1], true
	}

	return parts[0], "", false
}

func (p *Parser) resource(resource *v1.APIResource) (string, bool) {
	r, _, b := p.subResource(resource)
	return r, !b
}

// copyGroupVersion copies the group and version from src to dest if either is missing from dest
// If the src group is empty and the dest group is empty, sets the dest group to "core"
func (*Parser) copyGroupVersion(src *v1.APIResource, dest *v1.APIResource) {
	if len(dest.Group) == 0 {
		dest.Group = src.Group
	}
	if len(dest.Group) == 0 {
		dest.Group = "core"
	}
	if len(dest.Version) == 0 {
		dest.Version = src.Version
	}
}

// splitGroupVersion splits the groupVersion string into the group and version components
func (*Parser) splitGroupVersion(groupVersion string) (string, string) {
	parts := strings.Split(groupVersion, "/")
	var group, version string

	if len(parts) > 1 {
		group = parts[0]
	}

	if len(parts) > 1 {
		version = parts[1]
	} else if len(parts) > 0 {
		version = parts[0]
	} else {
		version = "v1"
	}

	return group, version
}

func (p *Parser) indexResources(gvs []*v1.APIResourceList) (
	map[string][]*Resource,
	map[schema.GroupVersionResource]*Resource) {

	resources := map[string][]*Resource{}
	byGVR := map[schema.GroupVersionResource]*Resource{}

	// Find all resources
	for _, gv := range gvs {
		group, version := p.splitGroupVersion(gv.GroupVersion)
		if !p.isGroupVersionMatch(group, version) {
			continue
		}

		for _, r := range gv.APIResources {
			p.defaultGroupVersion(&r, group, version)

			name, isResource := p.resource(&r)
			if !isResource {
				glog.Infof("skipping non-resource %s %s/%s/%s", r.Name, group, version, r.Kind)
				continue
			}

			openapiSchema, found := p.getOpenAPI(group, version, r.Kind)
			if !found {
				glog.Infof("openapi schema not found for %s/%s/%s", group, version, r.Kind)
				continue
			}

			value := &Resource{
				Resource:        r,
				ApiGroupVersion: schema.GroupVersion{Group: group, Version: version},
				Schema:          openapiSchema,
			}

			byGVR[schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: r.Name,
			}] = value
			// Add to the list of resources under that name
			resources[name] = append(resources[name], value)
		}
	}

	return resources, byGVR
}

func (p *Parser) attachSubResources(
	gvs []*v1.APIResourceList,
	resources map[string][]*Resource,
	byGVR map[schema.GroupVersionResource]*Resource) error {

	// Find all subresources and attach to parents
	for _, gv := range gvs {
		group, version := p.splitGroupVersion(gv.GroupVersion)
		if !p.isGroupVersionMatch(group, version) {
			continue
		}

		for _, r := range gv.APIResources {
			p.defaultGroupVersion(&r, group, version)

			resourceName, _, isSubResource := p.subResource(&r)
			if !isSubResource {
				continue
			}

			openapiSchema, found := p.getOpenAPI(group, version, r.Kind)
			if !found {
				continue
			}

			// Make sure the Parent resources wasn't filtered out
			gvr := schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: resourceName,
			}
			if _, found := byGVR[gvr]; !found {
				continue
			}
			parent := byGVR[gvr]

			sub := &SubResource{
				Resource:        r,
				Parent:          parent,
				ApiGroupVersion: schema.GroupVersion{Group: group, Version: version},
				Schema:          openapiSchema,
			}

			parent.SubResources = append(parent.SubResources, sub)
		}
	}
	return nil
}

func (p *Parser) getOpenAPI(group, version, kind string) (openapi.Schema, bool) {
	openapiSchema := p.resources.LookupResource(schema.GroupVersionKind{
		group,
		version,
		kind,
	})
	if openapiSchema == nil {
		return nil, false
	}
	return openapiSchema, true
}

func (p *Parser) defaultGroupVersion(resource *v1.APIResource, group, version string) {
	// Set the group and version to the API groupVersion if missing
	if len(resource.Group) == 0 {
		resource.Group = group
	}
	if len(resource.Version) == 0 {
		resource.Version = version
	}
}

func (p *Parser) isGroupVersionMatch(group, version string) bool {
	if len(p.apiGroup) > 0 && p.apiGroup != group {
		return false
	}
	if len(p.apiVersion) > 0 && p.apiVersion != version {
		return false
	}
	return true
}
