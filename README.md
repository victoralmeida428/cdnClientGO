# 📦 Projeto CDN Client em Go

Este projeto é um cliente CDN simples implementado em Go. Ele inicializa um serviço de CDN local, registra um arquivo e expõe uma rota HTTP (`/view`) para visualizar arquivos via ID.

## 🚀 Requisitos

- Go 1.18 ou superior
- Um serviço CDN disponível localmente em `http://localhost:8080/cdn`
- Arquivo `./declaracao_de_residencia.pdf` presente no diretório raiz

## 📂 Estrutura

```bash
.
├── main.go
├── declaracao_de_residencia.pdf
└── cdn_client/
    └── cdn/
        ├── config.go
        ├── cdn.go
        └── ...
```

> A pasta `cdn_client/cdn` deve conter a interface `ICDN` e suas implementações, incluindo os métodos `New()`, `SetDownload()`, `AddFile()`, e `View()`.

## 🛠️ Como executar

1. **Inicie o servidor CDN** (certifique-se de que o endpoint `http://localhost:8080/cdn` esteja ativo).
2. **Execute o servidor Go:**

```bash
go run main.go
```

3. O servidor HTTP estará disponível em:

```
http://localhost:8082
```

## 🌐 Endpoints

### `GET /view?id_file={id}`

Renderiza o arquivo registrado com o ID fornecido.

#### Parâmetros:
- `id_file` (obrigatório): ID do arquivo registrado no CDN.

#### Exemplo:

```bash
curl http://localhost:8082/view?id_file=1
```

## 📄 Observações

- O ID do arquivo é um número inteiro atribuído internamente na ordem de registro.
- O arquivo `declaracao_de_residencia.pdf` é adicionado automaticamente ao iniciar o servidor.
- A função `SetDownload(false)` garante que o arquivo será exibido inline no navegador, e não baixado.

## 📌 Licença

Este projeto é de uso interno/educacional. Adapte conforme necessário para sua aplicação.
