// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (ar *ActionsRuleQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ActionsRuleQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ar = ar.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ar
}

func (ar *ActionsRuleQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ActionsRuleQuery {
	return ar
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (cli *CheckListItemQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *CheckListItemQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		cli = cli.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return cli
}

func (cli *CheckListItemQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CheckListItemQuery {
	return cli
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (clid *CheckListItemDefinitionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *CheckListItemDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		clid = clid.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return clid
}

func (clid *CheckListItemDefinitionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CheckListItemDefinitionQuery {
	return clid
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CommentQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *CommentQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		c = c.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return c
}

func (c *CommentQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CommentQuery {
	return c
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CustomerQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *CustomerQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		c = c.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return c
}

func (c *CustomerQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *CustomerQuery {
	return c
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (e *EquipmentQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		e = e.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return e
}

func (e *EquipmentQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentQuery {
	return e
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (ec *EquipmentCategoryQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentCategoryQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ec = ec.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ec
}

func (ec *EquipmentCategoryQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentCategoryQuery {
	return ec
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPortQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentPortQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ep = ep.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPortQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortQuery {
	return ep
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPortDefinitionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentPortDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		epd = epd.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPortDefinitionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortDefinitionQuery {
	return epd
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (ept *EquipmentPortTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentPortTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ept = ept.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ept
}

func (ept *EquipmentPortTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPortTypeQuery {
	return ept
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (ep *EquipmentPositionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentPositionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		ep = ep.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return ep
}

func (ep *EquipmentPositionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionQuery {
	return ep
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (epd *EquipmentPositionDefinitionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentPositionDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		epd = epd.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return epd
}

func (epd *EquipmentPositionDefinitionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentPositionDefinitionQuery {
	return epd
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (et *EquipmentTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *EquipmentTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		et = et.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return et
}

func (et *EquipmentTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *EquipmentTypeQuery {
	return et
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (f *FileQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *FileQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		f = f.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return f
}

func (f *FileQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FileQuery {
	return f
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (fp *FloorPlanQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *FloorPlanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fp = fp.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fp
}

func (fp *FloorPlanQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanQuery {
	return fp
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (fprp *FloorPlanReferencePointQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *FloorPlanReferencePointQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fprp = fprp.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fprp
}

func (fprp *FloorPlanReferencePointQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanReferencePointQuery {
	return fprp
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (fps *FloorPlanScaleQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *FloorPlanScaleQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		fps = fps.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return fps
}

func (fps *FloorPlanScaleQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *FloorPlanScaleQuery {
	return fps
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (h *HyperlinkQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *HyperlinkQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		h = h.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return h
}

func (h *HyperlinkQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *HyperlinkQuery {
	return h
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LinkQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *LinkQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		l = l.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return l
}

func (l *LinkQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LinkQuery {
	return l
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (l *LocationQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *LocationQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		l = l.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return l
}

func (l *LocationQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LocationQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "cellData":
			l = l.WithCellScan(func(query *SurveyCellScanQuery) {
				query.withField(reqctx, field)
			})
		case "children":
			l = l.WithChildren(func(query *LocationQuery) {
				query.withField(reqctx, field)
			})
		case "equipments":
			l = l.WithEquipment(func(query *EquipmentQuery) {
				query.withField(reqctx, field)
			})
		case "files", "images":
			l = l.WithFiles(func(query *FileQuery) {
				query.withField(reqctx, field)
			})
		case "floorPlans":
			l = l.WithFloorPlans(func(query *FloorPlanQuery) {
				query.withField(reqctx, field)
			})
		case "hyperlinks":
			l = l.WithHyperlinks(func(query *HyperlinkQuery) {
				query.withField(reqctx, field)
			})
		case "parentLocation":
			l = l.WithParent(func(query *LocationQuery) {
				query.withField(reqctx, field)
			})
		case "properties":
			l = l.WithProperties(func(query *PropertyQuery) {
				query.withField(reqctx, field)
			})
		case "surveys":
			l = l.WithSurvey(func(query *SurveyQuery) {
				query.withField(reqctx, field)
			})
		case "locationType":
			l = l.WithType(func(query *LocationTypeQuery) {
				query.withField(reqctx, field)
			})
		case "wifiData":
			l = l.WithWifiScan(func(query *SurveyWiFiScanQuery) {
				query.withField(reqctx, field)
			})
		case "workOrders":
			l = l.WithWorkOrders(func(query *WorkOrderQuery) {
				query.withField(reqctx, field)
			})
		}
	}
	return l
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (lt *LocationTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *LocationTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		lt = lt.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return lt
}

func (lt *LocationTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *LocationTypeQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "locations":
			lt = lt.WithLocations(func(query *LocationQuery) {
				query.withField(reqctx, field)
			})
		case "propertyTypes":
			lt = lt.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.withField(reqctx, field)
			})
		case "surveyTemplateCategories":
			lt = lt.WithSurveyTemplateCategories(func(query *SurveyTemplateCategoryQuery) {
				query.withField(reqctx, field)
			})
		}
	}
	return lt
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *ProjectQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ProjectQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pr = pr.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pr
}

func (pr *ProjectQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ProjectQuery {
	return pr
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *ProjectTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ProjectTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pt = pt.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pt
}

func (pt *ProjectTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ProjectTypeQuery {
	return pt
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *PropertyQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *PropertyQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pr = pr.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pr
}

func (pr *PropertyQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *PropertyQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "equipmentValue":
			pr = pr.WithEquipment(func(query *EquipmentQuery) {
				query.withField(reqctx, field)
			})
		case "locationValue":
			pr = pr.WithLocation(func(query *LocationQuery) {
				query.withField(reqctx, field)
			})
		case "serviceValue":
			pr = pr.WithService(func(query *ServiceQuery) {
				query.withField(reqctx, field)
			})
		case "propertyType":
			pr = pr.WithType(func(query *PropertyTypeQuery) {
				query.withField(reqctx, field)
			})
		}
	}
	return pr
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *PropertyTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *PropertyTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		pt = pt.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return pt
}

func (pt *PropertyTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *PropertyTypeQuery {
	return pt
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (s *ServiceQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ServiceQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		s = s.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return s
}

func (s *ServiceQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceQuery {
	return s
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (se *ServiceEndpointQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ServiceEndpointQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		se = se.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return se
}

func (se *ServiceEndpointQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceEndpointQuery {
	return se
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (st *ServiceTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *ServiceTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		st = st.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return st
}

func (st *ServiceTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *ServiceTypeQuery {
	return st
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (s *SurveyQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		s = s.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return s
}

func (s *SurveyQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyQuery {
	return s
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (scs *SurveyCellScanQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyCellScanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		scs = scs.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return scs
}

func (scs *SurveyCellScanQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyCellScanQuery {
	return scs
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (sq *SurveyQuestionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyQuestionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		sq = sq.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return sq
}

func (sq *SurveyQuestionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyQuestionQuery {
	return sq
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (stc *SurveyTemplateCategoryQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyTemplateCategoryQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		stc = stc.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return stc
}

func (stc *SurveyTemplateCategoryQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateCategoryQuery {
	return stc
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (stq *SurveyTemplateQuestionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyTemplateQuestionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		stq = stq.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return stq
}

func (stq *SurveyTemplateQuestionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyTemplateQuestionQuery {
	return stq
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (swfs *SurveyWiFiScanQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *SurveyWiFiScanQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		swfs = swfs.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return swfs
}

func (swfs *SurveyWiFiScanQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *SurveyWiFiScanQuery {
	return swfs
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TechnicianQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *TechnicianQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		t = t.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return t
}

func (t *TechnicianQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *TechnicianQuery {
	return t
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (wo *WorkOrderQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *WorkOrderQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wo = wo.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wo
}

func (wo *WorkOrderQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderQuery {
	return wo
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (wod *WorkOrderDefinitionQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *WorkOrderDefinitionQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wod = wod.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wod
}

func (wod *WorkOrderDefinitionQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderDefinitionQuery {
	return wod
}

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (wot *WorkOrderTypeQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *WorkOrderTypeQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		wot = wot.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return wot
}

func (wot *WorkOrderTypeQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *WorkOrderTypeQuery {
	return wot
}
