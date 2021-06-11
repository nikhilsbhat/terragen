package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
	"github.com/nikhilsbhat/terragen/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	dataSourceTemp = `{{ .AutoGenMessage }}
package {{ .Provider }}

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasource_{{ (index .DataSource .Index) }}() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasource_{{ (index .DataSource .Index) }}Read,

		Schema: map[string]*schema.Schema{},
	}
}

func datasource_{{ (index .DataSource .Index) }}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}
`
)

func (i *Input) CreateDataSource(cmd *cobra.Command, args []string) {
	i.DataSource = args
	i.AutoGenMessage = autoGenMessage
	i.Path = i.getPath()
	i.getTemplate()

	//if oldVer, newVer, lock, err := i.lockTerragenExecution(); lock == true {
	//	if err != nil {
	//		log.Fatalf(ui.Error(err.Error()))
	//	}
	//	log.Fatalf("terragen version %v or greater is required\n cannot scaffold more with terragen version '%v', it breaks the project", oldVer, newVer)
	//}

	if !i.providerScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use"+
			" `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if i.dataSourceScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for data_source '%s' was already generated\n\t use"+
			" `terragen edit datasource` to edit one \n\t run `terragen edit datasource -h` for more info", i.DataSource[0])))
	}

	if err := i.createDataSource(); err != nil {
		log.Fatal(ui.Error(err.Error()))
	}

	if err := i.updateProvider(); err != nil {
		log.Fatal(ui.Error(fmt.Sprintf("updated provider '%s' errored with datasource '%s' with: %v", i.Provider, i.DataSource[0], err)))
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (i *Input) createDataSource() error {
	for index, dataSource := range i.DataSource {
		dataSourceFileName := fmt.Sprintf("datasource_%s.go", dataSource)
		dataSourceFilePath := filepath.Join(i.Path, i.Provider)
		dataSourceFile := filepath.Join(dataSourceFilePath, dataSourceFileName)

		log.Println(ui.Info(fmt.Sprintf("scaffolds for data-source '%s' would be generated under: '%s'", dataSource, dataSourceFilePath)))

		type datasoure struct {
			*Input
			Index int
		}

		data := &datasoure{i, index}

		dataSourceData, err := renderTemplate(terragenDataSource, i.TemplateRaw.DataTemp, data)
		if err != nil {
			return fmt.Errorf("oops rendering data_source %s errored with: %v ", dataSource, err)
		}

		if i.DryRun {
			log.Println(ui.Info("contents of data source looks like"))
			fmt.Println(string(dataSourceData))
		} else {
			if err = terragenFileCreate(dataSourceFile); err != nil {
				return fmt.Errorf("oops creating data source errored with: %v ", err)
			}
			if err = ioutil.WriteFile(dataSourceFile, dataSourceData, 0700); err != nil { //nolint:gosec
				return fmt.Errorf("oops scaffolding data_source %s errored with: %v ", dataSource, err)
			}
		}
	}
	return nil
}

func (i *Input) dataSourceScaffolded() bool {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	currentMetaData, err := i.getCurrentMetadata()
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return false
	}

	for _, dataSource := range i.DataSource {
		if utils.Contains(currentMetaData.DataSources, dataSource) {
			return true
		}
	}
	return false
}
