package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/victoralmeida428/cdnClientGO/cdn"
)

func init() {
	cdnCfg := cdn.GetInstanceConfig()
	// Mantido o localhost para o caso de estar testando diretamente na máquina host
	cdnCfg.Create("http://localhost:8080/cdn")
}

func main() {
	handlePost := func(w http.ResponseWriter, req *http.Request) {
		// 1. Abrir o arquivo README.md local
		file, err := os.Open("README.md")
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao abrir o arquivo README.md: %v", err), http.StatusInternalServerError)
			return
		}
		// Garante que o arquivo será fechado ao final da execução
		defer file.Close()

		// 2. Instanciar o cliente CDN
		clientCdn := cdn.New()

		// 3. Criar os dados de usuário (exigido pela interface do AddFile)
		dadosUsuario := map[string]interface{}{
			"descricao": "Upload de teste do README",
			"origem":    "script local",
		}

		// 4. Fazer o upload chamando AddFile
		response, err := clientCdn.AddFile(file, "README.md", dadosUsuario)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao enviar para o CDN: %v", err), http.StatusInternalServerError)
			return
		}

		// 5. Retornar sucesso mostrando o ID que o CDN gerou
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Upload realizado com sucesso!\nID gerado no CDN: %d\n", response.GetField())
	}


	http.HandleFunc("GET /view/{id}", func(w http.ResponseWriter, req *http.Request) {
		// 1. Extrair e validar o ID da URL
		idParam := req.PathValue("id")
		idFile, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID inválido. Deve ser um número inteiro.", http.StatusBadRequest)
			return
		}

		// 2. Instanciar o cliente CDN
		clientCdn := cdn.New()

		// Opcional: Se quiser forçar o download em vez de exibir no navegador, descomente a linha abaixo
		// clientCdn.SetDownload(true)

		// 3. Fazer o pedido de visualização ao CDN
		ok, err := clientCdn.View(idFile, w)
		if err != nil {
			// Apenas registamos o erro no terminal, sem derrubar o servidor
			log.Printf("Erro ao buscar o arquivo %d: %v", idFile, err)
			http.Error(w, "Erro interno ao comunicar com o CDN", http.StatusInternalServerError)
			return
		}

		// 4. Tratar o caso de o ficheiro não existir (ID não encontrado)
		if !ok {
			http.Error(w, "Ficheiro não encontrado", http.StatusNotFound)
			return
		}

	})


	http.HandleFunc("/", handlePost)


	fmt.Println("Listening on port http://localhost:5000")
	err := http.ListenAndServe("localhost:5000", nil)
	if err != nil {
		panic(err)
	}
}