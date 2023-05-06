package config

type Config struct {
	DB struct {
		Host     string `default:"db" env:"DB_HOST"`
		Port     int    `default:"3306" env:"DB_PORT"`
		Name     string `default:"dtv" env:"DB_NAME"`
		User     string `default:"dtv-discord" env:"DB_USER"`
		Password string `default:"dtv-discord" env:"DB_PASSWORD"`
	}
	Redis struct {
		Host string `default:"redis" env:"REDIS_HOST"`
		Port uint   `default:"6379" env:"REDIS_PORT"`
	}
	Meili struct {
		Host string `default:"meilisearch" env:"MEILI_HOST"`
		Port int    `default:"7700" env:"MEILI_PORT"`
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
		Ext                   string `required:"true" env:"ENCODING_EXT"`
		EncodeCommandTemplate string `required:"true" env:"ENCODING_COMMAND"`
		DeleteOriginalFile    bool   `required:"true" env:"ENCODING_DELETE_ORIGINAL_FILE"`
	}
	Transcription struct {
		Enabled    bool   `required:"true" env:"TRANSCRIPTION_ENABLED"`
		BasePath   string `required:"true" env:"TRANSCRIPTION_BASE_PATH"`
		Ext        string `required:"true" env:"TRANSCRIPTION_EXT"`
		Type       string `default:"local" env:"TRANSCRIPTION_TYPE"` // local or api
		ScriptPath string `required:"true" env:"TRANSCRIPTION_SCRIPT_PATH"`
		ModelSize  string `required:"true" env:"TRANSCRIPTION_MODEL_SIZE"`
	}
	OpenAI struct {
		Enabled bool   `default:"false" env:"PARSE_TITLE_WITH_GPT"`
		Token   string `default:"" env:"OPENAI_TOKEN"`
	}
	TaskQueue struct {
		DefaultQueueName    string `default:"default" env:"TASK_QUEUE_DEFAULT_QUEUE_NAME"`
		EncodeQueueName     string `default:"encode" env:"TASK_QUEUE_ENCODE_QUEUE_NAME"`
		TranscribeQueueName string `default:"transcribe" env:"TASK_QUEUE_TRANSCRIBE_QUEUE_NAME"`
	}
}
