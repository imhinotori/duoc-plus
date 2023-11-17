package config

type Config struct {
	General General `koanf:"general"`
	HTTP    HTTP    `koanf:"http"`
	JWT     JWT     `koanf:"jwt"`
	Duoc    Duoc    `koanf:"duoc"`
}

type General struct {
	Debug bool `koanf:"debug"`
}

type HTTP struct {
	Address string `koanf:"address"`
	Port    int    `koanf:"port"`
	SSL     bool   `koanf:"ssl"`
	SSLCert string `koanf:"ssl-cert"`
	SSLKey  string `koanf:"ssl-key"`
}

type JWT struct {
	PrivateKey   string `koanf:"private-key"`
	PublicKey    string `koanf:"public-key"`
	AutoGenerate bool   `koanf:"auto-generate"`
}

type Duoc struct {
	SSOURL       string `koanf:"sso-url"`
	MobileAPIUrl string `koanf:"mobile-api-url"`
	WebAPIUrl    string `koanf:"web-api-url"`
	ClientSecret string `koanf:"client-secret"`
	ClientId     string `koanf:"client-id"`
	GrantType    string `koanf:"grant-type"`
}

func Default() *Config {
	return &Config{
		General: General{
			Debug: false,
		},
		HTTP: HTTP{
			Address: "0.0.0.0",
			Port:    80,
			SSL:     false,
		},
		JWT: JWT{
			PrivateKey:   "./data/private.pem",
			PublicKey:    "./data/public.pem",
			AutoGenerate: true,
		},
	}
}

func Apply(cfg *Config, opts ...Option) *Config {
	for _, op := range opts {
		cfg = op(cfg)
	}

	return cfg
}

func Build(opts ...Option) *Config {
	return Apply(Default(), opts...)
}
