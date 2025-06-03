package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

//go:generate minimock -i UserVerifier -o ./mocks/user_verifier_mock.go -n UserVerifier -p mocks
type UserVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error)
}

//go:generate minimock -i UserReader -o ./mocks/user_reader_mock.go -n UserReader -p mocks
type UserReader interface {
	GetUser(ctx context.Context, id entity.UserId) (entity.User, error)
}

//go:generate minimock -i UserWriter -o ./mocks/user_reader_mock.go -n UserWriter -p mocks
type UserWriter interface {
	AddUser(ctx context.Context, user entity.User) error
}

//go:generate minimock -i GroupReader -o ./mocks/group_reader_mock.go -n GroupReaderMock -p mocks
type GroupReader interface {
	GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error)
	ListGroups(ctx context.Context, id entity.UserId) ([]entity.Group, error)
}

//go:generate minimock -i GroupWriter -o ./mocks/group_writer_mock.go -n GroupWriterMock -p mocks
type GroupWriter interface {
	AddGroup(ctx context.Context, group entity.Group) (entity.GroupId, error)
	UpdateGroup(ctx context.Context, group entity.Group) error
	DeleteGroup(ctx context.Context, groupId entity.GroupId) error
}

type GroupServiceDeps struct {
	GroupReader  GroupReader
	GroupWriter  GroupWriter
	UserReader   UserReader
	UserWriter   UserWriter
	UserVerifier UserVerifier
}

type GroupService struct {
	GroupServiceDeps
}

func NewGroupService(deps GroupServiceDeps) *GroupService {
	return &GroupService{
		GroupServiceDeps: deps,
	}
}

func (g *GroupService) CreateGroup(
	ctx context.Context,
	name, description string,
	visibility entity.GroupVisibility,
) (entity.GroupId, error) {
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
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

func (g *GroupService) GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	group, err := g.GroupReader.GetGroup(ctx, groupId)
	if err != nil {
		return entity.Group{}, err
	}
	if err = g.getReadGroupAccessByGroup(ctx, group); err != nil {
		return entity.Group{}, err
	}
	return group, nil
}

func (g *GroupService) List(ctx context.Context) ([]entity.Group, error) {
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := g.GroupReader.ListGroups(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (g *GroupService) UpdateGroup(ctx context.Context, updateGroup entity.UpdateGroup) error {
	group, err := g.GroupReader.GetGroup(ctx, updateGroup.Id)

	if err != nil {
		return err
	}

	if err = g.getWriteGroupAccessByGroup(ctx, group); err != nil {
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

func (g *GroupService) DeleteGroup(ctx context.Context, groupId entity.GroupId) error {
	if err := g.CheckWriteGroupAccess(ctx, groupId); err != nil {
		return err
	}
	if err := g.GroupWriter.DeleteGroup(ctx, groupId); err != nil {
		return err
	}
	return nil
}

func (g *GroupService) CheckReadGroupAccess(ctx context.Context, groupId entity.GroupId) error {
	group, err := g.GroupReader.GetGroup(ctx, groupId)

	if err != nil {
		return err
	}
	return g.getReadGroupAccessByGroup(ctx, group)
}

func (g *GroupService) getReadGroupAccessByGroup(ctx context.Context, group entity.Group) error {
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	if userId != group.OwnerId &&
		(group.Visibility == entity.GROUP_VISIBILITY_PRIVATE ||
			group.Visibility == entity.GROUP_VISIBILITY_NULL) {
		return entity.ErrGroupReadAccessDenied
	}
	return nil
}

func (g *GroupService) CheckWriteGroupAccess(ctx context.Context, groupId entity.GroupId) error {
	group, err := g.GroupReader.GetGroup(ctx, groupId)

	if err != nil {
		return err
	}
	return g.getWriteGroupAccessByGroup(ctx, group)
}

func (g *GroupService) getWriteGroupAccessByGroup(ctx context.Context, group entity.Group) error {
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	if userId != group.OwnerId {
		return entity.ErrGroupWriteAccessDenied
	}
	return nil
}
