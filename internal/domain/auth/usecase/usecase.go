package usecase

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/repository/postgres"
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	cfg  *config.ConfigModel
	log  *zap.Logger
	Repo *postgres.Repository
}

const clientRole = `CLIENT`

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository, cfg *config.ConfigModel) (*Usecase, error) {
	return &Usecase{
		cfg:  cfg,
		log:  logger,
		Repo: Repo,
	}, nil
}

func (uc *Usecase) GetUserToken(ctx context.Context, user *entities.User, password string) (string, error) {
	if exist, err := uc.Repo.IsUserExist(ctx, user); err != nil || !exist {
		return "", errors.New("user does not exist")
	}

	err := uc.Repo.GetUser(ctx, user)
	if err != nil {
		uc.log.Error("fail to GetUser", zap.Error(err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		uc.log.Error("Invalid Password", zap.Error(err))
		return "", err
	}

	if user.Token != "" { //если токен есть в бд то
		return user.Token, nil
	}

	if err := uc.createUserToken(user); err != nil {
		uc.log.Error("Invalid Password", zap.Error(err))
		return "", errors.New("fail to generate user token:" + err.Error())
	}

	if err := uc.Repo.UpdateUserToken(ctx, user); err != nil {
		return "", err
	}

	return user.Token, nil
}

func (uc *Usecase) CreateUser(ctx context.Context, user *entities.User, password string) (uint64, error) {
	if exist, err := uc.Repo.IsUserExist(ctx, user); exist || err != nil {
		uc.log.Error("fail to create user: user exist", zap.Error(err))
		return 0, errors.New("user exist")
	}

	if exist, err := uc.Repo.IsRoleExist(ctx, user.Role); !exist || err != nil {
		uc.log.Error("fail to create user: role does not exist", zap.Error(err))
		return 0, errors.New("role does not exist")
	}
	passwordHash, err := uc.encryptPassword(password)
	if err != nil {
		uc.log.Error("fail to create passwordHash", zap.Error(err))
		return 0, err
	}
	user.PasswordHash = passwordHash

	userID, err := uc.Repo.CreateUser(ctx, user)
	if err != nil {
		uc.log.Error("fail to insert user into DB", zap.Error(err))
		return 0, err
	}
	return userID, nil
}

func (uc *Usecase) UpdateUserPassword(ctx context.Context, user *entities.User, oldPassword, newPassword string) error {
	if exist, err := uc.Repo.IsUserExistByID(ctx, user); err != nil || !exist {
		return errors.New("user does not exist")
	}

	
	err := uc.Repo.GetUserByID(ctx, user)
	if err != nil {
		uc.log.Error("fail to GetUser", zap.Error(err))
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(oldPassword)); err != nil {
		uc.log.Error("Invalid Password", zap.Error(err))
		return err
	}

	newPasswordHash, err := uc.encryptPassword(newPassword)
	if err != nil {
		uc.log.Error("fail to create passwordHash", zap.Error(err))
		return err
	}
	user.PasswordHash = newPasswordHash
	if err := uc.Repo.UpdateUserPassword(ctx, user); err != nil {
		uc.log.Error("fail to update passwordHash", zap.Error(err))
		return err
	}

	return nil
}

func (uc *Usecase) createUserToken(user *entities.User) error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Mail
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour*10000).Unix()

	tokenString, err := token.SignedString([]byte(uc.cfg.Secret))
	if err != nil {
		return err
	}

	user.Token = tokenString
	return nil
}

// func (uc *Usecase) parseUserToken(user *entities.User) error {
// 	claims := jwt.MapClaims{}
// 	token, err := jwt.ParseWithClaims(user.Token, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(uc.cfg.Secret), nil
// 	})
// 	if err != nil {
// 		uc.log.Error("Fail to parse Token: ", zap.Error(err))
// 		return err
// 	}
// 	fmt.Println(token)
// 	// do something with decoded claims
// 	for key, val := range claims {
// 		fmt.Printf("Key: %v, value: %v\n", key, val)
// 	}
// 	return nil
// }

func (uc *Usecase) encryptPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Error("Fail to encrypt password: ", zap.Error(err))
		return nil, err
	}
	return passwordHash, nil
}