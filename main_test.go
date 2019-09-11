package templaid

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func CopyFs(sourcePath string, sourceFS afero.Fs, targetFS afero.Fs) {
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

func GetTree(fs afero.Fs) string {
	paths := []string{}
	afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
		// bytes, _ := afero.ReadFile(memFs, path)
		// paths = append(paths, fmt.Sprintf("%s : %s", path, bytes))

		paths = append(paths, path)
		return nil
	})

	result := strings.Join(paths, "\n") + "\n"
	return result
}

func TestGetDestinationFilePath(t *testing.T) {
	result := GetDestinationFilePath(GetDestinationFilePathProps{
		TemplatePath:    "/templates/complex",
		DestinationPath: "/output/result",
		File:            "/templates/complex/{{template.name}}/{{template.name}}-a.md",
		Data:            map[string]string{"template.name": "NewProject"}},
	)
	assert.Equal(t, result, "/output/result/NewProject/NewProject-a.md")
}

func TestRenderTemplate(t *testing.T) {
	osFs := afero.NewOsFs()
	memFs := afero.NewMemMapFs()

	// memFs.MkdirAll("complex", os.ModePerm)
	// memFs.MkdirAll("complex/folder-b", os.ModePerm)
	// afero.WriteFile(memFs, "complex/{{template.name}}-b.md", []byte("template contents"), os.ModePerm)

	CopyFs("/Users/michael/Workspace/go/src/github.com/fusepilot/templaid/examples", osFs, memFs)

	// templatePath, _ := filepath.Abs("./examples/complex")
	// outputPath, _ := filepath.Abs("./output")

	RenderTemplate(RenderTemplateProps{
		Fs: memFs,
		// TemplatePath:    templatePath,
		// DestinationPath: outputPath,
		TemplatePath:    "complex",
		DestinationPath: "output",
		Data:            map[string]string{"template.name": "NewProject"},
	})

	paths := GetTree(memFs)

	println(paths)
}
