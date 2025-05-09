// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/google/wire"
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/cockroach"
	"template-golang/modules/cockroach/handlers"
	"template-golang/modules/cockroach/repositories"
	"template-golang/modules/cockroach/usecases"
)

// Injectors from wire.go:

func Wire() (Server, error) {
	configConfig := config.Provide()
	postgresDatabase := database.Provide(configConfig)
	cockroachPostgresRepository := repositories.ProvidePostgresRepository(postgresDatabase)
	cockroachFCMMessaging := repositories.ProvideFCMMessaging()
	cockroachUsecaseImpl := usecases.Provide(cockroachPostgresRepository, cockroachFCMMessaging)
	cockroachHttpHandler := handlers.Provide(cockroachUsecaseImpl)
	cockroachCockroach := &cockroach.Cockroach{
		Handler:    cockroachHttpHandler,
		Repository: cockroachPostgresRepository,
		Messaging:  cockroachFCMMessaging,
		Usecase:    cockroachUsecaseImpl,
	}
	serverGinServer := Provide(configConfig, cockroachCockroach)
	return serverGinServer, nil
}

// wire.go:

var ProviderSet = wire.NewSet(
	Provide, wire.Bind(new(Server), new(*ginServer)), config.ProviderSet, database.ProviderSet, cockroach.ProviderSet,
)
