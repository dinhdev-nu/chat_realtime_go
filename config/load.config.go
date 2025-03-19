package config

import (
	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/spf13/viper"
)

func LoadConfig() {

	viper := viper.New()
	viper.AddConfigPath("./environment")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	if err:= viper.ReadInConfig(); err != nil {
		panic("Error reading config file" + err.Error())
	}

	// chuyển dữ liệu từ file local.yaml vào biến gobal.Config
	if err:= viper.Unmarshal(&global.Config); err != nil {
		panic("Error unmarshal config" + err.Error())
	}
}