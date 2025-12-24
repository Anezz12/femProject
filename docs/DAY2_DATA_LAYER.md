# Summary Pembelajaran Hari ke-2: Data Layer

**Tanggal:** 24 Desember 2025  
**Materi:** Data Layer - Database Integration, Migrations, dan Repository Pattern  
**Project:** fem_project - Workout API

---

## 📚 Table of Contents

1. [Pengenalan Data Layer](#1-pengenalan-data-layer)
2. [Database Setup dengan Docker](#2-database-setup-dengan-docker)
3. [PostgreSQL Driver](#3-postgresql-driver)
4. [Database Connection](#4-database-connection)
5. [Database Migrations dengan Goose](#5-database-migrations-dengan-goose)
6. [Embed Migrations](#6-embed-migrations)
7. [Database Schema Design](#7-database-schema-design)
8. [Repository Pattern](#8-repository-pattern)
9. [Database Transactions](#9-database-transactions)
10. [Store Implementation](#10-store-implementation)
11. [Integration dengan Application Layer](#11-integration-dengan-application-layer)
12. [Best Practices](#12-best-practices)

---

## 1. Pengenalan Data Layer

### Apa itu Data Layer?

**Data Layer** adalah lapisan dalam aplikasi yang bertanggung jawab untuk:

- 🗄️ Mengelola akses ke database
- 📝 Melakukan operasi CRUD (Create, Read, Update, Delete)
- 🔄 Menangani transaksi database
- 🎯 Abstraksi detail implementasi database dari business logic

### Separation of Concerns

```
┌─────────────────────┐
│   HTTP Handlers     │  ← API Layer (routes, handlers)
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  Business Logic     │  ← Application Layer
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│    Data Layer       │  ← Store/Repository Layer
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│     Database        │  ← PostgreSQL
└─────────────────────┘
```

---

## 2. Database Setup dengan Docker

### Docker Compose Configuration

**File:** `docker-compose.yml`

```yaml
version: "3.8"

services:
  db:
    container_name: "workoutDB"
    image: postgres:12.4-alpine
    volumes:
      - "./database/postgres-data:/var/lib/postgresql/data:rw"
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: "postgres"
      POSTGRES_USER: "postgres"
      POSTGRES_PASSWORD: "postgres"
    restart: unless-stopped

  test_db:
    container_name: "workoutDB_test"
    image: postgres:12.4-alpine
    volumes:
      - "./database/postgres-test-data:/var/lib/postgresql/data:rw"
    ports:
      - "5433:5432"
    environment:
      POSTGRES_DB: "postgres"
      POSTGRES_USER: "postgres"
      POSTGRES_PASSWORD: "postgres"
    restart: unless-stopped
```

### Konsep yang Dipelajari:

#### 1. **Docker Services**

- `db`: Database untuk development
- `test_db`: Database terpisah untuk testing

#### 2. **Container Configuration**

```yaml
container_name: "workoutDB"
image: postgres:12.4-alpine
```

- **container_name**: Nama container untuk mudah diidentifikasi
- **image**: PostgreSQL versi 12.4 dengan Alpine Linux (lightweight)

#### 3. **Volumes**

```yaml
volumes:
  - "./database/postgres-data:/var/lib/postgresql/data:rw"
```

- **Volume mounting**: Persist data di host machine
- **Path mapping**: `host:container:permissions`
- **rw**: Read-Write access
- Data tidak hilang saat container restart

#### 4. **Port Mapping**

```yaml
ports:
  - "5432:5432" # Development DB
  - "5433:5432" # Test DB
```

- Format: `host_port:container_port`
- Development DB: accessible di `localhost:5432`
- Test DB: accessible di `localhost:5433`

#### 5. **Environment Variables**

```yaml
environment:
  POSTGRES_DB: "postgres"
  POSTGRES_USER: "postgres"
  POSTGRES_PASSWORD: "postgres"
```

- Initial database configuration
- Credentials untuk koneksi

#### 6. **Restart Policy**

```yaml
restart: unless-stopped
```

- Container auto-restart jika crash
- Tidak restart jika manually stopped

### Commands:

```bash
# Start database containers
docker-compose up -d

# Stop containers
docker-compose down

# View logs
docker-compose logs db

# Connect to database
docker exec -it workoutDB psql -U postgres
```

---

## 3. PostgreSQL Driver

### PGX Driver

**Package:** `github.com/jackc/pgx/v5/stdlib`

```go
import _ "github.com/jackc/pgx/v5/stdlib"
```

### Konsep yang Dipelajari:

#### 1. **Blank Import (`_`)**

```go
import _ "github.com/jackc/pgx/v5/stdlib"
```

- **Blank identifier**: Import tanpa menggunakan package secara langsung
- **Side effect import**: Driver register dirinya sendiri ke `database/sql`
- Driver otomatis available untuk `sql.Open()`

#### 2. **Kenapa PGX?**

✅ **Kelebihan:**

- Pure Go driver (no C dependencies)
- Excellent performance
- Support PostgreSQL specific features
- Active development dan maintenance
- Type-safe operations

**Alternatif Driver:**

- `lib/pq`: Older, but stable
- `pgx/v5`: Direct driver (without database/sql)

#### 3. **Driver Registration**

```go
sql.Open("pgx", connectionString)
```

- String `"pgx"` identifies driver yang telah di-register
- Connection string mengikuti format driver

---

## 4. Database Connection

### File: `internal/store/database.go`

```go
package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx",
		"host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db: ping %w", err)
	}

	fmt.Println("Database connected...")
	return db, nil
}
```

### Konsep yang Dipelajari:

#### 1. **Connection String**

```go
"host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
```

**Format PostgreSQL DSN (Data Source Name):**

- `host`: Database server address
- `user`: Database username
- `password`: User password
- `dbname`: Database name
- `port`: PostgreSQL port (default: 5432)
- `sslmode`: SSL connection mode
  - `disable`: No SSL (development)
  - `require`: SSL required (production)

#### 2. **sql.Open()**

```go
db, err := sql.Open("pgx", connectionString)
```

- Return `*sql.DB`: Database connection pool
- **Tidak langsung connect** - hanya validate format
- Lazy connection: connect on first use

#### 3. **db.Ping()**

```go
err = db.Ping()
if err != nil {
    return nil, fmt.Errorf("db: ping %w", err)
}
```

- **Verify connection**: Test koneksi ke database
- Memastikan credentials valid
- Database accessible
- Network reachable

#### 4. **Error Wrapping**

```go
return nil, fmt.Errorf("db: open %w", err)
```

- `%w` verb: Wrap error untuk error chain
- Prefix `"db: open"`: Context information
- Allows `errors.Is()` dan `errors.As()` untuk unwrap

#### 5. **Connection Pool**

`*sql.DB` adalah connection pool, bukan single connection:

- Automatically manages multiple connections
- Thread-safe
- Reuses connections
- Handles connection lifecycle

**Pool Configuration (Optional):**

```go
db.SetMaxOpenConns(25)        // Max open connections
db.SetMaxIdleConns(25)        // Max idle connections
db.SetConnMaxLifetime(5 * time.Minute)  // Connection lifetime
```

---

## 5. Database Migrations dengan Goose

### Apa itu Database Migration?

**Migration** adalah version control untuk database schema:

- 📝 Track perubahan schema
- 🔄 Rollback changes jika diperlukan
- 🚀 Deploy schema changes consistently
- 👥 Share schema dengan tim

### Goose Migration Tool

**Package:** `github.com/pressly/goose/v3`

```go
import "github.com/pressly/goose/v3"
```

### Migration Functions

#### 1. **Migrate() Function**

```go
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: set dialect %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: up %w", err)
	}

	fmt.Println("Database migrated...")
	return nil
}
```

**Konsep:**

- `goose.SetDialect("postgres")`: Set database type
- `goose.Up(db, dir)`: Run all pending migrations
- Migrations run in order (00001, 00002, 00003, ...)

#### 2. **MigrateFS() Function**

```go
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
```

**Konsep:**

- `fs.FS`: File system interface (Go 1.16+)
- `goose.SetBaseFS()`: Use embedded file system
- `defer`: Cleanup after function returns
- Clean separation: embedded vs disk migrations

### Migration File Format

**Naming Convention:** `00001_description.sql`

- Sequential number: `00001`, `00002`, `00003`
- Underscore separator
- Descriptive name
- `.sql` extension

**Migration Structure:**

```sql
-- +goose Up
-- SQL statements for upgrading

-- +goose Down
-- SQL statements for downgrading
```

---

## 6. Embed Migrations

### Go Embed Directive

**File:** `migration/fs.go`

```go
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
```

### Konsep yang Dipelajari:

#### 1. **Go Embed Feature (Go 1.16+)**

```go
//go:embed *.sql
var FS embed.FS
```

- **Directive**: `//go:embed` (no space!)
- **Pattern**: `*.sql` embeds all SQL files
- **Type**: `embed.FS` implements `fs.FS` interface
- Files embedded at compile time

#### 2. **Keuntungan Embed:**

✅ Single binary deployment  
✅ No external file dependencies  
✅ Files included in compiled binary  
✅ Simpler deployment process  
✅ No "file not found" errors

#### 3. **Usage dalam Application**

```go
import migrations "github.com/Anezz12/femProject/migration"

// In app initialization
err = store.MigrateFS(pgDB, migrations.FS, ".")
```

- Import package dengan alias `migrations`
- Access embedded FS via `migrations.FS`
- Directory `"."` karena files di root embedded FS

---

## 7. Database Schema Design

### Table: Users

**File:** `migration/00001_user.sql`

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  bio TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
```

### Table: Workouts

**File:** `migration/00002_workouts.sql`

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workouts (
  id BIGSERIAL PRIMARY KEY,
  -- user_id (will be added later with foreign key)
  title VARCHAR(255) NOT NULL,
  description TEXT,
  duration_minutes INTEGER NOT NULL,
  calories_burned INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workouts;
-- +goose StatementEnd
```

### Table: Workout Entries

**File:** `migration/00003_workout_entries.sql`

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
  id BIGSERIAL PRIMARY KEY,
  workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
  exercise_name VARCHAR(255) NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER,
  duration_seconds INTEGER,
  weight DECIMAL(5, 2),
  notes TEXT,
  order_index INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT valid_workout_entry CHECK (
    (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
    (reps IS NULL OR duration_seconds IS NULL)
  )
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_entries;
-- +goose StatementEnd
```

### Konsep yang Dipelajari:

#### 1. **Data Types**

**BIGSERIAL:**

```sql
id BIGSERIAL PRIMARY KEY
```

- Auto-incrementing 64-bit integer
- Range: 1 to 9,223,372,036,854,775,807
- Automatically creates sequence

**VARCHAR:**

```sql
username VARCHAR(50)
email VARCHAR(255)
```

- Variable-length string
- Number = max length
- Storage efficient

**TEXT:**

```sql
bio TEXT
description TEXT
```

- Unlimited length text
- No length limit
- Use untuk content fields

**INTEGER:**

```sql
duration_minutes INTEGER
sets INTEGER
```

- 32-bit integer
- Range: -2,147,483,648 to 2,147,483,647

**DECIMAL:**

```sql
weight DECIMAL(5, 2)
```

- Fixed precision decimal
- `(5, 2)`: 5 total digits, 2 after decimal
- Example: 999.99

**TIMESTAMP WITH TIME ZONE:**

```sql
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
```

- Stores date and time with timezone
- UTC storage recommended
- Automatic timezone conversion

#### 2. **Constraints**

**PRIMARY KEY:**

```sql
id BIGSERIAL PRIMARY KEY
```

- Unique identifier untuk row
- Automatically creates index
- Cannot be NULL

**UNIQUE:**

```sql
username VARCHAR(50) UNIQUE NOT NULL
```

- Ensures uniqueness across table
- Creates unique index
- Can be NULL (unless NOT NULL specified)

**NOT NULL:**

```sql
title VARCHAR(255) NOT NULL
```

- Field must have value
- Cannot insert NULL

**FOREIGN KEY:**

```sql
workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE
```

- Links to another table
- `REFERENCES workouts(id)`: Points to workouts table
- `ON DELETE CASCADE`: Delete entries when parent deleted
- Maintains referential integrity

**CHECK Constraint:**

```sql
CONSTRAINT valid_workout_entry CHECK (
  (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
  (reps IS NULL OR duration_seconds IS NULL)
)
```

- Custom validation logic
- Ensures either `reps` OR `duration_seconds` exists
- Not both, not neither
- Business rule enforcement

#### 3. **DEFAULT Values**

```sql
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
```

- Automatic value if not provided
- `CURRENT_TIMESTAMP`: Server time saat insert

#### 4. **Relationship Design**

```
users (1) ──────── (N) workouts
                       │
                       │
                       └─ (1) ──────── (N) workout_entries
```

- **One-to-Many**: One workout has many entries
- **CASCADE**: Delete entries when workout deleted
- **Referential Integrity**: Enforced by database

---

## 8. Repository Pattern

### Interface Definition

```go
type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
	GetWorkoutOwner(id int64) (int, error)
}
```

### Konsep yang Dipelajari:

#### 1. **Interface Abstraction**

- Define operations tanpa implementation
- Decouples business logic dari database details
- Allows multiple implementations (PostgreSQL, MySQL, MongoDB, etc.)
- Easier testing dengan mock implementations

#### 2. **CRUD Operations**

- **Create**: `CreateWorkout(*Workout) (*Workout, error)`
- **Read**: `GetWorkoutByID(id int64) (*Workout, error)`
- **Update**: `UpdateWorkout(*Workout) error`
- **Delete**: `DeleteWorkout(id int64) error`

#### 3. **Custom Queries**

```go
GetWorkoutOwner(id int64) (int, error)
```

- Business-specific operations
- Beyond basic CRUD
- Domain logic queries

---

## 9. Database Transactions

### Transaction Implementation

```go
func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert workout
	query := `
		INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description,
		workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	// Insert entries
	for _, entry := range workout.Entries {
		query := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps,
				duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets,
			entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes,
			entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}
```

### Konsep yang Dipelajari:

#### 1. **Transaction Basics**

```go
tx, err := pg.db.Begin()
if err != nil {
    return nil, err
}
defer tx.Rollback()
```

**Transaction Properties (ACID):**

- **Atomicity**: All or nothing
- **Consistency**: Valid state transitions
- **Isolation**: Concurrent execution isolation
- **Durability**: Committed changes persist

#### 2. **defer tx.Rollback()**

```go
defer tx.Rollback()
```

- Cleanup function
- Runs when function exits
- Rollback jika error occurs
- Safe to call after commit (no-op)

#### 3. **Why Transactions?**

**Scenario tanpa transaction:**

```
1. Insert workout ✅
2. Insert entry 1 ✅
3. Insert entry 2 ❌ (error)
Result: Orphaned workout! 😱
```

**Scenario dengan transaction:**

```
1. Begin transaction
2. Insert workout
3. Insert entry 1
4. Insert entry 2 (error)
5. Rollback - semua changes dibatalkan ✅
```

#### 4. **Commit Transaction**

```go
err = tx.Commit()
if err != nil {
    return nil, err
}
```

- Finalizes all changes
- Makes changes permanent
- Releases locks
- Only call if all operations success

#### 5. **Parameterized Queries**

```go
query := `
	INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
`
err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description,
	workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
```

**Placeholders:** `$1`, `$2`, `$3`, etc.

- PostgreSQL placeholder syntax
- Prevents SQL injection
- Type-safe query execution

**RETURNING Clause:**

```sql
RETURNING id
```

- PostgreSQL feature
- Returns values after INSERT/UPDATE
- Get generated ID without extra query

#### 6. **QueryRow() and Scan()**

```go
err = tx.QueryRow(query, args...).Scan(&workout.ID)
```

- `QueryRow()`: Execute query expecting single row
- `Scan(&variable)`: Read result into variable
- Automatic type conversion
- Error if no rows or multiple rows

---

## 10. Store Implementation

### Model Structs

```go
type Workout struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}
```

### Repository Implementation

```go
type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}
```

### Konsep yang Dipelajari:

#### 1. **Struct Tags**

```go
ID int `json:"id"`
```

- Metadata attached to struct fields
- `json:"id"`: JSON encoding/decoding name
- Used by `encoding/json` package
- snake_case untuk database, camelCase untuk API

#### 2. **Pointer Fields**

```go
Reps            *int     `json:"reps"`
DurationSeconds *int     `json:"duration_seconds"`
Weight          *float64 `json:"weight"`
```

**Why Pointers?**

- Represent NULL in database
- `nil` = NULL
- Non-nil = has value
- Distinguish between zero value and NULL

**Example:**

```go
var reps *int = nil          // NULL in database
var reps *int = new(int)     // 0 in database
*reps = 10                   // 10 in database
```

#### 3. **Nested Structs**

```go
Entries []WorkoutEntry `json:"entries"`
```

- Array of related entities
- Represents one-to-many relationship
- Populated with JOIN queries or separate queries

#### 4. **Constructor Pattern**

```go
func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}
```

- Dependency injection
- Accepts `*sql.DB` connection
- Returns initialized store
- Convention: `New<Type>` function

---

## 11. Integration dengan Application Layer

### File: `internal/app/app.go`

```go
package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Anezz12/femProject/internal/api"
	"github.com/Anezz12/femProject/internal/store"
	migrations "github.com/Anezz12/femProject/migration"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	// 1. Open database connection
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	// 2. Run migrations
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// 3. Initialize logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// 4. Initialize handlers
	workoutHandler := api.NewWorkoutHandler()

	// 5. Create application
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}

	return app, nil
}
```

### Konsep yang Dipelajari:

#### 1. **Initialization Order**

```
1. Database Connection  ← Must be first
2. Migrations          ← Setup schema
3. Logger              ← For debugging
4. Stores              ← Data access
5. Handlers            ← Business logic
6. Application         ← Wire everything
```

#### 2. **Error Handling vs Panic**

```go
if err != nil {
    return nil, err  // Return error - caller handles
}

if err != nil {
    panic(err)  // Fatal error - cannot continue
}
```

**When to panic:**

- Application cannot function
- Configuration error
- Missing critical dependency

**When to return error:**

- Recoverable errors
- Request-level errors
- Expected failures

#### 3. **Dependency Flow**

```
DB Connection
    ↓
Store (needs DB)
    ↓
Handler (needs Store)
    ↓
Application (wires everything)
```

---

## 12. Best Practices

### 1. **Database Connection**

✅ Use connection pooling (`*sql.DB`)  
✅ Always call `db.Ping()` to verify connection  
✅ Configure pool limits untuk production  
✅ Handle connection errors gracefully  
✅ Close connections pada shutdown

### 2. **Migrations**

✅ Sequential numbering: 00001, 00002, etc.  
✅ Always include Up and Down  
✅ Use transactions dalam migrations  
✅ Test rollback scenarios  
✅ Keep migrations small dan focused  
✅ Embed migrations dalam binary

### 3. **SQL Queries**

✅ Always use parameterized queries (`$1, $2`)  
✅ Never concatenate user input  
✅ Use transactions untuk multiple operations  
✅ Handle NULL values dengan pointers  
✅ Use `RETURNING` untuk get generated IDs

### 4. **Error Handling**

✅ Wrap errors dengan context  
✅ Rollback transactions on error  
✅ Use `defer tx.Rollback()` pattern  
✅ Return meaningful error messages  
✅ Log database errors

### 5. **Repository Pattern**

✅ Define interface untuk abstraction  
✅ One repository per aggregate  
✅ Separate database logic dari business logic  
✅ Use dependency injection  
✅ Keep repositories focused

### 6. **Data Modeling**

✅ Use appropriate data types  
✅ Add constraints untuk data integrity  
✅ Use foreign keys untuk relationships  
✅ Index frequently queried columns  
✅ Normalize untuk reduce redundancy

### 7. **Transaction Management**

✅ Use transactions untuk related operations  
✅ Keep transactions short  
✅ Always rollback on error  
✅ Commit only on success  
✅ Handle transaction deadlocks

---

## 📝 Commands Cheatsheet

### Docker Commands

```bash
# Start database
docker-compose up -d

# Stop database
docker-compose down

# View logs
docker-compose logs -f db

# Connect to database
docker exec -it workoutDB psql -U postgres

# Restart database
docker-compose restart db
```

### PostgreSQL Commands

```sql
-- List databases
\l

-- Connect to database
\c postgres

-- List tables
\dt

-- Describe table
\d users
\d workouts
\d workout_entries

-- Show table data
SELECT * FROM users;
SELECT * FROM workouts;
SELECT * FROM workout_entries;

-- Check migrations
SELECT * FROM goose_db_version;
```

### Go Commands

```bash
# Install dependencies
go get github.com/jackc/pgx/v5/stdlib
go get github.com/pressly/goose/v3

# Run application
go run main.go

# Build application
go build

# Clean modules
go mod tidy
```

---

## 🧪 Testing Database Operations

### Test Connection

```bash
# From terminal
psql -h localhost -U postgres -d postgres -p 5432

# From Go
go run main.go
# Should see: "Database connected..."
# Should see: "Database migrated..."
```

### Verify Schema

```sql
-- Check tables created
\dt

-- Expected output:
-- users
-- workouts
-- workout_entries
-- goose_db_version
```

### Test Insert

```sql
-- Insert test user
INSERT INTO users (username, email, password_hash)
VALUES ('testuser', 'test@example.com', 'hashedpassword');

-- Insert test workout
INSERT INTO workouts (title, description, duration_minutes, calories_burned)
VALUES ('Morning Run', 'Quick 5K run', 30, 250);

-- Verify
SELECT * FROM users;
SELECT * FROM workouts;
```

---

## 🎯 Next Steps

### 1. **Complete CRUD Operations**

- Implement GetWorkoutByID
- Implement UpdateWorkout
- Implement DeleteWorkout
- Implement ListWorkouts

### 2. **Advanced Queries**

- JOIN operations
- Filtering dan sorting
- Pagination
- Search functionality

### 3. **Error Handling**

- Custom error types
- Database error mapping
- Transaction error handling

### 4. **Testing**

- Unit tests untuk store
- Integration tests dengan test database
- Mock database interfaces
- Test fixtures

### 5. **Performance**

- Query optimization
- Index strategy
- Connection pool tuning
- Query caching

### 6. **Security**

- SQL injection prevention
- Prepared statements
- Input validation
- Secure credentials management

---

## 📚 Key Concepts Summary

### Database Connection

```go
*sql.DB               // Connection pool, not single connection
sql.Open()            // Create pool (lazy)
db.Ping()             // Verify connection
```

### Migrations

```go
//go:embed *.sql      // Embed migrations
goose.Up()            // Run migrations
goose.Down()          // Rollback migrations
```

### Transactions

```go
tx.Begin()            // Start transaction
tx.Rollback()         // Cancel changes
tx.Commit()           // Finalize changes
defer tx.Rollback()   // Safety net
```

### SQL Operations

```go
tx.QueryRow()         // Single row query
tx.Query()            // Multiple rows
tx.Exec()             // No result set
Scan()                // Read query results
```

### Repository Pattern

```go
Interface             // Define operations
Implementation        // Concrete database logic
Dependency Injection  // Pass DB to repository
```

---

## ✅ Checklist Pembelajaran

- [x] Setup PostgreSQL dengan Docker
- [x] Install dan configure PGX driver
- [x] Implement database connection
- [x] Setup Goose migrations
- [x] Embed migration files
- [x] Design database schema
- [x] Create migration files
- [x] Implement repository interface
- [x] Create model structs
- [x] Implement database transactions
- [x] Integrate dengan application layer
- [x] Test database connection
- [x] Verify migrations running

---

## 🎓 Key Takeaways

1. **Data Layer Separation**: Isolate database logic dari business logic
2. **Connection Pooling**: `*sql.DB` manages connections efficiently
3. **Migrations**: Version control untuk database schema
4. **Embed Migrations**: Single binary deployment
5. **Transactions**: ACID compliance untuk data integrity
6. **Parameterized Queries**: SQL injection prevention
7. **Repository Pattern**: Clean abstraction layer
8. **Pointers untuk NULL**: Represent database NULL values
9. **Error Wrapping**: Context dalam error messages
10. **Docker Compose**: Consistent development environment

---

**Status:** ✅ Data Layer Complete  
**Ready for:** Implementing full CRUD operations dan business logic

---

_Dibuat pada: 24 Desember 2025_  
_Project: fem_project - Workout API - Day 2_
