package templaid

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

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
	memFs := afero.NewMemMapFs()

	writeFs(memFs, map[string]string{
		"/templates/complex/":                                                  "",
		"/templates/complex/{{template.name}}/":                                "",
		"/templates/complex/{{template.name}}/{{template.name}}-a.md.template": "{{template.name}} that should be parsed\n",
		"/templates/complex/{{template.name}}/{{template.name}}-b.md":          "{{template.name}} that shouldnt be parsed\n",
		"/templates/complex/{{template.name}}/file-c.md":                       "file c content\n",
		"/templates/complex/folder-b/":                                         "",
		"/templates/complex/folder-b/.hidden.json":                             "{\n  \"foo\": \"bar\"\n}\n",
		"/templates/complex/folder-b/{{template.name}}-folder-c/":              "",
		"/templates/complex/folder-b/{{template.name}}-folder-d/":              "",
		"/templates/complex/folder-b/{{template.name}}-folder-d/file-e.md":     "{{template.name}} shouldnt be replaced. file-e\n",
		"/templates/complex/folder-b/file-d.md":                                "file d content\n",
		"/templates/complex/folder-e/":                                         "",
		"/templates/complex/folder-e/.gitkeep":                                 "",
	})

	RenderTemplate(RenderTemplateProps{
		Fs:              memFs,
		TemplatePath:    "/templates/complex",
		DestinationPath: "/output",
		IgnorePatterns:  []string{".gitkeep"},
		Data:            map[string]string{"template.name": "NewProject"},
	})

	assertFsEqual(t, memFs, "/output", map[string]string{
		"/output/":            "",
		"/output/NewProject/": "",
		"/output/NewProject/NewProject-a.md.template":    "NewProject that should be parsed\n",
		"/output/NewProject/NewProject-b.md":             "{{template.name}} that shouldnt be parsed\n",
		"/output/NewProject/file-c.md":                   "file c content\n",
		"/output/folder-b/":                              "",
		"/output/folder-b/.hidden.json":                  "{\n  \"foo\": \"bar\"\n}\n",
		"/output/folder-b/NewProject-folder-c/":          "",
		"/output/folder-b/NewProject-folder-d/":          "",
		"/output/folder-b/NewProject-folder-d/file-e.md": "{{template.name}} shouldnt be replaced. file-e\n",
		"/output/folder-b/file-d.md":                     "file d content\n",
		"/output/folder-e/":                              "",
	})

}
