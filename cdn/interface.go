package cdn

import "net/http"

type IResponse interface {
	GetField() int
	SetField(int)
	GetRawContent() IRawContentFile
	SetRawContent(IRawContentFile)
	GetHttpCode() int
	SetHttpCode(int)
}

type ICDNConfig interface {
	Create(server string)
	GetURLServer() string
}

type IRawContentFile interface {
	FillAttr(data map[string]interface{})
	GetFileName() string
	GetFileSize() int64
	GetClientMimeType() string
}

type ICDN interface {
	AddFile(fullFilePath, fileName string, dadosUsuario map[string]interface{}) (IResponse, error)
	SetCheckVirus(checkVirus bool)
	View(idFile int, w http.ResponseWriter) (bool, error)
	ExistsFile(idFile int) (bool, error)
	SetDownload(download bool)
	GetRawContentFile() IRawContentFile
}
