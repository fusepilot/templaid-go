package templaid

import (
	"regexp"

	"github.com/spf13/afero"
)

func stringContainsToken(s string) bool {
	contains, err := regexp.Match(`{{[^}]*}}`, []byte(s))
	if err != nil {
		return false
	}
	return contains
}

func checkPathExists(fs afero.Fs, path string) bool {
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return false
	} else {
		return exists
	}
}
