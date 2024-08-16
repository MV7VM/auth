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
               WHERE login = $1);
`

func (r *Repository) IsUserExist(ctx context.Context, user *entities.User) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetUserID, user.Phone).Scan(&res)
	if err != nil {
		r.log.Error("fail to check user exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryCheckUserByID = `
	SELECT EXISTS (SELECT id
               FROM users
               WHERE id = $1);
`

func (r *Repository) IsUserExistByID(ctx context.Context, user *entities.User) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryCheckUserByID, user.ID).Scan(&res)
	if err != nil {
		r.log.Error("fail to check user exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetRoleID = `
	SELECT EXISTS (SELECT id
               FROM roles
               WHERE (role = $1))
`

func (r *Repository) IsRoleExist(ctx context.Context, role string) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetRoleID, role).Scan(&res)
	if err != nil {
		r.log.Error("fail to check role exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetUserRole = `select role 
from roles 
where id = (select roleid 
            from users 
            where login = $1 and password_hash = $2)`

func (r *Repository) GetUserRole(ctx context.Context, login string, passwordHash []byte) (string, error) {
	var role string
	err := r.DB.QueryRow(ctx, queryGetUserRole, login, passwordHash).Scan(&role)
	if err != nil {
		r.log.Error("fail to get user role from DB", zap.Error(err))
		return "", err
	}
	return role, nil
}

const queryGetUserToken = `SELECT COALESCE(token, '') AS token
					FROM users
					WHERE login = $1 and password_hash = $2`

func (r *Repository) GetUserToken(ctx context.Context, login string, passHash []byte) (string, error) {
	var token string
	err := r.DB.QueryRow(ctx, queryGetUserToken, login, passHash).Scan(&token)
	if err != nil {
		r.log.Error("fail to get user")
		return "", err
	}
	return token, nil
}

const queryCreateUser = `insert into users (login, password_hash, roleid, mail) 
					values ($1, $2, (select id from roles where role = $3), $4)
					returning id`

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (uint64, error) {
	var id uint64
	err := r.DB.QueryRow(ctx, queryCreateUser, user.Phone, user.PasswordHash, user.Role, user.Mail).Scan(&id)
	if err != nil {
		r.log.Error("fail to create client", zap.Error(err))
		return 0, err
	}
	return id, nil
}

const queryUpdateUserPassword = `update users 
									set password_hash = $1
									where id = $2`

func (r *Repository) UpdateUserPassword(ctx context.Context, user *entities.User) error {
	_, err := r.DB.Exec(ctx, queryUpdateUserPassword, user.PasswordHash, user.ID)
	if err != nil {
		r.log.Error("fail to update user password", zap.Error(err))
		return err
	}
	return nil
}

const queryUpdateUserToken = `update users
								set token = $1
								where login = $2`

func (r *Repository) UpdateUserToken(ctx context.Context, user *entities.User) error {
	_, err := r.DB.Exec(ctx, queryUpdateUserToken, user.Token, user.Phone)
	if err != nil {
		r.log.Error("fail to exec user token", zap.Error(err))
		return err
	}
	return nil
}

const queryGetUser = `
SELECT
    users.id,
    users.login,
    users.password_hash,
    COALESCE(users.token, '') AS token,
    users.mail,
    roles.role
FROM
    users
INNER JOIN
    roles ON users.roleID = roles.id
WHERE
	users.login = $1;
`



func (r *Repository) GetUser(ctx context.Context, user *entities.User) error {
	var userID int
	var login, token, mail, role string
	var passwordHash []byte

	err := r.DB.QueryRow(ctx, queryGetUser, user.Phone).Scan(&userID, &login, &passwordHash, &token, &mail, &role)
	if err != nil {
		r.log.Error("fail to select data from users: ", zap.Error(err))
		return err
	}
	user.ID = uint64(userID)
	user.Mail = mail
	user.Phone = login
	user.PasswordHash = passwordHash 
	user.Role = role
	user.Token = token
	fmt.Println(user.Token)
	return nil
}

const queryGetUserByID = `
SELECT
    users.id,
    users.login,
    users.password_hash,
    COALESCE(users.token, '') AS token,
    users.mail,
    roles.role
FROM
    users
INNER JOIN
    roles ON users.roleID = roles.id
WHERE
	users.id = $1;
`

func (r *Repository) GetUserByID(ctx context.Context, user *entities.User) error {
	var userID int
	var login, token, mail, role string
	var passwordHash []byte

	err := r.DB.QueryRow(ctx, queryGetUserByID, user.ID).Scan(&userID, &login, &passwordHash, &token, &mail, &role)
	if err != nil {
		r.log.Error("fail to select data from users: ", zap.Error(err))
		return err
	}
	user.ID = uint64(userID)
	user.Mail = mail
	user.Phone = login
	user.PasswordHash = passwordHash 
	user.Role = role
	user.Token = token
	return nil
}

const queryGetUsersByRole = `
SELECT
    users.id
FROM
    users
INNER JOIN
    roles ON users.roleID = roles.id
WHERE
    roles.role = $1;
`
func (r *Repository) GetUsersByRole(ctx context.Context, role string) ([]uint64, error) {
	var users []uint64
	rows, err := r.DB.Query(ctx, queryGetUsersByRole, role)
	if err != nil {
		r.log.Error("fail to select users.id by roles.role: ", zap.Error(err))
		return nil, err	
	}
	for rows.Next() {
		var userID uint64
		err := rows.Scan(&userID)
		if err != nil {
			r.log.Error("fail to scan row", zap.Error(err))
			return nil, err
		}

		users = append(users, userID)
	}

	if rows.Err() != nil {
		r.log.Error("rows iteration error", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return users, nil

}