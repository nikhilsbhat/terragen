package gen

import (
	_ "embed"
	"fmt"
	"github.com/nikhilsbhat/neuron/cli/ui"
	"log"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/makefile.tmpl
	makefileTemplate string
)

type Make struct {
	DryRun   bool
	Path     string
	Provider string
}

func (m *Make) Create() error {
	makeFile := filepath.Join(m.Path, terragenMakefile)
	makeFileData, err := renderTemplate(terragenMakefile, makefileTemplate, m)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %v ", terragenMakefile, err)
	}

	if m.DryRun {
		log.Print(ui.Info(fmt.Sprintf("Makefile would be created under %s", m.Path)))
		log.Println(ui.Info("contents of Makefile source looks like"))
		printData(makeFileData)
		return nil
	} else {
		if err = terragenFileCreate(makeFile); err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(m.Path, terragenMakefile), makeFileData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding povider component %s errored with: %v ", terragenMakefile, err)
		}
	}

	return nil
}

func NewMake(i *Input) *Make {
	return &Make{
		DryRun:   i.DryRun,
		Provider: i.Provider,
		Path:     i.Path,
	}
}
