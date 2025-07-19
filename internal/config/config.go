package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	CfgMinio     MinioConfig
	CfgRedis     RedisConfig
	Env          string `yaml:"env" env-default:"local"`
	StoragePath  string `yaml:"storage_path" env-required:"true"`
	JWTSecretKey string `yaml:"jwt_secret_key" env-required:"true"`
	HTTPServer   `yaml:"http_server"`
}

type MinioConfig struct {
	MinioPort         string
	MinioEndpoint     string
	BucketName        string
	MinioRootUser     string
	MinioRootPassword string
	MinioUseSSL       bool
}

type RedisConfig struct {
	Address     string
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

var CfgMinio *MinioConfig
var CfgRedis *RedisConfig

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("cannot read file: %w", err)
	}

	CfgMinio = &MinioConfig{
		MinioPort:         getEnv("PORT", "8083"),
		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
		BucketName:        getEnv("MINIO_BUCKET_NAME", "defaultbucket"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		MinioUseSSL:       getEnvAsBool("MINIO_USE_SSL", false),
	}

	CfgRedis = &RedisConfig{
		Address: "localhost:6379",
	}

	Config := Config{
		CfgMinio: *CfgMinio,
		CfgRedis: *CfgRedis,
	}

	if err := cleanenv.ReadConfig(configPath, &Config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	os.Setenv("JWT_SECRET_KEY", Config.JWTSecretKey)

	return &Config
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
