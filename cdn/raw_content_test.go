package cdn

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillAttr(t *testing.T) {
	
	t.Run("fill all attrs", func(t *testing.T) {
		data := map[string]interface{}{
			"c_time":                    1,
			"base_name":                 "base name",
			"file_name":                 "file name",
			"client_size":               1,
			"client_mime_type":          "client mime type",
			"full_path_sistema":         "full_path_sistema",
			"client_original_name":      "client_originalName",
			"client_original_extension": "client_original_extension",
		}
		
		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)
		
		rawContent, err := NewRawContentFile(jsonData)
		
		assert.NoError(t, err)
		
		rawContent.FillAttr(data)
		
		assert.Equal(t, "base name", rawContent.BaseName)
		assert.Equal(t, "file name", rawContent.FileName)
		assert.Equal(t, "client mime type", rawContent.ClientMimeType)
		assert.Equal(t, "client mime type", rawContent.ClientMimeType)
		assert.Equal(t, "full_path_sistema", rawContent.FullPathSistema)
		assert.Equal(t, "client_originalName", rawContent.ClientOriginalName)
		assert.Equal(t, "client_original_extension", rawContent.ClientOriginalExtension)
	})
	
	t.Run("panic with wrong type", func(t *testing.T) {
		data := map[string]interface{}{
			"c_time":    1,
			"base_name": 233,
		}
		
		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)
		
		assert.Panics(t, func() {
			_, _ = NewRawContentFile(jsonData)
		})
		
	})
	
}
