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

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

type Factory struct {
	sync.Once
	config     *rest.Config
	discovery  discovery.DiscoveryInterface
	resources  openapi.Resources
	rest       rest.Interface
	apiGroup   *string
	apiVersion *string
}

func NewFactory() *Factory {
	f := &Factory{}
	f.inject()
	return f
}

func (c *Factory) inject() *rest.Config {
	c.Do(func() {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		c.apiGroup = flag.String("api-group", "", "use only API group")
		c.apiVersion = flag.String("api-version", "", "use only this API version")
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
	return *f.apiGroup
}

func (f *Factory) GetApiVersion() string {
	return *f.apiVersion
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
