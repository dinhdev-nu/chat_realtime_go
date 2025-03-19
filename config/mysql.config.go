package config

import (
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	msql := global.Config.MySql

	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, msql.Username, msql.Password, msql.Host, msql.Port, msql.Dbname)

	db, err := gorm.Open(mysql.Open(s), &gorm.Config{
		SkipDefaultTransaction: true, // close transaction
	})
	if err != nil {
		global.Log.Error("Failed to connect to MySQL" + err.Error())
		panic(err)
	}

	global.Mdb = db
	global.Log.Info("Connected to MySQL successfully")

	setPoolSize()
}

// InnitMysql().setPoolSize(global.Mdb) // method chaining

func setPoolSize() {
	mdb := global.Mdb
	msql := global.Config.MySql
	sqlDB, err := mdb.DB()
	if err != nil {
		fmt.Println("Mysql pool err: ", err)
	}
	sqlDB.SetConnMaxIdleTime(time.Duration(msql.MaxIdleConns)) // Thời gian tối đa mà một kết nối có thể ở trong pool trước khi bị đóng
	sqlDB.SetMaxOpenConns(msql.MaxOpenConns)                   // Số lượng kết nối tối đa mà có thể mở
	sqlDB.SetConnMaxLifetime(time.Duration(msql.MaxLifetime))  // Thời gian tối đa mà một kết nối có thể sống
}
