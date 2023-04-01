package config

type Config struct {
	DB struct {
		Host     string `default:"db" env:"DB_HOST"`
		Name     string `default:"dtv" env:"DB_NAME"`
		User     string `default:"dtv-discord" env:"DB_USER"`
		Password string `default:"dtv-discord" env:"DB_PASSWORD"`
	}
	Redis struct {
		Host string `default:"redis" env:"REDIS_HOST"`
		Port uint   `default:"6379" env:"REDIS_PORT"`
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
	Recording struct {
		BasePath            string `required:"true" env:"RECORDING_BASE_PATH"`
		ContentPathTemplate string `required:"true" env:"CONTENT_PATH_TEMPLATE"`
	}
	Encoding struct {
		Enabled               bool   `required:"true" env:"ENCODING_ENABLED"`
		BasePath              string `required:"true" env:"ENCODING_BASE_PATH"`
		OutputPathTemplate    string `required:"true" env:"ENCODING_OUTPUT_PATH_TEMPLATE"`
		EncodeCommandTemplate string `required:"true" env:"ENCODING_COMMAND"`
	}
	Match struct {
		KanaMatch bool `default:"true" env:"KANA_MATCH"`
	}
}
