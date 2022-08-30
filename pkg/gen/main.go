package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
)

var mainTemp = `{{ .AutoGenMessage }}
package main

import (
	{{- range $index, $element := .Dependents }}
	"{{- $element }}"
	{{- end }}
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: {{ .Provider }}.Provider})
}`

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
		return fmt.Errorf("oops rendering povider component %s errored with: %v ", terragenMain, err)
	}

	if m.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terragenMain, m.Path)))
		log.Println(ui.Info("contents of main.go looks like"))
		printData(mainData)
	} else {
		if err = terragenFileCreate(mainFile); err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(m.Path, terragenMain), mainData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding povider component %s errored with: %v ", terragenMain, err)
		}
	}
	return nil
}

func (m *Main) Scaffolded() bool {
	return false
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
