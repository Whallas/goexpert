# Multithreading - Busca de CEP

Desafio de implementação de busca de CEP utilizando duas APIs concorrentes, retornando o resultado da mais rápida.

## Descrição

O programa realiza requisições simultâneas para duas APIs de consulta de CEP e retorna o resultado da API que responder primeiro, descartando a resposta mais lenta.

## APIs Utilizadas

- **BrasilAPI**: `https://brasilapi.com.br/api/cep/v1/{cep}`
- **ViaCEP**: `https://viacep.com.br/ws/{cep}/json/`

## Funcionalidades

- Consulta concorrente em duas APIs de CEP
- Timeout de 1 segundo para as requisições
- Exibe qual API respondeu primeiro
- Retorna os dados do endereço em formato JSON

## Como Executar

```bash
cd multithreading
go run main.go <cep>
```

### Exemplo

```bash
go run main.go 01310100
```

### Saída esperada

```
ViaCep got first:
{"cep":"01310-100","logradouro":"Avenida Paulista","complemento":"...", ...}
```

ou

```
BrasilApi got first:
{"cep":"01310100","state":"SP","city":"São Paulo","neighborhood":"Bela Vista", ...}
```

## Conceitos Aplicados

- **Goroutines**: Execução concorrente das requisições HTTP
- **Channels**: Comunicação entre goroutines
- **Select**: Aguarda a primeira resposta disponível
- **Context com Timeout**: Limita o tempo máximo de espera
- **Race condition controlada**: A primeira API a responder "vence"

## Estrutura

```
multithreading/
├── main.go      # Código principal
└── README.md
```

## Observações

- Se nenhuma API responder em 1 segundo, será exibido "timeout"
- Erros de requisição são exibidos no stderr
- O CEP deve ser passado como argumento na linha de comando
