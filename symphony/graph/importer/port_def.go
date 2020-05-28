// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"

	"go.uber.org/zap"
)

// processPortDefinitionsCSV imports port types and assign them to equipments (from CSV file to DB)
func (m *importer) processPortDefinitionsCSV(w http.ResponseWriter, r *http.Request) {
	log := m.logger.For(r.Context())
	log.Debug("PortDefinitions- started")
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

		portNameIndex := findIndex(firstLine, "Port_ID")
		if portNameIndex == -1 {
			errorReturn(w, "Couldn't find 'Port_ID' title", log, nil)
			return
		}
		portTypeIndex := findIndex(firstLine, "Port_Type")
		portBWIndex := findIndex(firstLine, "Port_Bandwidth")
		if portBWIndex == -1 {
			errorReturn(w, "Couldn't find 'Port_Bandwidth' title", log, nil)
			return
		}
		portLabelIndex := findIndex(firstLine, "Port_Visible_Label")
		if portLabelIndex == -1 {
			errorReturn(w, "Couldn't find 'Port_Visible_Label' title", log, nil)
			return
		}
		equipmentTypeNameIndex := findIndex(firstLine, "Equipment_Type")
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
			name := line[portNameIndex]
			equipmentTypeName := line[equipmentTypeNameIndex]
			if strings.HasPrefix(name, "<") {
				log.Info("skipping dynamic port config", zap.String("equipment", equipmentTypeName))
				continue
			}
			var portType string
			if portTypeIndex != -1 {
				portType = line[portTypeIndex]
			}

			var potyTypeObj *ent.EquipmentPortType
			potyTypeObj, err = m.ClientFrom(ctx).EquipmentPortType.Query().
				Where(equipmentporttype.Name(portType)).
				Only(ctx)
			if err != nil && !ent.IsNotFound(err) {
				log.Info("cant fetch port Type", zap.String("portType", portType))
				continue
			}
			if potyTypeObj == nil {
				potyTypeObj, err = m.ClientFrom(ctx).
					EquipmentPortType.
					Create().
					SetName(portType).
					Save(ctx)
				if err != nil {
					log.Info("cant create port Type", zap.String("portType", portType), zap.Error(err))
					continue
				}
			}

			equipTypeID := equipmentTypeNameToID[equipmentTypeName]
			if equipTypeID == 0 {
				log.Warn("cannot find equipment of port - creating new",
					zap.String("name", name),
					zap.String("equipment", equipmentTypeName),
				)
				equipTypeID = m.getOrCreateEquipmentType(ctx, equipmentTypeName, 0, "", 0, nil).ID
				equipmentTypeNameToID[equipmentTypeName] = equipTypeID
				continue
			}
			if !m.ClientFrom(ctx).EquipmentType.Query().
				Where(equipmenttype.ID(equipTypeID)).
				QueryPortDefinitions().
				Where(equipmentportdefinition.Name(name)).
				ExistX(ctx) {
				if _, err := m.ClientFrom(ctx).EquipmentPortDefinition.
					Create().
					SetName(name).
					SetEquipmentPortType(potyTypeObj).
					SetVisibilityLabel(line[portLabelIndex]).
					SetBandwidth(line[portBWIndex]).
					SetEquipmentTypeID(equipTypeID).
					Save(ctx); err != nil {
					log.Warn("cannot save port",
						zap.String("name", name),
						zap.String("type", portType),
						zap.Error(err),
					)
				}
			} else {
				log.Debug("port definition exists",
					zap.String("name", name),
					zap.Int("type", equipTypeID),
				)
			}
		}
	}
	log.Debug("PortDefinitions- Done")
}
