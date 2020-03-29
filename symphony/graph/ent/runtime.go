// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/schema"

	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"

	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"

	"github.com/facebookincubator/symphony/graph/ent/comment"

	"github.com/facebookincubator/symphony/graph/ent/customer"

	"github.com/facebookincubator/symphony/graph/ent/equipment"

	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"

	"github.com/facebookincubator/symphony/graph/ent/equipmentport"

	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"

	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"

	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"

	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"

	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"

	"github.com/facebookincubator/symphony/graph/ent/file"

	"github.com/facebookincubator/symphony/graph/ent/floorplan"

	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"

	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"

	"github.com/facebookincubator/symphony/graph/ent/hyperlink"

	"github.com/facebookincubator/symphony/graph/ent/link"

	"github.com/facebookincubator/symphony/graph/ent/location"

	"github.com/facebookincubator/symphony/graph/ent/locationtype"

	"github.com/facebookincubator/symphony/graph/ent/project"

	"github.com/facebookincubator/symphony/graph/ent/projecttype"

	"github.com/facebookincubator/symphony/graph/ent/property"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"

	"github.com/facebookincubator/symphony/graph/ent/reportfilter"

	"github.com/facebookincubator/symphony/graph/ent/service"

	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"

	"github.com/facebookincubator/symphony/graph/ent/servicetype"

	"github.com/facebookincubator/symphony/graph/ent/survey"

	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"

	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"

	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"

	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"

	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"

	"github.com/facebookincubator/symphony/graph/ent/technician"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/ent/usersgroup"

	"github.com/facebookincubator/symphony/graph/ent/workorder"

	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	actionsruleMixin := schema.ActionsRule{}.Mixin()
	actionsruleMixinFields := [...][]ent.Field{
		actionsruleMixin[0].Fields(),
	}
	actionsruleFields := schema.ActionsRule{}.Fields()
	_ = actionsruleFields
	// actionsruleDescCreateTime is the schema descriptor for create_time field.
	actionsruleDescCreateTime := actionsruleMixinFields[0][0].Descriptor()
	// actionsrule.DefaultCreateTime holds the default value on creation for the create_time field.
	actionsrule.DefaultCreateTime = actionsruleDescCreateTime.Default.(func() time.Time)
	// actionsruleDescUpdateTime is the schema descriptor for update_time field.
	actionsruleDescUpdateTime := actionsruleMixinFields[0][1].Descriptor()
	// actionsrule.DefaultUpdateTime holds the default value on creation for the update_time field.
	actionsrule.DefaultUpdateTime = actionsruleDescUpdateTime.Default.(func() time.Time)
	// actionsrule.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	actionsrule.UpdateDefaultUpdateTime = actionsruleDescUpdateTime.UpdateDefault.(func() time.Time)
	checklistcategoryMixin := schema.CheckListCategory{}.Mixin()
	checklistcategoryMixinFields := [...][]ent.Field{
		checklistcategoryMixin[0].Fields(),
	}
	checklistcategoryFields := schema.CheckListCategory{}.Fields()
	_ = checklistcategoryFields
	// checklistcategoryDescCreateTime is the schema descriptor for create_time field.
	checklistcategoryDescCreateTime := checklistcategoryMixinFields[0][0].Descriptor()
	// checklistcategory.DefaultCreateTime holds the default value on creation for the create_time field.
	checklistcategory.DefaultCreateTime = checklistcategoryDescCreateTime.Default.(func() time.Time)
	// checklistcategoryDescUpdateTime is the schema descriptor for update_time field.
	checklistcategoryDescUpdateTime := checklistcategoryMixinFields[0][1].Descriptor()
	// checklistcategory.DefaultUpdateTime holds the default value on creation for the update_time field.
	checklistcategory.DefaultUpdateTime = checklistcategoryDescUpdateTime.Default.(func() time.Time)
	// checklistcategory.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	checklistcategory.UpdateDefaultUpdateTime = checklistcategoryDescUpdateTime.UpdateDefault.(func() time.Time)
	checklistitemdefinitionMixin := schema.CheckListItemDefinition{}.Mixin()
	checklistitemdefinitionMixinFields := [...][]ent.Field{
		checklistitemdefinitionMixin[0].Fields(),
	}
	checklistitemdefinitionFields := schema.CheckListItemDefinition{}.Fields()
	_ = checklistitemdefinitionFields
	// checklistitemdefinitionDescCreateTime is the schema descriptor for create_time field.
	checklistitemdefinitionDescCreateTime := checklistitemdefinitionMixinFields[0][0].Descriptor()
	// checklistitemdefinition.DefaultCreateTime holds the default value on creation for the create_time field.
	checklistitemdefinition.DefaultCreateTime = checklistitemdefinitionDescCreateTime.Default.(func() time.Time)
	// checklistitemdefinitionDescUpdateTime is the schema descriptor for update_time field.
	checklistitemdefinitionDescUpdateTime := checklistitemdefinitionMixinFields[0][1].Descriptor()
	// checklistitemdefinition.DefaultUpdateTime holds the default value on creation for the update_time field.
	checklistitemdefinition.DefaultUpdateTime = checklistitemdefinitionDescUpdateTime.Default.(func() time.Time)
	// checklistitemdefinition.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	checklistitemdefinition.UpdateDefaultUpdateTime = checklistitemdefinitionDescUpdateTime.UpdateDefault.(func() time.Time)
	commentMixin := schema.Comment{}.Mixin()
	commentMixinFields := [...][]ent.Field{
		commentMixin[0].Fields(),
	}
	commentFields := schema.Comment{}.Fields()
	_ = commentFields
	// commentDescCreateTime is the schema descriptor for create_time field.
	commentDescCreateTime := commentMixinFields[0][0].Descriptor()
	// comment.DefaultCreateTime holds the default value on creation for the create_time field.
	comment.DefaultCreateTime = commentDescCreateTime.Default.(func() time.Time)
	// commentDescUpdateTime is the schema descriptor for update_time field.
	commentDescUpdateTime := commentMixinFields[0][1].Descriptor()
	// comment.DefaultUpdateTime holds the default value on creation for the update_time field.
	comment.DefaultUpdateTime = commentDescUpdateTime.Default.(func() time.Time)
	// comment.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	comment.UpdateDefaultUpdateTime = commentDescUpdateTime.UpdateDefault.(func() time.Time)
	customerMixin := schema.Customer{}.Mixin()
	customerMixinFields := [...][]ent.Field{
		customerMixin[0].Fields(),
	}
	customerFields := schema.Customer{}.Fields()
	_ = customerFields
	// customerDescCreateTime is the schema descriptor for create_time field.
	customerDescCreateTime := customerMixinFields[0][0].Descriptor()
	// customer.DefaultCreateTime holds the default value on creation for the create_time field.
	customer.DefaultCreateTime = customerDescCreateTime.Default.(func() time.Time)
	// customerDescUpdateTime is the schema descriptor for update_time field.
	customerDescUpdateTime := customerMixinFields[0][1].Descriptor()
	// customer.DefaultUpdateTime holds the default value on creation for the update_time field.
	customer.DefaultUpdateTime = customerDescUpdateTime.Default.(func() time.Time)
	// customer.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	customer.UpdateDefaultUpdateTime = customerDescUpdateTime.UpdateDefault.(func() time.Time)
	// customerDescName is the schema descriptor for name field.
	customerDescName := customerFields[0].Descriptor()
	// customer.NameValidator is a validator for the "name" field. It is called by the builders before save.
	customer.NameValidator = customerDescName.Validators[0].(func(string) error)
	// customerDescExternalID is the schema descriptor for external_id field.
	customerDescExternalID := customerFields[1].Descriptor()
	// customer.ExternalIDValidator is a validator for the "external_id" field. It is called by the builders before save.
	customer.ExternalIDValidator = customerDescExternalID.Validators[0].(func(string) error)
	equipmentMixin := schema.Equipment{}.Mixin()
	equipmentMixinFields := [...][]ent.Field{
		equipmentMixin[0].Fields(),
	}
	equipmentFields := schema.Equipment{}.Fields()
	_ = equipmentFields
	// equipmentDescCreateTime is the schema descriptor for create_time field.
	equipmentDescCreateTime := equipmentMixinFields[0][0].Descriptor()
	// equipment.DefaultCreateTime holds the default value on creation for the create_time field.
	equipment.DefaultCreateTime = equipmentDescCreateTime.Default.(func() time.Time)
	// equipmentDescUpdateTime is the schema descriptor for update_time field.
	equipmentDescUpdateTime := equipmentMixinFields[0][1].Descriptor()
	// equipment.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipment.DefaultUpdateTime = equipmentDescUpdateTime.Default.(func() time.Time)
	// equipment.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipment.UpdateDefaultUpdateTime = equipmentDescUpdateTime.UpdateDefault.(func() time.Time)
	// equipmentDescName is the schema descriptor for name field.
	equipmentDescName := equipmentFields[0].Descriptor()
	// equipment.NameValidator is a validator for the "name" field. It is called by the builders before save.
	equipment.NameValidator = equipmentDescName.Validators[0].(func(string) error)
	// equipmentDescDeviceID is the schema descriptor for device_id field.
	equipmentDescDeviceID := equipmentFields[2].Descriptor()
	// equipment.DeviceIDValidator is a validator for the "device_id" field. It is called by the builders before save.
	equipment.DeviceIDValidator = equipmentDescDeviceID.Validators[0].(func(string) error)
	equipmentcategoryMixin := schema.EquipmentCategory{}.Mixin()
	equipmentcategoryMixinFields := [...][]ent.Field{
		equipmentcategoryMixin[0].Fields(),
	}
	equipmentcategoryFields := schema.EquipmentCategory{}.Fields()
	_ = equipmentcategoryFields
	// equipmentcategoryDescCreateTime is the schema descriptor for create_time field.
	equipmentcategoryDescCreateTime := equipmentcategoryMixinFields[0][0].Descriptor()
	// equipmentcategory.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentcategory.DefaultCreateTime = equipmentcategoryDescCreateTime.Default.(func() time.Time)
	// equipmentcategoryDescUpdateTime is the schema descriptor for update_time field.
	equipmentcategoryDescUpdateTime := equipmentcategoryMixinFields[0][1].Descriptor()
	// equipmentcategory.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentcategory.DefaultUpdateTime = equipmentcategoryDescUpdateTime.Default.(func() time.Time)
	// equipmentcategory.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentcategory.UpdateDefaultUpdateTime = equipmentcategoryDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmentportMixin := schema.EquipmentPort{}.Mixin()
	equipmentportMixinFields := [...][]ent.Field{
		equipmentportMixin[0].Fields(),
	}
	equipmentportFields := schema.EquipmentPort{}.Fields()
	_ = equipmentportFields
	// equipmentportDescCreateTime is the schema descriptor for create_time field.
	equipmentportDescCreateTime := equipmentportMixinFields[0][0].Descriptor()
	// equipmentport.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentport.DefaultCreateTime = equipmentportDescCreateTime.Default.(func() time.Time)
	// equipmentportDescUpdateTime is the schema descriptor for update_time field.
	equipmentportDescUpdateTime := equipmentportMixinFields[0][1].Descriptor()
	// equipmentport.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentport.DefaultUpdateTime = equipmentportDescUpdateTime.Default.(func() time.Time)
	// equipmentport.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentport.UpdateDefaultUpdateTime = equipmentportDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmentportdefinitionMixin := schema.EquipmentPortDefinition{}.Mixin()
	equipmentportdefinitionMixinFields := [...][]ent.Field{
		equipmentportdefinitionMixin[0].Fields(),
	}
	equipmentportdefinitionFields := schema.EquipmentPortDefinition{}.Fields()
	_ = equipmentportdefinitionFields
	// equipmentportdefinitionDescCreateTime is the schema descriptor for create_time field.
	equipmentportdefinitionDescCreateTime := equipmentportdefinitionMixinFields[0][0].Descriptor()
	// equipmentportdefinition.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentportdefinition.DefaultCreateTime = equipmentportdefinitionDescCreateTime.Default.(func() time.Time)
	// equipmentportdefinitionDescUpdateTime is the schema descriptor for update_time field.
	equipmentportdefinitionDescUpdateTime := equipmentportdefinitionMixinFields[0][1].Descriptor()
	// equipmentportdefinition.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentportdefinition.DefaultUpdateTime = equipmentportdefinitionDescUpdateTime.Default.(func() time.Time)
	// equipmentportdefinition.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentportdefinition.UpdateDefaultUpdateTime = equipmentportdefinitionDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmentporttypeMixin := schema.EquipmentPortType{}.Mixin()
	equipmentporttypeMixinFields := [...][]ent.Field{
		equipmentporttypeMixin[0].Fields(),
	}
	equipmentporttypeFields := schema.EquipmentPortType{}.Fields()
	_ = equipmentporttypeFields
	// equipmentporttypeDescCreateTime is the schema descriptor for create_time field.
	equipmentporttypeDescCreateTime := equipmentporttypeMixinFields[0][0].Descriptor()
	// equipmentporttype.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentporttype.DefaultCreateTime = equipmentporttypeDescCreateTime.Default.(func() time.Time)
	// equipmentporttypeDescUpdateTime is the schema descriptor for update_time field.
	equipmentporttypeDescUpdateTime := equipmentporttypeMixinFields[0][1].Descriptor()
	// equipmentporttype.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentporttype.DefaultUpdateTime = equipmentporttypeDescUpdateTime.Default.(func() time.Time)
	// equipmentporttype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentporttype.UpdateDefaultUpdateTime = equipmentporttypeDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmentpositionMixin := schema.EquipmentPosition{}.Mixin()
	equipmentpositionMixinFields := [...][]ent.Field{
		equipmentpositionMixin[0].Fields(),
	}
	equipmentpositionFields := schema.EquipmentPosition{}.Fields()
	_ = equipmentpositionFields
	// equipmentpositionDescCreateTime is the schema descriptor for create_time field.
	equipmentpositionDescCreateTime := equipmentpositionMixinFields[0][0].Descriptor()
	// equipmentposition.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentposition.DefaultCreateTime = equipmentpositionDescCreateTime.Default.(func() time.Time)
	// equipmentpositionDescUpdateTime is the schema descriptor for update_time field.
	equipmentpositionDescUpdateTime := equipmentpositionMixinFields[0][1].Descriptor()
	// equipmentposition.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentposition.DefaultUpdateTime = equipmentpositionDescUpdateTime.Default.(func() time.Time)
	// equipmentposition.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentposition.UpdateDefaultUpdateTime = equipmentpositionDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmentpositiondefinitionMixin := schema.EquipmentPositionDefinition{}.Mixin()
	equipmentpositiondefinitionMixinFields := [...][]ent.Field{
		equipmentpositiondefinitionMixin[0].Fields(),
	}
	equipmentpositiondefinitionFields := schema.EquipmentPositionDefinition{}.Fields()
	_ = equipmentpositiondefinitionFields
	// equipmentpositiondefinitionDescCreateTime is the schema descriptor for create_time field.
	equipmentpositiondefinitionDescCreateTime := equipmentpositiondefinitionMixinFields[0][0].Descriptor()
	// equipmentpositiondefinition.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmentpositiondefinition.DefaultCreateTime = equipmentpositiondefinitionDescCreateTime.Default.(func() time.Time)
	// equipmentpositiondefinitionDescUpdateTime is the schema descriptor for update_time field.
	equipmentpositiondefinitionDescUpdateTime := equipmentpositiondefinitionMixinFields[0][1].Descriptor()
	// equipmentpositiondefinition.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmentpositiondefinition.DefaultUpdateTime = equipmentpositiondefinitionDescUpdateTime.Default.(func() time.Time)
	// equipmentpositiondefinition.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmentpositiondefinition.UpdateDefaultUpdateTime = equipmentpositiondefinitionDescUpdateTime.UpdateDefault.(func() time.Time)
	equipmenttypeMixin := schema.EquipmentType{}.Mixin()
	equipmenttypeMixinFields := [...][]ent.Field{
		equipmenttypeMixin[0].Fields(),
	}
	equipmenttypeFields := schema.EquipmentType{}.Fields()
	_ = equipmenttypeFields
	// equipmenttypeDescCreateTime is the schema descriptor for create_time field.
	equipmenttypeDescCreateTime := equipmenttypeMixinFields[0][0].Descriptor()
	// equipmenttype.DefaultCreateTime holds the default value on creation for the create_time field.
	equipmenttype.DefaultCreateTime = equipmenttypeDescCreateTime.Default.(func() time.Time)
	// equipmenttypeDescUpdateTime is the schema descriptor for update_time field.
	equipmenttypeDescUpdateTime := equipmenttypeMixinFields[0][1].Descriptor()
	// equipmenttype.DefaultUpdateTime holds the default value on creation for the update_time field.
	equipmenttype.DefaultUpdateTime = equipmenttypeDescUpdateTime.Default.(func() time.Time)
	// equipmenttype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	equipmenttype.UpdateDefaultUpdateTime = equipmenttypeDescUpdateTime.UpdateDefault.(func() time.Time)
	fileMixin := schema.File{}.Mixin()
	fileMixinFields := [...][]ent.Field{
		fileMixin[0].Fields(),
	}
	fileFields := schema.File{}.Fields()
	_ = fileFields
	// fileDescCreateTime is the schema descriptor for create_time field.
	fileDescCreateTime := fileMixinFields[0][0].Descriptor()
	// file.DefaultCreateTime holds the default value on creation for the create_time field.
	file.DefaultCreateTime = fileDescCreateTime.Default.(func() time.Time)
	// fileDescUpdateTime is the schema descriptor for update_time field.
	fileDescUpdateTime := fileMixinFields[0][1].Descriptor()
	// file.DefaultUpdateTime holds the default value on creation for the update_time field.
	file.DefaultUpdateTime = fileDescUpdateTime.Default.(func() time.Time)
	// file.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	file.UpdateDefaultUpdateTime = fileDescUpdateTime.UpdateDefault.(func() time.Time)
	// fileDescSize is the schema descriptor for size field.
	fileDescSize := fileFields[2].Descriptor()
	// file.SizeValidator is a validator for the "size" field. It is called by the builders before save.
	file.SizeValidator = fileDescSize.Validators[0].(func(int) error)
	floorplanMixin := schema.FloorPlan{}.Mixin()
	floorplanMixinFields := [...][]ent.Field{
		floorplanMixin[0].Fields(),
	}
	floorplanFields := schema.FloorPlan{}.Fields()
	_ = floorplanFields
	// floorplanDescCreateTime is the schema descriptor for create_time field.
	floorplanDescCreateTime := floorplanMixinFields[0][0].Descriptor()
	// floorplan.DefaultCreateTime holds the default value on creation for the create_time field.
	floorplan.DefaultCreateTime = floorplanDescCreateTime.Default.(func() time.Time)
	// floorplanDescUpdateTime is the schema descriptor for update_time field.
	floorplanDescUpdateTime := floorplanMixinFields[0][1].Descriptor()
	// floorplan.DefaultUpdateTime holds the default value on creation for the update_time field.
	floorplan.DefaultUpdateTime = floorplanDescUpdateTime.Default.(func() time.Time)
	// floorplan.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	floorplan.UpdateDefaultUpdateTime = floorplanDescUpdateTime.UpdateDefault.(func() time.Time)
	floorplanreferencepointMixin := schema.FloorPlanReferencePoint{}.Mixin()
	floorplanreferencepointMixinFields := [...][]ent.Field{
		floorplanreferencepointMixin[0].Fields(),
	}
	floorplanreferencepointFields := schema.FloorPlanReferencePoint{}.Fields()
	_ = floorplanreferencepointFields
	// floorplanreferencepointDescCreateTime is the schema descriptor for create_time field.
	floorplanreferencepointDescCreateTime := floorplanreferencepointMixinFields[0][0].Descriptor()
	// floorplanreferencepoint.DefaultCreateTime holds the default value on creation for the create_time field.
	floorplanreferencepoint.DefaultCreateTime = floorplanreferencepointDescCreateTime.Default.(func() time.Time)
	// floorplanreferencepointDescUpdateTime is the schema descriptor for update_time field.
	floorplanreferencepointDescUpdateTime := floorplanreferencepointMixinFields[0][1].Descriptor()
	// floorplanreferencepoint.DefaultUpdateTime holds the default value on creation for the update_time field.
	floorplanreferencepoint.DefaultUpdateTime = floorplanreferencepointDescUpdateTime.Default.(func() time.Time)
	// floorplanreferencepoint.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	floorplanreferencepoint.UpdateDefaultUpdateTime = floorplanreferencepointDescUpdateTime.UpdateDefault.(func() time.Time)
	floorplanscaleMixin := schema.FloorPlanScale{}.Mixin()
	floorplanscaleMixinFields := [...][]ent.Field{
		floorplanscaleMixin[0].Fields(),
	}
	floorplanscaleFields := schema.FloorPlanScale{}.Fields()
	_ = floorplanscaleFields
	// floorplanscaleDescCreateTime is the schema descriptor for create_time field.
	floorplanscaleDescCreateTime := floorplanscaleMixinFields[0][0].Descriptor()
	// floorplanscale.DefaultCreateTime holds the default value on creation for the create_time field.
	floorplanscale.DefaultCreateTime = floorplanscaleDescCreateTime.Default.(func() time.Time)
	// floorplanscaleDescUpdateTime is the schema descriptor for update_time field.
	floorplanscaleDescUpdateTime := floorplanscaleMixinFields[0][1].Descriptor()
	// floorplanscale.DefaultUpdateTime holds the default value on creation for the update_time field.
	floorplanscale.DefaultUpdateTime = floorplanscaleDescUpdateTime.Default.(func() time.Time)
	// floorplanscale.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	floorplanscale.UpdateDefaultUpdateTime = floorplanscaleDescUpdateTime.UpdateDefault.(func() time.Time)
	hyperlinkMixin := schema.Hyperlink{}.Mixin()
	hyperlinkMixinFields := [...][]ent.Field{
		hyperlinkMixin[0].Fields(),
	}
	hyperlinkFields := schema.Hyperlink{}.Fields()
	_ = hyperlinkFields
	// hyperlinkDescCreateTime is the schema descriptor for create_time field.
	hyperlinkDescCreateTime := hyperlinkMixinFields[0][0].Descriptor()
	// hyperlink.DefaultCreateTime holds the default value on creation for the create_time field.
	hyperlink.DefaultCreateTime = hyperlinkDescCreateTime.Default.(func() time.Time)
	// hyperlinkDescUpdateTime is the schema descriptor for update_time field.
	hyperlinkDescUpdateTime := hyperlinkMixinFields[0][1].Descriptor()
	// hyperlink.DefaultUpdateTime holds the default value on creation for the update_time field.
	hyperlink.DefaultUpdateTime = hyperlinkDescUpdateTime.Default.(func() time.Time)
	// hyperlink.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	hyperlink.UpdateDefaultUpdateTime = hyperlinkDescUpdateTime.UpdateDefault.(func() time.Time)
	linkMixin := schema.Link{}.Mixin()
	linkMixinFields := [...][]ent.Field{
		linkMixin[0].Fields(),
	}
	linkFields := schema.Link{}.Fields()
	_ = linkFields
	// linkDescCreateTime is the schema descriptor for create_time field.
	linkDescCreateTime := linkMixinFields[0][0].Descriptor()
	// link.DefaultCreateTime holds the default value on creation for the create_time field.
	link.DefaultCreateTime = linkDescCreateTime.Default.(func() time.Time)
	// linkDescUpdateTime is the schema descriptor for update_time field.
	linkDescUpdateTime := linkMixinFields[0][1].Descriptor()
	// link.DefaultUpdateTime holds the default value on creation for the update_time field.
	link.DefaultUpdateTime = linkDescUpdateTime.Default.(func() time.Time)
	// link.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	link.UpdateDefaultUpdateTime = linkDescUpdateTime.UpdateDefault.(func() time.Time)
	locationMixin := schema.Location{}.Mixin()
	locationMixinFields := [...][]ent.Field{
		locationMixin[0].Fields(),
	}
	locationFields := schema.Location{}.Fields()
	_ = locationFields
	// locationDescCreateTime is the schema descriptor for create_time field.
	locationDescCreateTime := locationMixinFields[0][0].Descriptor()
	// location.DefaultCreateTime holds the default value on creation for the create_time field.
	location.DefaultCreateTime = locationDescCreateTime.Default.(func() time.Time)
	// locationDescUpdateTime is the schema descriptor for update_time field.
	locationDescUpdateTime := locationMixinFields[0][1].Descriptor()
	// location.DefaultUpdateTime holds the default value on creation for the update_time field.
	location.DefaultUpdateTime = locationDescUpdateTime.Default.(func() time.Time)
	// location.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	location.UpdateDefaultUpdateTime = locationDescUpdateTime.UpdateDefault.(func() time.Time)
	// locationDescName is the schema descriptor for name field.
	locationDescName := locationFields[0].Descriptor()
	// location.NameValidator is a validator for the "name" field. It is called by the builders before save.
	location.NameValidator = locationDescName.Validators[0].(func(string) error)
	// locationDescLatitude is the schema descriptor for latitude field.
	locationDescLatitude := locationFields[2].Descriptor()
	// location.DefaultLatitude holds the default value on creation for the latitude field.
	location.DefaultLatitude = locationDescLatitude.Default.(float64)
	// location.LatitudeValidator is a validator for the "latitude" field. It is called by the builders before save.
	location.LatitudeValidator = locationDescLatitude.Validators[0].(func(float64) error)
	// locationDescLongitude is the schema descriptor for longitude field.
	locationDescLongitude := locationFields[3].Descriptor()
	// location.DefaultLongitude holds the default value on creation for the longitude field.
	location.DefaultLongitude = locationDescLongitude.Default.(float64)
	// location.LongitudeValidator is a validator for the "longitude" field. It is called by the builders before save.
	location.LongitudeValidator = locationDescLongitude.Validators[0].(func(float64) error)
	// locationDescSiteSurveyNeeded is the schema descriptor for site_survey_needed field.
	locationDescSiteSurveyNeeded := locationFields[4].Descriptor()
	// location.DefaultSiteSurveyNeeded holds the default value on creation for the site_survey_needed field.
	location.DefaultSiteSurveyNeeded = locationDescSiteSurveyNeeded.Default.(bool)
	locationtypeMixin := schema.LocationType{}.Mixin()
	locationtypeMixinFields := [...][]ent.Field{
		locationtypeMixin[0].Fields(),
	}
	locationtypeFields := schema.LocationType{}.Fields()
	_ = locationtypeFields
	// locationtypeDescCreateTime is the schema descriptor for create_time field.
	locationtypeDescCreateTime := locationtypeMixinFields[0][0].Descriptor()
	// locationtype.DefaultCreateTime holds the default value on creation for the create_time field.
	locationtype.DefaultCreateTime = locationtypeDescCreateTime.Default.(func() time.Time)
	// locationtypeDescUpdateTime is the schema descriptor for update_time field.
	locationtypeDescUpdateTime := locationtypeMixinFields[0][1].Descriptor()
	// locationtype.DefaultUpdateTime holds the default value on creation for the update_time field.
	locationtype.DefaultUpdateTime = locationtypeDescUpdateTime.Default.(func() time.Time)
	// locationtype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	locationtype.UpdateDefaultUpdateTime = locationtypeDescUpdateTime.UpdateDefault.(func() time.Time)
	// locationtypeDescSite is the schema descriptor for site field.
	locationtypeDescSite := locationtypeFields[0].Descriptor()
	// locationtype.DefaultSite holds the default value on creation for the site field.
	locationtype.DefaultSite = locationtypeDescSite.Default.(bool)
	// locationtypeDescMapZoomLevel is the schema descriptor for map_zoom_level field.
	locationtypeDescMapZoomLevel := locationtypeFields[3].Descriptor()
	// locationtype.DefaultMapZoomLevel holds the default value on creation for the map_zoom_level field.
	locationtype.DefaultMapZoomLevel = locationtypeDescMapZoomLevel.Default.(int)
	// locationtypeDescIndex is the schema descriptor for index field.
	locationtypeDescIndex := locationtypeFields[4].Descriptor()
	// locationtype.DefaultIndex holds the default value on creation for the index field.
	locationtype.DefaultIndex = locationtypeDescIndex.Default.(int)
	projectMixin := schema.Project{}.Mixin()
	projectMixinFields := [...][]ent.Field{
		projectMixin[0].Fields(),
	}
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescCreateTime is the schema descriptor for create_time field.
	projectDescCreateTime := projectMixinFields[0][0].Descriptor()
	// project.DefaultCreateTime holds the default value on creation for the create_time field.
	project.DefaultCreateTime = projectDescCreateTime.Default.(func() time.Time)
	// projectDescUpdateTime is the schema descriptor for update_time field.
	projectDescUpdateTime := projectMixinFields[0][1].Descriptor()
	// project.DefaultUpdateTime holds the default value on creation for the update_time field.
	project.DefaultUpdateTime = projectDescUpdateTime.Default.(func() time.Time)
	// project.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	project.UpdateDefaultUpdateTime = projectDescUpdateTime.UpdateDefault.(func() time.Time)
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectFields[0].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	projecttypeMixin := schema.ProjectType{}.Mixin()
	projecttypeMixinFields := [...][]ent.Field{
		projecttypeMixin[0].Fields(),
	}
	projecttypeFields := schema.ProjectType{}.Fields()
	_ = projecttypeFields
	// projecttypeDescCreateTime is the schema descriptor for create_time field.
	projecttypeDescCreateTime := projecttypeMixinFields[0][0].Descriptor()
	// projecttype.DefaultCreateTime holds the default value on creation for the create_time field.
	projecttype.DefaultCreateTime = projecttypeDescCreateTime.Default.(func() time.Time)
	// projecttypeDescUpdateTime is the schema descriptor for update_time field.
	projecttypeDescUpdateTime := projecttypeMixinFields[0][1].Descriptor()
	// projecttype.DefaultUpdateTime holds the default value on creation for the update_time field.
	projecttype.DefaultUpdateTime = projecttypeDescUpdateTime.Default.(func() time.Time)
	// projecttype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	projecttype.UpdateDefaultUpdateTime = projecttypeDescUpdateTime.UpdateDefault.(func() time.Time)
	// projecttypeDescName is the schema descriptor for name field.
	projecttypeDescName := projecttypeFields[0].Descriptor()
	// projecttype.NameValidator is a validator for the "name" field. It is called by the builders before save.
	projecttype.NameValidator = projecttypeDescName.Validators[0].(func(string) error)
	propertyMixin := schema.Property{}.Mixin()
	propertyMixinFields := [...][]ent.Field{
		propertyMixin[0].Fields(),
	}
	propertyFields := schema.Property{}.Fields()
	_ = propertyFields
	// propertyDescCreateTime is the schema descriptor for create_time field.
	propertyDescCreateTime := propertyMixinFields[0][0].Descriptor()
	// property.DefaultCreateTime holds the default value on creation for the create_time field.
	property.DefaultCreateTime = propertyDescCreateTime.Default.(func() time.Time)
	// propertyDescUpdateTime is the schema descriptor for update_time field.
	propertyDescUpdateTime := propertyMixinFields[0][1].Descriptor()
	// property.DefaultUpdateTime holds the default value on creation for the update_time field.
	property.DefaultUpdateTime = propertyDescUpdateTime.Default.(func() time.Time)
	// property.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	property.UpdateDefaultUpdateTime = propertyDescUpdateTime.UpdateDefault.(func() time.Time)
	propertytypeMixin := schema.PropertyType{}.Mixin()
	propertytypeMixinFields := [...][]ent.Field{
		propertytypeMixin[0].Fields(),
	}
	propertytypeFields := schema.PropertyType{}.Fields()
	_ = propertytypeFields
	// propertytypeDescCreateTime is the schema descriptor for create_time field.
	propertytypeDescCreateTime := propertytypeMixinFields[0][0].Descriptor()
	// propertytype.DefaultCreateTime holds the default value on creation for the create_time field.
	propertytype.DefaultCreateTime = propertytypeDescCreateTime.Default.(func() time.Time)
	// propertytypeDescUpdateTime is the schema descriptor for update_time field.
	propertytypeDescUpdateTime := propertytypeMixinFields[0][1].Descriptor()
	// propertytype.DefaultUpdateTime holds the default value on creation for the update_time field.
	propertytype.DefaultUpdateTime = propertytypeDescUpdateTime.Default.(func() time.Time)
	// propertytype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	propertytype.UpdateDefaultUpdateTime = propertytypeDescUpdateTime.UpdateDefault.(func() time.Time)
	// propertytypeDescIsInstanceProperty is the schema descriptor for is_instance_property field.
	propertytypeDescIsInstanceProperty := propertytypeFields[12].Descriptor()
	// propertytype.DefaultIsInstanceProperty holds the default value on creation for the is_instance_property field.
	propertytype.DefaultIsInstanceProperty = propertytypeDescIsInstanceProperty.Default.(bool)
	// propertytypeDescEditable is the schema descriptor for editable field.
	propertytypeDescEditable := propertytypeFields[13].Descriptor()
	// propertytype.DefaultEditable holds the default value on creation for the editable field.
	propertytype.DefaultEditable = propertytypeDescEditable.Default.(bool)
	// propertytypeDescMandatory is the schema descriptor for mandatory field.
	propertytypeDescMandatory := propertytypeFields[14].Descriptor()
	// propertytype.DefaultMandatory holds the default value on creation for the mandatory field.
	propertytype.DefaultMandatory = propertytypeDescMandatory.Default.(bool)
	// propertytypeDescDeleted is the schema descriptor for deleted field.
	propertytypeDescDeleted := propertytypeFields[15].Descriptor()
	// propertytype.DefaultDeleted holds the default value on creation for the deleted field.
	propertytype.DefaultDeleted = propertytypeDescDeleted.Default.(bool)
	reportfilterMixin := schema.ReportFilter{}.Mixin()
	reportfilterMixinFields := [...][]ent.Field{
		reportfilterMixin[0].Fields(),
	}
	reportfilterFields := schema.ReportFilter{}.Fields()
	_ = reportfilterFields
	// reportfilterDescCreateTime is the schema descriptor for create_time field.
	reportfilterDescCreateTime := reportfilterMixinFields[0][0].Descriptor()
	// reportfilter.DefaultCreateTime holds the default value on creation for the create_time field.
	reportfilter.DefaultCreateTime = reportfilterDescCreateTime.Default.(func() time.Time)
	// reportfilterDescUpdateTime is the schema descriptor for update_time field.
	reportfilterDescUpdateTime := reportfilterMixinFields[0][1].Descriptor()
	// reportfilter.DefaultUpdateTime holds the default value on creation for the update_time field.
	reportfilter.DefaultUpdateTime = reportfilterDescUpdateTime.Default.(func() time.Time)
	// reportfilter.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	reportfilter.UpdateDefaultUpdateTime = reportfilterDescUpdateTime.UpdateDefault.(func() time.Time)
	// reportfilterDescName is the schema descriptor for name field.
	reportfilterDescName := reportfilterFields[0].Descriptor()
	// reportfilter.NameValidator is a validator for the "name" field. It is called by the builders before save.
	reportfilter.NameValidator = reportfilterDescName.Validators[0].(func(string) error)
	// reportfilterDescFilters is the schema descriptor for filters field.
	reportfilterDescFilters := reportfilterFields[2].Descriptor()
	// reportfilter.DefaultFilters holds the default value on creation for the filters field.
	reportfilter.DefaultFilters = reportfilterDescFilters.Default.(string)
	serviceMixin := schema.Service{}.Mixin()
	serviceMixinFields := [...][]ent.Field{
		serviceMixin[0].Fields(),
	}
	serviceFields := schema.Service{}.Fields()
	_ = serviceFields
	// serviceDescCreateTime is the schema descriptor for create_time field.
	serviceDescCreateTime := serviceMixinFields[0][0].Descriptor()
	// service.DefaultCreateTime holds the default value on creation for the create_time field.
	service.DefaultCreateTime = serviceDescCreateTime.Default.(func() time.Time)
	// serviceDescUpdateTime is the schema descriptor for update_time field.
	serviceDescUpdateTime := serviceMixinFields[0][1].Descriptor()
	// service.DefaultUpdateTime holds the default value on creation for the update_time field.
	service.DefaultUpdateTime = serviceDescUpdateTime.Default.(func() time.Time)
	// service.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	service.UpdateDefaultUpdateTime = serviceDescUpdateTime.UpdateDefault.(func() time.Time)
	// serviceDescName is the schema descriptor for name field.
	serviceDescName := serviceFields[0].Descriptor()
	// service.NameValidator is a validator for the "name" field. It is called by the builders before save.
	service.NameValidator = serviceDescName.Validators[0].(func(string) error)
	// serviceDescExternalID is the schema descriptor for external_id field.
	serviceDescExternalID := serviceFields[1].Descriptor()
	// service.ExternalIDValidator is a validator for the "external_id" field. It is called by the builders before save.
	service.ExternalIDValidator = serviceDescExternalID.Validators[0].(func(string) error)
	serviceendpointMixin := schema.ServiceEndpoint{}.Mixin()
	serviceendpointMixinFields := [...][]ent.Field{
		serviceendpointMixin[0].Fields(),
	}
	serviceendpointFields := schema.ServiceEndpoint{}.Fields()
	_ = serviceendpointFields
	// serviceendpointDescCreateTime is the schema descriptor for create_time field.
	serviceendpointDescCreateTime := serviceendpointMixinFields[0][0].Descriptor()
	// serviceendpoint.DefaultCreateTime holds the default value on creation for the create_time field.
	serviceendpoint.DefaultCreateTime = serviceendpointDescCreateTime.Default.(func() time.Time)
	// serviceendpointDescUpdateTime is the schema descriptor for update_time field.
	serviceendpointDescUpdateTime := serviceendpointMixinFields[0][1].Descriptor()
	// serviceendpoint.DefaultUpdateTime holds the default value on creation for the update_time field.
	serviceendpoint.DefaultUpdateTime = serviceendpointDescUpdateTime.Default.(func() time.Time)
	// serviceendpoint.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	serviceendpoint.UpdateDefaultUpdateTime = serviceendpointDescUpdateTime.UpdateDefault.(func() time.Time)
	servicetypeMixin := schema.ServiceType{}.Mixin()
	servicetypeMixinFields := [...][]ent.Field{
		servicetypeMixin[0].Fields(),
	}
	servicetypeFields := schema.ServiceType{}.Fields()
	_ = servicetypeFields
	// servicetypeDescCreateTime is the schema descriptor for create_time field.
	servicetypeDescCreateTime := servicetypeMixinFields[0][0].Descriptor()
	// servicetype.DefaultCreateTime holds the default value on creation for the create_time field.
	servicetype.DefaultCreateTime = servicetypeDescCreateTime.Default.(func() time.Time)
	// servicetypeDescUpdateTime is the schema descriptor for update_time field.
	servicetypeDescUpdateTime := servicetypeMixinFields[0][1].Descriptor()
	// servicetype.DefaultUpdateTime holds the default value on creation for the update_time field.
	servicetype.DefaultUpdateTime = servicetypeDescUpdateTime.Default.(func() time.Time)
	// servicetype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	servicetype.UpdateDefaultUpdateTime = servicetypeDescUpdateTime.UpdateDefault.(func() time.Time)
	// servicetypeDescHasCustomer is the schema descriptor for has_customer field.
	servicetypeDescHasCustomer := servicetypeFields[1].Descriptor()
	// servicetype.DefaultHasCustomer holds the default value on creation for the has_customer field.
	servicetype.DefaultHasCustomer = servicetypeDescHasCustomer.Default.(bool)
	surveyMixin := schema.Survey{}.Mixin()
	surveyMixinFields := [...][]ent.Field{
		surveyMixin[0].Fields(),
	}
	surveyFields := schema.Survey{}.Fields()
	_ = surveyFields
	// surveyDescCreateTime is the schema descriptor for create_time field.
	surveyDescCreateTime := surveyMixinFields[0][0].Descriptor()
	// survey.DefaultCreateTime holds the default value on creation for the create_time field.
	survey.DefaultCreateTime = surveyDescCreateTime.Default.(func() time.Time)
	// surveyDescUpdateTime is the schema descriptor for update_time field.
	surveyDescUpdateTime := surveyMixinFields[0][1].Descriptor()
	// survey.DefaultUpdateTime holds the default value on creation for the update_time field.
	survey.DefaultUpdateTime = surveyDescUpdateTime.Default.(func() time.Time)
	// survey.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	survey.UpdateDefaultUpdateTime = surveyDescUpdateTime.UpdateDefault.(func() time.Time)
	surveycellscanMixin := schema.SurveyCellScan{}.Mixin()
	surveycellscanMixinFields := [...][]ent.Field{
		surveycellscanMixin[0].Fields(),
	}
	surveycellscanFields := schema.SurveyCellScan{}.Fields()
	_ = surveycellscanFields
	// surveycellscanDescCreateTime is the schema descriptor for create_time field.
	surveycellscanDescCreateTime := surveycellscanMixinFields[0][0].Descriptor()
	// surveycellscan.DefaultCreateTime holds the default value on creation for the create_time field.
	surveycellscan.DefaultCreateTime = surveycellscanDescCreateTime.Default.(func() time.Time)
	// surveycellscanDescUpdateTime is the schema descriptor for update_time field.
	surveycellscanDescUpdateTime := surveycellscanMixinFields[0][1].Descriptor()
	// surveycellscan.DefaultUpdateTime holds the default value on creation for the update_time field.
	surveycellscan.DefaultUpdateTime = surveycellscanDescUpdateTime.Default.(func() time.Time)
	// surveycellscan.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	surveycellscan.UpdateDefaultUpdateTime = surveycellscanDescUpdateTime.UpdateDefault.(func() time.Time)
	surveyquestionMixin := schema.SurveyQuestion{}.Mixin()
	surveyquestionMixinFields := [...][]ent.Field{
		surveyquestionMixin[0].Fields(),
	}
	surveyquestionFields := schema.SurveyQuestion{}.Fields()
	_ = surveyquestionFields
	// surveyquestionDescCreateTime is the schema descriptor for create_time field.
	surveyquestionDescCreateTime := surveyquestionMixinFields[0][0].Descriptor()
	// surveyquestion.DefaultCreateTime holds the default value on creation for the create_time field.
	surveyquestion.DefaultCreateTime = surveyquestionDescCreateTime.Default.(func() time.Time)
	// surveyquestionDescUpdateTime is the schema descriptor for update_time field.
	surveyquestionDescUpdateTime := surveyquestionMixinFields[0][1].Descriptor()
	// surveyquestion.DefaultUpdateTime holds the default value on creation for the update_time field.
	surveyquestion.DefaultUpdateTime = surveyquestionDescUpdateTime.Default.(func() time.Time)
	// surveyquestion.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	surveyquestion.UpdateDefaultUpdateTime = surveyquestionDescUpdateTime.UpdateDefault.(func() time.Time)
	surveytemplatecategoryMixin := schema.SurveyTemplateCategory{}.Mixin()
	surveytemplatecategoryMixinFields := [...][]ent.Field{
		surveytemplatecategoryMixin[0].Fields(),
	}
	surveytemplatecategoryFields := schema.SurveyTemplateCategory{}.Fields()
	_ = surveytemplatecategoryFields
	// surveytemplatecategoryDescCreateTime is the schema descriptor for create_time field.
	surveytemplatecategoryDescCreateTime := surveytemplatecategoryMixinFields[0][0].Descriptor()
	// surveytemplatecategory.DefaultCreateTime holds the default value on creation for the create_time field.
	surveytemplatecategory.DefaultCreateTime = surveytemplatecategoryDescCreateTime.Default.(func() time.Time)
	// surveytemplatecategoryDescUpdateTime is the schema descriptor for update_time field.
	surveytemplatecategoryDescUpdateTime := surveytemplatecategoryMixinFields[0][1].Descriptor()
	// surveytemplatecategory.DefaultUpdateTime holds the default value on creation for the update_time field.
	surveytemplatecategory.DefaultUpdateTime = surveytemplatecategoryDescUpdateTime.Default.(func() time.Time)
	// surveytemplatecategory.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	surveytemplatecategory.UpdateDefaultUpdateTime = surveytemplatecategoryDescUpdateTime.UpdateDefault.(func() time.Time)
	surveytemplatequestionMixin := schema.SurveyTemplateQuestion{}.Mixin()
	surveytemplatequestionMixinFields := [...][]ent.Field{
		surveytemplatequestionMixin[0].Fields(),
	}
	surveytemplatequestionFields := schema.SurveyTemplateQuestion{}.Fields()
	_ = surveytemplatequestionFields
	// surveytemplatequestionDescCreateTime is the schema descriptor for create_time field.
	surveytemplatequestionDescCreateTime := surveytemplatequestionMixinFields[0][0].Descriptor()
	// surveytemplatequestion.DefaultCreateTime holds the default value on creation for the create_time field.
	surveytemplatequestion.DefaultCreateTime = surveytemplatequestionDescCreateTime.Default.(func() time.Time)
	// surveytemplatequestionDescUpdateTime is the schema descriptor for update_time field.
	surveytemplatequestionDescUpdateTime := surveytemplatequestionMixinFields[0][1].Descriptor()
	// surveytemplatequestion.DefaultUpdateTime holds the default value on creation for the update_time field.
	surveytemplatequestion.DefaultUpdateTime = surveytemplatequestionDescUpdateTime.Default.(func() time.Time)
	// surveytemplatequestion.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	surveytemplatequestion.UpdateDefaultUpdateTime = surveytemplatequestionDescUpdateTime.UpdateDefault.(func() time.Time)
	surveywifiscanMixin := schema.SurveyWiFiScan{}.Mixin()
	surveywifiscanMixinFields := [...][]ent.Field{
		surveywifiscanMixin[0].Fields(),
	}
	surveywifiscanFields := schema.SurveyWiFiScan{}.Fields()
	_ = surveywifiscanFields
	// surveywifiscanDescCreateTime is the schema descriptor for create_time field.
	surveywifiscanDescCreateTime := surveywifiscanMixinFields[0][0].Descriptor()
	// surveywifiscan.DefaultCreateTime holds the default value on creation for the create_time field.
	surveywifiscan.DefaultCreateTime = surveywifiscanDescCreateTime.Default.(func() time.Time)
	// surveywifiscanDescUpdateTime is the schema descriptor for update_time field.
	surveywifiscanDescUpdateTime := surveywifiscanMixinFields[0][1].Descriptor()
	// surveywifiscan.DefaultUpdateTime holds the default value on creation for the update_time field.
	surveywifiscan.DefaultUpdateTime = surveywifiscanDescUpdateTime.Default.(func() time.Time)
	// surveywifiscan.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	surveywifiscan.UpdateDefaultUpdateTime = surveywifiscanDescUpdateTime.UpdateDefault.(func() time.Time)
	technicianMixin := schema.Technician{}.Mixin()
	technicianMixinFields := [...][]ent.Field{
		technicianMixin[0].Fields(),
	}
	technicianFields := schema.Technician{}.Fields()
	_ = technicianFields
	// technicianDescCreateTime is the schema descriptor for create_time field.
	technicianDescCreateTime := technicianMixinFields[0][0].Descriptor()
	// technician.DefaultCreateTime holds the default value on creation for the create_time field.
	technician.DefaultCreateTime = technicianDescCreateTime.Default.(func() time.Time)
	// technicianDescUpdateTime is the schema descriptor for update_time field.
	technicianDescUpdateTime := technicianMixinFields[0][1].Descriptor()
	// technician.DefaultUpdateTime holds the default value on creation for the update_time field.
	technician.DefaultUpdateTime = technicianDescUpdateTime.Default.(func() time.Time)
	// technician.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	technician.UpdateDefaultUpdateTime = technicianDescUpdateTime.UpdateDefault.(func() time.Time)
	// technicianDescName is the schema descriptor for name field.
	technicianDescName := technicianFields[0].Descriptor()
	// technician.NameValidator is a validator for the "name" field. It is called by the builders before save.
	technician.NameValidator = technicianDescName.Validators[0].(func(string) error)
	// technicianDescEmail is the schema descriptor for email field.
	technicianDescEmail := technicianFields[1].Descriptor()
	// technician.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	technician.EmailValidator = technicianDescEmail.Validators[0].(func(string) error)
	userMixin := schema.User{}.Mixin()
	userMixinFields := [...][]ent.Field{
		userMixin[0].Fields(),
	}
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreateTime is the schema descriptor for create_time field.
	userDescCreateTime := userMixinFields[0][0].Descriptor()
	// user.DefaultCreateTime holds the default value on creation for the create_time field.
	user.DefaultCreateTime = userDescCreateTime.Default.(func() time.Time)
	// userDescUpdateTime is the schema descriptor for update_time field.
	userDescUpdateTime := userMixinFields[0][1].Descriptor()
	// user.DefaultUpdateTime holds the default value on creation for the update_time field.
	user.DefaultUpdateTime = userDescUpdateTime.Default.(func() time.Time)
	// user.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	user.UpdateDefaultUpdateTime = userDescUpdateTime.UpdateDefault.(func() time.Time)
	// userDescAuthID is the schema descriptor for auth_id field.
	userDescAuthID := userFields[0].Descriptor()
	// user.AuthIDValidator is a validator for the "auth_id" field. It is called by the builders before save.
	user.AuthIDValidator = userDescAuthID.Validators[0].(func(string) error)
	// userDescFirstName is the schema descriptor for first_name field.
	userDescFirstName := userFields[1].Descriptor()
	// user.FirstNameValidator is a validator for the "first_name" field. It is called by the builders before save.
	user.FirstNameValidator = userDescFirstName.Validators[0].(func(string) error)
	// userDescLastName is the schema descriptor for last_name field.
	userDescLastName := userFields[2].Descriptor()
	// user.LastNameValidator is a validator for the "last_name" field. It is called by the builders before save.
	user.LastNameValidator = userDescLastName.Validators[0].(func(string) error)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[3].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	usersgroupMixin := schema.UsersGroup{}.Mixin()
	usersgroupMixinFields := [...][]ent.Field{
		usersgroupMixin[0].Fields(),
	}
	usersgroupFields := schema.UsersGroup{}.Fields()
	_ = usersgroupFields
	// usersgroupDescCreateTime is the schema descriptor for create_time field.
	usersgroupDescCreateTime := usersgroupMixinFields[0][0].Descriptor()
	// usersgroup.DefaultCreateTime holds the default value on creation for the create_time field.
	usersgroup.DefaultCreateTime = usersgroupDescCreateTime.Default.(func() time.Time)
	// usersgroupDescUpdateTime is the schema descriptor for update_time field.
	usersgroupDescUpdateTime := usersgroupMixinFields[0][1].Descriptor()
	// usersgroup.DefaultUpdateTime holds the default value on creation for the update_time field.
	usersgroup.DefaultUpdateTime = usersgroupDescUpdateTime.Default.(func() time.Time)
	// usersgroup.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	usersgroup.UpdateDefaultUpdateTime = usersgroupDescUpdateTime.UpdateDefault.(func() time.Time)
	// usersgroupDescName is the schema descriptor for name field.
	usersgroupDescName := usersgroupFields[0].Descriptor()
	// usersgroup.NameValidator is a validator for the "name" field. It is called by the builders before save.
	usersgroup.NameValidator = usersgroupDescName.Validators[0].(func(string) error)
	workorderMixin := schema.WorkOrder{}.Mixin()
	workorderMixinFields := [...][]ent.Field{
		workorderMixin[0].Fields(),
	}
	workorderFields := schema.WorkOrder{}.Fields()
	_ = workorderFields
	// workorderDescCreateTime is the schema descriptor for create_time field.
	workorderDescCreateTime := workorderMixinFields[0][0].Descriptor()
	// workorder.DefaultCreateTime holds the default value on creation for the create_time field.
	workorder.DefaultCreateTime = workorderDescCreateTime.Default.(func() time.Time)
	// workorderDescUpdateTime is the schema descriptor for update_time field.
	workorderDescUpdateTime := workorderMixinFields[0][1].Descriptor()
	// workorder.DefaultUpdateTime holds the default value on creation for the update_time field.
	workorder.DefaultUpdateTime = workorderDescUpdateTime.Default.(func() time.Time)
	// workorder.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	workorder.UpdateDefaultUpdateTime = workorderDescUpdateTime.UpdateDefault.(func() time.Time)
	// workorderDescName is the schema descriptor for name field.
	workorderDescName := workorderFields[0].Descriptor()
	// workorder.NameValidator is a validator for the "name" field. It is called by the builders before save.
	workorder.NameValidator = workorderDescName.Validators[0].(func(string) error)
	// workorderDescStatus is the schema descriptor for status field.
	workorderDescStatus := workorderFields[1].Descriptor()
	// workorder.DefaultStatus holds the default value on creation for the status field.
	workorder.DefaultStatus = workorderDescStatus.Default.(string)
	// workorderDescPriority is the schema descriptor for priority field.
	workorderDescPriority := workorderFields[2].Descriptor()
	// workorder.DefaultPriority holds the default value on creation for the priority field.
	workorder.DefaultPriority = workorderDescPriority.Default.(string)
	workorderdefinitionMixin := schema.WorkOrderDefinition{}.Mixin()
	workorderdefinitionMixinFields := [...][]ent.Field{
		workorderdefinitionMixin[0].Fields(),
	}
	workorderdefinitionFields := schema.WorkOrderDefinition{}.Fields()
	_ = workorderdefinitionFields
	// workorderdefinitionDescCreateTime is the schema descriptor for create_time field.
	workorderdefinitionDescCreateTime := workorderdefinitionMixinFields[0][0].Descriptor()
	// workorderdefinition.DefaultCreateTime holds the default value on creation for the create_time field.
	workorderdefinition.DefaultCreateTime = workorderdefinitionDescCreateTime.Default.(func() time.Time)
	// workorderdefinitionDescUpdateTime is the schema descriptor for update_time field.
	workorderdefinitionDescUpdateTime := workorderdefinitionMixinFields[0][1].Descriptor()
	// workorderdefinition.DefaultUpdateTime holds the default value on creation for the update_time field.
	workorderdefinition.DefaultUpdateTime = workorderdefinitionDescUpdateTime.Default.(func() time.Time)
	// workorderdefinition.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	workorderdefinition.UpdateDefaultUpdateTime = workorderdefinitionDescUpdateTime.UpdateDefault.(func() time.Time)
	workordertypeMixin := schema.WorkOrderType{}.Mixin()
	workordertypeMixinFields := [...][]ent.Field{
		workordertypeMixin[0].Fields(),
	}
	workordertypeFields := schema.WorkOrderType{}.Fields()
	_ = workordertypeFields
	// workordertypeDescCreateTime is the schema descriptor for create_time field.
	workordertypeDescCreateTime := workordertypeMixinFields[0][0].Descriptor()
	// workordertype.DefaultCreateTime holds the default value on creation for the create_time field.
	workordertype.DefaultCreateTime = workordertypeDescCreateTime.Default.(func() time.Time)
	// workordertypeDescUpdateTime is the schema descriptor for update_time field.
	workordertypeDescUpdateTime := workordertypeMixinFields[0][1].Descriptor()
	// workordertype.DefaultUpdateTime holds the default value on creation for the update_time field.
	workordertype.DefaultUpdateTime = workordertypeDescUpdateTime.Default.(func() time.Time)
	// workordertype.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	workordertype.UpdateDefaultUpdateTime = workordertypeDescUpdateTime.UpdateDefault.(func() time.Time)
}
