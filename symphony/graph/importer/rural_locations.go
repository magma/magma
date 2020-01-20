// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint
package importer

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// ProcessRuralLegacyLocationsCSV imports locations data from CSV file to DB
func (m *importer) ProcessRuralLocationsCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	log.Debug("ProcessRuralRanLocationsCSV -started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusUnprocessableEntity)
		return
	}
	ctx = m.CloneContext(ctx)

	for fileName := range r.MultipartForm.File {
		_, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusUnprocessableEntity)
			return
		}

		DepartamentoLocationType := m.getOrCreateLocationType(ctx, "Departamento", nil)
		m.updateMapTypeForLocationType(ctx, DepartamentoLocationType.ID, "map", 8)

		ProvinciaLocationType := m.getOrCreateLocationType(ctx, "Provincia", nil)
		m.updateMapTypeForLocationType(ctx, ProvinciaLocationType.ID, "map", 10)

		DistritoLocationType := m.getOrCreateLocationType(ctx, "Distrito", nil)
		m.updateMapTypeForLocationType(ctx, DistritoLocationType.ID, "map", 12)

		CentroPobladoLocationType := m.getOrCreateLocationType(ctx, "Centro Poblado", nil)
		m.updateMapTypeForLocationType(ctx, CentroPobladoLocationType.ID, "map", 13)

		EstacionLocationType := m.getOrCreateLocationType(ctx, "Estacion", []*models.PropertyTypeInput{
			{
				Name: "Codigo Unico Estacion",
				Type: "string",
			},
			{
				Name: "Direccion",
				Type: "string",
			},
		})

		m.updateMapTypeForLocationType(ctx, EstacionLocationType.ID, "map", 14)

		cuePTypeID := m.getLocPropTypeID(ctx, "Codigo Unico Estacion", EstacionLocationType.ID)
		dPTypeID := m.getLocPropTypeID(ctx, "Direccion", EstacionLocationType.ID)

		instance := true
		rowID := 0
		for {
			rowID++
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn(fmt.Sprintf("cannot read row %d", rowID), zap.Error(err))
				continue
			}
			isNonEmptyRow := false
			for _, value := range line {
				if len(value) != 0 {
					isNonEmptyRow = true
					break
				}
			}
			if !isNonEmptyRow {
				continue
			}

			CodigoUnicoEStacionVal := line[0]
			NombreEstacionVal := line[14]
			DepartamentoVal := line[3]
			ProvinciaVal := line[4]
			DistritoVal := line[5]
			NombreCentroPobladoVal := line[6]
			DireccionVal := line[7]
			latVal := line[8]
			longVal := line[9]

			parsedLat, _ := strconv.ParseFloat(latVal, 64)
			parsedLong, _ := strconv.ParseFloat(longVal, 64)

			Departamento, _ := m.getOrCreateLocation(ctx, DepartamentoVal, 0.0, 0.0, DepartamentoLocationType, nil, nil, nil)
			Provincia, _ := m.getOrCreateLocation(ctx, ProvinciaVal, 0.0, 0.0, ProvinciaLocationType, &Departamento.ID, nil, nil)
			Distrito, _ := m.getOrCreateLocation(ctx, DistritoVal, 0.0, 0.0, DistritoLocationType, &Provincia.ID, nil, nil)
			CentroPoblado, _ := m.getOrCreateLocation(ctx, NombreCentroPobladoVal, 0.0, 0.0, CentroPobladoLocationType, &Distrito.ID, nil, nil)

			estPropertyInput := []*models.PropertyInput{
				{
					PropertyTypeID:     cuePTypeID,
					StringValue:        &CodigoUnicoEStacionVal,
					IsInstanceProperty: &instance,
				}, {
					PropertyTypeID:     dPTypeID,
					StringValue:        &DireccionVal,
					IsInstanceProperty: &instance,
				},
			}
			l, wasNew := m.getOrCreateLocation(ctx, NombreEstacionVal, parsedLat, parsedLong, EstacionLocationType, &CentroPoblado.ID, estPropertyInput, &CodigoUnicoEStacionVal)
			if wasNew {
				oldLocations, err := EstacionLocationType.QueryLocations().Where(location.HasPropertiesWith(property.HasTypeWith(propertytype.Name("Codigo Unico Estacion")), property.StringVal(CodigoUnicoEStacionVal)), location.IDNEQ(l.ID)).All(ctx)

				if err != nil {
					m.log.For(ctx).Error("Failed query location")
					continue
				}

				for _, oldLocation := range oldLocations {
					m.log.For(ctx).Debug("Duplicate location found", zap.String("ID", oldLocation.ID), zap.String("CodigoUnico", CodigoUnicoEStacionVal))
				}
			}

		}
		log.Debug("Done!!")
		w.WriteHeader(http.StatusOK)
	}
}
