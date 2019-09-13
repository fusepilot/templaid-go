package templaid

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func writeFs(fs afero.Fs, files map[string]string) {
	for fileName, fileContents := range files {
		if strings.HasSuffix(fileName, "/") {
			fs.MkdirAll(fileName, os.ModePerm)
		} else {
			file, err := fs.Create(fileName)
			defer file.Close()
			if err != nil {
				panic(err)
			}
			file.WriteString(fileContents)
		}
	}
}

func copyFs(sourcePath string, sourceFS afero.Fs, targetFS afero.Fs) {
	afero.Walk(sourceFS, sourcePath, func(path string, info os.FileInfo, err error) error {
		relPath, _ := filepath.Rel(sourcePath, path)
		if info.IsDir() {
			err := targetFS.MkdirAll(relPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			data, err := afero.ReadFile(sourceFS, path)
			if err != nil {
				return err
			}
			file, err := targetFS.Create(relPath)
			if err != nil {
				return err
			}
			_, err = file.Write(data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func assertFsEqual(t *testing.T, fs afero.Fs, root string, fileMap map[string]string) {
	assert.Equal(t, fileMap, getTreeMap(fs, root))
}
