package gen

import (
	// git template has to be sourced from template.
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

//go:embed templates/git.tmpl
var gitignore string

type Git struct {
	DryRun    bool
	Provider  string
	Path      string
	GitIgnore string
	logger    *logrus.Logger
}

// Create scaffolds Git as per the requirements.
func (g *Git) Create() error {
	mainFile := filepath.Join(g.Path, terrgenGitIgnore)
	gitIgnoreData, err := renderTemplate(terrgenGitIgnore, g.GitIgnore, g)
	if err != nil {
		return fmt.Errorf("oops rendering povider component %s errored with: %w ", terrgenGitIgnore, err)
	}

	if g.DryRun {
		g.logger.Infof("%s would be created under %s", terrgenGitIgnore, g.Path)
		g.logger.Infof("contents of gitignore looks like")
		printData(gitIgnoreData)

		return nil
	}

	if err = terragenFileCreate(mainFile); err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(g.Path, terrgenGitIgnore), gitIgnoreData, scaffoldPerm); err != nil {
		return fmt.Errorf("oops scaffolding povider component %s errored with: %w ", terrgenGitIgnore, err)
	}

	if _, err = git.PlainInit(g.Path, false); err != nil {
		return err
	}

	return nil
}

func (g *Git) Scaffolded() bool {
	return false
}

func (g *Git) GetUpdated() error {
	return nil
}

func (g *Git) Update() error {
	return nil
}

//nolint:revive
func (g *Git) Get(currentContent []byte) ([]byte, error) {
	return nil, nil
}

func NewGit(i *Input) *Git {
	return &Git{
		DryRun:    i.DryRun,
		Provider:  i.Provider,
		Path:      i.Path,
		GitIgnore: i.TemplateRaw.GitIgnore,
		logger:    i.logger,
	}
}
