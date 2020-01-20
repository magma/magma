// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"math"
	"sync"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"

	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

type locationTypeResolver struct{}

func (locationTypeResolver) PropertyTypes(ctx context.Context, obj *ent.LocationType) ([]*ent.PropertyType, error) {
	return obj.QueryPropertyTypes().All(ctx)
}

func (locationTypeResolver) NumberOfLocations(ctx context.Context, obj *ent.LocationType) (int, error) {
	return obj.QueryLocations().Count(ctx)
}

func (locationTypeResolver) Locations(ctx context.Context, typ *ent.LocationType, enforceHasLatLong *bool) (*ent.LocationConnection, error) {
	query := typ.QueryLocations()
	if enforceHasLatLong != nil && *enforceHasLatLong {
		query = query.Where(location.LatitudeNEQ(0), location.LongitudeNEQ(0))
	}
	return query.Paginate(ctx, nil, nil, nil, nil)
}

func (locationTypeResolver) SurveyTemplateCategories(ctx context.Context, obj *ent.LocationType) ([]*ent.SurveyTemplateCategory, error) {
	return obj.QuerySurveyTemplateCategories().All(ctx)
}

type locationResolver struct{}

func (r locationResolver) Surveys(ctx context.Context, obj *ent.Location) ([]*ent.Survey, error) {
	return obj.QuerySurvey().All(ctx)
}

func (r locationResolver) WifiData(ctx context.Context, obj *ent.Location) ([]*ent.SurveyWiFiScan, error) {
	return obj.QueryWifiScan().All(ctx)
}

func (r locationResolver) CellData(ctx context.Context, obj *ent.Location) ([]*ent.SurveyCellScan, error) {
	return obj.QueryCellScan().All(ctx)
}

func (locationResolver) LocationType(ctx context.Context, obj *ent.Location) (*ent.LocationType, error) {
	return obj.QueryType().Only(ctx)
}

func (r locationResolver) FloorPlans(ctx context.Context, obj *ent.Location) ([]*ent.FloorPlan, error) {
	return obj.QueryFloorPlans().All(ctx)
}

func (locationResolver) ParentLocation(ctx context.Context, obj *ent.Location) (*ent.Location, error) {
	parent, err := obj.QueryParent().Only(ctx)
	if ent.IsNotFound(err) {
		err = nil
	}
	return parent, err
}

func (locationResolver) Children(ctx context.Context, obj *ent.Location) ([]*ent.Location, error) {
	return obj.QueryChildren().All(ctx)
}

func (locationResolver) NumChildren(ctx context.Context, obj *ent.Location) (int, error) {
	return obj.QueryChildren().Count(ctx)
}

func (locationResolver) Equipments(ctx context.Context, obj *ent.Location) ([]*ent.Equipment, error) {
	return obj.QueryEquipment().All(ctx)
}

func (locationResolver) Properties(ctx context.Context, obj *ent.Location) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

func (locationResolver) Images(ctx context.Context, obj *ent.Location) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeImage.String())).All(ctx)
}

func (locationResolver) Files(ctx context.Context, obj *ent.Location) ([]*ent.File, error) {
	return obj.QueryFiles().Where(file.Type(models.FileTypeFile.String())).All(ctx)
}

type topologist struct {
	equipment sync.Map
	links     sync.Map
	sem       *semaphore.Weighted
	maxDepth  int
}

func (*topologist) rootNode(ctx context.Context, eq *ent.Equipment) *ent.Equipment {
	parent := eq
	for parent != nil {
		p, err := parent.QueryParentPosition().QueryParent().Only(ctx)
		if err != nil {
			break
		}

		parent = p
	}

	return parent
}

func (t *topologist) nestedNodes(ctx context.Context, eq *ent.Equipment, depth int) ([]*ent.Equipment, error) {
	if depth >= 5 {
		return nil, nil
	}

	posEqs, err := eq.QueryPositions().QueryAttachment().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed querying positions")
	}

	posEqs = append(posEqs, eq)
	for _, posEq := range posEqs {
		nestedEqs, err := t.nestedNodes(ctx, posEq, depth+1)
		if err != nil {
			return nil, err
		}
		posEqs = append(posEqs, nestedEqs...)
	}

	return posEqs, nil
}

func (*topologist) hkey(id1, id2 string) string {
	if id1 < id2 {
		return id1 + id2
	}
	return id2 + id1
}

func (t *topologist) build(ctx context.Context, eq *ent.Equipment, depth int) error {
	if err := t.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer t.sem.Release(1)

	t.equipment.Store(eq.ID, eq)
	if depth >= t.maxDepth {
		return nil
	}

	subTree, err := t.nestedNodes(ctx, eq, 0)
	if err != nil {
		return errors.Wrap(err, "failed querying nested equipment")
	}

	g := ctxgroup.WithContext(ctx)
	for _, neq := range subTree {
		leqs, err := neq.QueryPorts().
			QueryLink().
			QueryPorts().
			QueryParent().
			Where(equipment.IDNEQ(eq.ID)).
			All(ctx)
		if err != nil {
			return errors.Wrap(err, "querying equipment links")
		}

		for _, leq := range leqs {
			root := t.rootNode(ctx, leq)
			key := t.hkey(eq.ID, root.ID)
			value := &models.TopologyLink{Type: models.TopologyLinkTypePhysical, Source: eq, Target: root}
			if _, loaded := t.links.LoadOrStore(key, value); !loaded {
				g.Go(func(ctx context.Context) error {
					return t.build(ctx, root, depth+1)
				})
			}
		}
	}
	return g.Wait()
}

func (t *topologist) topology() *models.NetworkTopology {
	var nodes []ent.Noder
	t.equipment.Range(func(_, value interface{}) bool {
		nodes = append(nodes, value.(*ent.Equipment))
		return true
	})
	var links []*models.TopologyLink
	t.links.Range(func(_, value interface{}) bool {
		links = append(links, value.(*models.TopologyLink))
		return true
	})
	return &models.NetworkTopology{Nodes: nodes, Links: links}
}

// Need to deal with positions
func (locationResolver) Topology(ctx context.Context, loc *ent.Location, depth *int) (*models.NetworkTopology, error) {
	if depth == nil {
		return nil, errors.New("depth not supplied")
	}
	eqs, err := loc.QueryEquipment().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed querying location root equipment")
	}

	t := &topologist{
		sem:      semaphore.NewWeighted(32),
		maxDepth: *depth,
	}
	g := ctxgroup.WithContext(ctx)
	for _, eq := range eqs {
		eq := eq
		g.Go(func(ctx context.Context) error {
			return t.build(ctx, eq, 0)
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return t.topology(), nil
}

func (locationResolver) LocationHierarchy(ctx context.Context, l *ent.Location) ([]*ent.Location, error) {
	var locs []*ent.Location
	for l != nil {
		pl, err := l.QueryParent().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, errors.Wrapf(err, "querying parent location of %v", l.ID)
		}
		if pl != nil {
			locs = append([]*ent.Location{pl}, locs...)
		}
		l = pl
	}
	return locs, nil
}

func (locationResolver) DistanceKm(ctx context.Context, location *ent.Location, latitude float64, longitude float64) (float64, error) {
	const Radian = math.Pi / 180
	const EarthRadiusKm = 6371

	locLat, locLong := location.Latitude, location.Longitude
	a := 0.5 - math.Cos((latitude-locLat)*Radian)/2 +
		math.Cos(locLat*Radian)*math.Cos(latitude*Radian)*
			(1-math.Cos((longitude-locLong)*Radian))/2
	return EarthRadiusKm * 2 * math.Asin(math.Sqrt(a)), nil
}
