package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config เก็บค่าการตั้งค่าทั้งหมดของแอปพลิเคชัน
type Config struct {
	ServerPort               string `mapstructure:"PORT"`
	Environment              string `mapstructure:"ENV"`
	MongoURI                 string `mapstructure:"MONGO_URI"`
	MongoDBName              string `mapstructure:"MONGO_DB_NAME"`
	RedisAddr                string `mapstructure:"REDIS_ADDR"`
	RedisPassword            string `mapstructure:"REDIS_PASSWORD"`
	JwtAccessSecret          string `mapstructure:"JWT_ACCESS_SECRET"`
	JwtAccessExpirationMins  int    `mapstructure:"JWT_ACCESS_EXPIRATION_MINS"`
	JwtRefreshSecret         string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtRefreshExpirationMins int    `mapstructure:"JWT_REFRESH_EXPIRATION_MINS"`
	AiEstateRagUri           string `mapstructure:"AI_ESTATE_RAG_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath("/secrets")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_ACCESS_EXPIRATION_MINS", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRATION_MINS", 10080)
	viper.SetDefault("AI_ESTATE_RAG_URI", "localhost:50052")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Printf("⚠️ Warning: .env file not found, fallback to system environment variables: %v", err)
	}

	bindEnvKeys()

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("❌ Unable to decode configuration into struct: %v", err)
		return
	}

	return config, nil
}

func bindEnvKeys() {
	keys := []string{
		"PORT",
		"ENV",
		"MONGO_URI",
		"MONGO_DB_NAME",
		"REDIS_ADDR",
		"REDIS_PASSWORD",
		"JWT_ACCESS_SECRET",
		"JWT_ACCESS_EXPIRATION_MINS",
		"JWT_REFRESH_SECRET",
		"JWT_REFRESH_EXPIRATION_MINS",
		"AI_ESTATE_RAG_URI",
	}
	for _, key := range keys {
		_ = viper.BindEnv(key)
	}
}
