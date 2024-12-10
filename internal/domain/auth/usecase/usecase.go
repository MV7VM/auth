package usecase

import (
	"auth/config"
	"auth/internal/domain/auth/common"
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
	ok, err := uc.isAdmin(ctx, user, password)
	if err != nil {
		uc.log.Error("failed to check admin", zap.Error(err))
		return "", err
	}

	if ok {
		user.Role = adminRole
	} else {
		err := uc.client(ctx, user, password)
		if err != nil {
			return "", err
		}
		user.Role = clientRole
	}

	err = uc.createUserToken(user)
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

func (uc *Usecase) ChangePassword(ctx context.Context, user *entities.User) error {
	err := uc.Repo.UpdateUserPassword(ctx, user)
	if err != nil {
		uc.log.Error("failed to change user password", zap.Error(err))
		return err
	}

	return nil
}

func (uc *Usecase) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	users, err := uc.Repo.GetAllUsers(ctx)
	if err != nil {
		uc.log.Error("failed to get all users", zap.Error(err))
		return nil, err
	}

	return users, err
}

func (uc *Usecase) CreateUser() {

}

func (uc *Usecase) UpdateUser() {

}

func (uc *Usecase) isAdmin(ctx context.Context, user *entities.User, password string) (bool, error) {
	if user.Phone == adminRole {
		exist, err := uc.Repo.IsUserExist(ctx, user)
		if err != nil {
			uc.log.Error("failed to get user exist", zap.Error(err))
			return false, err
		}

		if exist {
			err = uc.Repo.GetUser(ctx, user)
			if err != nil {
				return false, err
			}
			if user.Role == adminRole {
				if password == string(user.PasswordHash) {
					return true, nil
				} else {
					return true, common.UserPasswordError
				}
			}
			return false, nil
		} else {
			user.Role = adminRole
			user.PasswordHash = []byte(password)
			userID, err := uc.Repo.CreateUser(ctx, user)
			if err != nil {
				return false, err
			}

			user.ID = int(userID)
			return true, nil
		}
	}

	return false, nil
}

func (uc *Usecase) client(ctx context.Context, user *entities.User, password string) error {
	exist, err := uc.Repo.IsUserExist(ctx, user)
	if err != nil {
		uc.log.Error("failed to get user exist", zap.Error(err))
		return err
	}

	if exist {
		err := uc.Repo.GetUser(ctx, user)
		if err != nil {
			uc.log.Error("failed to get user", zap.Error(err))
			return err
		}

		if string(user.PasswordHash) == password {
			return nil
		} else {
			return common.UserPasswordError
		}
	} else {
		//user.Role = clientRole
		//user.PasswordHash = []byte(password)
		//user.ID, err = uc.Repo.CreateUser(ctx, user)
		//if err != nil {
		//uc.log.Error("failed to create user", zap.Error(err))
		return common.UncreatedUserError
		//}
	}

	return nil
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
