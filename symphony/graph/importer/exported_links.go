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
	log := m.log.For(ctx)
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
	err := r.ParseForm()
	if err != nil {
		errorReturn(w, "can't parse form", log, err)
		return
	}
	skipLines, err := getLinesToSkip(r)
	if err != nil {
		errorReturn(w, "can't parse skipped lines", log, err)
		return
	}
	if len(skipLines) > 0 {
		nextLineToSkipIndex = 0
	}

	verifyBeforeCommit, err := getVerifyBeforeCommitParam(r)
	if err != nil {
		errorReturn(w, "can't parse verify_before_commit param", log, err)
		return
	}

	if *verifyBeforeCommit {
		commitRuns = []bool{false, true}
	} else {
		commitRuns = []bool{true}
	}

	for fileName := range r.MultipartForm.File {
		first, _, err := m.newReader(fileName, r)
		importHeader := NewImportHeader(first, ImportEntityLink)
		if err != nil {
			errorReturn(w, fmt.Sprintf("cannot handle file: %q", fileName), log, err)
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
				if id == "" {
					client := m.ClientFrom(ctx)
					var linkPropertyInputs []*models.PropertyInput
					linkInput := make(map[int]*models.LinkSide, 2)

					for i, portRecord := range []ImportRecord{*portARecord, *portBRecord} {
						etn := portRecord.PortEquipmentTypeName()
						defName := portRecord.Name()
						en := portRecord.PortEquipmentName()

						equipmentType, err := client.EquipmentType.Query().Where(equipmenttype.Name(etn)).Only(ctx)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("getting equipment type: %v", etn)})
							continue
						}
						portDef, err := equipmentType.QueryPortDefinitions().Where(equipmentportdefinition.Name(defName)).Only(ctx)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("getting port definition %v under equipment type %v", defName, etn)})
							continue
						}

						parentLoc, err := m.verifyOrCreateLocationHierarchy(ctx, portRecord, commit)

						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/verifying location hierarchy"})
							continue
						} else if parentLoc == nil && !commit {
							continue
						}

						parentEquipmentID, positionDefinitionID, err := m.getPositionDetailsIfExists(ctx, parentLoc, portRecord, false)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "fetching equipment and positions hierarchy"})
							continue
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
								errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "creating equipment position"})
								continue
							}
						}
						var equipment *ent.Equipment
						if commit {
							equipment, _, err = m.getOrCreateEquipment(ctx, m.r.Mutation(), en, equipmentType, nil, parentLoc, pos, nil)

						} else {
							equipment, err = m.getEquipmentIfExist(ctx, m.r.Mutation(), en, equipmentType, nil, parentLoc, pos, nil)
							if equipment == nil && err == nil {
								continue
							}
						}
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "creating/fetching equipment"})
							continue
						}
						var propInputs []*models.PropertyInput
						if importLine.Len() > importHeader.PropertyStartIdx() {
							portType, err := portDef.QueryEquipmentPortType().Only(ctx)
							if ent.MaskNotFound(err) != nil {
								errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("can't fetch port type %v", portDef.Name)})
								continue
							}
							if portType != nil {
								propInputs, err = m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityLink)
								if err != nil {
									errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating property for type %v", portType.Name)})
									continue
								}
								linkPropertyInputs = append(linkPropertyInputs, propInputs...)
							}
						}
						linkInput[i] = &models.LinkSide{Equipment: equipment.ID, Port: portDef.ID}
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
					//edit existing link - only properties
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
					var allPropInputs []*models.PropertyInput
					ports, err := l.QueryPorts().All(ctx)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("querying link ports: id %v", id)})
						continue
					}
					for _, port := range ports {
						definition := port.QueryDefinition().OnlyX(ctx)
						portType, _ := definition.QueryEquipmentPortType().Only(ctx)
						if portType != nil && importLine.Len() > importHeader.PropertyStartIdx() {
							portProps, err := m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityLink)
							if err != nil {
								errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating property for type %v.", portType.Name)})
								continue
							}
							allPropInputs = append(allPropInputs, portProps...)
						}
					}
					if len(allPropInputs) == 0 {
						log.Info(fmt.Sprintf("(row #%d) [SKIPING]no port types or link properties", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
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
						log.Info(fmt.Sprintf("(row #%d) editing link", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
					}
				}
			}
		}
		log.Debug("Exported links - Done")
		w.WriteHeader(http.StatusOK)
		err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0)
		if err != nil {
			errorReturn(w, "cannot marshal message", log, err)
			return
		}
	}
}

func (m *importer) getTwoPortRecords(importLine ImportRecord) (*ImportRecord, *ImportRecord, error) {
	header := importLine.Header()
	headerSlices := header.LinkGetTwoPortsSlices()
	ahead, bhead := headerSlices[0], headerSlices[1]
	headerA := NewImportHeader(ahead, ImportEntityPortInLink)
	headerB := NewImportHeader(bhead, ImportEntityPortInLink)

	portsSlices := importLine.LinkGetTwoPortsSlices()
	portASlice, portBSlice := portsSlices[0], portsSlices[1]
	if equal(portASlice, portBSlice) {
		return nil, nil, errors.New("ports are identical")
	}

	portA := NewImportRecord(portASlice, headerA)
	portB := NewImportRecord(portBSlice, headerB)
	return &portA, &portB, nil
}

func (m *importer) validateLineForExistingLink(ctx context.Context, linkID string, importLine ImportRecord) (*ent.Link, error) {
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
	ha := NewImportHeader(portsSlices[0], ImportEntityPortInLink)

	locStart, _ := ha.LocationsRangeIdx()
	if !equal(ha.line[:locStart], []string{"Port A Name", "Equipment A Name", "Equipment A Type"}) {
		return errors.New("first line misses sequence; 'Port A Name', 'Equipment A Name' or 'Equipment A Type' ")
	}
	err := m.validateAllLocationTypeExist(ctx, locStart, ha.LocationTypesRangeArr(), false)
	if err != nil {
		return err
	}
	hb := NewImportHeader(portsSlices[1], ImportEntityPortInLink)
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

func (m *importer) validateServicesForLinks(ctx context.Context, line ImportRecord) ([]string, error) {
	serviceNamesMap := make(map[string]bool)
	var serviceIds []string
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
