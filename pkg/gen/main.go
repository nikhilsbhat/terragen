package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
)

var (
	mainTemp = `{{ .AutoGenMessage }}
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
)

func (i *Input) CreateMain() error {
	mainFile := filepath.Join(i.Path, terragenMain)
	mainData, err := renderTemplate(terragenMain, i.TemplateRaw.RootTemp, i)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %v ", terragenMain, err)
	}

	if i.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terragenMain, i.Path)))
		log.Println(ui.Info("contents of main.go looks like"))
		printData(mainData)
	} else {
		if err = terragenFileCreate(mainFile); err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(i.Path, terragenMain), mainData, 0700); err != nil { //nolint:gosec
			return fmt.Errorf("oops scaffolding povider component %s errored with: %v ", terragenMain, err)
		}
	}
	return nil
}
