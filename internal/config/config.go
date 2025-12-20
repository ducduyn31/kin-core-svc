package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	S3         S3Config         `mapstructure:"s3"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Telemetry  TelemetryConfig  `mapstructure:"telemetry"`
	Presence   PresenceConfig   `mapstructure:"presence"`
	Pagination PaginationConfig `mapstructure:"pagination"`
}

type ServerConfig struct {
	Port             int           `mapstructure:"port"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout  time.Duration `mapstructure:"shutdown_timeout"`
	EnableReflection bool          `mapstructure:"enable_reflection"` // Enable gRPC reflection (development only)
}

type DatabaseConfig struct {
	WriteURL        string        `mapstructure:"write_url"`
	ReadURL         string        `mapstructure:"read_url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type RedisConfig struct {
	URL          string        `mapstructure:"url"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type S3Config struct {
	Endpoint     string `mapstructure:"endpoint"`
	Region       string `mapstructure:"region"`
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	Bucket       string `mapstructure:"bucket"`
	UsePathStyle bool   `mapstructure:"use_path_style"`
}

type AuthConfig struct {
	Domain      string `mapstructure:"domain"`
	Audience    string `mapstructure:"audience"`
	TokenLookup string `mapstructure:"token_lookup"`
	AuthScheme  string `mapstructure:"auth_scheme"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type TelemetryConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	OTLPEndpoint   string `mapstructure:"otlp_endpoint"`
}

type PresenceConfig struct {
	TTL             time.Duration `mapstructure:"ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type PaginationConfig struct {
	DefaultLimit int `mapstructure:"default_limit"`
	MaxLimit     int `mapstructure:"max_limit"`
}

func Load() (*Config, error) {
	env := os.Getenv("KIN_ENV")
	if env == "" {
		env = "development"
	}

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read %s config: %w", env, err)
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindEnvVars(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	setDefaults(&cfg)

	return &cfg, nil
}

// bindEnvVars binds specific environment variables to viper configuration keys.
// It maps DATABASE_WRITE_URL/DATABASE_READ_URL -> database.write_url/read_url; REDIS_URL -> redis.url; S3_ENDPOINT/S3_REGION/S3_ACCESS_KEY/S3_SECRET_KEY/S3_BUCKET -> s3.endpoint/region/access_key/secret_key/bucket; AUTH0_DOMAIN/AUTH0_AUDIENCE -> auth.domain/audience; PORT -> server.port; and OTEL_ENABLED/OTEL_EXPORTER_OTLP_ENDPOINT -> telemetry.enabled/telemetry.otlp_endpoint.
func bindEnvVars(v *viper.Viper) {
	_ = v.BindEnv("database.write_url", "DATABASE_WRITE_URL")
	_ = v.BindEnv("database.read_url", "DATABASE_READ_URL")

	_ = v.BindEnv("redis.url", "REDIS_URL")

	_ = v.BindEnv("s3.endpoint", "S3_ENDPOINT")
	_ = v.BindEnv("s3.region", "S3_REGION")
	_ = v.BindEnv("s3.access_key", "S3_ACCESS_KEY")
	_ = v.BindEnv("s3.secret_key", "S3_SECRET_KEY")
	_ = v.BindEnv("s3.bucket", "S3_BUCKET")

	_ = v.BindEnv("auth.domain", "AUTH0_DOMAIN")
	_ = v.BindEnv("auth.audience", "AUTH0_AUDIENCE")

	_ = v.BindEnv("server.port", "PORT")

	_ = v.BindEnv("telemetry.enabled", "OTEL_ENABLED")
	_ = v.BindEnv("telemetry.otlp_endpoint", "OTEL_EXPORTER_OTLP_ENDPOINT")
}

// setDefaults fills zero-valued fields of cfg with sensible defaults for server,
// database, redis, presence, pagination, authentication, S3, and telemetry
// configuration, mutating cfg in place.
func setDefaults(cfg *Config) {
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 * time.Second
	}
	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 10 * time.Second
	}

	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}

	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 10
	}
	if cfg.Redis.MaxRetries == 0 {
		cfg.Redis.MaxRetries = 3
	}

	if cfg.Presence.TTL == 0 {
		cfg.Presence.TTL = 5 * time.Minute
	}
	if cfg.Presence.CleanupInterval == 0 {
		cfg.Presence.CleanupInterval = 1 * time.Minute
	}

	if cfg.Pagination.DefaultLimit == 0 {
		cfg.Pagination.DefaultLimit = 20
	}
	if cfg.Pagination.MaxLimit == 0 {
		cfg.Pagination.MaxLimit = 100
	}

	if cfg.Auth.TokenLookup == "" {
		cfg.Auth.TokenLookup = "header:Authorization"
	}
	if cfg.Auth.AuthScheme == "" {
		cfg.Auth.AuthScheme = "Bearer"
	}

	if cfg.S3.Region == "" {
		cfg.S3.Region = "us-east-1"
	}

	if cfg.Telemetry.ServiceName == "" {
		cfg.Telemetry.ServiceName = "kin-core-svc"
	}
	if cfg.Telemetry.ServiceVersion == "" {
		cfg.Telemetry.ServiceVersion = "0.1.0"
	}
	if cfg.Telemetry.OTLPEndpoint == "" {
		cfg.Telemetry.OTLPEndpoint = "localhost:4317"
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
}

func (c *ServerConfig) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}