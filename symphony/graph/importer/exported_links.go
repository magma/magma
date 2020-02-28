// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var fixedFirstPortLink = []string{"Link ID", "Port A Name", "Equipment A Name", "Equipment A Type"}
var fixedSecondPortLink = []string{"Port B Name", "Equipment B Name", "Equipment B Type"}

func minimalLinksLineLength() int {
	return len(fixedFirstPortLink) + len(fixedSecondPortLink) + 1 + maxEquipmentParents*2*2
}

// processExportedLinks imports links csv generated from the export feature
// nolint: staticcheck
func (m *importer) processExportedLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.logger.For(ctx)
	var (
		commitRuns             []bool
		errs                   Errors
		modifiedCount, numRows int
	)
	nextLineToSkipIndex := -1
	log.Debug("exported links-started")
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

	if *verifyBeforeCommit {
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
		importHeader, err := NewImportHeader(first, ImportEntityLink)
		if err != nil {
			errorReturn(w, "error on header", log, err)
			return
		}
		if err = m.inputValidationsLinks(ctx, importHeader); err != nil {
			errorReturn(w, "first line validation error", log, err)
			return
		}

		for _, commit := range commitRuns {
			// if we encounter errors on the "verifyBefore" flow - don't run the commit=true phase
			if commit && *verifyBeforeCommit && len(errs) != 0 {
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

				ln := m.trimLine(untrimmedLine)
				importLine := NewImportRecord(ln, importHeader)
				portARecord, portBRecord, err := m.getTwoPortRecords(importLine)

				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "getting two ports"})
					continue
				}
				id := importLine.ID()
				if id == 0 {
					client := m.ClientFrom(ctx)
					var linkPropertyInputs []*models.PropertyInput
					linkInput := make(map[int]*models.LinkSide, 2)

					for i, portRecord := range []ImportRecord{*portARecord, *portBRecord} {
						side, msg, propertyInputs, err := m.getLinkSide(ctx, client, portRecord, importLine, importHeader, commit)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: msg})
							break
						}
						if side == nil {
							if !commit {
								break
							}
							errs = append(errs, ErrorLine{Line: numRows, Error: "failed to create link side", Message: fmt.Sprintf("port: %v", portRecord.Name())})
							break
						}
						linkPropertyInputs = append(linkPropertyInputs, propertyInputs...)
						linkInput[i] = side
					}
					if linkInput[0] == nil || linkInput[1] == nil {
						continue
					}
					serviceIds, err := m.validateServicesForLinks(ctx, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "validating services where the link is a part of them"})
						continue
					}
					if commit {
						_, err = m.r.Mutation().AddLink(ctx, models.AddLinkInput{
							Sides: []*models.LinkSide{
								linkInput[0],
								linkInput[1],
							},
							Properties: linkPropertyInputs,
							ServiceIds: serviceIds,
						})
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "creating/fetching link"})
							continue
						}
						modifiedCount++
						log.Info(fmt.Sprintf("(row #%d) creating link", numRows))
					}
				} else {
					// edit existing link - only properties
					l, err := m.validateLineForExistingLink(ctx, id, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating existing port: id %v", id)})
						continue
					}
					serviceIds, err := m.validateServicesForLinks(ctx, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "validating services where the link is a part of them: id %v"})
						continue
					}

					allPropInputs, msg, err := m.getLinkPropertyInputs(ctx, importLine, importHeader, l)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: msg})
						continue
					}
					if len(allPropInputs) == 0 {
						modifiedCount++
						log.Info(fmt.Sprintf("(row #%d) [SKIPING]no port types or link properties", numRows), zap.String("name", importLine.Name()), zap.Int("id", importLine.ID()))
						continue
					}

					if commit {
						_, err = m.r.Mutation().EditLink(ctx, models.EditLinkInput{
							ID:         id,
							Properties: allPropInputs,
							ServiceIds: serviceIds,
						})
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("saving link: id %v", id)})
							continue
						}
						modifiedCount++
						log.Info(fmt.Sprintf("(row #%d) editing link", numRows), zap.String("name", importLine.Name()), zap.Int("id", importLine.ID()))
					}
				}
			}
		}
		log.Debug("Exported links - Done")
		w.WriteHeader(http.StatusOK)
		err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0, startSaving)
		if err != nil {
			errorReturn(w, "cannot marshal message", log, err)
			return
		}
	}
}

func (m *importer) getLinkPropertyInputs(ctx context.Context, importLine ImportRecord, importHeader ImportHeader, l *ent.Link) ([]*models.PropertyInput, string, error) {
	var allPropInputs []*models.PropertyInput
	ports, err := l.QueryPorts().All(ctx)
	if err != nil {
		return nil, fmt.Sprintf("querying link ports: id %v", l.ID), err
	}
	var portTypes []interface{}
	for _, port := range ports {
		definition := port.QueryDefinition().OnlyX(ctx)
		portType, _ := definition.QueryEquipmentPortType().Only(ctx)

		if portType != nil && importLine.Len() > importHeader.PropertyStartIdx() {
			portTypes = append(portTypes, portType)
			portProps, err := m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityLink)
			if err != nil {
				return nil, fmt.Sprintf("validating property for type %v.", portType.Name), err
			}
			allPropInputs = append(allPropInputs, portProps...)
		}
	}
	if len(portTypes) != 0 {
		err = importLine.validatePropertiesMismatch(ctx, portTypes)
	}
	if err != nil {
		return nil, "", err
	}
	return allPropInputs, "", nil
}

func (m *importer) getLinkSide(ctx context.Context, client *ent.Client, portRecord, linkRecord ImportRecord, linkHeader ImportHeader, commit bool) (*models.LinkSide, string, []*models.PropertyInput, error) {
	etn := portRecord.PortEquipmentTypeName()
	defName := portRecord.Name()
	en := portRecord.PortEquipmentName()

	equipmentType, err := client.EquipmentType.Query().Where(equipmenttype.Name(etn)).Only(ctx)
	if err != nil {
		return nil, fmt.Sprintf("getting equipment type: %v", etn), nil, err
	}
	portDef, err := equipmentType.QueryPortDefinitions().Where(equipmentportdefinition.Name(defName)).Only(ctx)
	if err != nil {
		return nil, fmt.Sprintf("getting port definition %v under equipment type %v", defName, etn), nil, err
	}

	parentLoc, err := m.verifyOrCreateLocationHierarchy(ctx, portRecord, commit, nil)

	if err != nil {
		return nil, "error while creating/verifying location hierarchy", nil, err
	} else if parentLoc == nil && !commit {
		return nil, "", nil, nil
	}

	parentEquipmentID, positionDefinitionID, err := m.getPositionDetailsIfExists(ctx, parentLoc, portRecord, false)
	if err != nil {
		return nil, "fetching equipment and positions hierarchy", nil, err
	}
	var pos *ent.EquipmentPosition

	if parentEquipmentID != nil && positionDefinitionID != nil {
		parentLoc = nil
		if commit {
			pos, err = resolverutil.GetOrCreatePosition(ctx, m.ClientFrom(ctx), parentEquipmentID, positionDefinitionID, false)
		} else {
			pos, err = resolverutil.ValidateAndGetPositionIfExists(ctx, client, parentEquipmentID, positionDefinitionID, false)
		}
		if err != nil {
			return nil, "creating equipment position", nil, err
		}
	}
	var equipment *ent.Equipment
	if commit {
		equipment, _, err = m.getOrCreateEquipment(ctx, m.r.Mutation(), en, equipmentType, nil, parentLoc, pos, nil)

	} else {
		equipment, err = m.getEquipmentIfExist(ctx, en, equipmentType, parentLoc, pos)
		if equipment == nil && err == nil {
			return nil, "", nil, nil
		}
	}
	if err != nil {
		return nil, "creating/fetching equipment", nil, err
	}
	var propInputs []*models.PropertyInput
	if linkRecord.Len() > linkHeader.PropertyStartIdx() {
		portType, err := portDef.QueryEquipmentPortType().Only(ctx)
		if ent.MaskNotFound(err) != nil {
			return nil, fmt.Sprintf("can't fetch port type %v", portDef.Name), nil, err
		}
		if portType != nil {
			propInputs, err = m.validatePropertiesForPortType(ctx, linkRecord, portType, ImportEntityLink)
			if err != nil {
				return nil, fmt.Sprintf("validating property for type %v", portType.Name), nil, err
			}
		}
	}
	return &models.LinkSide{Equipment: equipment.ID, Port: portDef.ID}, "", propInputs, nil
}

func (m *importer) getTwoPortRecords(importLine ImportRecord) (*ImportRecord, *ImportRecord, error) {
	header := importLine.Header()
	headerSlices := header.LinkGetTwoPortsSlices()
	ahead, bhead := headerSlices[0], headerSlices[1]
	headerA, err := NewImportHeader(ahead, ImportEntityPortInLink)
	if err != nil {
		return nil, nil, err
	}
	headerB, err := NewImportHeader(bhead, ImportEntityPortInLink)
	if err != nil {
		return nil, nil, err
	}
	portsSlices := importLine.LinkGetTwoPortsSlices()
	portASlice, portBSlice := portsSlices[0], portsSlices[1]
	if equal(portASlice, portBSlice) {
		return nil, nil, errors.New("ports are identical")
	}

	portA := NewImportRecord(portASlice, headerA)
	portB := NewImportRecord(portBSlice, headerB)
	return &portA, &portB, nil
}

func (m *importer) validateLineForExistingLink(ctx context.Context, linkID int, importLine ImportRecord) (*ent.Link, error) {
	link, err := m.ClientFrom(ctx).Link.Query().Where(link.ID(linkID)).Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching link")
	}
	ports, err := link.QueryPorts().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching link ports")
	}
	if len(ports) != 2 {
		return nil, errors.New("link must have two ports")
	}
	portAFromFile, portBFromFile, err := m.getTwoPortRecords(importLine)
	if err != nil {
		return nil, err
	}

	var linkPropNames []string
	for _, port := range ports {
		def, err := port.QueryDefinition().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't fetch port definition")
		}
		equip, err := port.QueryParent().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't fetch port equipment parent")
		}
		equipType, err := equip.QueryType().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't fetch equipment type")
		}

		if !(def.Name == portAFromFile.Name() &&
			equip.Name == portAFromFile.PortEquipmentName() &&
			equipType.Name == portAFromFile.PortEquipmentTypeName()) && !(def.Name == portBFromFile.Name() &&
			equip.Name == portBFromFile.PortEquipmentName() &&
			equipType.Name == portBFromFile.PortEquipmentTypeName()) {
			return nil, errors.Errorf("port doesn't match line: %v, %v, %v", def.Name, equip.Name, equipType.Name)
		}
		// TODO Validate location and position (currently not editing it, therefor not validating)

		portType, err := def.QueryEquipmentPortType().Only(ctx)
		if ent.MaskNotFound(err) != nil {
			return nil, errors.Wrapf(err, "fetching equipment port type")
		}
		if portType != nil {
			lps, err := portType.QueryLinkPropertyTypes().All(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "fetching links port type properties")
			}
			for _, value := range lps {
				linkPropNames = append(linkPropNames, value.Name)
			}
		}
	}
	for propTypName, value := range importLine.PropertiesMap() {
		if value != "" {
			if findIndex(linkPropNames, propTypName) == -1 {
				return nil, errors.Errorf("link property %v does not exist on either portType", propTypName)
			}
		}
	}
	return link, nil
}

func (m *importer) inputValidationsLinks(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	if len(firstLine) < minimalLinksLineLength() {
		return errors.Errorf("first line too short. should include: %q and location/position data  for both sides", fixedFirstPortLink)
	}
	if firstLine[0] != "Link ID" {
		return errors.Errorf("first cell should be 'Link ID' ")
	}
	portsSlices := importHeader.LinkGetTwoPortsSlices()
	ha, err := NewImportHeader(portsSlices[0], ImportEntityPortInLink)
	if err != nil {
		return err
	}
	locStart, _ := ha.LocationsRangeIdx()
	if !equal(ha.line[:locStart], []string{"Port A Name", "Equipment A Name", "Equipment A Type"}) {
		return errors.New("first line misses sequence; 'Port A Name', 'Equipment A Name' or 'Equipment A Type' ")
	}
	err = m.validateAllLocationTypeExist(ctx, locStart, ha.LocationTypesRangeArr(), false)
	if err != nil {
		return err
	}
	hb, err := NewImportHeader(portsSlices[1], ImportEntityPortInLink)
	if err != nil {
		return err
	}
	locStart, _ = hb.LocationsRangeIdx()
	if !equal(hb.line[:locStart], []string{"Port B Name", "Equipment B Name", "Equipment B Type"}) {
		return errors.New("first line misses sequence; 'Port B Name', 'Equipment B Name' or 'Equipment B Type' ")
	}
	err = m.validateAllLocationTypeExist(ctx, locStart, hb.LocationTypesRangeArr(), false)
	if err != nil {
		return err
	}
	if !equal(ha.line[ha.prnt3Idx:importHeader.LinkSecondPortStartIdx()-1], []string{"Parent Equipment (3) A", "Position (3) A", "Parent Equipment (2) A", "Position (2) A", "Parent Equipment A", "Equipment Position A"}) {
		return errors.New("First port on first line misses sequence: 'Parent Equipment (3) A', 'Position (3) A', 'Parent Equipment (2) A', 'Position (2) A', 'Parent Equipment A' or 'Equipment Position A'")
	}
	if !equal(hb.line[hb.prnt3Idx:], []string{"Parent Equipment (3) B", "Position (3) B", "Parent Equipment (2) B", "Position (2) B", "Parent Equipment B", "Equipment Position B"}) {
		return errors.New("second port on first line misses sequence: 'Parent Equipment (3) B', 'Position (3) B', 'Parent Equipment (2) B', 'Position (2) B', 'Parent Equipment B' or 'Equipment Position B'")
	}
	if importHeader.ServiceNamesIdx() == -1 {
		return errors.New("column 'Service Names' is missing")
	}
	return nil
}

func (m *importer) validateServicesForLinks(ctx context.Context, line ImportRecord) ([]int, error) {
	serviceNamesMap := make(map[string]bool)
	var serviceIds []int
	serviceNames := strings.Split(line.ServiceNames(), ";")
	for _, serviceName := range serviceNames {
		if serviceName != "" {
			serviceID, err := m.validateServiceExistsAndUnique(ctx, serviceNamesMap, serviceName)
			if err != nil {
				return nil, err
			}
			serviceIds = append(serviceIds, serviceID)
		}
	}
	return serviceIds, nil
}
