# realtime_chat_app (Go) ![Go Version](https://img.shields.io/badge/Go-v1.23-blue)

A real-time chat application built with Golang, GIN, GORM, Sqlc, Redis, Mysql and Websockets.

## Features

- JWT authentication with email OTP verification
- Real-time chat (1-1 and group) via WebSocket
- Redis caching for sessions or online presence
- SQLC for type-safe database operations
- RESTful API endpoints built with Gin

## Tech Stack

- **Golang**: - Programming language.
- **Go** – Core programming language
- **Gin** – HTTP framework
- **GORM** – ORM for MySQL
- **SQLC** – SQL → Go code generation
- **Redis** – In-memory cache
- **WebSocket** – Real-time bi-directional communication
- **JWT** – Token-based authentication

## Project Structure

```
chat_realtime_go/
├── cmd/
│ └── server/ # Main server entry point
├── config/ # Configuration loader (MySQL, Redis, etc.)
├── db/ # SQL schema, migrations
├── global/ # Global variables/constants
├── internal/
│ ├── controller/ # HTTP + WebSocket handlers
│ ├── dto/ # Request/response DTOs
│ ├── model/ # GORM models
│ ├── repository/ # DB interaction logic
│ ├── router/ # Gin router setup
│ ├── service/ # Business logic
│ ├── websocket/ # WebSocket core logic
│ └── utils/ # Helper utilities
├── pkg/
│ ├── logger/ # Logging setup
│ ├── middleware/ # Custom middleware (auth, logging)
│ └── response/ # Standardized API response format
├── settings/ # Env variable binding
├── sqlc/ # SQLC-generated query code
├── templates-email/ # Email OTP HTML templates
├── test/ # Unit and integration tests
├── .env.example # Sample environment config
├── sqlc.yaml # SQLC config file
├── go.mod
├── go.sum
└── README.md
```

## Purpose

- **Learn Golang**: Understand and practice backend development using Golang.
- **Learn Websockets**: Implement real-time communication features.
- **Modern Technologies**: Implement JWT, OAuth2, and realtime features.
- **Hands-On Experience**: Build a working system and enhance your development skills.

## Getting Started

```bash
git clone https://github.com/dinhdev-nu/chat_realtime_go.git
cd chat_realtime_go
cp .env.example .env
go mode download
sqlc generate
go run ./cmd/server/main.go
```
