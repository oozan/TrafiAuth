Project: Secure Microservices with Traefik Gateway and Go Authentication Service

This project sets up a reverse proxy with Traefik and demonstrates how to implement an authentication system using Go, Docker, and Docker Compose. The architecture leverages Traefik as the entry point for routing, with integrated backend services secured by a Go-based authentication system. MongoDB and Redis are used for data storage and session management, respectively.
Features

    Traefik Reverse Proxy: Routes HTTP and HTTPS traffic to backend services while acting as a gateway.
    Authentication Service: Handles user registration, login, and token validation using MongoDB and Redis.
    Secure Backends: Backend services accessible only through Traefik after authentication.
    Middleware for Security: Ensures only authenticated requests can access backend services.

Architecture
Services Overview

    Traefik: A reverse proxy for request routing, load balancing, and secure communication.
    MongoDB: Stores user credentials and authentication data.
    Redis: Manages sessions and caching for faster authentication workflows.
    Authentication Service: Provides endpoints for user-related operations and token handling.
    Backend Services: Two separate HTTP-based backends to serve distinct content securely.

Middleware Integration

    Authentication Middleware: Intercepts requests to validate JWT tokens before passing them to backend services.

Docker Compose Configuration
Services

    traefik:
        Routes requests to backend services via predefined rules.
        Exposes ports 80 (HTTP), 443 (HTTPS), and 8080 (for the Traefik dashboard).

    mongo:
        A database service for storing user details.
        Exposes port 27018 and persists data with volumes.

    redis:
        Provides a caching layer for refresh tokens and session data.
        Exposes port 6380 for Redis clients.

    auth-service:
        Developed in Go, manages authentication flows.
        Exposes port 8080 and communicates with MongoDB and Redis.

    backend1 and backend2:
        Two separate HTTP services, accessible only through authentication.
        Each backend is routed via unique host rules defined in Traefik.

    auth-middleware:
        Validates JWT tokens using the authentication service's /validate endpoint.

Network Configuration

    traefik-net: Custom Docker bridge network with a defined IP range, enabling seamless communication between containers.

Usage
Prerequisites

    Install Docker and Docker Compose.

Running the Project

    Start the containers:

    docker-compose up --build

    Access the Traefik dashboard:

    Traefik Dashboard: http://localhost:8080

    Explore backend services (secured via authentication):
        Backend 1: http://backend1.localhost
        Backend 2: http://backend2.localhost

Authentication System

The authentication service is implemented in Go and manages user registration, login, token issuance, and validation using JWT. MongoDB stores user credentials, and Redis handles token sessions.
Key Endpoints
1. Register a New User

    Endpoint: POST /register

    Description: Registers a user with their email and password.

    Request:

    {
      "email": "user@example.com",
      "password": "secure_password"
    }

    Response:
        201 Created: User successfully registered.
        409 Conflict: User already exists.
        400 Bad Request: Invalid input.
        500 Internal Server Error: Registration failed.

2. User Login

    Endpoint: POST /login

    Description: Authenticates the user and issues JWT tokens.

    Request:

    {
      "email": "user@example.com",
      "password": "secure_password"
    }

    Response:
        200 OK: Tokens issued and set as cookies.
        401 Unauthorized: Invalid credentials.
        400 Bad Request: Input validation failed.

3. Validate an Access Token

    Endpoint: GET /validate

    Description: Checks if an access token is valid.

    Headers:

    Authorization: Bearer <access_token>

    Response:
        200 OK: Token is valid.
        401 Unauthorized: Token is expired or invalid.

4. Refresh an Access Token

    Endpoint: POST /refresh

    Description: Renews an expired access token using a valid refresh token.

    Cookies:

    refresh_token=<refresh_token>

    Response:
        200 OK: New access token issued.
        401 Unauthorized: Invalid or expired refresh token.

Example Usage

    Register a User:

curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"email":"user@example.com", "password":"secure_password"}'

Login:

curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"email":"user@example.com", "password":"secure_password"}'

Validate a Token:

curl -X GET http://localhost:8080/validate \
-H "Authorization: Bearer <access_token>"

Refresh a Token:

    curl -X POST http://localhost:8080/refresh --cookie "refresh_token=<refresh_token>"

Environment Variables
Variable	Description
MONGODB_URL	MongoDB connection string.
REDIS_ADDR	Redis connection string.
JWT_KEY	Secret key for signing tokens.
PORT	Port for the auth-service.
Traefik Configuration

    EntryPoints:
        web: HTTP traffic on port 80.
        websecure: HTTPS traffic on port 443.

    Middleware:
        Authentication middleware uses the auth-service /validate endpoint for validating JWT tokens.

Conclusion

This project is a scalable and secure solution for managing microservices with Traefik and Go. By leveraging token-based authentication, it ensures that only authorized users can access backend services, making it ideal for modern web applications requiring robust security.