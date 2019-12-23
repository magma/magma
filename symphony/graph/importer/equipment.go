// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// processEquipmentCSV  imports equipment from CSV file to DB
func (m *importer) processEquipmentCSV(w http.ResponseWriter, r *http.Request) {
	log := m.log.For(r.Context())
	log.Debug("Equipment- started")
	instance := true
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	uploadLocationProps := false
	var err error
	boolData, ok := r.URL.Query()["uploadLocationProps"]
	if ok {
		uploadLocationProps, err = strconv.ParseBool(boolData[0])
		if err != nil {
			uploadLocationProps = false
		}
	}

	ctx := r.Context()
	for fileName := range r.MultipartForm.File {
		firstLine, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusInternalServerError)
			return
		}
		parentLocationIndex := findIndex(firstLine, "Location_ID")
		fullLocationPath := parentLocationIndex == -1
		equipmentNameIdx := findIndex(firstLine, "Equipment Name")
		equipmentTypeNameIdx := findIndex(firstLine, "Equipment Type")

		m.populateIndexToLocationTypeMap(ctx, firstLine, false)
		_ = m.populateEquipmentTypeNameToIDMap(ctx, NewImportHeader(firstLine, ImportEntityEquipment), true)
		ic := getImportContext(ctx)
		equipmentTypeIDToProperties := ic.equipmentTypeIDToProperties
		equipmentTypeNameToID := ic.equipmentTypeNameToID
		i := 0
		for {
			untrimmedLine, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				http.Error(w, fmt.Sprintf("cannot read row #%d", i), http.StatusInternalServerError)
				return
			}
			i++
			line := m.trimLine(untrimmedLine)
			if equipmentNameIdx == -1 || equipmentTypeNameIdx == -1 {
				log.Warn("'Equipment Name' or 'equipTypeName' columns are missing")
				http.Error(w, "'Equipment Name' or 'equipTypeName' columns are missing", http.StatusInternalServerError)
				return
			}
			equipName := line[equipmentNameIdx]
			equipTypeName := line[equipmentTypeNameIdx]
			equipTypeID := equipmentTypeNameToID[equipTypeName]

			if exists, err := m.ClientFrom(ctx).EquipmentType.Query().
				Where(equipmenttype.ID(equipTypeID)).
				Exist(ctx); err != nil || !exists {
				log.Warn("cannot find equipment type", zap.String("equip type name", equipTypeName), zap.Error(err))
				http.Error(w, fmt.Sprintf("row %d: cannot find equipment type %q", i, equipTypeName), http.StatusInternalServerError)
				return
			}
			var locationID string
			if fullLocationPath {
				locationID, err = m.getOrCreateEquipmentLocationByFullPath(ctx, line, firstLine, uploadLocationProps)
				if err != nil {
					log.Warn("cannot find or create location for equipment by full path", zap.String("name", equipName), zap.Error(err))
					http.Error(w, fmt.Sprintf("row %d: cannot find or create location for equipment '%q' by full path", i, equipName), http.StatusInternalServerError)
					return
				}
			} else {
				locName := line[parentLocationIndex]
				locationID, err = m.getLocationIDByName(ctx, locName)
				if err != nil {
					log.Warn("[SKIP]cannot find one location for equipment", zap.String("name", equipName), zap.Error(err), zap.String("location name", locName))
					http.Error(w, fmt.Sprintf("row %d: [SKIP]cannot find one location ('%q) for equipment '%q' by full path", i, locName, equipName), http.StatusInternalServerError)
					return
				}
			}
			var propertyInput []*models.PropertyInput
			propsKeys := equipmentTypeIDToProperties[equipTypeID]
			for _, key := range propsKeys {
				valIndex := findIndex(firstLine, key)
				if valIndex == -1 {
					continue
				}

				ptype, err := m.getOrCreatePropTypeForEquipment(ctx, equipTypeID, key)
				if err != nil {
					continue
				}
				// TODO T40408163. get "Type" from model
				propertyInput = append(propertyInput, &models.PropertyInput{
					PropertyTypeID:     ptype.ID,
					StringValue:        &line[valIndex],
					IsInstanceProperty: &instance,
				})
			}

			et, err := m.ClientFrom(ctx).EquipmentType.Query().Where(equipmenttype.ID(equipTypeID)).Only(ctx)
			if err != nil {
				log.Warn("[SKIP]cannot find equipment type by ID", zap.String("equip type id", equipTypeID), zap.Error(err))
				http.Error(w, fmt.Sprintf("row %d: [SKIP]cannot find equipment type by ID (%q)", i, equipTypeID), http.StatusInternalServerError)
				return
			}
			loc, err := m.ClientFrom(ctx).Location.Get(ctx, locationID)
			if err != nil {
				log.Warn("[SKIP]cannot find location type by ID", zap.String("loc type id", locationID), zap.Error(err))
				http.Error(w, fmt.Sprintf("row %d: [SKIP]cannot find location type by ID (%q)", i, locationID), http.StatusInternalServerError)
				return
			}

			_, _, err = m.getOrCreateEquipment(ctx, m.r.Mutation(), equipName, et, nil, loc, nil, propertyInput)
			if err != nil {
				log.Warn("error while creating equipment", zap.String("equipment name", equipName), zap.Error(err))
				http.Error(w, fmt.Sprintf("row %d: error while creating equipment (namew=%s): %s", i, equipName, err.Error()), http.StatusInternalServerError)
				return
			}
		}
	}
	log.Debug("Equipment- Done")
	w.WriteHeader(http.StatusOK)
}
