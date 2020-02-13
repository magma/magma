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
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ar = ar.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ar
}

func (ar *ActionsRuleQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ActionsRuleQuery {
	return ar
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (cli *CheckListItemQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListItemQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		cli = cli.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return cli
}

func (cli *CheckListItemQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CheckListItemQuery {
	return cli
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (clid *CheckListItemDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *CheckListItemDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		clid = clid.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return clid
}

func (clid *CheckListItemDefinitionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CheckListItemDefinitionQuery {
	return clid
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CommentQuery) CollectFields(ctx context.Context, satisfies ...string) *CommentQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		c = c.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return c
}

func (c *CommentQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CommentQuery {
	return c
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CustomerQuery) CollectFields(ctx context.Context, satisfies ...string) *CustomerQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		c = c.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return c
}

func (c *CustomerQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CustomerQuery {
	return c
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (e *EquipmentQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		e = e.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return e
}

func (e *EquipmentQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "files":
			e = e.WithFiles(func(query *FileQuery) {
				query.collectField(reqctx, field)
			})
		case "hyperlinks":
			e = e.WithHyperlinks(func(query *HyperlinkQuery) {
				query.collectField(reqctx, field)
			})
		case "parentLocation":
			e = e.WithLocation(func(query *LocationQuery) {
				query.collectField(reqctx, field)
			})
		case "parentPosition":
			e = e.WithParentPosition(func(query *EquipmentPositionQuery) {
				query.collectField(reqctx, field)
			})
		case "ports":
			e = e.WithPorts(func(query *EquipmentPortQuery) {
				query.collectField(reqctx, field)
			})
		case "positions":
			e = e.WithPositions(func(query *EquipmentPositionQuery) {
				query.collectField(reqctx, field)
			})
		case "properties":
			e = e.WithProperties(func(query *PropertyQuery) {
				query.collectField(reqctx, field)
			})
		case "equipmentType":
			e = e.WithType(func(query *EquipmentTypeQuery) {
				query.collectField(reqctx, field)
			})
		case "workOrder":
			e = e.WithWorkOrder(func(query *WorkOrderQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return e
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ec *EquipmentCategoryQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentCategoryQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ec = ec.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ec
}

func (ec *EquipmentCategoryQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentCategoryQuery {
	return ec
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPortQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ep = ep.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPortQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "definition":
			ep = ep.WithDefinition(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(reqctx, field)
			})
		case "serviceEndpoints":
			ep = ep.WithEndpoints(func(query *ServiceEndpointQuery) {
				query.collectField(reqctx, field)
			})
		case "link":
			ep = ep.WithLink(func(query *LinkQuery) {
				query.collectField(reqctx, field)
			})
		case "parentEquipment":
			ep = ep.WithParent(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		case "properties":
			ep = ep.WithProperties(func(query *PropertyQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return ep
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPortDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		epd = epd.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPortDefinitionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortDefinitionQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "portType":
			epd = epd.WithEquipmentPortType(func(query *EquipmentPortTypeQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return epd
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ept *EquipmentPortTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPortTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ept = ept.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ept
}

func (ept *EquipmentPortTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortTypeQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "linkPropertyTypes":
			ept = ept.WithLinkPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(reqctx, field)
			})
		case "numberOfPortDefinitions":
			ept = ept.WithPortDefinitions(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(reqctx, field)
			})
		case "propertyTypes":
			ept = ept.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return ept
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPositionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPositionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ep = ep.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPositionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "attachedEquipment":
			ep = ep.WithAttachment(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		case "definition":
			ep = ep.WithDefinition(func(query *EquipmentPositionDefinitionQuery) {
				query.collectField(reqctx, field)
			})
		case "parentEquipment":
			ep = ep.WithParent(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return ep
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPositionDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentPositionDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		epd = epd.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPositionDefinitionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionDefinitionQuery {
	return epd
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (et *EquipmentTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *EquipmentTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		et = et.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return et
}

func (et *EquipmentTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentTypeQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "category":
			et = et.WithCategory(func(query *EquipmentCategoryQuery) {
				query.collectField(reqctx, field)
			})
		case "equipments":
			et = et.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		case "portDefinitions":
			et = et.WithPortDefinitions(func(query *EquipmentPortDefinitionQuery) {
				query.collectField(reqctx, field)
			})
		case "positionDefinitions":
			et = et.WithPositionDefinitions(func(query *EquipmentPositionDefinitionQuery) {
				query.collectField(reqctx, field)
			})
		case "propertyTypes":
			et = et.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return et
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (f *FileQuery) CollectFields(ctx context.Context, satisfies ...string) *FileQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		f = f.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return f
}

func (f *FileQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FileQuery {
	return f
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fp *FloorPlanQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fp = fp.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fp
}

func (fp *FloorPlanQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanQuery {
	return fp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fprp *FloorPlanReferencePointQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanReferencePointQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fprp = fprp.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fprp
}

func (fprp *FloorPlanReferencePointQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanReferencePointQuery {
	return fprp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (fps *FloorPlanScaleQuery) CollectFields(ctx context.Context, satisfies ...string) *FloorPlanScaleQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fps = fps.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fps
}

func (fps *FloorPlanScaleQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanScaleQuery {
	return fps
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (h *HyperlinkQuery) CollectFields(ctx context.Context, satisfies ...string) *HyperlinkQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		h = h.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return h
}

func (h *HyperlinkQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *HyperlinkQuery {
	return h
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LinkQuery) CollectFields(ctx context.Context, satisfies ...string) *LinkQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		l = l.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return l
}

func (l *LinkQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LinkQuery {
	return l
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LocationQuery) CollectFields(ctx context.Context, satisfies ...string) *LocationQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		l = l.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return l
}

func (l *LocationQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LocationQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "cellData":
			l = l.WithCellScan(func(query *SurveyCellScanQuery) {
				query.collectField(reqctx, field)
			})
		case "children":
			l = l.WithChildren(func(query *LocationQuery) {
				query.collectField(reqctx, field)
			})
		case "equipments":
			l = l.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		case "files", "images":
			l = l.WithFiles(func(query *FileQuery) {
				query.collectField(reqctx, field)
			})
		case "floorPlans":
			l = l.WithFloorPlans(func(query *FloorPlanQuery) {
				query.collectField(reqctx, field)
			})
		case "hyperlinks":
			l = l.WithHyperlinks(func(query *HyperlinkQuery) {
				query.collectField(reqctx, field)
			})
		case "parentLocation":
			l = l.WithParent(func(query *LocationQuery) {
				query.collectField(reqctx, field)
			})
		case "properties":
			l = l.WithProperties(func(query *PropertyQuery) {
				query.collectField(reqctx, field)
			})
		case "surveys":
			l = l.WithSurvey(func(query *SurveyQuery) {
				query.collectField(reqctx, field)
			})
		case "locationType":
			l = l.WithType(func(query *LocationTypeQuery) {
				query.collectField(reqctx, field)
			})
		case "wifiData":
			l = l.WithWifiScan(func(query *SurveyWiFiScanQuery) {
				query.collectField(reqctx, field)
			})
		case "workOrders":
			l = l.WithWorkOrders(func(query *WorkOrderQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return l
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (lt *LocationTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *LocationTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		lt = lt.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return lt
}

func (lt *LocationTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LocationTypeQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "locations":
			lt = lt.WithLocations(func(query *LocationQuery) {
				query.collectField(reqctx, field)
			})
		case "propertyTypes":
			lt = lt.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(reqctx, field)
			})
		case "surveyTemplateCategories":
			lt = lt.WithSurveyTemplateCategories(func(query *SurveyTemplateCategoryQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return lt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *ProjectQuery) CollectFields(ctx context.Context, satisfies ...string) *ProjectQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pr = pr.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pr
}

func (pr *ProjectQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ProjectQuery {
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *ProjectTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *ProjectTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pt = pt.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pt
}

func (pt *ProjectTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ProjectTypeQuery {
	return pt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *PropertyQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pr = pr.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pr
}

func (pr *PropertyQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *PropertyQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "equipmentValue":
			pr = pr.WithEquipment(func(query *EquipmentQuery) {
				query.collectField(reqctx, field)
			})
		case "locationValue":
			pr = pr.WithLocation(func(query *LocationQuery) {
				query.collectField(reqctx, field)
			})
		case "serviceValue":
			pr = pr.WithService(func(query *ServiceQuery) {
				query.collectField(reqctx, field)
			})
		case "propertyType":
			pr = pr.WithType(func(query *PropertyTypeQuery) {
				query.collectField(reqctx, field)
			})
		}
	}
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *PropertyTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pt = pt.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pt
}

func (pt *PropertyTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *PropertyTypeQuery {
	return pt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (s *ServiceQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		s = s.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return s
}

func (s *ServiceQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceQuery {
	return s
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (se *ServiceEndpointQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceEndpointQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		se = se.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return se
}

func (se *ServiceEndpointQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceEndpointQuery {
	return se
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (st *ServiceTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *ServiceTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		st = st.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return st
}

func (st *ServiceTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceTypeQuery {
	return st
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (s *SurveyQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		s = s.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return s
}

func (s *SurveyQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyQuery {
	return s
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (scs *SurveyCellScanQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyCellScanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		scs = scs.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return scs
}

func (scs *SurveyCellScanQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyCellScanQuery {
	return scs
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (sq *SurveyQuestionQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyQuestionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		sq = sq.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return sq
}

func (sq *SurveyQuestionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyQuestionQuery {
	return sq
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (stc *SurveyTemplateCategoryQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyTemplateCategoryQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		stc = stc.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return stc
}

func (stc *SurveyTemplateCategoryQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateCategoryQuery {
	return stc
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (stq *SurveyTemplateQuestionQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyTemplateQuestionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		stq = stq.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return stq
}

func (stq *SurveyTemplateQuestionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateQuestionQuery {
	return stq
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (swfs *SurveyWiFiScanQuery) CollectFields(ctx context.Context, satisfies ...string) *SurveyWiFiScanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		swfs = swfs.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return swfs
}

func (swfs *SurveyWiFiScanQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyWiFiScanQuery {
	return swfs
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TechnicianQuery) CollectFields(ctx context.Context, satisfies ...string) *TechnicianQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		t = t.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return t
}

func (t *TechnicianQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *TechnicianQuery {
	return t
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wo *WorkOrderQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wo = wo.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wo
}

func (wo *WorkOrderQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderQuery {
	return wo
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wod *WorkOrderDefinitionQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wod = wod.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wod
}

func (wod *WorkOrderDefinitionQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderDefinitionQuery {
	return wod
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (wot *WorkOrderTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *WorkOrderTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wot = wot.collectField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wot
}

func (wot *WorkOrderTypeQuery) collectField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderTypeQuery {
	return wot
}
