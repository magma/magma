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

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const minimalPortsLineLength = 13

// processExportedPorts imports ports csv generated from the export feature
// nolint: staticcheck
func (m *importer) processExportedPorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)

	log.Debug("exported ports-started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	count, numRows := 0, 0

	for fileName := range r.MultipartForm.File {
		first, reader, err := m.newReader(fileName, r)
		importHeader := NewImportHeader(first, ImportEntityPort)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file: %q. file name: %q", err, fileName), http.StatusInternalServerError)
			return
		}
		//
		//	populating, but not using:
		//	indexToLocationTypeID
		//
		if err = m.inputValidationsPorts(ctx, importHeader); err != nil {
			log.Warn("first line validation error", zap.Error(err))
			http.Error(w, fmt.Sprintf("first line validation error: %q", err), http.StatusBadRequest)
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
				log.Warn("supporting only port editing", zap.Error(err))
				http.Error(w, fmt.Sprintf("supporting only port editing (row #%d)", numRows), http.StatusBadRequest)
				return
			} else {
				//edit existing  port
				port, err := m.validateLineForExistingPort(ctx, id, importLine)
				if err != nil {
					log.Warn("validating existing port", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("%q: validating existing port: id %q (row #%d)", err, id, numRows), http.StatusBadRequest)
					return
				}
				var propInputs []*models.PropertyInput
				parent := port.QueryParent().OnlyX(ctx)
				definition := port.QueryDefinition().OnlyX(ctx)
				portType, _ := definition.QueryEquipmentPortType().Only(ctx)
				if portType != nil && importLine.Len() > importHeader.PropertyStartIdx() {
					propInputs, err = m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityPort)
					if err != nil {
						log.Warn("validating property for type", zap.Error(err))
						http.Error(w, fmt.Sprintf("validating property for type %q (row #%d). %q", portType.Name, numRows, err.Error()), http.StatusBadRequest)
						return
					}

					_, err = m.r.Mutation().EditEquipmentPort(ctx, models.EditEquipmentPortInput{
						Side: &models.LinkSide{
							Equipment: parent.ID,
							Port:      definition.ID,
						},
						Properties: propInputs,
					})
					if err != nil {
						log.Warn("saving port", zap.Error(err), importLine.ZapField())
						http.Error(w, fmt.Sprintf("%q: saving port: id %q (row #%d)", err, id, numRows), http.StatusBadRequest)
						return
					}
					count++
					log.Info(fmt.Sprintf("(row #%d) editing port", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
				} else {
					log.Info(fmt.Sprintf("(row #%d) [SKIPING]no port type or properties", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
				}
			}
		}
	}
	log.Debug("Exported ports - Done")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Edited %q instances, out of %q", strconv.FormatInt(int64(count), 10), strconv.FormatInt(int64(numRows), 10))
	w.Write([]byte(msg))
}

func (m *importer) validateLineForExistingPort(ctx context.Context, portID string, importLine ImportRecord) (*ent.EquipmentPort, error) {
	port, err := m.ClientFrom(ctx).EquipmentPort.Query().Where(equipmentport.ID(portID)).Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment port")
	}
	def, err := port.QueryDefinition().Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment port definition")
	}
	if def.Name != importLine.Name() {
		return nil, errors.Wrapf(err, "wrong port type. should be %q, but %q", importLine.TypeName(), def.Name)
	}
	portType, err := def.QueryEquipmentPortType().Only(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, errors.Wrapf(err, "fetching equipment port type")
	}
	var tempPortType string
	if ent.IsNotFound(err) {
		tempPortType = ""
	} else {
		tempPortType = portType.Name
	}
	if tempPortType != importLine.TypeName() {
		return nil, errors.Wrapf(err, "wrong port type. should be %q, but %q", importLine.TypeName(), tempPortType)
	}

	equipment, err := port.QueryParent().Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment for port")
	}
	if equipment.Name != importLine.PortEquipmentName() {
		return nil, errors.Wrapf(err, "wrong equipment. should be %q, but %q", importLine.PortEquipmentName(), equipment.Name)
	}
	equipmentType, err := equipment.QueryType().Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment type for equipment")
	}
	if equipmentType.Name != importLine.PortEquipmentTypeName() {
		return nil, errors.Wrapf(err, "wrong equipment type. should be %q, but %q", importLine.PortEquipmentTypeName(), equipmentType.Name)
	}
	err = m.verifyPositionHierarchy(ctx, equipment, importLine)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching positions hierarchy")
	}
	err = m.validateLocationHierarchy(ctx, equipment, importLine)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func (m *importer) inputValidationsPorts(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	prnt3Idx := importHeader.prnt3Idx
	if len(firstLine) < minimalPortsLineLength {
		return errors.New("first line too short. should include: 'Port ID','Port Name','Port Type','Equipment Name','Equipment Type', location types, parents and link data")
	}
	locStart, _ := importHeader.LocationsRangeIdx()
	if !equal(firstLine[:locStart], []string{"Port ID", "Port Name", "Port Type", "Equipment Name", "Equipment Type"}) {
		return errors.New("first line misses sequence; 'Port ID','Port Name','Port Type','Equipment Name' or 'Equipment Type'")
	}
	if !equal(firstLine[prnt3Idx:importHeader.PropertyStartIdx()], []string{"Parent Equipment (3)", "Parent Equipment (2)", "Parent Equipment", "Equipment Position", "Linked Port ID", "Linked Port Name", "Linked Equipment ID", "Linked Equipment"}) {
		return errors.New("first line should include: 'Parent Equipment (3)', 'Parent Equipment (2)', 'Parent Equipment', 'Equipment Position' 'Linked Port ID', 'Linked Port Name', 'Linked Equipment ID', 'Linked Equipment'")
	}
	err := m.validateAllLocationTypeExist(ctx, 5, importHeader.LocationTypesRangeArr(), false)
	return err
}
