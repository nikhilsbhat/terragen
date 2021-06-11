package gen

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

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

func (i *Input) lockTerragenExecution() (old, new float64, lock bool, err error) {
	runnigTerragenVersion, err := getTerragenVersion(version.Version)
	if err != nil {
		return 0, runnigTerragenVersion, true, err
	}
	log.Printf("current version %s", runnigTerragenVersion)

	metadata, err := i.getCurrentMetadata()
	if err != nil || len(metadata.Version) == 0 {
		return 0, runnigTerragenVersion, true, err
	}
	log.Printf("old version %s", metadata.Version)

	oldTerragenVersion, err := getTerragenVersion(metadata.Version)
	if err != nil {
		return oldTerragenVersion, runnigTerragenVersion, true, err
	}

	log.Printf("old version %s", oldTerragenVersion)
	if runnigTerragenVersion < oldTerragenVersion {
		return oldTerragenVersion, runnigTerragenVersion, true, nil
	}
	return oldTerragenVersion, runnigTerragenVersion, false, nil
}

func getTerragenVersion(v string) (float64, error) {
	fmt.Println(v)
	ver, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}
	return ver, nil
}
