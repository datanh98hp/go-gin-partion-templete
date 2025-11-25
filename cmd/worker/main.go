package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMq    rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	config      *config.Config
	logger      *zerolog.Logger
}

func NewWorker(cfg *config.Config) *Worker {
	logs := utils.NewLoggerWithPath("worker.log", "info")
	/// RabbitMQ service
	rabbitMq, err := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		logs,
	)
	if err != nil {
		logs.Fatal().Err(err).Msg("Failed to initialize RabbitMQ service in worker")
	}
	// Init mail service
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory in worker")
		return nil
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to initialize mail service in worker")
		return nil
	}
	return &Worker{
		rabbitMq:    rabbitMq,
		mailService: mailService,
		config:      cfg,
		logger:      logs,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	// consume
	const emailQueueName = "auth_mail_queue"

	handler := func(body []byte) error {
		//handle with message ( send mail)

		w.logger.Info().Msgf("Received message: %s", string(body))
		var email mail.Email
		err := json.Unmarshal([]byte(body), &email)
		if err != nil {
			w.logger.Error().Err(err).Msg("Failed to unmarshal email from queue in worker")
			return err
		}
		if err := w.mailService.SendMail(ctx, &email); err != nil {
			return utils.NewWrapError("Failed to send mail", utils.ErrorCodeInternal, err)
		}

		w.logger.Info().Msgf("Email sent successfully: %v", email.To)
		return nil
	}

	// consume message
	if err := w.rabbitMq.Consume(ctx, emailQueueName, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to consume message in worker")
		return err
	}
	w.logger.Info().Msgf("Worker started, listening for messages in queue '%s'", emailQueueName)

	// wait for context to be canceled
	<-ctx.Done()
	w.logger.Info().Msgf("Worker stopped, consume due to context cancel in queue '%s", emailQueueName)
	return ctx.Err()

}
func (w *Worker) Shutdown(ctx context.Context) error {
	w.logger.Info().Msg("Shutdowning worker ...")
	if err := w.rabbitMq.Close(); err != nil {
		w.logger.Error().Err(err).Msg("Failed to close RabbitMQ in worker")
		return err
	}
	w.logger.Info().Msg("RabbitMQ closed successfully ...")
	// handle multiple message queue
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			w.logger.Warn().Msg("Shutdown timeout exceeded in worker")
		}
		return ctx.Err()
	default:
		w.logger.Warn().Msg("Shutdown completed in worker...")
		return nil
	}
}
func main() {
	rootDir := utils.MushGetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/worker.log")
	logger.InitLogger(logger.LoggerConfig{
		Level:     "info",
		FileName:  logFile,
		MaxSize:   2,
		MaxBackUp: 5,
		MaxAge:    5,
		Compress:  true,
		IsDev:     utils.GetEnv("APP_ENV", "development"),
	})
	//Load .env file
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Printf("Error loading .env file")
		logger.Log.Warn().Msg("Error loading .env file")
	} else {
		log.Printf("Loaded .env file in worker")
		logger.Log.Info().Msg("Loaded .env file in worker")
	}
	//initialize the config
	cfg := config.NewConfig()
	worker := NewWorker(cfg)
	if worker == nil {
		logger.Log.Fatal().Msg("Failed to create worker")
	}

	// block until we receive our signal.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := worker.Start(ctx); err != nil && err != context.Canceled {
			logger.Log.Fatal().Err(err).Msg("Failed to start worker")
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("Received signal, shutting down worker...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := worker.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to shutdown worker")
	}
	wg.Wait()
	logger.Log.Info().Msg("Worker stopped successfully ")
}
