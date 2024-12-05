package config

type AppConfig struct {
	Auth struct {
		JWTSecretKey string `env:"JWT_SECRET_KEY"` // double check with consts.JWTSecretKey. It should be the same
	} `yaml:"auth"`
}
