package global

import "github.com/dinhdev-nu/realtime_auth_go/pkg/logger"

var (
	Config Confg
	Log    *logger.LoggerZap
)

type Confg struct {
	Server Server `mapstructure:"server"`
	MySql  MySql  `mapstructure:"mysql"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type MySql struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Dbname       string `mapstructure:"dbname"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxLifetime  int    `mapstructure:"maxLifetime"`
}