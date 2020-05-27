/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"strings"

	"magma/orc8r/cloud/go/services/state"
)

// Subscription denotes a set of primary keys.
type Subscription struct {
	Type       string
	KeyMatcher KeyMatcher
}

// KeyMatcher indicates whether a particular key matches some pattern.
type KeyMatcher interface {
	Match(s string) bool
}

func (s Subscription) Match(id state.ID) bool {
	if typeMatch := s.Type == id.Type; !typeMatch {
		return false
	}
	return s.KeyMatcher.Match(id.DeviceID)
}

type matchAll struct{}
type matchExact struct{ exact string }
type matchPrefix struct{ prefix string }

// MatchAll is a singleton key matcher for matching all keys.
var MatchAll KeyMatcher = &matchAll{}

// NewMatchExact returns a new KeyMatcher that matches keys exactly matching exact.
func NewMatchExact(exact string) KeyMatcher { return &matchExact{exact: exact} }

// NewMatchPrefix returns a new KeyMatcher that matches keys prefixed with prefix.
func NewMatchPrefix(prefix string) KeyMatcher { return &matchPrefix{prefix: prefix} }

func (m *matchAll) Match(s string) bool    { return true }
func (m *matchExact) Match(s string) bool  { return s == m.exact }
func (m *matchPrefix) Match(s string) bool { return strings.HasPrefix(s, m.prefix) }
