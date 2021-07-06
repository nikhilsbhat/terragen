// Package gen is the core of terragen, where the template generation happens.
package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	// Version of terragen used for generating scaffolds. Updates only when higher version of terragen used.
	Version string `json:"version" yaml:"version"`
	// RepoGroup to which the project is part of.
	RepoGroup string `json:"repo-group" yaml:"repo-group"`
	// ProjectModule represents the module of the project
	ProjectModule string `json:"project-module" yaml:"project-module"`
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
			"{{ $resource }}": {{ toCamel $resource }}(),
		{{- end }}
		},

		DataSourcesMap: map[string]*schema.Resource{
		{{- range $datasource := .DataSource }}
			"{{ $datasource }}": {{ toCamel $datasource }}(),
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
	i.setProvider()
	i.Provider = args[0]
	i.mod = i.setMod()
	provideFile := filepath.Join(i.Path, i.Provider, terragenProvider)

	log.Println(ui.Info(fmt.Sprintf("go module for scaffold would be: %s", i.mod)))
	if i.providerScaffolded() {
		log.Fatal(ui.Error(fmt.Sprintf("scaffolds for provider '%s' was already generated\n\t use"+
			" `terragen create -h` or `terragen edit -h` for more info", i.Provider)))
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

	providerData, err := renderTemplate(terragenProvider, i.TemplateRaw.ProviderTemp, i)
	if err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops rendering provider %s errored with: %v ", i.Provider, err)))
	}

	if i.DryRun {
		log.Println(ui.Info(fmt.Sprintf("provider '%s' would be created under '%s'", i.Provider, i.Path)))
		log.Println(ui.Info("contents of provider looks like"))
		printData(providerData)
	} else {
		if err = ioutil.WriteFile(provideFile, providerData, 0700); err != nil { //nolint:gosec
			log.Fatalf(ui.Error(fmt.Sprintf("oops scaffolding provider %s errored with: %v ", i.Provider, err)))
		}
	}

	if err = i.createOtherComponents(); err != nil {
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

func (i *Input) setProvider() {
	i.setDefaults()
	i.getTemplate()
	i.enrichNames()
	i.Path = i.getPath()
	i.AutoGenMessage = autoGenMessage
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

func (i *Input) createOtherComponents() error {
	if err := i.CreateMain(); err != nil {
		return fmt.Errorf("oops generating main.go scaffolds for provider %s failed with error: %v", i.Provider, err)
	}

	if i.ResourceRequired {
		if err := i.createResource(); err != nil {
			log.Fatal(ui.Error(err.Error()))
		}
	}

	if i.DatasourceRequired {
		if err := i.createDataSource(); err != nil {
			log.Fatal(ui.Error(err.Error()))
		}
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
	if err = ioutil.WriteFile(providerFile, updateData, 0700); err != nil { //nolint:gosec
		return err
	}
	return nil
}

func (i *Input) getUpdatedProviderData(currentProvider []byte) ([]byte, error) {
	updatedProvider, err := renderTemplate(terragenProvider, providerTemp, i)
	if err != nil {
		return nil, err
	}

	dmp := diffmatchpatch.New()
	providerDiff := dmp.DiffMain(string(currentProvider), string(updatedProvider), false)
	return []byte(dmp.DiffText2(providerDiff)), nil
}
