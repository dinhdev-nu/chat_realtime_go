package test

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}  `mapstructure:"server"`
	Databases []struct {
		User string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host string `mapstructure:"host"`
	} `mapstructure:"databases"`
} 

func Viper() {
	viper:= viper.New()
	viper.AddConfigPath("./config")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	
	err:= viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config

	if err:= viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Printf("Server: %d\n", config.Server.Port)
	for _, db:= range config.Databases {
		fmt.Printf("User: %s, Password: %s, Host: %s\n", db.User, db.Password, db.Host)
	}

}
