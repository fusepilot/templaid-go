package templaid

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type GetDestinationFilePathProps struct {
	TemplatePath    string
	DestinationPath string
	File            string
	Data            map[string]string
}

func GetDestinationFilePath(props GetDestinationFilePathProps) string {
	renderedPath := GetRenderedPath(props.File, props.TemplatePath, props.DestinationPath)

	result, err := GetTemplatedFilePath(renderedPath, props.Data)
	if err != nil {
		panic(err)
	}
	return result
}

type RenderTemplateProps struct {
	Fs              afero.Fs
	TemplatePath    string
	DestinationPath string
	IgnorePattern   string
	Data            map[string]string
	TemplatePattern string
}

func checkPathExists(fs afero.Fs, path string) bool {
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return false
	} else {
		return exists
	}
}

func RenderTemplate(props RenderTemplateProps) error {
	if !checkPathExists(props.Fs, props.TemplatePath) {
		return fmt.Errorf(`template path "%s" doesnt not exist`, props.TemplatePath)
	}
	fullPath := filepath.Join(props.DestinationPath, "a/c", "b.test")
	dirPath := filepath.Dir(fullPath)

	props.Fs.MkdirAll(dirPath, os.ModePerm)
	err := afero.WriteFile(props.Fs, filepath.Join(props.DestinationPath, "a/c", "b.test"), []byte("file b"), 0644)
	if err != nil {
		panic(err)
	}
	return nil
}
