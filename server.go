package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eztrade/kpi/graph/generated"
	"github.com/eztrade/kpi/graph/logengine"
	postgres "github.com/eztrade/kpi/graph/postgres"
	"github.com/eztrade/kpi/graph/resolvers"
	"github.com/eztrade/login/auth"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {

	err := godotenv.Load(".dev.env")
	if err != nil {
		// log.Println("error loading env", err)
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	postgres.InitDbPool()
	pool := postgres.GetPool()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "*"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		goc := graphql.GetOperationContext(ctx)
		if goc.OperationName != "IntrospectionQuery" {
			if goc.Operation.Operation == "query" {
				logengine.GetTelemetryClient().TrackEvent(string(goc.Operation.Operation) + " " + goc.RawQuery)
			} else {
				logengine.GetTelemetryClient().TrackEvent(goc.RawQuery)
			}
		}
		return next(ctx)
	})

	router.Use(auth.AuthMiddleWare(pool))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	router.Handle("/v2/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}
