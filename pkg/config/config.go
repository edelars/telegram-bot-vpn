package config

type Environment struct {
	LogLevel string `long:"log-level" env:"LOG_LEVEL" required:"false" default:"debug"`
	TgToken  string `long:"tg-token" env:"TG_TOKEN" required:"false" default:"5284896259:AAHOGV3H_46GvzjXEFEzVIf3hXnnw3aNxFo"`
	MysqlEnvironment
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
}
