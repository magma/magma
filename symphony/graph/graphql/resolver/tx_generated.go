// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by github.com/99designs/gqlgen, DO NOT EDIT.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

// txResolver wraps a mutation resolver and executes every mutation under a transaction.
type txResolver struct {
	generated.MutationResolver
}

func (tr txResolver) WithTransaction(ctx context.Context, f func(context.Context, generated.MutationResolver) error) error {
	tx, err := ent.FromContext(ctx).Tx(ctx)
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	ctx = ent.NewContext(ctx, tx.Client())
	if err := f(ctx, tr.MutationResolver); err != nil {
		if r := tx.Rollback(); r != nil {
			err = fmt.Errorf("rolling back transaction: %v", r)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (tr txResolver) EditUser(ctx context.Context, input models.EditUserInput) (*ent.User, error) {
	var result, zero *ent.User
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditUser(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) CreateSurvey(ctx context.Context, data models.SurveyCreateData) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.CreateSurvey(ctx, data)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddLocation(ctx context.Context, input models.AddLocationInput) (*ent.Location, error) {
	var result, zero *ent.Location
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddLocation(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditLocation(ctx context.Context, input models.EditLocationInput) (*ent.Location, error) {
	var result, zero *ent.Location
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditLocation(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveLocation(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveLocation(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddLocationType(ctx context.Context, input models.AddLocationTypeInput) (*ent.LocationType, error) {
	var result, zero *ent.LocationType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddLocationType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditLocationType(ctx context.Context, input models.EditLocationTypeInput) (*ent.LocationType, error) {
	var result, zero *ent.LocationType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditLocationType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveLocationType(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveLocationType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddEquipment(ctx context.Context, input models.AddEquipmentInput) (*ent.Equipment, error) {
	var result, zero *ent.Equipment
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddEquipment(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditEquipment(ctx context.Context, input models.EditEquipmentInput) (*ent.Equipment, error) {
	var result, zero *ent.Equipment
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditEquipment(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveEquipment(ctx context.Context, id int, workOrderID *int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveEquipment(ctx, id, workOrderID)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddEquipmentType(ctx context.Context, input models.AddEquipmentTypeInput) (*ent.EquipmentType, error) {
	var result, zero *ent.EquipmentType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddEquipmentType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditEquipmentType(ctx context.Context, input models.EditEquipmentTypeInput) (*ent.EquipmentType, error) {
	var result, zero *ent.EquipmentType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditEquipmentType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveEquipmentType(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveEquipmentType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddEquipmentPortType(ctx context.Context, input models.AddEquipmentPortTypeInput) (*ent.EquipmentPortType, error) {
	var result, zero *ent.EquipmentPortType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddEquipmentPortType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditEquipmentPortType(ctx context.Context, input models.EditEquipmentPortTypeInput) (*ent.EquipmentPortType, error) {
	var result, zero *ent.EquipmentPortType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditEquipmentPortType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveEquipmentPortType(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveEquipmentPortType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddLink(ctx context.Context, input models.AddLinkInput) (*ent.Link, error) {
	var result, zero *ent.Link
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddLink(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditLink(ctx context.Context, input models.EditLinkInput) (*ent.Link, error) {
	var result, zero *ent.Link
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditLink(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveLink(ctx context.Context, id int, workOrderID *int) (*ent.Link, error) {
	var result, zero *ent.Link
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveLink(ctx, id, workOrderID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddService(ctx context.Context, data models.ServiceCreateData) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddService(ctx, data)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditService(ctx context.Context, data models.ServiceEditData) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditService(ctx, data)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddServiceLink(ctx context.Context, id int, linkID int) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddServiceLink(ctx, id, linkID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveServiceLink(ctx context.Context, id int, linkID int) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveServiceLink(ctx, id, linkID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddServiceEndpoint(ctx context.Context, input models.AddServiceEndpointInput) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddServiceEndpoint(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveServiceEndpoint(ctx context.Context, serviceEndpointID int) (*ent.Service, error) {
	var result, zero *ent.Service
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveServiceEndpoint(ctx, serviceEndpointID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddServiceType(ctx context.Context, data models.ServiceTypeCreateData) (*ent.ServiceType, error) {
	var result, zero *ent.ServiceType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddServiceType(ctx, data)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditServiceType(ctx context.Context, data models.ServiceTypeEditData) (*ent.ServiceType, error) {
	var result, zero *ent.ServiceType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditServiceType(ctx, data)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveEquipmentFromPosition(ctx context.Context, positionID int, workOrderID *int) (*ent.EquipmentPosition, error) {
	var result, zero *ent.EquipmentPosition
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveEquipmentFromPosition(ctx, positionID, workOrderID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) MoveEquipmentToPosition(ctx context.Context, parentEquipmentID *int, positionDefinitionID *int, equipmentID int) (*ent.EquipmentPosition, error) {
	var result, zero *ent.EquipmentPosition
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.MoveEquipmentToPosition(ctx, parentEquipmentID, positionDefinitionID, equipmentID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddComment(ctx context.Context, input models.CommentInput) (*ent.Comment, error) {
	var result, zero *ent.Comment
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddComment(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddImage(ctx context.Context, input models.AddImageInput) (*ent.File, error) {
	var result, zero *ent.File
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddImage(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddHyperlink(ctx context.Context, input models.AddHyperlinkInput) (*ent.Hyperlink, error) {
	var result, zero *ent.Hyperlink
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddHyperlink(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) DeleteHyperlink(ctx context.Context, id int) (*ent.Hyperlink, error) {
	var result, zero *ent.Hyperlink
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.DeleteHyperlink(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) DeleteImage(ctx context.Context, entityType models.ImageEntity, entityID int, id int) (*ent.File, error) {
	var result, zero *ent.File
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.DeleteImage(ctx, entityType, entityID, id)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveWorkOrder(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveWorkOrder(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) ExecuteWorkOrder(ctx context.Context, id int) (*models.WorkOrderExecutionResult, error) {
	var result, zero *models.WorkOrderExecutionResult
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.ExecuteWorkOrder(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) RemoveWorkOrderType(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveWorkOrderType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) MarkSiteSurveyNeeded(ctx context.Context, locationID int, needed bool) (*ent.Location, error) {
	var result, zero *ent.Location
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.MarkSiteSurveyNeeded(ctx, locationID, needed)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveService(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveService(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) RemoveServiceType(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveServiceType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) EditLocationTypeSurveyTemplateCategories(ctx context.Context, id int, surveyTemplateCategories []*models.SurveyTemplateCategoryInput) ([]*ent.SurveyTemplateCategory, error) {
	var result, zero []*ent.SurveyTemplateCategory
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditLocationTypeSurveyTemplateCategories(ctx, id, surveyTemplateCategories)
		return
	}); err != nil {
		return zero, err
	}
	for i := range result {
		result[i] = result[i].Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditEquipmentPort(ctx context.Context, input models.EditEquipmentPortInput) (*ent.EquipmentPort, error) {
	var result, zero *ent.EquipmentPort
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditEquipmentPort(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) MarkLocationPropertyAsExternalID(ctx context.Context, propertyName string) (string, error) {
	var result, zero string
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.MarkLocationPropertyAsExternalID(ctx, propertyName)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) RemoveSiteSurvey(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveSiteSurvey(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddWiFiScans(ctx context.Context, data []*models.SurveyWiFiScanData, locationID int) ([]*ent.SurveyWiFiScan, error) {
	var result, zero []*ent.SurveyWiFiScan
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddWiFiScans(ctx, data, locationID)
		return
	}); err != nil {
		return zero, err
	}
	for i := range result {
		result[i] = result[i].Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddCellScans(ctx context.Context, data []*models.SurveyCellScanData, locationID int) ([]*ent.SurveyCellScan, error) {
	var result, zero []*ent.SurveyCellScan
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddCellScans(ctx, data, locationID)
		return
	}); err != nil {
		return zero, err
	}
	for i := range result {
		result[i] = result[i].Unwrap()
	}
	return result, nil
}

func (tr txResolver) MoveLocation(ctx context.Context, locationID int, parentLocationID *int) (*ent.Location, error) {
	var result, zero *ent.Location
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.MoveLocation(ctx, locationID, parentLocationID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditLocationTypesIndex(ctx context.Context, locationTypesIndex []*models.LocationTypeIndex) ([]*ent.LocationType, error) {
	var result, zero []*ent.LocationType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditLocationTypesIndex(ctx, locationTypesIndex)
		return
	}); err != nil {
		return zero, err
	}
	for i := range result {
		result[i] = result[i].Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddTechnician(ctx context.Context, input models.TechnicianInput) (*ent.Technician, error) {
	var result, zero *ent.Technician
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddTechnician(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddWorkOrder(ctx context.Context, input models.AddWorkOrderInput) (*ent.WorkOrder, error) {
	var result, zero *ent.WorkOrder
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddWorkOrder(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditWorkOrder(ctx context.Context, input models.EditWorkOrderInput) (*ent.WorkOrder, error) {
	var result, zero *ent.WorkOrder
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditWorkOrder(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddWorkOrderType(ctx context.Context, input models.AddWorkOrderTypeInput) (*ent.WorkOrderType, error) {
	var result, zero *ent.WorkOrderType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddWorkOrderType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditWorkOrderType(ctx context.Context, input models.EditWorkOrderTypeInput) (*ent.WorkOrderType, error) {
	var result, zero *ent.WorkOrderType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditWorkOrderType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) CreateProjectType(ctx context.Context, input models.AddProjectTypeInput) (*ent.ProjectType, error) {
	var result, zero *ent.ProjectType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.CreateProjectType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditProjectType(ctx context.Context, input models.EditProjectTypeInput) (*ent.ProjectType, error) {
	var result, zero *ent.ProjectType
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditProjectType(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) DeleteProjectType(ctx context.Context, id int) (bool, error) {
	var result, zero bool
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.DeleteProjectType(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) CreateProject(ctx context.Context, input models.AddProjectInput) (*ent.Project, error) {
	var result, zero *ent.Project
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.CreateProject(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditProject(ctx context.Context, input models.EditProjectInput) (*ent.Project, error) {
	var result, zero *ent.Project
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditProject(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) DeleteProject(ctx context.Context, id int) (bool, error) {
	var result, zero bool
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.DeleteProject(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddCustomer(ctx context.Context, input models.AddCustomerInput) (*ent.Customer, error) {
	var result, zero *ent.Customer
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddCustomer(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveCustomer(ctx context.Context, id int) (int, error) {
	var result, zero int
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveCustomer(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddFloorPlan(ctx context.Context, input models.AddFloorPlanInput) (*ent.FloorPlan, error) {
	var result, zero *ent.FloorPlan
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddFloorPlan(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) DeleteFloorPlan(ctx context.Context, id int) (bool, error) {
	var result, zero bool
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.DeleteFloorPlan(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) AddActionsRule(ctx context.Context, input models.AddActionsRuleInput) (*ent.ActionsRule, error) {
	var result, zero *ent.ActionsRule
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddActionsRule(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) EditActionsRule(ctx context.Context, id int, input models.AddActionsRuleInput) (*ent.ActionsRule, error) {
	var result, zero *ent.ActionsRule
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.EditActionsRule(ctx, id, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) RemoveActionsRule(ctx context.Context, id int) (bool, error) {
	var result, zero bool
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.RemoveActionsRule(ctx, id)
		return
	}); err != nil {
		return zero, err
	}
	return result, nil
}

func (tr txResolver) TechnicianWorkOrderCheckIn(ctx context.Context, workOrderID int) (*ent.WorkOrder, error) {
	var result, zero *ent.WorkOrder
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.TechnicianWorkOrderCheckIn(ctx, workOrderID)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}

func (tr txResolver) AddReportFilter(ctx context.Context, input models.ReportFilterInput) (*ent.ReportFilter, error) {
	var result, zero *ent.ReportFilter
	if err := tr.WithTransaction(ctx, func(ctx context.Context, mr generated.MutationResolver) (err error) {
		result, err = mr.AddReportFilter(ctx, input)
		return
	}); err != nil {
		return zero, err
	}
	if result != nil {
		result = result.Unwrap()
	}
	return result, nil
}
