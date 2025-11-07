package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/pgx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var DB sqlc.Querier
var DBpool *pgxpool.Pool

func InitializeDatabase() error {
	connStr := config.NewConfig().DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("Error parsing config : %v", err)
	}

	//log
	sqlLogger := utils.NewLoggerWithPath("sql.log", "infor")
	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxZerologTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Microsecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}
	// Set connection pool settings
	conf.MaxConns = 50
	conf.MinConns = 10
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.MaxConnLifetime = 30 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	// Establish the connection pool
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBpool, err = pgxpool.NewWithConfig(context, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %v", err)
	}
	// Create a new Queries instance using the established connection pool
	DB = sqlc.New(DBpool) // Assuming sqlc package is imported correctly

	// Ping the database to ensure the connection is established
	if err := DBpool.Ping(context); err != nil {
		return fmt.Errorf("Unable to ping database: %v", err)
	}

	log.Println("Connected")
	return nil
}
