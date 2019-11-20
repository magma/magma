// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"errors"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
)

// nolint: unparam
func (m *importer) validateAllLocationTypeExist(ctx context.Context, offset int, locations []string, ignoreHierarchy bool) error {
	currIndex := -1
	ic := getImportContext(ctx)
	for i, locName := range locations {
		lt, err := m.ClientFrom(ctx).LocationType.Query().Where(locationtype.Name(locName)).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errors.New("location type not found, create it: + " + locName)
			}
			return err
		}
		if !ignoreHierarchy {
			if currIndex >= lt.Index {
				return errors.New("location types are not in the right order on the first line. edit the index and export again")
			}
			currIndex = lt.Index
		}
		ic.indexToLocationTypeID[offset+i] = lt.ID
	}
	return nil
}
