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

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// ProcessRuralRanCSV imports RAN equipment from CSV file to DB
func (m *importer) ProcessRuralRanCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	log.Debug("ProcessRuralRanCSV -started")
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

		DepartamentoLocationType := m.getOrCreateLocationType(ctx, "Departamento", nil)
		m.updateMapTypeForLocationType(ctx, DepartamentoLocationType.ID, "map", 8)

		ProvinciaLocationType := m.getOrCreateLocationType(ctx, "Provincia", nil)
		m.updateMapTypeForLocationType(ctx, ProvinciaLocationType.ID, "map", 10)

		DistritoLocationType := m.getOrCreateLocationType(ctx, "Distrito", nil)
		m.updateMapTypeForLocationType(ctx, DistritoLocationType.ID, "map", 12)

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

		AntennaEquipmentType := m.getOrCreateEquipmentType(ctx, "Antenna", 0, "", 1, []*models.PropertyTypeInput{
			{
				Name: "Numero",
				Type: "string",
			},
			{
				Name: "Altura",
				Type: "string",
			},
			{
				Name: "Marca",
				Type: "string",
			},
			{
				Name: "Nombre Modelo",
				Type: "string",
			},
			{
				Name: "Azimuth",
				Type: "int",
			},
			{
				Name: "Amplitud Beam",
				Type: "string",
			},
			{
				Name: "Ganancia",
				Type: "string",
			},
			{
				Name: "Tilt Electrico",
				Type: "string",
			},
			{
				Name: "Tilt Mecanico",
				Type: "string",
			},
			{
				Name: "Modelo Jumper",
				Type: "string",
			},
			{
				Name: "Longitud Jumper",
				Type: "string",
			},
			{
				Name: "Perdida Conector",
				Type: "string",
			},
		})

		antNumeroPTypeID := m.getEquipPropTypeID(ctx, "Numero", AntennaEquipmentType.ID)
		antAlturaPTypeID := m.getEquipPropTypeID(ctx, "Altura", AntennaEquipmentType.ID)
		antMarcaPTypeID := m.getEquipPropTypeID(ctx, "Marca", AntennaEquipmentType.ID)
		antNombreModelPTypeID := m.getEquipPropTypeID(ctx, "Nombre Modelo", AntennaEquipmentType.ID)
		antAzimuthPTypeID := m.getEquipPropTypeID(ctx, "Azimuth", AntennaEquipmentType.ID)
		antAmplitudPTypeID := m.getEquipPropTypeID(ctx, "Amplitud Beam", AntennaEquipmentType.ID)
		antGananciaPTypeID := m.getEquipPropTypeID(ctx, "Ganancia", AntennaEquipmentType.ID)
		antTiltElectricoPTypeID := m.getEquipPropTypeID(ctx, "Tilt Electrico", AntennaEquipmentType.ID)
		antTiltMecanicoPTypeID := m.getEquipPropTypeID(ctx, "Tilt Mecanico", AntennaEquipmentType.ID)
		antModeloJumperPTypeID := m.getEquipPropTypeID(ctx, "Modelo Jumper", AntennaEquipmentType.ID)
		antLongitudJumperPTypeID := m.getEquipPropTypeID(ctx, "Longitud Jumper", AntennaEquipmentType.ID)
		antPerdidaConectorPTypeID := m.getEquipPropTypeID(ctx, "Perdida Conector", AntennaEquipmentType.ID)

		TorreEquipmentType := m.getOrCreateEquipmentType(ctx, "Torre", 0, "", 0, []*models.PropertyTypeInput{
			{
				Name: "Tipo Zona",
				Type: "string",
			},
			{
				Name: "Clasificacion Instalacion",
				Type: "string",
			},
			{
				Name: "Estacion Msnm",
				Type: "string",
			},
			{
				Name: "Tipo",
				Type: "string",
			},
			{
				Name: "Propietario",
				Type: "string",
			},
			{
				Name: "Altura",
				Type: "string",
			},
			{
				Name: "Altura Edificio",
				Type: "string",
			},
		})

		torTipoZonaPTypeID := m.getEquipPropTypeID(ctx, "Tipo Zona", TorreEquipmentType.ID)
		torClasificacionInstalacionPTypeID := m.getEquipPropTypeID(ctx, "Clasificacion Instalacion", TorreEquipmentType.ID)
		torEstacionMsnmPTypeID := m.getEquipPropTypeID(ctx, "Estacion Msnm", TorreEquipmentType.ID)
		torTipoPTypeID := m.getEquipPropTypeID(ctx, "Tipo", TorreEquipmentType.ID)
		torPropietarioPTypeID := m.getEquipPropTypeID(ctx, "Propietario", TorreEquipmentType.ID)
		torAlturaPTypeID := m.getEquipPropTypeID(ctx, "Altura", TorreEquipmentType.ID)
		torAlturaEdificioPTypeID := m.getEquipPropTypeID(ctx, "Altura Edificio", TorreEquipmentType.ID)

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

			func() {
				defer func() {
					if v := recover(); v != nil {
						log.Error("process line panicked", zap.Strings("line", line), zap.Any("error", v))
					}
				}()
				CodigoUnicoEStacionVal := line[3]
				NombreEstacionVal := line[4]
				DepartamentoVal := line[6]
				ProvinciaVal := line[7]
				DistritoVal := line[8]
				NombreCentroPobladoVal := line[9]
				DireccionVal := line[10]
				latVal := line[11]
				longVal := line[12]

				torTipoZonaVal := line[13]
				torClasificacionInstalacionVal := line[14]
				torEstacionMsnmVal := line[15]
				torTipoTorreVal := line[16]
				torPropietarioTorreVal := line[17]
				torAlturaTorreVal := line[18]
				torAlturaEdificioVal := line[19]

				CodigoUnicoCeldaVal := line[20]
				NombreNodoVal := line[21]
				EtiquetaNodoVal := line[22]
				VendorRANVal := line[24]
				ModeloEquipoNodoVal := line[25]
				TipoEstacionBaseVal := line[27]
				TecnologiaVal := line[28]
				ControladorVal := line[29]
				LacTacVal := line[30]

				antNumeroVal := line[40]
				antAlturaVal := line[41]
				antMarcaVal := line[42]
				antNombreModeloVal := line[43]

				antAzimuthVal := line[45]
				antAmplitudBeamVal := line[46]
				antGananciaVal := line[47]
				antTiltElectricoVal := line[48]
				antTiltMecanicoVal := line[49]
				antModeloJumperVal := line[50]
				antLongitudJumperVal := line[50]
				antPerdidaConectorVal := line[50]

				parsedLat, _ := strconv.ParseFloat(latVal, 64)
				parsedLong, _ := strconv.ParseFloat(longVal, 64)

				if len(ModeloEquipoNodoVal) == 0 {
					ModeloEquipoNodoVal = "BTS"
				}

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
				Estacion = m.ClientFrom(ctx).Location.UpdateOne(Estacion).SetExternalID(CodigoUnicoEStacionVal).SaveX(ctx)

				BTSEquipmentType := m.getOrCreateEquipmentType(ctx, ModeloEquipoNodoVal, 0, "", 4, []*models.PropertyTypeInput{
					{
						Name: "Codigo Unico Celda",
						Type: "string",
					},
					{
						Name: "Nombre Nodo",
						Type: "string",
					},
					{
						Name: "Etiqueta Nodo",
						Type: "string",
					},
					{
						Name: "Vendor RAN",
						Type: "string",
					},
					{
						Name: "Modelo Equipo Nodo",
						Type: "string",
					},
					{
						Name: "Tipo Estacion Base",
						Type: "string",
					},
					{
						Name: "Tecnologia",
						Type: "string",
					},
					{
						Name: "Controlador",
						Type: "string",
					},
					{
						Name: "LacTac",
						Type: "string",
					},
				})

				cucPTypeID := m.getEquipPropTypeID(ctx, "Codigo Unico Celda", BTSEquipmentType.ID)
				nnPTypeID := m.getEquipPropTypeID(ctx, "Nombre Nodo", BTSEquipmentType.ID)
				enPTypeID := m.getEquipPropTypeID(ctx, "Etiqueta Nodo", BTSEquipmentType.ID)
				vrPTypeID := m.getEquipPropTypeID(ctx, "Vendor RAN", BTSEquipmentType.ID)
				menPTypeID := m.getEquipPropTypeID(ctx, "Modelo Equipo Nodo", BTSEquipmentType.ID)
				tebPTypeID := m.getEquipPropTypeID(ctx, "Tipo Estacion Base", BTSEquipmentType.ID)
				tPTypeID := m.getEquipPropTypeID(ctx, "Tecnologia", BTSEquipmentType.ID)
				cPTypeID := m.getEquipPropTypeID(ctx, "Controlador", BTSEquipmentType.ID)
				ltPTypeID := m.getEquipPropTypeID(ctx, "LacTac", BTSEquipmentType.ID)
				btsPropertyInput := []*models.PropertyInput{
					{
						PropertyTypeID:     cucPTypeID,
						StringValue:        &CodigoUnicoCeldaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     nnPTypeID,
						StringValue:        &NombreNodoVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     enPTypeID,
						StringValue:        &EtiquetaNodoVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     vrPTypeID,
						StringValue:        &VendorRANVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     menPTypeID,
						StringValue:        &ModeloEquipoNodoVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     tebPTypeID,
						StringValue:        &TipoEstacionBaseVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     tPTypeID,
						StringValue:        &TecnologiaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     cPTypeID,
						StringValue:        &ControladorVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     ltPTypeID,
						StringValue:        &LacTacVal,
						IsInstanceProperty: &instance,
					},
				}

				btsName := CodigoUnicoCeldaVal
				if len(btsName) == 0 {
					m.deleteEquipmentIfExists(ctx, mr, btsName, BTSEquipmentType, Estacion, nil)
					btsName = TipoEstacionBaseVal
				}

				// create specific BTS Equipment types
				btsEquipment, _, _ := m.getOrCreateEquipment(ctx, mr, btsName, BTSEquipmentType, nil, Estacion, nil, btsPropertyInput)

				// create Antenna instance
				antPropertyInput := []*models.PropertyInput{
					{
						PropertyTypeID:     antNumeroPTypeID,
						StringValue:        &antNumeroVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antAlturaPTypeID,
						StringValue:        &antAlturaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antMarcaPTypeID,
						StringValue:        &antMarcaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antNombreModelPTypeID,
						StringValue:        &antNombreModeloVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antAzimuthPTypeID,
						StringValue:        &antAzimuthVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antAmplitudPTypeID,
						StringValue:        &antAmplitudBeamVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antGananciaPTypeID,
						StringValue:        &antGananciaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antTiltElectricoPTypeID,
						StringValue:        &antTiltElectricoVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antTiltMecanicoPTypeID,
						StringValue:        &antTiltMecanicoVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antModeloJumperPTypeID,
						StringValue:        &antModeloJumperVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antLongitudJumperPTypeID,
						StringValue:        &antLongitudJumperVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     antPerdidaConectorPTypeID,
						StringValue:        &antPerdidaConectorVal,
						IsInstanceProperty: &instance,
					},
				}

				antName := antNombreModeloVal
				if len(antName) == 0 {
					antName = "Antenna"
				}

				m.deleteEquipmentIfExists(ctx, mr, antName, AntennaEquipmentType, Estacion, nil)

				if antAlturaVal != "" {
					antName = antName + "_" + string(antAlturaVal)
				} else {
					antName = antName + "_0"
				}

				if antAzimuthVal != "" {
					antName = antName + "_" + string(antAzimuthVal)
				} else {
					antName = antName + "_0"
				}

				antennaEquipment, _, _ := m.getOrCreateEquipment(ctx, mr, antName, AntennaEquipmentType, nil, Estacion, nil, antPropertyInput)

				m.addLinkBetweenEquipments(ctx, mr, *btsEquipment, *antennaEquipment)

				// create Torre instance
				torPropertyInput := []*models.PropertyInput{
					{
						PropertyTypeID:     torTipoZonaPTypeID,
						StringValue:        &torTipoZonaVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torClasificacionInstalacionPTypeID,
						StringValue:        &torClasificacionInstalacionVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torEstacionMsnmPTypeID,
						StringValue:        &torEstacionMsnmVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torTipoPTypeID,
						StringValue:        &torTipoTorreVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torPropietarioPTypeID,
						StringValue:        &torPropietarioTorreVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torAlturaPTypeID,
						StringValue:        &torAlturaTorreVal,
						IsInstanceProperty: &instance,
					},
					{
						PropertyTypeID:     torAlturaEdificioPTypeID,
						StringValue:        &torAlturaEdificioVal,
						IsInstanceProperty: &instance,
					},
				}
				m.getOrCreateEquipment(ctx, mr, "Torre", TorreEquipmentType, nil, Estacion, nil, torPropertyInput)
			}()
		}
		w.WriteHeader(http.StatusOK)
	}
}
