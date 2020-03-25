// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/gqlerror"
	"go.uber.org/zap"
)

type workOrderDefinitionResolver struct{}

func (workOrderTypeResolver) CheckListDefinitions(ctx context.Context, obj *ent.WorkOrderType) ([]*ent.CheckListItemDefinition, error) {
	return obj.QueryCheckListDefinitions().All(ctx)
}

func (workOrderTypeResolver) CheckListCategories(ctx context.Context, obj *ent.WorkOrderType) ([]*ent.CheckListCategory, error) {
	return obj.QueryCheckListCategories().All(ctx)
}

func (workOrderDefinitionResolver) Type(ctx context.Context, obj *ent.WorkOrderDefinition) (*ent.WorkOrderType, error) {
	return obj.QueryType().Only(ctx)
}

type workOrderTypeResolver struct{}

func (workOrderTypeResolver) PropertyTypes(ctx context.Context, obj *ent.WorkOrderType) ([]*ent.PropertyType, error) {
	return obj.QueryPropertyTypes().All(ctx)
}

func (workOrderTypeResolver) NumberOfWorkOrders(ctx context.Context, obj *ent.WorkOrderType) (int, error) {
	return obj.QueryWorkOrders().Count(ctx)
}

type workOrderResolver struct{}

func (workOrderResolver) WorkOrderType(ctx context.Context, obj *ent.WorkOrder) (*ent.WorkOrderType, error) {
	return obj.QueryType().Only(ctx)
}

func (workOrderResolver) Location(ctx context.Context, obj *ent.WorkOrder) (*ent.Location, error) {
	l, err := obj.QueryLocation().Only(ctx)
	return l, ent.MaskNotFound(err)
}

func (workOrderResolver) Project(ctx context.Context, obj *ent.WorkOrder) (*ent.Project, error) {
	p, err := obj.QueryProject().Only(ctx)
	return p, ent.MaskNotFound(err)
}

func (workOrderResolver) CreationDate(_ context.Context, obj *ent.WorkOrder) (int, error) {
	secs := int(obj.CreationDate.Unix())
	return secs, nil
}

func (workOrderResolver) InstallDate(_ context.Context, obj *ent.WorkOrder) (*int, error) {
	secs := int(obj.InstallDate.Unix())
	return &secs, nil
}

func (workOrderResolver) Status(_ context.Context, obj *ent.WorkOrder) (models.WorkOrderStatus, error) {
	return models.WorkOrderStatus(obj.Status), nil
}

func (workOrderResolver) EquipmentToAdd(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Equipment, error) {
	return obj.QueryEquipment().Where(equipment.FutureState(models.FutureStateInstall.String())).All(ctx)
}

func (workOrderResolver) EquipmentToRemove(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Equipment, error) {
	return obj.QueryEquipment().Where(equipment.FutureState(models.FutureStateRemove.String())).All(ctx)
}

func (workOrderResolver) LinksToAdd(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Link, error) {
	return obj.QueryLinks().Where(link.FutureState(models.FutureStateInstall.String())).All(ctx)
}

func (workOrderResolver) LinksToRemove(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Link, error) {
	return obj.QueryLinks().Where(link.FutureState(models.FutureStateRemove.String())).All(ctx)
}

func (workOrderResolver) Properties(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

func (workOrderResolver) CheckList(ctx context.Context, obj *ent.WorkOrder) ([]*ent.CheckListItem, error) {
	return obj.QueryCheckListItems().All(ctx)
}

func (workOrderResolver) CheckListCategories(ctx context.Context, obj *ent.WorkOrder) ([]*ent.CheckListCategory, error) {
	return obj.QueryCheckListCategories().All(ctx)
}

func (workOrderResolver) Priority(_ context.Context, obj *ent.WorkOrder) (models.WorkOrderPriority, error) {
	return models.WorkOrderPriority(obj.Priority), nil
}

func (workOrderResolver) Images(ctx context.Context, obj *ent.WorkOrder) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeImage.String())).All(ctx)
}

func (workOrderResolver) Files(ctx context.Context, obj *ent.WorkOrder) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeFile.String())).All(ctx)
}

func (workOrderResolver) Hyperlinks(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Hyperlink, error) {
	return obj.QueryHyperlinks().All(ctx)
}

func (workOrderResolver) Comments(ctx context.Context, obj *ent.WorkOrder) ([]*ent.Comment, error) {
	return obj.QueryComments().All(ctx)
}

func (workOrderResolver) OwnerName(ctx context.Context, obj *ent.WorkOrder) (string, error) {
	owner, err := obj.QueryOwner().Only(ctx)
	if err != nil {
		return "", err
	}
	return owner.Email, nil
}

func (workOrderResolver) Assignee(ctx context.Context, obj *ent.WorkOrder) (*string, error) {
	assignee, err := obj.QueryAssignee().Only(ctx)
	if err != nil {
		return nil, ent.MaskNotFound(err)
	}
	return &assignee.Email, nil
}

func (r mutationResolver) AddWorkOrder(
	ctx context.Context,
	input models.AddWorkOrderInput,
) (*ent.WorkOrder, error) {
	return r.internalAddWorkOrder(ctx, input, false)
}

func (r mutationResolver) internalAddWorkOrder(
	ctx context.Context,
	input models.AddWorkOrderInput,
	skipMandatoryPropertiesCheck bool,
) (*ent.WorkOrder, error) {
	c := r.ClientFrom(ctx)
	propInput, err := r.validatedPropertyInputsFromTemplate(ctx, input.Properties, input.WorkOrderTypeID, models.PropertyEntityWorkOrder, skipMandatoryPropertiesCheck)
	if err != nil {
		return nil, fmt.Errorf("validating property for template : %w", err)
	}
	mutation := r.ClientFrom(ctx).
		WorkOrder.Create().
		SetName(input.Name).
		SetTypeID(input.WorkOrderTypeID).
		SetNillableProjectID(input.ProjectID).
		SetNillableLocationID(input.LocationID).
		SetNillableDescription(input.Description).
		SetCreationDate(time.Now()).
		SetNillableIndex(input.Index)
	if input.Status != nil {
		mutation.SetStatus(input.Status.String())
		if *input.Status == models.WorkOrderStatusDone {
			mutation.SetCloseDate(time.Now())
		}
	}
	if input.Priority != nil {
		mutation.SetPriority(input.Priority.String())
	}
	if input.Assignee != nil && *input.Assignee != "" {
		assigneeID, err := c.User.Query().Where(user.AuthID(*input.Assignee)).OnlyID(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching assignee user: %w", err)
		} else {
			mutation = mutation.SetAssigneeID(assigneeID)
		}
	}
	var owner *ent.User
	if input.OwnerName != nil && *input.OwnerName != "" {
		owner, err = c.User.Query().Where(user.AuthID(*input.OwnerName)).Only(ctx)
	} else {
		owner, err = viewer.UserFromContext(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("fetching own user: %w", err)
	}
	mutation = mutation.SetOwner(owner)
	for _, clInput := range input.CheckListCategories {
		checkListCategory, err := r.createOrUpdateCheckListCategory(ctx, clInput)
		if err != nil {
			return nil, errors.Wrap(err, "creating check list category")
		}
		mutation = mutation.AddCheckListCategories(checkListCategory)
	}
	wo, err := mutation.Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating work order")
	}
	if _, err := r.AddProperties(propInput,
		resolverutil.AddPropertyArgs{
			Context:    ctx,
			EntSetter:  func(b *ent.PropertyCreate) { b.SetWorkOrderID(wo.ID) },
			IsTemplate: pointer.ToBool(true),
		},
	); err != nil {
		return nil, errors.Wrap(err, "creating work order properties")
	}
	for _, clInput := range input.CheckList {
		if _, err = c.CheckListItem.Create().
			SetTitle(clInput.Title).
			SetType(clInput.Type.String()).
			SetNillableIndex(clInput.Index).
			SetNillableEnumValues(clInput.EnumValues).
			SetNillableHelpText(clInput.HelpText).
			SetNillableChecked(clInput.Checked).
			SetNillableStringVal(clInput.StringValue).
			SetWorkOrderID(wo.ID).
			Save(ctx); err != nil {
			return nil, errors.Wrap(err, "creating check list item")
		}
	}
	return wo, nil
}

func (r mutationResolver) EditWorkOrder(
	ctx context.Context, input models.EditWorkOrderInput,
) (*ent.WorkOrder, error) {
	client := r.ClientFrom(ctx)
	wo, err := client.WorkOrder.Get(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrap(err, "querying work order")
	}
	mutation := client.WorkOrder.
		UpdateOne(wo).
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetStatus(input.Status.String()).
		SetPriority(input.Priority.String()).
		SetNillableIndex(input.Index)
	if input.OwnerName != nil && *input.OwnerName != "" {
		ownerID, err := client.User.Query().Where(user.AuthID(*input.OwnerName)).OnlyID(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching owner user: %w", err)
		}
		mutation = mutation.SetOwnerID(ownerID)
	}
	if input.Assignee != nil && *input.Assignee != "" {
		assigneeID, err := client.User.Query().Where(user.AuthID(*input.Assignee)).OnlyID(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching assignee user: %w", err)
		}
		mutation = mutation.SetAssigneeID(assigneeID)
	}
	if input.Status == models.WorkOrderStatusDone {
		mutation.SetCloseDate(time.Now())
	} else {
		mutation.ClearCloseDate()
	}
	if input.InstallDate != nil {
		mutation.SetInstallDate(*input.InstallDate)
	} else {
		mutation.ClearInstallDate()
	}
	if input.ProjectID != nil {
		mutation.SetProjectID(*input.ProjectID)
	} else {
		mutation.ClearProject()
	}
	if input.LocationID != nil {
		mutation.SetLocationID(*input.LocationID)
	} else {
		mutation.ClearLocation()
	}

	for _, pInput := range input.Properties {
		err := r.updateProperty(ctx, wo.QueryProperties(), pInput)
		if err != nil {
			return nil, errors.Wrap(err, "updating work order property value")
		}
	}

	ids := make([]int, 0, len(input.CheckList))
	for _, clInput := range input.CheckList {
		cli, err := r.createOrUpdateCheckListItem(ctx, clInput)
		if err != nil {
			return nil, err
		}
		ids = append(ids, cli.ID)
	}
	currentCL, err := wo.QueryCheckListItems().IDs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "querying checklist items of work order %q", wo.ID)
	}
	addedCLIds, deletedCLIds := resolverutil.GetDifferenceBetweenSlices(currentCL, ids)
	mutation.
		RemoveCheckListItemIDs(deletedCLIds...).
		AddCheckListItemIDs(addedCLIds...)

	ids = make([]int, 0, len(input.CheckListCategories))
	for _, clInput := range input.CheckListCategories {
		cli, err := r.createOrUpdateCheckListCategory(ctx, clInput)
		if err != nil {
			return nil, err
		}
		ids = append(ids, cli.ID)
	}
	currentCL, err = wo.QueryCheckListCategories().IDs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "querying checklist categories of work order %q", wo.ID)
	}
	addedCLIds, deletedCLIds = resolverutil.GetDifferenceBetweenSlices(currentCL, ids)
	mutation.
		RemoveCheckListCategoryIDs(deletedCLIds...).
		AddCheckListCategoryIDs(addedCLIds...)

	return mutation.Save(ctx)
}

func (r mutationResolver) updateProperty(
	ctx context.Context,
	query *ent.PropertyQuery,
	input *models.PropertyInput) error {
	propertyQuery := query.
		Where(property.HasTypeWith(propertytype.ID(input.PropertyTypeID)))
	if input.ID != nil {
		propertyQuery = propertyQuery.
			Where(property.ID(*input.ID))
	}
	existingProperty, err := propertyQuery.Only(ctx)
	if err != nil {
		if input.ID == nil {
			return errors.Wrapf(err, "querying property type %q", input.PropertyTypeID)
		}
		return errors.Wrapf(err, "querying property type %q and id %q", input.PropertyTypeID, *input.ID)
	}
	client := r.ClientFrom(ctx)
	typ, err := client.PropertyType.Get(ctx, input.PropertyTypeID)
	if err != nil {
		return errors.Wrapf(err, "querying property type %q", input.PropertyTypeID)
	}
	if typ.Editable && typ.IsInstanceProperty {
		existingPropQuery := client.Property.
			Update().
			Where(property.ID(existingProperty.ID))

		if _, err = updatePropValues(input, existingPropQuery).Save(ctx); err != nil {
			return errors.Wrap(err, "saving work order property value update")
		}
	}
	return nil
}

func (r mutationResolver) createOrUpdateCheckListCategory(
	ctx context.Context,
	clInput *models.CheckListCategoryInput) (*ent.CheckListCategory, error) {
	client := r.ClientFrom(ctx)
	cl := client.CheckListCategory
	var clc *ent.CheckListCategory
	var err error
	if clInput.ID == nil {
		clc, err = cl.Create().
			SetTitle(clInput.Title).
			SetNillableDescription(clInput.Description).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "creating check list category")
		}
	} else {
		clc, err = cl.UpdateOneID(*clInput.ID).
			SetTitle(clInput.Title).
			SetNillableDescription(clInput.Description).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "updating check list category")
		}
	}
	mutation := cl.UpdateOneID(clc.ID)
	addedCLIds, deletedCLIds, err := r.createOrUpdateCheckListItems(ctx, clc, clInput.CheckList)
	if err != nil {
		return nil, errors.Wrap(err, "updating check list category items")
	}
	return mutation.
		RemoveCheckListItemIDs(deletedCLIds...).
		AddCheckListItemIDs(addedCLIds...).
		Save(ctx)
}

func (r mutationResolver) createOrUpdateCheckListItems(
	ctx context.Context,
	clc *ent.CheckListCategory,
	inputs []*models.CheckListItemInput) ([]int, []int, error) {
	ids := make([]int, 0, len(inputs))
	for _, input := range inputs {
		cli, err := r.createOrUpdateCheckListItem(ctx, input)
		if err != nil {
			return nil, nil, err
		}
		if cli != nil {
			ids = append(ids, cli.ID)
		}
	}
	currentCLIds := clc.QueryCheckListItems().IDsX(ctx)
	addedCLIds, deletedCLIds := resolverutil.GetDifferenceBetweenSlices(currentCLIds, ids)
	return addedCLIds, deletedCLIds, nil
}

func (r mutationResolver) createOrUpdateCheckListItem(
	ctx context.Context,
	input *models.CheckListItemInput) (*ent.CheckListItem, error) {
	client := r.ClientFrom(ctx)
	cl := client.CheckListItem
	var selectionMode *string
	var yesNoVal *checklistitem.YesNoVal
	if input.EnumSelectionMode != nil {
		selectionMode = pointer.ToString(input.EnumSelectionMode.String())
	}
	if input.YesNoResponse != nil {
		var yesNo checklistitem.YesNoVal
		if *input.YesNoResponse == models.YesNoResponseYes {
			yesNo = checklistitem.YesNoValYES
		} else {
			yesNo = checklistitem.YesNoValNO
		}
		yesNoVal = &yesNo
	}

	var cli *ent.CheckListItem
	var err error
	if input.ID == nil {
		cli, err = cl.Create().
			SetTitle(input.Title).
			SetType(input.Type.String()).
			SetNillableIndex(input.Index).
			SetNillableEnumValues(input.EnumValues).
			SetNillableHelpText(input.HelpText).
			SetNillableChecked(input.Checked).
			SetNillableStringVal(input.StringValue).
			SetNillableEnumSelectionMode(selectionMode).
			SetNillableSelectedEnumValues(input.SelectedEnumValues).
			SetNillableYesNoVal(yesNoVal).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "creating check list item")
		}
	} else {
		cli, err = cl.UpdateOneID(*input.ID).
			SetTitle(input.Title).
			SetType(input.Type.String()).
			SetNillableIndex(input.Index).
			SetNillableEnumValues(input.EnumValues).
			SetNillableHelpText(input.HelpText).
			SetNillableChecked(input.Checked).
			SetNillableStringVal(input.StringValue).
			SetNillableEnumSelectionMode(selectionMode).
			SetNillableSelectedEnumValues(input.SelectedEnumValues).
			SetNillableYesNoVal(yesNoVal).
			Save(ctx)
	}
	if err != nil {
		return nil, errors.Wrap(err, "updating check list item")
	}

	return r.createOrUpdateCheckListItemFiles(ctx, cli, input.Files)
}

func toIDSet(ids []int) map[int]bool {
	idSet := make(map[int]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	return idSet
}

func (r mutationResolver) deleteRemovedCheckListItemFiles(ctx context.Context, item *ent.CheckListItem, currentFileIDs []int, inputFileIDs []int) (*ent.CheckListItem, map[int]bool, error) {
	client := r.ClientFrom(ctx)
	_, deletedFileIDs := resolverutil.GetDifferenceBetweenSlices(currentFileIDs, inputFileIDs)
	deletedIDSet := toIDSet(deletedFileIDs)

	for _, fileID := range deletedFileIDs {
		err := client.File.DeleteOneID(fileID).Exec(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("deleting checklist file: file=%q: %w", fileID, err)
		}
	}

	item, err := client.CheckListItem.UpdateOne(item).RemoveFileIDs(deletedFileIDs...).Save(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("removing checklist files, item: %q: %w", item, err)
	}

	return item, deletedIDSet, nil
}

func (r mutationResolver) createAddedCheckListItemFiles(ctx context.Context, item *ent.CheckListItem, fileInputs []*models.FileInput) (*ent.CheckListItem, error) {
	client := r.ClientFrom(ctx)
	var addedFiles []*ent.File
	for _, input := range fileInputs {
		if input.ID != nil {
			continue
		}
		f, err := r.createImage(
			ctx,
			&models.AddImageInput{
				ImgKey:   input.StoreKey,
				FileName: input.FileName,
				FileSize: func() int {
					if input.SizeInBytes != nil {
						return *input.SizeInBytes
					}
					return 0
				}(),
				Modified:    time.Now(),
				ContentType: models.FileTypeFile.String(),
			},
		)
		if err != nil {
			return nil, err
		}

		addedFiles = append(addedFiles, f)
	}

	if len(addedFiles) > 0 {
		if item, err := client.CheckListItem.
			UpdateOne(item).
			AddFiles(addedFiles...).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("adding checklist file item=%q %w", item.ID, err)
		}
	}

	return item, nil
}

func (r mutationResolver) createOrUpdateCheckListItemFiles(ctx context.Context, item *ent.CheckListItem, fileInputs []*models.FileInput) (*ent.CheckListItem, error) {
	client := r.ClientFrom(ctx)
	currentFileIDs, err := client.CheckListItem.QueryFiles(item).IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying checklist files, item=%q: %w", item.ID, err)
	}
	inputFileIDs := make([]int, 0, len(fileInputs))
	for _, fileInput := range fileInputs {
		if fileInput.ID == nil {
			continue
		}
		inputFileIDs = append(inputFileIDs, *fileInput.ID)
	}

	item, deletedIDSet, err := r.deleteRemovedCheckListItemFiles(ctx, item, currentFileIDs, inputFileIDs)
	if err != nil {
		return nil, fmt.Errorf("deleting checklist files, item=%q: %w", item.ID, err)
	}

	item, err = r.createAddedCheckListItemFiles(ctx, item, fileInputs)
	if err != nil {
		return nil, fmt.Errorf("creating checklist files, item=%q: %w", item.ID, err)
	}

	for _, input := range fileInputs {
		if input.ID == nil {
			continue
		}
		if _, ok := deletedIDSet[*input.ID]; ok {
			continue
		}

		existingFile, err := client.File.Get(ctx, *input.ID)
		if err != nil {
			return nil, fmt.Errorf("querying file: file=%q: %w", *input.ID, err)
		}
		if existingFile.Name == input.FileName {
			continue
		}
		_, err = client.File.UpdateOne(existingFile).SetName(input.FileName).SetModifiedAt(time.Now()).Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("updating file name: file=%q: %w", existingFile.ID, err)
		}
	}

	return item, nil
}

func (r mutationResolver) AddWorkOrderType(
	ctx context.Context, input models.AddWorkOrderTypeInput) (*ent.WorkOrderType, error) {
	props, err := r.AddPropertyTypes(ctx, input.Properties...)
	if err != nil {
		return nil, err
	}

	client := r.ClientFrom(ctx)
	mutation := client.WorkOrderType.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		AddPropertyTypes(props...)
	for _, clInput := range input.CheckListCategories {
		checkListCategory, err := client.CheckListCategory.Create().
			SetTitle(clInput.Title).
			SetNillableDescription(clInput.Description).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "creating check list category")
		}
		mutation = mutation.AddCheckListCategories(checkListCategory)
	}
	typ, err := mutation.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A work order type with the name %v already exists", input.Name)
		}
		return nil, errors.Wrap(err, "creating work order type")
	}

	for _, def := range input.CheckList {
		def := def
		if _, err = client.CheckListItemDefinition.Create().
			SetTitle(def.Title).
			SetType(def.Type.String()).
			SetNillableIndex(def.Index).
			SetNillableHelpText(def.HelpText).
			SetNillableEnumValues(def.EnumValues).
			SetWorkOrderType(typ).
			Save(ctx); err != nil {
			return nil, errors.Wrap(err, "creating check list item")
		}
	}
	return typ, nil
}

func (r mutationResolver) EditWorkOrderType(
	ctx context.Context, input models.EditWorkOrderTypeInput,
) (*ent.WorkOrderType, error) {
	client := r.ClientFrom(ctx)
	wot, err := client.WorkOrderType.Get(ctx, input.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, gqlerror.Errorf("A work order template with id=%q does not exist", input.ID)
		}
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A work order template with the name %v already exists", input.Name)
		}
		return nil, errors.Wrapf(err, "updating work order template: id=%q", input.ID)
	}
	for _, p := range input.Properties {
		if p.ID == nil {
			err = r.validateAndAddNewPropertyType(ctx, p, func(b *ent.PropertyTypeUpdateOne) { b.SetWorkOrderTypeID(input.ID) })
		} else {
			err = r.updatePropType(ctx, p)
		}
		if err != nil {
			return nil, err
		}
	}

	mutation := client.WorkOrderType.
		UpdateOneID(input.ID).
		SetName(input.Name).
		SetNillableDescription(input.Description)

	currentCL := wot.QueryCheckListDefinitions().IDsX(ctx)
	ids := make([]int, 0, len(input.CheckList))
	for _, clInput := range input.CheckList {
		cli, err := r.createOrUpdateCheckListDefinition(ctx, clInput, input.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, cli.ID)
	}
	_, deletedCLIds := resolverutil.GetDifferenceBetweenSlices(currentCL, ids)
	mutation.RemoveCheckListDefinitionIDs(deletedCLIds...)
	return mutation.Save(ctx)
}

func (r mutationResolver) createOrUpdateCheckListDefinition(
	ctx context.Context,
	clInput *models.CheckListDefinitionInput,
	wotID int) (*ent.CheckListItemDefinition, error) {
	client := r.ClientFrom(ctx)
	cl := client.CheckListItemDefinition
	if clInput.ID == nil {
		cli, err := cl.Create().
			SetTitle(clInput.Title).
			SetType(clInput.Type.String()).
			SetNillableIndex(clInput.Index).
			SetNillableEnumValues(clInput.EnumValues).
			SetNillableHelpText(clInput.HelpText).
			SetWorkOrderTypeID(wotID).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "creating check list definition")
		}
		return cli, nil
	}

	cli, err := cl.UpdateOneID(*clInput.ID).
		SetTitle(clInput.Title).
		SetType(clInput.Type.String()).
		SetNillableIndex(clInput.Index).
		SetNillableEnumValues(clInput.EnumValues).
		SetNillableHelpText(clInput.HelpText).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "updating check list definition")
	}
	return cli, nil
}

func (r mutationResolver) RemoveWorkOrderType(ctx context.Context, id int) (int, error) {
	client, logger := r.ClientFrom(ctx), r.logger.For(ctx).With(zap.Int("id", id))
	switch count, err := client.WorkOrderType.Query().
		Where(workordertype.ID(id)).
		QueryWorkOrders().
		Count(ctx); {
	case err != nil:
		logger.Error("cannot query work order count of type", zap.Error(err))
		return id, fmt.Errorf("querying work orders for type: %w", err)
	case count > 0:
		logger.Warn("work order type has existing work orders", zap.Int("count", count))
		return id, gqlerror.Errorf("cannot delete work order type with %d existing work orders", count)
	}
	if _, err := client.PropertyType.Delete().
		Where(propertytype.HasWorkOrderTypeWith(workordertype.ID(id))).
		Exec(ctx); err != nil {
		logger.Error("cannot delete properties of work order type", zap.Error(err))
		return id, fmt.Errorf("deleting work order property types: %w", err)
	}
	switch err := client.WorkOrderType.DeleteOneID(id).Exec(ctx); err.(type) {
	case nil:
		logger.Info("deleted work order type")
		return id, nil
	case *ent.NotFoundError:
		err := gqlerror.Errorf("work order type not found")
		logger.Error(err.Message)
		return id, err
	default:
		logger.Error("cannot delete work order type", zap.Error(err))
		return id, fmt.Errorf("deleting work order type: %w", err)
	}
}
