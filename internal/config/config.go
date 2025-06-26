package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	App      AppConfig      `mapstructure:"app"`
	Security SecurityConfig `mapstructure:"security"`
}

type ServerConfig struct {
	Port            string `mapstructure:"port"`
	Host            string `mapstructure:"host"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	IdleTimeout     int    `mapstructure:"idle_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
	StaticDir   string `mapstructure:"static_dir"`
	TemplateDir string `mapstructure:"template_dir"`
}

type SecurityConfig struct {
	TrustedProxies []string `mapstructure:"trusted_proxies"`
	RateLimitRPM   int      `mapstructure:"rate_limit_rpm"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("$HOME/.gohtmx")

	// Set defaults
	setDefaults()

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("GOHTMX")

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "3000")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 120)
	viper.SetDefault("server.shutdown_timeout", 30)

	// App defaults
	viper.SetDefault("app.name", "go-htmx")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.static_dir", "internal/static")
	viper.SetDefault("app.template_dir", "internal/template")

	// Security defaults
	viper.SetDefault("security.trusted_proxies", []string{})
	viper.SetDefault("security.rate_limit_rpm", 60)
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}
