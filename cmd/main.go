package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/p12s/csv-create-api/internal/repository"
	"github.com/p12s/csv-create-api/internal/service"
	"github.com/p12s/csv-create-api/internal/transport/rest"
	"github.com/sirupsen/logrus"
)

// @title Product app REST-API
// @version 0.0.1
// @description Simple product application for adding/getting products and download CSV-file
// @host localhost:8010
// @BasePath /
func main() {
	runtime.GOMAXPROCS(1)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s\n", err.Error())
	}

	db, err := repository.NewSqlite3DB(repository.Config{
		Driver: os.Getenv("DB_DRIVER"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s\n", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := rest.NewHandler(services)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handlers.InitRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("error while running http server: %s\n", err.Error())
		}
	}()
	logrus.Print("app started with port: ", os.Getenv("PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("app shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
}
