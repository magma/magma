// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"

	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/volatiletech/authboss"
)

type (
	// User models a user.
	User struct {
		*ent.User
		mutator userMutator
	}

	// userMutator builds / saves users.
	userMutator interface {
		SetPID(string)
		SetPassword(string)
		Save(context.Context) (*ent.User, error)
	}
)

// GetPID gets user primary id.
func (u *User) GetPID() string {
	return u.Email
}

// PutPID sets user primary id.
func (u *User) PutPID(pid string) {
	u.Email = pid
	u.mutator.SetPID(pid)
}

// GetPassword gets user password.
func (u *User) GetPassword() string {
	return u.Password
}

// PutPassword sets user password.
func (u *User) PutPassword(pwd string) {
	u.Password = pwd
	u.mutator.SetPassword(pwd)
}

// save persists user puts to ent client.
func (u *User) save(ctx context.Context) error {
	user, err := u.mutator.Save(ctx)
	if err == nil {
		*u = User{user, userUpdater{user.Update()}}
	}
	return err
}

// check if user implements necessary interfaces.
var (
	_ authboss.User         = (*User)(nil)
	_ authboss.AuthableUser = (*User)(nil)
)

// userCreator implements userMutator for creation.
type userCreator struct{ *ent.UserCreate }

func (u userCreator) SetPID(pid string)      { u.UserCreate.SetEmail(pid) }
func (u userCreator) SetPassword(pwd string) { u.UserCreate.SetPassword(pwd) }

// userUpdater implements userMutator for update.
type userUpdater struct{ *ent.UserUpdateOne }

func (u userUpdater) SetPID(pid string)      { u.UserUpdateOne.SetEmail(pid) }
func (u userUpdater) SetPassword(pwd string) { u.UserUpdateOne.SetPassword(pwd) }
