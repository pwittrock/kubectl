/*
Copyright 2016 The Kubernetes Authors.
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

// Note: the example only works with the code within the same release/branch.
package kubecmd

import (
    "flag"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "fmt"
    "strings"
    "sort"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/apimachinery/pkg/util/sets"
)

func main() {
    var kubeconfig *string
    if home := homeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // use the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    // create the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    discovery := clientset.Discovery()
    list, err := discovery.ServerResources()
    if err != nil {
        panic(err.Error())
    }

    openapischema, err := discovery.OpenAPISchema()
    if err != nil {
        panic(err.Error())
    }

    s, err := openapi.NewOpenAPIData(openapischema)
    if err != nil {
        panic(err.Error())
    }

    p := []string{}
    resources := []v1.APIResource{}
    for _, i := range list {
        for _, r := range i.APIResources {
            if strings.Contains(r.Name, "/") && !strings.HasSuffix(r.Name, "/status") {
                group := r.Group
                //if len(group) == 0 {
                //    group = "core"
                //}
                version := r.Version
                //if len(version) == 0 {
                //    version = "v1"
                //}

                gvk := schema.GroupVersionKind{group, version, r.Kind}
                s := s.LookupResource(gvk)
                if s != nil {
                    fmt.Printf("Schema found %v %v\n", r.Name, s)
                    kv := &kindVisitor{emptyVisitor{}, nil}
                    s.Accept(kv)
                    for _, f := range kv.GetFlags() {
                        fmt.Printf("Flag %+v\n", f)
                    }

                    p = append(p, fmt.Sprintf("%v/%v/%v %v %v", group, version, r.Kind, r.Name, strings.Join(r.Verbs, ",")))
                    resources = append(resources, r)
                }
            }
        }
    }
    sort.Strings(p)
    fmt.Println(strings.Join(p, "\n"))

}

type emptyVisitor struct {}
func (*emptyVisitor) VisitArray(*openapi.Array)          {}
func (*emptyVisitor) VisitMap(*openapi.Map)              {}
func (*emptyVisitor) VisitPrimitive(*openapi.Primitive)  {}
func (*emptyVisitor) VisitKind(k *openapi.Kind)          {}
func (*emptyVisitor) VisitReference(r openapi.Reference) {}

type kindVisitor struct {
    emptyVisitor
    flags []openAPIFlag
}

var blacklisted = sets.NewString("apiVersion", "metadata", "kind", "status")

func (visitor *kindVisitor) GetFlags() []openAPIFlag {
    return visitor.flags
}

func (visitor *kindVisitor) VisitKind(k *openapi.Kind) {
    for k, v := range k.Fields {
        if blacklisted.Has(k) {
            continue
        }
        fv := &fieldVisitor{visitor.emptyVisitor, k, nil, nil}
        v.Accept(fv)
        visitor.flags = append(visitor.flags, fv.GetFlags()...)
    }
}

func (visitor *kindVisitor) VisitReference(r openapi.Reference) {
    r.SubSchema().Accept(visitor)
}

type fieldVisitor struct {
    emptyVisitor
    name string
    path []string
    flags []openAPIFlag
}

type openAPIFlag struct {
    name string
    path []string
    description string
    flagType string
}

func (visitor *fieldVisitor) GetFlags() []openAPIFlag {
    return visitor.flags
}

func (visitor *fieldVisitor) VisitKind(k *openapi.Kind) {
    if visitor.name != "spec" {
        return
    }

    for k, v := range k.Fields {
        fv := &fieldVisitor{visitor.emptyVisitor, k, nil, nil}
        v.Accept(fv)
        visitor.flags = append(visitor.flags, fv.GetFlags()...)
    }
}

func (visitor *fieldVisitor) VisitPrimitive(p *openapi.Primitive) {
    visitor.flags = append(visitor.flags, openAPIFlag{visitor.name,visitor.path, p.Description, p.Type})
}

func (visitor *fieldVisitor) VisitReference(r openapi.Reference) {
    r.SubSchema().Accept(visitor)
}


func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}