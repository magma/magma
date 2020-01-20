// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signtest

import (
	"context"

	"github.com/facebookincubator/symphony/store/sign"

	"github.com/stretchr/testify/mock"
)

// MockSigner defines a mock signer.
type MockSigner struct {
	mock.Mock
}

// Sign implements sign.Signer interface.
func (m *MockSigner) Sign(ctx context.Context, op sign.Operation, key, filename string) (string, error) {
	args := m.Called(ctx, op, key, filename)
	return args.String(0), args.Error(1)
}
