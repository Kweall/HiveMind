package main

import (
	"log"
	"net/http"

	"hivemind/graph/generated"
	"hivemind/graph/resolver"
	"hivemind/internal/config"
	"hivemind/internal/db"
	"hivemind/internal/memory"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	cfg := config.Load()
	var res *resolver.Resolver

	switch cfg.StorageType {
	case "postgres":
		dbConn, err := db.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}
		res = resolver.NewResolver(dbConn)
	case "memory":
		res = resolver.NewResolver(memory.NewMemoryStorage())
	default:
		log.Fatalf("unknown storage type: %s", cfg.StorageType)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
