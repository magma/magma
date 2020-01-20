// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint
package importer

import (
	"fmt"
	"io"
	"net/http"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// ProcessRuralTransportCSV imports rural sites transport equipment from CSV file to DB
func (m *importer) ProcessRuralTransportCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	log.Debug("ProcessRuralTransportCSV -started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusUnprocessableEntity)
		return
	}
	ctx = m.CloneContext(ctx)
	mr := m.r.Mutation()

	for fileName := range r.MultipartForm.File {
		_, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusUnprocessableEntity)
			return
		}

		energyEquipmentType := m.getOrCreateEquipmentType(ctx, "Energia", 0, "", 0, []*models.PropertyTypeInput{
			{
				Name: "Tipo Energia",
				Type: "string",
			},
			{
				Name: "Detalle Equipos Energia",
				Type: "string",
			},
		})
		energyTipEnergiaPTypeID := m.getEquipPropTypeID(ctx, "Tipo Energia", energyEquipmentType.ID)
		energyDetalleEquiposEnergiaPTypeID := m.getEquipPropTypeID(ctx, "Detalle Equipos Energia", energyEquipmentType.ID)

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

			nombreNodoVal := line[3]
			energyType := line[9]
			energyEquipmentDetails := line[10]
			transportType := line[11]
			transportState := line[12]
			transportNeUltimeMilla := line[13]
			transportBackhaulEbc := line[14]
			transportTecnologia := line[15]
			transportBanda := line[16]
			transportModem := line[17]
			transportHub := line[18]

			if transportState == "SIP" {
				continue
			}

			estacion, err := m.getLocationByName(ctx, nombreNodoVal)
			if estacion == nil {
				log.Warn(fmt.Sprintf("cannot find site %s on row %d", nombreNodoVal, rowID), zap.Error(err))
				continue
			}

			switch transportType {
			case "SAT":
				vsatTypeName := transportModem
				if len(vsatTypeName) < 2 {
					vsatTypeName = transportHub
				}
				if len(vsatTypeName) < 2 {
					continue
				}

				vsatEquipmentType := m.getOrCreateEquipmentType(ctx, vsatTypeName, 0, "", 0, []*models.PropertyTypeInput{
					{
						Name: "Technologia",
						Type: "string",
					},
					{
						Name: "Banda",
						Type: "string",
					},
					{
						Name: "Modem",
						Type: "string",
					},
					{
						Name: "Hub",
						Type: "string",
					},
				})

				vsatTechnologiaPTypeID := m.getEquipPropTypeID(ctx, "Technologia", vsatEquipmentType.ID)
				vsatBandaPTypeID := m.getEquipPropTypeID(ctx, "Banda", vsatEquipmentType.ID)
				vsatModemPTypeID := m.getEquipPropTypeID(ctx, "Modem", vsatEquipmentType.ID)
				vsatHubPTypeID := m.getEquipPropTypeID(ctx, "Hub", vsatEquipmentType.ID)
				vsatPropertyInput := []*models.PropertyInput{
					{
						PropertyTypeID:     vsatTechnologiaPTypeID,
						StringValue:        &transportTecnologia,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     vsatBandaPTypeID,
						StringValue:        &transportBanda,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     vsatModemPTypeID,
						StringValue:        &transportModem,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     vsatHubPTypeID,
						StringValue:        &transportHub,
						IsInstanceProperty: &instance,
					},
				}
				// create specific BTS Equipment types
				m.getOrCreateEquipment(ctx, mr, nombreNodoVal+"_SAT", vsatEquipmentType, nil, estacion, nil, vsatPropertyInput)
			case "FO":
				fallthrough
			case "MW":
				foTypeName := transportNeUltimeMilla
				if len(foTypeName) < 2 {
					foTypeName = transportBackhaulEbc
				}
				if len(foTypeName) < 2 {
					continue
				}

				foEquipmentType := m.getOrCreateEquipmentType(ctx, foTypeName, 0, "", 0, []*models.PropertyTypeInput{
					{
						Name: "Ne Ultima Milla",
						Type: "string",
					},
					{
						Name: "Backhaul EBC",
						Type: "string",
					},
				})

				foNeUltimeMillaPTypeID := m.getEquipPropTypeID(ctx, "Ne Ultima Milla", foEquipmentType.ID)
				foBackhaulEbcPTypeID := m.getEquipPropTypeID(ctx, "Backhaul EBC", foEquipmentType.ID)
				vsatPropertyInput := []*models.PropertyInput{
					{
						PropertyTypeID:     foNeUltimeMillaPTypeID,
						StringValue:        &transportNeUltimeMilla,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     foBackhaulEbcPTypeID,
						StringValue:        &transportBackhaulEbc,
						IsInstanceProperty: &instance,
					},
				}
				// create specific BTS Equipment types
				m.getOrCreateEquipment(ctx, mr, nombreNodoVal+"_"+transportType, foEquipmentType, nil, estacion, nil, vsatPropertyInput)
			}

			energyPropertyInput := []*models.PropertyInput{
				{
					PropertyTypeID:     energyTipEnergiaPTypeID,
					StringValue:        &energyType,
					IsInstanceProperty: &instance,
				},
				{
					PropertyTypeID:     energyDetalleEquiposEnergiaPTypeID,
					StringValue:        &energyEquipmentDetails,
					IsInstanceProperty: &instance,
				},
			}

			m.getOrCreateEquipment(ctx, mr, energyType+" Energia", energyEquipmentType, nil, estacion, nil, energyPropertyInput)
		}
		log.Debug("Done!!")
		w.WriteHeader(http.StatusOK)
	}
}
