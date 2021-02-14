// Package version powers the versioning of terragen.
package version

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	// Versn specifies the version of the application and cannot be changed by end user.
	Versn string

	// Env tells end user that what variant (here we use the name of the git branch to make it simple)
	// of application is he using.
	Env string
)

// GetVersion returns the version with variant, and the value will be used to bersion the application.
// The same version will be displayed at both CLI and app(API).
func GetVersion() string {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "v%s", Versn)
	if strings.ToLower(Env) != "production" {
		if Env == "" {
			Env = "alfa"
		}
		fmt.Fprintf(&versionString, "-%s", Env)
	}

	return versionString.String()
}
