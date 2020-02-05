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

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const minimalPortsLineLength = 15

// processExportedPorts imports ports csv generated from the export feature
// nolint: staticcheck
func (m *importer) processExportedPorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)

	nextLineToSkipIndex := -1
	log.Debug("exported ports-started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	var (
		commitRuns             []bool
		errs                   Errors
		modifiedCount, numRows int
	)

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

		for _, commit := range commitRuns {
			// if we encounter errors on the "verifyBefore" flow - don't run the commit=true phase
			if commit && *verifyBeforeCommit && len(errs) != 0 {
				break
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

				id := importLine.ID()
				if id == "" {
					errs = append(errs, ErrorLine{Line: numRows, Error: "no id provided for row", Message: "supporting only port editing"})
					continue
				} else {
					//edit existing  port
					port, err := m.validateLineForExistingPort(ctx, id, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%v: validating existing port: id %v", err.Error(), id)})
						continue
					}
					consumerServiceIds, providerServiceIds, err := m.validateServicesForPortEndpoints(ctx, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%v: validating services where the port is to be endpoint for the services: id %v", err.Error(), id)})
						continue
					}
					if commit {
						err = m.editServiceEndpoints(ctx, port, consumerServiceIds, models.ServiceEndpointRoleConsumer)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%v: editing services where the port is to be consumer endpoint for the services: id %v", err.Error(), id)})
							continue
						}
						err = m.editServiceEndpoints(ctx, port, providerServiceIds, models.ServiceEndpointRoleProvider)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%q: editing services where the port is to be provider endpoint for the services: id %v", err.Error(), id)})
							continue
						}
					}
					var propInputs []*models.PropertyInput
					parent := port.QueryParent().OnlyX(ctx)
					definition := port.QueryDefinition().OnlyX(ctx)
					portType, _ := definition.QueryEquipmentPortType().Only(ctx)
					if portType != nil && importLine.Len() > importHeader.PropertyStartIdx() {
						propInputs, err = m.validatePropertiesForPortType(ctx, importLine, portType, ImportEntityPort)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%v: validating property for type %v", err.Error(), portType.Name)})
							continue
						}
						if commit {
							_, err = m.r.Mutation().EditEquipmentPort(ctx, models.EditEquipmentPortInput{
								Side: &models.LinkSide{
									Equipment: parent.ID,
									Port:      definition.ID,
								},
								Properties: propInputs,
							})

							if err != nil {
								errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("%v: saving port: id %v", err.Error(), id)})
								continue
							}
							modifiedCount++
							log.Info(fmt.Sprintf("(row #%d) editing port", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
						}
					} else {
						log.Info(fmt.Sprintf("(row #%d) [SKIPING]no port type or properties", numRows), zap.String("name", importLine.Name()), zap.String("id", importLine.ID()))
					}
				}
			}
		}
	}
	log.Debug("Exported ports - Done")
	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0)

	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
}

func (m *importer) validateLineForExistingPort(ctx context.Context, portID string, importLine ImportRecord) (*ent.EquipmentPort, error) {
	port, err := m.ClientFrom(ctx).EquipmentPort.Query().Where(equipmentport.ID(portID)).Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment port")
	}
	portData, err := importLine.PortData(nil)
	if err != nil {
		return nil, errors.New("error while calculating port data")
	}
	err = m.validatePort(ctx, *portData, *port)

	if err != nil {
		return nil, err
	}
	equipment, err := port.QueryParent().Only(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching equipment for port")
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
	if !equal(firstLine[prnt3Idx:importHeader.PropertyStartIdx()], []string{"Parent Equipment (3)", "Parent Equipment (2)", "Parent Equipment", "Equipment Position", "Linked Port ID", "Linked Port Name", "Linked Equipment ID", "Linked Equipment", "Consumer Endpoint for These Services", "Provider Endpoint for These Services"}) {
		return errors.New("first line should include: 'Parent Equipment (3)', 'Parent Equipment (2)', 'Parent Equipment', 'Equipment Position' 'Linked Port ID', 'Linked Port Name', 'Linked Equipment ID', 'Linked Equipment', 'Consumer Endpoint for These Services', 'Provider Endpoint for These Services'")
	}
	err := m.validateAllLocationTypeExist(ctx, 5, importHeader.LocationTypesRangeArr(), false)
	return err
}

func (m *importer) validateServicesForPortEndpoints(ctx context.Context, line ImportRecord) ([]string, []string, error) {
	serviceNamesMap := make(map[string]bool)
	var consumerServiceIds []string
	var providerServiceIds []string
	consumerServiceNames := strings.Split(line.ConsumerPortsServices(), ";")
	for _, serviceName := range consumerServiceNames {
		if serviceName != "" {
			serviceID, err := m.validateServiceExistsAndUnique(ctx, serviceNamesMap, serviceName)
			if err != nil {
				return nil, nil, err
			}
			consumerServiceIds = append(consumerServiceIds, serviceID)
		}

	}
	providerServiceNames := strings.Split(line.ProviderPortsServices(), ";")
	for _, serviceName := range providerServiceNames {
		if serviceName != "" {
			serviceID, err := m.validateServiceExistsAndUnique(ctx, serviceNamesMap, serviceName)
			if err != nil {
				return nil, nil, err
			}
			providerServiceIds = append(providerServiceIds, serviceID)
		}
	}
	return consumerServiceIds, providerServiceIds, nil
}

func (m *importer) editServiceEndpoints(ctx context.Context, port *ent.EquipmentPort, serviceIds []string, role models.ServiceEndpointRole) error {
	mutation := m.r.Mutation()
	currentServiceIds, err := port.QueryEndpoints().Where(serviceendpoint.Role(role.String())).QueryService().IDs(ctx)
	if err != nil {
		return err
	}
	addedServiceIds, deletedServiceIds := resolverutil.GetDifferenceBetweenSlices(currentServiceIds, serviceIds)
	for _, serviceID := range addedServiceIds {
		if _, err := mutation.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
			ID:     serviceID,
			PortID: port.ID,
			Role:   role,
		}); err != nil {
			return err
		}
	}
	for _, serviceID := range deletedServiceIds {
		serviceEndpointID, err := port.QueryEndpoints().Where(serviceendpoint.HasServiceWith(service.ID(serviceID))).OnlyID(ctx)
		if err != nil {
			return err
		}
		if _, err := mutation.RemoveServiceEndpoint(ctx, serviceEndpointID); err != nil {
			return err
		}
	}
	return nil
}
