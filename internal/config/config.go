package config

import "github.com/imhinotori/duoc-plus/internal/envutil"

type Config struct {
	General General `koanf:"general"`
	HTTP    HTTP    `koanf:"http"`
	JWT     JWT     `koanf:"jwt"`
	Duoc    Duoc    `koanf:"duoc"`
}

type General struct {
	Debug     bool `koanf:"debug"`
	Cache     bool `koanf:"cache"`
	CacheTime int  `koanf:"cache-time"`
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

func LoadFromEnvironment() *Config {
	return &Config{
		General: General{
			Debug:     envutil.GetEnvBool("debug_mode"),
			Cache:     !envutil.GetEnvBool("cache_disabled"),
			CacheTime: envutil.GetEnvInt("cache_time", 60),
		},
		HTTP: HTTP{
			Address: envutil.GetEnv("address", "0.0.0.0"),
			Port:    envutil.GetEnvInt("port", 80),
			SSL:     envutil.GetEnvBool("ssl_enabled"),
		},
		JWT: JWT{
			PrivateKey:   envutil.GetEnv("jwt_private_key", "./private.pem"),
			PublicKey:    envutil.GetEnv("jwt_public_key", "./public.pem"),
			AutoGenerate: envutil.GetEnvBool("jwt_auto_generate"),
		},
		Duoc: Duoc{
			SSOURL:       envutil.GetEnv("duoc_sso_url", "duoc_sso_url"),
			MobileAPIUrl: envutil.GetEnv("duoc_mobile_api_url", "duoc_mobile_api_url"),
			WebAPIUrl:    envutil.GetEnv("duoc_web_api_url", "duoc_web_api_url"),
			ClientSecret: envutil.GetEnv("duoc_client_secret", "secret"),
			ClientId:     envutil.GetEnv("duoc_client_id", "client"),
			GrantType:    envutil.GetEnv("duoc_grant_type", "password"),
		},
	}
}
