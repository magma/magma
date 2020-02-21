// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/gqlerror"
	"golang.org/x/xerrors"
)

type (
	projectTypeResolver struct{}
	projectResolver     struct{}
)

var (
	errNoProjectType = gqlerror.Errorf("project type doesn't exist")
	errNoProject     = gqlerror.Errorf("project doesn't exist")
)

func (r queryResolver) ProjectType(ctx context.Context, id string) (*ent.ProjectType, error) {
	noder, err := r.Node(ctx, id)
	if err != nil {
		return nil, err
	}
	typ, _ := noder.(*ent.ProjectType)
	return typ, nil
}

func (projectTypeResolver) NumberOfProjects(ctx context.Context, obj *ent.ProjectType) (int, error) {
	return obj.QueryProjects().Count(ctx)
}

func (projectTypeResolver) Projects(ctx context.Context, obj *ent.ProjectType) ([]*ent.Project, error) {
	projects, err := obj.QueryProjects().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying projects: %w", err)
	}
	return projects, nil
}

func (projectTypeResolver) Properties(ctx context.Context, obj *ent.ProjectType) ([]*ent.PropertyType, error) {
	properties, err := obj.QueryProperties().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying properties: %w", err)
	}
	return properties, nil
}

func (projectTypeResolver) WorkOrders(ctx context.Context, obj *ent.ProjectType) ([]*ent.WorkOrderDefinition, error) {
	properties, err := obj.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying work order definitions: %w", err)
	}
	return properties, nil
}

func (r mutationResolver) CreateProjectType(ctx context.Context, input models.AddProjectTypeInput) (*ent.ProjectType, error) {
	properties, err := r.AddPropertyTypes(ctx, input.Properties...)
	if err != nil {
		return nil, xerrors.Errorf("creating properties: %w", err)
	}
	client := r.ClientFrom(ctx)
	typ, err := client.
		ProjectType.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		AddProperties(properties...).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("Project type %q already exists", input.Name)
		}
		return nil, xerrors.Errorf("creating project type: %w", err)
	}
	for _, wo := range input.WorkOrders {
		if _, err = client.WorkOrderDefinition.Create().
			SetNillableIndex(wo.Index).
			SetTypeID(wo.Type).
			SetProjectType(typ).
			Save(ctx); err != nil {
			return nil, xerrors.Errorf("creating work orders", err)
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
			err = r.validateAndAddNewPropertyType(ctx, p, func(b *ent.PropertyTypeUpdateOne) { b.SetProjectTypeID(pt.ID) })
		} else {
			err = r.updatePropType(ctx, p)
		}
		if err != nil {
			return nil, err
		}
	}

	var ids []string
	for _, wo := range input.WorkOrders {
		if wo.ID == nil {
			def, err := client.WorkOrderDefinition.Create().
				SetNillableIndex(wo.Index).
				SetTypeID(wo.Type).
				SetProjectType(pt).
				Save(ctx)
			if err != nil {
				return nil, xerrors.Errorf("creating work orders", err)
			}
			ids = append(ids, def.ID)
		} else {
			_, err := client.WorkOrderDefinition.UpdateOneID(*wo.ID).
				SetNillableIndex(wo.Index).
				SetTypeID(wo.Type).
				Save(ctx)
			if err != nil {
				return nil, xerrors.Errorf("creating work orders", err)
			}
			ids = append(ids, *wo.ID)
		}
	}
	ids, err = pt.QueryWorkOrders().Where(workorderdefinition.Not(workorderdefinition.IDIn(ids...))).IDs(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fetching work orders", err)
	}
	_, err = client.WorkOrderDefinition.Delete().Where(workorderdefinition.IDIn(ids...)).Exec(ctx)
	if err != nil {
		return nil, xerrors.Errorf("removing work orders", err)
	}
	return pt, nil
}

func (r mutationResolver) DeleteProjectType(ctx context.Context, id string) (bool, error) {
	client := r.ClientFrom(ctx)
	switch count, err := client.ProjectType.Query().Where(projecttype.ID(id)).QueryProjects().Count(ctx); {
	case err != nil:
		return false, xerrors.Errorf("cannot query project count for project type: %w", err)
	case count > 0:
		return false, gqlerror.Errorf("project type contains %d associated project", count)
	}
	if _, err := client.PropertyType.Delete().Where(propertytype.HasProjectTypeWith(projecttype.ID(id))).Exec(ctx); err != nil {
		return false, xerrors.Errorf("deleting project type properties: %w", err)
	}
	if err := client.ProjectType.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, errNoProjectType
		}
		return false, xerrors.Errorf("deleting project type: %w", err)
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

func (projectResolver) Type(ctx context.Context, obj *ent.Project) (*ent.ProjectType, error) {
	typ, err := obj.QueryType().Only(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying project type: %w", err)
	}
	return typ, nil
}

func (projectResolver) Location(ctx context.Context, obj *ent.Project) (*ent.Location, error) {
	l, err := obj.QueryLocation().Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, xerrors.Errorf("querying location: %w", err)
	}
	return l, nil
}

func (projectResolver) WorkOrders(ctx context.Context, obj *ent.Project) ([]*ent.WorkOrder, error) {
	wo, err := obj.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying work orders: %w", err)
	}
	return wo, nil
}

func (projectResolver) Properties(ctx context.Context, obj *ent.Project) ([]*ent.Property, error) {
	properties, err := obj.QueryProperties().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("querying properties: %w", err)
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
	properties, err := r.AddProperties(input.Properties, resolverutil.AddPropertyArgs{Context: ctx, IsTemplate: pointer.ToBool(true)})
	if err != nil {
		return nil, xerrors.Errorf("creating properties: %w", err)
	}
	client := r.ClientFrom(ctx)
	pt, err := client.ProjectType.Get(ctx, input.Type)
	if err != nil {
		return nil, xerrors.Errorf("fetching template", err)
	}
	proj, err := client.
		Project.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetTypeID(input.Type).
		SetNillableCreator(input.Creator).
		SetNillableLocationID(input.Location).
		AddProperties(properties...).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("Project %q already exists", input.Name)
		}
		return nil, xerrors.Errorf("creating project: %w", err)
	}
	wos, err := pt.QueryWorkOrders().All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fetching work orders templates", err)
	}
	for _, wo := range wos {
		wot := wo.QueryType().FirstX(ctx)
		w, err := r.AddWorkOrder(ctx, models.AddWorkOrderInput{
			Name:            wot.Name,
			Description:     &wot.Description,
			WorkOrderTypeID: wot.ID,
			ProjectID:       &proj.ID,
			LocationID:      input.Location,
			Index:           &wo.Index,
		})
		if err != nil {
			return nil, xerrors.Errorf("creating work order", err)
		}
		props, err := wot.QueryPropertyTypes().All(ctx)
		if err != nil {
			return nil, xerrors.Errorf("fetching work order properties", err)
		}
		for _, p := range props {
			var stringValue *string = nil
			if p.Type != models.PropertyKindEnum.String() {
				stringValue = &p.StringVal
			}

			newProp := &models.PropertyInput{
				PropertyTypeID: p.ID,
				StringValue:    stringValue,
				IntValue:       &p.IntVal,
				BooleanValue:   &p.BoolVal,
				FloatValue:     &p.FloatVal,
				LatitudeValue:  &p.LatitudeVal,
				LongitudeValue: &p.LongitudeVal,
				RangeFromValue: &p.RangeFromVal,
				RangeToValue:   &p.RangeToVal,
			}
			addPropertyArgs := resolverutil.AddPropertyArgs{
				Context:    ctx,
				EntSetter:  func(b *ent.PropertyCreate) { b.SetWorkOrder(w) },
				IsTemplate: pointer.ToBool(true),
			}

			_, err = r.AddProperty(newProp, addPropertyArgs)
			if err != nil {
				return nil, xerrors.Errorf("creating work order properties", err)
			}
		}
	}
	return proj, nil
}

func (r mutationResolver) DeleteProject(ctx context.Context, id string) (bool, error) {
	client := r.ClientFrom(ctx)
	if _, err := client.Property.Delete().Where(property.HasProjectWith(project.ID(id))).Exec(ctx); err != nil {
		return false, xerrors.Errorf("deleting project properties: %w", err)
	}
	if err := client.Project.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return false, errNoProject
		}
		return false, xerrors.Errorf("deleting project: %w", err)
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
	if input.Creator != nil {
		mutation.SetCreator(*input.Creator)
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
			query := client.Property.
				Update().
				Where(property.ID(existingProperty.ID))
			if _, err := updatePropValues(pInput, query).Save(ctx); err != nil {
				return nil, errors.Wrap(err, "updating property values")
			}
		}
	}
	return mutation.Save(ctx)
}
