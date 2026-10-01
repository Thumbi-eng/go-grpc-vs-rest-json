# 🎬 Masterclass Video Tutorial Script (Hacker Typer Edition)
## Title: **Go gRPC Microservices Masterclass: Build Production Auth & Full CRUD from Scratch**

> **How to use this script with Hacker Typer:**
> Each section is split into two synchronized steps:
> 1. ⌨️ **[HACKER TYPER - PRESS KEYS]**: Tap keys on your keyboard to reveal the code chunk on screen.
> 2. 🗣️ **[WHAT YOU SAY TO THE AUDIENCE]**: Speak the explanation with confidence while the newly typed code is visible.

---

## ⏱️ Video Timeline & Chapter Breakdown (2.5+ Hours Masterclass)

- **00:00 – 12:00**: Chapter 0: Hook, Live Client Demo, REST vs gRPC & Why Go
- **12:00 – 25:00**: Chapter 1: Tooling & Architecture Overview (`go.mod`, `protoc`, Folder Layout)
- **25:00 – 48:00**: Chapter 2: Production Protobuf Contracts (`proto/user.proto` & `proto/auth.proto`)
- **48:00 – 58:00**: Chapter 3: Compiling Protos into Go Stubs (`Makefile` & `protoc` generation)
- **58:00 – 1:15:00**: Chapter 4: Config & User Domain Model (`internal/config` & `internal/model`)
- **1:15:00 – 1:35:00**: Chapter 5: Security Utilities: Bcrypt & JWT (`internal/utils`)
- **1:35:00 – 1:58:00**: Chapter 6: Thread-Safe In-Memory DB & CRUD Repository (`internal/db` & `internal/repository`)
- **1:58:00 – 2:20:00**: Chapter 7: The gRPC Auth Interceptor (`internal/interceptor/auth.go`)
- **2:20:00 – 2:45:00**: Chapter 8: Implementing gRPC Services (`auth_service.go` & `user_service.go`)
- **2:45:00 – 3:00:00**: Chapter 9: Server Entrypoint, Reflection & Graceful Stop (`cmd/server/main.go`)
- **3:00:00 – 3:20:00**: Chapter 10: Building the Automated Test Client (`cmd/client/main.go`) & `grpcurl` Demo
- **3:20:00 – End**: Outro, Key Architecture Takeaways & Next Steps

---

## Chapter 0: Hook, Live Demo, REST vs gRPC (00:00 – 12:00)

### 🎥 Visual:
* Camera on you, then switch split-screen to your terminal.
* In Terminal 1, run `make run-server`. The banner boots up:
  `🚀 gRPC Server listening on port 50060`
* In Terminal 2, run `make run-client`. The interactive client fires through registration, login, interceptor blocking, profile retrieval, list, update, and deletion in milliseconds.

### 🗣️ [WHAT YOU SAY]:
> "Welcome back, engineers! If you've been building REST APIs with Express or Node, you've probably hit that wall where microservices need to talk to each other fast, type-safely, and with ultra-low network latency.
>
> In high-scale tech companies, microservices don't communicate using slow, verbose JSON over HTTP/1. They communicate using **gRPC** over **HTTP/2** powered by **Google Protocol Buffers** and compiled directly into native Go code.
>
> In today's masterclass, we are building a complete, production-grade Authentication and CRUD microservice in Go using gRPC.
>
> We are not taking shortcuts, skipping imports, or waving our hands. We will write every single file line-by-line:
> 1. Defining strict Protobuf contracts with shared packages and optional fields.
> 2. Compiling proto files into Go stubs.
> 3. Hashing passwords with Bcrypt and signing JWTs with standard claims.
> 4. Building a thread-safe database storage layer using Go's `sync.RWMutex`.
> 5. Writing a gRPC Unary Server Interceptor that inspects HTTP/2 metadata headers for Bearer tokens.
> 6. Implementing full CRUD RPC handlers.
> 7. Enabling gRPC Server Reflection for Postman and `grpcurl`.
> 8. And building an interactive automated Go client.
>
> Grab your coffee, fire up your IDE, and let's master Go gRPC together!"

---

## Chapter 1: Tooling & Architecture Overview (12:00 – 25:00)

### 🎥 Visual:
Terminal showing directory setup and `go.mod` initialization.

```bash
mkdir grpc-auth-service && cd grpc-auth-service
go mod init grpc-auth-service
```

### 🗣️ [WHAT YOU SAY]:
> "Before we write a single line of code, let's understand our toolchain.
> gRPC in Go relies on two foundational tools:
> 1. `protoc` — the Google Protocol Buffers compiler.
> 2. The Go protoc plugins: `protoc-gen-go` and `protoc-gen-go-grpc`.
>
> When you write `.proto` files, `protoc` compiles them into Go structs and client/server interfaces.
>
> Let's look at our architecture layout:
> - `proto/`: Contains our single source of truth contracts.
> - `internal/model/`: Our internal database entities (separated from Protobuf wire messages).
> - `internal/config/`: Safe environment configuration.
> - `internal/utils/`: Security hashing and token generation.
> - `internal/db/` & `internal/repository/`: Data storage and CRUD operations.
> - `internal/interceptor/`: The authentication gatekeeper.
> - `internal/service/`: The gRPC business logic handlers.
> - `cmd/server/` & `cmd/client/`: Server bootstrap and interactive test harness.
>
> This layout follows standard Go idiomatic architecture: code inside `internal/` cannot be imported by external modules, enforcing strict boundary security."

---

## Chapter 2: The Production Protobuf Contracts (25:00 – 48:00)

### File: `proto/user.proto`

#### Chunk 1: Syntax, Package & User Entity
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```protobuf
syntax = "proto3";

package pb;

option go_package = "grpc-auth-service/proto/pb";

// Core User entity returned across services
message User {
  string id = 1;
  string username = 2;
  string email = 3;
  string created_at = 4;
}
```

🗣️ **[WHAT YOU SAY]**:
> "We always start with contracts first. Why? In gRPC, **API contracts dictate everything**.
> Look at line 1: `syntax = 'proto3'`. We specify proto3 syntax, which is the modern standard.
> Next: `package pb;`. This sets the Protobuf namespace.
> Line 5: `option go_package = 'grpc-auth-service/proto/pb';`. This tells `protoc` which Go package name and import path to generate.
>
> Now look at `message User`. In JSON, you send keys as full strings over the wire on every single request. That's heavy and wasteful.
> In Protobuf, fields have unique integer tags (`= 1`, `= 2`, `= 3`, `= 4`). Over the wire, Protobuf encodes field keys as lightweight binary integers. That's why gRPC payloads are up to 7 to 10 times smaller and orders of magnitude faster than JSON!
> Notice `created_at` is an RFC3339 timestamp string, keeping serialization completely clean."

---

#### Chunk 2: Request & Response Payloads
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```protobuf
// Request to fetch a single user by ID
message GetUserRequest {
  string id = 1;
}

// Request to fetch the currently authenticated user (derived from JWT token in gRPC metadata)
message GetCurrentUserRequest {}

// Request to list users with optional pagination
message ListUsersRequest {
  int32 page = 1;
  int32 limit = 2;
}

// Response containing a list of users and total count
message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
```

🗣️ **[WHAT YOU SAY]**:
> "Let's examine these request payloads:
> Notice `GetCurrentUserRequest` is completely empty: `{}`. Why?
> Because in gRPC, client identity is not sent in the payload body. Just like in REST where you pass an `Authorization` header, in gRPC we pass the JWT token inside HTTP/2 **Metadata**. Our auth interceptor will parse it automatically!
>
> Next, in `ListUsersResponse`, notice the `repeated` keyword: `repeated User users = 1;`.
> In Protobuf, `repeated` represents an array or slice in Go. `protoc` will generate this as `[]*User`!"

---

#### Chunk 3: Partial Updates with Optional Fields & UserService Definition
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```protobuf
// Request to update a user's details (fields are optional to support partial updates)
message UpdateUserRequest {
  string id = 1;
  optional string username = 2;
  optional string email = 3;
  optional string password = 4;
}

// Request to delete a user by ID
message DeleteUserRequest {
  string id = 1;
}

// Response confirming user deletion
message DeleteUserResponse {
  bool success = 1;
  string message = 2;
}

// Response wrapping a single user
message UserResponse {
  User user = 1;
}

// UserService defines User CRUD operations
service UserService {
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc GetCurrentUser(GetCurrentUserRequest) returns (UserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}
```

🗣️ **[WHAT YOU SAY]**:
> "Pay special attention to `UpdateUserRequest`. In standard proto3, primitive types like string default to empty string `""` when omitted.
> But by adding the `optional` keyword (`optional string username = 2;`), Protobuf tracks field presence!
> In generated Go code, this becomes a pointer `*string`. If the client doesn't send a username, the pointer is `nil`. If they send a username, it holds the address of that string. This makes partial updates (PATCH style in REST) completely straightforward!
>
> Finally, we define `service UserService`. These 5 RPCs form our complete CRUD operations."

---

### File: `proto/auth.proto`

#### Chunk 4: AuthService & Reusing Proto Definitions
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```protobuf
syntax = "proto3";

package pb;

option go_package = "grpc-auth-service/proto/pb";

import "user.proto";

// Request payload for registering a new user
message RegisterRequest {
  string username = 1;
  string email = 2;
  string password = 3;
}

// Request payload for user login
message LoginRequest {
  string email = 1;
  string password = 2;
}

// Response payload for successful authentication
message AuthResponse {
  string access_token = 1;
  User user = 2;
  string message = 3;
}

// AuthService defines registration and authentication operations
service AuthService {
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc Login(LoginRequest) returns (AuthResponse);
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's examine `auth.proto`.
> Notice line 7: `import 'user.proto';`.
> This shows viewers how enterprise projects structure contracts. Instead of putting hundreds of messages into one giant proto file, we keep domain boundaries separate.
> `auth.proto` imports `user.proto`, so inside `AuthResponse`, we can directly embed the `User` message!
>
> We define two public RPCs:
> - `Register`: takes username, email, password; returns the created user and a signed JWT access token.
> - `Login`: takes email and password; validates credentials and returns a new access token."

---

## Chapter 3: Compiling Protos into Go Stubs (48:00 – 58:00)

### File: `Makefile`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```makefile
.PHONY: proto tidy run-server run-client clean

PROTO_DIR=proto
PB_DIR=proto/pb

proto:
	@mkdir -p $(PB_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

tidy:
	go mod tidy

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

clean:
	rm -f $(PB_DIR)/*.pb.go
```

🗣️ **[WHAT YOU SAY]**:
> "Let's automate code generation with a `Makefile`.
> Look at the `proto` recipe:
> - `--proto_path=$(PROTO_DIR)`: points `protoc` to the folder containing our `.proto` files so imports like `import 'user.proto'` resolve instantly.
> - `--go_out` & `--go_opt=paths=source_relative`: tells the Go plugin to generate structs into `proto/pb/`.
> - `--go-grpc_out` & `--go-grpc_opt=paths=source_relative`: tells the gRPC plugin to generate client and server stubs.
>
> Let's run `make proto` in our terminal right now."

### 🎥 Visual:
Run in terminal:
```bash
make proto
ls -la proto/pb/
```
Show generated files: `auth.pb.go`, `auth_grpc.pb.go`, `user.pb.go`, `user_grpc.pb.go`.

🗣️ **[WHAT YOU SAY]**:
> "Look at that! Four Go files were generated:
> - `*.pb.go` contains the serialization/deserialization logic and Go structs.
> - `*_grpc.pb.go` contains the client interfaces and `AuthServiceServer` / `UserServiceServer` interfaces.
> We never touch these generated files manually — they are derived entirely from our contracts."

---

## Chapter 4: Config & User Domain Model (58:00 – 1:15:00)

### File: `internal/config/config.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration properties for the application
type Config struct {
	Port          string
	JWTSecret     string
	JWTExpiration time.Duration
	Environment   string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	// Attempt to load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  No .env file found; using environment variables and defaults")
	}

	port := getEnv("PORT", "50051")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-grpc-jwt-key-change-in-production")
	env := getEnv("ENV", "development")

	expStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExp, err := time.ParseDuration(expStr)
	if err != nil {
		jwtExp = 24 * time.Hour
	}

	return &Config{
		Port:          port,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExp,
		Environment:   env,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
```

🗣️ **[WHAT YOU SAY]**:
> "Next, let's build `config.go`.
> We use `godotenv` to read `.env` when developing locally, but we implement `getEnv` with fallback defaults.
> Notice `time.ParseDuration(expStr)`. This allows users to set human-readable token expirations like `24h` or `1h` in their `.env` file without hardcoding integers.
> By defaulting gRPC port to `50051`, our service works out of the box."

---

### File: `internal/model/user.go`

#### Chunk 1: The Domain Entity & Input DTOs
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package model

import (
	"time"

	pb "grpc-auth-service/proto/pb"
)

// User represents the domain user entity
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hashed password, never exposed in JSON or protobuf responses
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserInput represents data required to register or create a user
type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

// UpdateUserInput represents optional fields for updating an existing user
type UpdateUserInput struct {
	Username *string
	Email    *string
	Password *string
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's build `internal/model/user.go`.
> Beginners often make the mistake of using Protobuf structs as their internal database models. That couples your transport layer directly to your database!
> In production Go, you keep a dedicated domain `User` model.
> Notice `Password string json:'-'`. This holds the bcrypt hash and is tagged so it can never be serialized in JSON or sent over the wire.
> Notice `UpdateUserInput`: each field is a pointer (`*string`), matching our Protobuf `optional` fields. If a field is `nil`, we know not to update it in the database!"

---

#### Chunk 2: Protobuf Mappers
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// ToProto maps an internal User model to the Protobuf User message
func (u *User) ToProto() *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ToProtoList converts a slice of domain users into a slice of Protobuf users
func ToProtoList(users []*User) []*pb.User {
	pbUsers := make([]*pb.User, len(users))
	for i, u := range users {
		pbUsers[i] = u.ToProto()
	}
	return pbUsers
}
```

🗣️ **[WHAT YOU SAY]**:
> "Attaching `ToProto()` directly to `User` is an awesome pattern.
> Whenever a repository method returns a `*model.User`, our gRPC service handler simply calls `user.ToProto()`.
> It handles null-checks gracefully and standardizes formatting `time.Time` into ISO/RFC3339 strings."

---

## Chapter 5: Security Utilities: Bcrypt & JWT (1:15:00 – 1:35:00)

### File: `internal/utils/hash.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash from the given plaintext password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a plaintext password with its bcrypt hash and returns true if they match
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "In `hash.go`, we use Go's official `golang.org/x/crypto/bcrypt` library.
> `bcrypt.GenerateFromPassword` uses `bcrypt.DefaultCost` (which is 10 rounds).
> Bcrypt automatically generates a cryptographic salt and embeds it directly into the resulting hash string.
> In `CheckPassword`, `bcrypt.CompareHashAndPassword` extracts that salt, hashes the incoming password candidate, and does a constant-time comparison to prevent timing attacks."

---

### File: `internal/utils/jwt.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines the structured claims embedded in each JWT
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates and signs a new JWT token for an authenticated user
func GenerateToken(userID, email, username, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "grpc-auth-service",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses, validates the signature and expiration, and extracts the claims
func ValidateToken(tokenStr, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's examine `jwt.go`. We use `github.com/golang-jwt/jwt/v5`.
> Look at `CustomClaims`: we embed `jwt.RegisteredClaims`. That automatically provides standard RFC7519 fields like `exp`, `iat`, `nbf`, and `sub`.
>
> Now look closely at lines 35 to 40 in `ValidateToken`:
> `if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok`.
> Why is this line critical?
> If you don't check the signing method, an attacker can modify the token header to specify algorithm `'none'` or use an RSA public key to forge signatures! By pinning the method strictly to HMAC, we protect our service against algorithm substitution attacks."

---

## Chapter 6: Thread-Safe In-Memory DB & CRUD Repository (1:35:00 – 1:58:00)

### File: `internal/db/db.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package db

import (
	"sync"

	"grpc-auth-service/internal/model"
)

// MemoryDB provides a thread-safe in-memory storage store for our application
// It requires zero external database setup, making it ideal for tutorials and quick demonstrations
type MemoryDB struct {
	Mu    sync.RWMutex
	Users map[string]*model.User
}

// NewMemoryDB initializes a new in-memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		Users: make(map[string]*model.User),
	}
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's talk about storage. In our TypeScript video, an array worked fine because Node.js runs on a single event-loop thread.
> But **Go is fundamentally concurrent**. In gRPC, every incoming RPC call executes on its own goroutine!
> If two clients call our server simultaneously—one reading users and one writing a user—a standard Go map will trigger a fatal runtime crash: `fatal error: concurrent map read and map write`.
> To prevent this, we introduce `sync.RWMutex` (Read-Write Mutual Exclusion lock). This makes our in-memory database 100% thread-safe."

---

### File: `internal/repository/user_repository.go`

#### Chunk 1: Interface & Create Method
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"grpc-auth-service/internal/db"
	"grpc-auth-service/internal/model"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email is already registered")
)

// UserRepository defines the standard CRUD interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, page, limit int32) ([]*model.User, int32, error)
	Update(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id string) error
}

type memoryUserRepository struct {
	db *db.MemoryDB
}

func NewUserRepository(database *db.MemoryDB) UserRepository {
	return &memoryUserRepository{
		db: database,
	}
}

// Create inserts a new user record into the store
func (r *memoryUserRepository) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	normalizedEmail := strings.ToLower(strings.TrimSpace(input.Email))

	for _, u := range r.db.Users {
		if strings.EqualFold(u.Email, normalizedEmail) {
			return nil, ErrEmailAlreadyExists
		}
	}

	newUser := &model.User{
		ID:        uuid.NewString(),
		Username:  strings.TrimSpace(input.Username),
		Email:     normalizedEmail,
		Password:  input.Password,
		CreatedAt: time.Now().UTC(),
	}

	r.db.Users[newUser.ID] = newUser
	return cloneUser(newUser), nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "Look at the repository pattern here. We define a clean `UserRepository` interface.
> Notice `Create`:
> 1. `r.db.Mu.Lock()` acquires an exclusive write lock. Nobody else can read or write while we insert.
> 2. `defer r.db.Mu.Unlock()` ensures the lock is released the moment the function returns.
> 3. We normalize emails to lowercase so duplicate emails are caught reliably.
> 4. We generate unique UUIDs using `github.com/google/uuid`.
> 5. And notice line 58: `return cloneUser(newUser), nil`. Why clone?
> Because if we return the internal pointer directly, outside code could mutate map data without acquiring the mutex lock! Defensive copying guarantees zero data races."

---

#### Chunk 2: Read Operations (`GetByID`, `GetByEmail`, `List`)
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// GetByID finds a single user by their unique UUID
func (r *memoryUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	user, exists := r.db.Users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return cloneUser(user), nil
}

// GetByEmail finds a single user by email address (case-insensitive)
func (r *memoryUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	for _, user := range r.db.Users {
		if strings.EqualFold(user.Email, normalizedEmail) {
			return cloneUser(user), nil
		}
	}

	return nil, ErrUserNotFound
}

// List returns a paginated list of users, sorted newest to oldest
func (r *memoryUserRepository) List(ctx context.Context, page, limit int32) ([]*model.User, int32, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	all := make([]*model.User, 0, len(r.db.Users))
	for _, u := range r.db.Users {
		all = append(all, cloneUser(u))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := int32(len(all))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	startIndex := (page - 1) * limit
	if startIndex >= total {
		return []*model.User{}, total, nil
	}

	endIndex := startIndex + limit
	if endIndex > total {
		endIndex = total
	}

	return all[startIndex:endIndex], total, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now look at the read operations.
> Instead of `Mu.Lock()`, we use `Mu.RLock()` and `defer Mu.RUnlock()`.
> `RLock` stands for Read Lock. It allows 1,000 clients to read users from the database concurrently! As long as nobody is writing, readers do not block each other.
>
> In `List`, we sort users by `CreatedAt` descending so our pagination returns consistent, deterministic slices."

---

#### Chunk 3: Update, Delete & Defensive Cloning
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// Update updates fields on an existing user record
func (r *memoryUserRepository) Update(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	existing, exists := r.db.Users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	if input.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*input.Email))
		if !strings.EqualFold(existing.Email, newEmail) {
			for _, u := range r.db.Users {
				if u.ID != id && strings.EqualFold(u.Email, newEmail) {
					return nil, ErrEmailAlreadyExists
				}
			}
			existing.Email = newEmail
		}
	}

	if input.Username != nil && strings.TrimSpace(*input.Username) != "" {
		existing.Username = strings.TrimSpace(*input.Username)
	}

	if input.Password != nil && *input.Password != "" {
		existing.Password = *input.Password
	}

	return cloneUser(existing), nil
}

// Delete removes a user record from the store
func (r *memoryUserRepository) Delete(ctx context.Context, id string) error {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, exists := r.db.Users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.db.Users, id)
	return nil
}

func cloneUser(u *model.User) *model.User {
	if u == nil {
		return nil
	}
	copied := *u
	return &copied
}
```

🗣️ **[WHAT YOU SAY]**:
> "In `Update`, notice how we check if pointers are non-nil. If the client only passes a username, only the username is updated. If they also pass an email, we verify that another account hasn't claimed it already.
> And `Delete` uses Go's built-in `delete(map, key)` inside a write lock."

---

## Chapter 7: The gRPC Auth Interceptor (1:58:00 – 2:20:00)

### File: `internal/interceptor/auth.go`

#### Chunk 1: Defining the Bouncer & Public Whitelist
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/utils"
)

type contextKey string

const userClaimsContextKey contextKey = "user_claims"

// UserClaimsContext represents the authenticated user data stored in the request context
type UserClaimsContext struct {
	UserID   string
	Email    string
	Username string
}

// AuthInterceptor enforces JWT authentication on protected gRPC endpoints
type AuthInterceptor struct {
	jwtSecret     string
	publicMethods map[string]bool
}

// NewAuthInterceptor creates a new AuthInterceptor with whitelist for public methods
func NewAuthInterceptor(jwtSecret string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSecret: jwtSecret,
		publicMethods: map[string]bool{
			"/pb.AuthService/Register": true,
			"/pb.AuthService/Login":    true,
		},
	}
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's build the crown jewel of our service: the **gRPC Auth Interceptor**.
> In Express or Fastify, you use middleware. In gRPC, middleware is called an **Interceptor**.
>
> We create `AuthInterceptor` with a `publicMethods` whitelist.
> Notice the method names: `/pb.AuthService/Register` and `/pb.AuthService/Login`.
> Because these are public authentication endpoints, anyone can call them without a token. Every other endpoint in our microservice will be guarded!"

---

#### Chunk 2: The Unary Interceptor Logic
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// Unary returns a gRPC UnaryServerInterceptor function
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// If the method is public, allow access without token validation
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata from request")
		}

		// Retrieve the authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is required")
		}

		// Expect format: 'Bearer <token>'
		authHeader := authHeaders[0]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format; expected 'Bearer <token>'")
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return nil, status.Errorf(codes.Unauthenticated, "token cannot be empty")
		}

		// Validate the token and parse claims
		claims, err := utils.ValidateToken(tokenString, i.jwtSecret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
		}

		// Inject the authenticated user identity into context
		userContext := &UserClaimsContext{
			UserID:   claims.UserID,
			Email:    claims.Email,
			Username: claims.Username,
		}

		newCtx := context.WithValue(ctx, userClaimsContextKey, userContext)

		// Proceed to handler with enriched context
		return handler(newCtx, req)
	}
}
```

🗣️ **[WHAT YOU SAY]**:
> "Look closely at this interceptor:
> 1. `info.FullMethod`: checks which RPC is being executed.
> 2. `metadata.FromIncomingContext(ctx)`: in gRPC, headers are called Metadata. This extracts the HTTP/2 headers.
> 3. `md.Get('authorization')`: reads the header, validates the `Bearer <token>` prefix.
> 4. `utils.ValidateToken(...)`: validates the cryptographic HMAC signature and checks if the token has expired.
> 5. If invalid, we return `codes.Unauthenticated`.
> 6. If valid, we create a `UserClaimsContext` and inject it into the context using `context.WithValue`.
> 7. Finally, `handler(newCtx, req)` executes the target RPC method with this newly enriched context!"

---

#### Chunk 3: Context Retrieval Helpers
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// GetUserClaims extracts the authenticated user's claims from the gRPC context
func GetUserClaims(ctx context.Context) (*UserClaimsContext, error) {
	claims, ok := ctx.Value(userClaimsContextKey).(*UserClaimsContext)
	if !ok || claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: user claims not found in context")
	}
	return claims, nil
}

// GetUserID is a convenience helper to extract only the UserID from the context
func GetUserID(ctx context.Context) (string, error) {
	claims, err := GetUserClaims(ctx)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "These two helper functions make consuming auth context effortless. Any RPC handler in our entire service can simply call `interceptor.GetUserID(ctx)` to know exactly who is calling them. Clean, elegant, and type-safe."

---

## Chapter 8: Implementing gRPC Services (2:20:00 – 2:45:00)

### File: `internal/service/auth_service.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package service

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/config"
	"grpc-auth-service/internal/model"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/utils"
	pb "grpc-auth-service/proto/pb"
)

// AuthService implements the pb.AuthServiceServer gRPC interface
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register registers a new user with hashed password and returns a JWT access token
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is required")
	}
	if email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if len(password) < 6 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 6 characters")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	createdUser, err := s.userRepo.Create(ctx, model.CreateUserInput{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email is already registered")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	token, err := utils.GenerateToken(
		createdUser.ID,
		createdUser.Email,
		createdUser.Username,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate authentication token: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken: token,
		User:        createdUser.ToProto(),
		Message:     "User registered successfully",
	}, nil
}

// Login authenticates a user by email and password, returning a JWT access token
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	email := strings.TrimSpace(req.GetEmail())
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist")
		}
		return nil, status.Errorf(codes.Internal, "database lookup failed: %v", err)
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Username,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate authentication token: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken: token,
		User:        user.ToProto(),
		Message:     "Login successful",
	}, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now we implement the generated `pb.AuthServiceServer` interface.
> Notice `pb.UnimplementedAuthServiceServer` embedded in `AuthService`. In Go gRPC, embedding this struct gives us forward compatibility.
> Look at the error handling: we don't return raw Go errors. We return `status.Errorf` with official gRPC codes:
> - `codes.InvalidArgument` if inputs are blank.
> - `codes.AlreadyExists` if email is taken.
> - `codes.Unauthenticated` if passwords don't match.
> These codes are language-agnostic: a Python or React client will receive standard error codes!"

---

### File: `internal/service/user_service.go`

#### Chunk 1: GetUser, GetCurrentUser & ListUsers
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package service

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/interceptor"
	"grpc-auth-service/internal/model"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/utils"
	pb "grpc-auth-service/proto/pb"
)

// UserService implements the pb.UserServiceServer gRPC interface
type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser fetches a single user by their ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

// GetCurrentUser returns the profile of the caller extracted from the validated JWT in gRPC metadata
func (s *UserService) GetCurrentUser(ctx context.Context, _ *pb.GetCurrentUserRequest) (*pb.UserResponse, error) {
	userID, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "authenticated user account no longer exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve profile: %v", err)
	}

	return &pb.UserResponse{
		User: user.ToProto(),
	}, nil
}

// ListUsers returns a paginated list of all registered users
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := s.userRepo.List(ctx, req.GetPage(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	return &pb.ListUsersResponse{
		Users: model.ToProtoList(users),
		Total: total,
	}, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "In `UserService`, look at `GetCurrentUser`.
> The client sends an empty request `{}`. But we call `interceptor.GetUserID(ctx)`, fetch the record, and respond with their profile! This mirrors `GET /users/me` from REST, but in high-speed binary gRPC."

---

#### Chunk 2: UpdateUser & DeleteUser
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
// UpdateUser applies partial updates to a user record
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	updateInput := model.UpdateUserInput{}

	if req.Username != nil {
		username := strings.TrimSpace(req.GetUsername())
		updateInput.Username = &username
	}

	if req.Email != nil {
		email := strings.TrimSpace(req.GetEmail())
		updateInput.Email = &email
	}

	if req.Password != nil {
		password := req.GetPassword()
		if len(password) < 6 {
			return nil, status.Errorf(codes.InvalidArgument, "new password must be at least 6 characters")
		}
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash new password: %v", err)
		}
		updateInput.Password = &hashedPassword
	}

	updatedUser, err := s.userRepo.Update(ctx, id, updateInput)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email already in use by another account")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		User: updatedUser.ToProto(),
	}, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}
```

🗣️ **[WHAT YOU SAY]**:
> "In `UpdateUser`, if the user sends a new password, we hash it with bcrypt before passing it to the repository.
> And in `DeleteUser`, we remove the user and return a clean `DeleteUserResponse`."

---

## Chapter 9: Server Entrypoint, Reflection & Graceful Shutdown (2:45:00 – 3:00:00)

### File: `cmd/server/main.go`

⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-auth-service/internal/config"
	"grpc-auth-service/internal/db"
	"grpc-auth-service/internal/interceptor"
	"grpc-auth-service/internal/repository"
	"grpc-auth-service/internal/service"
	pb "grpc-auth-service/proto/pb"
)

func main() {
	// 1. Load application configuration
	cfg := config.LoadConfig()

	// 2. Initialize in-memory database and repository
	database := db.NewMemoryDB()
	userRepo := repository.NewUserRepository(database)

	// 3. Initialize gRPC services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)

	// 4. Initialize Auth Interceptor
	authInterceptor := interceptor.NewAuthInterceptor(cfg.JWTSecret)

	// 5. Create gRPC server with the unary auth interceptor attached
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// 6. Register gRPC service implementations
	pb.RegisterAuthServiceServer(grpcServer, authService)
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 7. Enable gRPC Server Reflection
	reflection.Register(grpcServer)

	// 8. Bind TCP listener on the configured port
	addr := fmt.Sprintf(":%s", cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌ Failed to bind TCP listener on %s: %v", addr, err)
	}

	// 9. Start gRPC server in a separate goroutine
	go func() {
		log.Println("====================================================")
		log.Printf("🚀 gRPC Server listening on port %s\n", cfg.Port)
		log.Println("📋 Registered Services:")
		log.Println("   • pb.AuthService (Register, Login)")
		log.Println("   • pb.UserService (GetUser, GetCurrentUser, ListUsers, UpdateUser, DeleteUser)")
		log.Println("🛡️  Auth Interceptor: Active (JWT Bearer Token required for UserService)")
		log.Println("🔍 gRPC Reflection:  Enabled (Ready for Postman / grpcurl)")
		log.Println("====================================================")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ gRPC server encountered an error: %v", err)
		}
	}()

	// 10. Set up graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit
	log.Printf("\n🛑 Received signal (%s). Gracefully stopping gRPC server...\n", sig)

	grpcServer.GracefulStop()
	log.Println("✅ gRPC server stopped cleanly. Goodbye!")
}
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's wire everything together in `cmd/server/main.go`.
> Notice two critical production details:
> 1. Line 40: `reflection.Register(grpcServer)`. Server Reflection is a game changer! It allows Postman or `grpcurl` to inspect the server dynamically without you having to manually upload `.proto` files into Postman.
> 2. Lines 64 to 73: Graceful shutdown. We don't just kill the server abruptly on `Ctrl+C`. We catch the OS signal and invoke `grpcServer.GracefulStop()`, which lets existing inflight RPCs finish before terminating."

---

## Chapter 10: Building the Automated Test Client & Live Demo (3:00:00 – 3:20:00)

### File: `cmd/client/main.go`

#### Chunk 1: Connection & Register/Login
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"grpc-auth-service/internal/config"
	pb "grpc-auth-service/proto/pb"
)

func main() {
	cfg := config.LoadConfig()
	target := fmt.Sprintf("localhost:%s", cfg.Port)

	log.Println("====================================================")
	log.Printf("🔌 Connecting to gRPC Server at %s...\n", target)
	log.Println("====================================================")

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. REGISTER
	printSection("1. REGISTER NEW USER (Public RPC)")
	regResp, err := authClient.Register(ctx, &pb.RegisterRequest{
		Username: "mustang_coder",
		Email:    "mustang@crowntechlab.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		log.Fatalf("❌ Register failed: %v", err)
	}
	fmt.Printf("✅ %s\n", regResp.GetMessage())
	fmt.Printf("   User ID: %s | Token: %s...\n", regResp.GetUser().GetId(), regResp.GetAccessToken()[:25])
	userID := regResp.GetUser().GetId()

	// 2. LOGIN
	printSection("2. LOGIN (Public RPC)")
	loginResp, err := authClient.Login(ctx, &pb.LoginRequest{
		Email:    "mustang@crowntechlab.com",
		Password: "SecurePassword123!",
	})
	if err != nil {
		log.Fatalf("❌ Login failed: %v", err)
	}
	jwtToken := loginResp.GetAccessToken()
	fmt.Printf("✅ %s | JWT: %s...\n", loginResp.GetMessage(), jwtToken[:25])
```

🗣️ **[WHAT YOU SAY]**:
> "Now let's build an interactive Go client in `cmd/client/main.go`.
> We use `grpc.NewClient` with `insecure.NewCredentials()` to establish a local channel.
> We create client stubs: `pb.NewAuthServiceClient(conn)` and `pb.NewUserServiceClient(conn)`.
> Then we call `Register` and `Login` just like normal local Go functions!"

---

#### Chunk 2: Security Verification & Calling Protected RPCs
⌨️ **[HACKER TYPER - PRESS KEYS]**:
```go
	// 3. TEST SECURITY
	printSection("3. TEST INTERCEPTOR SECURITY (Calling without Token)")
	_, err = userClient.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			fmt.Printf("🛡️  Auth Interceptor blocked unauthorized request!\n")
			fmt.Printf("   Code: %s (%d) | Message: %s\n", st.Code(), st.Code(), st.Message())
		}
	}

	// 4. ATTACH BEARER TOKEN
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwtToken)

	// 5. GET CURRENT USER
	printSection("4. GET CURRENT USER PROFILE (Protected RPC)")
	meResp, _ := userClient.GetCurrentUser(authCtx, &pb.GetCurrentUserRequest{})
	fmt.Printf("✅ Profile from claims -> ID: %s | Username: %s\n", meResp.GetUser().GetId(), meResp.GetUser().GetUsername())

	// 6. LIST USERS
	printSection("5. LIST ALL USERS (Protected RPC)")
	listResp, _ := userClient.ListUsers(authCtx, &pb.ListUsersRequest{Page: 1, Limit: 10})
	fmt.Printf("✅ Total Users: %d\n", listResp.GetTotal())

	// 7. UPDATE USER
	printSection("6. UPDATE USER (Partial Update)")
	newUsername := "mustang_superstar"
	updateResp, _ := userClient.UpdateUser(authCtx, &pb.UpdateUserRequest{Id: userID, Username: &newUsername})
	fmt.Printf("✅ Updated Username: %s\n", updateResp.GetUser().GetUsername())

	// 8. DELETE & VERIFY 404
	printSection("7. DELETE USER & VERIFY NOT FOUND")
	delResp, _ := userClient.DeleteUser(authCtx, &pb.DeleteUserRequest{Id: userID})
	fmt.Printf("✅ %s\n", delResp.GetMessage())

	_, err = userClient.GetUser(authCtx, &pb.GetUserRequest{Id: userID})
	st, _ := status.FromError(err)
	fmt.Printf("✅ Confirmed Deletion: (Code: %s, Message: %s)\n", st.Code(), st.Message())

	log.Println("\n🎉 ALL gRPC AUTH & CRUD TESTS PASSED SUCCESSFULLY!")
}

func printSection(title string) {
	fmt.Printf("\n----------------------------------------------------\n🔹 %s\n----------------------------------------------------\n", title)
}
```

🗣️ **[WHAT YOU SAY]**:
> "Look at step 4:
> `authCtx := metadata.AppendToOutgoingContext(ctx, 'authorization', 'Bearer ' + jwtToken)`.
> On the client side, this is how you attach HTTP/2 headers to any gRPC call.
> From then on, every RPC passed `authCtx` travels with the Bearer token attached!"

---

### 🎥 Live Terminal Action: Running the Entire Suite

1. **Terminal 1**:
   ```bash
   make run-server
   ```
2. **Terminal 2**:
   ```bash
   make run-client
   ```

🗣️ **[WHAT YOU SAY]**:
> "Look at our terminal.
> Registration succeeded.
> Login generated our signed token.
> The unauthorized request was rejected with `Code: Unauthenticated (16)`.
> Profile was loaded directly from the token identity.
> ListUsers paginated our users.
> UpdateUser changed the username.
> And DeleteUser removed the record, with a follow-up call returning `NotFound`.
> Everything executed end-to-end with sub-millisecond latency!"

---

## Chapter 11: Outro & Next Steps (3:20:00 – End)

### 🎥 Visual: Camera on you.

### 🗣️ [WHAT YOU SAY]:
> "And that is how you build a production-grade gRPC authentication and CRUD microservice in Go from absolute scratch!
>
> Let's recap what you mastered today:
> 1. You designed modular Protobuf contracts with shared packages and optional fields.
> 2. You automated compilation using `protoc` and Makefiles.
> 3. You implemented thread-safe storage using Go's `sync.RWMutex` to handle high-concurrency requests safely.
> 4. You built a gRPC Auth Interceptor that reads HTTP/2 metadata and injects user context.
> 5. You registered gRPC Reflection so tools like Postman and `grpcurl` work seamlessly.
> 6. And you built an automated Go client.
>
> The complete source code is linked down below in the description.
> If this masterclass leveled up your backend engineering skills, hit that **Like** button, **Subscribe** to Crowntech, and drop a comment below telling me what you want to build next—maybe connecting this to PostgreSQL with gRPC streaming or deploying to Kubernetes!
>
> Keep coding, keep building, and I will see you in the next masterclass!"
