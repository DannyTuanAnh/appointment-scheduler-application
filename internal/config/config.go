package config

import (
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration

	MaxHeaderBytes int
}

type Config struct {
	DB     DatabaseConfig
	Server ServerConfig
}

func NewConfigServer() *Config {
	return &Config{
		Server: ServerConfig{
			Port:              utils.GetEnv("PORT", "8080"),                              // port which server listens on, default is 8080
			ReadTimeout:       utils.GetEnvTime("SV_READTIMEOUT", 5) * time.Second,       // maximum duration for reading the entire request, including the body (default is 5 seconds)
			ReadHeaderTimeout: utils.GetEnvTime("SV_READHEADERTIMEOUT", 3) * time.Second, // maximum duration for reading the headers of the request (default is 3 seconds)
			WriteTimeout:      utils.GetEnvTime("SV_WRITETIMEOUT", 10) * time.Second,     // maximum duration before timing out writes of the response (default is 10 seconds)
			IdleTimeout:       utils.GetEnvTime("SV_IDLETIMEOUT", 120) * time.Second,     // maximum amount of time to wait for the next request when keep-alives are enabled (default is 120 seconds)
			ShutdownTimeout:   utils.GetEnvTime("SV_SHUTDOWNTIMEOUT", 5) * time.Second,   // maximum amount of time to wait for the server to shutdown gracefully (default is 5 seconds)

			MaxHeaderBytes: utils.GetEnvInt("SV_MAXHEADERBYTES", 16) << 10, // maximum size of request headers in bytes (default is 16 KB)
		},
	}
}

func NewConfigDB() *Config {
	return &Config{
		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "myapp"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (c *Config) DB_DNS() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DBName, c.DB.SSLMode)
}
