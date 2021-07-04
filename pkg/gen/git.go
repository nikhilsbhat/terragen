package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
	"gopkg.in/src-d/go-git.v4"
)

var (
	gitignore = `
# dropping IDE data's
.vscode
.idea

.DS_Store

# dropping vendor and other data directories
vendor/

# dropping local built binaries
terraform-provider-demo
terraform-provider-{{ .Provider }}

# dropping terraform plans and states
terraform.tfstate*
*.tfplan
.terraform
*.log
`
)

// createGitIgnore scaffolds the provider and its components as per the requirements.
func (i *Input) createGitIgnore() error {
	mainFile := filepath.Join(i.Path, terrgenGitIgnore)
	gitIgnoreData, err := renderTemplate(terrgenGitIgnore, i.TemplateRaw.GitIgnore, i)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %v ", terrgenGitIgnore, err)
	}

	if i.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terrgenGitIgnore, i.Path)))
		log.Println(ui.Info("contents of gitignore looks like"))
		printData(gitIgnoreData)
	} else {
		if err = terragenFileCreate(mainFile); err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(i.Path, terrgenGitIgnore), gitIgnoreData, 0700); err != nil { //nolint:gosec
			return fmt.Errorf("oops scaffolding povider component %s errored with: %v ", terrgenGitIgnore, err)
		}
	}

	if err := i.initGit(); err != nil {
		return err
	}

	return nil
}

func (i *Input) initGit() error {
	_, err := git.PlainInit(i.Path, false)
	if err != nil {
		return err
	}
	return nil
}
