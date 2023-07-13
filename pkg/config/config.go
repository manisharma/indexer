package config

import (
	"fmt"
	"net/url"

	"github.com/ardanlabs/conf/v2"
	"github.com/joho/godotenv"
)

type AppCfg struct {
	ClientURL string `conf:"required"`
	Postgres  PgCfg
}

func Parse() (*AppCfg, error) {
	cfg := AppCfg{}
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("godotenv.Load() failed: %w", err)
	}
	_, err = conf.Parse("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// represents postgres db connection credentials
type PgCfg struct {
	Host       string `conf:"required"`
	Name       string `conf:"required"`
	User       string `conf:"required"`
	Password   string `conf:"required"`
	DisableTLS bool   `conf:"required"`
}

// create postgres db connection string
func (p PgCfg) String() string {
	sslMode := "require"
	if p.DisableTLS {
		sslMode = "disable"
	}
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.User, p.Password),
		Host:     p.Host,
		Path:     p.Name,
		RawQuery: q.Encode(),
	}
	return u.String()
}
