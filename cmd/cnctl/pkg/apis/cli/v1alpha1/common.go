/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

// Flag defines a cli flag that should be registered and available in request / output templates
type Flag struct {
	Type FlagType `json:"type"`

	Name string `json:"name"`

	Description string `json:"description"`

	// +optional
	StringValue string `json:"stringValue,omitempty"`

	// +optional
	StringSliceValue []string `json:"stringSliceValue,omitempty"`

	// +optional
	BoolValue bool `json:"boolValue,omitempty"`

	// +optional
	IntValue int32 `json:"intValue,omitempty"`

	// +optional
	FloatValue float64 `json:"floatValue,omitempty"`
}

// ResponseValue defines a value that should be parsed from a response and available in output templates
type ResponseValue struct {
	Name     string `json:"name"`
	JsonPath string `json:"jsonPath"`
}

type FlagType string

const (
	STRING       FlagType = "String"
	BOOL                  = "Bool"
	FLOAT                 = "Float"
	INT                   = "Int"
	STRING_SLICE          = "StringSlice"
)

type Command struct {
	// Use is the one-line usage message.
	Use string `json:"use"`

	// Path is the path to the sub-command.  Omit if the command is directly under the root command.
	// +optional
	Path []string `json:"path,omitempty"`

	// Short is the short description shown in the 'help' output.
	// +optional
	Short string `json:"short,omitempty"`

	// Long is the long message shown in the 'help <this-command>' output.
	// +optional
	Long string `json:"long,omitempty"`

	// Example is examples of how to use the command.
	// +optional
	Example string `json:"example,omitempty"`

	// Deprecated defines, if this command is deprecated and should print this string when used.
	// +optional
	Deprecated string `json:"deprecated,omitempty"`

	// Flags are the command line flags.
	//
	// Example:
	// 		  - name: namespace
	//    		type: String
	//    		stringValue: "default"
	//    		description: "deployment namespace"
	//
	// +optional
	Flags []Flag `json:"flags,omitempty"`

	// SuggestFor is an array of command names for which this command will be suggested -
	// similar to aliases but only suggests.
	SuggestFor []string `json:"suggestFor,omitempty"`

	// Aliases is an array of aliases that can be used instead of the first word in Use.
	Aliases []string `json:"aliases,omitempty"`

	// Version defines the version for this command. If this value is non-empty and the command does not
	// define a "version" flag, a "version" boolean flag will be added to the command and, if specified,
	// will print content of the "Version" variable.
	// +optional
	Version string `json:"version,omitempty"`
}
