package cmd

import (
	"os"

	"github.com/nikhilsbhat/neuron/cli/ui"
)

type cliMeta struct {
	*ui.NeuronUi
}

var cm = &cliMeta{}

func init() {
	nui := ui.NeuronUi{
		UiWriter: &ui.UiWriter{
			Writer: os.Stdout,
		},
	}
	cm = &cliMeta{NeuronUi: &nui}
}
