package main

import (
	"errors"

	graphqlhandler "github.com/dgrng/graphql-handler"
	"github.com/gofiber/fiber/v2"
	"github.com/graph-gophers/graphql-go"
)

type query struct{}

func (_ query) Hello() string {
	return "Hello world"
}

func main() {
	app := fiber.New()
	s := `
	type Query {
			hello: String!
	}`
	schema := graphql.MustParseSchema(s, &query{}, graphql.UseFieldResolvers())

	onErr := errors.New("Bad Request")
	h := graphqlhandler.New(schema, onErr)
	h.RegisterFiber("/graphql", app, nil)

	app.Listen(":8000")
	/*
		visit http://localhost:8000/graphql
	*/
}
