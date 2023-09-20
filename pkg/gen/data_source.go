package gen

import (
	// datasource template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikhilsbhat/terragen/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed templates/datasource.tmpl
	dataSourceTemp string
	//go:embed templates/datasource_v2.tmpl
	dataSourceV2Temp string
)

type DataSource struct {
	Path               string
	DryRun             bool
	Force              bool
	SkipProviderUpdate bool
	AutoGenMessage     string
	Provider           string
	DataTemp           string
	Name               []string
	logger             *logrus.Logger
}

func (i *Input) GenerateDataSource(dataSources []string) error {
	i.DataSource = dataSources
	i.AutoGenMessage = autoGenMessage
	i.enrichNames()
	i.Provider = strings.ReplaceAll(i.Provider, "-", "_")
	i.Path = filepath.Join(i.getPath(), i.Provider)
	i.getTemplate()
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)

	dataSource := NewDataSource(i)
	if !NewProvider(i).Scaffolded() {
		i.logger.Errorf("scaffolds for provider '%s' was not generated earlier", i.Provider)
		i.logger.Errorf("use `terragen create provider` to create one or run `terragen create provider -h` for more info")

		return fmt.Errorf("errored while validating 'provider' scaffolds")
	}

	if dataSource.Scaffolded() {
		i.logger.Errorf("scaffolds for data_source '%s' was already generated", i.DataSource[0])
		i.logger.Fatal("use `terragen edit datasource` to edit one or run `terragen edit datasource -h` for more info")

		return fmt.Errorf("errored while validating 'datasource' scaffolds")
	}

	metadata, err := getCurrentMetadata(i.metaDataPath)
	if err != nil {
		return err
	}

	if oldVer, newVer, lock, err := lockTerragenExecution(metadata.Version, dataSource.Force); lock {
		if err != nil {
			return err
		}

		i.logger.Fatalf("terragen version '%v' or greater is required cannot scaffold more with terragen version '%v', "+
			"it might breaks the project", oldVer, newVer)

		return fmt.Errorf("errored while checking the version compatibility of Terragen on the scaffolds")
	}

	if err := dataSource.Create(); err != nil {
		return err
	}

	if !i.SkipProviderUpdate {
		if err := NewProvider(i).Update(); err != nil {
			i.logger.Errorf("updating provider '%s' errored with datasource '%s' with: %v", i.Provider, i.DataSource[0], err)

			return fmt.Errorf("errored while updating provider")
		}
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		i.logger.Errorf("oops creating/updating metadata errored out with %v", err)

		return fmt.Errorf("errored while updating metadata")
	}

	return nil
}

func (d *DataSource) Create() error {
	for index, currentDataSource := range d.Name {
		dataSourceFile := filepath.Join(d.Path, "internal", fmt.Sprintf("%s.go", currentDataSource))

		d.logger.Infof("scaffolds for data-source '%s' would be generated under: '%s'", currentDataSource, d.Path)

		type dataSoure struct {
			*DataSource
			Index int
		}

		data := &dataSoure{d, index}

		dataSourceData, err := renderTemplate(terragenDataSource, d.DataTemp, data)
		if err != nil {
			return fmt.Errorf("oops rendering data_source %s errored with: %w ", currentDataSource, err)
		}

		if d.DryRun {
			d.logger.Infof("contents of data source looks like")
			printData(dataSourceData)

			return nil
		}

		if err = terragenFileCreate(dataSourceFile); err != nil {
			return fmt.Errorf("oops creating data source errored with: %w ", err)
		}

		if err = os.WriteFile(dataSourceFile, dataSourceData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding data_source %s errored with: %w ", currentDataSource, err)
		}
	}

	return nil
}

func (d *DataSource) Scaffolded() bool {
	currentMetaData, err := getCurrentMetadata(filepath.Join(d.Path, terragenMetadata))
	if err != nil {
		d.logger.Error(err.Error())

		return false
	}

	for _, dataSource := range d.Name {
		if utils.Contains(currentMetaData.DataSources, dataSource) {
			return true
		}
	}

	return false
}

func (d *DataSource) GetUpdated() error {
	return nil
}

func (d *DataSource) Update() error {
	return nil
}

//nolint:revive
func (d *DataSource) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewDataSource(i *Input) *DataSource {
	return &DataSource{
		Path:               i.Path,
		DryRun:             i.DryRun,
		Force:              i.Force,
		AutoGenMessage:     i.AutoGenMessage,
		Provider:           i.Provider,
		DataTemp:           i.TemplateRaw.DataTemp,
		Name:               i.DataSource,
		SkipProviderUpdate: i.SkipProviderUpdate,
		logger:             i.logger,
	}
}
