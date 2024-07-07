package usecase

import (
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/repository/postgres"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Usecase struct {
	log  *zap.Logger
	Repo *postgres.Repository
}

const clientRole = `CLIENT`

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository) (*Usecase, error) {
	return &Usecase{
		log:  logger,
		Repo: Repo,
	}, nil
}

func (uc *Usecase) GetUserToken(ctx context.Context, user *entities.User) (string, string, error) {
	if exist, err := uc.Repo.IsUserExist(ctx, user); err != nil || !exist {
		return "", "", errors.New("user does not exist")
	}

	userRole, err := uc.Repo.GetUserRole(ctx, user.Phone, user.Password)
	if err != nil {
		uc.log.Error("fail too get User Token", zap.Error(err))
		return "", "", err
	}

	userToken, err := uc.Repo.GetUserToken(ctx, user.Phone, user.Password) // sqlNoRows
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

	err = uc.Repo.UpdateUserToken(ctx, userToken, user)
	if err != nil {
		return "", "", err
	}

	return userToken, userRole, nil
}

func (uc *Usecase) CreateUser(ctx context.Context, user *entities.User) (uint64, error) {
	pass, err := generatePassword()
	if err != nil {
		uc.log.Error("fail to generate password", zap.Error(err))
		return 0, err
	}

	pass, err = uc.Crypt(pass)
	if err != nil {
		uc.log.Error("fail to crypt password", zap.Error(err))
		return 0, err
	}

	user.Password = pass
	user.Role = clientRole

	userID, err := uc.Repo.CreateUser(ctx, user)
	if err != nil {
		uc.log.Error("fail to insert user into DB", zap.Error(err))
		return 0, err
	}
	return userID, nil
}

func (uc *Usecase) UpdateUserPassword(ctx context.Context, user *entities.User) error {
	if exist, err := uc.Repo.IsUserExist(ctx, &entities.User{ID: user.ID}); err != nil || !exist {
		return errors.New("user does not exist")
	}

	cryptedPass, err := uc.Crypt(user.Password)
	if err != nil {
		uc.log.Error("fail to crypt password", zap.Error(err))
		return errors.New("fail to crypt password" + err.Error())
	}

	user.Password = cryptedPass

	err = uc.Repo.UpdateUserPassword(ctx, user)
	if err != nil {
		uc.log.Error("fail to update user password into DB", zap.Error(err))
		return errors.New("fail to update user password into DB" + err.Error())
	}
	return nil
}

func generatePassword() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func (uc *Usecase) Crypt(pass string) (string, error) {
	return pass, nil
}

func (uc *Usecase) createUserToken() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}
