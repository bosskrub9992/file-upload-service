package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bosskrub9992/file-upload-service/configs"
	"github.com/bosskrub9992/file-upload-service/emails"
	"github.com/bosskrub9992/file-upload-service/handlers"
	"github.com/bosskrub9992/file-upload-service/objectstores"
	"github.com/bosskrub9992/file-upload-service/services"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	cfg := configs.New()
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("key.json"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer storageClient.Close()
	gcs := objectstores.NewGoogleCloudStorage(storageClient)
	dialer := gomail.NewDialer(
		cfg.Secret.Email.Host,
		cfg.Secret.Email.Port,
		cfg.Secret.Email.Username,
		cfg.Secret.Email.Password,
	)
	emailSender := emails.NewSender(cfg, dialer)
	service := services.New(cfg, gcs, emailSender)
	handler := handlers.New(service)

	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/files", handler.Upload,
		middleware.BodyLimit(cfg.UploadFile.MaxFileSizeMBLimit),
	)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
