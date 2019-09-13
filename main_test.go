package templaid

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

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

func TestGetDestinationFilePath(t *testing.T) {
	result := GetDestinationFilePath(
		"/templates/complex",
		"/output/result",
		"/templates/complex/{{template.name}}/{{template.name}}-a.md",
		map[string]string{"template.name": "NewProject"},
	)
	assert.Equal(t, result, "/output/result/NewProject/NewProject-a.md")
}

func TestRenderTemplate(t *testing.T) {
	osFs := afero.NewOsFs()
	memFs := afero.NewMemMapFs()

	copyFs("examples", osFs, memFs)

	RenderTemplate(RenderTemplateProps{
		Fs:              memFs,
		TemplatePath:    "complex",
		DestinationPath: "output",
		IgnorePatterns:  []string{".gitkeep", "complex/folder-c"},
		Data:            map[string]string{"template.name": "NewProject"},
	})

	assert.Equal(t, map[string]string{"output": "",
		"output/NewProject":                             "",
		"output/NewProject/NewProject-a.md.template":    "NewProject that should be parsed\n",
		"output/NewProject/NewProject-b.md":             "{{template.name}} that shouldnt be parsed\n",
		"output/NewProject/file-c.md":                   "file c content\n",
		"output/folder-b":                               "",
		"output/folder-b/.hidden.json":                  "{\n  \"foo\": \"bar\"\n}\n",
		"output/folder-b/NewProject-folder-c":           "",
		"output/folder-b/NewProject-folder-d":           "",
		"output/folder-b/NewProject-folder-d/file-e.md": "{{template.name}} shouldnt be replaced. file-e\n",
		"output/folder-b/file-d.md":                     "file d content\n",
	}, getTreeMap(memFs, "output"))
}
