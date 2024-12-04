package server

import (
	"auth/config"
	"auth/internal/domain/auth/delivery/http/middleware"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/usecase"
	protos "auth/pkg/proto/gen/go"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const statusOK = `OK`

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	serv    *gin.Engine
	Usecase *usecase.Usecase
	mdlware *middleware.Middleware
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase, mdlware *middleware.Middleware) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		serv:    gin.Default(),
		Usecase: uc,
		mdlware: mdlware,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	go func() {
		s.logger.Debug("serv started")
		if err := s.serv.Run(s.cfg.Server.Host + ":" + s.cfg.Server.Port); err != nil {
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

func (s *Server) Time(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"time": s.Usecase.GetTime().String()})
}

func (s *Server) Amin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": s.Usecase.Admin()})
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
