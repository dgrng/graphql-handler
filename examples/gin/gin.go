package main

import (
	"errors"

	graphqlhandler "github.com/dgrng/graphql-handler"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
)

type query struct{}

func (_ query) Hello() string {
	return "Hello world"
}

func main() {
	r := gin.Default()
	s := `
	type Query {
			hello: String!
	}`
	schema := graphql.MustParseSchema(s, &query{}, graphql.UseFieldResolvers())

	onErr := errors.New("Bad Request")
	h := graphqlhandler.New(schema, onErr)
	h.RegisterGin("/graphql", r, nil)
	r.Run(":8000")
	/*
		visit http://localhost:8000/graphql
	*/
}
