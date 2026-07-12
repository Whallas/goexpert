package graph

import (
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/usecase"
	"github.com/graphql-go/graphql"
)

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"price":       &graphql.Field{Type: graphql.Float},
		"tax":         &graphql.Field{Type: graphql.Float},
		"final_price": &graphql.Field{Type: graphql.Float},
	},
})

func NewSchema(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) (graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type:        graphql.NewList(orderType),
				Description: "List all orders",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					orders, err := listUC.Execute()
					if err != nil {
						return nil, err
					}
					result := make([]map[string]interface{}, 0, len(orders))
					for _, o := range orders {
						result = append(result, map[string]interface{}{
							"id":          o.ID,
							"price":       o.Price,
							"tax":         o.Tax,
							"final_price": o.FinalPrice,
						})
					}
					return result, nil
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createOrder": &graphql.Field{
				Type:        orderType,
				Description: "Create a new order",
				Args: graphql.FieldConfigArgument{
					"price": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
					"tax":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					price, _ := p.Args["price"].(float64)
					tax, _ := p.Args["tax"].(float64)
					output, err := createUC.Execute(usecase.OrderInputDTO{
						Price: price,
						Tax:   tax,
					})
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"id":          output.ID,
						"price":       output.Price,
						"tax":         output.Tax,
						"final_price": output.FinalPrice,
					}, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
