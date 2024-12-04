package usecase

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/repository/postgres"
	"context"
	"fmt"
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

const (
	clientRole = `CLIENT`
	adminRole  = `ADMIN`
)

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository, cfg *config.ConfigModel) (*Usecase, error) {
	return &Usecase{
		cfg:  cfg,
		log:  logger,
		Repo: Repo,
	}, nil
}

func (uc *Usecase) GetUserToken(ctx context.Context, user *entities.User, password string) (string, error) {
	if user.Phone == "admin" && password == "admin" {
		user.Role = adminRole
	} else {
		user.Role = clientRole
	}

	err := uc.createUserToken(user)
	if err != nil {
		uc.log.Error("failed to create user token", zap.Error(err))
		return "", err
	}

	return user.Token, nil
}

func (uc *Usecase) GetTime() time.Time {
	return time.Now()
}

func (uc *Usecase) Admin() string {
	return fmt.Sprintf("By admin: time = %s", time.Now().UTC())
}

func (uc *Usecase) createUserToken(user *entities.User) error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Mail
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 10000).Unix()

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
