
# Dynamic Command APIs

## Command requests target a Resource in the K8S cluster

- `pkg/apis/cli/v1alpha1/resource_command_types.go`

## Command requests target a Service in the K8S cluster

- `pkg/apis/cli/v1alpha1/service_command_types.go`

# Dynamic Command Examples

- `sample/resource_cmd.yaml`
- `sample/service_cmd.yaml`

# Install the CRD with the dynamic commands

- `kubectl apply -f sample/cli_v1alpha1_clitestresource.yaml`

# Try the Create command

- `go run main.go` // Shows the list of commands
- `go run main.go create deployment --name foo --namespace default --image nginx` // Create a Deployment from the Dynamic template

# Install the Service for the Service command

This will create the Service that backs the Service command

- `docker build . -t <tag>`
- `docker push <tag>`
- `kubectl run --image pwittrock/go-server nginx --labels app=nginx`
- `kubectl expose deploy nginx --port 80`

# Try the Random command

- `go run main.go` // Shows the list of commands
- `go run main.go create deployment --name foo --namespace default --image nginx` // Create a Deployment from the Dynamic template
