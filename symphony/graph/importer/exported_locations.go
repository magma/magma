// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

const minimalLocationLineLength = 4

// processExportedLocation imports location csv generated from the export feature
func (m *importer) processExportedLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.logger.For(ctx)
	nextLineToSkipIndex := -1
	client := m.ClientFrom(ctx)

	log.Debug("Exported location - started")
	var (
		err                    error
		modifiedCount, numRows int
		errs                   Errors
		commitRuns             []bool
	)
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
			errorReturn(w, fmt.Sprintf("cannot handle file: %q", fileName), log, err)
			return
		}
		importHeader, err := NewImportHeader(first, ImportEntityLocation)
		if err != nil {
			errorReturn(w, "error on header", log, err)
			return
		}

		//	populating:
		//	indexToLocationTypeID
		if err = m.inputValidationsLocation(ctx, importHeader); err != nil {
			errorReturn(w, "first line validation error", log, err)
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

				var parentLoc *ent.Location
				currLocIndex, err := m.getCurrentLocationIndex(importLine)
				if err != nil || currLocIndex <= 0 {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "getting relevant location index"})
					continue
				}
				parentIndex, err := m.getParentOfLocationIndex(importLine)
				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "getting relevant location index"})
					continue
				}

				shouldHaveParent := parentIndex != -1
				var parentLocID *int
				if shouldHaveParent {
					parentLoc, err = m.verifyOrCreateLocationHierarchy(ctx, importLine, commit, pointer.ToInt(currLocIndex-1))
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "getting location path and parents"})
						continue
					}
					if parentLoc != nil {
						parentLocID = &parentLoc.ID
					}
				}
				var externalID string
				if importLine.ExternalID() != "" {
					externalID = importLine.ExternalID()
				}
				locName := importLine.line[currLocIndex]
				lat, long, msg, err := m.getLatLong(importLine)
				if err != nil {
					errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: msg})
					continue
				}

				id := importLine.ID()
				if id == 0 {
					// new location
					typName := importHeader.line[currLocIndex] // the actual index
					locType, err := client.LocationType.Query().Where(locationtype.Name(typName)).Only(ctx)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("querying location type %v", typName)})
						continue
					}

					var propInputs []*models.PropertyInput
					if importLine.Len() > importHeader.PropertyStartIdx() {
						propInputs, err = m.validatePropertiesForLocationType(ctx, importLine, locType)
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("validating property for type %q", locType.Name)})
							continue
						}
					}

					if shouldHaveParent && parentLoc == nil {
						if commit {
							errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: fmt.Sprintf("failed getting parent location %v", importLine.line[parentIndex])})
							log.Info("Row " + strconv.FormatInt(int64(numRows), 10) + ": failed getting parent location")
							continue
						} else {
							modifiedCount++
							log.Info("Row " + strconv.FormatInt(int64(numRows), 10) + ": skipping row, parent locations does not exist")
							// can't continue checks
							continue
						}
					}
					var created bool
					var loc *ent.Location
					if commit {
						_, created, err = m.getOrCreateLocation(ctx, locName, lat, long, locType, parentLocID, propInputs, &externalID)
						if created {
							modifiedCount++
						} else if err == nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: "location already exists"})
							log.Info("Row " + strconv.FormatInt(int64(numRows), 10) + ": location already exists under location/position")
							continue
						}
					} else { // no commit
						loc, err = m.queryLocationForTypeAndParent(ctx, locName, locType, parentLocID)
						err = ent.MaskNotFound(err)
						if loc != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: "", Message: "location already exists"})
							continue
						}
						modifiedCount++
					}
					if err != nil { // both commit and dry-run
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error while creating/fetching location"})
						continue
					}
				} else {
					// existing location
					location, err := m.validateLineForExistingLocation(ctx, id, importLine)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: "error validating location line"})
						continue
					}
					inputs, msg, err := m.getLocationPropertyInputs(ctx, importLine, location)
					if err != nil {
						errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: msg})
						continue
					}
					if commit {
						_, err = m.r.Mutation().EditLocation(ctx, models.EditLocationInput{
							ID:         id,
							Name:       locName,
							Properties: inputs,
							ExternalID: &externalID,
							Latitude:   lat,
							Longitude:  long,
						})
						if err != nil {
							errs = append(errs, ErrorLine{Line: numRows, Error: err.Error(), Message: fmt.Sprintf("editing location: id %v", id)})
							continue
						}
					}
					modifiedCount++
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	err = writeSuccessMessage(w, modifiedCount, numRows, errs, !*verifyBeforeCommit || len(errs) == 0, startSaving)
	if err != nil {
		errorReturn(w, "cannot marshal message", log, err)
		return
	}
	log.Debug("Exported location - Done", zap.Any("errors list", errs), zap.Int("all_lines", numRows), zap.Int("edited_added_rows", modifiedCount))
}

func (m *importer) getCurrentLocationIndex(importLine ImportRecord) (int, error) {
	header := importLine.Header()
	i := header.ExternalIDIdx() - 1
	locIndexStart, _ := header.LocationsRangeIdx()

	for ; i >= locIndexStart; i-- {
		if importLine.line[i] != "" {
			return i, nil
		}
	}
	return -1, fmt.Errorf("no location names specified in row")
}

func (m *importer) getParentOfLocationIndex(importLine ImportRecord) (int, error) {
	header := importLine.Header()
	currLocationIndex, err := m.getCurrentLocationIndex(importLine)
	if err != nil {
		return -1, err
	}
	locIndexStart, _ := header.LocationsRangeIdx()
	i := currLocationIndex - 1
	for ; i >= locIndexStart; i-- {
		if importLine.line[i] != "" {
			return i, nil
		}
	}
	// no parent
	return -1, nil
}

func (m *importer) getLocationPropertyInputs(ctx context.Context, importLine ImportRecord, location *ent.Location) ([]*models.PropertyInput, string, error) {
	typ := location.QueryType().OnlyX(ctx)
	propTypes, err := typ.QueryPropertyTypes().All(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, fmt.Sprintf("can't query property types for location type %v", typ.Name), err
	}

	var inputs []*models.PropertyInput
	for _, propType := range propTypes {
		propName := propType.Name
		inp, err := importLine.GetPropertyInput(m.ClientFrom(ctx), ctx, typ, propName)
		if inp == nil {
			continue
		}
		propType := typ.QueryPropertyTypes().Where(propertytype.Name(propName)).OnlyX(ctx)
		if err != nil {
			return nil, fmt.Sprintf("getting property input: prop %v", propName), err
		}
		propID, err := location.QueryProperties().Where(property.HasTypeWith(propertytype.ID(propType.ID))).OnlyID(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, fmt.Sprintf("property fetching error: property name %v", propName), err
			}
		} else {
			inp.ID = &propID
		}
		inputs = append(inputs, inp)
	}
	return inputs, "", nil
}

func (m *importer) validateLineForExistingLocation(ctx context.Context, locationID int, importLine ImportRecord) (*ent.Location, error) {
	location, err := m.ClientFrom(ctx).Location.Get(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("fetching location: %w", err)
	}
	typ, err := location.QueryType().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching location type: %w", err)
	}
	currLocIndex, err := m.getCurrentLocationIndex(importLine)
	if err != nil {
		return nil, err
	}
	typNameFromHeader := importLine.title.line[currLocIndex]
	if typ.Name != typNameFromHeader {
		return nil, fmt.Errorf("wrong location type. should be %v, but %v", typ.Name, typNameFromHeader)
	}
	err = m.validateHierarchyForLocation(ctx, location, importLine)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (m *importer) validateHierarchyForLocation(ctx context.Context, location *ent.Location, importLine ImportRecord) error {
	hierarchy, err := m.r.Location().LocationHierarchy(ctx, location)
	if err != nil {
		return err
	}
	prevIdx := 0
	for _, loc := range hierarchy {
		currIdx := findIndex(importLine.line, strings.Trim(loc.Name, " "))
		if currIdx == -1 {
			return fmt.Errorf("missing location from hierarchy (%v)", loc.Name)
		}
		if prevIdx > currIdx {
			return fmt.Errorf("location not in the right order (%v)", loc.Name)
		}
		prevIdx = currIdx
	}
	return nil
}

func (m *importer) inputValidationsLocation(ctx context.Context, importHeader ImportHeader) error {
	firstLine := importHeader.line
	if len(firstLine) < minimalLocationLineLength {
		return errors.New("first line too short. should include: 'Location ID', all location types, 'External ID', 'Latitude' and 'Longitude'")
	}
	locStart, _ := importHeader.LocationsRangeIdx()
	if !equal(firstLine[:locStart], []string{"Location ID"}) {
		return errors.New("first line should begin with 'Location ID'")
	}
	if !equal(firstLine[importHeader.ExternalIDIdx():importHeader.PropertyStartIdx()], []string{"External ID", "Latitude", "Longitude"}) {
		return errors.New("first line misses sequence: 'External ID', 'Latitude' and 'Longitude'")
	}
	err := m.validateAllLocationTypeExist(ctx, 1, importHeader.LocationTypesRangeArr(), false)
	return err
}

func (m *importer) validatePropertiesForLocationType(ctx context.Context, line ImportRecord, locType *ent.LocationType) ([]*models.PropertyInput, error) {
	var pInputs []*models.PropertyInput
	propTypes, err := locType.QueryPropertyTypes().All(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, fmt.Errorf("can't query property types for location type: %w", err)
	}
	for _, propType := range propTypes {
		ptypeName := propType.Name
		pInput, err := line.GetPropertyInput(m.ClientFrom(ctx), ctx, locType, ptypeName)
		if err != nil {
			return nil, err
		}
		if pInput != nil {
			pInputs = append(pInputs, pInput)
		}
	}
	return pInputs, nil
}

func (m *importer) getLatLong(importLine ImportRecord) (float64, float64, string, error) {
	var lat, long float64
	var err error
	if importLine.Latitude() != "" || importLine.Longitude() != "" {
		lat, err = strconv.ParseFloat(importLine.Latitude(), 64)
		if err != nil {
			return lat, long, "wrong latitude", err
		}

		long, err = strconv.ParseFloat(importLine.Longitude(), 64)
		if err != nil {
			return lat, long, "wrong longitude", err
		}
	}
	return lat, long, "", nil
}
