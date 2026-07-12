# Stress Test CLI

Ferramenta de teste de carga via linha de comando escrita em Go. Dispara um número
configurável de requisições HTTP contra uma URL alvo usando um pool de workers e,
ao final, imprime um relatório com o tempo total, a quantidade de requisições e a
distribuição dos códigos de status HTTP.

## Parâmetros

| Flag            | Obrigatório | Padrão | Descrição                                    |
|-----------------|-------------|--------|----------------------------------------------|
| `--url`         | sim         | —      | URL do serviço a ser testado                 |
| `--requests`    | sim         | —      | Número total de requisições a realizar       |
| `--concurrency` | não         | `1`    | Número de chamadas simultâneas               |
| `--timeout`     | não         | `10s`  | Timeout por requisição (ex: `5s`, `500ms`)   |

Restrições: `--requests > 0`, `--concurrency > 0` e `--concurrency <= --requests`.
O número total de requisições é sempre cumprido de forma exata, independentemente
do nível de concorrência.

## Buildar a imagem Docker

```bash
docker build -t stress-test .
```

## Executar o teste

```bash
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

Exemplo com timeout customizado:

```bash
docker run stress-test --url=https://google.com --requests=500 --concurrency=20 --timeout=3s
```

## Exemplo de saída

```
==================== Load Test Report ====================
Total time taken:          4.231s
Total requests performed:  1000
HTTP 200 responses:        985
Status code distribution:
  HTTP 200: 985
  HTTP 429: 12
  HTTP 500: 3
=========================================================
```

Se uma requisição nunca receber resposta (erro de transporte, falha de DNS,
timeout), ela é contabilizada em `Failed (no response)` em vez de um código de
status HTTP.

## Executar localmente (sem Docker)

```bash
go run . --url=http://google.com --requests=1000 --concurrency=10
```

## Estrutura do projeto

```
.
├── main.go                     # parsing de flags e entrypoint
├── internal/loadtest/
│   ├── runner.go               # config, validação, pool de workers
│   └── report.go               # agregação e impressão do relatório
├── Dockerfile                  # build multi-stage → imagem scratch
└── README.md
```
