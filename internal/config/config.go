package config

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	GatewayPort     int    `yaml:"gateway_port"`     // http gateway server port
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // in ms
	EnableGateway   bool   `yaml:"enable_gateway"`   // turn on http gateway server
	EnableSwagger   bool   `yaml:"enable_swagger"`
}

func NewConfig() *Config {
	return &Config{}
}
