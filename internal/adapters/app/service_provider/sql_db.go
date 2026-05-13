package service_provider

import (
	"SmartLeague/pkg/closer"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func (s *ServiceProvider) SQLDB() *sql.DB {
	if s.sqlDB == nil {
		s.Logger().Debugf("Connecting to SQL database (dsn=%s)", s.PGConfig().DSN())

		db, err := sql.Open("postgres", s.PGConfig().DSN())
		if err != nil {
			s.Logger().Panicf("failed to open sql database: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			s.Logger().Panicf("failed to connect to sql database: %v", err)
		}

		closer.Add(func() error {
			s.Logger().Info("Closing SQL database connection")
			return db.Close()
		})

		s.sqlDB = db
	}

	return s.sqlDB
}

