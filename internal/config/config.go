package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App
	HTTP
	Postgres
	JWT
	Telegram
	Tracing
}

type App struct {
	Env      string
	LogLevel string
	NodeID   string
}

type HTTP struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type Postgres struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWT struct {
	Secret         string
	AccessTokenTTL time.Duration
}

type Telegram struct {
	APIID   int
	APIHash string
}

type Tracing struct {
	Enabled      bool
	ServiceName  string
	OTLPEndpoint string
}

func Load() (*Config, error) {
	httpPort, err := getInt("HTTP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP port: %v", err)
	}

	readTimeout, err := getDuration("HTTP_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP read timeout: %v", err)
	}

	writeTimeout, err := getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP request timeout: %v", err)
	}

	shutdownTimeout, err := getDuration("HTTP_SHUTDOWN_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP shutdown timeout: %v", err)
	}

	maxOpenConns, err := getInt("POSTGRES_MAX_OPEN_CONNS", 20)
	if err != nil {
		return nil, fmt.Errorf("error parsing max open connections: %v", err)
	}

	maxIdleConns, err := getInt("POSTGRES_MAX_IDLE_CONNS", 5)
	if err != nil {
		return nil, fmt.Errorf("error parsing max idle connections: %v", err)
	}

	connMaxLifetime, err := getDuration("POSTGRES_CONN_MAX_LIFETIME", time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error parsing conn max lifetime: %v", err)
	}

	accessTokenTTL, err := getDuration("JWT_ACCESS_TOKEN_TTL", 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT access token TTL: %v", err)
	}

	apiID, err := getInt("TELEGRAM_API_ID", 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing TELEGRAM_API_ID: %v", err)
	}

	tracingEnabled, err := getBool("TRACING_ENABLED", true)
	if err != nil {
		return nil, fmt.Errorf("error parsing TRACING_ENABLED: %v", err)
	}

	cfg := Config{
		App: App{
			Env:      getString("APP_ENV", "local"),
			LogLevel: getString("APP_LOG_LEVEL", "info"),
			NodeID:   getString("APP_NODE_ID", hostname()),
		},
		HTTP: HTTP{
			Host:            getString("HTTP_HOST", "0.0.0.0"),
			Port:            httpPort,
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Postgres: Postgres{
			DSN:             getString("POSTGRES_DSN", ""),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		JWT: JWT{
			Secret:         getString("JWT_SECRET", ""),
			AccessTokenTTL: accessTokenTTL,
		},
		Telegram: Telegram{
			APIID:   apiID,
			APIHash: getString("TELEGRAM_API_HASH", ""),
		},
		Tracing: Tracing{
			Enabled:      tracingEnabled,
			ServiceName:  getString("TRACING_SERVICE_NAME", "telegram-message-collector"),
			OTLPEndpoint: getString("TRACING_OTLP_ENDPOINT", ""),
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) Validate() error {
	if c.Postgres.DSN == "" {
		return fmt.Errorf("POSTGRES_DSN is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Telegram.APIID == 0 {
		return fmt.Errorf("TELEGRAM_API_ID is required")
	}
	if c.Telegram.APIHash == "" {
		return fmt.Errorf("TELEGRAM_API_HASH is required")
	}
	return nil
}

func getString(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return n, nil
}

func getBool(key string, fallback bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s: %w", key, err)
	}
	return b, nil
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return d, nil
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
