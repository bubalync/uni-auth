package v1

import (
	"context"
	authv1 "github.com/bubalync/uni-auth/internal/proto/v1"
	"github.com/bubalync/uni-auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	authv1.UnimplementedAuthServiceServer
	as service.Auth
}

func NewAuthServer(gRPCServer *grpc.Server, as service.Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverApi{as: as})
}

func (s *serverApi) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if req.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	claims, err := s.as.ParseToken(req.GetAccessToken())
	if err != nil {
		return &authv1.ValidateTokenResponse{IsValid: false}, nil
	}

	return &authv1.ValidateTokenResponse{
		IsValid: true,
		UserId:  claims.UserId.String(),
		Email:   claims.Email,
	}, nil
}
