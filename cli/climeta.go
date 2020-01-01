package cli

import (
	"github.com/nikhilsbhat/neuron/cli/ui"
	"os"
)

type cliMeta struct {
	*ui.NeuronUi
}

var (
	cm = &cliMeta{}
)

func init() {

	nui := ui.NeuronUi{&ui.UiWriter{os.Stdout}}
	cm = &cliMeta{&nui}

}
