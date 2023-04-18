package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	dummyv1alpha1 "github.com/upbound/provider-dummy/apis/iam/v1alpha1"
	"sigs.k8s.io/yaml"
)

var (
	Colors = []string{"red", "green", "blue", "yellow", "orange", "purple", "black", "white"}
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
	result, err := yaml.Marshal(obj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v", err)
		os.Exit(1)
	}
	fmt.Print(string(result))
}
