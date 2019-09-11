package templaid

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
