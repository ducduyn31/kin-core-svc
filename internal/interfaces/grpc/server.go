package grpc

import (
	"log/slog"
	"net"

	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/internal/application/circle"
	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/handlers"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Logger           *slog.Logger
	Auth0Validator   *auth.Auth0Validator
	UserService      *user.Service
	CircleService    *circle.Service
	EnableReflection bool
}

type Server struct {
	grpcServer *grpc.Server
	logger     *slog.Logger
}

func NewServer(cfg ServerConfig) *Server {
	recoveryInterceptor := interceptors.NewRecoveryInterceptor(cfg.Logger)
	authInterceptor := interceptors.NewAuthInterceptor(cfg.Auth0Validator, cfg.UserService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor.Unary(),
			authInterceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			recoveryInterceptor.Stream(),
			authInterceptor.Stream(),
		),
	)

	userHandler := handlers.NewUserHandler(cfg.UserService)
	circleHandler := handlers.NewCircleHandler(cfg.CircleService)

	kinv1.RegisterUserServiceServer(grpcServer, userHandler)
	kinv1.RegisterCircleServiceServer(grpcServer, circleHandler)

	if cfg.EnableReflection {
		reflection.Register(grpcServer)
	}

	return &Server{
		grpcServer: grpcServer,
		logger:     cfg.Logger,
	}
}

func (s *Server) Serve(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.logger.Info("gRPC server listening", "address", address)
	return s.grpcServer.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.logger.Info("gracefully stopping gRPC server...")
	s.grpcServer.GracefulStop()
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

func (s *Server) ServeWithListener(lis net.Listener) error {
	s.logger.Info("gRPC server listening", "address", lis.Addr().String())
	return s.grpcServer.Serve(lis)
}
