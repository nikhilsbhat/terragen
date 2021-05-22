package gen

import (
	"html/template"
	"io"
	"log"
	"os"

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
	var fileWriter io.Writer
	if i.DryRun {
		log.Println(ui.Info("contents of gitignore looks like"))
		fileWriter = os.Stdout
	} else {
		file, err := terragenWriter(i.Path, terrgenGitIgnore)
		if err != nil {
			return err
		}
		defer file.Close()
		fileWriter = file
	}

	if len(i.TemplateRaw.DataTemp) != 0 {
		tmpl := template.Must(template.New(terrgenGitIgnore).Parse(i.TemplateRaw.GitIgnore))
		if err := tmpl.Execute(fileWriter, i); err != nil {
			return err
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
