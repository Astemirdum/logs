package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
	"time"

	"github.com/Astemirdum/logs/internal/config"
	"github.com/Astemirdum/logs/migrations"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type DB struct {
	master *sqlx.DB
	slaves []*sqlx.DB
}

func NewDB(cfg *config.DB) (*DB, error) {
	var (
		err error
		db  DB
	)
	db.master, err = sqlx.Open("pgx", newDSN(cfg, 0))
	if err != nil {
		return nil, err
	}

	db.slaves = make([]*sqlx.DB, len(cfg.Hosts)-1)
	for i := range db.slaves {
		db.slaves[i], err = sqlx.Open("pgx", newDSN(cfg, i+1))
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}

func (db *DB) Close() error {
	err := db.master.Close()
	for i := range db.slaves {
		err = multierr.Append(err, db.slaves[i].Close())
	}
	return err
}

func NewPostgresDB(cfg *config.DB) (*DB, error) {
	for i := range cfg.Hosts {
		if err := MigrateSchema(cfg, i); err != nil {
			return nil, err
		}
	}
	return NewDB(cfg)
}

func newDSN(cfg *config.DB, instance int) string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Hosts[instance], cfg.Ports[instance], cfg.Username, cfg.NameDB, cfg.Password)
}

const (
	pubName = "log_pub"
	subName = "log_sub_%d"
)

func MigrateSchema(cfg *config.DB, instance int) error {
	dsn := newDSN(cfg, instance)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("migrateSchema ping: %w", err)
	}

	goose.SetBaseFS(migrations.MigrationFiles)

	if err = goose.Up(db, "."); err != nil {
		return errors.Wrap(err, "goose run()")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if instance == 0 {
		pubQuery := fmt.Sprintf(`CREATE PUBLICATION %s FOR TABLE %s
			WITH (publish = 'insert,update');`, pubName, logTable)
		if _, err = db.ExecContext(ctx, pubQuery); err != nil {
			logrus.Warnf("pub creation %s", err.Error())
		}
	} else {
		subQuery := fmt.Sprintf(`CREATE SUBSCRIPTION %s
			CONNECTION '%s'
		PUBLICATION %s;`,
			fmt.Sprintf(subName, instance),
			newDSN(cfg, 0),
			pubName)
		if _, err = db.ExecContext(ctx, subQuery); err != nil {
			logrus.Warnf("pub creation %s", err.Error())
		}
	}
	return nil
}
