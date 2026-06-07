package main

import (
	"context"
	"log"

	"github.com/banggibima/go-english-test-platform/config"
	"github.com/banggibima/go-english-test-platform/internal/jobs"
	"github.com/banggibima/go-english-test-platform/internal/results"
	"github.com/banggibima/go-english-test-platform/pkg/database"
	"github.com/banggibima/go-english-test-platform/pkg/logger"
	"github.com/banggibima/go-english-test-platform/pkg/queue"
)

const scoreAttemptQueue = "score-attempt"

func main() {
	cfg := config.Load()
	logg := logger.New()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	if err := rabbit.DeclareQueue(scoreAttemptQueue); err != nil {
		log.Fatal(err)
	}

	resultRepository := results.NewRepository(db)
	jobRepository := jobs.NewRepository(db)
	jobService := jobs.NewService(jobRepository, resultRepository)

	messages, err := rabbit.Consume(scoreAttemptQueue)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("worker started", "queue", scoreAttemptQueue)

	for message := range messages {
		if err := jobService.HandleScoreAttempt(context.Background(), message.Body); err != nil {
			logg.Error("failed to process score attempt job", "error", err)
			_ = message.Nack(false, true)
			continue
		}

		_ = message.Ack(false)
		logg.Info("score attempt job processed")
	}
}
