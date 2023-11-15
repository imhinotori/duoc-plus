package config

type Config struct {
	General General `koanf:"general"`
	HTTP    HTTP    `koanf:"http"`
	JWT     JWT     `koanf:"jwt"`
}

type General struct {
	Debug bool `koanf:"debug"`
}

type HTTP struct {
	Address string `koanf:"address"`
	Port    int    `koanf:"port"`
	SSL     bool   `koanf:"ssl"`
}

type JWT struct {
	PrivateKey string `koanf:"private-key"`
	PublicKey  string `koanf:"public-key"`
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
			PrivateKey: "./keys/private.pem",
			PublicKey:  "./keys/public.pem",
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
