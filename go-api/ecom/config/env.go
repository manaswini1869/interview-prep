package config

import "os"

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() *Config {

	return &Config{
		PublicHost: getEnv("PUBLIC_HOST", "db"),
		Port:       getEnv("PORT", "5432"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "something"),
		DBName:     getEnv("DB_NAME", "ecom"),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
