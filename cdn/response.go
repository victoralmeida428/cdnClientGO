package cdn

type Response struct {
	fileID         int
	httpCode       int
	rawContentFile IRawContentFile
}

func (r Response) GetRawContent() IRawContentFile {
	return r.rawContentFile
}

func (r *Response) SetRawContent(file IRawContentFile) {
	r.rawContentFile = file
}

func (r Response) GetHttpCode() int {
	return r.httpCode
}

func (r *Response) SetHttpCode(i int) {
	r.httpCode = i
}

func (r Response) GetField() int {
	//TODO implement me
	return r.fileID
}

func (r *Response) SetField(fileID int) {
	//TODO implement me
	r.fileID = fileID
}
