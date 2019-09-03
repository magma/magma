/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors

import (
	"errors"
	"fmt"
)

// Client APIs should raise ErrNotFound to indicate that a resource requested
// in a get/read does not exist.
var ErrNotFound = errors.New("Not found")
var ErrAlreadyExists = errors.New("Already exists")

func NewInitError(err error, service string) error {
	return ClientInitError{Err: err, Service: service}
}

// ClientInitError is a custom Go error type to indicate that initializing a
// client connection to a service failed.
type ClientInitError struct {
	Err     error
	Service string
}

func (c ClientInitError) Error() string {
	return fmt.Sprintf("%s service client initialization error: %s", c.Service, c.Err)
}
