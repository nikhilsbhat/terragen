package gen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/terragen/pkg/decode"
	"gopkg.in/yaml.v2"
)

func (i *Input) CreateOrUpdateMetadata() error {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)

	if i.DryRun {
		log.Print(fmt.Sprintf("metadata would be generated under %s", i.metaDataPath))
		return nil
	}
	var metadata *Metadata
	if i.scaffoldStatus() {
		currentMetaData, err := i.getCurrentMetadata()
		if err != nil {
			return err
		}

		newMetaData := i.getMetadata()
		if newMetaData.Provider != currentMetaData.Provider {
			return fmt.Errorf("renaming provider is not supported once scaffolds are created")
		}
		if hasChange(currentMetaData.Resources, newMetaData.Resources) {
			currentMetaData.Resources = append(currentMetaData.Resources, newMetaData.Resources...)
		}
		if hasChange(currentMetaData.DataSources, newMetaData.DataSources) {
			currentMetaData.DataSources = append(currentMetaData.DataSources, newMetaData.DataSources...)
		}
		if hasChange(currentMetaData.Importers, newMetaData.Importers) {
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
	writer, err := i.getMetaWriter()
	if err != nil {
		return err
	}
	_, err = writer.Write(metaData)
	if err != nil {
		return err
	}
	return nil
}

func (i *Input) scaffoldStatus() bool {
	if _, fileErr := os.Stat(i.metaDataPath); os.IsNotExist(fileErr) {
		return false
	}
	return true
}

func (i *Input) getCurrentMetadata() (*Metadata, error) {
	var metadata Metadata
	meta, err := decode.ReadFile(i.metaDataPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(meta, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

func (i *Input) getMetadata() *Metadata {
	return &Metadata{
		Provider:     i.Provider,
		ProviderPath: i.Path,
		Resources:    i.Resource,
		DataSources:  i.DataSource,
		Importers:    []string{i.Importer},
	}
}

func (i *Input) getMetaWriter() (*os.File, error) {
	if _, fileErr := os.Stat(i.metaDataPath); os.IsNotExist(fileErr) {
		return os.Create(i.metaDataPath)
	}
	return os.OpenFile(i.metaDataPath, os.O_WRONLY, os.ModeAppend)
}
