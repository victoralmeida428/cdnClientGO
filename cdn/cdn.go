package cdn

import "C"
import (
	"bytes"
	"cdn_client/utils"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Check bool

func (c Check) String() string {
	if c {
		return "1"
	}
	return "0"
}

type CDN struct {
	Server         string
	rawContentFile IRawContentFile
	HttpCode       int
	CheckVirus     Check
	Download       Check
}

func New() *CDN {
	cfg := GetInstanceConfig()
	if cfg == nil {
		panic("cdn.New: instance config is nil")
	}
	
	cdn := CDN{
		Server: cfg.GetURLServer(),
	}
	
	if !utils.ValidateURL(cdn.Server) {
		panic("cdn.New: invalid server url")
	}
	
	return &cdn
	
}

func (c *CDN) addCurl(fullFilePath, fileName string, dadosUsuario map[string]interface{}) (IResponse, error) {
	
	// Abre o arquivo
	file, err := os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	// Prepara o body da requisição multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Adiciona o campo rawContentUser
	userData, err := json.Marshal(dadosUsuario)
	if err != nil {
		return nil, err
	}
	if err = writer.WriteField("rawContentUser", string(userData)); err != nil {
		return nil, err
	}
	
	// Adiciona o campo checkVirus (assumindo que getCheckVirus() retorna bool)
	if err = writer.WriteField("checkVirus", c.CheckVirus.String()); err != nil {
		return nil, err
	}
	
	// Adiciona o arquivo
	if fileName == "" {
		fileName = filepath.Base(fullFilePath)
	}
	part, err := writer.CreateFormFile("file", fileName)
	
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	
	// Fecha o writer para finalizar o body
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	
	// Cria a requisição
	resp, err := c.sendCurl(writer, body, "/add/")
	defer resp.Body.Close()
	c.HttpCode = resp.StatusCode
	
	// Processa a resposta (ajuste conforme sua estrutura de Response)
	var output struct {
		Created        map[string]interface{} `json:"created"`
		RawContentFile map[string]interface{} `json:"raw_content_file"`
		RawContentUser map[string]interface{} `json:"raw_content_user"`
		IdFile         int                    `json:"id_arquivo"`
	}
	// Decodifica o JSON da resposta se necessário
	if err = json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	
	var rawContent RawContentFile
	
	rawContent.FillAttr(output.RawContentFile)
	
	return c.SetResponse(output.IdFile, &rawContent), err
	
}

func (c CDN) sendCurl(writer *multipart.Writer, body *bytes.Buffer, path string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.Server+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	
	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	
	return resp, err
}

func (c *CDN) SetResponse(field int, rawContent IRawContentFile) IResponse {
	var response Response
	response.SetField(field)
	response.SetRawContent(rawContent)
	response.SetHttpCode(c.HttpCode)
	return &response
}

func (c CDN) AddFile(fullFilePath, fileName string, dadosUsuario map[string]interface{}) (IResponse, error) {
	
	if fileName == "" {
		fileName = fullFilePath
	}
	return c.addCurl(fullFilePath, fileName, dadosUsuario)
}

func (c *CDN) SetCheckVirus(checkVirus bool) {
	c.CheckVirus = Check(checkVirus)
}

//func (c *CDN) setFileInfoData(wrapperData http.Header) {
//	dados := make(map[string]interface{})
//	for key, row := range wrapperData {
//		dados[key] = row[0]
//	}
//	rawContent, err := NewRawContentFile(dados)
//	if err != nil {
//		panic(err)
//	}
//	c.rawContentFile = rawContent
//}

func (c *CDN) setHeader(header http.Header) {
	if header == nil {
		panic("cdn.SetHeader: req is nil")
	}
	var disposition string
	fileName := c.rawContentFile.GetFileName()
	if c.Download {
		disposition = "attachment"
	} else {
		disposition = "inline"
	}
	
	header.Set("Content-Description", "File Transfer")
	header.Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, fileName))
	header.Set("Content-Type", c.rawContentFile.GetClientMimeType())
	header.Set("Content-Length", strconv.FormatInt(c.rawContentFile.GetFileSize(), 64))
	header.Set("Expires", "0")
	header.Set("Pragma", "public")
	header.Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	
}

func (c CDN) View(idFile int, w http.ResponseWriter) (bool, error) {
	// Primeiro verificamos a existência do arquivo
	exists, err := c.ExistsFile(idFile)
	if err != nil {
		return false, fmt.Errorf("erro na verificação inicial: %v", err)
	}
	if !exists {
		return false, nil
	}
	
	url := fmt.Sprintf("%s/view/?id_arquivo=%d", c.Server, idFile)
	
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("erro na requisição ao CDN: %v", err)
	}
	defer resp.Body.Close()
	
	// Copia os headers
	utils.CopyHeaders(w.Header(), resp.Header)
	name := resp.Header.Get("File_name")
	
	if c.Download {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	}
	
	// Define o status code
	w.WriteHeader(resp.StatusCode)
	
	// Stream direto do conteúdo
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return false, fmt.Errorf("erro no streaming do conteúdo: %v", err)
	}
	
	return true, nil
}

func (c *CDN) existFileByRefArquivo(refArquivo int) (bool, error) {
	// Prepara o body da requisição
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Adiciona o campo id_arquivo
	if err := writer.WriteField("id_arquivo", strconv.Itoa(refArquivo)); err != nil {
		return false, fmt.Errorf("failed to write id_arquivo field: %v", err)
	}
	
	// Fecha o writer para finalizar o body
	if err := writer.Close(); err != nil {
		return false, fmt.Errorf("failed to close multipart writer: %v", err)
	}
	
	// Envia a requisição usando o método existente
	resp, err := c.sendCurl(writer, body, "/exist/")
	defer resp.Body.Close()
	c.HttpCode = resp.StatusCode
	
	var output string
	
	if err = json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}
	
	return utils.IsTrue(output), nil
}

func (c CDN) ExistsFile(idFile int) (bool, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	if err := writer.WriteField("id_arquivo", strconv.Itoa(idFile)); err != nil {
		return false, fmt.Errorf("failed to write id_arquivo field: %v", err)
	}
	
	// Fecha o writer para finalizar o body
	if err := writer.Close(); err != nil {
		return false, fmt.Errorf("failed to close multipart writer: %v", err)
	}
	
	exists, err := c.existFileByRefArquivo(idFile)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *CDN) SetDownload(download bool) {
	c.Download = Check(download)
}

func (c CDN) GetRawContentFile() IRawContentFile {
	return c.rawContentFile
}
