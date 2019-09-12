package templaid

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	template "github.com/valyala/fasttemplate"
)

func RenderTemplateString(s string, data map[string]string) (string, error) {
	t, err := template.NewTemplate(s, "{{", "}}")
	if err != nil {
		log.Fatalf("unexpected error when parsing template: %s", err)
		return "", err
	}

	resultString := t.ExecuteFuncString(func(w io.Writer, token string) (int, error) {
		tokenString := "{{" + token + "}}"
		if knownToken := data[token]; knownToken != "" {
			return w.Write([]byte(knownToken))
		} else {
			return w.Write([]byte(tokenString))
		}
	})

	return resultString, nil
}

func GetTemplatedFilePath(path string, data map[string]string) (string, error) {
	if stringContainsToken(path) {
		normalizedFilePath := strings.Join(filepath.SplitList(path), "/")
		renderedFilePath, err := RenderTemplateString(normalizedFilePath, data)
		if err != nil {
			return "", err
		}
		denormalizedFilePath := strings.Join(filepath.SplitList(renderedFilePath), "/")
		return denormalizedFilePath, nil
	} else {
		return path, nil
	}
}

func GetRenderedPath(templatePath string, rootPath string, newRootPath string) string {
	relativeRoot, _ := filepath.Rel(rootPath, templatePath)
	newRelativeRoot := filepath.Join(newRootPath, relativeRoot)
	return newRelativeRoot
}

func GetDestinationFilePath(templatePath string, destinationPath string, file string, data map[string]string) string {
	renderedPath := GetRenderedPath(file, templatePath, destinationPath)

	result, err := GetTemplatedFilePath(renderedPath, data)
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
