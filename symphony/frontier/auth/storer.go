// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"

	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/token"
	"github.com/facebookincubator/symphony/frontier/ent/user"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
)

// UserStorer loads and stores users.
type UserStorer struct {
	client *ent.Client
	logger log.Logger
}

// NewUserStorer creates ent based user storer.
func NewUserStorer(client *ent.Client, logger log.Logger) *UserStorer {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &UserStorer{
		client: client,
		logger: logger,
	}
}

// userPredicate returns user predicate for attached context.
func userPredicate(ctx context.Context, email string) predicate.User {
	return user.And(
		user.Email(email),
		user.Tenant(
			CurrentTenant(ctx).Name,
		),
	)
}

// load user by key for attached tenant.
func (s *UserStorer) load(ctx context.Context, key string) (*ent.User, error) {
	return s.client.User.
		Query().
		Where(
			userPredicate(ctx, key),
		).
		Only(ctx)
}

// Load user by email for attached tenant.
func (s *UserStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	logger := s.logger.For(ctx).
		With(zap.String("user", key))
	switch u, err := s.load(ctx, key); err.(type) {
	case nil:
		logger.Debug("loaded user")
		return &User{u, userUpdater{u.Update()}}, nil
	case *ent.NotFoundError:
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
	var e *ent.ConstraintError
	if errors.As(err, &e) {
		return authboss.ErrUserFound
	}
	return fmt.Errorf("cannot save user: %w", err)
}

// New creates a blank user.
func (s *UserStorer) New(ctx context.Context) authboss.User {
	creator := s.client.User.Create().
		SetTenant(CurrentTenant(ctx).Name).
		SetNetworks([]string{})
	return &User{&ent.User{}, userCreator{creator}}
}

// Create persists blank user to database.
func (s *UserStorer) Create(ctx context.Context, user authboss.User) error {
	return s.Save(ctx, user)
}

// Close underlying ent client.
func (s *UserStorer) Close() error {
	return s.client.Close()
}

// hashToken hashes a token to a string.
func hashToken(token string) string {
	hash := fnv.New64()
	_, _ = io.WriteString(hash, token)
	return hex.EncodeToString(hash.Sum(nil))
}

// AddRememberToken adds remember token to user.
func (s *UserStorer) AddRememberToken(ctx context.Context, pid, value string) error {
	logger := s.logger.For(ctx).
		With(zap.String("user", pid))
	u, err := s.load(ctx, pid)
	if err != nil {
		logger.Error("cannot load user", zap.Error(err))
		return fmt.Errorf("cannot load user: %w", err)
	}
	logger = logger.With(zap.String("token", hashToken(value)))
	switch _, err := s.client.Token.
		Create().
		SetUser(u).
		SetValue(value).
		Save(ctx); err.(type) {
	case nil:
		logger.Debug("saved remember token")
		return nil
	case *ent.ConstraintError:
		logger.Warn("remember token already exists")
		return nil
	default:
		logger.Error("cannot save remember token", zap.Error(err))
		return fmt.Errorf("cannot save remember token: %w", err)
	}
}

// DelRememberTokens clears all remember tokens of user.
func (s *UserStorer) DelRememberTokens(ctx context.Context, pid string) error {
	logger := s.logger.For(ctx).
		With(zap.String("user", pid))
	cnt, err := s.client.Token.
		Delete().
		Where(
			token.HasUserWith(
				userPredicate(ctx, pid),
			),
		).
		Exec(ctx)
	if err != nil {
		logger.Error("cannot clear remember tokens", zap.Error(err))
		return fmt.Errorf("cannot clear remember tokens: %w", err)
	}
	logger.Debug("cleared remember tokens", zap.Int("count", cnt))
	return nil
}

// UseRememberToken clears a single remember token of user.
func (s *UserStorer) UseRememberToken(ctx context.Context, pid, value string) error {
	logger := s.logger.For(ctx).With(
		zap.String("user", pid),
		zap.String("token", hashToken(value)),
	)
	switch cnt, err := s.client.Token.
		Delete().
		Where(
			token.Value(value),
			token.HasUserWith(
				userPredicate(ctx, pid),
			),
		).
		Exec(ctx); {
	case err != nil:
		logger.Error("cannot clear remember token", zap.Error(err))
		return fmt.Errorf("cannot clear remember token: %w", err)
	case cnt == 0:
		logger.Warn("remember token not found")
		return authboss.ErrTokenNotFound
	default:
		logger.Debug("cleared remember token")
		return nil
	}
}

// check if user storer implements necessary interfaces.
var (
	_ authboss.ServerStorer            = (*UserStorer)(nil)
	_ authboss.CreatingServerStorer    = (*UserStorer)(nil)
	_ authboss.RememberingServerStorer = (*UserStorer)(nil)
)
