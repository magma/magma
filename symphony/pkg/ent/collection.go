// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ar *ActionsRuleQuery) CollectFields(ctx context.Context, satisfies ...string) *ActionsRuleQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ar = ar.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ar
}

func (ar *ActionsRuleQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ActionsRuleQuery {
	return ar
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (a *ActivityQuery) CollectFields(ctx context.Context, satisfies ...string) *ActivityQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		a = a.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return a
}

func (a *ActivityQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ActivityQuery {
	return a
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (clc *CheckListCategoryQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListCategoryQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		clc = clc.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return clc
}

func (clc *CheckListCategoryQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CheckListCategoryQuery {
	return clc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (clcd *CheckListCategoryDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListCategoryDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		clcd = clcd.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return clcd
}

func (clcd *CheckListCategoryDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CheckListCategoryDefinitionQuery {
	return clcd
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (cli *CheckListItemQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListItemQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		cli = cli.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return cli
}

func (cli *CheckListItemQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CheckListItemQuery {
	return cli
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (clid *CheckListItemDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListItemDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		clid = clid.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return clid
}

func (clid *CheckListItemDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CheckListItemDefinitionQuery {
	return clid
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CommentQuery) CollectFields(ctx context.Context, satisfies ...string) *CommentQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return c
}

func (c *CommentQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CommentQuery {
	return c
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CustomerQuery) CollectFields(ctx context.Context, satisfies ...string) *CustomerQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return c
}

func (c *CustomerQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *CustomerQuery {
	return c
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (e *EquipmentQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		e = e.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return e
}

func (e *EquipmentQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "files":
			e = e.WithFiles(func(query *FileQuery) {
				query.collectField(ctx, field)
			})
		case "hyperlinks":
			e = e.WithHyperlinks(func(query *HyperlinkQuery) {
				query.collectField(ctx, field)
			})
		case "parentLocation":
			e = e.WithLocation(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "parentPosition":
			e = e.WithParentPosition(func(query *EquipmentPositionQuery) {
				query.collectField(ctx, field)
			})
		case "ports":
			e = e.WithPorts(func(query *EquipmentPortQuery) {
				query.collectField(ctx, field)
			})
		case "positions":
			e = e.WithPositions(func(query *EquipmentPositionQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			e = e.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		case "equipmentType":
			e = e.WithType(func(query *EquipmentTypeQuery) {
				query.collectField(ctx, field)
			})
		case "workOrder":
			e = e.WithWorkOrder(func(query *WorkOrderQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return e
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ec *EquipmentCategoryQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentCategoryQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ec = ec.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ec
}

func (ec *EquipmentCategoryQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentCategoryQuery {
	return ec
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPortQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ep = ep.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPortQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "definition":
			ep = ep.WithDefinition(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "serviceEndpoints":
			ep = ep.WithEndpoints(func(query *ServiceEndpointQuery) {
				query.collectField(ctx, field)
			})
		case "link":
			ep = ep.WithLink(func(query *LinkQuery) {
				query.collectField(ctx, field)
			})
		case "parentEquipment":
			ep = ep.WithParent(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			ep = ep.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return ep
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPortDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		epd = epd.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPortDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortDefinitionQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "portType":
			epd = epd.WithEquipmentPortType(func(query *EquipmentPortTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return epd
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ept *EquipmentPortTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ept = ept.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ept
}

func (ept *EquipmentPortTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "linkPropertyTypes":
			ept = ept.WithLinkPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		case "numberOfPortDefinitions":
			ept = ept.WithPortDefinitions(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "propertyTypes":
			ept = ept.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return ept
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPositionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPositionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ep = ep.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPositionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "attachedEquipment":
			ep = ep.WithAttachment(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		case "definition":
			ep = ep.WithDefinition(func(query *EquipmentPositionDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "parentEquipment":
			ep = ep.WithParent(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return ep
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPositionDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPositionDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		epd = epd.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPositionDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionDefinitionQuery {
	return epd
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (et *EquipmentTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		et = et.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return et
}

func (et *EquipmentTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *EquipmentTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "category":
			et = et.WithCategory(func(query *EquipmentCategoryQuery) {
				query.collectField(ctx, field)
			})
		case "equipments":
			et = et.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		case "portDefinitions":
			et = et.WithPortDefinitions(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "positionDefinitions":
			et = et.WithPositionDefinitions(func(query *EquipmentPositionDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "propertyTypes":
			et = et.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		case "serviceEndpointDefinitions":
			et = et.WithServiceEndpointDefinitions(func(query *ServiceEndpointDefinitionQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return et
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (f *FileQuery) CollectFields(ctx context.Context, satisfies ...string) *FileQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		f = f.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return f
}

func (f *FileQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *FileQuery {
	return f
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fp *FloorPlanQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		fp = fp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return fp
}

func (fp *FloorPlanQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *FloorPlanQuery {
	return fp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fprp *FloorPlanReferencePointQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanReferencePointQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		fprp = fprp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return fprp
}

func (fprp *FloorPlanReferencePointQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *FloorPlanReferencePointQuery {
	return fprp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fps *FloorPlanScaleQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanScaleQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		fps = fps.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return fps
}

func (fps *FloorPlanScaleQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *FloorPlanScaleQuery {
	return fps
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (h *HyperlinkQuery) CollectFields(ctx context.Context, satisfies ...string) *HyperlinkQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		h = h.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return h
}

func (h *HyperlinkQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *HyperlinkQuery {
	return h
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LinkQuery) CollectFields(ctx context.Context, satisfies ...string) *LinkQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		l = l.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return l
}

func (l *LinkQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *LinkQuery {
	return l
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LocationQuery) CollectFields(ctx context.Context, satisfies ...string) *LocationQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		l = l.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return l
}

func (l *LocationQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *LocationQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "cellData":
			l = l.WithCellScan(func(query *SurveyCellScanQuery) {
				query.collectField(ctx, field)
			})
		case "children":
			l = l.WithChildren(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "equipments":
			l = l.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		case "files", "images":
			l = l.WithFiles(func(query *FileQuery) {
				query.collectField(ctx, field)
			})
		case "floorPlans":
			l = l.WithFloorPlans(func(query *FloorPlanQuery) {
				query.collectField(ctx, field)
			})
		case "hyperlinks":
			l = l.WithHyperlinks(func(query *HyperlinkQuery) {
				query.collectField(ctx, field)
			})
		case "parentLocation":
			l = l.WithParent(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			l = l.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		case "surveys":
			l = l.WithSurvey(func(query *SurveyQuery) {
				query.collectField(ctx, field)
			})
		case "locationType":
			l = l.WithType(func(query *LocationTypeQuery) {
				query.collectField(ctx, field)
			})
		case "wifiData":
			l = l.WithWifiScan(func(query *SurveyWiFiScanQuery) {
				query.collectField(ctx, field)
			})
		case "workOrders":
			l = l.WithWorkOrders(func(query *WorkOrderQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return l
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (lt *LocationTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *LocationTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		lt = lt.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return lt
}

func (lt *LocationTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *LocationTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "locations":
			lt = lt.WithLocations(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "propertyTypes":
			lt = lt.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		case "surveyTemplateCategories":
			lt = lt.WithSurveyTemplateCategories(func(query *SurveyTemplateCategoryQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return lt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pp *PermissionsPolicyQuery) CollectFields(ctx context.Context, satisfies ...string) *PermissionsPolicyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pp = pp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pp
}

func (pp *PermissionsPolicyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PermissionsPolicyQuery {
	return pp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *ProjectQuery) CollectFields(ctx context.Context, satisfies ...string) *ProjectQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pr
}

func (pr *ProjectQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ProjectQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "comments":
			pr = pr.WithComments(func(query *CommentQuery) {
				query.collectField(ctx, field)
			})
		case "createdBy":
			pr = pr.WithCreator(func(query *UserQuery) {
				query.collectField(ctx, field)
			})
		case "location":
			pr = pr.WithLocation(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			pr = pr.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		case "type":
			pr = pr.WithType(func(query *ProjectTypeQuery) {
				query.collectField(ctx, field)
			})
		case "workOrders":
			pr = pr.WithWorkOrders(func(query *WorkOrderQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *ProjectTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *ProjectTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pt = pt.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pt
}

func (pt *ProjectTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ProjectTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "projects":
			pt = pt.WithProjects(func(query *ProjectQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			pt = pt.WithProperties(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		case "workOrders":
			pt = pt.WithWorkOrders(func(query *WorkOrderDefinitionQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return pt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *PropertyQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pr
}

func (pr *PropertyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PropertyQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "equipmentValue":
			pr = pr.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(ctx, field)
			})
		case "locationValue":
			pr = pr.WithLocation(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "serviceValue":
			pr = pr.WithService(func(query *ServiceQuery) {
				query.collectField(ctx, field)
			})
		case "propertyType":
			pr = pr.WithType(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *PropertyTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pt = pt.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pt
}

func (pt *PropertyTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PropertyTypeQuery {
	return pt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rf *ReportFilterQuery) CollectFields(ctx context.Context, satisfies ...string) *ReportFilterQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rf = rf.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rf
}

func (rf *ReportFilterQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ReportFilterQuery {
	return rf
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (s *ServiceQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		s = s.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return s
}

func (s *ServiceQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ServiceQuery {
	return s
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (se *ServiceEndpointQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceEndpointQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		se = se.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return se
}

func (se *ServiceEndpointQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ServiceEndpointQuery {
	return se
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (sed *ServiceEndpointDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceEndpointDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		sed = sed.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return sed
}

func (sed *ServiceEndpointDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ServiceEndpointDefinitionQuery {
	return sed
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (st *ServiceTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		st = st.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return st
}

func (st *ServiceTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ServiceTypeQuery {
	return st
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (s *SurveyQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		s = s.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return s
}

func (s *SurveyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyQuery {
	return s
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (scs *SurveyCellScanQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyCellScanQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		scs = scs.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return scs
}

func (scs *SurveyCellScanQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyCellScanQuery {
	return scs
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (sq *SurveyQuestionQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyQuestionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		sq = sq.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return sq
}

func (sq *SurveyQuestionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyQuestionQuery {
	return sq
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (stc *SurveyTemplateCategoryQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyTemplateCategoryQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		stc = stc.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return stc
}

func (stc *SurveyTemplateCategoryQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateCategoryQuery {
	return stc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (stq *SurveyTemplateQuestionQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyTemplateQuestionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		stq = stq.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return stq
}

func (stq *SurveyTemplateQuestionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateQuestionQuery {
	return stq
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (swfs *SurveyWiFiScanQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyWiFiScanQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		swfs = swfs.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return swfs
}

func (swfs *SurveyWiFiScanQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *SurveyWiFiScanQuery {
	return swfs
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (u *UserQuery) CollectFields(ctx context.Context, satisfies ...string) *UserQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		u = u.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return u
}

func (u *UserQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *UserQuery {
	return u
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ug *UsersGroupQuery) CollectFields(ctx context.Context, satisfies ...string) *UsersGroupQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		ug = ug.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return ug
}

func (ug *UsersGroupQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *UsersGroupQuery {
	return ug
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wo *WorkOrderQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		wo = wo.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return wo
}

func (wo *WorkOrderQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *WorkOrderQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "activities":
			wo = wo.WithActivities(func(query *ActivityQuery) {
				query.collectField(ctx, field)
			})
		case "assignedTo":
			wo = wo.WithAssignee(func(query *UserQuery) {
				query.collectField(ctx, field)
			})
		case "checkListCategories":
			wo = wo.WithCheckListCategories(func(query *CheckListCategoryQuery) {
				query.collectField(ctx, field)
			})
		case "comments":
			wo = wo.WithComments(func(query *CommentQuery) {
				query.collectField(ctx, field)
			})
		case "hyperlinks":
			wo = wo.WithHyperlinks(func(query *HyperlinkQuery) {
				query.collectField(ctx, field)
			})
		case "location":
			wo = wo.WithLocation(func(query *LocationQuery) {
				query.collectField(ctx, field)
			})
		case "owner":
			wo = wo.WithOwner(func(query *UserQuery) {
				query.collectField(ctx, field)
			})
		case "project":
			wo = wo.WithProject(func(query *ProjectQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			wo = wo.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		case "workOrderTemplate":
			wo = wo.WithTemplate(func(query *WorkOrderTemplateQuery) {
				query.collectField(ctx, field)
			})
		case "workOrderType":
			wo = wo.WithType(func(query *WorkOrderTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return wo
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wod *WorkOrderDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderDefinitionQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		wod = wod.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return wod
}

func (wod *WorkOrderDefinitionQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *WorkOrderDefinitionQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "type":
			wod = wod.WithType(func(query *WorkOrderTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return wod
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wot *WorkOrderTemplateQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderTemplateQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		wot = wot.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return wot
}

func (wot *WorkOrderTemplateQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *WorkOrderTemplateQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "checkListCategoryDefinitions":
			wot = wot.WithCheckListCategoryDefinitions(func(query *CheckListCategoryDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "propertyTypes":
			wot = wot.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return wot
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wot *WorkOrderTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		wot = wot.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return wot
}

func (wot *WorkOrderTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *WorkOrderTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "checkListCategoryDefinitions":
			wot = wot.WithCheckListCategoryDefinitions(func(query *CheckListCategoryDefinitionQuery) {
				query.collectField(ctx, field)
			})
		case "propertyTypes":
			wot = wot.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return wot
}
