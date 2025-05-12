package main

import (
	"cdn_client/cdn"
	"fmt"
	"net/http"
	"strconv"
)

var Cdn cdn.ICDN

func init() {
	config := cdn.GetInstanceConfig()
	
	config.Create("http://localhost:8080/cdn")
	Cdn = cdn.New()
	Cdn.SetDownload(false)
	
	Cdn.AddFile("./declaracao_de_residencia.pdf", "", map[string]interface{}{"teste": "pdf no GO"})
}

func main() {
	// Configuração do CDN
	
	// Configuração do servidor HTTP
	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		// Obter o ID do arquivo da query string
		idFile := r.URL.Query().Get("id_file")
		if idFile == "" {
			http.Error(w, "ID do arquivo não fornecido", http.StatusBadRequest)
			return
		}
		
		// Converter para inteiro (você pode querer adicionar tratamento de erro aqui)
		
		idInt, err := strconv.Atoi(idFile)
		if err != nil {
			http.Error(w, "ID deve ser um número", http.StatusBadRequest)
			return
		}
		
		// Chamar a função View do CDN
		ok, err := Cdn.View(idInt, w)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if !ok {
			http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
		}
	})
	
	fmt.Println("Servidor iniciado em http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}
