## Graphql Handlers for popular golang frameworks

This library requires [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)

## Supported Frameworks:

| Framework | Version |
| --------- | ------- |
| Echo      | v4      |
| Gin       | v1      |
| Fiber     | v2      |

## Install:

`go get github.com/dgrng/graphql-handler`

## Example:

See [examples](https://github.com/dgrng/graphql-handler/tree/main/examples)

##### Echo:

```go
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
```
