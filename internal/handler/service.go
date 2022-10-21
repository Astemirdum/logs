package handler

import (
	"context"
	"github.com/Astemirdum/logs/models"
)

//go:generate go run github.com/golang/mock/mockgen -source=service.go -destination=mocks/mock.go

type Service interface {
	CreateLog(ctx context.Context, raw string) (int64, error)
	GetLog(ctx context.Context, id int) (models.Log, error)
	ListLogs(ctx context.Context) ([]models.Log, error)
}
