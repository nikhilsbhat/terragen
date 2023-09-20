package gen

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

func terragenFileCreate(path string) error {
	_, err := os.Create(path)

	return err
}

func (i *Input) setRequires() {
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
		return fmt.Errorf("running 'go mod init' errored with %w", err)
	}

	goFmt := exec.Command("goimports", "-w", i.Path) //nolint:gosec
	goFmt.Dir = i.Path
	if err := goFmt.Run(); err != nil {
		return fmt.Errorf("running 'goimports -w' errored with %w", err)
	}

	goTidy := exec.Command("go", "mod", "tidy")
	goTidy.Dir = i.Path
	if err := goTidy.Run(); err != nil {
		return fmt.Errorf("running 'go mod tidy' errored with %w", err)
	}

	goVnd := exec.Command("go", "mod", "vendor")
	goVnd.Dir = i.Path
	if err := goVnd.Run(); err != nil {
		return fmt.Errorf("running 'go mod vendor' errored with %w", err)
	}

	return nil
}

func (i *Input) genTerraDir() error {
	terraPath := filepath.Join(i.Path, "internal")
	err := os.MkdirAll(terraPath, scaffoldDirPerm)
	if err != nil {
		return err
	}

	return nil
}

// Set the templates to defaults if not specified.
func (i *Input) getTemplate() {
	if reflect.DeepEqual(i.TemplateRaw, TerraTemplate{}) {
		i.TemplateRaw.RootTemp = mainTemp
		i.TemplateRaw.GitIgnore = gitignore
		switch i.TerraformPluginFramework {
		case false:
			i.TemplateRaw.ProviderTemp = providerTemp
			i.TemplateRaw.DataTemp = dataSourceTemp
			i.TemplateRaw.ResourceTemp = resourceTemp
		case true:
			i.TemplateRaw.ProviderTemp = providerV2Temp
			i.TemplateRaw.DataTemp = dataSourceV2Temp
			i.TemplateRaw.ResourceTemp = resourceV2Temp
		}
	}
}

func (i *Input) getPath() string {
	if i.Path == "." {
		dir, err := os.Getwd()
		if err != nil {
			i.logger.Fatal(err.Error())
			os.Exit(1)
		}

		return dir
	}
	path, err := filepath.Abs(i.Path)
	if err != nil {
		i.logger.Fatal(err.Error())
		os.Exit(1)
	}

	return path
}

func renderTemplate(templateName, temp string, data interface{}) ([]byte, error) {
	var templateWriter bytes.Buffer
	tmpl := template.Must(template.New(templateName).Funcs(toCamel).Parse(temp))
	if err := tmpl.Execute(&templateWriter, data); err != nil {
		return nil, err
	}

	return templateWriter.Bytes(), nil
}

func (i *Input) validatePrerequisite() bool {
	success := true
	if goPath := exec.Command("go"); goPath.Err != nil {
		if !errors.Is(goPath.Err, exec.ErrDot) {
			i.logger.Error(goPath.Err.Error())
			i.logger.Error("terragen requires go to generate scaffolds")
			success = false
		}
	}

	if importsPath := exec.Command("goimports"); importsPath.Err != nil {
		if !errors.Is(importsPath.Err, exec.ErrDot) {
			i.logger.Error(importsPath.Err.Error())
			i.logger.Error("install goimports: go install goimports")
			success = false
		}
	}

	if fumptCmd := exec.Command("gofumpt"); fumptCmd.Err != nil {
		if !errors.Is(fumptCmd.Err, exec.ErrDot) {
			i.logger.Error(fumptCmd.Err.Error())
			i.logger.Error("install gofumpt: go install gofumpt")
			success = false
		}
	}

	if fmtCmd := exec.Command("gofmt"); fmtCmd.Err != nil {
		if !errors.Is(fmtCmd.Err, exec.ErrDot) {
			i.logger.Error(fmtCmd.Err.Error())
			i.logger.Error("install gofmt: go install gofmt")
			success = false
		}
	}

	if success {
		i.logger.Info("scaffolds would be generated with following golang version")
		out, err := exec.Command("go", "version").Output()
		if err != nil {
			i.logger.Error(err.Error())
		}
		i.logger.Info(string(bytes.TrimRight(out, "\n")))

		return success
	}

	return success
}
