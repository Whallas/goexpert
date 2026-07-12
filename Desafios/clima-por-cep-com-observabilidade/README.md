# Clima por CEP com Observabilidade

Sistema distribuído em Go com dois microsserviços que retornam o clima de uma
cidade a partir de um CEP, instrumentado com **OpenTelemetry** e **Zipkin** para
*distributed tracing*.

## Arquitetura

```
Cliente ──POST──▶ Serviço A ──GET /{cep}──▶ Serviço B ──▶ ViaCEP
                  (input)                   (orquestração)  WeatherAPI
                     │                          │
                     └──── spans (OTLP) ───────┘
                                  │
                          OTEL Collector ──▶ Zipkin
```

- **Serviço A** (`:8080`): recebe o POST, valida o CEP (8 dígitos) e encaminha para o B.
- **Serviço B** (`:8081`): descobre a cidade (ViaCEP), busca a temperatura (WeatherAPI) e converte para C/F/K.
- **OTEL Collector**: recebe spans via OTLP e exporta para o Zipkin.
- **Zipkin** (`:9411`): visualização dos traços.

## Provedores de clima

O Serviço B suporta múltiplos provedores, selecionados pela variável
`WEATHER_PROVIDER`:

| Provedor | Valor | Chave necessária |
|---|---|---|
| Open-Meteo (**padrão**) | `openmeteo` | não |
| WeatherAPI | `weatherapi` | sim (`WEATHER_API_KEY`) |

O Open-Meteo resolve a cidade em coordenadas (API de geocoding) e então busca a
temperatura, sem necessidade de chave de API.

## Pré-requisitos

- Docker + Docker Compose
- Chave da [WeatherAPI](https://www.weatherapi.com/) — **somente** se usar `weatherapi`

## Como executar

Com o provedor padrão (Open-Meteo) basta subir o ecossistema:

```bash
docker compose up --build
```

Para usar a WeatherAPI, crie o `.env` e informe a chave:

```bash
cp .env.example .env
# defina WEATHER_PROVIDER=weatherapi e WEATHER_API_KEY=<sua_chave>
docker compose up --build
```

## Como fazer a requisição (Serviço A)

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```

### Respostas

| Situação | HTTP | Corpo |
|---|---|---|
| Sucesso | `200` | `{"city":"...","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}` |
| CEP inválido (≠ 8 dígitos) | `422` | `invalid zipcode` |
| CEP não encontrado | `404` | `can not find zipcode` |

> Conversões: `F = C × 1.8 + 32` e `K = C + 273` (conforme fórmulas do desafio).

## Como visualizar os traços (Zipkin)

1. Acesse <http://localhost:9411>.
2. Clique em **Run Query** (ou filtre por `serviceName=service-a`).
3. Abra um trace para ver o fluxo `service-a → service-b` com os spans manuais:
   - `encaminha-servico-b` (Serviço A)
   - `busca-cep` (consulta ViaCEP)
   - `busca-temperatura` (consulta WeatherAPI)

## Portas

| Serviço | Porta |
|---|---|
| Serviço A | 8080 |
| Serviço B | 8081 |
| Zipkin | 9411 |
| OTEL Collector (OTLP gRPC / HTTP) | 4317 / 4318 |
