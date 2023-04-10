package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlGorm() *gorm.DB {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dbconfig := dbUser + ":" + dbPass + "@tcp(" + dbHost + ")/" + dbName + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dbconfig), &gorm.Config{})
	if err != nil {
		log.Printf("if db error :", err)
		return nil
	}

	return db
}

func RedisConnection() (*redis.Client, error) {
	address := os.Getenv("REDIS_DSN")
	password := os.Getenv("REDIS_PASSWORD")
	user := os.Getenv("REDIS_USER")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         address, //redis port
		Password:     password,
		Username:     user,
		DB:           db,
		PoolSize:     1000,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})
	return client, nil
}
