package format

import "go-framework/util/xconfig"

type Format struct {
	FileFormat map[string]fileType
}

type fileType interface {
	Load(content []byte, config *map[string]interface{}) error
}

func NewFileFormat() *Format {
	fileTypeMap := make(map[string]fileType)
	fileTypeMap[xconfig.Yaml] = &Yaml{}
	fileTypeMap[xconfig.Json] = &Json{}
	return &Format{
		FileFormat: fileTypeMap,
	}
}
