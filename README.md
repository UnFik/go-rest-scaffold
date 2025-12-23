# Golang Clean Architecture Template

RESTful API built with clean architecture principles for managing users, contacts, and addresses.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Makefile Commands](#makefile-commands)

## ✨ Features

- 🏗️ **Clean Architecture** - Separation of concerns with clear boundaries
- 🔐 **Authentication** - Token-based authentication system
- 📝 **CRUD Operations** - Complete CRUD for Users, Contacts, and Addresses
- 📚 **API Documentation** - Auto-generated Swagger/OpenAPI documentation
- ✅ **Input Validation** - Request validation using go-playground/validator
- 🔍 **Logging** - Structured logging with Logrus
- 🗄️ **Database Migration** - Version-controlled database schema
- 🧪 **Unit Testing** - Comprehensive test coverage
- 🚀 **Hot Reload** - Development mode with Air

## 🏛️ Architecture

![Clean Architecture](architecture.png)

### Architecture Flow

1. External system perform request (HTTP, gRPC, Messaging, etc)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system 
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

### Layer Responsibilities

- **Delivery Layer** (`internal/delivery/`) - HTTP handlers, middleware, routes
- **Use Case Layer** (`internal/usecase/`) - Business logic and orchestration
- **Repository Layer** (`internal/repository/`) - Database operations
- **Entity Layer** (`internal/entity/`) - Domain models
- **Model Layer** (`internal/model/`) - Request/Response DTOs

## 🛠️ Tech Stack

### Core

- **Go 1.25** - Programming language
- **PostgreSQL** - Primary database

### Frameworks & Libraries

| Library | Purpose | Link |
|---------|---------|------|
| GoFiber v2 | HTTP Framework | [github.com/gofiber/fiber](https://github.com/gofiber/fiber) |
| GORM | ORM | [github.com/go-gorm/gorm](https://github.com/go-gorm/gorm) |
| Viper | Configuration | [github.com/spf13/viper](https://github.com/spf13/viper) |
| Migrate | Database Migration | [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| Validator | Input Validation | [github.com/go-playground/validator](https://github.com/go-playground/validator) |
| Logrus | Logging | [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) |
| Swaggo | API Documentation | [github.com/swaggo/swag](https://github.com/swaggo/swag) |
| UUID | Unique Identifiers | [github.com/google/uuid](https://github.com/google/uuid) |

## 📦 Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile)

## 🚀 Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd go-rest-scaffold
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Install development tools

```bash
make install-tools
```

Or manually:

```bash
# Swagger documentation generator
go install github.com/swaggo/swag/cmd/swag@latest

# Database migration tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Hot reload (optional)
go install github.com/air-verse/air@latest
```

## ⚙️ Configuration

All configuration is in `config.json` file:

```json
{
  "app": {
    "name": "go-rest-scaffold"
  },
  "web": {
    "prefork": false,
    "port": 3000
  },
  "log": {
    "level": 6
  },
  "database": {
    "username": "postgres",
    "password": "postgres",
    "host": "localhost",
    "port": 5432,
    "name": "golang_clean_architecture",
    "pool": {
      "idle": 10,
      "max": 100,
      "lifetime": 300
    }
  }
}
```

## 🗄️ Database Setup

### Create Database

```bash
createdb -U postgres golang_clean_architecture
```

Or using Makefile:

```bash
make db-create
```

### Run Migrations

```bash
migrate -database "postgres://postgres:postgres@localhost:5432/golang_clean_architecture?sslmode=disable" -path db/migrations up
```

Or using Makefile:

```bash
make migrate-up
```

### Create New Migration

```bash
migrate create -ext sql -dir db/migrations create_table_xxx
```

Or using Makefile:

```bash
make migrate-create name=create_table_xxx
```

## 🏃 Running the Application

### Development Mode (with hot reload)

```bash
make dev
```

### Production Mode

```bash
make run
```

Or:

```bash
go run cmd/web/main.go
```

### Build Binary

```bash
make build
```

The application will start on `http://localhost:3000`

## 📚 API Documentation

### Swagger UI

Access the interactive API documentation at:

```
http://localhost:3000/swagger/index.html
```

### Generate/Update Documentation

After adding or modifying API endpoints:

```bash
make swagger
```

Or:

```bash
swag init -g cmd/web/main.go
```

### Authentication

1. Register a new user via `POST /api/users`
2. Login via `POST /api/users/_login` to get token
3. Click **"Authorize"** button in Swagger UI
4. Enter your token (without "Bearer" prefix)
5. All authenticated endpoints will now include the token

## 🧪 Testing

### Run All Tests

```bash
make test
```

Or:

```bash
go test -v ./test/
```

### Run Tests with Coverage

```bash
make test-coverage
```

This will generate a coverage report in `coverage.html`

## 📁 Project Structure

```
go-rest-scaffold/
│
├── cmd/                                    # Application entry points
│   └── web/
│       └── main.go                        # Main application (HTTP server)
│
├── internal/                              # Private application code
│   ├── config/                           # Configuration and setup
│   │   ├── app.go                       # Application bootstrap
│   │   ├── fiber.go                     # Fiber HTTP framework config
│   │   ├── gorm.go                      # GORM database config
│   │   ├── logger.go                    # Logrus logger config
│   │   ├── validator.go                 # Request validator config
│   │   └── viper.go                     # Viper configuration loader
│   │
│   ├── delivery/                         # Delivery layer (Interface Adapters)
│   │   └── http/                        # HTTP delivery
│   │       ├── address_controller.go    # Address HTTP handlers
│   │       ├── contact_controller.go    # Contact HTTP handlers
│   │       ├── user_controller.go       # User HTTP handlers
│   │       ├── middleware/              # HTTP middlewares
│   │       │   └── auth_middleware.go   # Authentication middleware
│   │       └── route/                   # Route definitions
│   │           └── route.go             # HTTP routes setup
│   │
│   ├── entity/                           # Entity layer (Enterprise Business Rules)
│   │   ├── address_entity.go            # Address domain entity
│   │   ├── contact_entity.go            # Contact domain entity
│   │   └── user_entity.go               # User domain entity
│   │
│   ├── model/                            # Models (DTOs)
│   │   ├── address_model.go             # Address request/response models
│   │   ├── contact_model.go             # Contact request/response models
│   │   ├── user_model.go                # User request/response models
│   │   ├── web_model.go                 # Common web response models
│   │   └── converter/                   # Entity-Model converters
│   │       ├── address_converter.go     # Address converter
│   │       ├── contact_converter.go     # Contact converter
│   │       └── user_converter.go        # User converter
│   │
│   ├── repository/                       # Repository layer (Data Access)
│   │   ├── address_repository.go        # Address database operations
│   │   ├── contact_repository.go        # Contact database operations
│   │   ├── user_repository.go           # User database operations
│   │   └── repository.go                # Base repository interface
│   │
│   └── usecase/                          # Use Case layer (Application Business Rules)
│       ├── address_usecase.go           # Address business logic
│       ├── contact_usecase.go           # Contact business logic
│       └── user_usecase.go              # User business logic
│
├── db/                                    # Database related files
│   └── migrations/                       # Database migrations
│       ├── 20231030144428_create_table_users.up.sql
│       ├── 20231030144428_create_table_users.down.sql
│       ├── 20231030144435_create_table_contacts.up.sql
│       ├── 20231030144435_create_table_contacts.down.sql
│       ├── 20231030144441_create_table_addresses.up.sql
│       └── 20231030144441_create_table_addresses.down.sql
│
├── docs/                                  # API documentation (auto-generated)
│   ├── docs.go                           # Swagger documentation
│   ├── swagger.json                      # OpenAPI JSON spec
│   └── swagger.yaml                      # OpenAPI YAML spec
│
├── test/                                  # Test files
│   ├── address_test.go                   # Address endpoint tests
│   ├── contact_test.go                   # Contact endpoint tests
│   ├── user_test.go                      # User endpoint tests
│   ├── helper_test.go                    # Test helpers
│   ├── init.go                           # Test initialization
│   ├── manual.http                       # Manual HTTP requests (REST Client)
│   └── http-client.env.json              # HTTP client environment variables
│
├── api/                                   # API specifications (optional)
│   └── api-spec.json                     # Manual API specification
│
├── .air.toml                             # Air configuration (hot reload)
├── .gitignore                            # Git ignore rules
├── config.json                           # Application configuration
├── go.mod                                # Go module definition
├── go.sum                                # Go module checksums
├── LICENSE.txt                           # License file
├── Makefile                              # Build automation
└── README.md                             # Project documentation

```

### 📂 Layer Details

#### 1. **cmd/** - Application Entry Points
Entry points for different applications. Currently contains only the web server.

#### 2. **internal/config/** - Configuration Layer
Handles all application configuration and initialization:
- **app.go**: Bootstraps the entire application
- **fiber.go**: Configures Fiber HTTP framework
- **gorm.go**: Sets up database connections
- **logger.go**: Initializes structured logging
- **validator.go**: Configures request validation
- **viper.go**: Loads configuration from files

#### 3. **internal/delivery/** - Delivery Layer
Handles external communication (HTTP, gRPC, etc.):
- **Controllers**: Handle HTTP requests and responses
- **Middleware**: Process requests before reaching controllers
- **Routes**: Define API endpoints and their handlers

#### 4. **internal/entity/** - Entity Layer
Core business entities representing domain models:
- Independent of external frameworks
- Contains business rules and logic
- Used throughout all layers

#### 5. **internal/model/** - Model Layer
DTOs (Data Transfer Objects) for API communication:
- **Request Models**: Validate incoming data
- **Response Models**: Structure outgoing data
- **Converters**: Transform entities to models and vice versa

#### 6. **internal/repository/** - Repository Layer
Abstracts data access and persistence:
- Database operations (CRUD)
- Query building
- Transaction management
- Independent of business logic

#### 7. **internal/usecase/** - Use Case Layer
Contains application business logic:
- Orchestrates operations between layers
- Implements business rules
- Coordinates repositories
- Validates business constraints

### 🔄 Data Flow

```
HTTP Request
    ↓
Controller (Delivery Layer)
    ↓
Use Case (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database
```

## 🔧 Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show available commands |
| `make build` | Build the application |
| `make run` | Run the application |
| `make dev` | Run with hot reload |
| `make test` | Run tests |
| `make test-coverage` | Run tests with coverage |
| `make swagger` | Generate Swagger docs |
| `make migrate-up` | Run database migrations |
| `make migrate-down` | Rollback migrations |
| `make migrate-create` | Create new migration |
| `make deps` | Download dependencies |
| `make install-tools` | Install dev tools |
| `make clean` | Clean build artifacts |
| `make fmt` | Format code |
| `make lint` | Run linter |
| `make setup` | First-time project setup |

## 📝 API Endpoints

### User Endpoints

- `POST /api/users` - Register new user
- `POST /api/users/_login` - Login user
- `GET /api/users/_current` - Get current user (authenticated)
- `PATCH /api/users/_current` - Update current user (authenticated)
- `DELETE /api/users` - Logout user (authenticated)

### Contact Endpoints

- `GET /api/contacts` - List contacts with pagination (authenticated)
- `POST /api/contacts` - Create contact (authenticated)
- `GET /api/contacts/:contactId` - Get contact by ID (authenticated)
- `PUT /api/contacts/:contactId` - Update contact (authenticated)
- `DELETE /api/contacts/:contactId` - Delete contact (authenticated)

### Address Endpoints

- `GET /api/contacts/:contactId/addresses` - List addresses (authenticated)
- `POST /api/contacts/:contactId/addresses` - Create address (authenticated)
- `GET /api/contacts/:contactId/addresses/:addressId` - Get address (authenticated)
- `PUT /api/contacts/:contactId/addresses/:addressId` - Update address (authenticated)
- `DELETE /api/contacts/:contactId/addresses/:addressId` - Delete address (authenticated)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- Fikri

## 🙏 Acknowledgments

- Clean Architecture Scaffolding by PZN
- Go community for excellent libraries