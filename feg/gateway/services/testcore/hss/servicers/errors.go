/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
)

// AuthDataUnavailableError indicates that an unexpectedly transient
// authentication failure occurs. See 3GPP TS 29.272 section 7.4.4.1.
type AuthDataUnavailableError struct {
	msg string
}

func (err AuthDataUnavailableError) Error() string {
	return fmt.Sprintf("Authentication data unavailable: %s", err.msg)
}

// NewAuthDataUnavailableError creates an AuthDataUnavailableError.
func NewAuthDataUnavailableError(msg string) AuthDataUnavailableError {
	return AuthDataUnavailableError{msg: msg}
}

// AuthRejectedError indicates that the HSS cannot return any authentication
// vectors due to unallowed attachment of the UE. See 3GPP TS 29.272 section 5.2.3.1.3.
type AuthRejectedError struct {
	msg string
}

func (err AuthRejectedError) Error() string {
	return fmt.Sprintf("Authentication rejected: %s", err.msg)
}

// NewAuthRejectedError creates an AuthRejectedError.
func NewAuthRejectedError(msg string) AuthRejectedError {
	return AuthRejectedError{msg: msg}
}
