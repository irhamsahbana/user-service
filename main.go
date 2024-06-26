package main

import (
	"database/sql"
	"user-service/src/util/config"
	"user-service/src/util/routes"

	"github.com/thedevsaddam/renderer"

	userUsecase "user-service/src/app/dto/users"
	userHandler "user-service/src/handlers/users"
	userStore "user-service/src/util/repository/users"

	integrationUseCase "user-service/src/app/dto/users/integrations"
	integrationHandler "user-service/src/handlers/users/integrations"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	sqlDb, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer sqlDb.Close()

	render := renderer.New()
	routes := setupRoutes(render, sqlDb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(render *renderer.Render, myDb *sql.DB) *routes.Routes {
	userStore := userStore.NewStore(myDb)
	userUsecase := userUsecase.NewUserUsecase(userStore)
	userHandler := userHandler.NewUserHandler(userUsecase, render)

	integrationUseCase := integrationUseCase.NewUserUsecase(userStore)
	integrationHandler := integrationHandler.NewHandler(render, userUsecase, integrationUseCase)

	return &routes.Routes{
		Integration: integrationHandler,
		User:        userHandler,
	}
}
