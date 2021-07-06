package gen

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"

	"github.com/nikhilsbhat/neuron/cli/ui"
	"github.com/nikhilsbhat/terragen/pkg/decode"
)

func newMetadata() *Metadata {
	return &Metadata{}
}

func newInput() *Input {
	return &Input{}
}

func terragenFileCreate(path string) error {
	_, err := os.Create(path)
	return err
}

func (i *Input) setDefaults() {
	if len(i.Resource) != 0 {
		i.ResourceRequired = true
	}
	if len(i.DataSource) != 0 {
		i.DatasourceRequired = true
	}
	if len(i.Importer) != 0 {
		i.ImporterRequired = true
	}
}

func (i *Input) setupTerragen() error {
	goInit := exec.Command("go", "mod", "init", i.mod) //nolint:gosec
	goInit.Dir = i.Path
	if err := goInit.Run(); err != nil {
		return err
	}

	goFmt := exec.Command("goimports", "-w", i.Path) //nolint:gosec
	goFmt.Dir = i.Path
	if err := goFmt.Run(); err != nil {
		return err
	}

	goVnd := exec.Command("go", "mod", "vendor") //nolint:gosec
	goVnd.Dir = i.Path
	if err := goVnd.Run(); err != nil {
		return err
	}
	return nil
}

func (i *Input) genTerraDir() error {
	err := os.MkdirAll(filepath.Join(i.Path, i.Provider), 0777)
	if err != nil {
		return err
	}
	return nil
}

// Set the templates to defaults if not specified.
func (i *Input) getTemplate() {
	if reflect.DeepEqual(i.TemplateRaw, TerraTemplate{}) {
		i.TemplateRaw.RootTemp = mainTemp
		i.TemplateRaw.ProviderTemp = providerTemp
		i.TemplateRaw.DataTemp = dataSourceTemp
		i.TemplateRaw.ResourceTemp = resourceTemp
		i.TemplateRaw.GitIgnore = gitignore
	}
}

func (i *Input) getPath() string {
	if i.Path == "." {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(ui.Error(decode.GetStringOfMessage(err)))
			os.Exit(1)
		}
		return dir
	}
	return path.Dir(i.Path)
}

func renderTemplate(templateName, temp string, data interface{}) ([]byte, error) {
	var templateWriter bytes.Buffer
	tmpl := template.Must(template.New(templateName).Funcs(toCamel).Parse(temp))
	if err := tmpl.Execute(&templateWriter, data); err != nil {
		return nil, err
	}
	return templateWriter.Bytes(), nil
}
