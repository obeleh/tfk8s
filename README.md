tfk8s
---

![](https://media.giphy.com/media/g8GfH3i5F0hby/giphy.gif)

`tfk8s` is a tool that makes it easier to work with the [Terraform Kubernetes Provider](https://github.com/hashicorp/terraform-provider-kubernetes-alpha).

If you want to copy examples from the Kubernetes documentation or migrate existing YAML manifests and use them with Terraform without having to convert YAML to HCL by hand, this tool is for you. 

## Demo 

[<img src="https://asciinema.org/a/jSmyAg4Ar6EcwKCTCXN8iAJM2.svg" width="250">](https://asciinema.org/a/jSmyAg4Ar6EcwKCTCXN8iAJM2)

## Features

- Convert a YAML file containing multiple manifests
- Strip out server side fields when piping `kubectl get $R -o yaml | tfk8s --strip`

## Install

```
go install github.com/jrhouston/tfk8s@latest
```

Alternatively, clone this repo and run:

```
make install
```

If Go's bin directory is not in your `PATH` you will need to add it:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Creating Terraform configurations

```
tfk8s -f input.yaml -o output.tf
```

or, using pipes: 
```
cat input.yaml | tfk8s > output.tf
```

**input.yaml**:
```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  TEST: test
```

✨✨ magically becomes ✨✨

**output.tf**:
```hcl
resource "kubernetes_manifest_hcl" "configmap_test" {
  manifest = {
    "apiVersion" = "v1"
    "data" = {
      "TEST" = "test"
    }
    "kind" = "ConfigMap"
    "metadata" = {
      "name" = "test"
    }
  }
}
```

### Common issues:

**error**:
```
'status' attribute key is not allowed in manifest configuration`
```
This means your yaml files contain the toplevel status field. 
It's not writable. Using `--strip` avoids this.

**solution**:
```
tfk8s -f input.yaml --strip -o output.tf
```

**error**:
```
The provider provider.kubernetes does not support resource type
"kubernetes_manifest".
```
**solution**:
You probably want to target the [kubernetes-alpha](https://github.com/hashicorp/terraform-provider-kubernetes-alpha) provider.
You can do that like so:
```
tfk8s -f input.yaml -p "kubernetes-alpha" -o output.tf
```


### Use with kubectl to output maps instead of YAML

```
kubectl get ns default -o yaml | tfk8s -M
```
```hcl
{
  "apiVersion" = "v1"
  "kind" = "Namespace"
  "metadata" = {
    "creationTimestamp" = "2020-05-02T15:01:32Z"
    "name" = "default"
    "resourceVersion" = "147"
    "selfLink" = "/api/v1/namespaces/default"
    "uid" = "6ac3424c-07a4-4a69-86ae-cc7a4ae72be3"
  }
  "spec" = {
    "finalizers" = [
      "kubernetes",
    ]
  }
  "status" = {
    "phase" = "Active"
  }
}
```
