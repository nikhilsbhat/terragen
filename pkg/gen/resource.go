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

var resourceTemp = `{{ .AutoGenMessage }}
package {{ .Provider }}

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func {{ toCamel (index .Name .Index) }}() *schema.Resource {
	return &schema.Resource{
		CreateContext: {{ toCamel (index .Name .Index) }}Create,
		ReadContext:   {{ toCamel (index .Name .Index) }}Read,
		DeleteContext: {{ toCamel (index .Name .Index) }}Delete,
		UpdateContext: {{ toCamel (index .Name .Index) }}Update,
		Schema: map[string]*schema.Schema{},
	}
}

func {{ toCamel (index .Name .Index) }}Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func {{ toCamel (index .Name .Index) }}Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func {{ toCamel (index .Name .Index) }}Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}

func {{ toCamel (index .Name .Index) }}Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics  {
	// Your code goes here
	return nil
}
`

type Resource struct {
	Path           string
	DryRun         bool
	Force          bool
	AutoGenMessage string
	Provider       string
	ResourceTemp   string
	Name           []string
}

func (i *Input) GenerateResource(cmd *cobra.Command, args []string) {
	i.Resource = args
	i.AutoGenMessage = autoGenMessage
	i.enrichNames()
	i.Path = i.getPath()
	i.getTemplate()
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)

	resource := NewResource(i)
	if !NewProvider(i).Scaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was not generated earlier\n\t use"+
			" `terragen create provider` to create one \n\t run `terragen create provider -h` for more info", i.Provider)))
	}

	if resource.Scaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for resource '%s' was already generated\n\t use"+
			" `terragen edit resource` to edit one \n\t run `terragen edit resource -h` for more info", i.Resource[0])))
	}

	metadata, err := getCurrentMetadata(i.metaDataPath)
	if err != nil {
		log.Fatalf(ui.Error(err.Error()))
	}

	if oldVer, newVer, lock, err := lockTerragenExecution(metadata.Version, resource.Force); lock {
		if err != nil {
			log.Fatalf(ui.Error(err.Error()))
		}
		log.Fatalf("terragen version %v or greater is required\n cannot scaffold more with terragen version '%v', "+
			" it breaks the project", oldVer, newVer)
	}

	if err := resource.Create(); err != nil {
		log.Fatal(ui.Error(err.Error()))
	}

	if err := NewProvider(i).Update(); err != nil {
		log.Fatal(ui.Error(fmt.Sprintf("updating provider '%s' errored with resource '%s' with: %v", i.Provider, i.Resource[0], err)))
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}
}

func (r *Resource) Create() error {
	for index, currentResource := range r.Name {
		resourceFilePath := filepath.Join(r.Path, r.Provider)
		resourceFile := filepath.Join(resourceFilePath, fmt.Sprintf("%s.go", currentResource))

		log.Println(ui.Info(fmt.Sprintf("scaffolds for resource '%s' would be generated under: '%s'", currentResource, resourceFilePath)))

		type resource struct {
			*Resource
			Index int
		}

		data := &resource{r, index}

		resourceData, err := renderTemplate(terragenResource, r.ResourceTemp, data)
		if err != nil {
			return fmt.Errorf("oops rendering resource %s errored with: %v ", currentResource, err)
		}

		if r.DryRun {
			log.Println(ui.Info("contents of resource looks like"))
			printData(resourceData)
		} else {
			if err = terragenFileCreate(resourceFile); err != nil {
				return fmt.Errorf("oops creating resource errored with: %v ", err)
			}
			if err = ioutil.WriteFile(resourceFile, resourceData, scaffoldPerm); err != nil {
				return fmt.Errorf("oops scaffolding resource %s errored with: %v ", currentResource, err)
			}
		}
	}
	return nil
}

func (r *Resource) Scaffolded() bool {
	currentMetaData, err := getCurrentMetadata(filepath.Join(r.Path, terragenMetadata))
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return false
	}

	for _, dataSource := range r.Name {
		if utils.Contains(currentMetaData.DataSources, dataSource) {
			return true
		}
	}

	return false
}

func NewResource(i *Input) *Resource {
	return &Resource{
		Path:           i.Path,
		DryRun:         i.DryRun,
		Force:          i.Force,
		AutoGenMessage: i.AutoGenMessage,
		Provider:       i.Provider,
		ResourceTemp:   i.TemplateRaw.ResourceTemp,
		Name:           i.Resource,
	}
}
