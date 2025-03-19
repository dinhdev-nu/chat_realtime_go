package global

import (
	"github.com/dinhdev-nu/realtime_auth_go/pkg/logger"
	"github.com/dinhdev-nu/realtime_auth_go/setting"
	"gorm.io/gorm"
)

var ( // Biến được khới tạo để chứa các dữ liệu trong suốt quá trình chạy
	Config setting.Confg     // Chứa các dữ liệu cấu hình
	Log    *logger.LoggerZap // chứa log là zap vd global.Log.Info("message")
	Mdb    *gorm.DB          // chứa kết nối đến database
)
