package config

type Config struct {
	Logger LoggerConfig
	Server ServerConfig
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type ServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	GatewayPort     int    `yaml:"gateway_port"`     // http gateway server port
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // in ms
	EnableGateway   bool   `yaml:"enable_gateway"`   // turn on http gateway server
	EnableSwagger   bool   `yaml:"enable_swagger"`
	SwaggerPath     string `yaml:"swagger_path"'`
}

func NewConfig() *Config {
	return &Config{}
}
