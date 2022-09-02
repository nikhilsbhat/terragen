package gen //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/nikhilsbhat/terragen/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
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
	t.Run("Should lock the execution of terragen due to version issues", func(t *testing.T) {
		oldTerragenVersion := "1.0.0"
		newTerragenVersion := "0.2.0"
		version.Version = newTerragenVersion
		metadataFromFile, err := mockGetMetadata("fixtures/terragen-test.yml")
		require.NoError(t, err)

		oldVersion, newVersion, lock, err := lockTerragenExecution(metadataFromFile.Version, false)
		require.NoError(t, err)
		assert.Equal(t, true, lock)
		assert.Equal(t, newTerragenVersion, newVersion)
		assert.Equal(t, oldTerragenVersion, oldVersion)
	})

	t.Run("Should not lock the execution of terragen", func(t *testing.T) {
		oldTerragenVersion := "1.0.0"
		newTerragenVersion := "1.2.0"
		version.Version = newTerragenVersion
		metadataFromFile, err := mockGetMetadata("fixtures/terragen-test.yml")
		require.NoError(t, err)

		oldVersion, newVersion, lock, err := lockTerragenExecution(metadataFromFile.Version, false)
		require.NoError(t, err)
		assert.Equal(t, false, lock)
		assert.Equal(t, newTerragenVersion, newVersion)
		assert.Equal(t, oldTerragenVersion, oldVersion)
	})
}

func Test_snakeCaseToCamelCase(t *testing.T) {
	t.Run("should be able to convert given snakeCase string to camelCase", func(t *testing.T) {
		snakeCaseString := "this_is_test"
		expected := "thisIsTest"

		actual := snakeCaseToCamelCase(snakeCaseString)
		assert.Equal(t, expected, actual)
	})
}

func mockGetMetadata(path string) (*Config, error) {
	metaData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	metadata := newMetadata()
	if yamlErr := yaml.Unmarshal(metaData, &metadata); yamlErr != nil {
		return nil, yamlErr
	}

	return metadata, nil
}
