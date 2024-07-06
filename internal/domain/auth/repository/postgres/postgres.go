package postgres

import (
	"auth/config"
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

const queryGetUser = `
	SELECT EXISTS (SELECT id
               FROM users
               WHERE login = $1 AND password = $2);
`

func (r *Repository) IsUserExist(log, pass string) (bool, error) {
	return false, nil
}

const queryGetUserRole = ``

func (r *Repository) GetUserRole(log, pas string) (string, error) {
	return "", nil
}

const queryGetUserToken = ``

func (r *Repository) GetUserToken(log, pas string) (string, error) {
	return "", nil
}
