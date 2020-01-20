// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volatiletech/authboss"
)

type testUserMutator struct {
	mock.Mock
}

func (t *testUserMutator) SetPID(pid string) {
	t.Called(pid)
}

func (t *testUserMutator) SetPassword(pwd string) {
	t.Called(pwd)
}

func (t *testUserMutator) Save(ctx context.Context) (*ent.User, error) {
	args := t.Called(ctx)
	user, _ := args.Get(0).(*ent.User)
	return user, args.Error(1)
}

func TestUserGettersSetters(t *testing.T) {
	pid := "tester@example.com"
	pwd := "testpassword"

	var m testUserMutator
	m.On("SetPID", pid).Once()
	m.On("SetPassword", pwd).Once()
	m.On("Save", mock.Anything).
		Return(nil, authboss.ErrUserFound).
		Once()
	defer m.AssertExpectations(t)

	user := User{&ent.User{}, &m}
	user.PutPID(pid)
	assert.Equal(t, pid, user.GetPID())
	user.PutPassword(pwd)
	assert.Equal(t, pwd, user.GetPassword())
	err := user.save(context.Background())
	assert.EqualError(t, err, authboss.ErrUserFound.Error())
}
