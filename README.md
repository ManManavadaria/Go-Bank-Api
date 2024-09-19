---

# Go Bank API

Go Bank API is a simple API built with Golang and PostgreSQL that simulates a basic banking system. This project allows users to create accounts, retrieve account details, update accounts, and delete accounts. JWT-based authentication is implemented for secure access.

## Features

- **Account Management**: Create, update, and delete accounts.
- **JWT Authentication**: Provides security for API access.
- **PostgreSQL Integration**: Persistent data storage using PostgreSQL.
- **RESTful API**: Standard RESTful endpoints for account operations.
- **Gorilla Mux**: Router for handling HTTP routes.

## Requirements

- Go 1.16+
- PostgreSQL
- [Gorilla Mux](https://github.com/gorilla/mux) package
- [JWT](https://github.com/golang-jwt/jwt/v5) package
- [pq](https://github.com/lib/pq) package (PostgreSQL driver)

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/ManManavadaria/Go-Bank-Api.git
   cd Go-Bank-Api
   ```

2. **Install dependencies**:

   You can install the necessary Go modules by running:

   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL**:

   Ensure that you have a running PostgreSQL instance. You can modify the connection string in the `NewPostgresConnection` function inside the `storage` package to match your PostgreSQL configuration.

   ```go
   connectionString := "user=Man_user dbname=go-bank password=password sslmode=disable"
   ```

4. **Run the application**:

   After setting up PostgreSQL, run the application by executing:

   ```bash
   go run main.go
   ```

   The server will run on port `7000` by default.

## API Endpoints

### JWT Authentication

To access the API, you need to provide a valid JWT token in the request headers.

```
x-jwt-token: <your_jwt_token>
```

### Account Endpoints

- **GET /account/{id}**: Retrieve an account by ID.
- **GET /account**: Retrieve all accounts.
- **POST /account**: Create a new account. Requires a JSON payload:
  
  ```json
  {
    "firstname": "John",
    "lastname": "Doe"
  }
  ```

- **PUT /account/{id}**: Update an account.
- **DELETE /account/{id}**: Delete an account.

## Database Schema

The PostgreSQL schema for accounts is created automatically by the application. The schema looks like this:

```sql
CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    account_number INTEGER NOT NULL UNIQUE, 
    balance REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Project Structure

```bash
.
├── main.go                   # Main entry point for the application
├── models                    # Package containing data models
│   └── account.go            # Account struct and helper functions
├── storage                   # Package containing database interaction logic
│   └── postgres.go           # PostgreSQL store and methods
├── go.mod                    # Module dependencies
├── go.sum                    # Module checksums
└── README.md                 # Project documentation
```

---
