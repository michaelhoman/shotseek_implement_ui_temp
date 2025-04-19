package config

import (
	"time"

	"github.com/michaelhoman/ShotSeek/internal/env"
)

// Config holds all application configurations
type Config struct {
	Addr          string
	Db            DBConfig
	Env           string
	ApiURL        string
	HttpsEnabled  bool
	HttpsKeyFile  string
	HttpsCertFile string
	Mail          MailConfig
	Auth          AuthConfig
}

// AuthConfig contains authentication-related settings
type AuthConfig struct {
	Basic        BasicConfig
	Token        TokenConfig
	RefreshToken TokenConfig
}

// TokenConfig defines JWT-related settings
type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
	Aud    string
}

// BasicConfig holds basic authentication credentials
type BasicConfig struct {
	User string
	Pass string
}

// MailConfig defines email-related configurations
type MailConfig struct {
	Exp time.Duration
}

// DBConfig contains database connection settings
type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// Load initializes the configuration by fetching values from environment variables
func Load() Config {
	return Config{
		Addr:          env.GetString("ADDR", ":8080"),
		ApiURL:        env.GetString("EXTERNAL_URL", "localhost:8080"),
		HttpsEnabled:  env.GetBool("HTTPS_ENABLED", false),
		HttpsKeyFile:  env.GetString("HTTPS_KEY_FILE", ""),
		HttpsCertFile: env.GetString("HTTPS_CERT_FILE", ""),
		Db: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/shotseek?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "development"),
		Mail: MailConfig{
			Exp: time.Hour * 1, // 1 hour
		},
		Auth: AuthConfig{
			Basic: BasicConfig{
				User: env.GetString("AUTH_BASIC_USER", "admin"),
				Pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			Token: TokenConfig{
				Secret: env.GetString("JWT_SIGNING_KEY", "example"),
				Exp:    time.Minute * 60, // 60 minute
				Iss:    "shotseek-auth-service",
				Aud:    "shotseek-api",
			},
			RefreshToken: TokenConfig{
				Secret: env.GetString("REFRESH_TOKEN_SECRET", "refresh-example"), // Secret for refresh token (optional)
				Exp:    time.Hour * 24 * 7,                                       // 7 days (refresh token expiration)
				Iss:    "shotseek-auth-service",
				Aud:    "shotseek-api-refresh", // Different audience for refresh token
			},
		},
	}
}
