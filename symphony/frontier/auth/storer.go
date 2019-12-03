// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/symphony/cloud/log"
	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/frontier/ent/user"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
)

// UserStorer loads and stores users.
type UserStorer struct {
	client *ent.UserClient
	logger log.Logger
}

// NewUserStorer creates ent based user storer.
func NewUserStorer(client *ent.UserClient, logger log.Logger) *UserStorer {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &UserStorer{
		client: client,
		logger: logger,
	}
}

// Load user by email for attached tenant.
func (s *UserStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	logger := s.logger.For(ctx).
		With(zap.String("user", key))
	switch u, err := s.client.Query().
		Where(
			user.Email(key),
			user.Tenant(
				CurrentTenant(ctx).Name,
			),
		).Only(ctx); err.(type) {
	case nil:
		logger.Debug("loaded user")
		return &User{u, userUpdater{u.Update()}}, nil
	case *ent.ErrNotFound:
		logger.Debug("user not found")
		return nil, authboss.ErrUserNotFound
	default:
		logger.Error("cannot load user", zap.Error(err))
		return nil, fmt.Errorf("cannot load user: %w", err)
	}
}

// Save persists user updates.
func (s *UserStorer) Save(ctx context.Context, user authboss.User) error {
	logger := s.logger.For(ctx).
		With(zap.String("user", user.GetPID()))
	err := user.(*User).save(ctx)
	if err == nil {
		logger.Debug("saved user")
		return nil
	}
	logger.Error("cannot save user", zap.Error(err))
	var e *ent.ErrConstraintFailed
	if errors.As(err, &e) {
		return authboss.ErrUserFound
	}
	return fmt.Errorf("cannot save user: %w", err)
}

// New creates a blank user.
func (s *UserStorer) New(ctx context.Context) authboss.User {
	creator := s.client.Create().
		SetTenant(CurrentTenant(ctx).Name).
		SetNetworks([]string{})
	return &User{&ent.User{}, userCreator{creator}}
}

// Create persists blank user to database.
func (s *UserStorer) Create(ctx context.Context, user authboss.User) error {
	return s.Save(ctx, user)
}

// check if user implements necessary interface.
var _ authboss.CreatingServerStorer = (*UserStorer)(nil)
