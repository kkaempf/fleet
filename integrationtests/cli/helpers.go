package cli

import (
	"os"
	"strings"

	"github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/yaml"

	"github.com/onsi/gomega/gbytes"
)

const (
	AssetsPath = "../assets/"
	separator  = "-------\n"
	apiVersion = "apiVersion: fleet.cattle.io/v1alpha1"
)

func GetBundleFromOutput(buf *gbytes.Buffer) (*v1alpha1.Bundle, error) {
	bundle := &v1alpha1.Bundle{}
	err := yaml.Unmarshal(buf.Contents(), bundle)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func GetBundleListFromOutput(buf *gbytes.Buffer) ([]*v1alpha1.Bundle, error) {
	bundles := []*v1alpha1.Bundle{}
	bundlesWithSeparator := strings.ReplaceAll(string(buf.Contents()), apiVersion, separator+apiVersion)
	bundlesStr := strings.Split(bundlesWithSeparator, separator)
	for _, bundleStr := range bundlesStr {
		if bundleStr != "" {
			bundle := &v1alpha1.Bundle{}
			err := yaml.Unmarshal([]byte(bundleStr), bundle)
			if err != nil {
				return nil, err
			}
			bundles = append(bundles, bundle)
		}
	}
	return bundles, nil
}

func IsResourcePresentInBundle(resourcePath string, resources []v1alpha1.BundleResource) (bool, error) {
	resourceFile, err := os.ReadFile(resourcePath)
	if err != nil {
		return false, err
	}

	for _, resource := range resources {
		if strings.ReplaceAll(resource.Content, "\n", "") == strings.ReplaceAll(string(resourceFile), "\n", "") {
			return true, nil
		}
	}

	return false, nil
}
