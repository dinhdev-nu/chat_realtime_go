package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
)

var user model.GoDbUser

type IAuthRepo interface {
	GetExistEmail(email string) bool
	AddOtp(email string, data map[string]interface{}, ttl int64) error
	GetOtp(email string) map[string]interface{}
	DelOtp(email string) error
}

type authRepo struct {
	ctx context.Context
}

func NewAuthRepo() IAuthRepo {
	ctx := context.Background()
	return &authRepo{
		ctx: ctx,
	}
}

// mysql
func (ar *authRepo) GetExistEmail(email string) bool {
	exist := global.Mdb.Model(&user).Where("usr_email = ?", email).First(&user).RowsAffected

	return exist != 0
}

// redis
func (ar *authRepo) AddOtp(email string, data map[string]interface{}, ttl int64) error {
	key := "otp:" + email + ":usr"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return global.Rdb.SetEx(ar.ctx, key, jsonData, time.Duration(ttl)*time.Second).Err()
}

func (ar *authRepo) GetOtp(email string) map[string]interface{} {
	key := "otp:" + email + ":usr"
	fmt.Println("key: ", key)
	jsonData := global.Rdb.Get(ar.ctx, key).Val()
	if jsonData == "" {
		return nil
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil
	}

	return data
}

func (ar *authRepo) DelOtp(email string) error {
	key := "otp:" + email + ":usr"
	return global.Rdb.Del(ar.ctx, key).Err()
}
