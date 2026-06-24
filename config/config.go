// Package config centraliza la configuración de la aplicación usando Viper.
// Sigue la convención de Twelve-Factor App: configuración desde archivo YAML
// sobreescribible con variables de entorno (viper.AutomaticEnv).
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN retorna la cadena de conexión para pgxpool.
func (d DatabaseConfig) DSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

// URL retorna la misma cadena para golang-migrate (espera el mismo formato).
func (d DatabaseConfig) URL() string {
	return d.DSN()
}

// Load lee el archivo YAML en la ruta dada y lo mapea a Config.
// Si DATABASE_URL está definida (Railway, Render, etc.), parsea la URL
// y sobreescribe los valores individuales de base de datos.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var cfgNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &cfgNotFound) {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		parsed, err := url.Parse(dbURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
		}

		portStr := parsed.Port()
		port := 5432
		if portStr != "" {
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return nil, fmt.Errorf("invalid port in DATABASE_URL: %w", err)
			}
		}

		sslMode := parsed.Query().Get("sslmode")
		if sslMode == "" {
			sslMode = "require"
		}

		dbName := strings.TrimLeft(parsed.Path, "/")

		cfg.Database = DatabaseConfig{
			Host:    parsed.Hostname(),
			Port:    port,
			DBName:  dbName,
			SSLMode: sslMode,
		}

		if parsed.User != nil {
			cfg.Database.User = parsed.User.Username()
			if pw, ok := parsed.User.Password(); ok {
				cfg.Database.Password = pw
			}
		}
	}

	return &cfg, nil
}
