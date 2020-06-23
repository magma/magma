// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"time"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/pkg/ent"
)

// event log events.
const (
	EntMutation = "ent/mutation"
)

// LogEntry holds an information on a single ent mutation that happened
type LogEntry struct {
	UserName  string    `json:"user_name"`
	UserID    *int      `json:"user_id"`
	Time      time.Time `json:"time"`
	Type      string    `json:"type"`
	Operation ent.Op    `json:"operation"`
	PrevState *ent.Node `json:"prevState"`
	CurrState *ent.Node `json:"currState"`
}

func FindEdge(edges []*ent.Edge, val string) (*ent.Edge, bool) {
	for _, edge := range edges {
		if edge.Name == val && len(edge.IDs) > 0 {
			return edge, true
		}
	}
	return nil, false
}

func FindField(fields []*ent.Field, val string) (*ent.Field, bool) {
	for _, f := range fields {
		if f.Name == val {
			return f, true
		}
	}
	return nil, false
}

func GetDiffOfUniqueEdge(entry *LogEntry, edge string) (*int, *int, bool) {
	var newIntVal, oldsIntVal *int
	newEdges := entry.CurrState.Edges
	oldEdges := entry.PrevState.Edges
	newEdge, newFound := FindEdge(newEdges, edge)
	oldEdge, oldFound := FindEdge(oldEdges, edge)
	if newFound && len(newEdge.IDs) > 0 {
		newIntVal = &newEdge.IDs[0]
	}
	if oldFound && len(oldEdge.IDs) > 0 {
		oldsIntVal = &oldEdge.IDs[0]
	}
	shouldUpdate := (newFound != oldFound) || (newIntVal != nil && oldsIntVal != nil && *newIntVal != *oldsIntVal)
	return newIntVal, oldsIntVal, shouldUpdate
}

func GetStringDiffValuesField(entry *LogEntry, field string) (*string, *string, bool) {
	var newStrVal, oldStrVal *string
	newFields := entry.CurrState.Fields
	oldFields := entry.PrevState.Fields
	newField, newFound := FindField(newFields, field)
	oldField, oldFound := FindField(oldFields, field)
	if newFound && newField != nil {
		newStrVal = pointer.ToString(newField.MustGetString())
	}
	if oldFound && oldField != nil {
		oldStrVal = pointer.ToString(oldField.MustGetString())
	}

	shouldUpdate := (newFound != oldFound) || (newStrVal != nil && oldStrVal != nil && *newStrVal != *oldStrVal)
	return newStrVal, oldStrVal, shouldUpdate
}
