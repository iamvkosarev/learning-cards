package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const groupsTraceName = "module.groups"

//go:generate minimock -i UserVerifier -o ./mocks/user_verifier_mock.go -n NewUserVerifierMock -p mocks
type UserVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID model.UserId, err error)
}

//go:generate minimock -i UserReader -o ./mocks/user_reader_mock.go -n UserReaderMock -p mocks
type UserReader interface {
	GetUser(ctx context.Context, id model.UserId) (model.User, error)
}

//go:generate minimock -i UserWriter -o ./mocks/user_reader_mock.go -n UserWriterMock -p mocks
type UserWriter interface {
	AddUser(ctx context.Context, user model.User) error
}

//go:generate minimock -i GroupReader -o ./mocks/group_reader_mock.go -n GroupReaderMock -p mocks
type GroupReader interface {
	GetGroup(ctx context.Context, groupId model.GroupId) (*model.Group, error)
	ListGroups(ctx context.Context, id model.UserId) ([]*model.Group, error)
}

//go:generate minimock -i GroupWriter -o ./mocks/group_writer_mock.go -n GroupWriterMock -p mocks
type GroupWriter interface {
	AddGroup(ctx context.Context, group *model.Group) (model.GroupId, error)
	UpdateGroup(ctx context.Context, group *model.Group) error
	DeleteGroup(ctx context.Context, groupId model.GroupId) error
}

type GroupsDeps struct {
	GroupReader  GroupReader
	GroupWriter  GroupWriter
	UserReader   UserReader
	UserWriter   UserWriter
	UserVerifier UserVerifier
}

type Groups struct {
	GroupsDeps
	tracer trace.Tracer
}

func NewGroups(deps GroupsDeps) *Groups {
	return &Groups{
		GroupsDeps: deps,
		tracer:     otel.Tracer(groupsTraceName),
	}
}

func (g *Groups) CreateGroup(
	ctx context.Context,
	name, description string,
	visibility model.GroupVisibility, cardSideTypes []model.CardSideType,
) (model.GroupId, error) {
	ctx, span := g.tracer.Start(ctx, "CreateGroup")
	defer span.End()

	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}
	_, err = g.UserReader.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			err = g.UserWriter.AddUser(
				ctx, model.User{
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

	if visibility == model.GROUP_VISIBILITY_NULL {
		visibility = model.GROUP_VISIBILITY_PRIVATE
	}

	group := &model.Group{
		Name:          name,
		Description:   description,
		Visibility:    visibility,
		OwnerId:       userId,
		CardSideTypes: cardSideTypes,
	}

	groupId, err := g.GroupWriter.AddGroup(ctx, group)
	if err != nil {
		return 0, err
	}

	return groupId, nil
}

func (g *Groups) GetGroup(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	ctx, span := g.tracer.Start(ctx, "GetGroup")
	defer span.End()
	group, err := g.GroupReader.GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if err = g.GetReadGroupAccessByGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (g *Groups) List(ctx context.Context) ([]*model.Group, error) {
	ctx, span := g.tracer.Start(ctx, "List")
	defer span.End()
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

func (g *Groups) UpdateGroup(ctx context.Context, updateGroup model.UpdateGroup) error {
	ctx, span := g.tracer.Start(ctx, "UpdateGroup")
	defer span.End()
	op := "module.Groups.UpdateGroup"
	group, err := g.GroupReader.GetGroup(ctx, updateGroup.Id)

	if err != nil {
		return err
	}

	if err = g.getWriteGroupAccessByGroup(ctx, group); err != nil {
		return err
	}

	if updateGroup.Visibility != model.GROUP_VISIBILITY_NULL {
		group.Visibility = updateGroup.Visibility
	}
	if updateGroup.Description != "" {
		group.Description = updateGroup.Description
	}
	if updateGroup.Name != "" {
		group.Name = updateGroup.Name
	}

	if len(updateGroup.CardSideType) > 0 {
		for sideI, updateSideType := range updateGroup.CardSideType {
			currentSideType := group.CardSideTypes[sideI]

			if currentSideType != updateSideType {
				if updateSideType != model.CARD_SIDE_TYPE_NULL &&
					currentSideType != model.CARD_SIDE_TYPE_NULL {
					return fmt.Errorf(
						"%s: current side type %v, update sie type %v: %w", op, currentSideType, updateSideType,
						model.ErrGroupModifyNotNullCardsSideType,
					)
				} else {
					group.CardSideTypes[sideI] = updateSideType
				}
			}
		}
	}

	err = g.GroupWriter.UpdateGroup(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func (g *Groups) DeleteGroup(ctx context.Context, groupId model.GroupId) error {
	ctx, span := g.tracer.Start(ctx, "DeleteGroup")
	defer span.End()
	if _, err := g.CheckWriteGroupAccess(ctx, groupId); err != nil {
		return err
	}
	if err := g.GroupWriter.DeleteGroup(ctx, groupId); err != nil {
		return err
	}
	return nil
}

func (g *Groups) CheckReadGroupAccess(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	ctx, span := g.tracer.Start(ctx, "CheckReadGroupAccess")
	defer span.End()
	group, err := g.GroupReader.GetGroup(ctx, groupId)

	if err != nil {
		return nil, err
	}
	return group, g.GetReadGroupAccessByGroup(ctx, group)
}

func (g *Groups) GetReadGroupAccessByGroup(ctx context.Context, group *model.Group) error {
	ctx, span := g.tracer.Start(ctx, "GetReadGroupAccessByGroup")
	defer span.End()
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	if userId != group.OwnerId &&
		(group.Visibility == model.GROUP_VISIBILITY_PRIVATE ||
			group.Visibility == model.GROUP_VISIBILITY_NULL) {
		return model.ErrGroupReadAccessDenied
	}
	return nil
}

func (g *Groups) CheckWriteGroupAccess(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	ctx, span := g.tracer.Start(ctx, "CheckWriteGroupAccess")
	defer span.End()
	group, err := g.GroupReader.GetGroup(ctx, groupId)

	if err != nil {
		return nil, err
	}
	return group, g.getWriteGroupAccessByGroup(ctx, group)
}

func (g *Groups) getWriteGroupAccessByGroup(ctx context.Context, group *model.Group) error {
	ctx, span := g.tracer.Start(ctx, "getWriteGroupAccessByGroup")
	defer span.End()
	userId, err := g.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	if userId != group.OwnerId {
		return model.ErrGroupWriteAccessDenied
	}
	return nil
}
