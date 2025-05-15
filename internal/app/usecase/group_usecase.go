package usecase

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupUseCaseDeps struct {
	GroupReader  contracts.GroupReader
	GroupWriter  contracts.GroupWriter
	AuthVerifier contracts.AuthVerifier
}

type GroupUseCase struct {
	GroupUseCaseDeps
}

func NewGroupUseCase(deps GroupUseCaseDeps) *GroupUseCase {
	return &GroupUseCase{
		GroupUseCaseDeps: deps,
	}
}

func (g *GroupUseCase) Create(
	ctx context.Context,
	name, description string,
	visibility entity.GroupVisibility,
) (entity.GroupId, error) {
	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}

	if visibility == entity.GROUP_VISIBILITY_NULL {
		visibility = entity.GROUP_VISIBILITY_PRIVATE
	}

	group := entity.Group{
		Name:        name,
		Description: description,
		Visibility:  visibility,
		OwnerId:     userId,
	}

	groupId, err := g.GroupWriter.Add(ctx, group)
	if err != nil {
		return 0, err
	}

	return groupId, nil
}

func (g *GroupUseCase) Get(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	op := "usecase.GroupUseCase.Get"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return entity.Group{}, err
	}

	group, err := g.GroupReader.Get(ctx, groupId)
	if err != nil {
		return entity.Group{}, err
	}

	if err := checkViewGroupAccess(userId, group, op); err != nil {
		return entity.Group{}, err
	}
	return group, nil
}

func (g *GroupUseCase) List(ctx context.Context) ([]entity.Group, error) {

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := g.GroupReader.ListByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (g *GroupUseCase) Update(ctx context.Context, updateGroup entity.UpdateGroup) error {
	op := "usecase.GroupUseCase.Update"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	group, err := g.GroupReader.Get(ctx, updateGroup.Id)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	if updateGroup.Visibility != entity.GROUP_VISIBILITY_NULL {
		group.Visibility = updateGroup.Visibility
	}
	if updateGroup.Description != "" {
		group.Description = updateGroup.Description
	}
	if updateGroup.Name != "" {
		group.Name = updateGroup.Name
	}

	err = g.GroupWriter.Update(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func (g *GroupUseCase) Delete(ctx context.Context, groupId entity.GroupId) error {
	op := "usecase.GroupUseCase.Delete"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	group, err := g.GroupReader.Get(ctx, groupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	err = g.GroupWriter.Delete(ctx, groupId)
	if err != nil {
		return err
	}

	return nil
}

func checkViewGroupAccess(userId entity.UserId, group entity.Group, op string) error {
	if userId != group.OwnerId &&
		(group.Visibility == entity.GROUP_VISIBILITY_PRIVATE ||
			group.Visibility == entity.GROUP_VISIBILITY_NULL) {
		message := fmt.Sprintf("%v: user (id:%v) not owner of card groups", op, userId)
		return entity.NewVerificationError(status.Error(codes.PermissionDenied, message))
	}
	return nil
}

func checkEditGroupAccess(userId entity.UserId, group entity.Group, op string) error {
	if userId != group.OwnerId {
		message := fmt.Sprintf("%v: user (id:%v) not owner of card groups", op, userId)
		return entity.NewVerificationError(status.Error(codes.PermissionDenied, message))
	}
	return nil
}
