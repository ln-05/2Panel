package config

type Docker struct {
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Version   string `mapstructure:"version" json:"version" yaml:"version"`
	TLSVerify bool   `mapstructure:"tls-verify" json:"tlsVerify" yaml:"tls-verify"`
	CertPath  string `mapstructure:"cert-path" json:"certPath" yaml:"cert-path"`
	Timeout   int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}