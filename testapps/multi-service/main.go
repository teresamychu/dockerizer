package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Postgres
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("postgres not available: %v", err)
	} else {
		defer conn.Close(context.Background())
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	defer rdb.Close()

	// RabbitMQ
	rmq, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Printf("rabbitmq not available: %v", err)
	} else {
		defer rmq.Close()
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	fmt.Printf("Starting server on :%s\n", port)
	r.Run(":" + port)
}
