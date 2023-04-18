# Prerequisites

* `kubectl`
* `docker`
* `go`

## Preparing the environment

1. Create a kind cluster.
    * ```bash
     kind create cluster --wait 5m
     ```
2. Install Crossplane v1.11.0 or later installed with composition functions
   feature flag enabled.
    * ```
     helm install crossplane --namespace crossplane-system crossplane-stable/crossplane \
     --create-namespace --wait \
     --set "args={--debug,--enable-composition-functions}" \
     --set "xfn.enabled=true" \
     --set "xfn.args={--debug}"
     ```

