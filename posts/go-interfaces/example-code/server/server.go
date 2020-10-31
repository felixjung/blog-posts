package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// User is a struct type for a user in the system. The ID is a UUID and
// stored as a string.
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

// UserStorage provides methods for creating, reading, updating, and deleting
// users in a storage dependency like a database or a file system.
type UserStorage interface {
	Create(user User) error
	Read(userID string) (User, error)
	Update(user User) (User, error)
	Delete(userID string) (User, error)
}

type Server struct {
	s      *http.Server
	Router *mux.Router
	us     UserStorage
}

func NewServer(s *http.Server, us UserStorage, r *mux.Router) *Server {
	return &Server{s, r, us}
}

func (s *Server) Listen() error {
	return s.s.ListenAndServe()
}

// createUserHandler returns the http.handlerFunc for the createUser operation.
// By creating a closure, the HandlerFunc has access to the server and thereby
// our UserStorage.
func (s *Server) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Process the request and write to the user. Here we use a dummy.
		routeParams := mux.Vars(r)
		// TODO: verify id route param.
		u := User{
			ID: routeParams["id"],
		}
		if err := s.us.Create(u); err != nil {
			// TODO: handle errors.
		}
		// TODO: write response.
	}
}
