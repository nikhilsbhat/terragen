package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	goVersion "github.com/hashicorp/go-version"
	"github.com/nikhilsbhat/terragen/version"
)

func (i *Input) getUpdatedResourceNDataSources() error {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	metadata, err := i.getCurrentMetadata()
	if err != nil {
		return err
	}
	i.DataSource = append(i.DataSource, metadata.DataSources...)
	i.Resource = append(i.Resource, metadata.Resources...)
	return nil
}

func (i *Input) setMod() string {
	if len(i.RepoGroup) == 0 {
		i.RepoGroup = i.Provider
	}
	return fmt.Sprintf("%s/terraform-provider-%s", i.RepoGroup, i.Provider)
}

func (i *Input) enrichNames() {
	if len(i.Resource) != 0 {
		resource := make([]string, 0)
		for _, rs := range i.Resource {
			resource = append(resource, fmt.Sprintf("resource_%s", rs))
		}
		i.Resource = resource
	}
	if len(i.DataSource) != 0 {
		datasource := make([]string, 0)
		for _, rs := range i.DataSource {
			datasource = append(datasource, fmt.Sprintf("datasource_%s", rs))
		}
		i.DataSource = datasource
	}
	if len(i.Importer) != 0 {
		// to be implemented once we hit on importers.
	}
}

func (i *Input) snakeCaseToCamelCase(input string) (camelCase string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			camelCase = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func lockTerragenExecution(currentVersion string) (old, new string, lock bool, err error) {
	terragenVersionUnkown := "unknown-version"

	if len(currentVersion) == 0 {
		return terragenVersionUnkown, terragenVersionUnkown, true, err
	}

	runnigTerragenVersion, err := goVersion.NewVersion(version.GetBuildInfo().Version)
	if err != nil {
		return terragenVersionUnkown, terragenVersionUnkown, true, err
	}

	oldTerragenVersion, err := goVersion.NewVersion(currentVersion)
	if err != nil {
		return terragenVersionUnkown, runnigTerragenVersion.String(), true, err
	}

	if runnigTerragenVersion.LessThan(oldTerragenVersion) {
		return oldTerragenVersion.String(), runnigTerragenVersion.String(), true, nil
	}
	return oldTerragenVersion.String(), runnigTerragenVersion.String(), false, nil
}

func printData(data []byte) {
	fmt.Println(string(data))
}
