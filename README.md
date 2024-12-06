# Online Store Backend

## Overview
This project implements a basic online store backend that allows users to register, log in, browse products by category, manage their shopping cart, and perform a checkout operation. The application follows a clean architecture pattern, making the code modular, testable, and maintainable. It is containerized using Docker and can be easily deployed alongside a PostgreSQL database.

---

## Architecture

**Clean Architecture Layers:**
1. **Domain (Entity/Model Layer)**  
   - Core business entities: `Customer`, `Product`, `CartItem`, `Order`
   - No external dependencies, purely business logic.

2. **Usecase (Application Logic)**  
   - Implements use cases: register/login, cart operations, product queries, checkout.
   - Interacts with domain models and repository interfaces.

3. **Repository (Data Access)**  
   - Interfaces for database operations.
   - Concrete implementations using `database/sql` and `pq` driver for PostgreSQL.

4. **Handler (Interface Adapters)**  
   - HTTP handlers using the Gin framework.
   - Translates HTTP requests to usecase inputs and outputs to HTTP responses.
   - Includes JWT-based authentication middleware.

5. **Infrastructure (Frameworks and Drivers)**  
   - Database connections, migrations, configuration (Viper), and Docker setup.

**Data Flow:**

Request → Handler → Usecase → Repository → Database Response ← Handler ← Usecase ← Repository ← Database


This ensures that the core business logic is independent of external frameworks, databases, or UI, providing a maintainable and extensible system.

---

## Features

**Authentication & User Management**
- **Register:** Create a new customer account.
- **Login:** Authenticate users and return a JWT token.

**Product Management**
- **View Products by Category:** Public endpoint to list products filtered by category.

**Shopping Cart**
- **Add to Cart (Protected):** Requires JWT, adds products to the user’s cart.
- **View Cart (Protected):** Lists all items in the authenticated user’s cart.
- **Remove from Cart (Protected):** Removes specific items from the cart.

**Checkout & Orders**
- **Checkout (Protected):** Converts cart items into an order, updates stock, sets order status to “paid” (placeholder payment).
- Database migrations and seeding scripts allow initializing schema and sample data.

**Health Check**
- **/health:** A simple endpoint to check the server’s readiness and uptime.

---

## Tech Stack

**Backend:**
- **Go (Golang):** High-performance, statically typed language.
- **Gin:** Lightweight, high-performance HTTP framework for routing, middleware.
- **JWT:** JSON Web Tokens for secure route protection.
- **PostgreSQL:** Relational database for durable, consistent storage.

**DevOps & Infrastructure:**
- **Docker & Docker Compose:** Containerize the app and the database environment.
- **golang-migrate:** Manage database schema changes through migration files.
- **Viper:** Manage configurations via `.env` and environment variables.

**Deployment:**
- **Docker Hub & Heroku:** Build and push Docker images to a registry and deploy to cloud platforms.
- **Local Development:** `docker-compose up` to start both app and database locally.

---

## Getting Started

**Prerequisites:**
- Go 1.20+
- Docker & Docker Compose installed
- (Optional) `migrate` CLI for running migrations locally

**Run Locally (Without Docker):**
1. Set environment variables in `.env` (PORT, DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET).
2. Run migrations: 
```bash
migrate -path migrations -database "postgres://admin:4444@localhost:5432/online-store?sslmode=disable" up
```
3. Seed the database 
```bash
go run scripts/seed.go
```
4. Start the server: 
```bash
go run cmd/server/main.go
```

Run with Docker:

1. Build and run: 
```bash
docker-compose up --build
```
2. Apply migrations (if using the migrate container approach): 
```bash
docker-compose run migrate
```
3. Seed the database: 
```bash
docker-compose run seed
```
4. Visit http://localhost:8080/health to confirm the server is running.
