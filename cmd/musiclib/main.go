package main

import (
	"go.uber.org/zap"
	"log"
	"musiclib/config"
	"musiclib/internal/app/httphandlers"
	"musiclib/internal/app/services/extraDataAPIProvider"
	"musiclib/pkg/databases/gormpostgres"
	"net/http"
)

func main() {
	conf := config.Configure()

	//set logger
	logConf := zap.NewProductionConfig()
	logConf.DisableStacktrace = true
	logLevel, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level, err: %v", err)
	}
	logConf.Level = logLevel
	logger, err := logConf.Build()
	if err != nil {
		log.Fatalf("Failed to build logger: %v", err)
	}
	sugar := logger.Sugar()
	sugar.Infof("Sugar logger was created with log level `%v`", logLevel)

	//set storage
	storage, err := gormpostgres.NewGormDB(conf.DBConnectionString)
	if err != nil {
		sugar.Fatalf("Failed to connect to database, err: %v", err)
	}
	err = storage.Migrate()
	if err != nil {
		sugar.Fatalf("Failed to migrate database, err: %v", err)
	}

	//set extra data provider
	provider := extraDataAPIProvider.NewExtraDataAPIProvider(conf.ExtraDataAPIAddress)

	//build and run server:
	r := httphandlers.NewHTTPRouter(sugar, storage, provider)
	sugar.Infof("Starting an HTTP server on address `%v`...", conf.ServerAddress)
	http.ListenAndServe(conf.ServerAddress, r)
}
