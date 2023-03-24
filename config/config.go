package config

type Config struct {
	DB struct {
		Host     string `default:"db"`
		Name     string `default:"dtv"`
		User     string `default:"dtv-discord" env:"DB_USER"`
		Password string `default:"dtv-discord" env:"DB_PASSWORD"`
	}
	Discord struct {
		Token string `required:"true" env:"DISCORD_TOKEN"`
	}
	Mirakc struct {
		Host string `default:"tuner" env:"MIRAKC_HOST"`
		Port uint   `default:"40772" env:"MIRAKC_PORT"`
	}
	Log struct {
		Level string `default:"INFO" env:"LOG_LEVEL"`
	}
}
