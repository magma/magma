// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorderdefinition"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type (
	projectTypeResolver struct{}
	projectResolver     struct{}
)

var (
	errNoProjectType = gqlerror.Errorf("project type doesn't exist")
	errNoProject     = gqlerror.Errorf("project doesn't exist")
)

func (projectTypeResolver) NumberOfProjects(ctx context.Context, obj *ent.ProjectType) (int, error) {
	return obj.QueryProjects().Count(ctx)
}

func (projectTypeResolver) Projects(ctx context.Context, obj *ent.ProjectType) ([]*ent.Project, error) {
	projects, err := obj.QueryProjects().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying projects: %w", err)
	}
	return projects, nil
}

func (projectTypeResolver) Properties(ctx context.Context, obj *ent.ProjectType) ([]*ent.PropertyType, error) {
	properties, err := obj.QueryProperties().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying properties: %w", err)
	}
	return properties, nil
}

func (projectTypeResolver) WorkOrders(ctx context.Context, obj *ent.ProjectType) ([]*ent.WorkOrderDefinition, error) {
	properties, err := obj.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying work order definitions: %w", err)
	}
	return properties, nil
}

func (r mutationResolver) CreateProjectType(ctx context.Context, input models.AddProjectTypeInput) (*ent.ProjectType, error) {
	client := r.ClientFrom(ctx)
	typ, err := client.
		ProjectType.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("Project type %q already exists", input.Name)
		}
		return nil, fmt.Errorf("creating project type: %w", err)
	}
	if err := r.AddPropertyTypes(ctx, func(ptc *ent.PropertyTypeCreate) {
		ptc.SetProjectTypeID(typ.ID)
	}, input.Properties...); err != nil {
		return nil, fmt.Errorf("creating properties: %w", err)
	}
	for _, wo := range input.WorkOrders {
		if _, err = client.WorkOrderDefinition.Create().
			SetNillableIndex(wo.Index).
			SetTypeID(wo.Type).
			SetProjectType(typ).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating work orders: %w", err)
		}
	}
	return typ, nil
}

func (r mutationResolver) EditProjectType(
	ctx context.Context, input models.EditProjectTypeInput,
) (*ent.ProjectType, error) {
	client := r.ClientFrom(ctx)
	pt, err := client.ProjectType.
		UpdateOneID(input.ID).
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, gqlerror.Errorf("Project template with id=%q does not exist", input.ID)
		}
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A project template with the name %v already exists", input.Name)
		}
		return nil, errors.Wrapf(err, "updating project template: id=%q", pt.ID)
	}
	for _, p := range input.Properties {
		if p.ID == nil {
			err = r.validateAndAddNewPropertyType(ctx, p, func(b *ent.PropertyTypeCreate) { b.SetProjectTypeID(pt.ID) })
		} else {
			err = r.updatePropType(ctx, p)
		}
		if err != nil {
			return nil, err
		}
	}

	var ids []int
	for _, wo := range input.WorkOrders {
		if wo.ID == nil {
			def, err := client.WorkOrderDefinition.Create().
				SetNillableIndex(wo.Index).
				SetTypeID(wo.Type).
				SetProjectType(pt).
				Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("creating work orders: %w", err)
			}
			ids = append(ids, def.ID)
		} else {
			_, err := client.WorkOrderDefinition.UpdateOneID(*wo.ID).
				SetNillableIndex(wo.Index).
				SetTypeID(wo.Type).
				Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("creating work orders: %w", err)
			}
			ids = append(ids, *wo.ID)
		}
	}
	ids, err = pt.QueryWorkOrders().Where(workorderdefinition.Not(workorderdefinition.IDIn(ids...))).IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching work orders: %w", err)
	}
	for _, id := range ids {
		if err := client.WorkOrderDefinition.DeleteOneID(id).Exec(ctx); err != nil {
			return nil, fmt.Errorf("removing work orders: %w", err)
		}
	}
	return pt, nil
}

func (r mutationResolver) DeleteProjectType(ctx context.Context, id int) (bool, error) {
	client := r.ClientFrom(ctx)
	switch count, err := client.ProjectType.Query().Where(projecttype.ID(id)).QueryProjects().Count(ctx); {
	case err != nil:
		return false, fmt.Errorf("cannot query project count for project type: %w", err)
	case count > 0:
		return false, gqlerror.Errorf("project type contains %d associated project", count)
	}
	pTypes, err := client.PropertyType.Query().Where(propertytype.HasProjectTypeWith(projecttype.ID(id))).All(ctx)
	if err != nil {
		return false, fmt.Errorf("querying project type properties: %w", err)
	}
	for _, pType := range pTypes {
		if err := client.PropertyType.DeleteOne(pType).Exec(ctx); err != nil {
			return false, fmt.Errorf("deleting project type property: %w", err)
		}
	}
	if err := client.ProjectType.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, errNoProjectType
		}
		return false, fmt.Errorf("deleting project type: %w", err)
	}
	return true, nil
}

func (r queryResolver) ProjectTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.ProjectTypeConnection, error) {
	return r.ClientFrom(ctx).ProjectType.Query().
		Paginate(ctx, after, first, before, last)
}

func (projectResolver) Creator(ctx context.Context, obj *ent.Project) (*string, error) {
	assignee, err := obj.QueryCreator().Only(ctx)
	if err != nil {
		return nil, ent.MaskNotFound(err)
	}
	return &assignee.Email, nil
}

func (projectResolver) CreatedBy(ctx context.Context, obj *ent.Project) (*ent.User, error) {
	c, err := obj.QueryCreator().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying creator: %w", err)
	}
	return c, nil
}

func (projectResolver) Type(ctx context.Context, obj *ent.Project) (*ent.ProjectType, error) {
	typ, err := obj.Edges.TypeOrErr()
	if ent.IsNotLoaded(err) {
		return obj.QueryType().Only(ctx)
	}
	return typ, err
}

func (projectResolver) Location(ctx context.Context, obj *ent.Project) (*ent.Location, error) {
	l, err := obj.QueryLocation().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying location: %w", err)
	}
	return l, nil
}

func (projectResolver) WorkOrders(ctx context.Context, obj *ent.Project) ([]*ent.WorkOrder, error) {
	wo, err := obj.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying work orders: %w", err)
	}
	return wo, nil
}

func (projectResolver) Properties(ctx context.Context, obj *ent.Project) ([]*ent.Property, error) {
	properties, err := obj.QueryProperties().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying properties: %w", err)
	}
	return properties, nil
}

func (projectResolver) Comments(ctx context.Context, obj *ent.Project) ([]*ent.Comment, error) {
	return obj.QueryComments().All(ctx)
}

func (projectResolver) NumberOfWorkOrders(ctx context.Context, obj *ent.Project) (int, error) {
	return obj.QueryWorkOrders().Count(ctx)
}

func (r mutationResolver) CreateProject(ctx context.Context, input models.AddProjectInput) (*ent.Project, error) {
	client := r.ClientFrom(ctx)
	creatorID, err := resolverutil.GetUserID(ctx, input.CreatorID, input.Creator)
	if err != nil {
		return nil, err
	}
	proj, err := client.
		Project.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetTypeID(input.Type).
		SetNillableLocationID(input.Location).
		SetNillableCreatorID(creatorID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("Project %q already exists", input.Name)
		}
		return nil, fmt.Errorf("creating project: %w", err)
	}
	propInput, err := r.validatedPropertyInputsFromTemplate(ctx, input.Properties, input.Type, models.PropertyEntityProject, false)
	if err != nil {
		return nil, fmt.Errorf("validating property for template : %w", err)
	}
	if _, err := r.AddProperties(propInput, resolverutil.AddPropertyArgs{
		Context:    ctx,
		IsTemplate: pointer.ToBool(true),
		EntSetter: func(create *ent.PropertyCreate) {
			create.SetProject(proj)
		}}); err != nil {
		return nil, fmt.Errorf("creating properties: %w", err)
	}
	pt, err := client.ProjectType.Get(ctx, input.Type)
	if err != nil {
		return nil, fmt.Errorf("fetching template: %w", err)
	}
	wos, err := pt.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching work orders templates: %w", err)
	}
	for _, wo := range wos {
		wot, err := wo.QueryType().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("query work order definition type: %w", err)
		}

		clCategoryDefs, err := wot.QueryCheckListCategoryDefinitions().WithCheckListItemDefinitions().All(ctx)
		if err != nil {
			return nil, fmt.Errorf("query work order checklist definitions: %w", err)
		}

		var categoryInputs []*models.CheckListCategoryInput
		for _, categoryDef := range clCategoryDefs {
			var clInputs []*models.CheckListItemInput
			for _, cliDef := range categoryDef.Edges.CheckListItemDefinitions {
				var enumSelectionMode *checklistitem.EnumSelectionModeValue
				if cliDef.EnumSelectionModeValue != "" {
					mode := checklistitem.EnumSelectionModeValue(cliDef.EnumSelectionModeValue.String())
					enumSelectionMode = &mode
				}

				clInputs = append(clInputs, &models.CheckListItemInput{
					Title:             cliDef.Title,
					Type:              models.CheckListItemType(cliDef.Type),
					Index:             pointer.ToInt(cliDef.Index),
					HelpText:          cliDef.HelpText,
					EnumValues:        cliDef.EnumValues,
					EnumSelectionMode: enumSelectionMode,
				})
			}

			categoryInputs = append(categoryInputs, &models.CheckListCategoryInput{
				Title:       categoryDef.Title,
				Description: pointer.ToString(categoryDef.Description),
				CheckList:   clInputs,
			})
		}

		_, err = r.internalAddWorkOrder(ctx, models.AddWorkOrderInput{
			Name:                wot.Name,
			Description:         &wot.Description,
			WorkOrderTypeID:     wot.ID,
			ProjectID:           &proj.ID,
			LocationID:          input.Location,
			Index:               &wo.Index,
			CheckListCategories: categoryInputs,
		}, true)
		if err != nil {
			return nil, fmt.Errorf("creating work order: %w", err)
		}
	}
	return proj, nil
}

func (r mutationResolver) DeleteProject(ctx context.Context, id int) (bool, error) {
	client := r.ClientFrom(ctx)
	props, err := client.Property.Query().Where(property.HasProjectWith(project.ID(id))).All(ctx)
	if err != nil {
		return false, fmt.Errorf("querying project properties: %w", err)
	}
	for _, prop := range props {
		if err := client.Property.DeleteOne(prop).Exec(ctx); err != nil {
			return false, fmt.Errorf("deleting project properties: %w", err)
		}
	}
	if err := client.Project.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, errNoProject
		}
		return false, fmt.Errorf("deleting project: %w", err)
	}
	return true, nil
}

func (r mutationResolver) EditProject(ctx context.Context, input models.EditProjectInput) (*ent.Project, error) {
	client := r.ClientFrom(ctx)
	proj, err := client.Project.Get(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying project: id=%q", input.ID)
	}

	mutation := client.Project.
		UpdateOne(proj).
		SetName(input.Name).
		SetNillableDescription(input.Description)

	creatorID, err := resolverutil.GetUserID(ctx, input.CreatorID, input.Creator)
	if err != nil {
		return nil, err
	}
	if creatorID != nil {
		mutation.SetCreatorID(*creatorID)
	} else {
		mutation.ClearCreator()
	}
	if input.Location != nil {
		mutation.SetLocationID(*input.Location)
	} else {
		mutation.ClearLocation()
	}
	for _, pInput := range input.Properties {
		propertyQuery := proj.QueryProperties().
			Where(property.HasTypeWith(propertytype.ID(pInput.PropertyTypeID)))
		if pInput.ID != nil {
			propertyQuery = propertyQuery.
				Where(property.ID(*pInput.ID))
		}
		existingProperty, err := propertyQuery.Only(ctx)
		if err != nil {
			if pInput.ID == nil {
				return nil, errors.Wrapf(err, "querying project property type %q", pInput.PropertyTypeID)
			}
			return nil, errors.Wrapf(err, "querying project property type %q and id %q", pInput.PropertyTypeID, *pInput.ID)
		}
		typ, err := client.PropertyType.Get(ctx, pInput.PropertyTypeID)
		if err != nil {
			return nil, errors.Wrapf(err, "querying property type %q", pInput.PropertyTypeID)
		}
		if typ.Editable && typ.IsInstanceProperty {
			updater := client.Property.UpdateOneID(existingProperty.ID)
			if r.updatePropValues(ctx, pInput, updater) != nil {
				return nil, errors.Wrap(err, "updating property values")
			}
		}
	}
	return mutation.Save(ctx)
}
