package usecase

import (
	"auth/internal/domain/auth/repository/postgres"
	"go.uber.org/zap"
)

type Usecase struct {
	log  *zap.Logger
	Repo *postgres.Repository
}

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository) (*Usecase, error) {
	return &Usecase{
		log:  logger,
		Repo: Repo,
	}, nil
}
