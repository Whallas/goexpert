# Rate Limiter

Um **middleware de rate limiter** HTTP escrito em Go. Ele controla o tráfego
por **IP do cliente** ou por **token de acesso**, persistindo os contadores e o
estado de bloqueio no **Redis**. A persistência é construída sobre o **padrão
Strategy**, de modo que o backend pode ser trocado sem alterar a lógica do
limiter.

## Funcionalidades

- Limita requisições **por segundo** por IP ou por token (header `API_KEY`).
- **Precedência Token > IP**: o limite de um token sempre se sobrepõe ao do IP.
- Sobrescrita de limite por token via configuração.
- **Tempo de bloqueio** configurável: ao exceder o limite, o infrator é
  bloqueado por uma janela fixa e toda requisição é rejeitada com `429`.
- Toda a configuração via **variáveis de ambiente** / `.env`.
- Persistência plugável (Redis incluído; implementação em memória disponível).
- Totalmente executável e testável usando **apenas Docker / Docker Compose**.

## Como funciona

```
Requisição HTTP
   │
   ▼
┌─────────────────────────┐   resolve a identidade (token ou IP) + seu limite
│  middleware.RateLimiter │ ──────────────────────────────────────────────┐
└─────────────────────────┘                                                │
   │ delega a decisão (desacoplado do HTTP)                                │
   ▼                                                                       │
┌─────────────────────────┐   contador de janela fixa + flag de bloqueio   │
│      limiter.Limiter    │                                                │
└─────────────────────────┘                                                │
   │ depende apenas da interface Strategy                                  │
   ▼                                                                       ▼
┌─────────────────────────┐   storage.Strategy: Increment / Block / IsBlocked
│  storage.RedisStrategy  │   (pode ser trocado por MemoryStrategy ou outro backend)
└─────────────────────────┘
```

- **Algoritmo**: contador de janela fixa (fixed-window). Cada IP/token tem um
  contador com TTL de 1 segundo. Quando o contador excede o limite, uma *flag
  de bloqueio* com TTL igual ao tempo de bloqueio configurado é definida, e
  todas as requisições seguintes retornam `429` até que ela expire.
- **Desacoplamento**: `middleware` lida apenas com HTTP; `limiter` contém as
  regras de negócio; `storage` cuida da persistência. Cada camada depende da
  camada abaixo através de uma interface.

### Resolução de limite (precedência)

1. Requisição com header `API_KEY` correspondente a uma entrada em
   `TOKEN_LIMITS` → o **limite sobrescrito** desse token.
2. Requisição com qualquer outro header `API_KEY` → seta **`TOKEN_RATE_LIMIT`**
   global.
3. Sem token → seta **`IP_RATE_LIMIT`** para o IP do cliente.

Tokens sempre vencem o IP, exatamente como a especificação exige.

## Configuração

Todas as definições vêm de variáveis de ambiente (veja `.env.example`). Os
valores padrão se aplicam quando uma variável não está definida.

| Variável                 | Padrão           | Descrição                                              |
|--------------------------|------------------|--------------------------------------------------------|
| `SERVER_PORT`            | `8080`           | Porta HTTP em que o servidor escuta.                   |
| `REDIS_ADDR`             | `localhost:6379` | Host:porta do Redis (`redis:6379` dentro do compose).  |
| `REDIS_PASSWORD`         | *(vazio)*        | Senha do Redis, se houver.                             |
| `REDIS_DB`               | `0`              | Número do banco de dados Redis.                        |
| `IP_RATE_LIMIT`          | `10`             | Máximo de requisições/segundo por IP.                  |
| `TOKEN_RATE_LIMIT`       | `100`            | Máximo padrão de requisições/segundo por token.        |
| `TOKEN_LIMITS`           | *(vazio)*        | Sobrescritas por token: `token:limite,token2:limite2`. |
| `BLOCK_DURATION_SECONDS` | `300`            | Por quanto tempo o infrator fica bloqueado após exceder. |

## Executando com Docker Compose

Requer apenas o Docker. A partir deste diretório:

```bash
docker compose up --build
```

Isso sobe:
- `app` — o rate limiter em **http://localhost:8080**
- `redis` — Redis 7 (interno à rede do compose)

Pare com `docker compose down`.

### Experimentando

```bash
# IP limitado a 10 req/s -> a 11ª dentro de um segundo retorna 429
for i in $(seq 1 12); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/; done

# Mesma rajada com um token cujo limite é 100 -> todas passam (Token > IP)
for i in $(seq 1 12); do curl -s -o /dev/null -w "%{http_code}\n" -H "API_KEY: abc123" http://localhost:8080/; done
```

Quando bloqueada, a resposta é HTTP `429` com o corpo exato:

```
you have reached the maximum number of requests or actions allowed within a certain time frame
```

## Executando os testes

Os testes usam a estratégia em memória, então **não precisam de Redis** e rodam
em um container Go descartável:

```bash
docker compose run --rm test
```

Ou diretamente:

```bash
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine go test ./... -v
```

Os testes cobrem: aplicação do limite, bloqueio após exceder, expiração do
bloqueio, o corpo exato do `429`, a **precedência Token > IP**, sobrescritas por
token e o comportamento *fail-closed* em caso de erro no backend.

## Estrutura do projeto

```
.
├── cmd/server/main.go              # ligação: config -> strategy -> limiter -> middleware
├── config/config.go                # configuração via variáveis de ambiente
├── internal/
│   ├── limiter/limiter.go          # lógica de negócio do rate limiter (sem HTTP)
│   ├── middleware/ratelimiter.go   # middleware HTTP (sem lógica de negócio)
│   └── storage/
│       ├── strategy.go             # interface Strategy
│       ├── redis.go                # implementação Redis (padrão)
│       └── memory.go               # implementação em memória
├── Dockerfile
├── docker-compose.yaml
└── .env.example
```
