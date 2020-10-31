package main

import (
	"log"
	"net/http"

	"github.com/felixjung/blog-post/go-interfaces/example-code/postgres"
	"github.com/felixjung/blog-post/go-interfaces/example-code/server"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func main() {
	host := ":8080"
	apiRoot := "/api"
	// FIXME: pass an actual config with valid values here. This currently does
	// not run.
	db, err := postgres.NewDB(&pgx.ConnConfig{})
	if err != nil {
		log.Fatalf("failed to connect to postgres DB: %v", err)
	}

	// Create our server.
	srv := server.NewServer(
		&http.Server{Addr: host},
		// postgres.DB implements user.UserStorage so we can use it here. 🙌
		db,
		mux.NewRouter().PathPrefix(apiRoot).Subrouter(),
	)

	// Register the createUser route handler. It would work the same for other
	// handlers.
	srv.Router.
		Path("/users").
		Methods("POST").
		HandlerFunc(srv.CreateUserHandler())

	// Start the server and listen for connections.
	if err := srv.Listen(); err != nil {
		log.Fatalf("failed to start or run server: %v", err)
	}
}
