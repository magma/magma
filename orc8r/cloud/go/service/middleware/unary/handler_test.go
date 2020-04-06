/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package unary

import (
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func simpleSuccessMockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return true, nil
}

func simpleErrorMockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return nil, errors.New("some error")
}

func simplePanicMockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	panic("failed")
}

func createFakeServerInfo() *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{
		FullMethod: "some method",
	}
}

func TestCallHandlerSimpleSuccess(t *testing.T) {
	resp, err := callHandler(
		context.Background(),
		nil,
		createFakeServerInfo(),
		simpleSuccessMockHandler,
	)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCallHandlerSimpleError(t *testing.T) {
	_, err := callHandler(
		context.Background(),
		nil,
		createFakeServerInfo(),
		simpleErrorMockHandler,
	)

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"some error",
	)
}

func TestCallHandlerPanics(t *testing.T) {
	_, err := callHandler(
		context.Background(),
		nil,
		createFakeServerInfo(),
		simplePanicMockHandler,
	)

	assert.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"rpc error: code = Unknown desc = Handler Panic: failed",
	)

	assert.EqualValues(
		t,
		1,
		testutil.ToFloat64(uncaughtCounterVec),
	)
}
