# Jǐn

> 瑾，瑾瑜，美玉也。

Jin is a HTTP web framework written in [Go](https://go.dev/) (Golang) 
with a slim core but limitless extensibility.

## Features

- **Middleware Support**: Easily add global or group-specific middleware.
- **Dependency Injection**: Built-in dependency injection for clean, testable code.
- **Routing**: Flexible routing with parameters and group support.
- **Integrate Non-intrusively**: Can be used as a standard `http.Handler`.

## Getting Started

### Installation

To install Jin, use `go get`:
```bash
go get -u github.com/juanjiTech/jin
```

### Hello, World

Create a `main.go` file with the following code:

```go
package main

import "github.com/juanjiTech/jin"

func main() {
	// Creates a default Jin engine
	r := jin.Default()

	// Define a route for GET requests to the root URL ("/")
	r.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("Hello, World!")
	})

	// Run the server on port 8080
	r.Run(":8080")
}
```

Run the application:
```bash
go run main.go
```
You should now be able to see "Hello, World!" when you navigate to `http://localhost:8080` in your browser.

## Features

### Dependency Injection

Jin has built-in dependency injection, allowing you to write clean and testable handlers.

```go
package main

import (
	"fmt"
	"github.com/juanjiTech/inject/v2"
	"github.com/juanjiTech/jin"
)

// A simple service we want to inject
type GreeterService struct {
	Greeting string
}

func (s *GreeterService) Greet(name string) string {
	return fmt.Sprintf("%s, %s!", s.Greeting, name)
}

func main() {
	r := jin.Default()

	// Map an instance of the GreeterService to the injector
	r.Map(&GreeterService{Greeting: "Hello"})

	// The GreeterService is automatically injected into the handler
	r.GET("/greet/:name", func(c *jin.Context, service *GreeterService) {
		name := c.Param("name")
		message := service.Greet(name)
		c.Writer.WriteString(message)
	})

	r.Run(":8080")
}
```
Visit `http://localhost:8080/greet/Jin` and you will see "Hello, Jin!".

### Routing with Parameters

Jin supports routing with named parameters.

```go
package main

import (
	"github.com/juanjiTech/jin"
)

func main() {
	r := jin.Default()

	// This handler will match /user/john but will not match /user/ or /user
	r.GET("/user/:name", func(c *jin.Context) {
		name := c.Param("name")
		c.Writer.WriteString("Hello, " + name)
	})

	r.Run(":8080")
}
```

### Route Grouping

You can group routes that share a common prefix or middleware.

```go
package main

import (
	"github.com/juanjiTech/jin"
	"log"
)

// A dummy authentication middleware
func AuthMiddleware() jin.HandlerFunc {
	return func(c *jin.Context) {
		log.Println("Authenticating request...")
		// In a real app, you'd check for a token or session
		c.Next()
	}
}

func main() {
	r := jin.Default()

	// Group for API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(AuthMiddleware()) // Apply auth middleware to all v1 routes
	{
		v1.GET("/users", func(c *jin.Context) {
			c.JSON(200, jin.H{"users": []string{"alice", "bob"}})
		})
		v1.GET("/products", func(c *jin.Context) {
			c.JSON(200, jin.H{"products": []string{"laptop", "mouse"}})
		})
	}

	r.Run(":8080")
}
```

### Middleware

You can easily add global middleware to your application.

```go
package main

import (
	"github.com/juanjiTech/jin"
	"log"
	"time"
)

func LoggerMiddleware() jin.HandlerFunc {
	return func(c *jin.Context) {
		start := time.Now()
		c.Next() // Process the request
		end := time.Now()
		log.Printf("[%s] %s %s %d",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			end.Sub(start),
		)
	}
}

func main() {
	r := jin.New()
	r.Use(LoggerMiddleware()) // Use the logger middleware globally
	r.Use(jin.Recovery())     // Use the built-in recovery middleware

	r.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("Hello with Middleware!")
	})

	r.Run(":8080")
}
```

## Performance

Due to the speed of `reflect.Call`, every inject process will take about
200ns, which means if the handler in handler-chain didn't support fast-invoke
will take about 200ns for dependency inject (on mac m2).

## Status

Alpha. Expect API changes and bug fixes.

## License

[MIT](./LICENSE)