// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ImportRecord struct {
	line  []string
	title ImportHeader
}

func NewImportRecord(line []string, title ImportHeader) ImportRecord {
	return ImportRecord{
		line:  line,
		title: title,
	}
}

func (l ImportRecord) ZapField() zap.Field {
	return zap.Strings("line", l.line)
}

func (l ImportRecord) Len() int {
	return len(l.line)
}

func (l ImportRecord) Header() ImportHeader {
	return l.title
}

func (l ImportRecord) GetPropertyInput(ctx context.Context, et *ent.EquipmentType, proptypeName string) (*models.PropertyInput, error) {
	ptyp, err := et.QueryPropertyTypes().Where(propertytype.Name(proptypeName)).Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "property type does not exist %q", proptypeName)
	}
	idx := l.title.Find(proptypeName)
	if idx == -1 {
		return nil, nil
	}
	return getPropInput(*ptyp, l.line[idx])
}

func (l ImportRecord) ID() string {
	return l.line[0]
}

func (l ImportRecord) Name() string {
	return l.line[1]
}

func (l ImportRecord) TypeName() string {
	return l.line[2]
}

func (l ImportRecord) ThirdParent() string {
	return l.line[l.title.ThirdParentIdx()]
}

func (l ImportRecord) SecondParent() string {
	return l.line[l.title.SecondParentIdx()]
}

func (l ImportRecord) DirectParent() string {
	return l.line[l.title.DirectParentIdx()]
}

func (l ImportRecord) Position() string {
	return l.line[l.title.PositionIdx()]
}

func (l ImportRecord) LocationsRangeArr() []string {
	s, e := l.title.LocationsRangeIdx()
	return l.line[s:e]
}

func (l ImportRecord) PropertiesSlice() []string {
	return l.line[l.title.PropertyStartIdx():]
}
