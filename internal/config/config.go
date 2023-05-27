package config

type Config struct {
	DB      DBConfig
	Storage StorageConfig
}

type DBConfig struct {
	DSN  string `envconfig:"DB_DSN" default:"test.db"`
	Type string `envconfig:"DB_TYPE" default:"sqlite"`
}

type StorageConfig struct {
	URL    string `envconfig:"STORAGE_URL" default:"http://hackathon2023.obs.ru-moscow-1.hc.sbercloud.ru"`
	Origin string `envconfig:"STORAGE_ORIGIN" default:"sbercloud.ru"`
}
