package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppPort string
	AppEnv  string
	DB      *sql.DB
	Logger  *logrus.Logger
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("no .env file found, reading from environment")
	}

	logger := buildLogger(getEnv("APP_ENV", "development"))
	db := connectDB(logger)

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
		DB:      db,
		Logger:  logger,
	}
}

func connectDB(log *logrus.Logger) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASS", "postgres"),
		getEnv("DB_NAME", "storyku_db"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.WithError(err).Fatal("failed to open database connection")
	}

	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("failed to ping database")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	log.Info("database connected successfully")
	return db
}

func buildLogger(env string) *logrus.Logger {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		time.Local = loc
	}

	log := logrus.New()

	if env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		log.SetLevel(logrus.DebugLevel)
	}

	hook := NewDailyErrorHook("logs")
	log.AddHook(hook)

	return log
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}