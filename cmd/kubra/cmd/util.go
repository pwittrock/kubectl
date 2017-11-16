package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	openapifw "k8s.io/kubectl/pkg/framework/openapi"
)

func add(parent *cobra.Command, verbs sets.String, name string) {
	cmdBuilder := openapifw.NewCmdBuilder()

	list, err := cmdBuilder.ListResources()
	if err != nil {
		panic(err)
	}
	for _, resource := range list {
		if cmdBuilder.IsCmd(resource) {
			// Don't expose multiple versions of the same resource
			if cmdBuilder.Seen(resource) {
				continue
			}

			// Setup the command
			cmd, err := cmdBuilder.BuildCmd(resource)
			if err != nil {
				panic(err)
			}
			parent.AddCommand(cmd)

			// Build the flags
			request, err := cmdBuilder.BuildFlags(cmd, resource)
			if err != nil {
				panic(err)
			}

			// Build the run function
			cmdBuilder.BuildRun(cmd, resource, request)

			//ver := sets.NewString(r.Verbs...)
			//if len(ver.Intersection(verbs)) == 0 {
			//	continue
			//}
		}
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
