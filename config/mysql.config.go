package config

import (
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func InitMysql() {
	msql := global.Config.MySql

	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, msql.Username, msql.Password, msql.Host, msql.Port, msql.Dbname)

	db, err := gorm.Open(mysql.Open(s), &gorm.Config{
		SkipDefaultTransaction: true, // Bỏ qua transaction mặc định để tăng hiệu suất
	})
	if err != nil {
		global.Log.Error("Failed to connect to MySQL" + err.Error())
		panic(err)
	}

	global.Mdb = db
	global.Log.Info("Connected to MySQL successfully")

	if err := setPoolSize(); err != nil {
		fmt.Println("::::::::::: Set pool size err: ", err)

	}

	// if err := migrateTables(); err != nil {
	// 	fmt.Println("::::::::::: Migrate tables err: ", err)
	// }
	// generatePo()
}

// InnitMysql().setPoolSize(global.Mdb) // method chaining

func setPoolSize() error {
	mdb := global.Mdb
	msql := global.Config.MySql
	sqlDB, err := mdb.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxIdleTime(time.Duration(msql.MaxIdleConns)) // Thời gian tối đa mà một kết nối có thể ở trong pool trước khi bị đóng
	sqlDB.SetMaxOpenConns(msql.MaxOpenConns)                   // Số lượng kết nối tối đa mà có thể mở
	sqlDB.SetConnMaxLifetime(time.Duration(msql.MaxLifetime))  // Thời gian tối đa mà một kết nối có thể sống

	return nil
}

func migrateTables() error {
	return global.Mdb.AutoMigrate(
		&model.GoDbUser{},
	)
}

func generatePo() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(global.Mdb)
	g.GenerateModel("go_db_user")
	g.Execute()
}
