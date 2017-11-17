
# Kubra

Kubra dynamically exposes commands by reading the subresources listed by the discovery service and exposing
them directly on the commandline by parsing their request schemas and exposing the fields as flags.

# Try it out

## Install

- setup a minikube cluster
- clone this [fork+branch](https://github.com/pwittrock/kubectl/tree/kubra)
- compile kubra with `go install k8s.io/kubectl/cmd/kubra`

## Start invoking subresources

Create the deployment

- `kubectl create deployment nginx`
- `kubectl get deployment nginx -o yaml`

Update the scale using the subresource

- `kubra do scale deployments --name nginx --replicas 3`
- `kubra read scale deployments --name nginx`
- `kubectl get deployment nginx -o yaml`

Change the container image, then roll it back

- `kubectl edit deployment nginx`
- `kubectl get deployment nginx -o yaml`
- `kubra do rollback deployments --name nginx`
- `kubectl get deployment nginx -o yaml`

