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

	log.Debug("exported links-started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	count, numRows := 0, 0

	for fileName := range r.MultipartForm.File {
		first, reader, err := m.newReader(fileName, r)
		importHeader := NewImportHeader(first, ImportEntityLink)
		if err != nil {
			errorReturn(w, fmt.Sprintf("cannot handle file: %q", fileName), log, err)
			return
		}

		if err = m.inputValidationsLinks(ctx, importHeader); err != nil {
			errorReturn(w, "first line validation error", log, err)
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
			importLine := NewImportRecord(m.trimLine(untrimmedLine), importHeader)

			id := importLine.ID()
			if id == "" {
				errorReturn(w, fmt.Sprintf("supporting only link property editing (row #%d)", numRows), log, err)
				return
			} else {
				//edit existing link - only properties
				link, err := m.validateLineForExistingLink(ctx, id, importLine)
				if err != nil {
					errorReturn(w, fmt.Sprintf("validating existing port: id %q (row #%d)", id, numRows), log, err)
					return
				}
				var allPropInputs []*models.PropertyInput
				ports, err := link.QueryPorts().All(ctx)
				if err != nil {
					errorReturn(w, fmt.Sprintf("querying link ports: id %q (row #%d)", id, numRows), log, err)
					return
				}
				for _, port := range ports {
					definition := port.QueryDefinition().OnlyX(ctx)
					portType, _ := definition.QueryEquipmentPortType().Only(ctx)
					if portType != nil && importLine.Len() > importHeader.PropertyStartIdx() {
						portProps, err := m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityLink)
						if err != nil {
							errorReturn(w, fmt.Sprintf("validating property for type %q (row #%d).", portType.Name, numRows), log, err)
							return
						}
						allPropInputs = append(allPropInputs, portProps...)
					}
				}
				if len(allPropInputs) == 0 {
					log.Info(fmt.Sprintf("(row #%d) [SKIPING]no port types or link properties", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
					continue
				}
				_, err = m.r.Mutation().EditLink(ctx, models.EditLinkInput{
					ID:         id,
					Properties: allPropInputs,
				})
				if err != nil {
					errorReturn(w, fmt.Sprintf("saving link: id %q (row #%d)", id, numRows), log, err)
					return
				}
				count++
				log.Info(fmt.Sprintf("(row #%d) editing link", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
			}
		}
	}
	log.Debug("Exported links - Done")
	w.WriteHeader(http.StatusOK)

	err := writeSuccessMessage(w, count, numRows)
	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
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
	header := importLine.Header()
	portsSlices := importLine.LinkGetTwoPortsSlices()
	portASlice, portBSlice := portsSlices[0], portsSlices[1]
	if equal(portASlice, portBSlice) {
		return nil, errors.New("ports are identical")
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

		if !(def.Name == portASlice[header.NameIdx()] &&
			equip.Name == portASlice[header.PortEquipmentNameIdx()] &&
			equipType.Name == portASlice[header.PortEquipmentTypeNameIdx()]) && !(def.Name == portBSlice[header.NameIdx()] &&
			equip.Name == portBSlice[header.PortEquipmentNameIdx()] &&
			equipType.Name == portBSlice[header.PortEquipmentTypeNameIdx()]) {
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
	return nil
}
