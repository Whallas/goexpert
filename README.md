# GoExpert - Desafios e Estudos

Repositório com uma coleção de desafios e projetos práticos desenvolvidos durante o curso **GoExpert**. Cada desafio demonstra diferentes conceitos e padrões em Go, desde arquitetura limpa até otimizações de performance.

## 📚 Projetos

### 1. **Clean Architecture - Sistema de Pedidos**
`Desafios/clean-architecture/`

Implementação do padrão Clean Architecture com múltiplas interfaces de comunicação simultâneas (REST, gRPC e GraphQL) para o mesmo Use Case.

**Stack**: Go, Docker, PostgreSQL, Migration
**Interfaces**: REST (8000), gRPC (50051), GraphQL (8080)

📖 [Veja o README](Desafios/clean-architecture/README.md)

---

### 2. **Client-Server API - Cotação de Câmbio**
`Desafios/client-server-api/`

Aplicação client-server que busca cotação de câmbio de uma API externa e salva o resultado em arquivo.

**Stack**: Go, HTTP API

📖 [Veja o README](Desafios/client-server-api/README.md)

---

### 3. **Clima por CEP com Observabilidade**
`Desafios/clima-por-cep-com-observabilidade/`

Sistema distribuído que busca clima por CEP com tracing distribuído completo usando OpenTelemetry.

**Stack**: Go, Docker, OpenTelemetry, gRPC, Observabilidade
**Serviços**: Service-A (8080), Service-B (8081)

📖 [Veja o README](Desafios/clima-por-cep-com-observabilidade/README.md)

---

### 4. **Multithreading**
`Desafios/multithreading/`

Exemplos práticos de programação concorrente em Go.

**Stack**: Go, Goroutines, Channels

📖 [Veja o README](Desafios/multithreading/README.md)

---

### 5. **Rate Limiter**
`Desafios/rate-limiter/`

Middleware HTTP de rate limiting com suporte a Redis e memória. Controla tráfego por IP ou token com persistência configurável.

**Stack**: Go, Redis, Docker, HTTP Middleware
**Algoritmo**: Fixed-window counter com bloqueio por duração configurável
**Ports**: 8888 (aplicação), 6379 (Redis)

📖 [Veja o README](Desafios/rate-limiter/README.md)

---

### 6. **Sistema de Temperatura por CEP**
`Desafios/sistema-temperatura-por-cep/`

API REST que busca temperatura pelo CEP, integrando-se com API de clima externa (OpenMeteo).

**Stack**: Go, HTTP, Docker, Testes Unitários
**Port**: 8080

📖 [Veja o README](Desafios/sistema-temperatura-por-cep/README.md)

---

### 7. **Stress Test**
`Desafios/stress-test/`

Ferramenta de teste de carga para validar performance e estabilidade de serviços HTTP.

**Stack**: Go, HTTP Client, Load Testing
**Port**: 8080 (alvo de teste)

📖 [Veja o README](Desafios/stress-test/README.md)

---

## 🚀 Como Executar

### Requisitos
- **Go** 1.20+
- **Docker** & **Docker Compose**
- **PostgreSQL** (alguns projetos)
- **Redis** (alguns projetos)

### Executar um projeto específico
```bash
cd Desafios/<nome-do-projeto>
docker compose up --build
```

### Executar testes
```bash
cd Desafios/<nome-do-projeto>
go test ./...
```

---

## 📋 Conceitos Abordados

- ✅ **Clean Architecture** - Separação de responsabilidades, interfaces
- ✅ **Design Patterns** - Strategy, Factory, Middleware
- ✅ **Concorrência** - Goroutines, Channels, Sincronização
- ✅ **APIs** - REST, gRPC, GraphQL
- ✅ **Bancos de Dados** - PostgreSQL, Redis
- ✅ **Cache & Rate Limiting** - Redis, Algoritmos de rate limiting
- ✅ **Observabilidade** - OpenTelemetry, Tracing Distribuído
- ✅ **Testes** - Unit tests, Integração
- ✅ **Docker** - Containerização, Docker Compose

---

## 📝 Estrutura Geral

```
goexpert/
└── Desafios/
    ├── clean-architecture/       # Clean arch + 3 interfaces
    ├── client-server-api/        # Cliente HTTP simples
    ├── clima-por-cep-com-observabilidade/  # Tracing distribuído
    ├── multithreading/           # Goroutines e channels
    ├── rate-limiter/             # Middleware de rate limit
    ├── sistema-temperatura-por-cep/  # API de clima
    └── stress-test/              # Teste de carga
```

---

## 🔧 Troubleshooting

### Porta já em uso
```bash
# Liberar porta (ex: 8080)
lsof -i :8080
kill -9 <PID>
```

### Docker Compose não inicia
```bash
# Limpar containers e volumes
docker compose down -v
docker compose up --build
```

### Erro de conexão ao banco
- Aguarde alguns segundos para o banco subir
- Verifique variáveis de ambiente em `.env` ou `docker-compose.yaml`

---

## 📚 Referências

- [GoExpert Curso](https://goexpert.fullcycle.com.br)
- [Go Documentation](https://golang.org/doc)
- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Última atualização**: Julho de 2026
