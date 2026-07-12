# Sistema de Temperatura por CEP

Serviço em Go que recebe um CEP, descobre a cidade via [ViaCEP](https://viacep.com.br/)
e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin. O provedor de clima
é configurável: [Open-Meteo](https://open-meteo.com/) (padrão, sem chave) ou
[WeatherAPI](https://www.weatherapi.com/). Containerizado e pronto para deploy no
Google Cloud Run.

## Provedores de clima

| Provedor     | `WEATHER_PROVIDER` | Chave           |
| ------------ | ------------------ | --------------- |
| Open-Meteo   | `openmeteo` (padrão) | não precisa     |
| WeatherAPI   | `weatherapi`       | `WEATHER_API_KEY` |

O Open-Meteo trabalha com coordenadas, então o provedor geocodifica a cidade
(API de geocoding do Open-Meteo) antes de consultar a previsão.

## URL do sistema (Cloud Run)

```
https://temperatura-por-cep-goexpert-1008656153135.southamerica-east1.run.app/weather/{cep}
```

Exemplo:

```bash
curl https://temperatura-por-cep-goexpert-1008656153135.southamerica-east1.run.app/weather/01001000
# {"temp_C":22.5,"temp_F":72.5,"temp_K":295.5}
```

## API

`GET /weather/{cep}`

| Cenário          | Status | Corpo                                            |
| ---------------- | ------ | ------------------------------------------------ |
| Sucesso          | `200`  | `{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }` |
| Formato inválido | `422`  | `{ "message": "invalid zipcode" }`               |
| Não encontrado   | `404`  | `{ "message": "can not find zipcode" }`          |

CEP válido = 8 dígitos numéricos.

### Conversões

- `F = C × 1.8 + 32`
- `K = C + 273`

## Rodando localmente via Docker

1. (Opcional) Para usar a WeatherAPI, obtenha uma chave gratuita em
   <https://www.weatherapi.com/>. Com o padrão Open-Meteo, pule esta etapa.
2. Copie o exemplo de variáveis:

   ```bash
   cp .env.example .env
   # padrão usa Open-Meteo (sem chave).
   # para usar WeatherAPI: WEATHER_PROVIDER=weatherapi e defina WEATHER_API_KEY
   ```

3. Suba o serviço:

   ```bash
   docker compose up --build
   ```

4. Teste:

   ```bash
   curl -i http://localhost:8080/weather/01001000
   ```

## Rodando os testes

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.22-alpine go test ./...
```

## Deploy no Google Cloud Run

```bash
gcloud run deploy temperatura-por-cep-goexpert \
  --source . \
  --region southamerica-east1 \
  --allow-unauthenticated
```

## Estrutura

```
cmd/server          ponto de entrada HTTP
internal/cep        cliente ViaCEP + validação de CEP
internal/weather    provedores de clima (Open-Meteo, WeatherAPI) + factory
internal/temperature  conversões C/F/K
internal/handler    orquestra CEP → clima → resposta
```
