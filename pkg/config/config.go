package config

type Environment struct {
	LogLevel string `long:"log-level" env:"LOG_LEVEL" required:"false" default:"debug"`
	TgToken  string `long:"tg-token" env:"TG_TOKEN" required:"false" default:"5284896259:AAHOGV3H_46GvzjXEFEzVIf3hXnnw3aNxFo"`
}
