// Package gen is the core of terragen, where the template generation happens.
package gen

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/neuron/cli/ui"
	"github.com/nikhilsbhat/terragen/pkg/decode"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

// Input holds the required values to generate the templates
type Input struct {
	// Resource to be created while generating scaffolds,
	// by passing a resource name here, it auto enabled ResourceRequired.
	// Provider name would be appended while constructing final resource name.
	// EX: resource 'create_cluster' for provider demo would become 'demo_create_cluster'.
	Resource []string
	// ResourceRequired determines if resource to be created while generating scaffolds.
	// Enabling this wth no resource name is not accepted.
	ResourceRequired bool
	// DataSource to be created while generating scaffolds,
	// by passing a resource name here, it auto enabled DatasourceRequired.
	// Provider name would be appended while constructing final data_source name.
	// EX: resource 'load_image' for provider demo would become 'demo_load_image'.
	DataSource []string
	// DatasourceRequired determines if data_source to be created while generating scaffolds.
	// Enabling this wth no data_source name is not accepted.
	DatasourceRequired bool
	// Importer to be created while generating scaffolds,
	// by passing a resource name here, it auto enabled ImporterRequired.
	Importer string
	// ImporterRequired determines if importer to be created while generating scaffolds.
	// Enabling this wth no importer name is not accepted.
	ImporterRequired bool
	// Provider name of which the scaffolds to be created, defaults to terraform-provider-demo
	Provider string
	// List of all the dependent packages for terraform, if not passed it picks default.
	Dependents []string
	// Path defines where the templates has to be generated.
	Path string
	// TemplateRaw consists of go-templates which are the core for terragen.
	TemplateRaw TerraTemplate
	// AutoGenMessage will be configured by terragen and cannot be overwritten.
	AutoGenMessage string
	// Description to be added to resource/datasource
	Description string
	// DryRun simulates scaffold creation by not creating one
	DryRun bool
	// RepoGroup is used while creating go mod. Defaults to 'github.com/test/'
	// For a given provider, repo group would be appended.
	// Ex: For provider 'demo' the go mod would looks 'github.com/test/demo'
	RepoGroup    string
	mod          string
	metaDataPath string
}

// TerraTemplate are the collections of go-templates which are used to generate terraform provider's base template.
type TerraTemplate struct {
	// ProviderTemp holds the template for provider
	ProviderTemp string `json:"provider-template" yaml:"provider-template"`
	// RootTemp holds the template for root file
	RootTemp string `json:"root-template" yaml:"provider-template"`
	// DataTemp holds the template for data
	DataTemp string `json:"data-template" yaml:"data-template"`
	// ResourceTemp holds the template for resource
	ResourceTemp string `json:"resource-template" yaml:"resource-template"`
	// GitIgnore that where scaffolded.
	GitIgnore string `json:"gitignore" yaml:"gitignore"`
}

// Metadata would be generated and stored by the utility for further references.
type Metadata struct {
	// Provider name that was scaffolded.
	Provider string `json:"provider" yaml:"provider"`
	// ProviderPath where scaffolds were created.
	ProviderPath string `json:"provider-path" yaml:"provider-path"`
	// Resources that where scaffolded.
	Resources []string `json:"resources" yaml:"resources"`
	// DataSources that where scaffolded.
	DataSources []string `json:"data-sources" yaml:"data-sources"`
	// Importers that where scaffolded.
	Importers []string `json:"importers" yaml:"importers"`
	// ProjectModule represents the module of the project
	ProjectModule string `json:"project-module" yaml:"project-module"`
}

var (
	providerTemp = `{{ .AutoGenMessage }}
package {{ .Provider }}

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
		{{- range $resource := .Resource }}
			"resource_{{ $resource }}": resource_{{ $resource }}(),
		{{- end }}
		},

		DataSourcesMap: map[string]*schema.Resource{
		{{- range $datasource := .DataSource }}
			"datasource_{{ $datasource }}": datasource_{{ $datasource }}(),
		{{- end }}
		},
	}
}
`

	autoGenMessage = `// ----------------------------------------------------------------------------
//
//     ***     TERRAGEN GENERATED CODE    ***    TERRAGEN GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file was auto generated by Terragen.
//     This autogenerated code has to be enhanced further to make it fully working terraform-provider.
//
//     Get more information on how terragen works.
//     https://github.com/nikhilsbhat/terragen
//
// ----------------------------------------------------------------------------`
)

// CreateProvider scaffolds the provider and its components as per the requirements.
func (i *Input) CreateProvider(cmd *cobra.Command, args []string) {
	i.setDefaults()
	i.getTemplate()
	i.Path = i.getPath()
	i.AutoGenMessage = autoGenMessage
	i.Provider = args[0]
	i.mod = i.setMod()

	log.Println(ui.Info(fmt.Sprintf("go module for scaffold would be: %s", i.mod)))
	if i.providerScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was already generated\n\t use `terragen create -h` or `terragen edit -h` for more info", i.Provider)))
	}

	log.Println(ui.Info(fmt.Sprintf("scaffolds for provider '%s' would be generated under: '%s'", i.Provider, i.Path)))

	if i.Dependents == nil {
		i.Dependents = []string{fmt.Sprintf("%s/%s", i.mod, i.Provider), "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"}
	}

	if !i.DryRun {
		if err := i.genTerraDir(); err != nil {
			log.Fatal(ui.Error(fmt.Sprintf("generating directories for scaffolds under %s failed with error: %v", i.Provider, err)))
		}
	}

	var fileWriter io.Writer
	if i.DryRun {
		log.Println(ui.Info(fmt.Sprintf("provider '%s' would be created under '%s'", i.Provider, i.Path)))
		log.Println(ui.Info("contents of provider looks like"))
		fileWriter = os.Stdout
	} else {
		providerFile, err := terragenWriter(filepath.Join(i.Path, i.Provider), terragenProvider)
		if err != nil {
			log.Fatal(ui.Error(err.Error()))
		}
		defer providerFile.Close()
		fileWriter = providerFile
	}

	if len(i.TemplateRaw.ProviderTemp) != 0 {
		tmpl := template.Must(template.New(terragenProvider).Parse(providerTemp))
		if err := tmpl.Execute(fileWriter, i); err != nil {
			log.Fatal(ui.Error(err.Error()))
		}
	}

	if err := i.createOtherComponents(); err != nil {
		log.Fatal(ui.Error(err.Error()))
	}

	// Setup the project to make it ready for development
	log.Println(ui.Info("terragen is in the process of making life simpler"))
	if !i.DryRun {
		if err := i.setupTerragen(); err != nil {
			log.Fatal(ui.Error(decode.GetStringOfMessage(err)))
		}
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops creating/updating metadata errored out with %v", err)))
	}

	log.Println(ui.Info("life is less complicated now ...!!"))
	log.Println(ui.Info(fmt.Sprintf("start enhancing terraform provider %s from the scaffold generated by terragen", i.Provider)))
}

func (i *Input) getPath() string {
	if i.Path == "." {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(ui.Error(decode.GetStringOfMessage(err)))
			os.Exit(1)
		}
		return dir
	}
	return path.Dir(i.Path)
}

func (i *Input) genTerraDir() error {
	err := os.MkdirAll(filepath.Join(i.Path, i.Provider), 0777)
	if err != nil {
		return err
	}
	return nil
}

// Set the templates to defaults if not specified.
func (i *Input) getTemplate() {
	if reflect.DeepEqual(i.TemplateRaw, TerraTemplate{}) {
		i.TemplateRaw.RootTemp = mainTemp
		i.TemplateRaw.ProviderTemp = providerTemp
		i.TemplateRaw.DataTemp = dataSourceTemp
		i.TemplateRaw.ResourceTemp = resourceTemp
		i.TemplateRaw.GitIgnore = gitignore
	}
}

func (i *Input) providerScaffolded() bool {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	if _, dirErr := os.Stat(filepath.Join(i.Path, terragenMetadata)); os.IsNotExist(dirErr) {
		return false
	}
	metadata, err := i.getCurrentMetadata()
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return true
	}
	return i.Provider == metadata.Provider
}

func (i *Input) providerExists() bool {
	if _, dirErr := os.Stat(filepath.Join(i.Path, terragenMetadata)); os.IsNotExist(dirErr) {
		return false
	}
	return true
}

func (i *Input) setDefaults() {
	if len(i.Resource) != 0 {
		i.ResourceRequired = true
	}
	if len(i.DataSource) != 0 {
		i.DatasourceRequired = true
	}
	if len(i.Importer) != 0 {
		i.ImporterRequired = true
	}
}

func (i *Input) setupTerragen() error {
	goInit := exec.Command("go", "mod", "init", i.mod) //nolint:gosec
	goInit.Dir = i.Path
	if err := goInit.Run(); err != nil {
		return err
	}

	goFmt := exec.Command("goimports", "-w", i.Path) //nolint:gosec
	goFmt.Dir = i.Path
	if err := goFmt.Run(); err != nil {
		return err
	}

	goVnd := exec.Command("go", "mod", "vendor") //nolint:gosec
	goVnd.Dir = i.Path
	if err := goVnd.Run(); err != nil {
		return err
	}
	return nil
}

func (i *Input) createOtherComponents() error {
	if err := i.CreateMain(); err != nil {
		return fmt.Errorf(fmt.Sprintf("oops generating main.go scaffolds for provider %s failed with error: %v", i.Provider, err))
	}

	if i.ResourceRequired {
		i.createResource()
	}

	if i.DatasourceRequired {
		i.createDataSource()
	}

	if err := i.createMakefile(); err != nil {
		return err
	}

	if err := i.createGitIgnore(); err != nil {
		log.Fatalf(ui.Error(err.Error()))
	}
	return nil
}

func (i *Input) updateProvider() error {
	currentProvider, err := ioutil.ReadFile(filepath.Join(i.Path, i.Provider, terragenProvider))
	if err != nil {
		return err
	}

	newIn := newInput()
	if err = copier.CopyWithOption(newIn, i, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	if err = newIn.getUpdatedResourceNDataSources(); err != nil {
		return err
	}

	updateData, err := newIn.getUpdatedProviderData(currentProvider)
	if err != nil {
		return err
	}

	providerFile := filepath.Join(newIn.Path, newIn.Provider, terragenProvider)
	if err = ioutil.WriteFile(providerFile, updateData, 0777); err != nil {
		return err
	}
	return nil
}

func (i *Input) getUpdatedProviderData(currentProvider []byte) ([]byte, error) {
	var updatedProvider bytes.Buffer
	tmpl := template.Must(template.New(terragenProvider).Parse(providerTemp))
	if err := tmpl.Execute(&updatedProvider, i); err != nil {
		return nil, err
	}

	dmp := diffmatchpatch.New()
	providerDiff := dmp.DiffMain(string(currentProvider), updatedProvider.String(), false)
	return []byte(dmp.DiffText2(providerDiff)), nil
}

func terragenWriter(path, file string) (*os.File, error) {
	fileToBeWritten, err := os.Create(filepath.Join(path, file))
	return fileToBeWritten, err
}

func (i *Input) getUpdatedResourceNDataSources() error {
	i.metaDataPath = filepath.Join(i.Path, terragenMetadata)
	metadata, err := i.getCurrentMetadata()
	if err != nil {
		return err
	}
	i.DataSource = append(i.DataSource, metadata.DataSources...)
	i.Resource = append(i.Resource, metadata.Resources...)
	return nil
}

func (i *Input) setMod() string {
	if len(i.RepoGroup) == 0 {
		i.RepoGroup = i.Provider
	}
	return fmt.Sprintf("%s/terraform-provider-%s", i.RepoGroup, i.Provider)
}

func newMetadata() *Metadata {
	return &Metadata{}
}

func newInput() *Input {
	return &Input{}
}
