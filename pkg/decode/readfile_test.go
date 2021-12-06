package decode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	t.Run("should be able to read the file and return the contents of it", func(t *testing.T) {
		expected := `version: "1.0"
repo-group: github.com/nikhilsbhat
project-module: github.com/nikhilsbhat/terraform-provider-hashicups
provider: hashicups
provider-path: path/to/terraform-provider-test
resources:
  - resource_hashicups_order
data-sources:
  - datasource_hashicups_coffees
  - datasource_hashicups_ingredients
  - datasource_hashicups_order
importers:
  - ""`

		path, err := getFixturesPath()
		assert.NoError(t, err)
		actual, err := ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("should return file not found error", func(t *testing.T) {
		path, err := ReadFile("gen/fixtures/terragen-test.yml")
		assert.EqualError(t, err, "stat gen/fixtures/terragen-test.yml: no such file or directory")
		assert.Nil(t, path)
	})
}

func getFixturesPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(dir), "gen/fixtures/terragen-test.yml"), nil
}
