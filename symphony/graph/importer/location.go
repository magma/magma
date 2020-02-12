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
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	log.Debug("Locations- started")
	instance := true
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusInternalServerError)
		return
	}

	for fileName := range r.MultipartForm.File {
		firstLine, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusInternalServerError)
			return
		}
		fullLocationPath := findIndex(firstLine, "Parent Location Name") == -1
		m.populateIndexToLocationTypeMap(ctx, firstLine, true)
		latIndx := findIndex(firstLine, "latitude")
		longIndx := findIndex(firstLine, "longitude")
		externalIDIndex := findIndexForSimilar(firstLine, "external id")
		indexToLocationTypeID := getImportContext(ctx).indexToLocationTypeID
		i := 0
		for {
			untrimmedLine, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				http.Error(w, fmt.Sprintf("cannot read row #%d", i), http.StatusInternalServerError)
				return
			}
			i++
			line := m.trimLine(untrimmedLine)
			if fullLocationPath && getImportContext(ctx).lowestHierarchyIndex == -1 {
				log.Warn("no location types on title")
				http.Error(w, "no location types on title", http.StatusInternalServerError)
				return
			}
			if fullLocationPath {
				lastPopulatedLocationIdx := getLowestLocationHierarchyIdxForRow(ctx, line)
				if lastPopulatedLocationIdx == -1 {
					log.Warn(fmt.Sprintf("invalid row (#%d). no location found", i))
					http.Error(w, fmt.Sprintf("invalid row (%d). no location found", i), http.StatusInternalServerError)
					return
				}
				var parentID *string
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
								http.Error(w, fmt.Sprintf("can't create property (%q) for location type. id=%q. row=%d", key, locationTypeID, i), http.StatusInternalServerError)
								return
							}
							inp, err := getPropInput(*ptype, currVal)
							if err != nil {
								log.Warn("error evaluating property input. Skipping", zap.Error(err), zap.String("property type", ptype.Name), zap.String("property value", currVal))
							}
							propertyInput = append(propertyInput, inp)
						}
						if externalIDIndex != -1 {
							externalID = line[externalIDIndex]
						}
						if latIndx != -1 && longIndx != -1 {
							lat, err = strconv.ParseFloat(line[latIndx], 64)
							if err != nil {
								log.Warn("no or wrong latitude", zap.Error(err))
								lat = 0.0
							}
							long, err = strconv.ParseFloat(line[longIndx], 64)
							if err != nil {
								log.Warn("no or wrong longitude", zap.Error(err))
								long = 0.0
							}
						}
					}

					q := client.Location.Query().
						Where(location.HasTypeWith(locationtype.ID(locationTypeID))).
						Where(location.Name(name))
					if parentID != nil {
						q = q.Where(location.HasParentWith(location.ID(*parentID)))
					} else {
						q = q.Where(location.Not(location.HasParent()))
					}
					id, err := q.FirstID(ctx)
					if ent.MaskNotFound(err) != nil {
						log.Warn("query location", zap.Error(err))
						http.Error(w, fmt.Sprintf("query location. name=%q. row=%d", name, i), http.StatusInternalServerError)
						return
					}
					if id == "" {
						ltyp := client.LocationType.Query().Where(locationtype.ID(locationTypeID)).OnlyX(ctx)
						l, _ := m.getOrCreateLocation(ctx, name, lat, long, ltyp, parentID, propertyInput, &externalID)
						id = l.ID
					} else if index == lastPopulatedLocationIdx && (lat != 0 || long != 0 || len(propertyInput) > 0 || externalID != "") {
						for _, inp := range propertyInput {
							ptype := m.ClientFrom(ctx).PropertyType.Query().Where(propertytype.ID(inp.PropertyTypeID)).OnlyX(ctx)
							propertyID, err := ptype.QueryProperties().Where(property.HasLocationWith(location.ID(id))).FirstID(ctx)
							if ent.MaskNotFound(err) != nil {
								log.Warn("can't find property for location", zap.Error(err))
								http.Error(w, fmt.Sprintf("can't find property (%q) for location. id=%q. row=%d", ptype.Name, id, i), http.StatusInternalServerError)
								return
							}
							if err == nil {
								inp.ID = &propertyID
							}
						}
						_, err := m.r.Mutation().EditLocation(ctx, models.EditLocationInput{
							ID: id, Name: name, Latitude: lat, Longitude: long, Properties: propertyInput, ExternalID: &externalID,
						})
						if err != nil {
							log.Warn("couldn't edit existing location", zap.Error(err))
							http.Error(w, fmt.Sprintf("couldn't edit existing location. id=%q. row=%d", id, i), http.StatusInternalServerError)
							return
						}
					}
					parentID = &id
				}
			} else { // if not full path
				typename := line[findIndex(firstLine, "Location Type")]
				locTyp, err := client.LocationType.Query().
					Where(locationtype.Name(typename)).
					Only(ctx)
				if err != nil {
					log.Warn(fmt.Sprintf("[SKIP]could not fetch location type. %s", typename), zap.Error(err))
					http.Error(w, fmt.Sprintf("[SKIP]could not fetch location type. %s, row %d", typename, i), http.StatusInternalServerError)
					return
				}
				propTypes := locTyp.QueryPropertyTypes().Where(propertytype.IsInstanceProperty(true)).AllX(ctx)
				parentName := line[findIndex(firstLine, "Parent Location Name")]
				parent, err := client.LocationType.Query().
					Where(locationtype.Not(locationtype.ID(locTyp.ID))).
					QueryLocations().
					Where(location.Name(parentName)).
					Only(ctx)
				if err != nil {
					log.Warn(fmt.Sprintf("[SKIP]could not fetch parent location. %s", parentName), zap.Error(err))
					http.Error(w, fmt.Sprintf("[SKIP]could not fetch parent location. %s, row %d", parentName, i), http.StatusInternalServerError)
					return
				}
				locName := line[findIndex(firstLine, "Location ID")]
				var propertyInput []*models.PropertyInput
				for _, ptype := range propTypes {
					currType := ptype
					idx := findIndex(firstLine, ptype.Name)
					if idx == -1 {
						continue
					}
					propertyInput = append(propertyInput, &models.PropertyInput{
						PropertyTypeID:     currType.ID,
						StringValue:        &line[idx],
						IsInstanceProperty: &instance,
					})
				}
				var lat, long float64
				if latIndx != -1 && longIndx != -1 {
					lat, err = strconv.ParseFloat(line[latIndx], 64)
					if err != nil {
						log.Warn("no or wrong latitude", zap.Error(err))
						lat = 0
					}
					long, err = strconv.ParseFloat(line[longIndx], 64)
					if err != nil {
						log.Warn("no or wrong longitude", zap.Error(err))
						long = 0
					}
				}
				var externalID string
				if externalIDIndex != -1 {
					externalID = line[externalIDIndex]
				}
				m.getOrCreateLocation(ctx, locName, lat, long, locTyp, &parent.ID, propertyInput, &externalID)
			}
		}
		err = writeSuccessMessage(w, i, i, []ErrorLine{}, true)
		if err != nil {
			errorReturn(w, "cannot marshal message", log, err)
			return
		}
	}
	log.Debug("Locations- Done")
	w.WriteHeader(http.StatusOK)
}

func getProperties(ctx context.Context, line []string, index int) map[string]string {
	ic := getImportContext(ctx)
	typeID := ic.indexToLocationTypeID[index]
	propsKeys := ic.typeIDsToProperties[typeID]
	propKeyValue := make(map[string]string)
	for _, propKey := range propsKeys {
		val := ""
		idx := ic.propNameToIndex[propKey]
		if idx != 0 {
			val = line[idx]
			propKeyValue[propKey] = val
		}
	}
	return propKeyValue
}
