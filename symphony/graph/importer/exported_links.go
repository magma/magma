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
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var fixedFirstLine = []string{"Link ID", "Port A ID", "Port A Name", "Port A Type", "Equipment A ID", "Equipment A Name", "Equipment A Type", "Port B ID", "Port B Name", "Port B Type", "Equipment B ID", "Equipment B Name", "Equipment B Type", "Service Names"}

func minimalLinksLineLength() int {
	return len(fixedFirstLine)
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
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file: %q. file name: %q", err, fileName), http.StatusInternalServerError)
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
				//edit existing link
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
				log.Info(fmt.Sprintf("(row #%d) editing port", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
			}
		}
	}
	log.Debug("Exported links - Done")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Edited %q instances, out of %q", strconv.FormatInt(int64(count), 10), strconv.FormatInt(int64(numRows), 10))
	w.Write([]byte(msg))
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
	portAData, err := importLine.PortData(pointer.ToString("A"))
	if err != nil {
		return nil, errors.New("error while calculating port A data")
	}
	portBData, err := importLine.PortData(pointer.ToString("B"))
	if err != nil {
		return nil, errors.New("error while calculating port B data")
	}

	if portAData.ID == portBData.ID {
		return nil, errors.New("same port for Port A and port B")
	}
	for _, port := range ports {
		switch port.ID {
		case portAData.ID:
			err = m.validatePort(ctx, *portAData, *port)
		case portBData.ID:
			err = m.validatePort(ctx, *portBData, *port)
		default:
			return nil, errors.Errorf("missing port %v on file for link %v", port.ID, linkID)
		}
		if err != nil {
			return nil, err
		}
		def, err := port.QueryDefinition().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "fetching equipment port definition")
		}
		portType, err := def.QueryEquipmentPortType().Only(ctx)
		if ent.MaskNotFound(err) != nil {
			return nil, errors.Wrapf(err, "fetching equipment port type")
		}
		if portType != nil {
			for propTypName, value := range importLine.PropertiesMap() {
				if value != "" {
					switch exist, err := portType.QueryLinkPropertyTypes().Where(propertytype.Name(propTypName)).Exist(ctx); {
					case err != nil:
						return nil, errors.Wrapf(err, "querying link properties for link %v", linkID)
					case !exist:
						return nil, errors.Errorf("link property %v does not exist on portType %v", propTypName, portType.Name)
					}
				}
			}
		}
	}
	return link, nil
}

func (m *importer) inputValidationsLinks(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	if len(firstLine) < minimalLinksLineLength() {
		return errors.Errorf("first line too short. should include: %q", fixedFirstLine)
	}
	propStart := importHeader.PropertyStartIdx()
	if !equal(firstLine[:propStart], fixedFirstLine) {
		return errors.Errorf("first line misses sequence: %q ", fixedFirstLine)
	}
	return nil
}
