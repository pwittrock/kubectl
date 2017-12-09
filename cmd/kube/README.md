# Kube

Kube is a commandline binary for talking with kubernetes.  It has all of the kubectl commands
embedded under `kube ctl`, but has other command groups as well for new commands developed outside
of the kubernetes/kubernetes repo.

Notably, the `kube ctl` command has both 1.8 and 1.7 versions of kubectl installed, and if the
apiserver publishes its version, `kube ctl` will use the version of `kubectl` matching the server.
If the server version is not a known or supported version, `kube ctl` defaults to 1.8.

## Instructions to install

### Create the direcoty

From your `GOPATH/src` directory:

- `mkdir -p k8s.io`
- `cd k8s.io`

### Clone this fork and branch

- `git clone https://github.com/pwittrock/kubectl.git --branch kubra --depth 1`

### Install kube

Build and install the `kube` binary under `GOBIN`

- `cd kubectl`
- `go install k8s.io/kubectl/cmd/kube`

