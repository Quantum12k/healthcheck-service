package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type (
	Config struct {
		DBName   string `yaml:"db_name"`
		Host     string `yaml:"db_host"`
		Port     string `yaml:"db_port"`
		User     string `yaml:"db_user"`
		Password string `yaml:"db_password"`
	}

	PostgreSQL struct {
		db *sqlx.DB
	}
)

func New(cfg Config) (*PostgreSQL, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open sql db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &PostgreSQL{db: sqlx.NewDb(db, "postgres")}, nil
}
