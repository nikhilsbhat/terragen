package gen

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/neuron/cli/ui"
	"gopkg.in/src-d/go-git.v4"
)

//go:embed templates/git.tmpl
var gitignore string

type Git struct {
	DryRun    bool
	Provider  string
	Path      string
	GitIgnore string
}

// Create scaffolds Git as per the requirements.
func (g *Git) Create() error {
	mainFile := filepath.Join(g.Path, terrgenGitIgnore)
	gitIgnoreData, err := renderTemplate(terrgenGitIgnore, g.GitIgnore, g)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %v ", terrgenGitIgnore, err)
	}

	if g.DryRun {
		log.Print(ui.Info(fmt.Sprintf("%s would be created under %s", terrgenGitIgnore, g.Path)))
		log.Println(ui.Info("contents of gitignore looks like"))
		printData(gitIgnoreData)
		return nil
	} else {
		if err = terragenFileCreate(mainFile); err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(g.Path, terrgenGitIgnore), gitIgnoreData, scaffoldPerm); err != nil {
			return fmt.Errorf("oops scaffolding povider component %s errored with: %v ", terrgenGitIgnore, err)
		}
	}

	_, err = git.PlainInit(g.Path, false)
	if err != nil {
		return err
	}

	return nil
}

func (g *Git) Scaffolded() bool {
	return false
}

func NewGit(i *Input) *Git {
	return &Git{
		DryRun:    i.DryRun,
		Provider:  i.Provider,
		Path:      i.Path,
		GitIgnore: i.TemplateRaw.GitIgnore,
	}
}
