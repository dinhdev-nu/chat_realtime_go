package config

import (
	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/logger"
)

func InitLogger() {
	global.Log = logger.NewLogger(global.Config.Logger)
}