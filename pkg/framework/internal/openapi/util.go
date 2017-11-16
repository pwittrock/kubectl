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

package openapi

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (builder *cmdBuilderImpl) IsCmd(resource v1.APIResource) bool {
	return strings.Contains(resource.Name, "/") && !strings.HasSuffix(resource.Name, "/status")
}

func (builder *cmdBuilderImpl) Seen(resource v1.APIResource) bool {
	return builder.seen[]
}

