package usecase

import (
	"auth/internal/domain/auth/repository/postgres"
	"errors"
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

func (uc *Usecase) GetUserToken(login, password string) (string, string, error) {
	if exist, err := uc.Repo.IsUserExist(login, password); err != nil || !exist {
		return "", "", errors.New("user does not exist")
	}

	userRole, err := uc.Repo.GetUserRole(login, password)
	if err != nil {
		uc.log.Error("fail too get User Token", zap.Error(err))
		return "", "", err
	}

	userToken, err := uc.Repo.GetUserToken(login, password)
	if err != nil {
		uc.log.Error("fail to get user token", zap.Error(err))
		return "", "", err
	}

	if userToken != "" { //если токен есть в бд то
		return userToken, userRole, nil
	}

	userToken, err = uc.createUserToken()
	if err != nil {
		return "", "", errors.New("fail to generate user token:" + err.Error())
	}

	return userToken, userRole, nil
}

func (uc *Usecase) createUserToken() (string, error) {
	return "", nil
}
