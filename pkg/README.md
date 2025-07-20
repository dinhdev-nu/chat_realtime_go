## Folder PKG

-- Dành cho các package dùng chung và thường xuyên sử dụng trong sự án
++ handle error
++ handle response
++ middleware
++ logers

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
go run ./cmd/cli/make/make.go
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
    Generate MySQL Tables -> Models struct in Go.

- **Chat**
  - Sqlc : convert sql query -> raw native query in go

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

- Goose : mirgation file .sql to MySql

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

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
