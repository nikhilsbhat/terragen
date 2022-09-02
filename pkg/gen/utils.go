package gen

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	goVersion "github.com/hashicorp/go-version"
	"github.com/nikhilsbhat/terragen/version"
)

var (
	toCamel = template.FuncMap{
		"toCamel": snakeCaseToCamelCase,
	}
	filePerm        = 0o700
	dirPerm         = 0o777
	scaffoldPerm    = os.FileMode(filePerm)
	scaffoldDirPerm = os.FileMode(dirPerm)
)

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
	if len(i.Importer) != 0 { //nolint:staticcheck
		// to be implemented once we hit on importers.
	}
}

func snakeCaseToCamelCase(input string) (camelCase string) {
	isToUpper := false
	for key, value := range input {
		if key == 0 { //nolint:nestif
			camelCase = strings.ToLower(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(value))
				isToUpper = false
			} else {
				if value == '_' {
					isToUpper = true
				} else {
					camelCase += string(value)
				}
			}
		}
	}

	return
}

func lockTerragenExecution(currentVersion string, force bool) (oldVer, newVer string, lock bool, err error) {
	if force {
		return "", "", false, nil
	}
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
