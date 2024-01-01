package edatgrpc

type ServerCfg struct {
	Network  string `envconfig:"NETWORK" default:"tcp"`
	Address  string `envconfig:"ADDRESS" default:":8000"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

type ClientCfg struct {
	URI      string `envconfig:"URI"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}
