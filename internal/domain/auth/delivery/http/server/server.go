package server

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/usecase"
	protos "auth/pkg/proto/gen/go"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
)

const statusOK = `OK`

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	serv    *gin.Engine
	Usecase *usecase.Usecase
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		serv:    gin.Default(),
		Usecase: uc,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Server.Host+":"+s.cfg.Server.Port)
	if err != nil {
		s.logger.Error("failed to listen: ", zap.Error(err))
		return fmt.Errorf("failed to listen:  %w", err)
	}

	go func() {
		s.logger.Debug("serv started")
		if err = s.serv.RunListener(lis); err != nil {
			s.logger.Error("failed to serve: " + err.Error())
		}
		return
	}()
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.logger.Debug("stop grps")
	//s.serv.GracefulStop()
	return nil
}

func (s *Server) GetUserToken(ctx *gin.Context) {
	request := protos.GetUserTokenRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to unmarshar request: %v", err)})
		return
	}

	token, err := s.Usecase.GetUserToken(
		ctx,
		convertToUserEntity(
			"",
			request.Login,
			nil,
			0,
			"",
		),
		request.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get token: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, &protos.GetUserTokenResponse{
		Token: token,
	})
	return
}

func (s *Server) CreateUser(ctx context.Context, request *protos.CreateUserRequest) (*protos.CreateUserResponse, error) {

	userID, err := s.Usecase.CreateUser(
		ctx,
		convertToUserEntity(request.GetMail(), request.GetPhone(), nil, 0, request.GetRole()),
		request.GetPassword(),
	)
	if err != nil {
		s.logger.Error("fail to create user", zap.Error(err))
		return nil, err
	}
	return &protos.CreateUserResponse{
		UserId: userID,
	}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, request *protos.UpdateUserPasswordRequest) (*protos.UpdateUserPasswordResponse, error) {
	if err := s.Usecase.UpdateUserPassword(ctx, convertToUserEntity("", "", nil, request.GetId(), ""), request.GetOldPassword(), request.GetNewPassword()); err != nil {
		return nil, err
	}
	return &protos.UpdateUserPasswordResponse{
		Status: statusOK,
	}, nil
}

func convertToUserEntity(mail, phone string, passwordHash []byte, ID uint64, role string) *entities.User {
	return &entities.User{
		ID:           ID,
		Mail:         mail,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role:         role,
	}
}
