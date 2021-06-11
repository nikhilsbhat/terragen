package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInput_enrichNames(t *testing.T) {
	input := Input{
		Resource:   []string{"hashicups_order"},
		DataSource: []string{"hashicups_ingredients", "hashicups_coffees"},
	}

	expected := Input{
		Resource:   []string{"resource_hashicups_order"},
		DataSource: []string{"datasource_hashicups_ingredients", "datasource_hashicups_coffees"},
	}

	input.enrichNames()
	assert.Equal(t, expected, input)
}

func TestInput_getUpdatedResourceNDataSources(t *testing.T) {
}

func TestInput_setMod(t *testing.T) {
	input := Input{
		RepoGroup: "github.com/nikhilsbhat",
		Provider:  "hashicups",
	}

	expected := "github.com/nikhilsbhat/terraform-provider-hashicups"
	actual := input.setMod()
	assert.Equal(t, expected, actual)
}

func TestInput_lockTerragenExecution(t *testing.T) {
	input := Input{
		Resource:   []string{"hashicups_order"},
		DataSource: []string{"hashicups_ingredients", "hashicups_coffees"},
	}

	_, _, lock, err := input.lockTerragenExecution()
	expected := true
	assert.NotNil(t, err)
	assert.Equal(t, expected, lock)
	//assert.NotNil(t, oldVer)
	//assert.NotNil(t, newVer)
}
