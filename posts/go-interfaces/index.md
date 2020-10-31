# JS to Go: Interfaces

> I have been working with Go in a service-oriented architecture over the past one and a half years. It has been a fun ride after years of working with and enjoying JavaScript. This post is part of a series on things I have learned and patterns I have picked up while learning a new language.

In object-oriented programming languages, interfaces are abstract types 
allowing developers to define a set of behaviors that a concrete type, such as a struct or class type, needs to exhibit to be used in a given context.

JavaScript has no concept of interfaces, which is not surprising because JavaScript also does not have a strong static type system. 🙃 If you are familiar with TypeScript you are likely to have come across interfaces, for example to describe the props of a React component. Interfaces in TypeScript are mighty and cover many of the different objects present in JavaScript: object literals, function types, indexable types, and classes.[^1] Typescript interfaces can also be implemented both implicitly and explicitly as shown below.

```ts
interface User {
  email: string;
}

// The object literal foo implicitly implements the User interface.
const foo = {
  email: 'foo@bar.com',
  firstName: 'Foo',
  lastName: 'Bar',
}

// Class SuperUser explicitly implements the User interface by
// declaring so with the 'implements' keyword.
class SuperUser implements User {
  email: string
  permissions: string

  constructor(email: string, permissions: string) {
    this.email = email
    this.permissions = permissions
  }
}
```

The [TypeScript handbook page on interfaces](https://www.typescriptlang.org/docs/handbook/interfaces.html "TypeScript Handbook — Interfaces") has a pretty extensive rundown of how the language supports them. For this post, the above explanations should suffice.

In Go, interfaces are much simpler than in TypeScript, but still incredibly powerful. [Jordan Orelli’s blog post “How to use interfaces in Go”](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go "Jordan Orelli — How to use interfaces in Go") covers a lot of the technical details and intricacies you may run into when working with interfaces in Go. Broadly speaking, interfaces in Go are a named set of methods. Interfaces are implemented implicitly. That is to say, if a given type implements all methods defined by an interface, it implements the interface. Let’s take a look at an example.

## A `UserStorage` Interface Example in Go

Consider the example of a user web service. Such a service could provide a CRUD API for users and depend on some storage mechanism supporting the required Create, Read, Update, and Delete operations. As developers, we do not necessarily care about _how_ these operations are implemented. In fact, we might want to postpone dealing with their concrete implementation until we have a better idea of what the usage patterns will look like.[^2] In such cases — and for other reasons like testing — we introduce an interface, like the `UserStorage` Go interface below.

```go
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
```

`UserStorage` describes all the behaviors we expect a dependency to implement so that we can provide a basic CRUD API in our user service. Let’s create a `Server` struct in a new `server` package with a handler for the create operation of our users API.

```go
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
```

As you can see we are able to write the code for a server without having considered the dependency we want to use to store our user data. Only when we initialize our server in `main()`, we will have to instantiate and pass an actual dependency which implements the `UserStorage` interface.[^3] However, this dependency could write and read users from JSON files in the local file system or from a Postgres database; it does not matter.

If we decide to use a Postgres database as storage for our user data, we implement the `UserStorage` interface with a `DB` struct in a `postgres` package.

```go
package postgres

import (
	"context"
	"fmt"

	"github.com/felixjung/blog-post/go-interfaces/example-code/server"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(config *pgx.ConnConfig) (*DB, error) {
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return &DB{conn}, nil
}

func (db *DB) Create(u server.User) error {
	// TODO: Create the user in postgres
	return nil
}

func (db *DB) Read(userID string) (server.User, error) {
	// TODO: Read the user from postgres
	return server.User{}, nil
}

func (db *DB) Update(u server.User) (server.User, error) {
	// TODO: Update the user in postgres.
	return server.User{}, nil
}

func (db *DB) Delete(userID string) (server.User, error) {
	// TODO: Delete the user in postgres.
	return server.User{}, nil
}
```

The `DB` struct has all the methods defined in the method set of our `UserStorage` interface. It therefore implements `UserStorage` and can be used with our server as shown in the continued example below.

```go
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
```

In the above code we have successfully introduced a `UserStorage` interface to separate the concerns of our users CRUD API and the dependency used to persist and read user data.

## Benefits of using Interfaces in Go

Leveraging interfaces in your Go code has several benefits. These two stand out for me.

1. By using an interface we have a much easier time testing our application business logic, for example in the route handlers, without having to run tests against an actual database. A tool like [GoMock](https://github.com/golang/mock "GoMock mocking framework on GitHub") will generate mocks for the interfaces in our code. In the above example, we would generate a mock for the `UserStorage` interface. The generated mock will be called something like `MockUserStorage`. `MockUserStorage` will implement the `UserStorage` interface so we can use it when instantiating our `server` struct in unit tests. Mocks generated by GoMock expose an `EXPECT()` method through which we can set expectations on the interface’s methods. For example, we can set the expectation that `CreateUser()` will be called by the `createUserHandler` `HandlerFunc`. We get fast unit tests that allow us to develop our handler in a test-driven way. Testing and using mocks will be the topic of another blog post.
2. We can postpone the decision for a specific storage layer until we know exactly what our requirements for that layer are. In the early stages of our project this means we can choose not to implement the storage layer at all and focus on something like the validation logic in our connection handlers. After that we could take things a step further with a prototype and implement the `UserStorage` interface with a simple NoSQL database. This would allow us to postpone decisions regarding the database schema. Once our schema stabilizes and we feel the need to work with a relational database, we may move on to something like our `postgres` package above. Our interface draws a nice architectural boundary that we can take advantage of.

##  My Top 3 Learnings from Using Interfaces in Go

1. **Define interfaces where you use them.** Doing so will allow you to keep them as narrow as you need them to be for the specific use case (i.e., interfaces will only include methods needed in the code that depends on them). It is much easier to reason about a narrow interface than a very wide one, for example when setting expectations against mocks in your tests. Having many narrow interfaces also helps with breaking up monolithic dependencies. Say for example your application relies on a NoSQL database for storage. In addition to the `User` model above, you have an `Order`, a `Product`, and a `Transaction` model. At some point you realize that  you are performing a lot of queries to fetch a list of users together with their orders and the respective product information. You come to the conclusion that a relational database would be better suited to handle these kinds of queries. With your NoSQL dependency, you had a `mongo` package exposing a `MongoDB` struct, which implemented the `UserStorage`, `OrderStorage`, `ProductStorage`, and `TransactionStorage` interfaces used in different parts of your code base. You could inject the same struct in all those places. Because of the implicit implementation of Go interfaces, all you need to do now to switch to a relational database is to create a new `postgres` package with a `PostgresDB` struct implementing the `UserStorage`, `OrderStorage`, and `ProductStorage` interfaces. You will then inject an instances of `PostgresDB` where you previously used an instance of `MongoDB` for these three interfaces. To take advantage of your new relational database, you can now create a new specialized interface to handle your more complex queries and implement that on `PostgresDB`.
2. **Create One Package for Every Dependency Implementation.** Every implementation of an interface should have its own package. By isolating dependencies, especially external dependencies like databases or message brokers, you are also able to test them in isolation. For example, a `postgres` package will contain all the methods for interacting with a Postgres database. You can unit-test those methods very efficiently using a library like [Dockertest](https://github.com/ory/dockertest "GitHub repository of the Dockertest library").
3. **Mock your interfaces with GoMock for great unit test performance.** I have already mentioned GoMock, but I would really like to stress how useful it is when writing unit tests. Use your mocks in tests. Every time you update an interface, regenerate the mocks.

## Conclusion

And that is it. JavaScript does not have interfaces. However, that does not mean you should be afraid of them when picking up Go. With a little bit of practice they will become your ally in building simple and maintainable applications. You can find the [source code for the above Go example on GitHub](https://github.com/felixjung/blog-posts/tree/main/posts/go-interfaces/example-code "Go user service example code on GitHub").

[^1]:	Remember that in JavaScript [almost everything is an object](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Objects/Basics "MDN web docs — JavaScript object basics").

[^2]:	In his classic book, “Clean Architecture”, Robert C. Martin deals with these ideas in Part V on Architecture.

[^3]:	Tip of the hat to Mat Ryer’s for the [pattern of creating route handler closures on the server](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html "Mat Ryer — How I write HTTP services after eight years.") struct.