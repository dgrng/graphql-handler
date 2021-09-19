package main

import (
	"errors"

	graphqlhandler "github.com/dgrng/graphql-handler"
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo/v4"
)

type query struct{}

func (_ query) Hello() string {
	return "Hello world"
}

func main() {
	e := echo.New()

	s := `
	type Query {
			hello: String!
	}
`
	schema := graphql.MustParseSchema(s, &query{}, graphql.UseFieldResolvers())
	onErr := errors.New("Bad Request")
	h := graphqlhandler.New(schema, onErr)
	h.RegisterEcho("/graphql", e, nil)
	e.Start(":8000")
	/*
		visit http://localhost:8000/graphql
	*/
}
