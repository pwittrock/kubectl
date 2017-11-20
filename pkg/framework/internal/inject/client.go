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

package inject

import (
	"flag"
	"path/filepath"
	"sync"

	"os"

	"strings"

	"github.com/googleapis/gnostic/OpenAPIv2"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

type Factory struct {
	sync.Once
	config        *rest.Config
	discovery     discovery.DiscoveryInterface
	resources     openapi.Resources
	openapischema *openapi_v2.Document
	rest          rest.Interface
	apiGroup      string
	apiVersion    string
}

var FactorySingleton = NewFactory()

func NewFactory() *Factory {
	f := &Factory{}
	f.inject()
	return f
}

// Hack to get around that we need to parse these flags before populating the flags that will be
// parsed by cobra - e.g. kubeconfig is needed to defined other flags which are based of the group and version
func getStringFlag(name, defaultVal string) string {
	found := false
	value := defaultVal
	for _, a := range os.Args {
		if found {
			value = a
			break
		}
		if strings.HasPrefix(a, "--") && strings.TrimPrefix(a, "--") == name {
			found = true
		}
	}
	return value
}

func (c *Factory) inject() *rest.Config {
	c.Do(func() {
		flag.Parse()
		var kubeconfig string
		if home := homeDir(); home != "" {
			kubeconfig = getStringFlag("kubeconfig", filepath.Join(home, ".kube", "config"))
		} else {
			kubeconfig = getStringFlag("kubeconfig", "")
		}
		c.apiGroup = getStringFlag("api-group", "")
		c.apiVersion = getStringFlag("api-version", "")
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		c.config = config

		clientset := kubernetes.NewForConfigOrDie(config)
		c.discovery = clientset.Discovery()

		openapischema, err := clientset.OpenAPISchema()
		if err != nil {
			panic(err.Error())
		}
		c.openapischema = openapischema

		resources, err := openapi.NewOpenAPIData(openapischema)
		if err != nil {
			panic(err.Error())
		}
		c.resources = resources

		c.rest = clientset.RESTClient()
	})
	return c.config
}

func (f *Factory) GetDiscovery() discovery.DiscoveryInterface {
	return f.discovery
}

func (f *Factory) GetResources() openapi.Resources {
	return f.resources
}

func (f *Factory) GetRest() rest.Interface {
	return f.rest
}

func (f *Factory) GetApiGroup() string {
	return f.apiGroup
}

func (f *Factory) GetApiVersion() string {
	return f.apiVersion
}

func (f *Factory) GetOpenapiSchema() *openapi_v2.Document {
	return f.openapischema
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
