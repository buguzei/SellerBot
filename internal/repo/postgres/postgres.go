package postgres

import (
	"bot/internal/config"
	"bot/internal/log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB     *sql.DB
	logger log.Logger
}

// NewPostgres is a constructor for Postgres
func NewPostgres(cfg config.DBConf, l log.Logger) Postgres {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName)

	fmt.Println(connStr)

	l = l.Named("postgres")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal("NewPostgres: error to open postgres", log.Fields{
			"error": err,
		})
	}

	err = db.Ping()
	if err != nil {
		l.Fatal("NewPostgres: error ping", log.Fields{
			"error": err,
		})
	}

	return Postgres{
		DB:     db,
		logger: l,
	}
}
