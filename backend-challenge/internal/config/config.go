package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                                int
	LogLevel                            string
	Enviornment                         string
	CouponCodeFolderPath                string
	CouponCodeFilePartialIndexChunkSize int
	CouponCodeFileConcurrentPoolSize    int
}

var AppConfig Config

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		Port:                                getEnvInt("PORT", 8080),
		LogLevel:                            strings.ToLower(getEnvString("LOG_LEVEL", "info")),
		Enviornment:                         strings.ToLower(getEnvString("ENVIRONMENT", "devlelopment")),
		CouponCodeFolderPath:                strings.ToLower(mustGetEnv("COUPON_CODE_FOLDER_PATH")),
		CouponCodeFilePartialIndexChunkSize: getEnvInt("COUPON_CODE_FILE_PARTIAL_INDEX_CHUNK_SIZE", 100000),
		CouponCodeFileConcurrentPoolSize:    getEnvInt("COUPON_CODE_FILE_CONCURRENT_POOL_SIZE", 5),
	}
}

func getEnvString(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Error in convertng the %s env variable to int: %v\n", key, err)
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		panic(fmt.Sprintf("required environment variable %s not set", key))
	}
}
