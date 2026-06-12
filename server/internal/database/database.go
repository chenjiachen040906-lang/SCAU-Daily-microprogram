package database

import (
	"context"
	"log"

	"scau-daily/internal/config"
	"scau-daily/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func Init() {
	initPostgres()
	initRedis()
}

func initPostgres() {
	var err error

	logLevel := logger.Silent
	if config.AppConfig.Port == "8080" {
		logLevel = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	err = DB.AutoMigrate(
		&model.User{},
		&model.Semester{},
		&model.Course{},
		&model.CourseSchedule{},
		&model.Todo{},
		&model.Reminder{},
		&model.Notification{},
		&model.ChatSession{},
		&model.ChatMessage{},
		&model.UserSettings{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("[DB] PostgreSQL connected and migrated")
}

func initRedis() {
	opts, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	Redis = redis.NewClient(opts)

	ctx := context.Background()
	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("[DB] Redis connected")
}

func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
	if Redis != nil {
		Redis.Close()
	}
}
