// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/ent"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// UserService is a user service.
	UserService struct{ Client UserProvider }

	// UserProvider returns an ent client given a context and tenant
	UserProvider func(context.Context, string) (*ent.Client, error)
)

// NewUserService create a new user service.
func NewUserService(provider UserProvider) UserService {
	return UserService{provider}
}

// Create a user by authID and tenantID.
func (s UserService) Create(ctx context.Context, input *UserInput) (*User, error) {
	if input.Tenant == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant")
	}
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing auth id")
	}

	client, err := s.Client(ctx, input.Tenant)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}

	u, err := client.User.Create().SetAuthID(input.Id).Save(ctx)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}

	return &User{Id: int64(u.ID)}, nil
}

// Delete a user by authID and tenantID.
func (s UserService) Delete(ctx context.Context, input *UserInput) (*empty.Empty, error) {
	if input.Tenant == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant")
	}
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing auth id")
	}

	client, err := s.Client(ctx, input.Tenant)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}

	err = client.User.Update().
		Where(user.AuthID(input.Id)).
		SetStatus(user.StatusDeleted).
		Exec(ctx)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}

	return &empty.Empty{}, nil
}
