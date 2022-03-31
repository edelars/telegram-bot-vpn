package config

type Environment struct {
	LogLevel string `long:"log-level" env:"LOG_LEVEL" required:"false" default:"debug"`
	TgToken  string `long:"tg-token" env:"TG_TOKEN" required:"false" default:"5284896259:AAHOGV3H_46GvzjXEFEzVIf3hXnnw3aNxFo"`
	MysqlEnvironment
	BussinesLogic
	Admins       string `long:"admins" env:"ADMINS" required:"true" default:"221953723"`
	OurServersIP string `long:"servers-sp" env:"SERVERS_IP" required:"true" default:"65.108.96.44"`
	Psk          string `long:"psk" env:"PSK" required:"true" default:"m9z6v3"`
	HttpPort     int    `long:"http-port" env:"HTTP_PORT" required:"true" default:"8080"`
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
	// In seconds, <= 0 - unlimited (default)  5212629202:AAEDdzItq_3pcp2V1Yq8L8KSIkEIGvaDcH0
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

//}fTNqXgQ*KHssO{ key1
//	?RL[Jz4r(P,*H_{ key2
