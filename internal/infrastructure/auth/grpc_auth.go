package auth

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/status"

	sso_pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcService struct {
	client sso_pb.SSOClient
}

func NewGRPCService(client sso_pb.SSOClient) *grpcService {
	return &grpcService{
		client: client,
	}
}

func (s *grpcService) VerifyUserByContext(ctx context.Context) (entity.UserId, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, entity.ErrMetadataIsEmpty
	}
	ctxWithMD := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.client.VerifyToken(ctxWithMD, &emptypb.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			status.New(st.Code(), st.Message())
			return 0, entity.NewVerificationError(
				status.Error(
					st.Code(), fmt.Sprintf("failed to verify token: %v", err),
				),
			)
		}
		return 0, entity.NewVerificationError(fmt.Errorf("failed to verify token: %w", err))
	}
	return entity.UserId(res.GetUserId()), nil
}
