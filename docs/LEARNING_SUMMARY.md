# Summary Pembelajaran Golang - Scaffolding an API Project

**Tanggal:** 23 Desember 2025  
**Project:** fem_project - Workout API  
**Module:** `github.com/Anezz12/femProject`

---

## 📚 Table of Contents

1. [Pengenalan Golang](#1-pengenalan-golang)
2. [Setup Project](#2-setup-project)
3. [Struktur Project](#3-struktur-project)
4. [Dependency Management](#4-dependency-management)
5. [Main Application](#5-main-application)
6. [Application Layer](#6-application-layer)
7. [HTTP Router dengan Chi](#7-http-router-dengan-chi)
8. [Routes Setup](#8-routes-setup)
9. [API Handlers](#9-api-handlers)
10. [Store Layer](#10-store-layer)
11. [Server Configuration](#11-server-configuration)
12. [Best Practices](#12-best-practices)

---

## 1. Pengenalan Golang

### Konsep Dasar

- **Bahasa:** Go (Golang) adalah bahasa pemrograman yang dikembangkan oleh Google
- **Karakteristik:**
  - Statically typed
  - Compiled language
  - Built-in concurrency support
  - Simple syntax
  - Fast compilation dan execution

### Package System

```go
package main  // Package declaration

import (      // Import dependencies
    "fmt"
    "net/http"
)
```

---

## 2. Setup Project

### Inisialisasi Go Module

```bash
go mod init github.com/Anezz12/femProject
```

**Penjelasan:**

- `go.mod` adalah file dependency management di Go
- Module path biasanya menggunakan format repository URL
- Go version: `1.25.5`

### File go.mod

```go
module github.com/Anezz12/femProject

go 1.25.5

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.2.3 // indirect
)
```

---

## 3. Struktur Project

### Arsitektur Directory

```
fem_project/
├── go.mod                          # Dependency management
├── main.go                         # Entry point application
├── docs/                           # Documentation
└── internal/                       # Private application code
    ├── api/                        # HTTP handlers
    │   └── workout_handler.go
    ├── app/                        # Application layer
    │   └── app.go
    ├── routes/                     # Route definitions
    │   └── routes.go
    └── store/                      # Data access layer
        └── workout_store.go
```

### Penjelasan Struktur

#### `/internal`

- Directory khusus Go yang bersifat **private**
- Code di dalam `/internal` tidak bisa di-import oleh package luar
- Best practice untuk application code yang tidak ingin di-expose

#### `/internal/api`

- Berisi HTTP handlers/controllers
- Menangani request dan response
- Business logic interface layer

#### `/internal/app`

- Application initialization dan configuration
- Dependency injection container
- Global application state

#### `/internal/routes`

- Route definitions dan mapping
- HTTP method dan path configuration

#### `/internal/store`

- Data access layer
- Database operations (CRUD)
- Repository pattern implementation

---

## 4. Dependency Management

### Dependencies yang Digunakan

#### 1. Chi Router (`github.com/go-chi/chi`)

```bash
go get github.com/go-chi/chi/v5
```

**Fitur:**

- Lightweight HTTP router
- Compatible dengan `net/http`
- URL parameters support
- Middleware support
- Fast dan minimal overhead

---

## 5. Main Application

### File: `main.go`

```go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Anezz12/femProject/internal/app"
	"github.com/Anezz12/femProject/internal/routes"
)

func main() {
	// 1. Command-line flags
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	// 2. Initialize application
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	// 3. Setup routes
	r := routes.SetupRoutes(app)

	// 4. Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 5. Start server
	app.Logger.Println(fmt.Sprintf("Application started at port %d", port))

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Error starting server: %s", err)
	}
}
```

### Konsep yang Dipelajari:

#### 1. **Command-line Flags**

```go
var port int
flag.IntVar(&port, "port", 8080, "Port to run the server on")
flag.Parse()
```

- Menggunakan package `flag` untuk command-line arguments
- Default port: `8080`
- Cara run: `go run main.go -port=3000`

#### 2. **Error Handling**

```go
app, err := app.NewApplication()
if err != nil {
    panic(err)
}
```

- Go menggunakan multiple return values untuk error handling
- Idiom Go: selalu check error setelah function call
- `panic()` untuk error yang fatal

#### 3. **HTTP Server Configuration**

```go
server := &http.Server{
    Addr:         fmt.Sprintf(":%d", port),
    Handler:      r,
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

**Timeout Configuration:**

- `IdleTimeout`: Waktu maksimal koneksi idle
- `ReadTimeout`: Waktu maksimal membaca request
- `WriteTimeout`: Waktu maksimal menulis response

---

## 6. Application Layer

### File: `internal/app/app.go`

```go
package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Anezz12/femProject/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	// Initialize logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize handlers
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Status is available\n")
}
```

### Konsep yang Dipelajari:

#### 1. **Struct Definition**

```go
type Application struct {
    Logger         *log.Logger
    WorkoutHandler *api.WorkoutHandler
}
```

- Struct adalah composite data type
- Field dengan huruf kapital = exported (public)
- Pointer `*` untuk reference types

#### 2. **Constructor Pattern**

```go
func NewApplication() (*Application, error) {
    // initialization
    return app, nil
}
```

- Convention: function `New<Type>` sebagai constructor
- Return pointer untuk efficiency
- Return error untuk error handling

#### 3. **Logger Setup**

```go
logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
```

- `os.Stdout`: output ke console
- `log.Ldate|log.Ltime`: format timestamp
- Bitwise OR `|` untuk combine flags

#### 4. **Method Definition**

```go
func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
    // method body
}
```

- Method dengan receiver `(app *Application)`
- Pointer receiver untuk modify state atau efficiency
- Signature HTTP handler: `func(w http.ResponseWriter, r *http.Request)`

#### 5. **Dependency Injection**

- Semua dependencies di-inject melalui `Application` struct
- Centralized initialization
- Easier testing dan maintenance

---

## 7. HTTP Router dengan Chi

### Kenapa Chi?

**Kelebihan:**

- ✅ Lightweight dan fast
- ✅ Compatible dengan standard library `net/http`
- ✅ Powerful URL patterns
- ✅ Middleware support
- ✅ Request context handling

### Import Chi

```go
import "github.com/go-chi/chi"
```

---

## 8. Routes Setup

### File: `internal/routes/routes.go`

```go
package routes

import (
	"github.com/Anezz12/femProject/internal/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Health check endpoint
	r.Get("/health", app.HealthCheck)

	// Workout endpoints
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	return r
}
```

### Konsep yang Dipelajari:

#### 1. **Router Initialization**

```go
r := chi.NewRouter()
```

- `chi.NewRouter()` membuat router instance
- Return type: `*chi.Mux` (HTTP request multiplexer)

#### 2. **Route Registration**

```go
r.Get("/path", handlerFunc)
r.Post("/path", handlerFunc)
```

- Method-specific routing
- First parameter: path/pattern
- Second parameter: handler function

#### 3. **URL Parameters**

```go
r.Get("/workouts/{id}", handler)
```

- `{id}` adalah URL parameter
- Dynamic segment dalam URL
- Dapat diakses di handler menggunakan `chi.URLParam(r, "id")`

#### 4. **RESTful API Pattern**

```
GET    /workouts/{id}  → Get workout by ID
POST   /workouts       → Create new workout
```

---

## 9. API Handlers

### File: `internal/api/workout_handler.go`

```go
package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (wh *WorkoutHandler) HandleGetWorkByID(w http.ResponseWriter, r *http.Request) {
	// 1. Extract URL parameter
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	// 2. Convert string to int
	workoutID, err := strconv.Atoi(paramsWorkoutID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 3. Send response
	fmt.Fprintf(w, "Workout ID: %d\n", workoutID)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a Workout\n")
}
```

### Konsep yang Dipelajari:

#### 1. **Handler Struct Pattern**

```go
type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
    return &WorkoutHandler{}
}
```

- Organize related handlers dalam struct
- Constructor untuk initialization
- Allows for dependency injection (future enhancement)

#### 2. **HTTP Handler Signature**

```go
func (wh *WorkoutHandler) HandleGetWorkByID(w http.ResponseWriter, r *http.Request)
```

- `w http.ResponseWriter`: untuk menulis response
- `r *http.Request`: berisi data request
- Method receiver `(wh *WorkoutHandler)` untuk attach ke struct

#### 3. **Extracting URL Parameters**

```go
paramsWorkoutID := chi.URLParam(r, "id")
```

- `chi.URLParam(r, "key")` untuk extract parameter dari URL
- Return type: `string`
- Empty string jika parameter tidak ada

#### 4. **Validation**

```go
if paramsWorkoutID == "" {
    http.NotFound(w, r)
    return
}
```

- Validasi input sebelum processing
- `http.NotFound()` untuk response 404
- Early return pattern untuk error handling

#### 5. **String Conversion**

```go
workoutID, err := strconv.Atoi(paramsWorkoutID)
if err != nil {
    http.NotFound(w, r)
    return
}
```

- `strconv.Atoi()`: convert string to integer
- Error handling untuk invalid conversion

#### 6. **Writing Response**

```go
fmt.Fprintf(w, "Workout ID: %d\n", workoutID)
fmt.Fprintln(w, "Create a Workout\n")
```

- `fmt.Fprintf()`: formatted output
- `fmt.Fprintln()`: print dengan newline
- `w` implements `io.Writer` interface

#### 7. **HTTP Status Codes**

- `http.NotFound(w, r)`: 404 Not Found
- Default: 200 OK (jika tidak set explicitly)

---

## 10. Store Layer

### File: `internal/store/workout_store.go`

```go
package store
```

**Status:** Template/placeholder untuk data access layer

### Konsep Future Implementation:

#### 1. **Repository Pattern**

```go
type WorkoutStore interface {
    GetByID(id int) (*Workout, error)
    Create(workout *Workout) error
    Update(workout *Workout) error
    Delete(id int) error
    List() ([]*Workout, error)
}
```

#### 2. **Database Operations**

- CRUD operations (Create, Read, Update, Delete)
- Database connection management
- Query execution
- Data mapping

#### 3. **Separation of Concerns**

- Handler layer: HTTP concerns
- Store layer: Data concerns
- Clean architecture principles

---

## 11. Server Configuration

### HTTP Server Best Practices

```go
server := &http.Server{
    Addr:         fmt.Sprintf(":%d", port),
    Handler:      r,
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

### Timeout Explanations:

#### 1. **IdleTimeout (1 minute)**

- Waktu maksimal koneksi idle sebelum ditutup
- Mencegah resource exhaustion
- Keep-alive connection management

#### 2. **ReadTimeout (10 seconds)**

- Waktu maksimal membaca entire request (including body)
- Protection against slow-read attacks
- Client harus mengirim request dalam waktu ini

#### 3. **WriteTimeout (30 seconds)**

- Waktu maksimal menulis response
- Includes handler execution time
- Protection against slow-write scenarios

### Why These Timeouts Matter:

- 🛡️ **Security**: Prevent DoS attacks
- 💰 **Resource Management**: Avoid connection leaks
- ⚡ **Performance**: Better resource utilization

---

## 12. Best Practices

### 1. **Project Structure**

✅ Gunakan `/internal` untuk private code  
✅ Separate concerns (handlers, routes, store)  
✅ Clear directory naming

### 2. **Error Handling**

✅ Selalu check error setelah function call  
✅ Return errors, don't ignore them  
✅ Use meaningful error messages

### 3. **HTTP Handlers**

✅ Validate input  
✅ Handle errors gracefully  
✅ Return appropriate status codes  
✅ Use pointer receivers untuk efficiency

### 4. **Code Organization**

✅ One handler per operation  
✅ Group related handlers dalam struct  
✅ Use constructor functions (`New*`)

### 5. **Server Configuration**

✅ Always set timeouts  
✅ Use command-line flags untuk configuration  
✅ Log important events

### 6. **Dependency Management**

✅ Use Go modules  
✅ Explicit version dalam `go.mod`  
✅ Regular `go mod tidy`

---

## 📝 Commands Cheatsheet

```bash
# Initialize module
go mod init github.com/username/project

# Install dependencies
go get github.com/go-chi/chi/v5

# Run application
go run main.go
go run main.go -port=3000

# Build binary
go build

# Clean up dependencies
go mod tidy

# Run tests
go test ./...

# Format code
go fmt ./...
```

---

## 🧪 Testing Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

**Response:**

```
Status is available
```

### Get Workout by ID

```bash
curl http://localhost:8080/workouts/123
```

**Response:**

```
Workout ID: 123
```

### Create Workout

```bash
curl -X POST http://localhost:8080/workouts
```

**Response:**

```
Create a Workout
```

---

## 🎯 Next Steps (Future Learning)

### 1. **Request Body Handling**

- JSON encoding/decoding
- Request validation
- Struct tags

### 2. **Database Integration**

- PostgreSQL/MySQL connection
- SQL queries
- ORM (GORM)
- Migrations

### 3. **Middleware**

- Logging middleware
- Authentication
- CORS
- Rate limiting

### 4. **Error Handling**

- Custom error types
- Error response formatting
- Logging errors

### 5. **Testing**

- Unit tests
- Integration tests
- Table-driven tests
- Mocking

### 6. **Configuration**

- Environment variables
- Configuration files
- Viper library

### 7. **Advanced Features**

- Context usage
- Graceful shutdown
- Metrics dan monitoring
- Docker containerization

---

## 📚 Resources

### Official Documentation

- [Go Documentation](https://go.dev/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [Effective Go](https://go.dev/doc/effective_go)

### Learning Materials

- Go Tour: https://go.dev/tour/
- Go by Example: https://gobyexample.com/
- Awesome Go: https://awesome-go.com/

---

## ✅ Key Takeaways

1. **Go Module System**: Dependency management dengan `go.mod`
2. **Project Structure**: Separation of concerns dengan `/internal`
3. **HTTP Server**: Configuration dengan timeouts
4. **Chi Router**: Lightweight routing library
5. **Handler Pattern**: Struct-based handlers untuk organization
6. **Error Handling**: Go's multiple return values pattern
7. **RESTful API**: Method-based routing (GET, POST)
8. **URL Parameters**: Dynamic routing dengan Chi
9. **Logging**: Centralized logging dengan standard library
10. **Application Layer**: Dependency injection container

---

**Project Status:** ✅ Scaffolding Complete  
**Ready for:** Database Integration & Business Logic Implementation

---

_Dibuat pada: 23 Desember 2025_  
_Project: fem_project - Workout API_
