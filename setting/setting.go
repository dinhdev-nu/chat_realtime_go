package setting

type Confg struct {
	Server Server `mapstructure:"server"`
	MySql  MySql  `mapstructure:"mysql"`
	Redis  Redis  `mapstructure:"redis"`
	Logger Logger `mapstructure:"log"`
	Jwt    Jwt    `mapstructure:"jwt"`
	Mail   Mail   `mapstructure:"mail"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Dbname       string `mapstructure:"dbname"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxLifetime  int    `mapstructure:"maxLifetime"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type Logger struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"maxsize"`
	MaxBackups int    `mapstructure:"maxbackups"`
	MaxAge     int    `mapstructure:"maxage"`
	Compress   bool   `mapstructure:"compress"`
}

type Jwt struct {
	JwtExpireTime int64  `mapstructure:"JwtExpireTime"`
	JwtSecret     string `mapstructure:"JwtSecret"`
}

type Mail struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	From     string `mapstructure:"from"`
	Password string `mapstructure:"password"`
}
