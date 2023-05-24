package config

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	DSN  string `envconfig:"DB_DSN"`
	Type string `envconfig:"DB_TYPE"`
}
