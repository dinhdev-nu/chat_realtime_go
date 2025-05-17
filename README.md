# realtime_auth_go
![Go Version](https://img.shields.io/badge/Go-v1.23-blue)

## Features
- **Authentication**:
  - JWT Authentication
  - OAuth2 Integration
- **Realtime**:
  - User Chat
  - Chatbot

## Purpose
- **Learn Golang**: Understand and practice backend development using Golang.
- **Modern Technologies**: Implement JWT, OAuth2, and realtime features.
- **Hands-On Experience**: Build a working system and enhance your development skills.

## Init 
### Framework 
- **GIN Framework**  
  ```bash
  go get -u github.com/gin-gonic/gin
  ```

### Run
```bash
make run
# or
go run ./cmd/server/main.go
```

## Database 
### MySQL 
- **Auth Service**  
  - GORM:  
    ```bash
    go get -u gorm.io/gorm
    go get -u gorm.io/gen
    ```
    Generate MySQL -> PO struct in Go.

- **Chat**  
  - Sqlc  

### Redis
```bash
go get github.com/redis/go-redis/v9
```

## Additional Tools
### Wire
- Handle error, response, etc.

### JWT
```bash
go get -u github.com/golang-jwt/jwt/v5
```

### Logger
- **Zap**  
  ```bash
  go get -u go.uber.org/zap
  ```
- **Segment Log File**  
  ```bash
  go get -u github.com/natefinch/lumberjack
  ```

### Environment
- **Viper**  
  ```bash
  go get github.com/spf13/viper
  ```
