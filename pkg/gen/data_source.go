package gen

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
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

func datasource_{{ (index .DataSource 0) }}() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasource_{{ (index .DataSource 0) }}Read,

		Schema: map[string]*schema.Schema{},
	}
}

func datasource_{{ (index .DataSource 0) }}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if !i.providerScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if i.dataSourceScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for data_source '%s' was already generated\n\t use `terragen edit datasource` to edit one \n\t run `terragen edit datasource -h` for more info", i.DataSource[0])))
	}

	i.createDataSource()

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (i *Input) createDataSource() {
	dataSourceFileName := fmt.Sprintf("datasource_%s.go", i.DataSource[0])
	dataSourceFilePath := filepath.Join(i.Path, i.Provider)

	log.Println(ui.Info(fmt.Sprintf("scaffolds for data-source '%s' would be generated under: '%s'", i.DataSource[0], dataSourceFilePath)))

	var fileWriter io.Writer
	if i.DryRun {
		log.Println(ui.Info("contents of data source looks like"))
		fileWriter = os.Stdout
	} else {
		file, err := terragenWriter(dataSourceFilePath, dataSourceFileName)
		if err != nil {
			log.Fatalf("oops creating data source errored with: %v ", err)
		}
		defer file.Close()
		fileWriter = file
	}

	if len(i.TemplateRaw.DataTemp) != 0 {
		tmpl := template.Must(template.New(terragenDataSource).Parse(i.TemplateRaw.DataTemp))
		if err := tmpl.Execute(fileWriter, i); err != nil {
			log.Fatalf(ui.Error(fmt.Sprintf("oops scaffolding data_source %s errored with: %v ", i.DataSource, err)))
		}
	}
}

func (i *Input) dataSourceScaffolded() bool {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	currentMetaData, err := i.getCurrentMetadata()
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return false
	}

	if contains(currentMetaData.DataSources, i.DataSource[0]) {
		return true
	}
	return false
}
