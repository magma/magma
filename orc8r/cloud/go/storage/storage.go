/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package storage contains common definitions to be used across service
// storage interfaces
package storage

import (
	"fmt"

	"github.com/google/uuid"
)

type TypeAndKey struct {
	Type string
	Key  string
}

type IsolationLevel int

// TxOptions specifies options for transactions
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

func (tk TypeAndKey) String() string {
	return fmt.Sprintf("%s-%s", tk.Type, tk.Key)
}

func IsTKLessThan(a TypeAndKey, b TypeAndKey) bool {
	return a.String() < b.String()
}

// IDGenerator is an interface which wraps the creation of unique IDs
type IDGenerator interface {
	// New returns a new unique ID
	New() string
}

// UUIDGenerator is an implementation of IDGenerator which uses uuidv4
type UUIDGenerator struct{}

func (*UUIDGenerator) New() string {
	return uuid.New().String()
}
