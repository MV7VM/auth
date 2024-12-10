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

const qGetUserRoles = `
select role, grade from roles`

func (r *Repository) GetUserRoles(ctx context.Context) map[string]int {
	rows, err := r.DB.Query(ctx, qGetUserRoles)
	if err != nil {
		return nil
	}

	Roles := make(map[string]int)

	for rows.Next() {
		var (
			grade int
			role  string
		)

		err := rows.Scan(&role, &grade)
		if err != nil {
			return nil
		}

		Roles[role] = grade
	}

	return Roles
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

const queryCreateUser = `insert into users (login, password, role_id) 
					values ($1, $2, (select id from roles where role = $3))
					returning id`

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) (uint64, error) {
	var id uint64
	err := r.DB.QueryRow(ctx, queryCreateUser, user.Phone, user.PasswordHash, user.Role).Scan(&id)
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
	_, err := r.DB.Exec(ctx, queryUpdateUserPassword, user.PasswordHash, user.ID)
	if err != nil {
		r.log.Error("fail to update user password", zap.Error(err))
		return err
	}
	return nil
}

const qGetAllUsers = `
SELECT
    users.id,
    users.login as phone,
    roles.role
FROM
    users
right join
    roles ON users.role_id = roles.id
`

func (r *Repository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	rows, err := r.DB.Query(ctx, qGetAllUsers)
	if err != nil {
		return nil, err
	}

	users := make([]entities.User, 0, 8)

	for rows.Next() {
		var user entities.User
		err = rows.Scan(&user.ID, &user.Phone, &user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, err
}

const queryGetUser = `
SELECT
    users.id,
    users.login,
    users.password,
    roles.role
FROM
    users
right join
    roles ON users.role_id = roles.id
WHERE
	users.login = $1;
`

func (r *Repository) GetUser(ctx context.Context, user *entities.User) error {
	var userID int
	var login, role string
	var passwordHash []byte

	err := r.DB.QueryRow(ctx, queryGetUser, user.Phone).Scan(&userID, &login, &passwordHash, &role)
	if err != nil {
		r.log.Error("fail to select data from users: ", zap.Error(err))
		return err
	}

	user.ID = userID
	user.Phone = login
	user.PasswordHash = passwordHash
	user.Role = role
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
	user.ID = userID
	//user.Mail = mail
	user.Phone = login
	user.PasswordHash = passwordHash
	user.Role = role
	//user.Token = token
	return nil
}
