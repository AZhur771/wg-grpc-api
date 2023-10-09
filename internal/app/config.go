package app

type Config struct {
	IsProduction bool     `env:"PRODUCTION"`
	Host         string   `env:"HOST" envDefault:"localhost"`
	Port         int      `env:"PORT" envDefault:"3000"`
	ServeSwagger bool     `env:"SWAGGER"`
	Tokens       []string `env:"TOKENS" envSeparator:","`

	DBHost               string `env:"DB_HOST,required"`
	DBPort               int    `env:"DB_PORT" envDefault:"5432"`
	DBName               string `env:"DB_NAME,required"`
	DBUsername           string `env:"DB_USERNAME,required"`
	DBPassword           string `env:"DB_PASSWORD,required"`
	DBTimeout            int    `env:"DB_TIMEOUT" envDefault:"5"`
	DBMaxOpenConnections int    `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"10"`
	DBMaxIdleConnections int    `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"10"`

	CaCert string `env:"CACERT"`
	Cert   string `env:"CERT"`
	Key    string `env:"KEY"`
}
