package config

import (
	"context"
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password, // no password set
		DB:       r.Database, // use default DB
		PoolSize: 10,         // connection pool size
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Log.Error("Failed to connect to Redis", zap.Error(err))
		return
	}

	global.Log.Info("Connected to Redis successfully")
	global.Rdb = rdb
	// TestRedis()
}

func TestRedis() {
	err := global.Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Print("Error 1 -------------: ", err)
		return
	}

	getVey, err := global.Rdb.Get(ctx, "key").Result()
	if err != nil {
		fmt.Print("Error 11 ---------: ", err)
		return
	}
	fmt.Println("Key:::::::::::::::::::::", getVey)

}
