package templaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDestinationFilePath(t *testing.T) {
	result := GetDestinationFilePath(GetDestinationFilePathProps{
		TemplatePath:    "/templates/complex",
		DestinationPath: "/output/result",
		File:            "/templates/complex/{{template.name}}/{{template.name}}-a.md",
		Data:            map[string]string{"template.name": "NewProject"}},
	)
	assert.Equal(t, result, "/output/result/NewProject/NewProject-a.md")
}
