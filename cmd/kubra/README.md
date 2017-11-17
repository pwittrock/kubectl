
# Kubra

Kubra dynamically exposes commands by reading the subresources through the discovery service and exposing them as
commands.

# Install

- setup a minikube cluster
- clone this fork+branch
- `go install k8s.io/kubectl/cmd/kubra`
- `kubectl create deployment nginx --replicas 3`
- `kubra set scale deployments --name nginx --replicas 3`
- `kubra get scale deployments --name nginx`

