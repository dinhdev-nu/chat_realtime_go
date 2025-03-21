# realtime_auth_go
![Go Version](https://img.shields.io/badge/Go-v1.23-blue)


**Features**
-- Authentication:
++  JWT Authentication
++  OAuth2 Integration
-- Realtime:
++  User Chat
++  Chatbot
**Purpose** 
-- Learn Golang: Understand and practice backend development using Golang.
-- Modern Technologies: Implement JWT, OAuth2, and realtime features.
-- Hands-On Experience: Build a working system and enhance your development skills.


## Init 
### Framework 
-- GIN Framework 
<code>$ go get -u github.com/gin-gonic/gin</code>

**RUN**
<code>make run</code>
or 
<code>go run ./cmd/server/main.go</code>
...

## Database 
### MySql 
**Auth service**
-- GORM: <code>go get -u gorm.io/gorm</code> 

**Chat** 
-- Sqlc 
### Redis
<code>go get github.com/redis/go-redis/v9</code>

### Interface

**Wire**
-- clean up dependencies from params to the function to use
<code></code>


-- Handle error, response ...
-- 

**Logger**
-- Zap 
<code>go get -u go.uber.org/zap</code>
-- Segment Log file
<code>go get -u github.com/natefinch/lumberjack</code>

**Environment**
--viper
<code>go get github.com/spf13/viper</code>
