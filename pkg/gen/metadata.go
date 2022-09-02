package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/terragen/pkg/decode"
	"github.com/nikhilsbhat/terragen/pkg/utils"
	"github.com/nikhilsbhat/terragen/version"
	"gopkg.in/yaml.v2"
)

type Metadata interface {
	Get(repoGroup, mod string) *Config
}

// Config would be generated and stored by the utility for further references.
type Config struct {
	// Version of terragen used for generating scaffolds. Updates only when higher version of terragen used.
	Version string `json:"version" yaml:"version"`
	// RepoGroup to which the project is part of.
	RepoGroup string `json:"repo-group" yaml:"repo-group"`
	// ProjectModule represents the module of the project
	ProjectModule string `json:"project-module" yaml:"project-module"`
	// Provider name that was scaffolded.
	Provider string `json:"provider" yaml:"provider"`
	// ProviderPath where scaffolds were created.
	ProviderPath string `json:"provider-path" yaml:"provider-path"`
	// Resources that where scaffolded.
	Resources []string `json:"resources" yaml:"resources"`
	// DataSources that where scaffolded.
	DataSources []string `json:"data-sources" yaml:"data-sources"`
	// Importers that where scaffolded.
	Importers []string `json:"importers" yaml:"importers"`
}

func (i *Input) CreateOrUpdateMetadata() error {
	metaDataPath := filepath.Join(i.Path, terragenMetadata)

	var metadata *Config
	if i.MetadataScaffolded() { //nolint:nestif
		currentMetaData, err := getCurrentMetadata(i.metaDataPath)
		if err != nil {
			return err
		}

		newMetaData := i.getMetadata()
		if newMetaData.Provider != currentMetaData.Provider {
			return fmt.Errorf("renaming provider is not supported once scaffolds are created") //nolint:goerr113
		}
		if utils.HasChange(currentMetaData.Resources, newMetaData.Resources) {
			currentMetaData.Resources = append(currentMetaData.Resources, newMetaData.Resources...)
		}
		if utils.HasChange(currentMetaData.DataSources, newMetaData.DataSources) {
			currentMetaData.DataSources = append(currentMetaData.DataSources, newMetaData.DataSources...)
		}
		if utils.HasChange(currentMetaData.Importers, newMetaData.Importers) {
			currentMetaData.Importers = append(currentMetaData.Importers, newMetaData.Importers...)
		}
		metadata = currentMetaData
	} else {
		metadata = i.getMetadata()
	}
	metaData, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}
	writer, err := getMetaWriter(metaDataPath)
	if err != nil {
		return err
	}
	_, err = writer.Write(metaData)
	if err != nil {
		return err
	}

	return nil
}

func (i *Input) MetadataScaffolded() bool {
	if _, fileErr := os.Stat(i.metaDataPath); os.IsNotExist(fileErr) {
		return false
	}

	return true
}

func (i *Input) getMetadata() *Config {
	return &Config{
		Version:       version.Version,
		RepoGroup:     i.RepoGroup,
		Provider:      i.Provider,
		ProviderPath:  i.Path,
		Resources:     i.Resource,
		DataSources:   i.DataSource,
		Importers:     []string{i.Importer},
		ProjectModule: i.mod,
	}
}

func getMetaWriter(path string) (*os.File, error) {
	if _, fileErr := os.Stat(path); os.IsNotExist(fileErr) {
		return os.Create(path)
	}

	return os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
}

func getCurrentMetadata(path string) (*Config, error) {
	var metadata Config
	meta, err := decode.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(meta, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func newMetadata() *Config {
	return &Config{}
}
