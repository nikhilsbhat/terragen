package gen

import (
	"fmt"
	"io/ioutil"
	"log"
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

func resource_{{ (index .Resource .Index) }}() *schema.Resource {
	return &schema.Resource{
		CreateContext: resource_{{ (index .Resource .Index) }}Create,
		ReadContext:   resource_{{ (index .Resource .Index) }}Read,
		DeleteContext: resource_{{ (index .Resource .Index) }}Delete,
		UpdateContext: resource_{{ (index .Resource .Index) }}Update,
		Schema: map[string]*schema.Schema{},
	}
}

func resource_{{ (index .Resource .Index) }}Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource .Index) }}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource .Index) }}Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func resource_{{ (index .Resource .Index) }}Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
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
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use"+
			" `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if i.resourceScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for resource '%s' was already generated\n\t use"+
			" `terragen edit resource` to edit one \n\t run `terragen edit resource -h` for more info", i.Resource[0])))
	}

	if err := i.createResource(); err != nil {
		log.Fatal(ui.Error(err.Error()))
	}

	if err := i.updateProvider(); err != nil {
		log.Fatal(ui.Error(fmt.Sprintf("updated provider '%s' errored with datasource '%s' with: %v", i.Provider, i.DataSource[0], err)))
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (i *Input) createResource() error {
	for index, currentResource := range i.Resource {
		resourceFileName := fmt.Sprintf("resource_%s.go", currentResource)
		resourceFilePath := filepath.Join(i.Path, i.Provider)
		resourceFile := filepath.Join(resourceFilePath, resourceFileName)

		log.Println(ui.Info(fmt.Sprintf("scaffolds for resource '%s' would be generated under: '%s'", currentResource, resourceFilePath)))

		type resource struct {
			*Input
			Index int
		}

		data := &resource{i, index}

		resourceData, err := renderTemplate(terragenResource, i.TemplateRaw.ResourceTemp, data)
		if err != nil {
			return fmt.Errorf("oops rendering resource %s errored with: %v ", currentResource, err)
		}

		if i.DryRun {
			log.Println(ui.Info("contents of resource looks like"))
			fmt.Println(string(resourceData))
		} else {
			if err = terragenFileCreate(resourceFile); err != nil {
				return fmt.Errorf("oops creating resource errored with: %v ", err)
			}
			if err = ioutil.WriteFile(resourceFile, resourceData, 0700); err != nil { //nolint:gosec
				return fmt.Errorf("oops scaffolding resource %s errored with: %v ", currentResource, err)
			}
		}
	}
	return nil
}

func (i *Input) resourceScaffolded() bool {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	currentMetaData, err := i.getCurrentMetadata()
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return false
	}

	for _, resource := range i.Resource {
		if contains(currentMetaData.Resources, resource) {
			return true
		}
	}
	return false
}
