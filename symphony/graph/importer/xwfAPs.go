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

// ProcessEquipmentCSV  imports equipment from CSV file to DB
func (m *importer) ProcessXwfApsCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	log.Debug("ProcessCoollink2CSV -started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusUnprocessableEntity)
		return
	}
	ctx, mr := m.CloneContext(ctx), m.r.Mutation()

	for fileName := range r.MultipartForm.File {
		_, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusUnprocessableEntity)
			return
		}

		rowID := 0
		unknownLocationType := m.getOrCreateLocationType(ctx, "Unknown", nil)
		hotspotLocationType := m.getOrCreateLocationType(ctx, "Hotspot", []*models.PropertyTypeInput{
			{
				Name: "XPP ID",
				Type: "string",
			},
			{
				Name: "XPP Launch Date",
				Type: "string",
			},
			{
				Name: "XPP Status",
				Type: "string",
			},
			{
				Name: "XPP City",
				Type: "string",
			},
			{
				Name: "XPP Street",
				Type: "string",
			},
			{
				Name: "XPP Zip",
				Type: "string",
			},
		})

		accesspointEquipmentType := m.getOrCreateEquipmentType(ctx, "FB-Access-point", 0, "", 0, []*models.PropertyTypeInput{
			{
				Name: "XPP ID",
				Type: "string",
			},
			{
				Name: "Mac",
				Type: "string",
			},
			{
				Name: "XPP IP",
				Type: "string",
			},
			{
				Name: "XPP Status",
				Type: "string",
			},
			{
				Name: "snNumber",
				Type: "string",
			},
			{
				Name: "Business",
				Type: "string",
			},
			{
				Name: "Retailer Id",
				Type: "string",
			},
			{
				Name: "Lat (coollink)",
				Type: "string",
			},
			{
				Name: "Long (coollink)",
				Type: "string",
			},
			{
				Name: "Root",
				Type: "string",
			},
			{
				Name: "Meshes",
				Type: "string",
			},
		})

		for {
			rowID++
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
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

			hotspotIdVal := line[0]
			hotspotLatVal := line[1]
			hotspotLongVal := line[2]
			hotspotLaunch_dateVal := line[3]
			hotspotNameVal := line[4]
			if hotspotNameVal == "" {
				hotspotNameVal = "no hotspot name"
			}
			hotspot_cityVal := line[5]
			if hotspot_cityVal == "" {
				hotspot_cityVal = "no city"
			}
			hotspotStreetVal := line[6]
			if hotspotStreetVal == "" {
				hotspotStreetVal = "no street"
			}
			hotspot_zipVal := line[7]
			if hotspot_zipVal == "" {
				hotspot_zipVal = "no zip"
			}
			hotspotStatusVal := line[8]
			apIdVal := line[9]
			apNameVal := line[10]
			if apNameVal == "" {
				apNameVal = "no ap name"
			}
			apMacVal := line[11]
			apIpVal := line[12]
			apStatusVal := line[13]

			hotspotLat, _ := strconv.ParseFloat(hotspotLatVal, 64)
			hotspotLong, _ := strconv.ParseFloat(hotspotLongVal, 64)

			unknownLoc, _ := m.getOrCreateLocation(ctx, "_FB XPP data", 0.0, 0.0, unknownLocationType, nil, nil, nil)

			xppIDPTypeID := m.getLocPropTypeID(ctx, "XPP ID", hotspotLocationType.ID)
			ldPTypeID := m.getLocPropTypeID(ctx, "XPP Launch Date", hotspotLocationType.ID)
			statusPTypeID := m.getLocPropTypeID(ctx, "XPP Status", hotspotLocationType.ID)
			cityPTypeID := m.getLocPropTypeID(ctx, "XPP City", hotspotLocationType.ID)
			streetPTypeID := m.getLocPropTypeID(ctx, "XPP Street", hotspotLocationType.ID)
			zipPTypeID := m.getLocPropTypeID(ctx, "XPP Zip", hotspotLocationType.ID)

			hotspotPropertyInput := []*models.PropertyInput{
				{
					PropertyTypeID: xppIDPTypeID,
					StringValue:    &hotspotIdVal,
				},
				{
					PropertyTypeID: ldPTypeID,
					StringValue:    &hotspotLaunch_dateVal,
				},
				{
					PropertyTypeID: statusPTypeID,
					StringValue:    &hotspotStatusVal,
				},
				{
					PropertyTypeID: cityPTypeID,
					StringValue:    &hotspot_cityVal,
				},
				{
					PropertyTypeID: streetPTypeID,
					StringValue:    &hotspotStreetVal,
				},
				{
					PropertyTypeID: zipPTypeID,
					StringValue:    &hotspot_zipVal,
				},
			}
			parent := unknownLoc
			if hotspotStatusVal == "shut_down" {
				continue
			}

			hotspot, _ := m.getOrCreateLocation(ctx, hotspotNameVal, hotspotLat, hotspotLong, hotspotLocationType, &parent.ID, hotspotPropertyInput, nil)

			xppIDPTypeID = m.getEquipPropTypeID(ctx, "XPP ID", accesspointEquipmentType.ID)
			macPTypeID := m.getEquipPropTypeID(ctx, "Mac", accesspointEquipmentType.ID)
			ipPTypeID := m.getEquipPropTypeID(ctx, "XPP IP", accesspointEquipmentType.ID)
			statusPTypeID = m.getEquipPropTypeID(ctx, "XPP Status", accesspointEquipmentType.ID)

			accessPointPropertyInput := []*models.PropertyInput{
				{
					PropertyTypeID: xppIDPTypeID,
					StringValue:    &apIdVal,
				},
				{
					PropertyTypeID: macPTypeID,
					StringValue:    &apMacVal,
				},
				{
					PropertyTypeID: ipPTypeID,
					StringValue:    &apIpVal,
				},
				{
					PropertyTypeID: statusPTypeID,
					StringValue:    &apStatusVal,
				},
			}
			m.getOrCreateEquipment(ctx, mr, "ap_"+apNameVal, accesspointEquipmentType, nil, hotspot, nil, accessPointPropertyInput)

		}
		log.Debug("Done!!")
		w.WriteHeader(http.StatusOK)
	}
}
