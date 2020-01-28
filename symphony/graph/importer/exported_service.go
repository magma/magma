// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"

	"io"
	"net/http"

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
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	nextLineToSkipIndex := -1
	log.Debug("Exported Service - started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}
	count, numRows := 0, 0
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

	for fileName := range r.MultipartForm.File {
		first, reader, err := m.newReader(fileName, r)
		importHeader := NewImportHeader(first, ImportEntityService)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file: %q. file name: %q", err, fileName), http.StatusInternalServerError)
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
				log.Warn("couldn't find service type", zap.Error(err), zap.String("service_type", serviceTypName))
				http.Error(w, fmt.Sprintf("couldn't find service type %q (row #%d). %q ", serviceTypName, numRows, err), http.StatusBadRequest)
				return
			}

			var customerID *string = nil
			customerName := importLine.CustomerName()
			if customerName != "" {
				customer, err := m.getOrCreateCustomer(ctx, m.r.Mutation(), customerName, importLine.CustomerExternalID())
				if err != nil {
					log.Error("add customer", zap.String("name", importLine.CustomerName()), zap.Error(err))
					http.Error(w, fmt.Sprintf("add customer with name %q (row #%d). %q", importLine.CustomerName(), numRows, err.Error()), http.StatusBadRequest)
					return
				}
				if customer != nil {
					customerID = &customer.ID
				}
			}

			externalID := pointer.ToStringOrNil(importLine.ServiceExternalID())

			status, err := m.getValidatedStatus(importLine)
			if err != nil {
				errorReturn(w, fmt.Sprintf("failed parsing status with value %q (row #%d)", importLine.Status(), numRows), log, nil)
				return
			}

			id := importLine.ID()
			var propInputs []*models.PropertyInput
			if importLine.Len() > importHeader.PropertyStartIdx() {
				propInputs, err = m.validatePropertiesForServiceType(ctx, importLine, serviceType)
				if err != nil {
					log.Warn("validating property for type", zap.Error(err))
					http.Error(w, fmt.Sprintf("validating property for type %q (row #%d). %q", serviceType.Name, numRows, err.Error()), http.StatusBadRequest)
					return
				}
			}
			if id == "" {
				service, created := m.getOrCreateService(ctx, m.r.Mutation(), name, serviceType, propInputs, customerID, externalID, *status)
				if created {
					count++
					log.Warn(fmt.Sprintf("(row #%d) creating service", numRows), zap.String("name", service.Name), zap.String("id", service.ID))
				} else {
					errorReturn(w, fmt.Sprintf("(row #%d) Service %v already exists under location/position (id=%v)", numRows, service.Name, service.ID), log, nil)
					return
				}
			} else {
				// existingService
				service, err := m.validateLineForExistingService(ctx, id, importLine)
				if err != nil {
					log.Warn("validating existing service", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("%q: validating existing service: id %q (row #%d)", err, id, numRows), http.StatusBadRequest)
					return
				}
				for _, propInput := range propInputs {
					propID, err := service.QueryProperties().Where(property.HasTypeWith(propertytype.ID(propInput.PropertyTypeID))).OnlyID(ctx)
					if err != nil {
						if !ent.IsNotFound(err) {
							log.Warn("property fetching error", zap.Error(err), importLine.ZapField())
							http.Error(w, fmt.Sprintf("%q: property fetching error: property type id %q (row #%d)", err, propInput.PropertyTypeID, numRows), http.StatusBadRequest)
							return
						}
					} else {
						propInput.ID = &propID
					}
				}
				_, err = m.r.Mutation().EditService(ctx, models.ServiceEditData{
					ID:         id,
					Name:       &name,
					Properties: propInputs,
					ExternalID: externalID,
					CustomerID: customerID,
					Status:     status,
				})
				if err != nil {
					log.Warn("editing service", zap.Error(err), importLine.ZapField())
					http.Error(w, fmt.Sprintf("editing service: id %q (row #%d). %q: ", id, numRows, err), http.StatusBadRequest)
					return
				}
			}
		}
	}
	log.Debug("Exported Service - Done")
	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, count, numRows, nil, true)

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
		pInputs = append(pInputs, pInput)
	}
	return pInputs, nil
}
