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
)

type ReleaseNLinter struct {
	DryRun       bool
	Provider     string
	Path         string
	GolangCILint string
	GoReleaser   string
	logger       *logrus.Logger
}

func (r *ReleaseNLinter) Create() error {
	releaserFile := filepath.Join(r.Path, terragenReleaser)
	ciLinterFile := filepath.Join(r.Path, terragenCILinter)

	releaserData := []byte(r.GoReleaser)
	ciLinterFileData := []byte(r.GolangCILint)

	if r.DryRun {
		r.logger.Infof("%s would be created under %s", terragenReleaser, r.Path)
		r.logger.Infof("contents of %s looks like", terragenReleaser)
		printData(releaserData)

		r.logger.Infof("%s would be created under %s", terragenCILinter, r.Path)
		r.logger.Infof("contents of %s looks like", terragenCILinter)
		printData(ciLinterFileData)

		return nil
	}

	if err := terragenFileCreate(releaserFile); err != nil {
		return err
	}

	if err := terragenFileCreate(ciLinterFile); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(r.Path, terragenReleaser), releaserData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenReleaser, err)
	}

	if err := os.WriteFile(filepath.Join(r.Path, terragenCILinter), ciLinterFileData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenCILinter, err)
	}

	return nil
}

func (r *ReleaseNLinter) Scaffolded() bool {
	return false
}

func (r *ReleaseNLinter) GetUpdated() error {
	return nil
}

func (r *ReleaseNLinter) Update() error {
	return nil
}

//nolint:revive
func (r *ReleaseNLinter) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewReleaseNLinter(i *Input) *ReleaseNLinter {
	return &ReleaseNLinter{
		DryRun:       i.DryRun,
		Path:         i.Path,
		Provider:     i.Provider,
		GoReleaser:   i.TemplateRaw.GoReleaser,
		GolangCILint: i.TemplateRaw.GolangCILint,
		logger:       i.logger,
	}
}
