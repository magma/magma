// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleActionsRule() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the actionsrule's edges.

	// create actionsrule vertex with its edges.
	ar := client.ActionsRule.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetTriggerID("string").
		SetRuleFilters(nil).
		SetRuleActions(nil).
		SaveX(ctx)
	log.Println("actionsrule created:", ar)

	// query edges.

	// Output:
}
func ExampleCheckListItem() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the checklistitem's edges.

	// create checklistitem vertex with its edges.
	cli := client.CheckListItem.
		Create().
		SetTitle("string").
		SetType("string").
		SetIndex(1).
		SetChecked(true).
		SetStringVal("string").
		SetEnumValues("string").
		SetHelpText("string").
		SaveX(ctx)
	log.Println("checklistitem created:", cli)

	// query edges.

	// Output:
}
func ExampleCheckListItemDefinition() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the checklistitemdefinition's edges.

	// create checklistitemdefinition vertex with its edges.
	clid := client.CheckListItemDefinition.
		Create().
		SetTitle("string").
		SetType("string").
		SetIndex(1).
		SetEnumValues("string").
		SetHelpText("string").
		SaveX(ctx)
	log.Println("checklistitemdefinition created:", clid)

	// query edges.

	// Output:
}
func ExampleComment() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the comment's edges.

	// create comment vertex with its edges.
	c := client.Comment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetAuthorName("string").
		SetText("string").
		SaveX(ctx)
	log.Println("comment created:", c)

	// query edges.

	// Output:
}
func ExampleCustomer() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the customer's edges.

	// create customer vertex with its edges.
	c := client.Customer.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SaveX(ctx)
	log.Println("customer created:", c)

	// query edges.

	// Output:
}
func ExampleEquipment() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipment's edges.
	et0 := client.EquipmentType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SaveX(ctx)
	log.Println("equipmenttype created:", et0)
	ep3 := client.EquipmentPosition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SaveX(ctx)
	log.Println("equipmentposition created:", ep3)
	ep4 := client.EquipmentPort.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SaveX(ctx)
	log.Println("equipmentport created:", ep4)
	wo5 := client.WorkOrder.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetStatus("string").
		SetPriority("string").
		SetDescription("string").
		SetOwnerName("string").
		SetInstallDate(time.Now()).
		SetCreationDate(time.Now()).
		SetAssignee("string").
		SetIndex(1).
		SaveX(ctx)
	log.Println("workorder created:", wo5)
	pr6 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr6)
	f7 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f7)

	// create equipment vertex with its edges.
	e := client.Equipment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetFutureState("string").
		SetDeviceID("string").
		SetExternalID("string").
		SetType(et0).
		AddPositions(ep3).
		AddPorts(ep4).
		SetWorkOrder(wo5).
		AddProperties(pr6).
		AddFiles(f7).
		SaveX(ctx)
	log.Println("equipment created:", e)

	// query edges.
	et0, err = e.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", et0)

	ep3, err = e.QueryPositions().First(ctx)
	if err != nil {
		log.Fatalf("failed querying positions: %v", err)
	}
	log.Println("positions found:", ep3)

	ep4, err = e.QueryPorts().First(ctx)
	if err != nil {
		log.Fatalf("failed querying ports: %v", err)
	}
	log.Println("ports found:", ep4)

	wo5, err = e.QueryWorkOrder().First(ctx)
	if err != nil {
		log.Fatalf("failed querying work_order: %v", err)
	}
	log.Println("work_order found:", wo5)

	pr6, err = e.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr6)

	f7, err = e.QueryFiles().First(ctx)
	if err != nil {
		log.Fatalf("failed querying files: %v", err)
	}
	log.Println("files found:", f7)

	// Output:
}
func ExampleEquipmentCategory() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentcategory's edges.

	// create equipmentcategory vertex with its edges.
	ec := client.EquipmentCategory.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SaveX(ctx)
	log.Println("equipmentcategory created:", ec)

	// query edges.

	// Output:
}
func ExampleEquipmentPort() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentport's edges.
	epd0 := client.EquipmentPortDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetBandwidth("string").
		SetVisibilityLabel("string").
		SaveX(ctx)
	log.Println("equipmentportdefinition created:", epd0)
	l2 := client.Link.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFutureState("string").
		SaveX(ctx)
	log.Println("link created:", l2)
	pr3 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr3)

	// create equipmentport vertex with its edges.
	ep := client.EquipmentPort.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetDefinition(epd0).
		SetLink(l2).
		AddProperties(pr3).
		SaveX(ctx)
	log.Println("equipmentport created:", ep)

	// query edges.
	epd0, err = ep.QueryDefinition().First(ctx)
	if err != nil {
		log.Fatalf("failed querying definition: %v", err)
	}
	log.Println("definition found:", epd0)

	l2, err = ep.QueryLink().First(ctx)
	if err != nil {
		log.Fatalf("failed querying link: %v", err)
	}
	log.Println("link found:", l2)

	pr3, err = ep.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr3)

	// Output:
}
func ExampleEquipmentPortDefinition() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentportdefinition's edges.
	ept0 := client.EquipmentPortType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SaveX(ctx)
	log.Println("equipmentporttype created:", ept0)

	// create equipmentportdefinition vertex with its edges.
	epd := client.EquipmentPortDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetBandwidth("string").
		SetVisibilityLabel("string").
		SetEquipmentPortType(ept0).
		SaveX(ctx)
	log.Println("equipmentportdefinition created:", epd)

	// query edges.
	ept0, err = epd.QueryEquipmentPortType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying equipment_port_type: %v", err)
	}
	log.Println("equipment_port_type found:", ept0)

	// Output:
}
func ExampleEquipmentPortType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentporttype's edges.
	pt0 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt0)
	pt1 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt1)

	// create equipmentporttype vertex with its edges.
	ept := client.EquipmentPortType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		AddPropertyTypes(pt0).
		AddLinkPropertyTypes(pt1).
		SaveX(ctx)
	log.Println("equipmentporttype created:", ept)

	// query edges.
	pt0, err = ept.QueryPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying property_types: %v", err)
	}
	log.Println("property_types found:", pt0)

	pt1, err = ept.QueryLinkPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying link_property_types: %v", err)
	}
	log.Println("link_property_types found:", pt1)

	// Output:
}
func ExampleEquipmentPosition() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentposition's edges.
	epd0 := client.EquipmentPositionDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetVisibilityLabel("string").
		SaveX(ctx)
	log.Println("equipmentpositiondefinition created:", epd0)
	e2 := client.Equipment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetFutureState("string").
		SetDeviceID("string").
		SetExternalID("string").
		SaveX(ctx)
	log.Println("equipment created:", e2)

	// create equipmentposition vertex with its edges.
	ep := client.EquipmentPosition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetDefinition(epd0).
		SetAttachment(e2).
		SaveX(ctx)
	log.Println("equipmentposition created:", ep)

	// query edges.
	epd0, err = ep.QueryDefinition().First(ctx)
	if err != nil {
		log.Fatalf("failed querying definition: %v", err)
	}
	log.Println("definition found:", epd0)

	e2, err = ep.QueryAttachment().First(ctx)
	if err != nil {
		log.Fatalf("failed querying attachment: %v", err)
	}
	log.Println("attachment found:", e2)

	// Output:
}
func ExampleEquipmentPositionDefinition() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmentpositiondefinition's edges.

	// create equipmentpositiondefinition vertex with its edges.
	epd := client.EquipmentPositionDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetVisibilityLabel("string").
		SaveX(ctx)
	log.Println("equipmentpositiondefinition created:", epd)

	// query edges.

	// Output:
}
func ExampleEquipmentType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the equipmenttype's edges.
	epd0 := client.EquipmentPortDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetBandwidth("string").
		SetVisibilityLabel("string").
		SaveX(ctx)
	log.Println("equipmentportdefinition created:", epd0)
	epd1 := client.EquipmentPositionDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetIndex(1).
		SetVisibilityLabel("string").
		SaveX(ctx)
	log.Println("equipmentpositiondefinition created:", epd1)
	pt2 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt2)
	ec4 := client.EquipmentCategory.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SaveX(ctx)
	log.Println("equipmentcategory created:", ec4)

	// create equipmenttype vertex with its edges.
	et := client.EquipmentType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		AddPortDefinitions(epd0).
		AddPositionDefinitions(epd1).
		AddPropertyTypes(pt2).
		SetCategory(ec4).
		SaveX(ctx)
	log.Println("equipmenttype created:", et)

	// query edges.
	epd0, err = et.QueryPortDefinitions().First(ctx)
	if err != nil {
		log.Fatalf("failed querying port_definitions: %v", err)
	}
	log.Println("port_definitions found:", epd0)

	epd1, err = et.QueryPositionDefinitions().First(ctx)
	if err != nil {
		log.Fatalf("failed querying position_definitions: %v", err)
	}
	log.Println("position_definitions found:", epd1)

	pt2, err = et.QueryPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying property_types: %v", err)
	}
	log.Println("property_types found:", pt2)

	ec4, err = et.QueryCategory().First(ctx)
	if err != nil {
		log.Fatalf("failed querying category: %v", err)
	}
	log.Println("category found:", ec4)

	// Output:
}
func ExampleFile() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the file's edges.

	// create file vertex with its edges.
	f := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f)

	// query edges.

	// Output:
}
func ExampleFloorPlan() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the floorplan's edges.
	l0 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l0)
	fprp1 := client.FloorPlanReferencePoint.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetX(1).
		SetY(1).
		SetLatitude(1).
		SetLongitude(1).
		SaveX(ctx)
	log.Println("floorplanreferencepoint created:", fprp1)
	fps2 := client.FloorPlanScale.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetReferencePoint1X(1).
		SetReferencePoint1Y(1).
		SetReferencePoint2X(1).
		SetReferencePoint2Y(1).
		SetScaleInMeters(1).
		SaveX(ctx)
	log.Println("floorplanscale created:", fps2)
	f3 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f3)

	// create floorplan vertex with its edges.
	fp := client.FloorPlan.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetLocation(l0).
		SetReferencePoint(fprp1).
		SetScale(fps2).
		SetImage(f3).
		SaveX(ctx)
	log.Println("floorplan created:", fp)

	// query edges.
	l0, err = fp.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l0)

	fprp1, err = fp.QueryReferencePoint().First(ctx)
	if err != nil {
		log.Fatalf("failed querying reference_point: %v", err)
	}
	log.Println("reference_point found:", fprp1)

	fps2, err = fp.QueryScale().First(ctx)
	if err != nil {
		log.Fatalf("failed querying scale: %v", err)
	}
	log.Println("scale found:", fps2)

	f3, err = fp.QueryImage().First(ctx)
	if err != nil {
		log.Fatalf("failed querying image: %v", err)
	}
	log.Println("image found:", f3)

	// Output:
}
func ExampleFloorPlanReferencePoint() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the floorplanreferencepoint's edges.

	// create floorplanreferencepoint vertex with its edges.
	fprp := client.FloorPlanReferencePoint.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetX(1).
		SetY(1).
		SetLatitude(1).
		SetLongitude(1).
		SaveX(ctx)
	log.Println("floorplanreferencepoint created:", fprp)

	// query edges.

	// Output:
}
func ExampleFloorPlanScale() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the floorplanscale's edges.

	// create floorplanscale vertex with its edges.
	fps := client.FloorPlanScale.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetReferencePoint1X(1).
		SetReferencePoint1Y(1).
		SetReferencePoint2X(1).
		SetReferencePoint2Y(1).
		SetScaleInMeters(1).
		SaveX(ctx)
	log.Println("floorplanscale created:", fps)

	// query edges.

	// Output:
}
func ExampleLink() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the link's edges.
	wo1 := client.WorkOrder.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetStatus("string").
		SetPriority("string").
		SetDescription("string").
		SetOwnerName("string").
		SetInstallDate(time.Now()).
		SetCreationDate(time.Now()).
		SetAssignee("string").
		SetIndex(1).
		SaveX(ctx)
	log.Println("workorder created:", wo1)
	pr2 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr2)

	// create link vertex with its edges.
	l := client.Link.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFutureState("string").
		SetWorkOrder(wo1).
		AddProperties(pr2).
		SaveX(ctx)
	log.Println("link created:", l)

	// query edges.

	wo1, err = l.QueryWorkOrder().First(ctx)
	if err != nil {
		log.Fatalf("failed querying work_order: %v", err)
	}
	log.Println("work_order found:", wo1)

	pr2, err = l.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr2)

	// Output:
}
func ExampleLocation() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the location's edges.
	lt0 := client.LocationType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetSite(true).
		SetName("string").
		SetMapType("string").
		SetMapZoomLevel(1).
		SetIndex(1).
		SaveX(ctx)
	log.Println("locationtype created:", lt0)
	l2 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l2)
	f3 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f3)
	e4 := client.Equipment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetFutureState("string").
		SetDeviceID("string").
		SetExternalID("string").
		SaveX(ctx)
	log.Println("equipment created:", e4)
	pr5 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr5)

	// create location vertex with its edges.
	l := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SetType(lt0).
		AddChildren(l2).
		AddFiles(f3).
		AddEquipment(e4).
		AddProperties(pr5).
		SaveX(ctx)
	log.Println("location created:", l)

	// query edges.
	lt0, err = l.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", lt0)

	l2, err = l.QueryChildren().First(ctx)
	if err != nil {
		log.Fatalf("failed querying children: %v", err)
	}
	log.Println("children found:", l2)

	f3, err = l.QueryFiles().First(ctx)
	if err != nil {
		log.Fatalf("failed querying files: %v", err)
	}
	log.Println("files found:", f3)

	e4, err = l.QueryEquipment().First(ctx)
	if err != nil {
		log.Fatalf("failed querying equipment: %v", err)
	}
	log.Println("equipment found:", e4)

	pr5, err = l.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr5)

	// Output:
}
func ExampleLocationType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the locationtype's edges.
	pt1 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt1)
	stc2 := client.SurveyTemplateCategory.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetCategoryTitle("string").
		SetCategoryDescription("string").
		SaveX(ctx)
	log.Println("surveytemplatecategory created:", stc2)

	// create locationtype vertex with its edges.
	lt := client.LocationType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetSite(true).
		SetName("string").
		SetMapType("string").
		SetMapZoomLevel(1).
		SetIndex(1).
		AddPropertyTypes(pt1).
		AddSurveyTemplateCategories(stc2).
		SaveX(ctx)
	log.Println("locationtype created:", lt)

	// query edges.

	pt1, err = lt.QueryPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying property_types: %v", err)
	}
	log.Println("property_types found:", pt1)

	stc2, err = lt.QuerySurveyTemplateCategories().First(ctx)
	if err != nil {
		log.Fatalf("failed querying survey_template_categories: %v", err)
	}
	log.Println("survey_template_categories found:", stc2)

	// Output:
}
func ExampleProject() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the project's edges.
	l1 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l1)
	c2 := client.Comment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetAuthorName("string").
		SetText("string").
		SaveX(ctx)
	log.Println("comment created:", c2)
	wo3 := client.WorkOrder.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetStatus("string").
		SetPriority("string").
		SetDescription("string").
		SetOwnerName("string").
		SetInstallDate(time.Now()).
		SetCreationDate(time.Now()).
		SetAssignee("string").
		SetIndex(1).
		SaveX(ctx)
	log.Println("workorder created:", wo3)
	pr4 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr4)

	// create project vertex with its edges.
	pr := client.Project.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		SetCreator("string").
		SetLocation(l1).
		AddComments(c2).
		AddWorkOrders(wo3).
		AddProperties(pr4).
		SaveX(ctx)
	log.Println("project created:", pr)

	// query edges.

	l1, err = pr.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l1)

	c2, err = pr.QueryComments().First(ctx)
	if err != nil {
		log.Fatalf("failed querying comments: %v", err)
	}
	log.Println("comments found:", c2)

	wo3, err = pr.QueryWorkOrders().First(ctx)
	if err != nil {
		log.Fatalf("failed querying work_orders: %v", err)
	}
	log.Println("work_orders found:", wo3)

	pr4, err = pr.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr4)

	// Output:
}
func ExampleProjectType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the projecttype's edges.
	pr0 := client.Project.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		SetCreator("string").
		SaveX(ctx)
	log.Println("project created:", pr0)
	pt1 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt1)
	wod2 := client.WorkOrderDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIndex(1).
		SaveX(ctx)
	log.Println("workorderdefinition created:", wod2)

	// create projecttype vertex with its edges.
	pt := client.ProjectType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		AddProjects(pr0).
		AddProperties(pt1).
		AddWorkOrders(wod2).
		SaveX(ctx)
	log.Println("projecttype created:", pt)

	// query edges.
	pr0, err = pt.QueryProjects().First(ctx)
	if err != nil {
		log.Fatalf("failed querying projects: %v", err)
	}
	log.Println("projects found:", pr0)

	pt1, err = pt.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pt1)

	wod2, err = pt.QueryWorkOrders().First(ctx)
	if err != nil {
		log.Fatalf("failed querying work_orders: %v", err)
	}
	log.Println("work_orders found:", wod2)

	// Output:
}
func ExampleProperty() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the property's edges.
	pt0 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt0)
	e8 := client.Equipment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetFutureState("string").
		SetDeviceID("string").
		SetExternalID("string").
		SaveX(ctx)
	log.Println("equipment created:", e8)
	l9 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l9)
	s10 := client.Service.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetStatus("string").
		SaveX(ctx)
	log.Println("service created:", s10)

	// create property vertex with its edges.
	pr := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SetType(pt0).
		SetEquipmentValue(e8).
		SetLocationValue(l9).
		SetServiceValue(s10).
		SaveX(ctx)
	log.Println("property created:", pr)

	// query edges.
	pt0, err = pr.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", pt0)

	e8, err = pr.QueryEquipmentValue().First(ctx)
	if err != nil {
		log.Fatalf("failed querying equipment_value: %v", err)
	}
	log.Println("equipment_value found:", e8)

	l9, err = pr.QueryLocationValue().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location_value: %v", err)
	}
	log.Println("location_value found:", l9)

	s10, err = pr.QueryServiceValue().First(ctx)
	if err != nil {
		log.Fatalf("failed querying service_value: %v", err)
	}
	log.Println("service_value found:", s10)

	// Output:
}
func ExamplePropertyType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the propertytype's edges.

	// create propertytype vertex with its edges.
	pt := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt)

	// query edges.

	// Output:
}
func ExampleService() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the service's edges.
	st0 := client.ServiceType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetHasCustomer(true).
		SaveX(ctx)
	log.Println("servicetype created:", st0)
	s2 := client.Service.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetStatus("string").
		SaveX(ctx)
	log.Println("service created:", s2)
	pr3 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr3)
	l4 := client.Link.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFutureState("string").
		SaveX(ctx)
	log.Println("link created:", l4)
	c5 := client.Customer.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SaveX(ctx)
	log.Println("customer created:", c5)
	se6 := client.ServiceEndpoint.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetRole("string").
		SaveX(ctx)
	log.Println("serviceendpoint created:", se6)

	// create service vertex with its edges.
	s := client.Service.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetStatus("string").
		SetType(st0).
		AddUpstream(s2).
		AddProperties(pr3).
		AddLinks(l4).
		AddCustomer(c5).
		AddEndpoints(se6).
		SaveX(ctx)
	log.Println("service created:", s)

	// query edges.
	st0, err = s.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", st0)

	s2, err = s.QueryUpstream().First(ctx)
	if err != nil {
		log.Fatalf("failed querying upstream: %v", err)
	}
	log.Println("upstream found:", s2)

	pr3, err = s.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr3)

	l4, err = s.QueryLinks().First(ctx)
	if err != nil {
		log.Fatalf("failed querying links: %v", err)
	}
	log.Println("links found:", l4)

	c5, err = s.QueryCustomer().First(ctx)
	if err != nil {
		log.Fatalf("failed querying customer: %v", err)
	}
	log.Println("customer found:", c5)

	se6, err = s.QueryEndpoints().First(ctx)
	if err != nil {
		log.Fatalf("failed querying endpoints: %v", err)
	}
	log.Println("endpoints found:", se6)

	// Output:
}
func ExampleServiceEndpoint() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the serviceendpoint's edges.
	ep0 := client.EquipmentPort.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SaveX(ctx)
	log.Println("equipmentport created:", ep0)

	// create serviceendpoint vertex with its edges.
	se := client.ServiceEndpoint.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetRole("string").
		SetPort(ep0).
		SaveX(ctx)
	log.Println("serviceendpoint created:", se)

	// query edges.
	ep0, err = se.QueryPort().First(ctx)
	if err != nil {
		log.Fatalf("failed querying port: %v", err)
	}
	log.Println("port found:", ep0)

	// Output:
}
func ExampleServiceType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the servicetype's edges.
	pt1 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt1)

	// create servicetype vertex with its edges.
	st := client.ServiceType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetHasCustomer(true).
		AddPropertyTypes(pt1).
		SaveX(ctx)
	log.Println("servicetype created:", st)

	// query edges.

	pt1, err = st.QueryPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying property_types: %v", err)
	}
	log.Println("property_types found:", pt1)

	// Output:
}
func ExampleSurvey() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the survey's edges.
	l0 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l0)
	f1 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f1)

	// create survey vertex with its edges.
	s := client.Survey.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetOwnerName("string").
		SetCompletionTimestamp(time.Now()).
		SetLocation(l0).
		SetSourceFile(f1).
		SaveX(ctx)
	log.Println("survey created:", s)

	// query edges.
	l0, err = s.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l0)

	f1, err = s.QuerySourceFile().First(ctx)
	if err != nil {
		log.Fatalf("failed querying source_file: %v", err)
	}
	log.Println("source_file found:", f1)

	// Output:
}
func ExampleSurveyCellScan() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the surveycellscan's edges.
	sq0 := client.SurveyQuestion.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFormName("string").
		SetFormDescription("string").
		SetFormIndex(1).
		SetQuestionType("string").
		SetQuestionFormat("string").
		SetQuestionText("string").
		SetQuestionIndex(1).
		SetBoolData(true).
		SetEmailData("string").
		SetLatitude(1).
		SetLongitude(1).
		SetLocationAccuracy(1).
		SetAltitude(1).
		SetPhoneData("string").
		SetTextData("string").
		SetFloatData(1).
		SetIntData(1).
		SetDateData(time.Now()).
		SaveX(ctx)
	log.Println("surveyquestion created:", sq0)
	l1 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l1)

	// create surveycellscan vertex with its edges.
	scs := client.SurveyCellScan.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetNetworkType("string").
		SetSignalStrength(1).
		SetTimestamp(time.Now()).
		SetBaseStationID("string").
		SetNetworkID("string").
		SetSystemID("string").
		SetCellID("string").
		SetLocationAreaCode("string").
		SetMobileCountryCode("string").
		SetMobileNetworkCode("string").
		SetPrimaryScramblingCode("string").
		SetOperator("string").
		SetArfcn(1).
		SetPhysicalCellID("string").
		SetTrackingAreaCode("string").
		SetTimingAdvance(1).
		SetEarfcn(1).
		SetUarfcn(1).
		SetLatitude(1).
		SetLongitude(1).
		SetSurveyQuestion(sq0).
		SetLocation(l1).
		SaveX(ctx)
	log.Println("surveycellscan created:", scs)

	// query edges.
	sq0, err = scs.QuerySurveyQuestion().First(ctx)
	if err != nil {
		log.Fatalf("failed querying survey_question: %v", err)
	}
	log.Println("survey_question found:", sq0)

	l1, err = scs.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l1)

	// Output:
}
func ExampleSurveyQuestion() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the surveyquestion's edges.
	s0 := client.Survey.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetOwnerName("string").
		SetCompletionTimestamp(time.Now()).
		SaveX(ctx)
	log.Println("survey created:", s0)
	f3 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f3)

	// create surveyquestion vertex with its edges.
	sq := client.SurveyQuestion.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFormName("string").
		SetFormDescription("string").
		SetFormIndex(1).
		SetQuestionType("string").
		SetQuestionFormat("string").
		SetQuestionText("string").
		SetQuestionIndex(1).
		SetBoolData(true).
		SetEmailData("string").
		SetLatitude(1).
		SetLongitude(1).
		SetLocationAccuracy(1).
		SetAltitude(1).
		SetPhoneData("string").
		SetTextData("string").
		SetFloatData(1).
		SetIntData(1).
		SetDateData(time.Now()).
		SetSurvey(s0).
		AddPhotoData(f3).
		SaveX(ctx)
	log.Println("surveyquestion created:", sq)

	// query edges.
	s0, err = sq.QuerySurvey().First(ctx)
	if err != nil {
		log.Fatalf("failed querying survey: %v", err)
	}
	log.Println("survey found:", s0)

	f3, err = sq.QueryPhotoData().First(ctx)
	if err != nil {
		log.Fatalf("failed querying photo_data: %v", err)
	}
	log.Println("photo_data found:", f3)

	// Output:
}
func ExampleSurveyTemplateCategory() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the surveytemplatecategory's edges.
	stq0 := client.SurveyTemplateQuestion.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetQuestionTitle("string").
		SetQuestionDescription("string").
		SetQuestionType("string").
		SetIndex(1).
		SaveX(ctx)
	log.Println("surveytemplatequestion created:", stq0)

	// create surveytemplatecategory vertex with its edges.
	stc := client.SurveyTemplateCategory.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetCategoryTitle("string").
		SetCategoryDescription("string").
		AddSurveyTemplateQuestions(stq0).
		SaveX(ctx)
	log.Println("surveytemplatecategory created:", stc)

	// query edges.
	stq0, err = stc.QuerySurveyTemplateQuestions().First(ctx)
	if err != nil {
		log.Fatalf("failed querying survey_template_questions: %v", err)
	}
	log.Println("survey_template_questions found:", stq0)

	// Output:
}
func ExampleSurveyTemplateQuestion() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the surveytemplatequestion's edges.

	// create surveytemplatequestion vertex with its edges.
	stq := client.SurveyTemplateQuestion.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetQuestionTitle("string").
		SetQuestionDescription("string").
		SetQuestionType("string").
		SetIndex(1).
		SaveX(ctx)
	log.Println("surveytemplatequestion created:", stq)

	// query edges.

	// Output:
}
func ExampleSurveyWiFiScan() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the surveywifiscan's edges.
	sq0 := client.SurveyQuestion.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetFormName("string").
		SetFormDescription("string").
		SetFormIndex(1).
		SetQuestionType("string").
		SetQuestionFormat("string").
		SetQuestionText("string").
		SetQuestionIndex(1).
		SetBoolData(true).
		SetEmailData("string").
		SetLatitude(1).
		SetLongitude(1).
		SetLocationAccuracy(1).
		SetAltitude(1).
		SetPhoneData("string").
		SetTextData("string").
		SetFloatData(1).
		SetIntData(1).
		SetDateData(time.Now()).
		SaveX(ctx)
	log.Println("surveyquestion created:", sq0)
	l1 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l1)

	// create surveywifiscan vertex with its edges.
	swfs := client.SurveyWiFiScan.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetSsid("string").
		SetBssid("string").
		SetTimestamp(time.Now()).
		SetFrequency(1).
		SetChannel(1).
		SetBand("string").
		SetChannelWidth(1).
		SetCapabilities("string").
		SetStrength(1).
		SetLatitude(1).
		SetLongitude(1).
		SetSurveyQuestion(sq0).
		SetLocation(l1).
		SaveX(ctx)
	log.Println("surveywifiscan created:", swfs)

	// query edges.
	sq0, err = swfs.QuerySurveyQuestion().First(ctx)
	if err != nil {
		log.Fatalf("failed querying survey_question: %v", err)
	}
	log.Println("survey_question found:", sq0)

	l1, err = swfs.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l1)

	// Output:
}
func ExampleTechnician() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the technician's edges.

	// create technician vertex with its edges.
	t := client.Technician.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetEmail("string").
		SaveX(ctx)
	log.Println("technician created:", t)

	// query edges.

	// Output:
}
func ExampleWorkOrder() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the workorder's edges.
	wot0 := client.WorkOrderType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		SaveX(ctx)
	log.Println("workordertype created:", wot0)
	f3 := client.File.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetSize(1).
		SetModifiedAt(time.Now()).
		SetUploadedAt(time.Now()).
		SetContentType("string").
		SetStoreKey("string").
		SetCategory("string").
		SaveX(ctx)
	log.Println("file created:", f3)
	l4 := client.Location.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetExternalID("string").
		SetLatitude(1).
		SetLongitude(1).
		SetSiteSurveyNeeded(true).
		SaveX(ctx)
	log.Println("location created:", l4)
	c5 := client.Comment.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetAuthorName("string").
		SetText("string").
		SaveX(ctx)
	log.Println("comment created:", c5)
	pr6 := client.Property.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetStringVal("string").
		SaveX(ctx)
	log.Println("property created:", pr6)
	cli7 := client.CheckListItem.
		Create().
		SetTitle("string").
		SetType("string").
		SetIndex(1).
		SetChecked(true).
		SetStringVal("string").
		SetEnumValues("string").
		SetHelpText("string").
		SaveX(ctx)
	log.Println("checklistitem created:", cli7)
	t8 := client.Technician.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetEmail("string").
		SaveX(ctx)
	log.Println("technician created:", t8)

	// create workorder vertex with its edges.
	wo := client.WorkOrder.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetStatus("string").
		SetPriority("string").
		SetDescription("string").
		SetOwnerName("string").
		SetInstallDate(time.Now()).
		SetCreationDate(time.Now()).
		SetAssignee("string").
		SetIndex(1).
		SetType(wot0).
		AddFiles(f3).
		SetLocation(l4).
		AddComments(c5).
		AddProperties(pr6).
		AddCheckListItems(cli7).
		SetTechnician(t8).
		SaveX(ctx)
	log.Println("workorder created:", wo)

	// query edges.
	wot0, err = wo.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", wot0)

	f3, err = wo.QueryFiles().First(ctx)
	if err != nil {
		log.Fatalf("failed querying files: %v", err)
	}
	log.Println("files found:", f3)

	l4, err = wo.QueryLocation().First(ctx)
	if err != nil {
		log.Fatalf("failed querying location: %v", err)
	}
	log.Println("location found:", l4)

	c5, err = wo.QueryComments().First(ctx)
	if err != nil {
		log.Fatalf("failed querying comments: %v", err)
	}
	log.Println("comments found:", c5)

	pr6, err = wo.QueryProperties().First(ctx)
	if err != nil {
		log.Fatalf("failed querying properties: %v", err)
	}
	log.Println("properties found:", pr6)

	cli7, err = wo.QueryCheckListItems().First(ctx)
	if err != nil {
		log.Fatalf("failed querying check_list_items: %v", err)
	}
	log.Println("check_list_items found:", cli7)

	t8, err = wo.QueryTechnician().First(ctx)
	if err != nil {
		log.Fatalf("failed querying technician: %v", err)
	}
	log.Println("technician found:", t8)

	// Output:
}
func ExampleWorkOrderDefinition() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the workorderdefinition's edges.
	wot0 := client.WorkOrderType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		SaveX(ctx)
	log.Println("workordertype created:", wot0)

	// create workorderdefinition vertex with its edges.
	wod := client.WorkOrderDefinition.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetIndex(1).
		SetType(wot0).
		SaveX(ctx)
	log.Println("workorderdefinition created:", wod)

	// query edges.
	wot0, err = wod.QueryType().First(ctx)
	if err != nil {
		log.Fatalf("failed querying type: %v", err)
	}
	log.Println("type found:", wot0)

	// Output:
}
func ExampleWorkOrderType() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the workordertype's edges.
	pt1 := client.PropertyType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetType("string").
		SetName("string").
		SetIndex(1).
		SetCategory("string").
		SetIntVal(1).
		SetBoolVal(true).
		SetFloatVal(1).
		SetLatitudeVal(1).
		SetLongitudeVal(1).
		SetStringVal("string").
		SetRangeFromVal(1).
		SetRangeToVal(1).
		SetIsInstanceProperty(true).
		SetEditable(true).
		SetMandatory(true).
		SetDeleted(true).
		SaveX(ctx)
	log.Println("propertytype created:", pt1)
	clid3 := client.CheckListItemDefinition.
		Create().
		SetTitle("string").
		SetType("string").
		SetIndex(1).
		SetEnumValues("string").
		SetHelpText("string").
		SaveX(ctx)
	log.Println("checklistitemdefinition created:", clid3)

	// create workordertype vertex with its edges.
	wot := client.WorkOrderType.
		Create().
		SetCreateTime(time.Now()).
		SetUpdateTime(time.Now()).
		SetName("string").
		SetDescription("string").
		AddPropertyTypes(pt1).
		AddCheckListDefinitions(clid3).
		SaveX(ctx)
	log.Println("workordertype created:", wot)

	// query edges.

	pt1, err = wot.QueryPropertyTypes().First(ctx)
	if err != nil {
		log.Fatalf("failed querying property_types: %v", err)
	}
	log.Println("property_types found:", pt1)

	clid3, err = wot.QueryCheckListDefinitions().First(ctx)
	if err != nil {
		log.Fatalf("failed querying check_list_definitions: %v", err)
	}
	log.Println("check_list_definitions found:", clid3)

	// Output:
}
