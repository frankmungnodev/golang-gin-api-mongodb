package envs

import (
	"os"
	"strconv"
)

type Envs struct{}

var (
	ENV string

	MONGO_DATABASE_URL string

	TURSO_DATABASE_URL string
	TURSO_AUTH_TOKEN   string

	JWT_SECRET         string
	JWT_EXPIRE_IN_DAYS int
)

func LoadEnvs() {
	ENV = os.Getenv("ENV")

	MONGO_DATABASE_URL = os.Getenv("MONGO_DATABASE_URL")

	TURSO_DATABASE_URL = os.Getenv("TURSO_DATABASE_URL")
	TURSO_AUTH_TOKEN = os.Getenv("TURSO_AUTH_TOKEN")

	expireInDaysEnv := os.Getenv("JWT_EXPIRE_IN_DAYS")
	expireInDays, err := strconv.Atoi(expireInDaysEnv)
	if err != nil {
		panic("Failed to get JWT_EXPIRE_IN_HOUR")
	}
	JWT_SECRET = os.Getenv("JWT_SECRET")
	JWT_EXPIRE_IN_DAYS = expireInDays
}
