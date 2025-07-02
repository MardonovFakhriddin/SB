# Simple Bank

Simple Bank is a RESTful API for a banking system built with Go, using PostgreSQL as the database, and leveraging the `sqlx` library for database interactions and `gin-gonic` for HTTP routing. The API is documented using Swagger for easy exploration and testing.

## Table of Contents
- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Running the Application](#running-the-application)
- [Swagger Documentation](#swagger-documentation)
- [Contributing](#contributing)
- [License](#license)

## Features
- User authentication (sign-up, sign-in)
- User management (update, delete, restore, find by name, get inactive users)
- Account management (create, update, delete, get by ID, get by user ID, get by currency, get inactive accounts, check balance)
- Money transfers between accounts
- API documentation with Swagger
- Token-based authentication (Bearer token)

## Technologies
- **Go**: Programming language for the backend
- **PostgreSQL**: Relational database for persistent storage
- **sqlx**: Go library for database operations
- **gin-gonic**: HTTP web framework for routing
- **Swagger**: API documentation and testing
- **Logger**: Custom logging for error tracking

## Installation
1. Ensure you have [Go](https://golang.org/dl/) (version 1.16 or higher) and [PostgreSQL](https://www.postgresql.org/download/) installed.
2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd simple-bank
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Set up the PostgreSQL database and update the configuration with your database credentials (see [Configuration](#configuration)).

## Configuration
Create a configuration file (e.g., `config.yaml`) in the `internal/configs` directory with the following structure:

```yaml
app:
  port: ":8080"
database:
  host: "localhost"
  port: 5432
  user: "your_username"
  password: "your_password"
  dbname: "simple_bank"
```

Update the database credentials as needed. The application reads these settings to connect to PostgreSQL and run the server.

## API Endpoints
The API is organized into several groups: general, authentication, users, accounts, and transfers. All endpoints except `/` and `/auth/*` require a Bearer token for authentication.

### General
- `GET /`: Ping the server to check if it's running.
  - Response: `{"message": "Server is up and running"}`

### Authentication
- `POST /auth/sign-up`: Register a new user.
- `POST /auth/sign-in`: Authenticate a user and return a token.

### Users (Authenticated)
- `PATCH /users`: Update user details.
- `DELETE /users/:id`: Delete a user by ID.
- `GET /users/:id`: Get user details by ID.
- `GET /users/inactive`: Get a list of inactive users.
- `POST /users/restore`: Restore a deleted user.
- `GET /users/find`: Find users by name.

### Accounts (Authenticated)
- `POST /accounts`: Create a new account.
- `PATCH /accounts`: Update account details.
- `DELETE /accounts/:id`: Delete an account by ID.
- `GET /accounts/:id`: Get account details by ID.
- `GET /accounts/users/:id`: Get accounts for a specific user ID.
- `GET /accounts/inactive`: Get a list of inactive accounts.
- `GET /accounts/currency`: Get accounts by currency.
- `GET /accounts/:id/balance`: Get the balance of an account.

### Transfers (Authenticated)
- `POST /transfers`: Create a money transfer between accounts.

## Running the Application
1. Ensure the PostgreSQL database is running and configured.
2. Run the application:
   ```bash
   go run .
   ```
3. The server will start on the port specified in the configuration (default: `:8080`).

## Swagger Documentation
The API includes Swagger documentation for easy exploration. Access it at:
```
http://localhost:8080/swagger/index.html
```
This provides an interactive interface to test all endpoints.

## Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Make your changes and commit (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

## License
This project is licensed under the Apache 2.0 License. See the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0.html) file for details.