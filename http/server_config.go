package edathttp

import "time"

type ServerCfg struct {
	Port              string        `envconfig:"PORT" default:":80"`
	CertPath          string        `envconfig:"CERT_PATH"`
	KeyPath           string        `envconfig:"KEY_PATH"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"1s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"1s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"30s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"2s"`
	RequestTimeout    time.Duration `envconfig:"REQUEST_TIMEOUT" default:"60s"`
}

type CorsCfg struct {
	Origins          []string `envconfig:"ORIGINS" default:"*"`
	AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS" default:"true"`
	MaxAge           int      `envconfig:"MAX_AGE" default:"300"`
}

