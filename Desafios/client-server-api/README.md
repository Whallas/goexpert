# Client-Server API - Cotação Dólar

Desafio de implementação de um sistema cliente-servidor para consulta de cotação do dólar (USD/BRL).

## Descrição

Este projeto consiste em dois componentes:

- **Server**: API HTTP que consulta a cotação do dólar em uma API externa e persiste os dados em um banco SQLite
- **Client**: Cliente HTTP que consome a API do servidor e salva o resultado em arquivo de texto

## Estrutura do Projeto

```
client-server-api/
├── server/
│   └── main.go          # Servidor HTTP
├── client/
│   ├── main.go          # Cliente HTTP
│   └── cotacao.txt      # Arquivo gerado com a cotação
└── go.mod
```

## Funcionalidades

### Server
- Expõe endpoint `/cotacao` na porta 8080
- Consulta a API https://economia.awesomeapi.com.br para obter cotação USD/BRL
- Timeout de 200ms para chamada à API externa
- Timeout de 10ms para persistência no banco de dados
- Armazena dados no SQLite usando GORM
- Retorna o valor `bid` da cotação em formato JSON

### Client
- Realiza requisição GET para `http://localhost:8080/cotacao`
- Timeout de 300ms para a requisição
- Salva a resposta no arquivo `cotacao.txt`

## Como Executar

### 1. Iniciar o servidor

```bash
cd server
go run main.go
```

O servidor estará disponível em `http://localhost:8080`

### 2. Executar o cliente

Em outro terminal:

```bash
cd client
go run main.go
```

O arquivo `cotacao.txt` será criado com a cotação obtida.

## Dependências

- Go 1.22.2+
- gorm.io/driver/sqlite
- gorm.io/gorm

## Observações

- O servidor cria automaticamente o banco de dados SQLite (`cotacoes.db`)
- Os timeouts são parte dos requisitos do desafio para gerenciamento de contexto
- O cliente sobrescreve o arquivo `cotacao.txt` a cada execução
