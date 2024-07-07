package postgres

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	ctx context.Context
	log *zap.Logger
	cfg *config.ConfigModel
	DB  *pgxpool.Pool
}

func NewRepository(log *zap.Logger, cfg *config.ConfigModel, ctx context.Context) (*Repository, error) {
	return &Repository{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}, nil
}

func (r *Repository) OnStart(_ context.Context) error {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		r.cfg.Postgres.Host,
		r.cfg.Postgres.Port,
		r.cfg.Postgres.User,
		r.cfg.Postgres.Password,
		r.cfg.Postgres.DBName,
		r.cfg.Postgres.SSLMode)
	pool, err := pgxpool.Connect(r.ctx, connectionUrl)
	if err != nil {
		return err
	}
	r.DB = pool
	return nil
}

func (r *Repository) OnStop(_ context.Context) error {
	r.DB.Close()
	return nil
}

const queryGetUserID = `
	SELECT EXISTS (SELECT id
               FROM users
               WHERE (login = $1 or $1 = '') AND (password = $2 or $2 = '') and (id = $3 or $3 = 0))
`

func (r *Repository) IsUserExist(ctx context.Context, user *entities.User) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetUserID, user.Phone, user.Password, user.ID).Scan(&res)
	if err != nil {
		r.log.Error("fail to check user exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetUserRole = `select role 
from roles 
where id = (select roleid 
            from users 
            where login = $1 and password = $2)`

func (r *Repository) GetUserRole(ctx context.Context, login, pass string) (string, error) {
	var role string
	err := r.DB.QueryRow(ctx, queryGetUserRole, login, pass).Scan(&role)
	if err != nil {
		r.log.Error("fail to get user role from DB", zap.Error(err))
		return "", err
	}
	return role, nil
}

const queryGetUserToken = `SELECT COALESCE(token, '') AS token
					FROM users
					WHERE login = $1 and password = $2`

func (r *Repository) GetUserToken(ctx context.Context, login, pass string) (string, error) {
	var token string
	err := r.DB.QueryRow(ctx, queryGetUserToken, login, pass).Scan(&token)
	if err != nil {
		r.log.Error("fail to get user")
		return "", err
	}
	return token, nil
}

const queryCreateUser = `insert into users (login, password, roleid, mail) 
					values ($1, $2, (select id from roles where role = $3), $4)
					returning id`

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (uint64, error) {
	var id uint64
	err := r.DB.QueryRow(ctx, queryCreateUser, user.Phone, user.Password, user.Role, user.Mail).Scan(&id)
	if err != nil {
		r.log.Error("fail to create client", zap.Error(err))
		return 0, err
	}
	return id, nil
}

const queryUpdateUserPassword = `update users 
									set password = $1
									where id = $2`

func (r *Repository) UpdateUserPassword(ctx context.Context, user *entities.User) error {
	_, err := r.DB.Exec(ctx, queryUpdateUserPassword, user.Password, user.ID)
	if err != nil {
		r.log.Error("fail to update user password", zap.Error(err))
		return err
	}
	return nil
}

const queryUpdateUserToken = `update users
								set token = $1
								where login = $2 and password = $3`

func (r *Repository) UpdateUserToken(ctx context.Context, token string, user *entities.User) error {
	_, err := r.DB.Exec(ctx, queryUpdateUserToken, token, user.Phone, user.Password)
	if err != nil {
		r.log.Error("fail to exec user token", zap.Error(err))
		return err
	}
	return nil
}
