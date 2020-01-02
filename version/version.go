// Package version powers the versioning of terragen.
package version

import (
	"bytes"
	"fmt"
)

// This specifies the version of the application and cannot be changed by end user.
const versn = "0.0.1"

// This tells end user that what variant (here we use the name of the git branch to make it simple)
// of application is he using.
var env = "alfa"

// GetVersion returns the version with variant, and the value will be used to bersion the application.
// The same version will be displayed at both CLI and app(API).
func GetVersion() string {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "v%s", versn)
	if env != "" {
		fmt.Fprintf(&versionString, "-%s", env)
	}

	return versionString.String()
}
