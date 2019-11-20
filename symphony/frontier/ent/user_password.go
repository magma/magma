// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

// HashPassword hashes password, returned value can be passed to user mutation.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return "", xerrors.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

// MustHashPassword calls HashPassword and panics on error.
func MustHashPassword(password string) string {
	hash, err := HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}

// ValidatePassword validates user password.
func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password), []byte(password),
	); err != nil {
		return xerrors.Errorf("validating user password: %w", err)
	}
	return nil
}
