package usecase

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
)

type GroupUseCaseDeps struct {
	contracts.GroupReader
	contracts.GroupWriter
	contracts.AuthVerifier
}

type GroupUseCase struct {
	deps GroupUseCaseDeps
}

func NewGroupUseCase(deps GroupUseCaseDeps) *GroupUseCase {
	return &GroupUseCase{
		deps: deps,
	}
}

func (uc *GroupUseCase) Create(
	ctx context.Context,
	name, description string,
	visibility entity.GroupVisibility,
) (entity.GroupId, error) {
	userId, err := uc.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}

	group := entity.Group{
		Name:        name,
		Description: description,
		Visibility:  visibility,
		OwnerId:     userId,
	}

	groupId, err := uc.deps.GroupWriter.Add(ctx, group)
	if err != nil {
		return 0, err
	}

	return groupId, nil
}

func (uc *GroupUseCase) Get(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	op := "usecase.GroupUseCase.Get"

	userId, err := uc.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return entity.Group{}, err
	}
	group, err := uc.deps.GroupReader.Get(ctx, groupId)
	if err != nil {
		return entity.Group{}, err
	}

	if userId != group.OwnerId && group.Visibility == entity.GROUP_VISIBILITY_PRIVATE {
		message := fmt.Sprintf("%v: user (id:%v) not owner of card groups", op, userId)
		return entity.Group{}, entity.NewVerificationError(message, codes.PermissionDenied)
	}
	return group, nil
}

func (uc *GroupUseCase) List(ctx context.Context) ([]entity.Group, error) {
	userId, err := uc.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := uc.deps.GroupReader.ListByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
