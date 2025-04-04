package verification

import (
	"context"
	"fmt"
	sso_pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSSOVerifier struct {
	client sso_pb.SSOClient
}

func NewGRPCVerifier(hostAddress string) (*GRPCSSOVerifier, error) {
	ssoConn, err := grpc.NewClient(hostAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error creating gRPC client: %w", err)
	}
	ssoClient := sso_pb.NewSSOClient(ssoConn)
	return &GRPCSSOVerifier{client: ssoClient}, nil
}

func (g *GRPCSSOVerifier) VerifyToken(ctx context.Context, token string) (int64, error) {
	res, err := g.client.VerifyToken(ctx, &sso_pb.VerifyTokenRequest{Token: token})
	if err != nil {
		return 0, err
	}
	return res.GetUserId(), nil
}
