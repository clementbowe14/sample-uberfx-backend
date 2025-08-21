package main

import (
	"fmt"
	"github.com/clementbowe14/sample-uberfx-backend/db"
	"github.com/clementbowe14/sample-uberfx-backend/handlers"
	"github.com/clementbowe14/sample-uberfx-backend/handlers/user"
	"github.com/clementbowe14/sample-uberfx-backend/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Errorf("failed to initialize zap logger: %v", err)
	}

	sugar := logger.Sugar()

	//setup containers
	fx.New(
		fx.Supply(sugar),
		fx.Provide(
			fx.Annotate(
				server.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			db.NewPostgreSQLConn,
			server.NewHttpServer,
			db.NewUserDb,
		),
		fx.Provide(
			server.AsRoute(handlers.NewHelloWorldHandler),
			server.AsRoute(user.NewCreateUserHandler),
			server.AsRoute(user.NewGetUserHandler),
			server.AsRoute(user.NewUpdateUserHandler),
			server.AsRoute(user.NewDeleteUserHandler),
		),

		fx.Invoke(func(srv *http.Server) {}),
		fx.Invoke(func(conn *pgxpool.Pool) {}),
	).Run()

}
