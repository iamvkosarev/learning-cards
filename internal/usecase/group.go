package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupUseCaseDeps struct {
	GroupReader  contracts.GroupReader
	GroupWriter  contracts.GroupWriter
	UserReader   contracts.UserReader
	UserWriter   contracts.UserWriter
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

func (g *GroupUseCase) CreateGroup(
	ctx context.Context,
	name, description string,
	visibility entity.GroupVisibility,
) (entity.GroupId, error) {
	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}
	_, err = g.UserReader.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			err = g.UserWriter.AddUser(
				ctx, entity.User{
					UserId: userId,
				},
			)
			if err != nil {
				return 0, fmt.Errorf("error creating user: %w", err)
			}
		} else {
			return 0, fmt.Errorf("error getting user: %w", err)
		}
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

	groupId, err := g.GroupWriter.AddGroup(ctx, group)
	if err != nil {
		return 0, err
	}

	return groupId, nil
}

func (g *GroupUseCase) GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	op := "usecase.GroupUseCase.GetGroup"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return entity.Group{}, err
	}

	group, err := getGroupAndCheckAccess(ctx, userId, groupId, op, g.GroupReader)
	if err != nil {
		return entity.Group{}, err
	}
	return group, nil
}

func getGroupAndCheckAccess(
	ctx context.Context, userId entity.UserId, groupId entity.GroupId, op string,
	r contracts.GroupReader,
) (entity.Group, error) {
	group, err := r.GetGroup(ctx, groupId)
	if err != nil {
		return entity.Group{}, err
	}

	if err = checkViewGroupAccess(userId, group, op); err != nil {
		return entity.Group{}, err
	}
	return group, nil
}

func (g *GroupUseCase) List(ctx context.Context) ([]entity.Group, error) {

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := g.GroupReader.ListGroups(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (g *GroupUseCase) UpdateGroup(ctx context.Context, updateGroup entity.UpdateGroup) error {
	op := "usecase.GroupUseCase.UpdateGroup"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	group, err := g.GroupReader.GetGroup(ctx, updateGroup.Id)

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

	err = g.GroupWriter.UpdateGroup(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func (g *GroupUseCase) DeleteGroup(ctx context.Context, groupId entity.GroupId) error {
	op := "usecase.GroupUseCase.DeleteGroup"

	userId, err := g.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	group, err := g.GroupReader.GetGroup(ctx, groupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	err = g.GroupWriter.DeleteGroup(ctx, groupId)
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
