# Clean Architecture - Orders System

Implementação do desafio Clean Architecture com três interfaces de comunicação simultâneas para o mesmo Use Case.

## Como executar

```bash
docker compose up --build
```

O Docker Compose sobe o banco de dados, aplica as migrations automaticamente e inicia a aplicação.

## Portas

| Interface | Porta | Descrição |
|-----------|-------|-----------|
| REST      | 8000  | `GET /order` e `POST /order` |
| gRPC      | 50051 | `OrderService` (com reflection) |
| GraphQL   | 8080  | `/graphql` (com GraphiQL UI) |

## Endpoints

### REST

```
POST http://localhost:8000/order   # criar pedido
GET  http://localhost:8000/order   # listar pedidos
```

### gRPC

Use `grpcurl` ou Evans:

```bash
# Listar pedidos
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders

# Criar pedido
grpcurl -plaintext -d '{"price": 100.50, "tax": 5.25}' localhost:50051 pb.OrderService/CreateOrder
```

### GraphQL

```
http://localhost:8080/graphql
```

Navegue para a URL acima para acessar o GraphiQL (UI interativa).

Queries disponíveis:

```graphql
# Criar pedido
mutation {
  createOrder(price: 100.50, tax: 5.25) {
    id price tax final_price
  }
}

# Listar pedidos
{
  listOrders {
    id price tax final_price
  }
}
```

## Arquivo api.http

O arquivo `api.http` na raiz contém as requisições prontas para REST e GraphQL, compatível com VS Code REST Client e JetBrains HTTP Client.

## Arquitetura

```
internal/
├── domain/entity/      # Entidade Order + interface do repositório
├── usecase/            # CreateOrderUseCase + ListOrdersUseCase
└── infra/
    ├── database/       # MySQL repository
    ├── web/            # REST handlers + webserver
    ├── grpc/           # gRPC service + proto
    └── graph/          # GraphQL schema + resolvers
```

Um único Use Case (`ListOrdersUseCase`) exposto por três interfaces diferentes — REST, gRPC e GraphQL — demonstrando o desacoplamento da Clean Architecture.
