package cmd

import (
	"fmt"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	//"sort"

	"flag"

	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	openapifw "k8s.io/kubectl/pkg/framework/openapi"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

type configInit struct {
	sync.Once
	config *rest.Config
}

func (c *configInit) GetConfig() *rest.Config {
	c.Do(func() {
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
		c.config = config
	})
	return c.config
}

var config = configInit{}

func add(cmd *cobra.Command, verbs sets.String, name string) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config.GetConfig())
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

	resources, err := openapi.NewOpenAPIData(openapischema)
	if err != nil {
		panic(err.Error())
	}

	discovery.RESTClient().Post().Name("").SubResource().Body()

	cmds := map[string]*cobra.Command{}
	seen := map[string]sets.String{}

	cmdBuilder := openapifw.NewCmdBuilder(resources)
	for _, i := range list {
		for _, resource := range i.APIResources {
			if cmdBuilder.IsCmd(resource) {

				// Setup the command
				cmd, err := cmdBuilder.BuildCmd(resource)
				if err != nil {
					panic(err)
				}

				// Build the flags
				request, err := cmdBuilder.BuildFlags(cmd, resource)
				if err != nil {
					panic(err)
				}

				// Build the run function
				cmdBuilder.BuildRun(cmd, resource, request)

				ver := sets.NewString(r.Verbs...)
				if len(ver.Intersection(verbs)) == 0 {
					continue
				}

				if _, found := cmds[operation]; !found {
					cmds[operation] = &cobra.Command{
						Use:   operation,
						Short: "",
						Long:  ``,
						Run: func(cmd *cobra.Command, args []string) {
						},
					}
				}
				parent := cmds[operation]

				if _, found := seen[operation]; !found {
					seen[operation] = sets.NewString()
				}
				e := seen[operation]

				if e.Has(resource) {
					continue
				}
				e.Insert(resource)

				//fmt.Printf("Schema found %v %v\n", r.Name, s)
				kv := &kindVisitor{emptyVisitor{}, nil}
				s.Accept(kv)
				request := map[string]interface{}{}
				request["apiVersion"] = fmt.Sprintf("%v/%v", group, version)
				request["kind"] = fmt.Sprintf("%v", r.Kind)
				child := &cobra.Command{
					Use:   fmt.Sprintf("%v", resource),
					Short: fmt.Sprintf("%v %v for %v/%v/%v", name, operation, group, version, resource),
					Long:  ``,
					Run: func(cmd *cobra.Command, args []string) {
						out, _ := yaml.Marshal(request)
						fmt.Printf("%s\n", out)
					},
				}
				for _, f := range kv.GetFlags() {
					switch f.flagType {
					case "integer":
						ref := child.Flags().Int32(f.name, 0, f.description)
						r := request
						f.path = f.path[:len(f.path)-1]
						for _, p := range f.path {
							if _, found := r[p]; !found {
								n := map[string]interface{}{}
								r[p] = n
								r = n
							}
						}
						r[f.name] = ref
					}
				}
				parent.AddCommand(child)
			}
		}
	}
	//sort.Strings(p)
	//fmt.Println(strings.Join(p, "\n"))

	for _, v := range cmds {
		cmd.AddCommand(v)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type emptyVisitor struct{}

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
		fv := &fieldVisitor{visitor.emptyVisitor, k, []string{k}, nil}
		v.Accept(fv)
		visitor.flags = append(visitor.flags, fv.GetFlags()...)
	}
}

func (visitor *kindVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}

type fieldVisitor struct {
	emptyVisitor
	name  string
	path  []string
	flags []openAPIFlag
}

type openAPIFlag struct {
	name        string
	path        []string
	description string
	flagType    string
}

func (visitor *fieldVisitor) GetFlags() []openAPIFlag {
	return visitor.flags
}

func (visitor *fieldVisitor) VisitKind(k *openapi.Kind) {
	if visitor.name != "spec" {
		return
	}

	for k, v := range k.Fields {
		fv := &fieldVisitor{visitor.emptyVisitor, k, append(visitor.path, k), nil}
		v.Accept(fv)
		visitor.flags = append(visitor.flags, fv.GetFlags()...)
	}
}

func (visitor *fieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	visitor.flags = append(visitor.flags,
		openAPIFlag{visitor.name, visitor.path, p.Description, p.Type},
	)
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
