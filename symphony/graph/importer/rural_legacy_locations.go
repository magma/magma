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

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// ProcessRuralLegacyLocationsCSV imports locations data from CSV file to DB
func (m *importer) ProcessRuralLegacyLocationsCSV(w http.ResponseWriter, r *http.Request) {
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
			{
				Name: "Nombre Centro Poblado",
				Type: "string",
			},
			{
				Name: "Nombre Estacion",
				Type: "string",
			},
		})

		m.updateMapTypeForLocationType(ctx, EstacionLocationType.ID, "map", 14)

		cuePTypeID := m.getLocPropTypeID(ctx, "Codigo Unico Estacion", EstacionLocationType.ID)
		dPTypeID := m.getLocPropTypeID(ctx, "Direccion", EstacionLocationType.ID)
		ncpPTypeID := m.getLocPropTypeID(ctx, "Nombre Centro Poblado", EstacionLocationType.ID)
		nePTypeID := m.getLocPropTypeID(ctx, "Nombre Estacion", EstacionLocationType.ID)

		instance := true
		rowID := 0

		estacions := make([]*ent.Location, 0)
		pobaldos := make([]*ent.Location, 0)

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

			CodigoUnicoEStacionVal := line[3]
			NombreEstacionVal := line[4]
			DepartamentoVal := line[6]
			ProvinciaVal := line[7]
			DistritoVal := line[8]
			NombreCentroPobladoVal := line[9]
			DireccionVal := line[10]
			latVal := line[11]
			longVal := line[12]

			parsedLat, _ := strconv.ParseFloat(latVal, 64)
			parsedLong, _ := strconv.ParseFloat(longVal, 64)

			Departamento, _ := m.getOrCreateLocation(ctx, DepartamentoVal, 0.0, 0.0, DepartamentoLocationType, nil, nil, nil)
			Provincia, _ := m.getOrCreateLocation(ctx, ProvinciaVal, 0.0, 0.0, ProvinciaLocationType, &Departamento.ID, nil, nil)
			Distrito, _ := m.getOrCreateLocation(ctx, DistritoVal, 0.0, 0.0, DistritoLocationType, &Provincia.ID, nil, nil)

			estPropertyInput := []*models.PropertyInput{
				{
					PropertyTypeID:     cuePTypeID,
					StringValue:        &CodigoUnicoEStacionVal,
					IsInstanceProperty: &instance,
				}, {
					PropertyTypeID:     dPTypeID,
					StringValue:        &DireccionVal,
					IsInstanceProperty: &instance,
				}, {
					PropertyTypeID:     ncpPTypeID,
					StringValue:        &NombreCentroPobladoVal,
					IsInstanceProperty: &instance,
				},
				{
					PropertyTypeID:     nePTypeID,
					StringValue:        &NombreEstacionVal,
					IsInstanceProperty: &instance,
				},
			}
			Estacion, _ := m.getOrCreateLocation(ctx, NombreEstacionVal, parsedLat, parsedLong, EstacionLocationType, &Distrito.ID, estPropertyInput, &CodigoUnicoEStacionVal)
			estacions = append(estacions, Estacion)

			CentroPoblado, _ := m.getOrCreateLocation(ctx, NombreCentroPobladoVal, 0.0, 0.0, CentroPobladoLocationType, &Distrito.ID, nil, nil)
			pobaldos = append(pobaldos, CentroPoblado)
		}

		for i, Estacion := range estacions {
			if err := m.ClientFrom(ctx).Location.UpdateOne(Estacion).ClearParent().SetParent(pobaldos[i]).Exec(ctx); err != nil {
				log.Error("updating location", zap.Error(err))
			}
		}
		log.Debug("Done!!")
		w.WriteHeader(http.StatusOK)
	}
}
