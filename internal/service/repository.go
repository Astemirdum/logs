package service

import (
	"context"
	"github.com/Astemirdum/logs/models"
)

type Repository interface {
	WriteLog(ctx context.Context, raw string) (int64, error)
	ReadLog(ctx context.Context, id int) (models.Log, error)
	ReadLogs(ctx context.Context) ([]models.Log, error)
}
