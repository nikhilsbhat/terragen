package gen

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
	"github.com/nikhilsbhat/terragen/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	//go:embed templates/datasource.tmpl
	dataSourceTemp string
	//go:embed templates/datasource_v2.tmpl
	dataSourceV2Temp string
)

type DataSource struct {
	Path           string
	DryRun         bool
	Force          bool
	AutoGenMessage string
	Provider       string
	DataTemp       string
	Name           []string
}

func (i *Input) GenerateDataSource(cmd *cobra.Command, args []string) {
	i.DataSource = args
	i.AutoGenMessage = autoGenMessage
	i.enrichNames()
	i.Path = i.getPath()
	i.getTemplate()
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)

	dataSource := NewDataSource(i)
	if !NewProvider(i).Scaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use"+
			" `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if dataSource.Scaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for data_source '%s' was already generated\n\t use"+
			" `terragen edit datasource` to edit one \n\t run `terragen edit datasource -h` for more info", i.DataSource[0])))
	}

	metadata, err := getCurrentMetadata(i.metaDataPath)
	if err != nil {
		log.Fatalf(ui.Error(err.Error()))
	}

	if oldVer, newVer, lock, err := lockTerragenExecution(metadata.Version, dataSource.Force); lock {
		if err != nil {
			log.Fatalf(ui.Error(err.Error()))
		}
		log.Fatalf("terragen version %v or greater is required\n cannot scaffold more with terragen version '%v', "+
			" it breaks the project", oldVer, newVer)
	}

	if err := dataSource.Create(); err != nil {
		log.Fatal(ui.Error(err.Error()))
	}

	if err := NewProvider(i).Update(); err != nil {
		log.Fatal(ui.Error(fmt.Sprintf("updating provider '%s' errored with datasource '%s' with: %v", i.Provider, i.DataSource[0], err)))
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (d *DataSource) Create() error {
	for index, currentDataSource := range d.Name {
		dataSourceFilePath := filepath.Join(d.Path, d.Provider)
		dataSourceFile := filepath.Join(dataSourceFilePath, fmt.Sprintf("%s.go", currentDataSource))

		log.Println(ui.Info(fmt.Sprintf("scaffolds for data-source '%s' would be generated under: '%s'", currentDataSource, dataSourceFilePath)))

		type datasoure struct {
			*DataSource
			Index int
		}

		data := &datasoure{d, index}

		dataSourceData, err := renderTemplate(terragenDataSource, d.DataTemp, data)
		if err != nil {
			return fmt.Errorf("oops rendering data_source %s errored with: %v ", currentDataSource, err)
		}

		if d.DryRun {
			log.Println(ui.Info("contents of data source looks like"))
			printData(dataSourceData)
			return nil
		} else {
			if err = terragenFileCreate(dataSourceFile); err != nil {
				return fmt.Errorf("oops creating data source errored with: %v ", err)
			}
			if err = os.WriteFile(dataSourceFile, dataSourceData, scaffoldPerm); err != nil {
				return fmt.Errorf("oops scaffolding data_source %s errored with: %v ", currentDataSource, err)
			}
		}
	}

	return nil
}

func (d *DataSource) Scaffolded() bool {
	currentMetaData, err := getCurrentMetadata(filepath.Join(d.Path, terragenMetadata))
	if err != nil {
		log.Println(ui.Error(err.Error()))
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

func (d *DataSource) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewDataSource(i *Input) *DataSource {
	return &DataSource{
		Path:           i.Path,
		DryRun:         i.DryRun,
		Force:          i.Force,
		AutoGenMessage: i.AutoGenMessage,
		Provider:       i.Provider,
		DataTemp:       i.TemplateRaw.DataTemp,
		Name:           i.DataSource,
	}
}
