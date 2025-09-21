package config

import (
	"os"
	"time"
)

type Config struct {
	Addr            string        // ex: ":8080"
	DatabaseURL     string        // ex: postgres://user:pass@localhost:5432/xeed?sslmode=disable
	ShutdownTimeout time.Duration // ex: 10s
	JWTSecret       string        // ← baru
	JWTTTL          time.Duration // ← baru
}

func FromEnv() Config {
	port := getenv("PORT", "8080")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// tetap boleh fallback, tapi kalau mau “wajib ada”, bisa panic/return error
		dsn = "postgres://postgres:postgres@127.0.0.1:5432/xeed?sslmode=disable"
	}

	ttl, _ := time.ParseDuration(getenv("JWT_TTL", "15m"))

	return Config{
		Addr:            ":" + port,
		DatabaseURL:     dsn,
		ShutdownTimeout: 10 * time.Second,
		JWTTTL:          ttl,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
