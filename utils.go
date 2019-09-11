package templaid

import (
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"

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
		if knownTag := data[token]; knownTag != "" {
			return w.Write([]byte(knownTag))
		} else {
			return w.Write([]byte(tokenString))
		}
	})

	return resultString, nil
}

func StringContainsToken(s string) bool {
	contains, err := regexp.Match(`{{[^}]*}}`, []byte(s))
	if err != nil {
		return false
	}
	return contains
}

func GetTemplatedFilePath(path string, data map[string]string) (string, error) {
	if StringContainsToken(path) {
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
