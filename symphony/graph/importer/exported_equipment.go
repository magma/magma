// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const minimalEquipmentLineLength = 9

// processExportedEquipment imports equipment csv generated from the export feature
// nolint: staticcheck, dupl
func (m *importer) processExportedEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	log.Debug("Exported Equipment - started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	count, numRows := 0, 0

	for fileName := range r.MultipartForm.File {
		first, reader, err := m.newReader(fileName, r)
		importHeader := NewImportHeader(first, ImportEntityEquipment)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file: %q. file name: %q", err, fileName), http.StatusInternalServerError)
			return
		}
		//
		//	populating:
		//	indexToLocationTypeID
		//
		if err = m.inputValidations(ctx, importHeader); err != nil {
			log.Warn("first line validation error", zap.Error(err))
			http.Error(w, fmt.Sprintf("first line validation error: %q", err), http.StatusBadRequest)
			return
		}
		//
		//	populating:
		//	equipmentTypeNameToID
		//	propNameToIndex
		//	equipmentTypeIDToProperties
		//
		err = m.populateEquipmentTypeNameToIDMap(ctx, importHeader, true)
		if err != nil {
			log.Warn("data fetching error", zap.Error(err))
			http.Error(w, fmt.Sprintf("data fetching error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		ic := getImportContext(ctx)
		for {
			untrimmedLine, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				continue
			}
			numRows++
			importLine := NewImportRecord(m.trimLine(untrimmedLine), importHeader)
			name := importLine.Name()
			equipTypName := importLine.TypeName()
			equipType, err := client.EquipmentType.Query().Where(equipmenttype.Name(equipTypName)).Only(ctx)
			if err != nil {
				log.Warn("couldn't find equipment type", zap.Error(err), zap.String("equipment_type", equipTypName))
				http.Error(w, fmt.Sprintf("couldn't find equipment type %q (row #%d). %q ", equipTypName, numRows, err), http.StatusBadRequest)
				return
			}

			externalID := importLine.ExternalID()
			id := importLine.ID()
			if id == "" {
				// new equip
				parentLoc, err := m.verifyOrCreateLocationHierarchy(ctx, importLine)
				if err != nil {
					log.Warn("creating location hierarchy", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("creating location hierarchy (row #%d). %q", numRows, err), http.StatusBadRequest)
					return
				}
				parentEquipmentID, positionDefinitionID, err := m.getPositionDetailsIfExists(ctx, parentLoc, importLine)
				if err != nil {
					log.Warn("creating equipment hierarchy", zap.Error(err), zap.Int("line_number", numRows), importLine.ZapField())
					http.Error(w, fmt.Sprintf("creating equipment hierarchy (row #%d). %q", numRows, err), http.StatusBadRequest)
					return
				}
				if parentEquipmentID != nil && positionDefinitionID != nil {
					parentLoc = nil
				}
				var propInputs []*models.PropertyInput
				if importLine.Len() > importHeader.PropertyStartIdx() {
					propInputs, err = m.validatePropertiesForEquipmentType(ctx, importLine, equipType)
					if err != nil {
						log.Warn("validating property for type", zap.Error(err))
						http.Error(w, fmt.Sprintf("validating property for type %q (row #%d). %q", equipType.Name, numRows, err.Error()), http.StatusBadRequest)
						return
					}
				}
				pos, err := resolverutil.GetOrCreatePosition(ctx, m.ClientFrom(ctx), parentEquipmentID, positionDefinitionID)
				if err != nil {
					log.Warn("creating equipment position", zap.Error(err), zap.Int("line_number", numRows), importLine.ZapField())
					http.Error(w, fmt.Sprintf("creating equipment position (row #%d). %q", numRows, err), http.StatusBadRequest)
					return
				}
				equip, created, err := m.getOrCreateEquipment(ctx, m.r.Mutation(), name, equipType, &externalID, parentLoc, pos, propInputs)
				if err != nil {
					log.Warn("creating/fetching equipment", zap.Error(err), zap.Int("line_number", numRows), importLine.ZapField())
					http.Error(w, fmt.Sprintf("creating/fetching equipment (row #%d). %q", numRows, err), http.StatusBadRequest)
					return
				}
				if created {
					count++
					log.Warn(fmt.Sprintf("(row #%d) creating equipment", numRows), zap.String("name", equip.Name), zap.String("id", equip.ID))
				} else {
					errorReturn(w, "Equipment "+equip.Name+" already exists under location/position", log, nil)
					return
				}
			} else {
				// existingEquip
				equipment, err := m.validateLineForExistingEquipment(ctx, id, importLine)
				if err != nil {
					log.Warn("validating existing equipment", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("%q: validating existing equipment: id %q (row #%d)", err, id, numRows), http.StatusBadRequest)
					return
				}
				typ := equipment.QueryType().OnlyX(ctx)
				props := ic.equipmentTypeIDToProperties[typ.ID]
				var inputs []*models.PropertyInput
				for _, propName := range props {
					inp, err := importLine.GetPropertyInput(m.ClientFrom(ctx), ctx, typ, propName)
					propType := typ.QueryPropertyTypes().Where(propertytype.Name(propName)).OnlyX(ctx)
					if err != nil {
						log.Warn("getting property input", zap.Error(err), importLine.ZapField())
						http.Error(w, fmt.Sprintf("%q: getting property input: prop %q (row #%d)", err, propName, numRows), http.StatusBadRequest)
						return
					}
					propID, err := equipment.QueryProperties().Where(property.HasTypeWith(propertytype.ID(propType.ID))).OnlyID(ctx)
					if err != nil {
						if !ent.IsNotFound(err) {
							log.Warn("property fetching error", zap.Error(err), importLine.ZapField())
							http.Error(w, fmt.Sprintf("%q: property fetching error: property name %q (row #%d)", err, propName, numRows), http.StatusBadRequest)
							return
						}
					} else {
						inp.ID = &propID
					}
					inputs = append(inputs, inp)
				}
				count++
				_, err = m.r.Mutation().EditEquipment(ctx, models.EditEquipmentInput{ID: id, Name: name, Properties: inputs, ExternalID: &externalID})
				if err != nil {
					log.Warn("editing equipment", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("editing equipment: id %q (row #%d). %q: ", id, numRows, err), http.StatusBadRequest)
					return
				}
			}
		}
	}
	log.Debug("Exported Equipment - Done")
	w.WriteHeader(http.StatusOK)
	err := writeSuccessMessage(w, count, numRows)
	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
}

func (m *importer) validateLineForExistingEquipment(ctx context.Context, equipID string, importLine ImportRecord) (*ent.Equipment, error) {
	equipment, err := m.r.Query().Equipment(ctx, equipID)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment")
	}
	typ := equipment.QueryType().OnlyX(ctx)
	if typ.Name != importLine.TypeName() {
		return nil, errors.Errorf("wrong equipment type. should be %v, but %v", importLine.TypeName(), typ.Name)
	}
	err = m.verifyPositionHierarchy(ctx, equipment, importLine)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching positions hierarchy")
	}
	err = m.validateLocationHierarchy(ctx, equipment, importLine)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (m *importer) inputValidations(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	prnt3Idx := importHeader.prnt3Idx
	if len(firstLine) < minimalEquipmentLineLength {
		return errors.New("first line too short. should include: 'Equipment ID', 'Equipment Name', 'Equipment Type', 'External ID' location types and parents")
	}
	locStart, _ := importHeader.LocationsRangeIdx()
	if !equal(firstLine[:locStart], []string{"Equipment ID", "Equipment Name", "Equipment Type", "External ID"}) {
		return errors.New("first line misses sequence; 'Equipment ID', 'Equipment Name' or 'Equipment Type' , 'External ID'")
	}
	if !equal(firstLine[prnt3Idx:importHeader.PropertyStartIdx()], []string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}) {
		return errors.New("first line misses sequence: 'Parent Equipment(3)', 'Position (3)', 'Parent Equipment (2)', 'Position (2)', 'Parent Equipment' or 'Equipment Position'")
	}
	err := m.validateAllLocationTypeExist(ctx, importHeader.ExternalIDIdx()+1, importHeader.LocationTypesRangeArr(), false)
	return err
}

func (m *importer) validatePropertiesForEquipmentType(ctx context.Context, line ImportRecord, equipType *ent.EquipmentType) ([]*models.PropertyInput, error) {
	ic := getImportContext(ctx)
	var pInputs []*models.PropertyInput
	propTypeNames := ic.equipmentTypeIDToProperties[equipType.ID]
	for _, ptypeName := range propTypeNames {
		pInput, err := line.GetPropertyInput(m.ClientFrom(ctx), ctx, equipType, ptypeName)
		if err != nil {
			return nil, err
		}
		pInputs = append(pInputs, pInput)
	}
	return pInputs, nil
}
