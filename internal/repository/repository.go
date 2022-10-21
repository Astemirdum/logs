package repository

import (
	"context"
	"database/sql"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	qb = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type LogRepository struct {
	db  *DB
	log *zap.Logger
}

func NewRepository(db *DB, log *zap.Logger) *LogRepository {
	return &LogRepository{
		db:  db,
		log: log.Named("repo"),
	}
}

type TxRunner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type TxFunc func(tx *sqlx.Tx) error

func RunTx(ctx context.Context, db TxRunner, fn TxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				err = multierr.Combine(err, rollErr)
			}
		}
	}()
	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
