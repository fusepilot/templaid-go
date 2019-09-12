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

const TemplateSuffix = ".template"

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

type RenderFolderProps struct {
	Fs              afero.Fs
	TemplatePath    string
	FolderPattern   string
	DestinationPath string
	IgnorePatterns  []string
	Data            map[string]string
}

func RenderFilesFolders(props RenderFolderProps) {
	afero.Walk(props.Fs, props.TemplatePath, func(path string, info os.FileInfo, err error) error {
		name := info.Name()

		for _, ignorePattern := range props.IgnorePatterns {
			if name == ignorePattern {
				return nil
			} else if match, err := filepath.Match(ignorePattern, path); match == true || err != nil {
				return nil
			}
		}

		if info.IsDir() {
			newFolderPath := GetDestinationFilePath(props.TemplatePath, props.DestinationPath, path, props.Data)
			props.Fs.MkdirAll(newFolderPath, os.ModePerm)
		} else {
			newFilePath := GetDestinationFilePath(props.TemplatePath, props.DestinationPath, path, props.Data)
			newFile, err := props.Fs.Create(newFilePath)
			if err != nil {
				return err
			}
			defer newFile.Close()

			bytes, err := afero.ReadFile(props.Fs, path)
			if err != nil {
				return err
			}

			if strings.HasSuffix(name, TemplateSuffix) {
				templateBytes, err := RenderTemplateString(string(bytes), props.Data)
				if err != nil {
					return err
				}
				newFile.WriteString(templateBytes)
			} else {
				newFile.Write(bytes)
			}

		}
		return nil
	})
}

type RenderTemplateProps struct {
	Fs              afero.Fs
	TemplatePath    string
	DestinationPath string
	IgnorePatterns  []string
	Data            map[string]string
	TemplatePattern string
}

func RenderTemplate(props RenderTemplateProps) error {
	if !checkPathExists(props.Fs, props.TemplatePath) {
		return fmt.Errorf(`template path "%s" doesnt not exist`, props.TemplatePath)
	}

	RenderFilesFolders(RenderFolderProps{
		Fs:              props.Fs,
		TemplatePath:    props.TemplatePath,
		DestinationPath: props.DestinationPath,
		Data:            props.Data,
		IgnorePatterns:  props.IgnorePatterns,
	})

	return nil
}
