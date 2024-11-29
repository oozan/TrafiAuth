# Secure Microservices with Traefik Gateway and Go Authentication Service

This project sets up a reverse proxy with Traefik and demonstrates how to implement an authentication system using Go, Docker, and Docker Compose. The architecture leverages Traefik as the entry point for routing, with integrated backend services secured by a Go-based authentication system. MongoDB and Redis are used for data storage and session management.

## Features

- **Traefik Reverse Proxy**: Routes HTTP and HTTPS traffic to backend services while acting as a gateway.
- **Authentication Service**: Handles user registration, login, and token validation using MongoDB and Redis.
- **Secure Backends**: Backend services accessible only through Traefik after authentication.
- **Middleware for Security**: Ensures only authenticated requests can access backend services.

## Architecture

### Services Overview

- **Traefik**: A reverse proxy for request routing, load balancing, and secure communication.
- **MongoDB**: Stores user credentials and authentication data.
- **Redis**: Manages sessions and caching for faster authentication workflows.
- **Authentication Service**: Provides endpoints for user-related operations and token handling.
- **Backend Services**: Two separate HTTP-based backends to serve distinct content securely.

### Middleware Integration

- **Authentication Middleware**: Intercepts requests to validate JWT tokens before passing them to backend services.

## Docker Compose Configuration

### Services

1. **traefik**:
   - Routes requests to backend services via predefined rules.
   - Exposes ports 80 (HTTP), 443 (HTTPS), and 8080 (for the Traefik dashboard).

2. **mongo**:
   - A database service for storing user details.
   - Exposes port 27018 and persists data with volumes.

3. **redis**:
   - Provides a caching layer for refresh tokens and session data.
   - Exposes port 6380 for Redis clients.

4. **auth-service**:
   - Developed in Go, manages authentication flows.
   - Exposes port 8080 and communicates with MongoDB and Redis.

5. **backend1** and **backend2**:
   - Two separate HTTP services, accessible only through authentication.
   - Each backend is routed via unique host rules defined in Traefik.

6. **auth-middleware**:
   - Validates JWT tokens using the authentication service's `/validate` endpoint.

### Network Configuration

- **traefik-net**: Custom Docker bridge network with a defined IP range, enabling seamless communication between containers.

## Usage

### Prerequisites

- Install Docker and Docker Compose.

### Running the Project

1. Start the containers:

   ```bash
   docker-compose up --build



# Traefik Dashboard and Authentication System

## Traefik Dashboard

    Access the Traefik dashboard:

    Open your browser and navigate to http://localhost:8080

Explore backend services (secured via authentication):
- **Backend 1**: [http://backend1.localhost](http://backend1.localhost)
- **Backend 2**: [http://backend2.localhost](http://backend2.localhost)

---

## Authentication System

The authentication service is implemented in **Go** and manages:
- User registration
- Login
- Token issuance and validation using JWT

**MongoDB** is used for storing user credentials, and **Redis** handles token sessions.

---

### Key Endpoints

#### 1. Register a New User
- **Endpoint**: `POST /register`
- **Description**: Registers a user with their email and password.
- **Request**:
    ```json
    {
      "email": "user@example.com",
      "password": "secure_password"
    }
    ```
- **Response**:
    - `201 Created`: User successfully registered.
    - `409 Conflict`: User already exists.
    - `400 Bad Request`: Invalid input.
    - `500 Internal Server Error`: Registration failed.
    Example:

```
curl -X POST http://localhost:3000/register -H "Content-Type: application/json" -d '{"email":"user@example.com", "password":"user_password"}'
 ```

#### 2. User Login
- **Endpoint**: `POST /login`
- **Description**: Authenticates the user and issues JWT tokens.
- **Request**:
    ```json
    {
      "email": "user@example.com",
      "password": "secure_password"
    }
    ```
- **Response**:
    - `200 OK`: Tokens issued and set as cookies.
    - `401 Unauthorized`: Invalid credentials.
    - `400 Bad Request`: Input validation failed.

    Example:
```
curl -X POST http://localhost:3000/login -H "Content-Type: application/json" -d '{"email":"user@example.com", "password":"user_password"}'
```


#### 3. Validate an Access Token
- **Endpoint**: `GET /validate`
- **Description**: Checks if an access token is valid.
- **Headers**:
    ```
    Authorization: Bearer <access_token>
    ```
- **Response**:
    - `200 OK`: Token is valid.
    - `401 Unauthorized`: Invalid or expired token.

    Example:
```
curl -X GET http://localhost:3000/validate -H "Authorization: Bearer <access_token>"
```

#### 4. Refresh an Access Token
- **Endpoint**: `POST /refresh`
- **Description**: Refreshes the JWT access token using the refresh token stored in cookies.
- **Request Cookies**:
    ```
    refresh_token: The refresh token issued during login.
    ```
- **Response**:
    - `200 OK`: A new access token is issued.
    - `401 Unauthorized`: Invalid or expired refresh token.
    - `500 Internal Server Error`: Error during token refresh.

    Example:
 ```
curl -X POST http://localhost:3000/refresh --cookie "refresh_token=<refresh_token>"
```

---

### Explanation of Key Components

#### Environment Variables
- **JWT_KEY**: Secret key for signing JWT tokens.
- **MONGODB_URL**: MongoDB connection string.
- **REDIS_ADDR**: Redis connection string.
- **PORT**: Port on which the authentication service listens.

#### Database Connections
- **MongoDB**: Connection is established via the `InitMongoDB` function in the authentication service.
- **Redis**: Connection is managed using the `InitRedis` function for caching refresh tokens.

#### Handlers
- **RegisterHandler**: Handles user registration by hashing passwords and storing user data in MongoDB.
- **LoginHandler**: Authenticates users, generates JWT tokens, and stores refresh tokens in Redis.
- **ValidateHandler**: Validates JWT access tokens to ensure they are active and correct.
- **RefreshHandler**: Generates new JWT access tokens using stored refresh tokens.

---

### JWT Token Management

#### Generation
- Uses **HS256** signing method to create tokens with embedded claims.
- Tokens contain user information (e.g., email) and an expiration time.

#### Validation
- Parses and validates the token using the shared secret key.

---

### Middleware Integration with Traefik
- **Traefik** utilizes the `/validate` endpoint of the authentication service to validate JWT tokens.
- If the validation fails, Traefik denies access to the backend services.
- Requests that pass validation are forwarded to the respective backend services.

---

## Conclusion

This setup ensures secure communication between clients and backend services by integrating **Traefik**, Go-based authentication, and Dockerized infrastructure. It serves as a robust example of how to build secure microservices with token-based authentication.
