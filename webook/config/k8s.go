//go:build k8s

package config

var Config = config{
	DB: DBconfig{
		DSN: "root:root@tcp(webook-record-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-record-redis:6380",
	},
}
