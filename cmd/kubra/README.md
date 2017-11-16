
# Kubra

Kubra dynamically exposes commands by reading the subresources through the discovery service and exposing them as
commands.

# Install

- setup a minikube cluster
- `go install k8s.io/kubectl/cmd/kubra`
- `kubra set scale deployments --replicas 3`

