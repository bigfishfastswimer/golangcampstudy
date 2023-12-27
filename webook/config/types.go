package config

type config struct {
	DB    DBconfig
	Redis RedisConfig
}

type DBconfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}
