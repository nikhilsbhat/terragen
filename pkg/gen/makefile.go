package gen

import (
	// makefile template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

//go:embed templates/makefile.tmpl
var makefileTemplate string

type Make struct {
	DryRun   bool
	Path     string
	Provider string
	logger   *logrus.Logger
}

func (m *Make) Create() error {
	makeFile := filepath.Join(m.Path, terragenMakefile)
	makeFileData, err := renderTemplate(terragenMakefile, makefileTemplate, m)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %w ", terragenMakefile, err)
	}

	if m.DryRun {
		m.logger.Infof("Makefile would be created under %s", m.Path)
		m.logger.Infof("contents of Makefile source looks like")
		printData(makeFileData)

		return nil
	}

	if err = terragenFileCreate(makeFile); err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(m.Path, terragenMakefile), makeFileData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenMakefile, err)
	}

	return nil
}

func (m *Make) Scaffolded() bool {
	return false
}

func (m *Make) GetUpdated() error {
	return nil
}

func (m *Make) Update() error {
	return nil
}

//nolint:revive
func (m *Make) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewMake(i *Input) *Make {
	return &Make{
		DryRun:   i.DryRun,
		Provider: i.Provider,
		Path:     i.Path,
		logger:   i.logger,
	}
}
