package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
)

func InitMysqlc() {

	m := global.Config.MySql

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Dbname,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		global.Log.Error("Failed to connect to MySQL: " + err.Error())
		panic(err)
	}

	global.Mdbc = db
	global.Log.Info("Connected to MySQLC successfully")

	setPool()

}

func setPool() {
	m := global.Config.MySql
	sqlDB, err := global.Mdb.DB()
	if err != nil {
		fmt.Println("Failed to get sqlDB: ", err)
	}

	sqlDB.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns)) // Thời gian tối đa mà một kết nối có thể ở trong pool trước khi bị đóng
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)                   // Số lượng kết nối tối đa mà có thể mở
	sqlDB.SetConnMaxLifetime(time.Duration(m.MaxLifetime))  // Thời gian tối đa mà một kết nối có thể sống

}
