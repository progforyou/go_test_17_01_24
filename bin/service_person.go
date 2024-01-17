package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
	"testing/service/person/controller"
	"testing/service/person/data"
	"testing/service/person/web"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.DebugLevel)
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Debug().Msgf("Start application on port %s", port)

	dsn := os.Getenv("DNS")
	if dsn == "" {
		dsn = "host=localhost user=nikolai password=nikolai dbname=persons"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to open database")
	}
	err = db.AutoMigrate(&data.Person{})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to migrate database")
	}
	httpLogger := log.With().Str("service", "http").Logger().Level(zerolog.InfoLevel)
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(httpLogger))
	r.Route("/person", web.NewCrudRouter(controller.NewPersonController(db, log.Logger)))

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal().Err(err).Msg("fail start server")
	}

}
