package v1

import (
	"context"
	"errors"
	"github.com/bubalync/uni-auth/internal/lib/jwtgen"
	"github.com/bubalync/uni-auth/internal/mocks/servicemocks"
	authv1 "github.com/bubalync/uni-auth/internal/proto/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

func startGRPCServer(t *testing.T, as *servicemocks.MockAuth) (*bufconn.Listener, *grpc.Server) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	NewAuthServer(s, as)

	go func() {
		err := s.Serve(lis)
		require.NoError(t, err)

	}()

	return lis, s
}

func newGRPCClient(t *testing.T, lis *bufconn.Listener) (*grpc.ClientConn, authv1.AuthServiceClient) {
	cc, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
	)
	require.NoError(t, err)

	client := authv1.NewAuthServiceClient(cc)
	return cc, client
}

func TestAuthGRPCRoutes_ValidateToken(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *authv1.ValidateTokenRequest
	}

	type MockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name          string
		args          args
		inputBody     string
		mockBehaviour MockBehaviour
		wantResponse  *authv1.ValidateTokenResponse
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				ctx:     context.Background(),
				request: &authv1.ValidateTokenRequest{AccessToken: "valid-token"},
			},
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				claims := &jwtgen.Claims{
					UserId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Email:  "test@example.com",
				}
				m.EXPECT().ParseToken(args.request.AccessToken).Return(claims, nil)
			},
			wantResponse: &authv1.ValidateTokenResponse{
				IsValid: true,
				UserId:  "00000000-0000-0000-0000-000000000001",
				Email:   "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "input token is empty",
			args: args{
				ctx:     context.Background(),
				request: &authv1.ValidateTokenRequest{AccessToken: ""},
			},
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {},
			wantResponse:  nil,
			wantErr:       true,
		},
		{
			name: "auth service: token is invalid",
			args: args{
				ctx:     context.Background(),
				request: &authv1.ValidateTokenRequest{AccessToken: "valid-token"},
			},
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().ParseToken(args.request.AccessToken).Return(nil, errors.New("some error"))
			},
			wantResponse: &authv1.ValidateTokenResponse{
				IsValid: false,
				UserId:  "",
				Email:   "",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			as := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(as, tc.args)

			// create grpc server
			lis, server := startGRPCServer(t, as)
			defer server.Stop()

			// create grpc client
			cc, client := newGRPCClient(t, lis)
			defer cc.Close()

			// execute request
			resp, err := client.ValidateToken(tc.args.ctx, tc.args.request)

			// check response
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantResponse.IsValid, resp.IsValid)
			assert.Equal(t, tc.wantResponse.UserId, resp.UserId)
			assert.Equal(t, tc.wantResponse.Email, resp.Email)
		})
	}
}
