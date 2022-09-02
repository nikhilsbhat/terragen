package gen

import (
	// main.go template has to be sourced from template.
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
)

//go:embed templates/main.tmpl
var mainTemp string

type Main struct {
	DryRun         bool
	Path           string
	RootTemp       string
	Provider       string
	AutoGenMessage string
	Dependents     []string
}

func (m *Main) Create() error {
	mainFile := filepath.Join(m.Path, terragenMain)
	mainData, err := renderTemplate(terragenMain, m.RootTemp, m)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %w ", terragenMain, err)
	}

	if m.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terragenMain, m.Path)))
		log.Println(ui.Info("contents of main.go looks like"))
		printData(mainData)

		return nil
	}
	if err = terragenFileCreate(mainFile); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(m.Path, terragenMain), mainData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terragenMain, err)
	}

	return nil
}

func (m *Main) Scaffolded() bool {
	return false
}

func (m *Main) GetUpdated() error {
	return nil
}

func (m *Main) Update() error {
	return nil
}

func (m *Main) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewMain(i *Input) *Main {
	return &Main{
		DryRun:         i.DryRun,
		Path:           i.Path,
		RootTemp:       i.TemplateRaw.RootTemp,
		Provider:       i.Provider,
		AutoGenMessage: i.AutoGenMessage,
		Dependents:     i.Dependents,
	}
}
