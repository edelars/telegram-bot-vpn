package config

type Environment struct {
	LogLevel string `long:"log-level" env:"LOG_LEVEL" required:"false" default:"debug"`
	TgToken  string `long:"tg-token" env:"TG_TOKEN" required:"false" default:"52"`
	MysqlEnvironment
	BussinesLogic
	Admins       string `long:"admins" env:"ADMINS" required:"true" default:""`
	OurServersIP string `long:"servers-sp" env:"SERVERS_IP" required:"true" default:""`
	Psk          string `long:"psk" env:"PSK" required:"true" default:""`
	HttpPort     int    `long:"http-port" env:"HTTP_PORT" required:"true" default:"8080"`
	FreeKassa
	qiwi
}

type MysqlEnvironment struct {
	DBHost string `long:"db-host" env:"DB_HOST" required:"true"`
	//DBPort     int    `long:"db-port" env:"DB_PORT" required:"true"`
	DBName     string `long:"db-name" env:"DB_NAME" required:"true"`
	DBUser     string `long:"db-user" env:"DB_USER" required:"true"`
	DBPassword string `long:"db-pass" env:"DB_PASS" required:"true"`
	// <= 0 - no idle conn; default = 2
	//DBMaxIdleConnCount int `long:"db-max-idle" env:"DB_MAX_IDLE" required:"false" default:"2"`
	// <= 0 - unlimited (default)
	//DBMaxConnCount int `long:"db-max-conn" env:"DB_MAX_CONN" required:"false"`
	// In seconds, <= 0 - unlimited (default)
	///DBMaxConnLifetime int `long:"db-max-conn-time" env:"DB_MAX_CONN_TIME" required:"false"`
	DBHostSW string `long:"db-host-sw" env:"DB_HOST_SW" required:"true"`
	//DBPort     int    `long:"db-port" env:"DB_PORT" required:"true"`
	DBNameSW     string `long:"db-name-sw" env:"DB_NAME_SW" required:"true"`
	DBUserSW     string `long:"db-user-sw" env:"DB_USER_SW" required:"true"`
	DBPasswordSW string `long:"db-pass-sw" env:"DB_PASS_SW" required:"true"`
}

type BussinesLogic struct {
	Price01 int `long:"price01" env:"PRICE01" required:"false" default:"179"`
	Price06 int `long:"price06" env:"PRICE06" required:"false" default:"959"`
	Price12 int `long:"price12" env:"PRICE12" required:"false" default:"1799"`
}

type FreeKassa struct {
	FKMName   string `long:"fk-m-name" env:"FK_M_NAME" required:"true" default:"OrionVPN"`
	FKId      string `long:"fk-m-id" env:"FK_M_ID" required:"true" default:"14684"`
	FKSecKey1 string `long:"fk-s-key1" env:"FK_S_KEY1" required:"true" default:"}fKHssO{"`
	FKSecKey2 string `long:"fk-s-key2" env:"FK_S_KEY2" required:"true" default:"?RH_{"`
}

type qiwi struct {
	QiwiSKey      string `long:"qiwi-s-key" env:"QIWI_S_KEY" required:"true" default:"="`
	QiwiThemeCode string `long:"qiwi-siteid" env:"QIWI_SITE_ID" required:"true" default:""`
}


