package gen

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

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
	var fileWriter io.Writer
	if i.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terragenMain, i.Path)))
		fmt.Println(ui.Info("contents of main.go looks like"))
		fileWriter = os.Stdout
	} else {
		file, err := terragenWriter(i.Path, terragenMain)
		if err != nil {
			return err
		}
		defer file.Close()
		fileWriter = file
	}

	if len(i.TemplateRaw.RootTemp) != 0 {
		tmpl := template.Must(template.New(terragenMain).Parse(i.TemplateRaw.RootTemp))
		if err := tmpl.Execute(fileWriter, i); err != nil {
			return err
		}
	}
	return nil
}
