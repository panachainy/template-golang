package main

import (
	"template-golang/config"
	"template-golang/database"
	dbsqlc "template-golang/db/sqlc"
	"template-golang/modules/auth"
	authHandler "template-golang/modules/auth/handlers"
	authMiddleware "template-golang/modules/auth/middlewares"
	authRepo "template-golang/modules/auth/repositories"
	authUsecase "template-golang/modules/auth/usecases"
	"template-golang/modules/cockroach"
	cockroachHandler "template-golang/modules/cockroach/handlers"
	cockroachRepo "template-golang/modules/cockroach/repositories"
	cockroachUsecase "template-golang/modules/cockroach/usecases"
	"template-golang/server"
)

func main() {
	cfg := config.NewConfig(config.NewConfigOption(".env"))

	// Setup database
	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		panic(err)
	}
	pool := db.GetPool()
	queries := dbsqlc.New(pool)

	// Auth module wiring
	authRepository := authRepo.NewAuthRepository(queries)
	jwtUsecase := authUsecase.NewJWTUsecase(cfg, authRepository)
	middleware := authMiddleware.NewAuthMiddleware(jwtUsecase)
	handler := authHandler.NewAuthHttpHandler(jwtUsecase, cfg, middleware, authRepository)
	authModule := &auth.Auth{
		Handler:    handler,
		Middleware: middleware,
	}

	// Cockroach module wiring
	cockroachRepository := cockroachRepo.NewPostgresRepository(queries)
	cockroachMessaging := cockroachRepo.NewFCMMessaging()
	cockroachUsecase := cockroachUsecase.NewCockroachUsecaseImpl(cockroachRepository, cockroachMessaging)
	cockroachHandler := cockroachHandler.NewCockroachHttpHandler(cockroachUsecase)
	cockroachModule := &cockroach.Cockroach{
		Handler:    cockroachHandler,
		Repository: cockroachRepository,
		Messaging:  cockroachMessaging,
		Usecase:    cockroachUsecase,
	}

	// Create server
	s := server.NewGin(cfg, cockroachModule, authModule)
	s.Start()
}
