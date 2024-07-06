package server

import (
	"auth/config"
	"auth/internal/domain/auth/usecase"
	protos "auth/pkg/proto/gen/go"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	RPC     *grpc.Server
	Usecase *usecase.Usecase
	protos.UnimplementedAuthServer
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		RPC:     grpc.NewServer(),
		Usecase: uc,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Server.Host+":"+s.cfg.Server.Port)
	if err != nil {
		s.logger.Error("failed to listen: ", zap.Error(err))
		return fmt.Errorf("failed to listen:  %w", err)
	}
	protos.RegisterAuthServer(s.RPC, s)
	reflection.Register(s.RPC) //по сети теперь видно все методы сети
	go func() {
		s.logger.Debug("grps serv started")
		if err = s.RPC.Serve(lis); err != nil {
			s.logger.Error("failed to serve: " + err.Error())
		}
		return
	}()
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.logger.Debug("stop grps")
	s.RPC.GracefulStop()
	return nil
}

func GetUserToken(context.Context, *protos.GetUserTokenRequest) (*protos.GetUserTokenResponse, error) {

}

func CreateUser(context.Context, *protos.CreateUserRequest) (*protos.CreateUserResponse, error) {

}

func UpdateUserPassword(context.Context, *protos.UpdateUserPasswordRequest) (*protos.UpdateUserPasswordResponse, error) {

}
