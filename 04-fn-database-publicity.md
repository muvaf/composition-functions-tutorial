
3. Install Google Cloud Platform (GCP) provider.
   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: pkg.crossplane.io/v1
   kind: Provider
   metadata:
     name: provider-gcp
   spec:
     package: xpkg.upbound.io/upbound/provider-gcp:v0.30.0
   ```