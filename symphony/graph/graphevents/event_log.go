// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphevents

import (
	"context"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/pkg/log"

	"go.uber.org/zap"
)

type eventLog struct {
	logger log.Logger
}

func getChangedFields(entry event.LogEntry) []string {
	values := make(map[string]*string)
	var fields []string
	if entry.PrevState != nil {
		for _, f := range entry.PrevState.Fields {
			values[f.Name] = &f.Value
		}
	}
	if entry.CurrState != nil {
		for _, f := range entry.CurrState.Fields {
			val, exists := values[f.Name]
			if exists {
				if f.Value != *val {
					fields = append(fields, f.Name)
				}
				values[f.Name] = nil
			} else {
				fields = append(fields, f.Name)
			}
		}
	}
	for name, val := range values {
		if val != nil {
			fields = append(fields, name)
		}
	}
	return fields
}

func getEntIdentifiers(entry event.LogEntry) (int, string) {
	if entry.PrevState != nil {
		return entry.PrevState.ID, entry.PrevState.Type
	}
	return entry.CurrState.ID, entry.CurrState.Type
}

func (e eventLog) Handle(ctx context.Context, entry event.LogEntry) {
	changedFields := getChangedFields(entry)
	id, typ := getEntIdentifiers(entry)

	e.logger.For(ctx).Info(
		"ent mutation",
		zap.String("user_name", entry.UserName),
		zap.Any("operation", entry.Operation),
		zap.Int("id", id),
		zap.String("type", typ),
		zap.Strings("changed_fields", changedFields))
}
