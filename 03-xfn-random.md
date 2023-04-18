# Building a Function that Sets Random Value

With composition, there is no way to set a random value to a field because every
time it reconciles, it'll override the value that was set in the previous
reconciliation pass. In this tutorial, we'll create a function that parses all
`Robot` objects in the desired state and sets a random color to them if they
don't have one because it is a required parameter.

Let's build on top of our no-op function.
```bash
cp -a xfn-noop xfn-random
cd xfn-random
```

Change the function name to `xfn-random` in all files.
```bash
# On Mac
sed -i 's/xfn-noop/xfn-random/g' *
# On Linux
sed -i 's/xfn-noop/xfn-random/g' *
```

NOTE: The rest of the guide assumes that you already have the `CompositeResourceDefinition`
and `Composition` created from the previous tutorial. If you don't, you can
go back [installation section](02-xfn-noop.md#installation) and create them.

### Parsing Input

Our function currently does nothing, hence it doesn't need to parse the input it
receives. We will first need to parse the input as proper objects so that we can
set values.

Let's import the type of the input object.
```bash
go get github.com/crossplane/crossplane
go get sigs.k8s.io/yaml
```

Let's update the `main.go` file to parse the input as a `FunctionIO` object but
still do nothing.
```go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"sigs.k8s.io/yaml"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v", err)
		os.Exit(1)
	}
	obj := &v1alpha1.FunctionIO{}
	if err := yaml.Unmarshal(b, obj); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal stdin: %v", err)
		os.Exit(1)
	}
	result, err := yaml.Marshal(obj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v", err)
		os.Exit(1)
	}
	fmt.Print(string(result))
}
```

Let's define an array of colors that we'll use to set the random value.
```go
var (
    Colors = []string{"red", "green", "blue", "yellow", "orange", "purple", "black", "white"}
)
```

### Manipulating the Desired State

Now, we have a `FunctionIO` object that contains the desired and observed `Robot`
objects. In order to decide whether we should generate and set a random value,
we need to check whether the created objects already got one.

Our function will work on only the `Robot` objects so let's import the type.
```bash
go get github.com/upbound/provider-dummy
```

Here we extract a list of existing resources that do not have their color set.
```go
	alreadySet := map[string]bool{}
	for _, observed := range obj.Observed.Resources {
		if observed.Resource.Object.GetObjectKind().GroupVersionKind() != dummyv1alpha1.SchemeGroupVersion.WithKind("Robot") {
			// skip if it is not a Robot
			continue
		}
		r := &dummyv1alpha1.Robot{}
		if err := yaml.Unmarshal(observed.Resource.Raw, r); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal observed resource: %v", err)
			os.Exit(1)
		}
		if r.Spec.ForProvider.Color != "" {
			alreadySet[observed.Name] = true
		}
	}
```

In the next loop, we skip all the entries that already have a color set and
generate a random color for the rest.

```go
	for i, desired := range obj.Desired.Resources {
		if desired.Resource.Object.GetObjectKind().GroupVersionKind() != dummyv1alpha1.SchemeGroupVersion.WithKind("Robot") {
			// skip if it is not a Robot
			continue
		}
		if alreadySet[desired.Name] {
			continue
		}
		r := &dummyv1alpha1.Robot{}
		if err := yaml.Unmarshal(desired.Resource.Raw, r); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal observed resource: %v", err)
			os.Exit(1)
		}
		r.Spec.ForProvider.Color = Colors[rand.Intn(len(Colors))]
		raw, err := yaml.Marshal(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal resource: %v", err)
			os.Exit(1)
		}
		obj.Desired.Resources[i].Resource.Raw = raw
	}
```

## Try it Out

Let's build and push the function.
```bash
# This is to make sure go.sum is tidied up after all the go get commands.
go mod tidy
```
```bash
docker build --tag muvaf/xfn-random:v0.1.0 .
docker push muvaf/xfn-random:v0.1.0
```

Set the new image on our `Composition` object with `kubectl edit`.
```yaml
  ...
  functions:
  - name: my-random-function
    type: Container
    container:
      image: muvaf/xfn-random:v0.1.0
```

Edit `Composition` to add a second `Robot` object but this time without its
color parameter set. The full `resources` array should look like the following:
```yaml
  resources:
  - name: one-robot
    base:
      apiVersion: iam.dummy.upbound.io/v1alpha1
      kind: Robot
      spec:
        forProvider:
          color: yellow
  - name: second-robot
    base:
      apiVersion: iam.dummy.upbound.io/v1alpha1
      kind: Robot
```

Let's create a new `RobotGroup` object and see what happens.
```bash
cat <<EOF | kubectl apply -f -
apiVersion: contribfest.crossplane.io/v1alpha1
kind: RobotGroup
metadata:
  name: my-robot-group
spec: {}
EOF
```