// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"io"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"

	"go.uber.org/zap"
)

// processPositionDefinitionsCSV imports position types and assign them to equipment types (from CSV file to DB)
func (m *importer) processPositionDefinitionsCSV(w http.ResponseWriter, r *http.Request) {
	log := m.logger.For(r.Context())
	log.Debug("PositionDefinitions- started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	for fileName := range r.MultipartForm.File {
		firstLine, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusUnprocessableEntity)
			return
		}

		_ = m.populateEquipmentTypeNameToIDMapGeneral(ctx, firstLine, false)
		equipmentTypeNameToID := getImportContext(ctx).equipmentTypeNameToID

		positionNameIndex := findIndexForSimilar(firstLine, "position_name")
		if positionNameIndex == -1 {
			errorReturn(w, "Couldn't find 'position_name' title", log, nil)
			return
		}
		positionLabelIndex := findIndexForSimilar(firstLine, "Position_Visible_Label")
		if positionLabelIndex == -1 {
			errorReturn(w, "Couldn't find 'Position_Visible_Label' title", log, nil)
			return
		}
		equipmentTypeNameIndex := findIndexForSimilar(firstLine, "Equipment_Type")
		if equipmentTypeNameIndex == -1 {
			errorReturn(w, "Couldn't find 'Equipment_Type' title", log, nil)
			return
		}
		for {
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				continue
			}
			name := line[positionNameIndex]
			equipmentTypeName := line[equipmentTypeNameIndex]
			positionLabel := line[positionLabelIndex]
			equipTypeID := equipmentTypeNameToID[equipmentTypeName]

			if equipTypeID == "" {
				log.Warn("cannot get equipment for position",
					zap.String("name", name),
					zap.String("equipment", equipmentTypeName),
				)
				continue
			}
			if !m.ClientFrom(ctx).EquipmentType.Query().
				Where(equipmenttype.ID(equipTypeID)).
				QueryPositionDefinitions().
				Where(equipmentpositiondefinition.Name(name)).
				ExistX(ctx) {
				if _, err := m.ClientFrom(ctx).EquipmentPositionDefinition.
					Create().
					SetName(name).
					SetVisibilityLabel(positionLabel).
					SetEquipmentTypeID(equipTypeID).
					Save(ctx); err != nil {
					log.Warn("cannot create position definition",
						zap.String("name", name),
						zap.String("label", positionLabel),
						zap.Error(err),
					)
				}
			} else {
				log.Debug("position definition exists",
					zap.String("name", name),
					zap.String("equipment", equipmentTypeName),
				)
			}
		}
	}
	log.Debug("PositionDefinitions- Done")
}
