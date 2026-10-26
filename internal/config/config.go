package config

import (
	"fmt"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	OTP     OTPConfig
	SMTP    SMTPConfig
	Storage StorageConfig
	CORS    CORSConfig
	Swagger SwaggerConfig
	Google  GoogleConfig
}

type AppConfig struct {
	Env         string `env:"APP_ENV" env-default:"development"`
	Name        string `env:"APP_NAME" env-default:"coursepack"`
	Port        string `env:"APP_PORT" env-default:"8000"`
	Host        string `env:"APP_HOST" env-default:"localhost"`
	AppURL      string `env:"APP_URL" env-default:"http://localhost:8000"`
	FrontendURL string `env:"FRONTEND_URL" env-default:"http://localhost:3000"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"debug"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:""`
	Name     string `env:"DB_NAME" env-default:"coursepack"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
	PoolSize int    `env:"REDIS_POOL_SIZE" env-default:"10"`
}

type JWTConfig struct {
	Secret                   string `env:"JWT_SECRET" env-required:"true"`
	ExpirationMinutes        int    `env:"JWT_EXPIRATION_MINUTES" env-default:"15"`
	RefreshExpirationMinutes int    `env:"JWT_REFRESH_EXPIRATION_MINUTES" env-default:"30"`
}

type OTPConfig struct {
	ExpirationMinutes int `env:"OTP_EXPIRATION_MINUTES" env-default:"10"`
	MaxAttempts       int `env:"OTP_MAX_ATTEMPTS" env-default:"5"`
}

type SMTPConfig struct {
	Host      string `env:"SMTP_HOST" env-default:"localhost"`
	Port      int    `env:"SMTP_PORT" env-default:"1025"`
	User      string `env:"SMTP_USER" env-default:""`
	Password  string `env:"SMTP_PASSWORD" env-default:""`
	FromEmail string `env:"SMTP_FROM_EMAIL" env-default:"noreply@example.com"`
	FromName  string `env:"SMTP_FROM_NAME" env-default:"Coursepack Support"`
}

type StorageConfig struct {
	Type        string `env:"STORAGE_TYPE" env-default:"local"`
	LocalPath   string `env:"STORAGE_LOCAL_PATH" env-default:"./uploads"`
	MaxFileSize int64  `env:"STORAGE_MAX_FILE_SIZE_BYTES" env-default:"20971520"` // 20 MB in bytes
}

type CORSConfig struct {
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" env-default:"http://localhost:8000"`
	AllowedMethods string `env:"CORS_ALLOWED_METHODS" env-default:"GET,POST,PATCH,DELETE,OPTIONS"`
	AllowedHeaders string `env:"CORS_ALLOWED_HEADERS" env-default:"Origin,Content-Type,Accept,Authorization"`
}

type SwaggerConfig struct {
	Enabled bool `env:"SWAGGER_ENABLED" env-default:"true"`
}

type GoogleConfig struct {
	ClientID     string `env:"GOOGLE_CLIENT_ID"     env-default:""`
	ClientSecret string `env:"GOOGLE_CLIENT_SECRET" env-default:""`
	RedirectURL  string `env:"GOOGLE_REDIRECT_URL"  env-default:""`
	FrontendURL  string `env:"GOOGLE_FRONTEND_URL"  env-default:""`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment: %w", err)
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET must not be empty")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate rejects settings that would only fail silently once deployed.
func (c *Config) Validate() error {
	if c.App.Env != "production" {
		return nil
	}
	if strings.HasSuffix(c.SMTP.FromEmail, "@example.com") {
		return fmt.Errorf("SMTP_FROM_EMAIL still points at example.com")
	}
	if c.CORS.AllowedOrigins == "*" {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS must list explicit origins in production")
	}
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode,
	)
}
