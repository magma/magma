// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

const maxEquipmentParents = 3

// ImportEntity specifies an entity that can be imported
type ImportEntity string

const (
	// ImportEntityEquipment specifies an equipment for import
	ImportEntityEquipment ImportEntity = "EQUIPMENT"
	// ImportEntityPort specifies a port for import
	ImportEntityPort ImportEntity = "PORT"
	// ImportEntityLink specifies a link for import
	ImportEntityLink ImportEntity = "LINK"
	// ImportEntityPortInLink specifies a port sub-slice inside a link entity for import
	ImportEntityPortInLink ImportEntity = "PORT_IN_LINK"
	// ImportEntityService specifies a service for import
	ImportEntityService ImportEntity = "SERVICE"
	// ImportEntityLocation specifies a location for import
	ImportEntityLocation ImportEntity = "LOCATION"
)

// SuccessMessage is the type returns to client on success import
type SuccessMessage struct {
	MessageCode  int  `json:"messageCode"`
	SuccessLines int  `json:"successLines"`
	AllLines     int  `json:"allLines"`
	Committed    bool `json:"committed"`
}

type ReturnMessage struct {
	Summary SuccessMessage `json:"summary"`
	Errors  Errors         `json:"errors"`
}

func writeSuccessMessage(w http.ResponseWriter, success, all int, errs Errors, isSuccess, startSaving bool) error {
	w.Header().Set("Content-Type", "application/json")
	messageCode := int(SuccessfullyUploaded)
	if !isSuccess {
		messageCode = int(FailedToUpload)
	}
	return json.NewEncoder(w).Encode(
		ReturnMessage{
			Summary: SuccessMessage{
				MessageCode:  messageCode,
				SuccessLines: success,
				AllLines:     all,
				Committed:    startSaving,
			},
			Errors: errs,
		},
	)
}

// ErrorLine represents a line which failed to validate
type ErrorLine struct {
	Line    int    `json:"line"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ErrorLine represents a summary of the errors while uploading a CSV file
type Errors []ErrorLine

func getLinesToSkip(r *http.Request) ([]int, error) {
	var skipLines []int
	arg := r.FormValue("skip_lines")
	if arg != "" {
		err := json.Unmarshal([]byte(arg), &skipLines)
		if err != nil {
			return nil, err
		}
	}
	if len(skipLines) > 0 {
		skipLines = sortSlice(skipLines, true)
	}
	return skipLines, nil
}

func shouldSkipLine(a []int, currRow, nextLineToSkipIndex int) bool {
	if nextLineToSkipIndex >= 0 && nextLineToSkipIndex < len(a) {
		return currRow == a[nextLineToSkipIndex]
	}
	return false
}

func getVerifyBeforeCommitParam(r *http.Request) (*bool, error) {
	verifyBeforeCommit := false
	commitParam := r.FormValue("verify_before_commit")
	if commitParam != "" {
		err := json.Unmarshal([]byte(commitParam), &verifyBeforeCommit)
		if err != nil {
			return nil, err
		}
	}
	return &verifyBeforeCommit, nil
}

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

// nolint: unparam
func (m *importer) verifyOrCreateLocationHierarchy(ctx context.Context, l ImportRecord, commit bool, limit *int) (*ent.Location, error) {
	var currParentID *int
	var loc *ent.Location
	ic := getImportContext(ctx)

	locStart, indexToStopLoop := l.Header().LocationsRangeIdx()
	if limit != nil {
		indexToStopLoop = *limit
	}

	for i, locName := range l.LocationsRangeArr() {
		if locName == "" {
			continue
		}
		if i >= indexToStopLoop {
			break
		}
		typID := ic.indexToLocationTypeID[i+locStart] // the actual index
		typ, err := m.r.Query().LocationType(ctx, typID)
		if err != nil {
			return nil, errors.Wrapf(err, "missing location type: id=%q", typID)
		}
		if commit {
			loc, _, err = m.getOrCreateLocation(ctx, locName, 0.0, 0.0, typ, currParentID, nil, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "querying or creating location: id=%v", typID)
			}
		} else {
			loc, err = m.queryLocationForTypeAndParent(ctx, locName, typ, currParentID)

			if loc == nil {
				if !ent.IsNotFound(err) {
					return nil, errors.Wrapf(err, "querying or creating location name: %v", locName)
				}
				// no location but no error (dry run mode)
				return nil, nil
			}
		}
		currParentID = &loc.ID
	}
	if loc == nil && limit != nil {
		return nil, errors.Errorf("equipment with no locations specified. id:%q, name: %q", l.ID(), l.Name())
	}
	return loc, nil
}

func (m *importer) validateLocationHierarchy(ctx context.Context, equipment *ent.Equipment, importLine ImportRecord) error {
	locs, err := m.r.Equipment().LocationHierarchy(ctx, equipment)
	if err != nil {
		return errors.Wrapf(err, "fetching location hierarchy")
	}
	prevIdx := 0
	for _, loc := range locs {
		currIdx := findIndex(importLine.line, strings.Trim(loc.Name, " "))
		if currIdx == -1 {
			return errors.Errorf("missing location from hierarchy (%q)", loc.Name)
		}
		if prevIdx > currIdx {
			return errors.Errorf("location not in the right order (%q)", loc.Name)
		}
		prevIdx = currIdx
	}
	return nil
}

func (m *importer) verifyPositionHierarchy(ctx context.Context, equipment *ent.Equipment, importLine ImportRecord) error {
	posHierarchy, err := m.r.Equipment().PositionHierarchy(ctx, equipment)
	if err != nil {
		return errors.Wrapf(err, "fetching positions hierarchy for equipment")
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
			return errors.Errorf("wrong position name. should be %v, but %v", importLine.Position(), defName)
		}
		pName := directPos.QueryParent().OnlyX(ctx).Name
		if pName != importLine.DirectParent() {
			return errors.Errorf("wrong equipment parent name. should be %v, but %v", importLine.DirectParent(), pName)
		}
	}
	return nil
}

func (m *importer) getPositionDetailsIfExists(ctx context.Context, parentLoc *ent.Location, importLine ImportRecord, mustBeEmpty bool) (*int, *int, error) {
	l := importLine.line
	title := importLine.title
	if importLine.Position() == "" {
		return nil, nil, nil
	}
	var (
		equip        *ent.Equipment
		err          error
		errMsg       error
		positionName string
	)
	for idx := title.prnt3Idx; idx < title.PositionIdx(); idx += 2 {
		if l[idx] == "" {
			continue
		}
		if equip == nil {
			equip, err = parentLoc.QueryEquipment().Where(equipment.Name(l[idx])).Only(ctx)
			errMsg = fmt.Errorf("equipment %q not found under location %q", l[idx], parentLoc.Name)
		} else {
			equip, err = equip.QueryPositions().
				Where(equipmentposition.HasDefinitionWith(equipmentpositiondefinition.Name(positionName))).
				QueryAttachment().
				Where(equipment.Name(l[idx])).Only(ctx)
			errMsg = fmt.Errorf("position %q not found under equipment %q", positionName, l[idx])
		}
		positionName = l[idx+1]
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, nil, errMsg
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

	if mustBeEmpty {
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
	}
	return &equip.ID, &def.ID, nil
}

func (m *importer) validatePropertiesForPortType(ctx context.Context, line ImportRecord, portType *ent.EquipmentPortType, entity ImportEntity) ([]*models.PropertyInput, error) {
	var (
		pInputs   []*models.PropertyInput
		propTypes []*ent.PropertyType
		err       error
	)

	switch entity {
	case ImportEntityPort:
		propTypes, err = portType.QueryPropertyTypes().All(ctx)
	case ImportEntityLink:
		propTypes, err = portType.QueryLinkPropertyTypes().All(ctx)
	default:
		return nil, errors.New(fmt.Sprintf("ImportEntity not supported %s", entity))
	}
	if ent.MaskNotFound(err) != nil {
		return nil, errors.Wrap(err, "can't query property types for port type")
	}
	for _, ptype := range propTypes {
		ptypeName := ptype.Name
		pInput, err := line.GetPropertyInput(m.ClientFrom(ctx), ctx, portType, ptypeName)
		if err != nil {
			return nil, err
		}
		if pInput != nil {
			pInputs = append(pInputs, pInput)
		}
	}
	return pInputs, nil
}

func (m *importer) validatePort(ctx context.Context, portData PortData, port ent.EquipmentPort) error {
	def, err := port.QueryDefinition().Only(ctx)
	if err != nil {
		return errors.Wrapf(err, "fetching equipment port definition")
	}
	if def.Name != portData.Name {
		return errors.Errorf("wrong port type. should be %q, but %q", def.Name, portData.Name)
	}
	portType, err := def.QueryEquipmentPortType().Only(ctx)
	if ent.MaskNotFound(err) != nil {
		return errors.Wrapf(err, "fetching equipment port type")
	}
	var tempPortType string
	if ent.IsNotFound(err) {
		tempPortType = ""
	} else {
		tempPortType = portType.Name
	}
	if tempPortType != portData.TypeName {
		return errors.Errorf("wrong port type. should be %q, but %q", tempPortType, portData.TypeName)
	}

	equipment, err := port.QueryParent().Only(ctx)
	if err != nil {
		return errors.Wrapf(err, "fetching equipment for port")
	}
	if equipment.Name != portData.EquipmentName {
		return errors.Errorf("wrong equipment. should be %q, but %q", equipment.Name, portData.EquipmentName)
	}
	equipmentType, err := equipment.QueryType().Only(ctx)
	if err != nil {
		return errors.Wrapf(err, "fetching equipment type for equipment")
	}
	if equipmentType.Name != portData.EquipmentTypeName {
		return errors.Errorf("wrong equipment type. should be %q, but %q", equipmentType.Name, portData.EquipmentTypeName)
	}
	return nil
}

func (m *importer) parseImportArgs(r *http.Request) ([]int, *bool, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, nil, err
	}

	skipLines, err := getLinesToSkip(r)
	if err != nil {
		return nil, nil, err
	}

	verifyBeforeCommit, err := getVerifyBeforeCommitParam(r)
	if err != nil {
		return nil, nil, err
	}
	return skipLines, verifyBeforeCommit, nil
}

func isEmptyRow(s []string) bool {
	for _, v := range s {
		if v != "" {
			return false
		}
	}
	return true
}
