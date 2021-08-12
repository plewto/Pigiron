/* 
** pigpath package provides utilities for for making filename substitutions.
** Specifically leading characters '~' and '!' are replaced with the user's
** home and pigiron configuration directories respectively.
**
*/

package pigpath

import (
	"os"
	"path/filepath"
)

func UserHomeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

func PigironConfigDir() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "pigiron")
}


// ResourceFilename returns filename relative to the resources directory.
// The resources directory is located at <config>/resources/
// On Linux this location is ~/.config/pigiron/resources/
//
// Example:
// ResourceFilename("foo", "bar.txt") --> ~/.config/pigiron/resources/foo/bar.txt
//
func ResourceFilename(elements ...string) string {
	acc := filepath.Join(PigironConfigDir(), "resources")
	for _, e := range elements {
		acc = filepath.Join(acc, e)
	}
	return acc
}


// SubSpecialDirectories replaces leading ~ and ! characters in filename
// '~' is replaced with user's home dir.
// '!' is replaced with pigiron configuration dir.
//
func SubSpecialDirectories (filename string) string {
	result := filename
	if len(filename) > 0 {
		switch {
		case filename[0] == '~':
			home := UserHomeDir()
			result = filepath.Join(home, filename[1:])
		case filename[0] == '!':
			cnfig := PigironConfigDir()
			result = filepath.Join(cnfig, filename[1:])
		default:
			// ignore
		}
	}
	return result
}


// Join concatenates filepath components into string.
// It is like the standard filepath.Join but substitutes leading
// ~ and ! characters with user's home and pigiron configuration
// directories respectively.
// 
func Join(base string, elements ...string) string {
	var acc = make([]string,1,12)
	acc[0] = SubSpecialDirectories(base)
	for _, e := range elements {
		acc = append(acc, e)
	}
	return filepath.Join(acc...)
}
	
