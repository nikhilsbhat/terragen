package gen

import (
	// resource template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikhilsbhat/terragen/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed templates/resource.tmpl
	resourceTemp string
	//go:embed templates/resource_v2.tmpl
	resourceV2Temp string
)

type Resource struct {
	Path               string
	DryRun             bool
	Force              bool
	SkipProviderUpdate bool
	AutoGenMessage     string
	Provider           string
	ResourceTemp       string
	Name               []string
	logger             *logrus.Logger
}

func (i *Input) GenerateResource(resources []string) error {
	i.Resource = resources
	i.AutoGenMessage = autoGenMessage
	i.enrichNames()
	i.Provider = strings.ReplaceAll(i.Provider, "-", "_")
	i.Path = filepath.Join(i.getPath(), i.Provider)
	i.getTemplate()
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)

	resource := NewResource(i)
	if !NewProvider(i).Scaffolded() {
		i.logger.Errorf("scaffolds for provider '%s' was not generated earlier", i.Provider)
		i.logger.Errorf("use `terragen create provider` to create one or run `terragen create provider -h` for more info")

		return fmt.Errorf("errored while validating 'provider' scaffolds")
	}

	if resource.Scaffolded() {
		i.logger.Errorf("scaffolds for resource '%s' was already generated", i.Resource[0])
		i.logger.Errorf("use `terragen edit resource` to edit one or run `terragen edit resource -h` for more info")

		return fmt.Errorf("errored while validating 'resource' scaffolds")
	}

	metadata, err := getCurrentMetadata(i.metaDataPath)
	if err != nil {
		return err
	}

	if oldVer, newVer, lock, err := lockTerragenExecution(metadata.Version, resource.Force); lock {
		if err != nil {
			return err
		}

		i.logger.Fatalf("terragen version '%v' or greater is required cannot scaffold more with terragen version '%v', it might breaks the project", oldVer, newVer)

		return fmt.Errorf("errored while checking the version compatibility of Terragen on the scaffolds")
	}

	if err := resource.Create(); err != nil {
		return err
	}

	if !i.SkipProviderUpdate {
		if err := NewProvider(i).Update(); err != nil {
			i.logger.Errorf("updating provider '%s' errored with resource '%s' with: %v", i.Provider, i.Resource[0], err)

			return fmt.Errorf("errored while updating provider")
		}
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		i.logger.Errorf("oops creating/updating metadata errored out with %v", err)

		return fmt.Errorf("errored while updating metadata")
	}

	return nil
}

func (r *Resource) Create() error {
	for index, currentResource := range r.Name {
		resourceFile := filepath.Join(r.Path, "internal", fmt.Sprintf("%s.go", currentResource))

		r.logger.Infof("scaffolds for resource '%s' would be generated under: '%s'", currentResource, r.Path)

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
			r.logger.Info("contents of resource looks like")
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
		r.logger.Fatal(err.Error())

		return false
	}

	for _, dataSource := range r.Name {
		if utils.Contains(currentMetaData.DataSources, dataSource) {
			return true
		}
	}

	for _, resource := range r.Name {
		if utils.Contains(currentMetaData.Resources, resource) {
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

//nolint:revive
func (r *Resource) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewResource(i *Input) *Resource {
	return &Resource{
		Path:               i.Path,
		DryRun:             i.DryRun,
		Force:              i.Force,
		AutoGenMessage:     i.AutoGenMessage,
		Provider:           i.Provider,
		ResourceTemp:       i.TemplateRaw.ResourceTemp,
		Name:               i.Resource,
		SkipProviderUpdate: i.SkipProviderUpdate,
		logger:             i.logger,
	}
}
