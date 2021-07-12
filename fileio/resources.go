package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)



// ResourceFilename returns filename relative to the resources directory.
// The resources directory is located at <config>/resources/
// On Linux this location is ~/.config/pigiron/resources/
//
// Returns non-nil error if resources directory can not be determined.
//
// Example:
// ResourceFilename("foo", "bar.txt") --> ~/.config/pigiron/resources/foo/bar.txt
//
func ResourceFilename(elements ...string) (string, error) {
	cfigdir, err := os.UserConfigDir()
	if err != nil {
		msg := "ERROR: Resource filename can not be determined.\n"
		msg += "ERROR: Can not determine configuration directory location.\n"
		msg += fmt.Sprintf("ERROR: %s\n", err)
		err = fmt.Errorf(msg)
		return "", err
	}
	acc := filepath.Join(cfigdir, "pigiron", "resources")
	for _, e := range elements {
		acc = filepath.Join(acc, e)
	}
	return acc, err
}
