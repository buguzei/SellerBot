package postgres

import (
	"bot/internal/config"
	log2 "bot/internal/log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Postgres struct {
	DB     *sql.DB
	logger log2.Logger
}

func NewPostgres(cfg config.DBConf, l log2.Logger) Postgres {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	l = l.Named("postgres")

	return Postgres{
		DB:     db,
		logger: l,
	}
}
