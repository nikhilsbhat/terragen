package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/nikhilsbhat/neuron/cli/ui"
)

var providerTemp = `{{ .AutoGenMessage }}
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

type Provider struct {
	Provider       string
	DryRun         bool
	Path           string
	Mod            string
	ProviderTemp   string
	MetaDataPath   string
	AutoGenMessage string
	Resource       []string
	DataSource     []string
	Importer       []string
}

func (p *Provider) Create() error {
	provideFile := filepath.Join(p.Path, p.Provider, terragenProvider)
	providerData, err := renderTemplate(terragenProvider, p.ProviderTemp, p)
	if err != nil {
		log.Fatalf(ui.Error(fmt.Sprintf("oops rendering provider %s errored with: %v ", p.Provider, err)))
	}

	if p.DryRun {
		log.Println(ui.Info(fmt.Sprintf("provider '%s' would be created under '%s'", p.Provider, p.Path)))
		log.Println(ui.Info("contents of provider looks like"))
		printData(providerData)
	} else {
		if err = ioutil.WriteFile(provideFile, providerData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding provider %s errored with: %v ", p.Provider, err)
		}
	}

	return nil
}

func (p *Provider) Update() error {
	currentProvider, err := ioutil.ReadFile(filepath.Join(p.Path, p.Provider, terragenProvider))
	if err != nil {
		return err
	}

	newIn := &Provider{}
	if err = copier.CopyWithOption(newIn, p, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	if err = newIn.GetUpdated(); err != nil {
		return err
	}

	updateData, err := newIn.Get(currentProvider)
	if err != nil {
		return err
	}

	providerFile := filepath.Join(newIn.Path, newIn.Provider, terragenProvider)
	if err = ioutil.WriteFile(providerFile, updateData, scaffoldPerm); err != nil {
		return err
	}
	return nil
}

func (p *Provider) Get(currentProvider []byte) ([]byte, error) {
	updatedProvider, err := renderTemplate(terragenProvider, p.ProviderTemp, p)
	if err != nil {
		return nil, err
	}

	dmp := diffmatchpatch.New()
	providerDiff := dmp.DiffMain(string(currentProvider), string(updatedProvider), false)
	return []byte(dmp.DiffText2(providerDiff)), nil
}

func (p *Provider) GetUpdated() error {
	metadata, err := getCurrentMetadata(filepath.Join(p.Path, terragenMetadata))
	if err != nil {
		return err
	}
	p.DataSource = append(p.DataSource, metadata.DataSources...)
	p.Resource = append(p.Resource, metadata.Resources...)
	return nil
}

func (p *Provider) Scaffolded() bool {
	if _, dirErr := os.Stat(filepath.Join(p.Path, terragenMetadata)); os.IsNotExist(dirErr) {
		return false
	}
	metadata, err := getCurrentMetadata(filepath.Join(p.Path, terragenMetadata))
	if err != nil {
		log.Println(ui.Error(err.Error()))
		return true
	}
	return p.Provider == metadata.Provider
}

func NewProvider(i *Input) *Provider {
	return &Provider{
		Provider:     i.Provider,
		DryRun:       i.DryRun,
		Path:         i.Path,
		Mod:          i.mod,
		ProviderTemp: i.TemplateRaw.ProviderTemp,
		Resource:     i.Resource,
		DataSource:   i.DataSource,
	}
}
