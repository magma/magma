// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/frontier/ent"
)

// The AuditLogFunc type is an adapter to allow the use of ordinary
// function as AuditLog mutator.
type AuditLogFunc func(context.Context, *ent.AuditLogMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f AuditLogFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.AuditLogMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.AuditLogMutation", m)
	}
	return f(ctx, mv)
}

// The TenantFunc type is an adapter to allow the use of ordinary
// function as Tenant mutator.
type TenantFunc func(context.Context, *ent.TenantMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TenantFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.TenantMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TenantMutation", m)
	}
	return f(ctx, mv)
}

// The TokenFunc type is an adapter to allow the use of ordinary
// function as Token mutator.
type TokenFunc func(context.Context, *ent.TokenMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TokenFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.TokenMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TokenMutation", m)
	}
	return f(ctx, mv)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *ent.UserMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.UserMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.UserMutation", m)
	}
	return f(ctx, mv)
}

// On executes the given hook only of the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
//
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Op().Is(op) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
//
func Reject(op ent.Op) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Op().Is(op) {
				return nil, fmt.Errorf("%s operation is not allowed", m.Op())
			}
			return next.Mutate(ctx, m)
		})
	}
}
