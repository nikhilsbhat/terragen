// Package gen is the core of terragen, where the template generation happens.
package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// Terragen implements methods to GET/CREATE/UPDATE the scaffolds for the Terraform provider.
type Terragen interface {
	Create() error
	Scaffolded() bool
	Update() error
	Get(currentContent []byte) ([]byte, error)
	GetUpdated() error
}

// Scaffold implement method that generates the scaffolds for the Terraform provider.
type Scaffold interface {
	Generate(providerName string) error
}

// Input holds the required values to generate the templates.
type Input struct {
	// DryRun simulates scaffold creation by not creating one
	DryRun bool
	// Force will forcefully scaffold the datasource/resource by not validating the terragen version.
	// Enabling this might tamper the scaffolds.
	Force bool
	// SkipValidation will skip validating all the prerequisites such as checking go,goimports etc.
	SkipValidation bool
	// TerraformPluginFramework would generate scaffolds
	// with terraform-plugin-framework(https://github.com/hashicorp/terraform-plugin-framework).
	TerraformPluginFramework bool
	// ResourceRequired determines if resource to be created while generating scaffolds.
	// Enabling this wth no resource name is not accepted.
	ResourceRequired bool
	// DatasourceRequired determines if data_source to be created while generating scaffolds.
	// Enabling this wth no data_source name is not accepted.
	DatasourceRequired bool
	// SkipProviderUpdate when set, would skip updating the provider after generating the `data-sources` or `resources`.
	// When it is enabled, updating the provider with newer `data-sources` or `resources` has to be done manually.
	SkipProviderUpdate bool
	// ImporterRequired determines if importer to be created while generating scaffolds.
	// Enabling this wth no importer name is not accepted.
	ImporterRequired bool
	// Provider name of which the scaffolds to be created, defaults to terraform-provider-demo
	Provider string
	// Path defines where the templates has to be generated.
	Path string
	// AutoGenMessage will be configured by terragen and cannot be overwritten.
	AutoGenMessage string
	// Description to be added to resource/datasource.
	Description string
	// RepoGroup is used while creating go mod. Defaults to 'github.com/test/'
	// For a given provider, repo group would be appended.
	// Ex: For provider 'demo' the go mod would look like 'github.com/test/demo'
	// Resource to be created while generating scaffolds,
	// by passing a resource name here, it auto enabled ResourceRequired.
	// Provider name would be appended while constructing final resource name.
	// EX: resource 'create_cluster' for provider demo would become 'demo_create_cluster'.
	Resource []string
	// DataSource to be created while generating scaffolds,
	// by passing a resource name here, it auto enabled DatasourceRequired.
	// Provider name would be appended while constructing final data_source name.
	// EX: resource 'load_image' for provider demo would become 'demo_load_image'.
	DataSource []string
	// List of all the dependent packages for terraform, if not passed it picks default.
	Dependents []string
	// TemplateRaw consists of go-templates which are the core for terragen.
	TemplateRaw TerraTemplate
	// RepoGroup to which the Terraform provider is to be part. This would be used as base name to set `go.mod`
	// ex: github.com/nikshilsbhat
	RepoGroup    string
	mod          string
	metaDataPath string
	logger       *logrus.Logger
}

// TerraTemplate are the collections of go-templates which are used to generate terraform provider's base template.
type TerraTemplate struct {
	// ProviderTemp holds the template for provider
	ProviderTemp string `json:"provider-template,omitempty" yaml:"provider-template,omitempty"`
	// RootTemp holds the template for root file
	RootTemp string `json:"root-template,omitempty" yaml:"provider-template,omitempty"`
	// DataTemp holds the template for data
	DataTemp string `json:"data-template,omitempty" yaml:"data-template,omitempty"`
	// ResourceTemp holds the template for resource
	ResourceTemp string `json:"resource-template,omitempty" yaml:"resource-template,omitempty"`
	// GitIgnore that where scaffolded.
	GitIgnore string `json:"gitignore,omitempty" yaml:"gitignore,omitempty"`
	// GolangCILint that where scaffolded.
	GolangCILint string `json:"golang-ci-lint,omitempty" yaml:"golang-ci-lint,omitempty"`
	// GoReleaser that where scaffolded.
	GoReleaser string `json:"go-releaser,omitempty" yaml:"go-releaser,omitempty"`
	// RegistryManifest is a terraform registry manifest that is to be scaffolded.
	RegistryManifest string `json:"registry-manifest,omitempty" yaml:"registry-manifest,omitempty"`
}

var autoGenMessage = `// ----------------------------------------------------------------------------
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

func (i *Input) Generate(providerName string) error {
	i.setInputs()
	i.Provider = strings.ReplaceAll(providerName, "-", "_")
	i.Path = filepath.Join(i.getPath(), i.Provider)
	i.mod = i.setMod()

	if i.TerraformPluginFramework {
		return fmt.Errorf("plugin-framework is enabled, " +
			"scaffold would be generated as per https://github.com/hashicorp/terraform-provider-scaffolding-framework")
	}

	if !i.SkipValidation {
		if !i.validatePrerequisite() {
			i.logger.Error("system validation failed; please install prerequisites for Terraform provider scaffolds to work perfectly.")
		}
	}

	if i.Dependents == nil {
		i.Dependents = []string{"github.com/hashicorp/terraform-plugin-sdk/v2/plugin", filepath.Join(i.mod, "internal")}
	}

	i.logger.Debugf("dependent packages would be: %v", i.Dependents)
	i.logger.Debugf("go module for scaffold would be: %s", i.mod)

	if NewProvider(i).Scaffolded() {
		return fmt.Errorf("scaffolds for provider '%s' was already generated; "+
			"use `terragen create -h` or `terragen edit -h` for more info", i.Provider)
	}

	i.logger.Infof("terragen is in the process of making life simpler")

	if !i.DryRun {
		if err := i.genTerraDir(); err != nil {
			return fmt.Errorf("generating directories for scaffolds under %s failed with error: %w", i.Provider, err)
		}
	}

	if err := NewProvider(i).Create(); err != nil {
		return fmt.Errorf("creating scaffolds for terraform provider errored with '%s'", err.Error())
	}

	if err := NewMain(i).Create(); err != nil {
		return fmt.Errorf("oops generating scaffold main.go for provider %s failed with error: %w", i.Provider, err)
	}

	if i.ResourceRequired {
		if err := NewResource(i).Create(); err != nil {
			return fmt.Errorf("creating scaffolds for terraform resource errored with '%s'", err.Error())
		}
	}

	if i.DatasourceRequired {
		if err := NewDataSource(i).Create(); err != nil {
			return fmt.Errorf("creating scaffolds for terraform data source errored with '%s'", err.Error())
		}
	}

	if err := NewMake(i).Create(); err != nil {
		return fmt.Errorf("creating scaffolds makefile errored with '%s'", err.Error())
	}

	if err := NewGit(i).Create(); err != nil {
		return fmt.Errorf("creating scaffolds for 'gitignore' errored with '%s'", err.Error())
	}

	if err := NewReleaseNLinter(i).Create(); err != nil {
		return fmt.Errorf("creating scaffolds for 'goreleaser' and 'golangci-lint' errored with '%s'", err.Error())
	}

	if !i.DryRun {
		if err := i.setupTerragen(); err != nil {
			return fmt.Errorf("setting up scaffolds post scaffold generation errored with '%s'", err.Error())
		}
	}

	if err := i.CreateOrUpdateMetadata(); err != nil {
		return fmt.Errorf("oops creating/updating metadata errored out with %s", err.Error())
	}

	i.logger.Infof("yay! life is less complicated now!")
	i.logger.Infof("start enhancing terraform provider '%s' from the scaffold generated by terragen", i.Provider)

	return nil
}

func (i *Input) setInputs() {
	i.setRequires()
	i.getTemplate()
	i.enrichNames()
	i.AutoGenMessage = autoGenMessage
}

func (i *Input) SetLogger(logger *logrus.Logger) {
	i.logger = logger
}
