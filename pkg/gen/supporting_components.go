package gen

import (
	// main.go template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	//go:embed templates/golangci-lint.yml
	golangCILint string
	//go:embed templates/goreleaser.yml
	goReleaser string
	//go:embed templates/terraform-registry-manifest.json
	registryManifest string
)

type SupportingComponents struct {
	DryRun           bool
	Provider         string
	Path             string
	GolangCILint     string
	GoReleaser       string
	RegistryManifest string
	logger           *logrus.Logger
}

func (r *SupportingComponents) Create() error {
	releaserFile := filepath.Join(r.Path, terragenReleaser)
	ciLinterFile := filepath.Join(r.Path, terragenCILinter)
	registryManifestFile := filepath.Join(r.Path, terragenManifestFile)

	releaserData := []byte(r.GoReleaser)
	ciLinterFileData := []byte(r.GolangCILint)
	registryManifestFileData := []byte(r.RegistryManifest)

	if r.DryRun {
		r.logger.Infof("%s would be created under %s", terragenReleaser, r.Path)
		r.logger.Infof("contents of %s looks like", terragenReleaser)
		printData(releaserData)

		r.logger.Infof("%s would be created under %s", terragenCILinter, r.Path)
		r.logger.Infof("contents of %s looks like", terragenCILinter)
		printData(ciLinterFileData)

		r.logger.Infof("%s would be created under %s", terragenManifestFile, r.Path)
		r.logger.Infof("contents of %s looks like", terragenManifestFile)
		printData(registryManifestFileData)

		return nil
	}

	if err := terragenFileCreate(releaserFile); err != nil {
		return err
	}

	if err := terragenFileCreate(ciLinterFile); err != nil {
		return err
	}

	if err := terragenFileCreate(registryManifestFile); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(r.Path, terragenReleaser), releaserData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenReleaser, err)
	}

	if err := os.WriteFile(filepath.Join(r.Path, terragenCILinter), ciLinterFileData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenCILinter, err)
	}

	if err := os.WriteFile(filepath.Join(r.Path, terragenManifestFile), registryManifestFileData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenManifestFile, err)
	}

	return nil
}

func (r *SupportingComponents) Scaffolded() bool {
	return false
}

func (r *SupportingComponents) GetUpdated() error {
	return nil
}

func (r *SupportingComponents) Update() error {
	return nil
}

//nolint:revive
func (r *SupportingComponents) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewReleaseNLinter(i *Input) *SupportingComponents {
	return &SupportingComponents{
		DryRun:           i.DryRun,
		Path:             i.Path,
		Provider:         i.Provider,
		GoReleaser:       i.TemplateRaw.GoReleaser,
		GolangCILint:     i.TemplateRaw.GolangCILint,
		RegistryManifest: i.TemplateRaw.RegistryManifest,
		logger:           i.logger,
	}
}
