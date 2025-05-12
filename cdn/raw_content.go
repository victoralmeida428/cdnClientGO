package cdn

import (
	"encoding/json"
	"errors"
	"github.com/victoralmeida428/cdnClientGO/utils"
	"reflect"
	"strings"
)

type RawContentFile struct {
	CTime                   int
	BaseName                string
	FileName                string
	ClientSize              int64
	ClientMimeType          string
	FullPathSistema         string
	ClientOriginalName      string
	ClientOriginalExtension string
}

func (r *RawContentFile) GetFileName() string {
	//TODO implement me
	return r.FileName
}

func (r *RawContentFile) GetFileSize() int64 {
	//TODO implement me
	return r.ClientSize
}

func (r *RawContentFile) GetClientMimeType() string {
	//TODO implement me
	return r.ClientMimeType
}

func NewRawContentFile(rawContent []byte) (*RawContentFile, error) {
	var rawContentFile RawContentFile
	var rawContentJson map[string]interface{}
	if err := json.Unmarshal(rawContent, &rawContentJson); err != nil {
		return nil, err
	}
	rawContentFile.FillAttr(rawContentJson)
	return &rawContentFile, nil
}

func (r *RawContentFile) FillAttr(data map[string]interface{}) {
	//TODO implement me
	val := reflect.ValueOf(r).Elem()

	for key, value := range data {
		camelKey := utils.ToCamelCase(key)

		field := val.FieldByName(camelKey)
		if !field.IsValid() {
			continue
		}

		if field.CanSet() {
			strValue, ok := value.(string)
			if ok {
				value = strings.TrimSpace(strValue)
			}
			fieldType := field.Type()
			valValue := reflect.ValueOf(value)
			if valValue.Type().ConvertibleTo(fieldType) {
				field.Set(valValue.Convert(fieldType))
			} else {
				panic(errors.New("value type is not convertible to field type"))
			}
		}

	}
}
