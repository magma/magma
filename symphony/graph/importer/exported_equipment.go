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

	"github.com/AlekSi/pointer"

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
	log := m.logger.For(ctx)
	nextLineToSkipIndex := -1
	client := m.ClientFrom(ctx)

	log.Debug("Exported Equipment - started")
	var (
		err                    error
		modifiedCount, numRows int
		errs                   Errors
		commitRuns             []bool
	)
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	skipLines, verifyBeforeCommit, err := m.parseImportArgs(r)
	if err != nil {
		errorReturn(w, "can't parse form or arguments", log, err)
		return
	}

	if pointer.GetBool(verifyBeforeCommit) {
		commitRuns = []bool{false, true}
	} else {
		commitRuns = []bool{true}
	}
	startSaving := false

	for fileName := range r.MultipartForm.File {
		first, _, err := m.newReader(fileName, r)
		if err != nil {
			errorReturn(w, fmt.Sprintf("cannot handle file: %q", fileName), log, err)
			return
		}
		importHeader, err := NewImportHeader(first, ImportEntityEquipment)
		if err != nil {
			errorReturn(w, fmt.Sprintf("error on header"), log, err)
			return
		}
		//
		//	populating:
		//	indexToLocationTypeID
		//
		if err = m.inputValidations(ctx, importHeader); err != nil {
			errorReturn(w, "first line validation error", log, err)
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
			errorReturn(w, "data fetching error", log, err)
			return
		}

		for _, commit := range commitRuns {
			// if we encounter errors on the "verifyBefore" flow - don't run the commit=true phase
			if commit && pointer.GetBool(verifyBeforeCommit) && len(errs) != 0 {
				break
			} else if commit && len(errs) == 0 {
				startSaving = true
			}
			if len(skipLines) > 0 {
				nextLineToSkipIndex = 0
			}

			numRows, modifiedCount = 0, 0
			_, reader, err := m.newReader(fileName, r)
			if err != nil {
				errorReturn(w, fmt.Sprintf("cannot handle file: %q", fileName), log, err)
				return
			}
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
				if shouldSkipLine(skipLines, numRows, nextLineToSkipIndex) {
					log.Warn("skipping line", zap.Error(err), zap.Int("line_number", numRows))
					nextLineToSkipIndex++
					continue
				}

				importLine := NewImportRecord(m.trimLine(untrimmedLine), importHeader)
				name := importLine.Name()
				equipTypName := importLine.TypeName()
				equipType, err := client.EquipmentType.Query().Where(equipmenttype.Name(equipTypName)).Only(ctx)
				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("couldn't find equipment type %q", equipTypName)})
					continue
				}

				externalID := importLine.ExternalID()
				id := importLine.ID()
				if id == "" {
					// new equip
					parentLoc, err := m.verifyOrCreateLocationHierarchy(ctx, importLine, commit, nil)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/verifying equipment location hierarchy"})
						continue
					} else if parentLoc == nil && !commit {
						continue
					}

					parentEquipmentID, positionDefinitionID, err := m.getPositionDetailsIfExists(ctx, parentLoc, importLine, true)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/verifying equipment hierarchy"})
						continue

					}
					if parentEquipmentID != nil && positionDefinitionID != nil {
						parentLoc = nil
					}
					var propInputs []*models.PropertyInput
					if importLine.Len() > importHeader.PropertyStartIdx() {
						propInputs, err = m.validatePropertiesForEquipmentType(ctx, importLine, equipType)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating property for type %q", equipType.Name)})
							continue
						}
					}

					var pos *ent.EquipmentPosition
					var created bool
					var equ *ent.Equipment
					if commit {
						pos, err = resolverutil.GetOrCreatePosition(ctx, m.ClientFrom(ctx), parentEquipmentID, positionDefinitionID, true)
					} else {
						pos, err = resolverutil.ValidateAndGetPositionIfExists(ctx, client, parentEquipmentID, positionDefinitionID, true)
					}
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/fetching equipment position"})
						continue
					}
					if commit {
						_, created, err = m.getOrCreateEquipment(ctx, m.r.Mutation(), name, equipType, &externalID, parentLoc, pos, propInputs)
						if created {
							modifiedCount++
						} else if err == nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: "Equipment already exists under location/position"})
							log.Info("Row " + strconv.FormatInt(int64(numRows), 10) + ": Equipment already exists under location/position")
							continue
						}
					} else {
						equ, err = m.getEquipmentIfExist(ctx, m.r.Mutation(), name, equipType, &externalID, parentLoc, pos, propInputs)
						if equ != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: "Equipment already exists under location/position"})
							continue
						}
					}
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/fetching equipment"})
						continue
					}
				} else {
					// existing equip
					equipment, err := m.validateLineForExistingEquipment(ctx, id, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error validating equipment line"})
						continue
					}
					inputs, msg, err := m.getEquipmentPropertyInputs(ctx, importLine, equipment)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: msg})
						continue
					}
					if commit {
						modifiedCount++
						_, err = m.r.Mutation().EditEquipment(ctx, models.EditEquipmentInput{ID: id, Name: name, Properties: inputs, ExternalID: &externalID})
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("editing equipment: id %v", id)})
							continue
						}
					}
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0, startSaving)
	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
	log.Debug("Exported Equipment - Done", zap.Any("errors list", errs), zap.Int("all_lines", numRows), zap.Int("edited_added_rows", modifiedCount))
}

func (m *importer) getEquipmentPropertyInputs(ctx context.Context, importLine ImportRecord, equipment *ent.Equipment) ([]*models.PropertyInput, string, error) {
	typ := equipment.QueryType().OnlyX(ctx)
	ic := getImportContext(ctx)
	props := ic.equipmentTypeIDToProperties[typ.ID]
	var inputs []*models.PropertyInput
	for _, propName := range props {
		inp, err := importLine.GetPropertyInput(m.ClientFrom(ctx), ctx, typ, propName)
		propType := typ.QueryPropertyTypes().Where(propertytype.Name(propName)).OnlyX(ctx)
		if err != nil {
			return nil, fmt.Sprintf("getting property input: prop %v", propName), err
		}
		propID, err := equipment.QueryProperties().Where(property.HasTypeWith(propertytype.ID(propType.ID))).OnlyID(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, fmt.Sprintf("property fetching error: property name %v", propName), err
			}
		} else {
			inp.ID = &propID
		}
		inputs = append(inputs, inp)
	}
	return inputs, "", nil
}

func (m *importer) validateLineForExistingEquipment(ctx context.Context, equipID string, importLine ImportRecord) (*ent.Equipment, error) {
	equipment, err := m.r.Query().Equipment(ctx, equipID)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment")
	}
	typ := equipment.QueryType().OnlyX(ctx)
	if typ.Name != importLine.TypeName() {
		return nil, errors.Errorf("wrong equipment type. should be %v, but %v", typ.Name, importLine.TypeName())
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
		if pInput != nil {
			pInputs = append(pInputs, pInput)
		}
	}
	return pInputs, nil
}
