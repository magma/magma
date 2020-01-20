// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent/property"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"go.uber.org/zap"
)

const (
	buildingNum        = "Building Asset Number"
	buildingNumType    = "string"
	siteNum            = "Site Asset Number"
	siteNumType        = "string"
	gbicName           = "GBIC"
	gbicType           = "bool"
	cardPortType       = "OLT"
	linkTagType        = "string"
	portDesc           = "Neighborhood"
	portDescType       = "string"
	portStatus         = "Status"
	portStatusType     = "enum"
	linkTag            = "Tag"
	linkContractor     = "Contractor Name"
	linkContractorType = "string"
	linkPlanNum        = "Plan Number"
	linkPlanNumType    = "string"
	slotPrefix         = "Slot "
	equipStatus        = "Status"
	equipStatusType    = "enum"
	eInstalled         = "Installed"
	ePlanned           = "Planned"
)

// ProcessEquipmentCSV  imports equipment from CSV file to DB
// nolint: staticcheck
func (m *importer) ProcessFTTHCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	log.Debug("ProcessFTTHCSV - started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, fmt.Sprintf("%q: cannot parse form", err), http.StatusInternalServerError)
		return
	}
	ctx, mr := m.CloneContext(ctx), m.r.Mutation()

	// prerequisites for processing the FFTH file.
	portType := m.ensurePortType(ctx, client, log)
	hasGibic := portType.QueryPropertyTypes().Where(propertytype.Name(gbicName)).OnlyX(ctx)
	ptPortDesc := portType.QueryPropertyTypes().Where(propertytype.Name(portDesc)).OnlyX(ctx)
	ptLinkTag := portType.QueryLinkPropertyTypes().Where(propertytype.Name(linkTag)).OnlyX(ctx)

	for fileName := range r.MultipartForm.File {
		_, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("%q: cannot handle file %q", err, fileName), http.StatusInternalServerError)
			return
		}

		typeLocationType := m.getOrCreateLocationType(ctx, "Type", nil)
		cityLocationType := m.getOrCreateLocationType(ctx, "City", nil)
		sitePropertyTypeInput := []*models.PropertyTypeInput{
			{
				Name: siteNum,
				Type: siteNumType,
			},
		}
		siteLocationType := m.getOrCreateLocationType(ctx, "Site", sitePropertyTypeInput)
		siteNumPtypeID, _ := siteLocationType.QueryPropertyTypes().Where(propertytype.Name(siteNum)).OnlyID(ctx)

		manholeLocationType := m.getOrCreateLocationType(ctx, "Manhole", nil)
		bldngPropertyTypeInput := []*models.PropertyTypeInput{
			{
				Name: buildingNum,
				Type: buildingNumType,
			},
		}

		bldgLocationType := m.getOrCreateLocationType(ctx, "Building", bldngPropertyTypeInput)
		bPTypeID, _ := bldgLocationType.QueryPropertyTypes().Where(propertytype.Name(buildingNum)).OnlyID(ctx)

		statusVal := "[\"Installed\",\"Planned\"]"
		statusPropertyTypeInput := []*models.PropertyTypeInput{
			{
				Name:        equipStatus,
				Type:        equipStatusType,
				StringValue: &statusVal,
			},
		}

		m.getOrCreateEquipmentType(ctx, "Nokia FX4", 4, slotPrefix, 0, statusPropertyTypeInput)
		m.getOrCreateEquipmentType(ctx, "Nokia FX8", 8, slotPrefix, 0, statusPropertyTypeInput)
		m.getOrCreateEquipmentType(ctx, "Dasan F-219", 2, slotPrefix, 0, statusPropertyTypeInput)
		m.getOrCreateEquipmentType(ctx, "Dasan F-1419", 14, slotPrefix, 0, statusPropertyTypeInput)
		m.getOrCreateEquipmentType(ctx, "Dasan MXK-198", 2, slotPrefix, 0, statusPropertyTypeInput)
		card16Type := m.getOrCreateCard16Type(ctx, "Card", portType, statusPropertyTypeInput)
		card16StatusPTypeID, _ := card16Type.QueryPropertyTypes().Where(propertytype.Name(equipStatus)).OnlyID(ctx)

		// Manhole Splitters: 1:x
		m.ensureSplitterType(ctx, "Splitter 1:2", 1, 2)
		m.ensureSplitterType(ctx, "Splitter 1:4", 1, 4)
		m.ensureSplitterType(ctx, "Splitter 1:8", 1, 8)
		m.ensureSplitterType(ctx, "Splitter 1:16", 1, 16)
		m.ensureSplitterType(ctx, "Splitter 1:32", 1, 32)

		// Building Splitters: 2:x
		m.ensureSplitterType(ctx, "Splitter 2:4", 2, 4)
		m.ensureSplitterType(ctx, "Splitter 2:8", 2, 8)
		m.ensureSplitterType(ctx, "Splitter 2:16", 2, 16)
		m.ensureSplitterType(ctx, "Splitter 2:32", 2, 32)

		for row := 1; ; row++ {
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Error("failed reading line", zap.Error(err), zap.String("filename", fileName), zap.Int("row", row))
				http.Error(w, fmt.Sprintf("failed reading line #%d", row), http.StatusInternalServerError)
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

			cityName := line[0]
			siteName := line[1]
			siteNum := line[2]

			if cityName == "" || siteName == "" {
				continue
			}
			func() {
				defer func() {
					if v := recover(); v != nil {
						log.Error("panic recovering for line", zap.Any("cause", v), zap.String("stack", string(debug.Stack())), zap.Int("row", row))
					}
				}()
				city, _ := m.getOrCreateLocation(ctx, cityName, 0.0, 0.0, cityLocationType, nil, nil, nil)
				siteCity, _ := m.getOrCreateLocation(ctx, "Sites", 0.0, 0.0, typeLocationType, &city.ID, nil, nil)
				customersCity, _ := m.getOrCreateLocation(ctx, "Buildings", 0.0, 0.0, typeLocationType, &city.ID, nil, nil)
				manholesCity, _ := m.getOrCreateLocation(ctx, "Manholes", 0.0, 0.0, typeLocationType, &city.ID, nil, nil)

				var props []*models.PropertyInput
				if siteNum != "" {
					siteNum = fmt.Sprintf("%05s", siteNum)
					props = append(props, &models.PropertyInput{
						PropertyTypeID: siteNumPtypeID,
						StringValue:    &siteNum,
					})
				}
				site, _ := m.getOrCreateLocation(ctx, siteName, 0.0, 0.0, siteLocationType, &siteCity.ID, props, nil)
				if siteNum != "" && site.ExternalID != siteNum {
					site = client.Location.UpdateOne(site).SetExternalID(siteNum).SaveX(ctx)
				}

				equipTypeName := line[3]
				equipTypePrefix := ""
				if strings.Contains(equipTypeName, "Dasan") {
					equipTypePrefix = "DS"
				} else if strings.Contains(equipTypeName, "Nokia") {
					equipTypePrefix = "NK"
				}
				equipName := "GP" + equipTypePrefix + siteNum + "-" + line[4]
				equipType, err := client.EquipmentType.Query().Where(equipmenttype.Name(equipTypeName)).Only(ctx)
				if err != nil {
					log.Warn("can not find type name", zap.Error(err), zap.String("equipTypeName", equipTypeName))
					return
				}
				equipStatusPTypeID, err := equipType.QueryPropertyTypes().Where(propertytype.Name(equipStatus)).OnlyID(ctx)
				var equipProps []*models.PropertyInput
				equipStatus := eInstalled
				if err == nil {
					equipProps = append(equipProps, &models.PropertyInput{
						PropertyTypeID: equipStatusPTypeID,
						StringValue:    &equipStatus,
					})
				}
				equip, _, _ := m.getOrCreateEquipment(ctx, mr, equipName, equipType, nil, site, nil, equipProps)

				var statusVal string
				if hasgbic := line[7] != ""; hasgbic {
					statusVal = eInstalled
				} else {
					statusVal = ePlanned
				}

				cardProps := []*models.PropertyInput{
					{
						PropertyTypeID: card16StatusPTypeID,
						StringValue:    &statusVal,
					},
				}
				posName := slotPrefix + line[5]
				cardName := "Card " + line[5]
				position := equip.
					QueryPositions().
					Where(equipmentposition.HasDefinitionWith(equipmentpositiondefinition.Name(posName))).
					OnlyX(ctx)
				card, _, _ := m.getOrCreateEquipment(ctx, mr, cardName, card16Type, nil, nil, position, cardProps)
				if card == nil {
					log.Warn("failed to create card", zap.Int("row", row), zap.String("cardName", cardName), zap.String("position", position.ID))
					return
				}

				cardPortName := line[6]
				cardPort := card.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.Name(cardPortName))).OnlyX(ctx)
				if hasgbic := line[7] != ""; hasgbic {
					_, err = cardPort.QueryProperties().Where(property.HasTypeWith(propertytype.ID(hasGibic.ID))).Only(ctx)
					if ent.IsNotFound(err) {
						client.Property.Create().SetType(hasGibic).SetBoolVal(true).SetEquipmentPortID(cardPort.ID).SaveX(ctx)
					}
				}
				// if there's a description information for this link,
				// add it to the link properties.
				if desc := line[8]; desc != "" {
					_, err = cardPort.QueryProperties().Where(property.HasTypeWith(propertytype.ID(ptPortDesc.ID))).Only(ctx)
					if ent.IsNotFound(err) {
						client.Property.Create().SetType(ptPortDesc).SetStringVal(desc).SetEquipmentPortID(cardPort.ID).SaveX(ctx)
					}
				}

				manholeName := line[20]
				if len(manholeName) == 0 {
					return
				}
				manhole, _ := m.getOrCreateLocation(ctx, manholeName, 0, 0, manholeLocationType, &manholesCity.ID, nil, nil)
				splitterName := "MSP" + line[21]
				if len(splitterName) == 0 {
					log.Warn("bad splitter name", zap.Int("row", row), zap.String("splitterName", splitterName))
					return
				}
				splitterSizeDefinition := line[22]
				if len(splitterSizeDefinition) == 0 {
					return
				}

				splitterSizeArr := strings.Split(splitterSizeDefinition, ":")
				if len(splitterSizeArr) != 2 {
					log.Warn("bad splitter definition", zap.Int("row", row), zap.String("splitterSizeDefinition", splitterSizeDefinition))
					return
				}
				splitterSizeStr := strings.Replace(splitterSizeArr[1], "0", "", -1)
				splitterSizeInt, err := strconv.Atoi(splitterSizeStr)
				if err != nil {
					log.Warn("bad splitter size", zap.Int("row", row), zap.String("splitterSizeStr", splitterSizeStr))
					return
				}

				splitterType := client.EquipmentType.Query().Where(equipmenttype.Name("Splitter 1:" + splitterSizeStr)).OnlyX(ctx)
				splitter, _, _ := m.getOrCreateEquipment(ctx, mr, splitterName, splitterType, nil, manhole, nil, nil)

				for i := 0; i < splitterSizeInt; i++ {
					splitterOutPortName := "out" + strconv.Itoa(i+1)
					splitterOutPort := splitter.
						QueryPorts().
						Where(equipmentport.HasDefinitionWith(equipmentportdefinition.Name(splitterOutPortName))).
						OnlyX(ctx)

					bldgAddress, bldgNum := m.getBuildingDetails(ctx, row, line, i)
					if bldgAddress == nil || *bldgAddress == "" {
						continue
					}
					propertyInput := []*models.PropertyInput{
						{
							PropertyTypeID: bPTypeID,
							StringValue:    bldgNum,
						},
					}

					bldgName := *bldgAddress
					bldg, _ := m.getOrCreateLocation(ctx, bldgName, 0, 0, bldgLocationType, &customersCity.ID, propertyInput, nil)
					if bldgNum != nil && bldg.ExternalID != *bldgNum {
						if _, err := strconv.Atoi(*bldgNum); err == nil {
							bldg = client.Location.UpdateOne(bldg).SetExternalID(*bldgNum).SaveX(ctx)
						}
					}

					_, err = splitterOutPort.QueryLink().Only(ctx)
					if !ent.IsNotFound(err) {
						continue
					}

					bldgSPSize := 64 / splitterSizeInt
					bldgPort := m.getOrCreateBuildingInPort(ctx, bldg, bldgSPSize, "BSP1")
					l, err := m.linkBuildingAndManholePorts(ctx, bldgPort, splitterOutPort)
					if err != nil {
						log.Warn("Failed linking spliter and building", zap.Error(err), zap.Int("row", row), zap.String("splitterOutPortName", splitterOutPortName), zap.String("bldgAddress", *bldgAddress), zap.String("bldgNum", *bldgNum))
						continue
					}
					if l == nil {
						// connect manhole port to a new blgd splitter
						// for cases where bldg has 2 splitters
						bldgPort := m.getOrCreateBuildingInPort(ctx, bldg, bldgSPSize, "BSP2")
						_, err = m.linkBuildingAndManholePorts(ctx, bldgPort, splitterOutPort)
						if err != nil {
							log.Warn("Failed linking spliter and building", zap.Error(err), zap.Int("row", row), zap.String("splitterOutPortName", splitterOutPortName), zap.String("bldgAddress", *bldgAddress), zap.String("bldgNum", *bldgNum))
							continue
						}
					}
				}

				// nolint: goconst
				splitterInPortName := "in"
				splitterInPort := splitter.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.Name(splitterInPortName))).OnlyX(ctx)
				_, err = cardPort.QueryLink().Only(ctx)

				if ent.IsNotFound(err) {
					creator := client.Link.Create().AddPorts(cardPort).AddPorts(splitterInPort)
					if tag := line[19]; tag != "" {
						linkProp := client.Property.Create().SetType(ptLinkTag).SetStringVal(tag).SaveX(ctx)
						creator.AddProperties(linkProp)
					}
					if _, err := creator.Save(ctx); err != nil {
						log.Warn("Failed linking card and splitter", zap.Error(err), zap.Int("row", row), zap.String("cardName", cardName), zap.String("cardPortName", cardPortName), zap.String("splitterName", splitterName))
					}
				}
			}()
		}

		log.Debug("Done uploading")
		w.WriteHeader(http.StatusOK)
	}
}

func (m importer) linkBuildingAndManholePorts(ctx context.Context, bldgPort *ent.EquipmentPort, splitterOutPort *ent.EquipmentPort) (*ent.Link, error) {
	if bldgPort == nil || splitterOutPort == nil {
		return nil, nil
	}
	client := m.ClientFrom(ctx)
	l, err := splitterOutPort.QueryLink().Only(ctx)
	if ent.IsNotFound(err) {
		// connect manhole port to building port
		_, err := bldgPort.QueryLink().Only(ctx)
		if ent.IsNotFound(err) {
			return client.Link.Create().AddPorts(splitterOutPort).AddPorts(bldgPort).Save(ctx)
		}
		// building is already connected to a different manhole
		return nil, nil
	}
	return l, nil
}

func (m *importer) getOrCreateBuildingInPort(ctx context.Context, bldg *ent.Location, sPSize int, name string) *ent.EquipmentPort {
	client, mr := m.ClientFrom(ctx), m.r.Mutation()
	bldgPortType := client.EquipmentType.Query().Where(equipmenttype.Name(fmt.Sprintf("Splitter 2:%d", sPSize))).OnlyX(ctx)
	bldgPortEquipment, _, _ := m.getOrCreateEquipment(ctx, mr, name, bldgPortType, nil, bldg, nil, nil)
	if bldgPortEquipment == nil {
		return nil
	}

	// nolint: goconst
	bldgPortName := "in1"
	bldgPort := bldgPortEquipment.
		QueryPorts().
		Where(equipmentport.HasDefinitionWith(equipmentportdefinition.Name(bldgPortName))).
		OnlyX(ctx)

	return bldgPort
}

func (m *importer) getBuildingDetails(ctx context.Context, rowID int, line []string, i int) (*string, *string) {
	log := m.log.For(ctx)
	bldgAddressLineIndex := 22 + 1 + i*2
	bldgNumLineIndex := 22 + 2 + i*2
	if len(line) <= bldgAddressLineIndex {
		log.Warn("Building address does not exist", zap.Int("rowID", rowID), zap.Int("i", i))
		return nil, nil
	}
	if len(line) <= bldgNumLineIndex {
		log.Warn("Building number does not exist", zap.Int("rowID", rowID), zap.Int("i", i))
		return nil, nil
	}
	bldgAddress := line[bldgAddressLineIndex]
	bldgNum := line[bldgNumLineIndex]

	if bldgAddress == "fb_import_empty" {
		if len(bldgNum) == 0 {
			return nil, nil
		}
		log.Warn("fb_import_empty but num not empty", zap.Int("rowID", rowID), zap.Int("i", i))
		return nil, nil
	}

	return &bldgAddress, &bldgNum
}

func (m *importer) getOrCreateCard16Type(ctx context.Context, name string, portType *ent.EquipmentPortType, props []*models.PropertyTypeInput) *ent.EquipmentType {
	log := m.log.For(ctx)
	client := m.ClientFrom(ctx)

	equipmentType, err := m.ClientFrom(ctx).EquipmentType.Query().Where(equipmenttype.Name(name)).Only(ctx)
	if equipmentType != nil {
		return equipmentType
	}
	if !ent.IsNotFound(err) {
		panic(err)
	}
	var proprArr []*ent.PropertyType
	for _, input := range props {
		propEnt := client.PropertyType.
			Create().
			SetName(input.Name).
			SetType(input.Type.String()).
			SetNillableStringVal(input.StringValue).
			SetNillableIntVal(input.IntValue).
			SetNillableBoolVal(input.BooleanValue).
			SetNillableFloatVal(input.FloatValue).
			SetNillableLatitudeVal(input.LatitudeValue).
			SetNillableLongitudeVal(input.LongitudeValue).
			SetNillableIsInstanceProperty(input.IsInstanceProperty).
			SetNillableEditable(input.IsEditable).
			SaveX(ctx)
		proprArr = append(proprArr, propEnt)
	}
	q := client.EquipmentType.Create().SetName(name).AddPropertyTypes(proprArr...)
	for i := 1; i <= 16; i++ {
		p := client.EquipmentPortDefinition.Create().
			SetName(strconv.Itoa(i)).
			SetEquipmentPortType(portType).
			SaveX(ctx)
		q.AddPortDefinitions(p)
	}
	log.Debug("Creating new card16 type", zap.String("name", name))
	return q.SaveX(ctx)
}

func (m *importer) ensurePortType(ctx context.Context, client *ent.Client, log *zap.Logger) *ent.EquipmentPortType {
	typ := client.EquipmentPortType.
		Query().
		Where(equipmentporttype.Name(cardPortType)).
		FirstX(ctx)
	if typ != nil {
		return typ
	}
	log.Debug("creating equipment port type for ftth")
	pGibic := client.PropertyType.
		Create().
		SetName(gbicName).
		SetType(gbicType).
		SetBoolVal(false).
		SaveX(ctx)
	pPortDesc := client.PropertyType.
		Create().
		SetName(portDesc).
		SetType(portDescType).
		SaveX(ctx)
	pPortStatus := client.PropertyType.
		Create().
		SetName(portStatus).
		SetType(portStatusType).
		SetStringVal("[\"Active\",\"Free\",\"Reserved\"]").
		SaveX(ctx)
	plinkTag := client.PropertyType.
		Create().
		SetName(linkTag).
		SetType(linkTagType).
		SaveX(ctx)
	plinkCtrctor := client.PropertyType.
		Create().
		SetName(linkContractor).
		SetType(linkContractorType).
		SaveX(ctx)
	plinkPlanNum := client.PropertyType.
		Create().
		SetName(linkPlanNum).
		SetType(linkPlanNumType).
		SaveX(ctx)
	return client.EquipmentPortType.
		Create().
		SetName(cardPortType).
		AddPropertyTypes(pGibic, pPortDesc, pPortStatus).
		AddLinkPropertyTypes(plinkTag, plinkCtrctor, plinkPlanNum).
		SaveX(ctx)
}
