package service

import (
	"context"
	"github.com/Astemirdum/logs/models"

	"github.com/Astemirdum/logs/internal/repository"
)

type Service struct {
	repo Repository
}

func NewService(repo *repository.LogRepository) *Service {
	return &Service{repo: repo}
}

func (svc *Service) CreateLog(ctx context.Context, raw string) (int64, error) {
	return svc.repo.WriteLog(ctx, raw)
}

func (svc *Service) GetLog(ctx context.Context, id int) (models.Log, error) {
	return svc.repo.ReadLog(ctx, id)
}

func (svc *Service) ListLogs(ctx context.Context) ([]models.Log, error) {
	return svc.repo.ReadLogs(ctx)
}
