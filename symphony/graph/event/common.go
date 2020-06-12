// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"go.uber.org/zap"
)

func (e *Eventer) hookWithLog(handler func(context.Context, pubsub.LogEntry) error, next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		v := viewer.FromContext(ctx)
		if m.Op().Is(ent.OpDelete) || m.Op().Is(ent.OpUpdate) || v == nil {
			return next.Mutate(ctx, m)
		}
		entry := pubsub.LogEntry{
			UserName:  v.Name(),
			Operation: m.Op(),
			Time:      time.Now(),
		}
		if v, ok := v.(*viewer.UserViewer); ok {
			entry.UserID = &v.User().ID
		}
		if !m.Op().Is(ent.OpCreate) {
			if prevNoder, ok := m.(ent.Noder); ok {
				node, err := prevNoder.Node(privacy.DecisionContext(ctx, privacy.Allow))
				if err != nil {
					if !ent.IsNotFound(err) {
						e.Logger.For(ctx).Error("query mutation previous value", zap.Error(err))
					}
					return next.Mutate(ctx, m)
				}
				entry.PrevState = node
			}
		}
		value, err := next.Mutate(ctx, m)
		if err != nil {
			return value, err
		}
		if !m.Op().Is(ent.OpDeleteOne) {
			if currNoder, ok := value.(ent.Noder); ok {
				node, err := currNoder.Node(privacy.DecisionContext(ctx, privacy.Allow))
				if err != nil {
					e.Logger.For(ctx).Error("query mutation current value", zap.Error(err))
					return value, err
				}
				entry.CurrState = node
			}
		}
		err = handler(ctx, entry)
		if err != nil {
			return value, err
		}
		return value, nil
	})
}

func findEdge(edges []*ent.Edge, val string) (*ent.Edge, bool) {
	for _, edge := range edges {
		if edge.Name == val && len(edge.IDs) > 0 {
			return edge, true
		}
	}
	return nil, false
}

func findField(fields []*ent.Field, val string) (*ent.Field, bool) {
	for _, f := range fields {
		if f.Name == val {
			return f, true
		}
	}
	return nil, false
}

func getDiffOfUniqueEdge(entry *pubsub.LogEntry, edge string) (*int, *int, bool) {
	var newIntVal, oldsIntVal *int
	newEdges := entry.CurrState.Edges
	oldEdges := entry.PrevState.Edges
	newEdge, newFound := findEdge(newEdges, edge)
	oldEdge, oldFound := findEdge(oldEdges, edge)
	if newFound && len(newEdge.IDs) > 0 {
		newIntVal = &newEdge.IDs[0]
	}
	if oldFound && len(oldEdge.IDs) > 0 {
		oldsIntVal = &oldEdge.IDs[0]
	}
	shouldUpdate := (newFound != oldFound) || (newIntVal != nil && oldsIntVal != nil && *newIntVal != *oldsIntVal)
	return newIntVal, oldsIntVal, shouldUpdate
}

func getStringDiffValuesField(entry *pubsub.LogEntry, field string) (*string, *string, bool) {
	var newStrVal, oldStrVal *string
	newFields := entry.CurrState.Fields
	oldFields := entry.PrevState.Fields
	newField, newFound := findField(newFields, field)
	oldField, oldFound := findField(oldFields, field)
	if newFound && newField != nil {
		newStrVal = pointer.ToString(newField.MustGetString())
	}
	if oldFound && oldField != nil {
		oldStrVal = pointer.ToString(oldField.MustGetString())
	}

	shouldUpdate := (newFound != oldFound) || (newStrVal != nil && oldStrVal != nil && *newStrVal != *oldStrVal)
	return newStrVal, oldStrVal, shouldUpdate
}
