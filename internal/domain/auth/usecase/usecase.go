package usecase

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/repository/postgres"
	"context"
	"errors"
	"fmt"
	// "strconv"
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

func (uc *Usecase) GetAllByRole(ctx context.Context, role string) ([]uint64, error) {
	if exist, err := uc.Repo.IsRoleExist(ctx, role); !exist || err != nil {
		uc.log.Error("role does not exist", zap.Error(err))
		return nil, errors.New("role does not exist")
	}

	users, err := uc.Repo.GetUsersByRole(ctx, role)
	if err != nil {
		uc.log.Error("fail to get roleID", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (uc *Usecase) CheckValidUserToken(ctx context.Context, user *entities.User, tokenString string) (error) {
	if err := uc.parseUserToken(user); err != nil {
		uc.log.Error("fail to parse token", zap.Error(err))
		return err
	}
	if exist, err := uc.Repo.IsUserExistByID(ctx, user); err != nil || !exist {
		return errors.New("user does not exist")
	}

	if tokenString != user.Token {
		return errors.New("tokens do not same")
	}

	return nil	
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

// type TokenUser struct {
//     UserID uint64   `json:"uid"`
//     Email string   `json:"email"`
//     Role string   `json:"role"`
//     Exp int64   `json:"exp"`
// }


// func (tu TokenUser) Valid() error {
//     return nil
// }
// func (tu TokenUser) GetAudience() (jwt.ClaimStrings, error) {
//     return nil, nil
// }

func (uc *Usecase) createUserToken(user *entities.User) error {
	token := jwt.New(jwt.SigningMethodHS256,
        // TokenUser{
		// 	UserID: user.ID,
		// 	Email: user.Mail,
		// 	Role: user.Role,
		// 	Exp: time.Now().Add(time.Hour*10000).Unix(),
        // },
	)
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

func (uc *Usecase) parseUserToken(user *entities.User) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(user.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.Secret), nil
	})
	if err != nil {
		uc.log.Error("Fail to parse Token: ", zap.Error(err))
		return err
	}
	// var token_claims TokenUser
    // token, err := jwt.ParseWithClaims(user.Token, &token_claims, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte(uc.cfg.Secret), nil
    // })
    // if err != nil || !token.Valid {
	// 	uc.log.Error("Fail to parse Token: ", zap.Error(err))
	// 	return err
    // }

    // uc.log.Info("User roles: %v", token_claims.Roles)


	// for key, val := range claims {
	// 	if key ==  "uid" {
	// 		user.ID = claims.GetSubject()
	// 	} else if key ==  "email"{
	// 		user.Mail = string(val)
	// 	} else if key ==  "role"{
	// 		user.Role = string(val)
	// 	} 
	// }
	// fmt.Println("ID: ", user.ID, "Mail: ", user.Mail, "Role: ", user.Role,)
	
	rawUID, ok := claims["uid"]
	if !ok {
		return errors.New("err")
	}
	uid, ok := rawUID.(float64)

	rawRole, ok := claims["role"]
	if !ok {
		return errors.New("err")
	}
	role, ok := rawRole.(string)
	
	rawMail, ok := claims["email"]
	if !ok {
		return errors.New("err")
	}
	mail, ok := rawMail.(string)
	
	user.Role = role
	user.ID = uint64(uid)
	user.Mail = mail

	// num, err := claims["uid"].(float64) 
	fmt.Println(user.ID, user.Role, user.Mail)
	return nil

}

func (uc *Usecase) encryptPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Error("Fail to encrypt password: ", zap.Error(err))
		return nil, err
	}
	return passwordHash, nil
}