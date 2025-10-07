package postgres

import "time"

const (
	defaultMaxConns          int32 = 10
	defaultMinConns          int32 = 2
	defaultMaxConnLifetime         = time.Hour
	defautlMinIdleConns      int32 = 2
	defaultMaxConnIdleTime         = 30 * time.Minute
	defaultHealthCheckPeriod       = time.Minute
	defaultSSLMode                 = "disable"
)

type Config struct {
	Host     string `mapstructure:"host" env:"DB_HOST"`
	Port     int    `mapstructure:"port" env:"DB_PORT"`
	Name     string `mapstructure:"name" env:"DB_NAME"`
	User     string `mapstructure:"user" env:"DB_USER"`
	Password string `mapstructure:"password" env:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"sslmode" env:"DB_SSLMODE"`

	// Параметры пула соединений
	MaxConns          int32         `mapstructure:"max_conns" env:"DB_MAX_CONNS"`
	MinConns          int32         `mapstructure:"min_conns" env:"DB_MIN_CONNS"`
	MinIdleConns      int32         `mapstructure:"min_idle_conns" env:"DB_MIN_IDLE_CONNS"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime" env:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time" env:"DB_MAX_CONN_IDLE_TIME"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period" env:"DB_HEALTH_CHECK_PERIOD"`
}
