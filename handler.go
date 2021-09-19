package graphqlhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo/v4"
)

type graphqlHandler struct {
	schema *graphql.Schema
	err    error
}

type params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type RegisterOption struct {
	EnableGraphiql bool
}

var (
	defaultRegisterOption = &RegisterOption{EnableGraphiql: true}
)

// New returns *GraphqlHandler
func New(schema *graphql.Schema, onError error) *graphqlHandler {
	if onError == nil {
		onError = errors.New("BadRequest")
	}
	return &graphqlHandler{
		schema,
		onError,
	}
}

func (h graphqlHandler) exec(ctx context.Context, p params) *graphql.Response {
	return h.schema.Exec(ctx, p.Query, p.OperationName, p.Variables)
}

// FiberHandler : Handler for fiber 2.0
func (h graphqlHandler) FiberHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p params
		if err := json.Unmarshal(c.Body(), &p); err != nil {
			return fiber.ErrBadRequest
		}
		response := h.exec(context.WithValue(c.Context(), "ctx", c), p)
		return c.JSON(response)
	}
}

// RegisterFiber : registers graphql to fiber app
func (h graphqlHandler) RegisterFiber(path string, app *fiber.App, opt *RegisterOption) {
	if opt == nil {
		opt = defaultRegisterOption
	}
	if opt.EnableGraphiql {
		page := graphiql(path)
		app.Get(path, func(c *fiber.Ctx) error {
			c.Context().SetContentType("text/html")
			return c.Send(page)
		})
	}
	app.Post(path, h.FiberHandler())
}

// EchoHandler : Handler for echo v4
func (h graphqlHandler) EchoHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var p params
		if err := json.NewDecoder(c.Request().Body).Decode(&p); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, h.err.Error())
		}
		return c.JSON(http.StatusOK, h.exec(context.WithValue(c.Request().Context(), "ctx", c), p))
	}
}

// RegisterEcho : registers graphql to echo instance
func (h graphqlHandler) RegisterEcho(path string, e *echo.Echo, opt *RegisterOption) {
	if opt == nil {
		opt = defaultRegisterOption
	}
	if opt.EnableGraphiql {
		page := graphiql(path)
		e.GET(path, func(c echo.Context) error {
			return c.HTMLBlob(http.StatusOK, page)
		})
	}
	e.POST(path, h.EchoHandler())
}

// GinHandler : handler for gin
func (h graphqlHandler) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p params
		if err := json.NewDecoder(c.Request.Body).Decode(&p); err != nil {
			c.AbortWithError(http.StatusBadRequest, h.err)
		}
		c.JSON(http.StatusOK, h.exec(context.WithValue(c.Request.Context(), "ctx", c), p))
	}
}

// RegisterGin : register graphql handlers to gin instance
func (h graphqlHandler) RegisterGin(path string, r *gin.Engine, opt *RegisterOption) {
	if opt == nil {
		opt = defaultRegisterOption
	}
	if opt.EnableGraphiql {
		r.GET(path, func(c *gin.Context) {
			page := graphiql(path)
			c.Data(http.StatusOK, "text/html", page)
		})
	}
	r.POST(path, h.GinHandler())
}

func graphiql(path string) []byte {
	t := template.New("graphiql")
	t, err := t.Parse(`
		<!DOCTYPE html>
	<html lang="en">
	  <head>
	    <meta charset="UTF-8" />
	    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
	    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	    <title>Graphiql</title>

	    <link
	      href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css"
	      rel="stylesheet"
	    />
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	  </head>
	  <body>
	    <body style="width: 100%; height: 100%; margin: 0; overflow: hidden">
	      <div id="graphiql" style="height: 100vh">Loading...</div>
	      <script>
	        function graphQLFetcher(graphQLParams) {
	          return fetch("{{ .path }}", {
	            method: "post",
	            body: JSON.stringify(graphQLParams),
	            credentials: "include",
	          })
	            .then(function (response) {
	              return response.text();
	            })
	            .then(function (responseBody) {
	              try {
	                return JSON.parse(responseBody);
	              } catch (error) {
	                return responseBody;
	              }
	            });
	        }
	        ReactDOM.render(
	          React.createElement(GraphiQL, { fetcher: graphQLFetcher }),
	          document.getElementById("graphiql")
	        );
	      </script>
	    </body>
	  </body>
	</html>
		`)
	if err != nil {
		log.Fatalln(err)
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, map[string]string{"path": path}); err != nil {
		log.Fatalln(err)
	}
	return buf.Bytes()
}
