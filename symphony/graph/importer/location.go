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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// processLocationsCSV imports locations from CSV file to DB
// nolint: staticcheck
func (m *importer) processLocationsCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := m.logger.For(ctx)
	var (
		errs               Errors
		verifyBeforeCommit *bool
		commitRuns         []bool
		numRows, validRows int
	)
	nextLineToSkipIndex := -1

	log.Debug("Locations- started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
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

	verifyBeforeCommit, err = getVerifyBeforeCommitParam(r)
	if err != nil {
		errorReturn(w, "can't parse skipped lines", log, err)
		return
	}

	if pointer.GetBool(verifyBeforeCommit) {
		commitRuns = []bool{false, true}
	} else {
		commitRuns = []bool{true}
	}
	startSaving := false

	for fileName := range r.MultipartForm.File {
		firstLine, _, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusInternalServerError)
			return
		}
		m.populateIndexToLocationTypeMap(ctx, firstLine, true)
		latIndx := findIndex(firstLine, "latitude")
		longIndx := findIndex(firstLine, "longitude")
		externalIDIndex := findIndexForSimilar(firstLine, "external id")
		if getImportContext(ctx).lowestHierarchyIndex == -1 {
			log.Warn("location types on title does not match")
			errorReturn(w, fmt.Sprintf("location types on title does not match schema"), log, err)
			return
		}

		for _, commit := range commitRuns {
			// if we encounter errors on the "verifyBefore" flow - don't run the commit=true phase
			if commit && *verifyBeforeCommit && len(errs) != 0 {
				break
			} else if commit && len(errs) == 0 {
				startSaving = true
			}
			numRows, validRows = 0, 0
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
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "cannot read row"})
					continue
				}
				numRows++
				if shouldSkipLine(skipLines, numRows, nextLineToSkipIndex) {
					log.Warn("skipping line", zap.Error(err), zap.Int("line_number", numRows))
					nextLineToSkipIndex++
					continue
				}
				line := m.trimLine(untrimmedLine)
				lastPopulatedLocationIdx := getLowestLocationHierarchyIdxForRow(ctx, line)
				if lastPopulatedLocationIdx == -1 {
					log.Warn(fmt.Sprintf("invalid row (#%d). no location found", numRows))
					errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: "invalid row. no location found"})
					continue
				}
				locationForRow, err := m.handleLocationRow(ctx, lastPopulatedLocationIdx, latIndx, longIndx, externalIDIndex, line, commit)
				if err != nil {
					log.Warn(fmt.Sprintf("invalid row (#%d). handling row failed", numRows))
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: ""})
					continue
				}
				validRows++
				if locationForRow == nil {
					continue
				}

			}
		}
	}
	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, validRows, numRows, errs, !*verifyBeforeCommit || len(errs) == 0, startSaving)
	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
	log.Debug("Locations- Done", zap.Any("errors list", errs), zap.Int("all_lines", numRows), zap.Int("edited_added_rows", validRows))
}

func getProperties(ctx context.Context, line []string, index int) map[string]string {
	ic := getImportContext(ctx)
	var propsKeys []string
	if typeID, ok := ic.indexToLocationTypeID[index]; ok {
		if val, ok := ic.typeIDsToProperties[typeID]; ok {
			propsKeys = val
		}
	}
	propKeyValue := make(map[string]string)
	for _, propKey := range propsKeys {
		if idx, ok := ic.propNameToIndex[propKey]; ok {
			val := line[idx]
			propKeyValue[propKey] = val
		}
	}
	return propKeyValue
}

func (m *importer) handleLocationRow(ctx context.Context, lastPopulatedLocationIdx, latIndx, longIndx, externalIDIndex int, line []string, commit bool) (*ent.Location, error) {
	var (
		err      error
		parentID *int
		locID    int
		log      = m.logger.For(ctx)
	)
	indexToLocationTypeID := getImportContext(ctx).indexToLocationTypeID
	for index, name := range line {
		if index > lastPopulatedLocationIdx {
			break
		}
		var propertyInput []*models.PropertyInput
		var externalID string
		lat, long := 0.0, 0.0

		locationTypeID := indexToLocationTypeID[index]
		if index == lastPopulatedLocationIdx {
			propertyMap := getProperties(ctx, line, index)
			for key, value := range propertyMap {
				currVal := value
				ptype, err := m.getOrCreatePropTypeForLocation(ctx, locationTypeID, key)
				if err != nil {
					log.Warn("can't create property for location type", zap.Error(err))
					return nil, fmt.Errorf("can't create property (%v) for location type. id=%v. error: %v", key, locationTypeID, err.Error())
				}
				inp, err := getPropInput(*ptype, currVal)
				if err != nil {
					log.Warn("error evaluating property input", zap.Error(err), zap.String("property type", ptype.Name), zap.String("property value", currVal))
					return nil, fmt.Errorf("error evaluating property input. error: %v", err.Error())
				}
				propertyInput = append(propertyInput, inp)
			}
			if externalIDIndex != -1 {
				externalID = line[externalIDIndex]
			}
			if latIndx != -1 && line[latIndx] != "" && longIndx != -1 && line[longIndx] != "" {
				lat, err = strconv.ParseFloat(line[latIndx], 64)
				if err != nil {
					log.Warn("wrong latitude", zap.Error(err))
					return nil, fmt.Errorf("wrong latitude value: %v", line[latIndx])
				}
				long, err = strconv.ParseFloat(line[longIndx], 64)
				if err != nil {
					log.Warn("wrong longitude", zap.Error(err))
					return nil, fmt.Errorf("wrong longitude value: %v", line[longIndx])
				}
			}
		}
		client := m.ClientFrom(ctx)
		q := client.Location.Query().
			Where(location.HasTypeWith(locationtype.ID(locationTypeID))).
			Where(location.Name(name))
		if parentID != nil {
			q = q.Where(location.HasParentWith(location.ID(*parentID)))
		} else {
			q = q.Where(location.Not(location.HasParent()))
		}
		locID, err := q.FirstID(ctx)
		if ent.MaskNotFound(err) != nil {
			log.Warn("query location", zap.Error(err))
			return nil, fmt.Errorf("query location. name=%v. error: %v", name, err.Error())
		}
		// nolint: gocritic
		if locID == 0 {
			ltyp, err := client.LocationType.Query().Where(locationtype.ID(locationTypeID)).Only(ctx)
			if err != nil {
				return nil, fmt.Errorf("no valid location type on column number %d. error: %v", index, err.Error())
			}

			var l *ent.Location
			if commit {
				l, _, err = m.getOrCreateLocation(ctx, name, lat, long, ltyp, parentID, propertyInput, &externalID)
				if err != nil {
					return nil, fmt.Errorf("query/creating location. name=%v. error: %v", name, err.Error())
				}
			} else {
				l, err = m.queryLocationForTypeAndParent(ctx, name, ltyp, parentID)
				if l == nil && ent.MaskNotFound(err) == nil {
					// no location but no error (dry run mode)
					return nil, nil
				} else if err != nil {
					return nil, fmt.Errorf("query location. name=%v. error: %v", name, err.Error())
				}
			}
			locID = l.ID
		} else if index == lastPopulatedLocationIdx && (lat != 0 || long != 0 || len(propertyInput) > 0 || externalID != "") {
			for _, inp := range propertyInput {
				ptype := m.ClientFrom(ctx).PropertyType.Query().Where(propertytype.ID(inp.PropertyTypeID)).OnlyX(ctx)
				propertyID, err := ptype.QueryProperties().Where(property.HasLocationWith(location.ID(locID))).FirstID(ctx)
				if ent.MaskNotFound(err) != nil {
					log.Warn("can't find property for location", zap.Error(err))
					return nil, fmt.Errorf("can't find property for location. error: %v", err.Error())
				}
				if err == nil {
					inp.ID = &propertyID
				}
			}
			if commit {
				_, err := m.r.Mutation().EditLocation(ctx, models.EditLocationInput{
					ID: locID, Name: name, Latitude: lat, Longitude: long, Properties: propertyInput, ExternalID: &externalID,
				})
				if err != nil {
					log.Warn("couldn't edit existing location", zap.Error(err))
					return nil, fmt.Errorf("couldn't edit existing location. error: %v", err.Error())
				}
			}
		}
		parentID = &locID
	}
	loc, err := m.r.Query().Location(ctx, locID)
	if err != nil {
		return nil, err
	}
	return loc, nil
}
