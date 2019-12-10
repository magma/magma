// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const minimalLineLength = 7

// processExportedEquipment imports equipment csv generated from the export feature
// nolint: staticcheck
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
		importHeader := NewImportHeader(first)
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
					propInputs, err = m.validatePropertiesForType(ctx, importLine, equipType)
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
				equip, created := m.getOrCreateEquipment(ctx, m.r.Mutation(), name, equipType, parentLoc, pos, propInputs)
				if created {
					count++
					log.Warn(fmt.Sprintf("(row #%d) creating equipment", numRows), zap.String("name", equip.Name), zap.String("id", equip.ID))
				} else {
					log.Warn(fmt.Sprintf("(row #%d) [SKIP]equipment existed under location", numRows), zap.String("name", equip.Name), zap.String("id", equip.ID))
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
					inp, err := importLine.GetPropertyInput(ctx, typ, propName)
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
				_, err = m.r.Mutation().EditEquipment(ctx, models.EditEquipmentInput{ID: id, Name: name, Properties: inputs})
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
	msg := fmt.Sprintf("Created %q instances, out of %q", strconv.FormatInt(int64(count), 10), strconv.FormatInt(int64(numRows), 10))
	w.Write([]byte(msg))
}

func (m *importer) validateLineForExistingEquipment(ctx context.Context, equipID string, importLine ImportRecord) (*ent.Equipment, error) {
	equipment, err := m.r.Query().Equipment(ctx, equipID)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment")
	}
	typ := equipment.QueryType().OnlyX(ctx)
	if typ.Name != importLine.TypeName() {
		return nil, errors.Wrapf(err, "wrong equipment type. should be %q, but %q", importLine.TypeName(), typ.Name)
	}
	posHierarchy, err := m.r.Equipment().PositionHierarchy(ctx, equipment)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching positions hierarchy")
	}
	length := len(posHierarchy)
	if length > 0 {
		if length > 4 {
			// getting the last 4 positions (we currently support 4 on export)
			posHierarchy = posHierarchy[(length - 4):]
		}
		directPos := posHierarchy[length-1]

		defName := directPos.QueryDefinition().OnlyX(ctx).Name
		if defName != importLine.Position() {
			return nil, errors.Errorf("wrong position name. should be %q, but %q", importLine.Position(), defName)
		}
		pName := directPos.QueryParent().OnlyX(ctx).Name
		if pName != importLine.DirectParent() {
			return nil, errors.Errorf("wrong equipment parent name. should be %q, but %q", importLine.DirectParent(), pName)
		}
	}
	locs, err := m.r.Equipment().LocationHierarchy(ctx, equipment)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching location hierarchy")
	}
	prevIdx := 0
	for _, loc := range locs {
		currIdx := findIndex(importLine.line, strings.Trim(loc.Name, " "))
		if currIdx == -1 {
			return nil, errors.Errorf("missing location from hierarchy (%q)", loc.Name)
		}
		if prevIdx > currIdx {
			return nil, errors.Errorf("location not in the right order (%q)", loc.Name)
		}
		prevIdx = currIdx
	}
	return equipment, nil
}

func (m *importer) inputValidations(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	prnt3Idx := importHeader.prnt3Idx
	if len(firstLine) < minimalLineLength {
		return errors.New("first line too short. should include: 'Equipment ID', 'Equipment Name' or 'Equipment Type', location types and parents")
	}
	locStart, _ := importHeader.LocationsRangeIdx()
	if !equal(firstLine[:locStart], []string{"Equipment ID", "Equipment Name", "Equipment Type"}) {
		return errors.New("first line misses sequence; 'Equipment ID', 'Equipment Name' or 'Equipment Type'")
	}
	if !equal(firstLine[prnt3Idx:importHeader.PropertyStartIdx()], []string{"Parent Equipment (3)", "Parent Equipment (2)", "Parent Equipment", "Equipment Position"}) {
		return errors.New("first line misses sequence: 'Parent Equipment(3)', 'Parent Equipment (2)', 'Parent Equipment' or 'Equipment Position'")
	}
	err := m.validateAllLocationTypeExist(ctx, 3, importHeader.LocationTypesRangeArr(), false)
	return err
}

// nolint: unparam
func (m *importer) verifyOrCreateLocationHierarchy(ctx context.Context, l ImportRecord) (*ent.Location, error) {
	var currParentID *string
	var loc *ent.Location
	ic := getImportContext(ctx)
	locStart, _ := l.Header().LocationsRangeIdx()
	for i, locName := range l.LocationsRangeArr() {
		if locName == "" {
			continue
		}
		typID := ic.indexToLocationTypeID[i+locStart] // the actual index
		typ, err := m.r.Query().LocationType(ctx, typID)
		if err != nil {
			return nil, errors.Wrapf(err, "missing location type: id=%q", typID)
		}
		loc, _ = m.getOrCreateLocation(ctx, locName, 0.0, 0.0, typ, currParentID, nil, nil)
		currParentID = &loc.ID
	}
	if loc == nil {
		return nil, errors.Errorf("equipment with no locations specified. id:%q, name: %q", l.ID(), l.Name())
	}
	return loc, nil
}

func (m *importer) getPositionDetailsIfExists(ctx context.Context, parentLoc *ent.Location, importLine ImportRecord) (*string, *string, error) {
	l := importLine.line
	title := importLine.title
	if importLine.Position() == "" {
		return nil, nil, nil
	}
	var (
		equip  *ent.Equipment
		err    error
		errMsg string
	)
	for idx := title.prnt3Idx; idx < title.PositionIdx(); idx++ {
		if l[idx] == "" {
			continue
		}
		if equip == nil {
			equip, err = parentLoc.QueryEquipment().Where(equipment.Name(l[idx])).Only(ctx)
			errMsg = fmt.Sprintf("equipment %q not found under location %q", l[idx], parentLoc.Name)
		} else {
			equip, err = equip.QueryPositions().QueryAttachment().Where(equipment.Name(l[idx])).Only(ctx)
			errMsg = fmt.Sprintf("empty position %q not found under equipment %q", l[idx], l[idx-1])
		}
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, nil, errors.New(errMsg)
			}
			return nil, nil, err
		}
	}
	if equip == nil {
		return nil, nil, errors.Errorf("location/equipment/position mismatch %q, %q, %q", parentLoc.Name, importLine.DirectParent(), importLine.Position())
	}
	def, err := equip.QueryType().QueryPositionDefinitions().Where(equipmentpositiondefinition.Name(importLine.Position())).Only(ctx)
	if err != nil {
		return nil, nil, err
	}
	hasAttachment, err := equip.QueryPositions().
		Where(equipmentposition.HasDefinitionWith(equipmentpositiondefinition.ID(def.ID))).
		QueryAttachment().
		Exist(ctx)
	if err != nil {
		return nil, nil, err
	}
	if hasAttachment {
		return nil, nil, errors.Errorf("position %q already has attachment", importLine.Position())
	}
	return &equip.ID, &def.ID, nil
}

func (m *importer) validatePropertiesForType(ctx context.Context, line ImportRecord, equipType *ent.EquipmentType) ([]*models.PropertyInput, error) {
	ic := getImportContext(ctx)
	var pInputs []*models.PropertyInput
	propTypeNames := ic.equipmentTypeIDToProperties[equipType.ID]
	for _, ptypeName := range propTypeNames {
		pInput, err := line.GetPropertyInput(ctx, equipType, ptypeName)
		if err != nil {
			return nil, err
		}
		pInputs = append(pInputs, pInput)
	}
	return pInputs, nil
}
