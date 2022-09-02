package gen

import (
	// resource template has to be sourced from template.
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
	//go:embed templates/resource.tmpl
	resourceTemp string
	//go:embed templates/resource_v2.tmpl
	resourceV2Temp string
)

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
			return fmt.Errorf("oops rendering resource %s errored with: %w ", currentResource, err)
		}

		if r.DryRun {
			log.Println(ui.Info("contents of resource looks like"))
			printData(resourceData)

			return nil
		}
		if err = terragenFileCreate(resourceFile); err != nil {
			return fmt.Errorf("oops creating resource errored with: %w ", err)
		}
		if err = os.WriteFile(resourceFile, resourceData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding resource %s errored with: %w ", currentResource, err)
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

func (r *Resource) GetUpdated() error {
	return nil
}

func (r *Resource) Update() error {
	return nil
}

func (r *Resource) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
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
