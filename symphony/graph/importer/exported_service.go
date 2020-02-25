// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"

	"io"
	"net/http"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const minimalLineLength = 6

// processExportedService imports service csv generated from the export feature
// nolint: staticcheck, dupl
func (m *importer) processExportedService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.logger.For(ctx)
	client := m.ClientFrom(ctx)
	var (
		err                    error
		commitRuns             []bool
		errs                   Errors
		modifiedCount, numRows int
	)

	nextLineToSkipIndex := -1
	log.Debug("Exported Service - started")
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
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file: %q. file name: %q", err, fileName), http.StatusInternalServerError)
			return
		}
		importHeader, err := NewImportHeader(first, ImportEntityService)
		if err != nil {
			errorReturn(w, "error on header", log, err)
			return
		}

		firstLine := importHeader.line
		if len(firstLine) < minimalLineLength {
			err := errors.New("first line too short. should include: 'Service ID', 'Service Name', 'Service Type', 'Service External ID', 'Customer Name' and 'Customer External ID'")
			log.Warn("first line validation error", zap.Error(err))
			http.Error(w, fmt.Sprintf("first line validation error: %q", err), http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Warn("data fetching error", zap.Error(err))
			http.Error(w, fmt.Sprintf("data fetching error: %s", err.Error()), http.StatusInternalServerError)
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
				serviceTypName := importLine.TypeName()
				serviceType, err := client.ServiceType.Query().Where(servicetype.Name(serviceTypName)).Only(ctx)
				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("couldn't find service type %v", serviceTypName)})
					continue
				}

				var customerID *string = nil
				customerName := importLine.CustomerName()
				if customerName != "" {
					var customer *ent.Customer
					if commit {
						customer, err = m.getOrCreateCustomer(ctx, m.r.Mutation(), customerName, importLine.CustomerExternalID())
					} else {
						customer, err = m.getCustomerIfExist(ctx, customerName)
					}
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("add customer with name %v", importLine.CustomerName())})
						continue
					}
					if customer != nil {
						customerID = &customer.ID
					}
				}

				externalID := pointer.ToStringOrNil(importLine.ServiceExternalID())

				status, err := m.getValidatedStatus(importLine)
				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("failed parsing status with value %v", importLine.Status())})
					continue
				}

				id := importLine.ID()
				var propInputs []*models.PropertyInput
				if importLine.Len() > importHeader.PropertyStartIdx() {
					propInputs, err = m.validatePropertiesForServiceType(ctx, importLine, serviceType)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating property for type %v", serviceType.Name)})
						continue
					}
				}
				if id == "" {
					var (
						created bool
						service *ent.Service
					)

					if commit {
						_, created, err = m.getOrCreateService(ctx, m.r.Mutation(), name, serviceType, propInputs, customerID, externalID, *status)
						if err == nil {
							if created {
								modifiedCount++
								log.Info(fmt.Sprintf("(row #%d) creating service", numRows), zap.String("name", name))
							} else {
								errs = append(errs, ErrorLine{Line: numRows, Error: "service exists", Message: fmt.Sprintf("service %v already exists under location/position (id=%v)", service.Name, service.ID)})
								continue
							}
						}
					} else {
						service, err = m.getServiceIfExist(ctx, m.r.Mutation(), name, serviceType, propInputs, customerID, externalID, *status)
						if service != nil {
							err = errors.Errorf("service %v already exists", name)
						}
					}
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/fetching service"})
						continue
					}
				} else {
					// existingService
					service, err := m.validateLineForExistingService(ctx, id, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating existing service: id %v", id)})
						continue
					}

					propertiesValid := false
					for i, propInput := range propInputs {
						propID, err := service.QueryProperties().Where(property.HasTypeWith(propertytype.ID(propInput.PropertyTypeID))).OnlyID(ctx)
						if err != nil {
							if !ent.IsNotFound(err) {
								errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("property fetching error: property type id %v", propInput.PropertyTypeID)})
								break
							}
						} else {
							propInput.ID = &propID
						}

						if i == len(propInputs)-1 {
							propertiesValid = true
						}
					}
					if !propertiesValid {
						continue
					}
					if commit {
						_, err = m.r.Mutation().EditService(ctx, models.ServiceEditData{
							ID:         id,
							Name:       &name,
							Properties: propInputs,
							ExternalID: externalID,
							CustomerID: customerID,
							Status:     status,
						})
						modifiedCount++
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("editing service: id %v", id)})
							continue
						}
					}
				}
			}
		}
	}
	log.Debug("Exported Service - Done")
	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0, startSaving)

	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
}

func (m *importer) validateLineForExistingService(ctx context.Context, serviceID string, importLine ImportRecord) (*ent.Service, error) {
	service, err := m.r.Query().Service(ctx, serviceID)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching service")
	}
	typ := service.QueryType().OnlyX(ctx)
	if typ.Name != importLine.TypeName() {
		return nil, errors.Errorf("wrong service type. should be %v, but %v", typ.Name, importLine.TypeName())
	}
	return service, nil
}

func (m *importer) getValidatedStatus(importLine ImportRecord) (*models.ServiceStatus, error) {
	statuses := make([]string, len(models.AllServiceStatus))
	for i, status := range models.AllServiceStatus {
		statuses[i] = status.String()
	}

	index := findIndexForSimilar(statuses, importLine.Status())
	if index == -1 {
		return nil, errors.Errorf("failed parse status %q", importLine.Status())
	}
	return &models.AllServiceStatus[index], nil
}

func (m *importer) validatePropertiesForServiceType(ctx context.Context, line ImportRecord, serviceType *ent.ServiceType) ([]*models.PropertyInput, error) {
	var pInputs []*models.PropertyInput
	propTypes, err := serviceType.QueryPropertyTypes().All(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, errors.Wrap(err, "can't query property types for service type")
	}
	for _, ptype := range propTypes {
		ptypeName := ptype.Name
		pInput, err := line.GetPropertyInput(m.ClientFrom(ctx), ctx, serviceType, ptypeName)
		if err != nil {
			return nil, err
		}
		if pInput != nil {
			pInputs = append(pInputs, pInput)
		}
	}
	return pInputs, nil
}
