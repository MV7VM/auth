package grpc

import (
	"auth/config"
	protos "auth/pkg/proto/gen/go"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	logger *zap.Logger
	cfg    *config.ConfigModel
	RPC    *grpc.Server
	protos.UnimplementedAuthServer
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel) (*Server, error) {
	return &Server{
		logger: logger,
		cfg:    cfg,
		RPC:    grpc.NewServer(),
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
