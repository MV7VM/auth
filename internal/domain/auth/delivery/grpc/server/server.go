package server

import (
	"auth/config"
	"auth/internal/domain/auth/entities"
	"auth/internal/domain/auth/usecase"
	protos "auth/pkg/proto/gen/go"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const statusOK = `OK`

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

// func convertUsersToProto(days []entities.User) []*protos.Day {
// 	Days := make([]*protos.Day, len(days))
// 	for i, _ := range days {
// 		Days[i] = &protos.Day{
// 			Date:     days[i].Date,
// 			Workouts: convertWorkoutsToProto(days[i].Workouts),
// 		}
// 	}
// 	return Days
// }

func (s *Server) GetAllByRole(ctx context.Context, request *protos.GetAllByRoleRequest) (*protos.GetAllByRoleResponse, error) {
	users, err := s.Usecase.GetAllByRole(ctx, request.GetRole())
	if err != nil {
		s.logger.Error("Fail to get users", zap.Error(err))
		return nil, err
	}
	return &protos.GetAllByRoleResponse{
		Id: users,
	}, nil
}

func (s *Server) CheckValidUserToken(ctx context.Context, request *protos.CheckValidUserTokenRequest) (*protos.CheckValidUserTokenResponse, error) {
 	user := convertToUserEntity("", "", nil, 0, "", request.GetToken())
	err := s.Usecase.CheckValidUserToken(ctx, user, request.GetToken())
	if err != nil {
		s.logger.Error("Token invalid", zap.Error(err))
		return nil, err
	}
	
	return &protos.CheckValidUserTokenResponse{
		Id: user.ID,
		Role: user.Role,
	}, nil
}

func (s *Server) GetUserToken(ctx context.Context, request *protos.GetUserTokenRequest) (*protos.GetUserTokenResponse, error) {
	token, err := s.Usecase.GetUserToken(
		ctx, 
		convertToUserEntity(
			"",
			request.GetLogin(),
			nil,
			0,
			"", 
			"",
		),
		request.GetPassword(),
	)
	if err != nil {
		return nil, err
	}
	return &protos.GetUserTokenResponse{
		Token:  token,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, request *protos.CreateUserRequest) (*protos.CreateUserResponse, error) {
	
	userID, err := s.Usecase.CreateUser(
		ctx, 
		convertToUserEntity(request.GetMail(), request.GetPhone(), nil, 0, request.GetRole(), ""), 
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
	if err := s.Usecase.UpdateUserPassword(ctx, convertToUserEntity("", "", nil, request.GetId(), "", ""), request.GetOldPassword(), request.GetNewPassword()); err != nil {
		return nil, err
	}
	return &protos.UpdateUserPasswordResponse{
		Status: statusOK,
	}, nil
}




func convertToUserEntity(mail, phone string, passwordHash []byte, ID uint64, role, token string) *entities.User {
	return &entities.User{
		ID:           ID,
		Mail:         mail,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role: 	      role,
		Token: 		  token,
	}
}
