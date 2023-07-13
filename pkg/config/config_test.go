package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPgCfg_String(t *testing.T) {
	type fields struct {
		Host       string
		Name       string
		User       string
		Password   string
		DisableTLS bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return postgres connection string with tls disabled",
			fields: fields{
				Host:       "localhost:5432",
				Name:       "database",
				User:       "user",
				Password:   "password",
				DisableTLS: true,
			},
			want: "postgres://user:password@localhost:5432/database?sslmode=disable&timezone=utc",
		},
		{
			name: "should return postgres connection string with tls enabled",
			fields: fields{
				Host:       "localhost:5432",
				Name:       "database",
				User:       "user",
				Password:   "password",
				DisableTLS: false,
			},
			want: "postgres://user:password@localhost:5432/database?sslmode=require&timezone=utc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PgCfg{
				Host:       tt.fields.Host,
				Name:       tt.fields.Name,
				User:       tt.fields.User,
				Password:   tt.fields.Password,
				DisableTLS: tt.fields.DisableTLS,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("PgCfg.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		dotEnv  []byte
		want    *AppCfg
		wantErr bool
	}{
		{
			name: "should return *AppCfg",
			dotEnv: []byte(`CLIENT_URL=https://dummy.client
			POSTGRES_HOST=localhost:5432
			POSTGRES_NAME=database
			POSTGRES_USER=user
			POSTGRES_PASSWORD=password
			POSTGRES_DISABLE_TLS=true`),
			want: &AppCfg{
				ClientURL: "https://dummy.client",
				Postgres: PgCfg{
					Host:       "localhost:5432",
					Name:       "database",
					User:       "user",
					Password:   "password",
					DisableTLS: true,
				},
			},
			wantErr: false,
		},
		{
			name:    "should return error for missing .env",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.dotEnv) > 0 {
				f, err := os.Create(".env")
				assert.Nil(t, err, "os.Create() must not fail")
				defer func() {
					f.Close()
					os.Remove(".env")
				}()
				_, err = f.Write(tt.dotEnv)
				assert.Nil(t, err, "write to .env must not fail")
			}
			got, err := Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
