package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func NewPostgreSQLConn(lifecycle fx.Lifecycle, log *zap.SugaredLogger) *pgxpool.Pool {

	pg_url := os.Getenv("POSTGRES_SQL_URL")
	config, err := pgxpool.ParseConfig(pg_url)

	log.Info("PostgreSQL URL", zap.String("url", pg_url))
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Start PostgreSQL connection")
			err := createTables(conn)
			if err != nil {
				log.Fatal("Error creating tables", zap.Error(err))
				return err
			}
			return nil
		},

		OnStop: func(context.Context) error {
			log.Info("Stopping database")
			conn.Close()
			return nil
		},
	})

	return conn
}

func createTables(conn *pgxpool.Pool) error {
	query := `
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name CHAR(64) NOT NULL,
    last_name CHAR(64) NOT NULL,
    email CHAR(64) NOT NULL UNIQUE, 
    password CHAR(64) NOT NULL
    );

CREATE TABLE IF NOT EXISTS todos(
    id SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users(id) ON DELETE CASCADE,
    description CHAR(255) NOT NULL,
    complete bool NOT NULL
);
`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}
