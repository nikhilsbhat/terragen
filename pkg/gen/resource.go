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
	resourceTemp = `{{ .AutoGenMessage }}
package {{ .Provider }}

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resource_{{ (index .Resource 0) }}() *schema.Resource {
	return &schema.Resource{
		CreateContext: resource_{{ (index .Resource 0) }}Create,
		ReadContext:   resource_{{ (index .Resource 0) }}Read,
		DeleteContext: resource_{{ (index .Resource 0) }}Delete,
		UpdateContext: resource_{{ (index .Resource 0) }}Update,
		Schema: map[string]*schema.Schema{},
	}
}

func resource_{{ (index .Resource 0) }}Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource 0) }}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource 0) }}Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource 0) }}Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}
`
)

func (i *Input) CreateResource(cmd *cobra.Command, args []string) {
	i.Resource = args
	i.AutoGenMessage = autoGenMessage
	i.Path = i.getPath()
	i.getTemplate()

	if !i.providerScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if i.resourceScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for resource '%s' was already generated\n\t use `terragen edit resource` to edit one \n\t run `terragen edit resource -h` for more info", i.Resource[0])))
	}

	i.createResource()

	if err := i.updateProvider(); err != nil {
		log.Fatal(ui.Error(fmt.Sprintf("updated provider '%s' errored with datasource '%s' with: %v", i.Provider, i.DataSource[0], err)))
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (i *Input) createResource() {
	resourceFileName := fmt.Sprintf("resource_%s.go", i.Resource[0])
	resourceFilePath := filepath.Join(i.Path, i.Provider)

	log.Println(ui.Info(fmt.Sprintf("scaffolds for resource '%s' would be generated under: '%s'", i.Resource[0], resourceFilePath)))

	var fileWriter io.Writer
	if i.DryRun {
		log.Println(ui.Info("contents of resource looks like"))
		fileWriter = os.Stdout
	} else {
		file, err := terragenWriter(resourceFilePath, resourceFileName)
		if err != nil {
			log.Fatal(ui.Error(fmt.Sprintf("oops creating data source errored with: %v ", err)))
		}
		defer file.Close()
		fileWriter = file
	}

	if len(i.TemplateRaw.ResourceTemp) != 0 {
		tmpl := template.Must(template.New(terragenResource).Parse(i.TemplateRaw.ResourceTemp))
		if err := tmpl.Execute(fileWriter, i); err != nil {
			log.Fatal(ui.Error(fmt.Sprintf("oops scaffolding resource %s errored with: %v ", i.Resource, err)))
		}
	}
}

func (i *Input) resourceScaffolded() bool {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	currentMetaData, err := i.getCurrentMetadata()
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return false
	}

	if contains(currentMetaData.Resources, i.Resource[0]) {
		return true
	}
	return false
}
