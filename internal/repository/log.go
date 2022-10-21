package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errr "github.com/Astemirdum/logs/internal/errs"
	"github.com/Astemirdum/logs/models"
	"go.uber.org/zap"
	"math/rand"
)

const (
	logTable = "log.logs"
)

var (
	logsSelect = []string{"id", "raw", "created_at"}
)

func (r *LogRepository) WriteLog(ctx context.Context, raw string) (int64, error) {
	b := qb.Insert(logTable).
		Columns("raw").
		Values(raw).Suffix("RETURNING \"id\"")

	query, args, err := b.ToSql()
	if err != nil {
		return 0, err
	}
	r.log.Debug("WriteLog",
		zap.String("query", query),
		zap.Any("args", args),
	)
	var id int64
	if err = r.db.master.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LogRepository) ReadLog(ctx context.Context, id int) (models.Log, error) {
	slaveKey := rand.Intn(len(r.db.slaves))
	b := qb.Select(logsSelect...).From(logTable).Where("id=$1", id)
	query, args, err := b.ToSql()
	if err != nil {
		return models.Log{}, err
	}
	r.log.Debug("ReadLog",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Int("slaveKey", slaveKey),
	)
	var log models.Log
	if err = r.db.slaves[slaveKey].GetContext(ctx, &log, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("ReadLog slave:%d : %s: %w", slaveKey, err.Error(), errr.ErrNotFound)
		}
		return models.Log{}, err
	}
	return log, nil
}

func (r *LogRepository) ReadLogs(ctx context.Context) ([]models.Log, error) {
	slaveKey := rand.Intn(len(r.db.slaves))
	tx, err := r.db.slaves[slaveKey].BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				r.log.Warn("rollback", zap.Error(rollErr))
			}
		}
	}()

	b := qb.Select("COUNT(*)").From(logTable)
	query, args, err := b.ToSql()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("readLogs COUNT: %s: %w", err.Error(), errr.ErrNotFound)
		}
		return nil, err
	}
	r.log.Debug("count logs",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Int("slaveKey", slaveKey),
	)
	var count int64
	if err = r.db.slaves[slaveKey].QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return nil, err
	}

	b = qb.Select(logsSelect...).From(logTable)
	query, args, err = b.ToSql()
	if err != nil {
		return nil, err
	}
	r.log.Debug("ReadLogs",
		zap.String("query", query),
		zap.Any("args", args),
	)
	logs := make([]models.Log, 0, count)
	if err = r.db.slaves[slaveKey].SelectContext(ctx, &logs, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("list logs: %s: %w", err.Error(), errr.ErrNotFound)
		}
		return nil, err
	}
	return logs, tx.Commit()

}
