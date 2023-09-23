package gen

import (
	// provider template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed templates/provider.tmpl
	providerTemp string
	//go:embed templates/provider_v2.tmpl
	providerV2Temp string
)

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
	logger         *logrus.Logger
}

func (p *Provider) Create() error {
	provideFile := filepath.Join(p.Path, "internal", terragenProvider)
	providerData, err := renderTemplate(terragenProvider, p.ProviderTemp, p)
	if err != nil {
		p.logger.Fatalf("oops rendering provider %s errored with: %v ", p.Provider, err)
	}

	p.logger.Infof("scaffolds for provider '%s' would be generated under: '%s'", p.Provider, p.Path)

	if p.DryRun {
		p.logger.Infof("provider '%s' would be created under '%s'", p.Provider, p.Path)
		p.logger.Infof("contents of provider looks like")
		printData(providerData)

		return nil
	}

	if err = os.WriteFile(provideFile, providerData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding provider %s errored with: %w ", p.Provider, err)
	}

	return nil
}

func (p *Provider) Update() error {
	currentProvider, err := os.ReadFile(filepath.Join(p.Path, "internal", terragenProvider))
	if err != nil {
		p.logger.Errorf("reading existing provider errored with '%s'", err.Error())

		return err
	}

	newIn := &Provider{}
	if err = copier.CopyWithOption(newIn, p, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	if err = newIn.GetUpdated(); err != nil {
		p.logger.Errorf("fetching updated provider config errored with '%s'", err.Error())

		return err
	}

	updateData, err := newIn.Get(currentProvider)
	if err != nil {
		p.logger.Errorf("getting updated provider config errored with '%s'", err.Error())

		return err
	}

	providerFile := filepath.Join(newIn.Path, "internal", terragenProvider)
	if err = os.WriteFile(providerFile, updateData, scaffoldPerm); err != nil {
		return fmt.Errorf("writing update provider config errored with '%s'", err.Error())
	}

	return nil
}

func (p *Provider) Get(currentContent []byte) ([]byte, error) {
	updatedProvider, err := renderTemplate(terragenProvider, p.ProviderTemp, p)
	if err != nil {
		return nil, err
	}

	dmp := diffmatchpatch.New()
	providerDiff := dmp.DiffMain(string(currentContent), string(updatedProvider), false)

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
	metadatataPath := filepath.Join(p.Path, terragenMetadata)
	if _, dirErr := os.Stat(metadatataPath); os.IsNotExist(dirErr) {
		p.logger.Warnf("unable to find the terragen.yml from '%s', ignore if this is first run", metadatataPath)

		return false
	}
	metadata, err := getCurrentMetadata(filepath.Join(p.Path, terragenMetadata))
	if err != nil {
		p.logger.Error(err.Error())

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
		logger:       i.logger,
	}
}
