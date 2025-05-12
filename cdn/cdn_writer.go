package cdn

import "net/http"

type CDNWriter struct {
	writer     http.ResponseWriter
	isDownload bool
}

func (c *CDNWriter) SetIsDownload(isDownload bool) {
	//TODO implement me
	c.isDownload = isDownload
}

func NewCDNWriter(w http.ResponseWriter) *CDNWriter {
	return &CDNWriter{writer: w, isDownload: false}
}

func (c CDNWriter) Write(header http.Header, body []byte) {
	//TODO implement me
	panic("implement me")
}
