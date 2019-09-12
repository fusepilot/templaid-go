package templaid

import (
	"os"
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

func getTreeMap(fs afero.Fs, path string) map[string]string {
	paths := map[string]string{}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		bytes, _ := afero.ReadFile(fs, path)
		paths[path] = string(bytes)
		return nil
	})

	return paths
}

func getFileTreeMap(fs afero.Fs, path string) map[string]string {
	paths := map[string]string{}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		bytes, _ := afero.ReadFile(fs, path)
		paths[path] = string(bytes)
		return nil
	})

	return paths
}
