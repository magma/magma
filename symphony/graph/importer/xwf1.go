// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//// nolint
package importer

import (
	"fmt"
	"io"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

// ProcessEquipmentCSV  imports equipment from CSV file to DB
func (m *importer) ProcessXwf1CSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	log.Debug("ProcessCoollink1CSV -started")
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

		//unknownLocationType := m.getOrCreateLocationType(ctx, "Unknown", nil)
		areaLocationType := m.getOrCreateLocationType(ctx, "Area", nil)
		addressLocationType := m.getOrCreateLocationType(ctx, "Address", nil)
		siteLocationType := m.getOrCreateLocationType(ctx, "Site", []*models.PropertyTypeInput{
			{
				Name: "Meshes",
				Type: "string",
			},
			{
				Name: "Address",
				Type: "string",
			},
		})
		inverter500EquipmentType := m.getOrCreateEquipmentType(ctx, "Inverter 500VA", 0, "", 0, nil)
		inverter900EquipmentType := m.getOrCreateEquipmentType(ctx, "Inverter 900VA", 0, "", 0, nil)
		battery100ahEquipmentType := m.getOrCreateEquipmentType(ctx, "Battery 100AH", 0, "", 0, nil)
		battery200ahEquipmentType := m.getOrCreateEquipmentType(ctx, "Battery 200AH", 0, "", 0, nil)
		radioModelCBNEquipmentType := m.getOrCreateEquipmentType(ctx, "Radio Model for CBN Root(Radwin 5000 Series)", 0, "", 0, nil)
		radioModelCLLEquipmentType := m.getOrCreateEquipmentType(ctx, "Radio Model for CLL Root(Cambium Force 200", 0, "", 0, nil)
		mikrotikRB2011CLLEquipmentType := m.getOrCreateEquipmentType(ctx, "Mikrotik RB 2011", 0, "", 0, nil)
		mikrotikRB750GCLLEquipmentType := m.getOrCreateEquipmentType(ctx, "Mikrotik RB 750G", 0, "", 0, nil)
		POEEquipmentType := m.getOrCreateEquipmentType(ctx, "Cambium E500", 0, "", 0, nil)
		poeEquipmentType := m.getOrCreateEquipmentType(ctx, "POE", 0, "", 0, nil)
		telephoneEquipmentType := m.getOrCreateEquipmentType(ctx, "telephoneTechnoWX3-7.0", 0, "", 0, nil)

		rowID := 0
		for {
			rowID++
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				panic("cannot read row")
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

			snNumberVal := line[0]
			businessVal := line[1]
			retailerIdVal := line[2]
			macIdVal := line[3]
			rootVal := line[4]
			meshesVal := line[5]
			latVal := line[6]
			longVal := line[7]
			// totalJanVal := line[8]
			// totalFebVal := line[9]
			// totalMarVal := line[10]
			inverter500val := line[11]
			inverter900val := line[12]
			battery100AHVal := line[13]
			battery200AHVal := line[14]
			radioModelCBNVal := line[15]
			radioModelCLLVal := line[16]
			mikrotikVal := line[17]
			cambiumE500Val := line[18]
			poeVal := line[19]
			telephoneVal := line[20]
			addressVal := line[21]
			areaVal := line[22]

			if len(macIdVal) == 0 {
				log.Warn("Empty mac, skipping", zap.String("macIdVal", macIdVal))
				continue
			}

			props, err := m.ClientFrom(ctx).EquipmentType.Query().
				QueryPropertyTypes().
				Where(propertytype.Name("Mac")).
				QueryProperties().
				Where(property.StringVal(macIdVal)).
				All(ctx)

			if err != nil {
				log.Warn("ERROR! while finding Coollink mac", zap.Error(err), zap.String("macIdVal", macIdVal))
				panic("error")
			}

			if len(props) > 1 {
				log.Warn("Found many macs", zap.String("macIdVal", macIdVal))
				panic("error")
			}

			if len(props) < 1 {
				log.Warn("cannot find mac", zap.String("macIdVal", macIdVal))
				//panic("error")
				continue
			}
			//log.Debug("Updating hotspot properties", zap.String("macIdVal", macIdVal))

			if len(areaVal) < 1 {
				log.Warn("Empty area name", zap.String("areaVal", areaVal))
				//panic("error")
				continue
			}
			areaLoc, _ := m.getOrCreateLocation(ctx, areaVal, 0.0, 0.0, areaLocationType, nil, nil, nil)
			addressLoc, _ := m.getOrCreateLocation(ctx, addressVal, 0.0, 0.0, addressLocationType, &areaLoc.ID, nil, nil)

			accessPoint := props[0].QueryEquipment().OnlyX(ctx)
			hotspot := accessPoint.QueryLocation().OnlyX(ctx)
			siteAddress := areaVal + " " + addressVal
			site, _ := m.getOrCreateLocation(ctx, accessPoint.Name, 0.0, 0.0, siteLocationType, &hotspot.ID, []*models.PropertyInput{
				{
					PropertyTypeID: m.getLocPropTypeID(ctx, "Meshes", siteLocationType.ID),
					StringValue:    &meshesVal,
				},
				{
					PropertyTypeID: m.getLocPropTypeID(ctx, "Address", siteLocationType.ID),
					StringValue:    &siteAddress,
				},
			}, nil)
			if hotspot.QueryChildren().CountX(ctx) <= 1 || len(meshesVal) > 0 {
				if hotspot.QueryParent().OnlyXID(ctx) != addressLoc.ID {
					qh := m.ClientFrom(ctx).Location.UpdateOne(hotspot)
					qh.ClearParent()
					qh.SetParent(addressLoc)
					qh.SaveX(ctx)
				}
			}

			qa := m.ClientFrom(ctx).Equipment.UpdateOne(accessPoint)
			if accessPoint.QueryLocation().OnlyXID(ctx) != site.ID {
				qa.ClearLocation()
				qa.SetLocation(site)
			}

			if !m.propExistsOnEquipment(ctx, accessPoint, "snNumber", snNumberVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "snNumber", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(snNumberVal).
					SaveX(ctx))
			}
			if !m.propExistsOnEquipment(ctx, accessPoint, "Business", businessVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "Business", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(businessVal).
					SaveX(ctx))
			}

			if !m.propExistsOnEquipment(ctx, accessPoint, "Retailer Id", retailerIdVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "Retailer Id", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(retailerIdVal).
					SaveX(ctx))
			}

			if !m.propExistsOnEquipment(ctx, accessPoint, "Lat (coollink)", latVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "Lat (coollink)", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(latVal).
					SaveX(ctx))
			}

			if !m.propExistsOnEquipment(ctx, accessPoint, "Long (coollink)", longVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "Long (coollink)", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(longVal).
					SaveX(ctx))
			}

			if !m.propExistsOnEquipment(ctx, accessPoint, "Root", rootVal) {
				ptypeID := m.getEquipPropTypeID(ctx, "Root", accessPoint.QueryType().OnlyX(ctx).ID)
				qa = qa.AddProperties(m.ClientFrom(ctx).Property.Create().
					SetTypeID(ptypeID).
					SetStringVal(rootVal).
					SaveX(ctx))
			}

			qa.SaveX(ctx)

			if inverter500val == "1" {
				m.getOrCreateEquipment(ctx, mr, "inverter500", inverter500EquipmentType, nil, site, nil, nil)
			}
			if inverter900val == "1" {
				m.getOrCreateEquipment(ctx, mr, "inverter900", inverter900EquipmentType, nil, site, nil, nil)
			}
			if battery100AHVal == "1" {
				m.getOrCreateEquipment(ctx, mr, "battery100ah", battery100ahEquipmentType, nil, site, nil, nil)
			}
			if battery200AHVal == "1" {
				m.getOrCreateEquipment(ctx, mr, "battery200ah", battery200ahEquipmentType, nil, site, nil, nil)
			}
			if radioModelCBNVal == "Installed" {
				m.getOrCreateEquipment(ctx, mr, "radioModelCBN", radioModelCBNEquipmentType, nil, site, nil, nil)
			}
			if radioModelCLLVal == "Installed" {
				m.getOrCreateEquipment(ctx, mr, "radioModelCLLN", radioModelCLLEquipmentType, nil, site, nil, nil)
			}
			if mikrotikVal == "RB2011" {
				m.getOrCreateEquipment(ctx, mr, "RB2011", mikrotikRB2011CLLEquipmentType, nil, site, nil, nil)
			}
			if mikrotikVal == "RB750G" {
				m.getOrCreateEquipment(ctx, mr, "RB750G", mikrotikRB750GCLLEquipmentType, nil, site, nil, nil)
			}
			if cambiumE500Val == "Installed" {
				m.getOrCreateEquipment(ctx, mr, "Cambium E500", POEEquipmentType, nil, site, nil, nil)
			}
			if poeVal == "Installed" {
				m.getOrCreateEquipment(ctx, mr, "POE", poeEquipmentType, nil, site, nil, nil)
			}
			if telephoneVal == "Allocated" {
				m.getOrCreateEquipment(ctx, mr, "Telephone", telephoneEquipmentType, nil, site, nil, nil)
			}

		}
		log.Debug("Done!!")
		w.WriteHeader(http.StatusOK)
	}
}
