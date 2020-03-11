// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/AlekSi/pointer"
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

func (locationTypeResolver) PropertyTypes(ctx context.Context, typ *ent.LocationType) ([]*ent.PropertyType, error) {
	if types, err := typ.Edges.PropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return types, err
	}
	return typ.QueryPropertyTypes().All(ctx)
}

func (locationTypeResolver) NumberOfLocations(ctx context.Context, obj *ent.LocationType) (int, error) {
	return obj.QueryLocations().Count(ctx)
}

func (locationTypeResolver) Locations(ctx context.Context, typ *ent.LocationType, enforceHasLatLong *bool) (*ent.LocationConnection, error) {
	query := typ.QueryLocations()
	if pointer.GetBool(enforceHasLatLong) {
		query = query.Where(location.LatitudeNEQ(0), location.LongitudeNEQ(0))
	}
	return query.Paginate(ctx, nil, nil, nil, nil)
}

func (locationTypeResolver) SurveyTemplateCategories(ctx context.Context, typ *ent.LocationType) ([]*ent.SurveyTemplateCategory, error) {
	if stc, err := typ.Edges.SurveyTemplateCategoriesOrErr(); !ent.IsNotLoaded(err) {
		return stc, err
	}
	return typ.QuerySurveyTemplateCategories().All(ctx)
}

type locationResolver struct{}

func (r locationResolver) Surveys(ctx context.Context, location *ent.Location) ([]*ent.Survey, error) {
	if surveys, err := location.Edges.SurveyOrErr(); !ent.IsNotLoaded(err) {
		return surveys, err
	}
	return location.QuerySurvey().All(ctx)
}

func (r locationResolver) WifiData(ctx context.Context, location *ent.Location) ([]*ent.SurveyWiFiScan, error) {
	if scans, err := location.Edges.WifiScanOrErr(); !ent.IsNotLoaded(err) {
		return scans, err
	}
	return location.QueryWifiScan().All(ctx)
}

func (r locationResolver) CellData(ctx context.Context, location *ent.Location) ([]*ent.SurveyCellScan, error) {
	if scans, err := location.Edges.CellScanOrErr(); !ent.IsNotLoaded(err) {
		return scans, err
	}
	return location.QueryCellScan().All(ctx)
}

func (locationResolver) LocationType(ctx context.Context, location *ent.Location) (*ent.LocationType, error) {
	if typ, err := location.Edges.TypeOrErr(); !ent.IsNotLoaded(err) {
		return typ, err
	}
	return location.QueryType().Only(ctx)
}

func (locationResolver) FloorPlans(ctx context.Context, location *ent.Location) ([]*ent.FloorPlan, error) {
	if plans, err := location.Edges.FloorPlansOrErr(); !ent.IsNotLoaded(err) {
		return plans, err
	}
	return location.QueryFloorPlans().All(ctx)
}

func (locationResolver) ParentLocation(ctx context.Context, location *ent.Location) (*ent.Location, error) {
	parent, err := location.Edges.ParentOrErr()
	if ent.IsNotLoaded(err) {
		parent, err = location.QueryParent().Only(ctx)
	}
	return parent, ent.MaskNotFound(err)
}

func (locationResolver) Children(ctx context.Context, location *ent.Location) ([]*ent.Location, error) {
	if children, err := location.Edges.ChildrenOrErr(); !ent.IsNotLoaded(err) {
		return children, err
	}
	return location.QueryChildren().All(ctx)
}

func (locationResolver) NumChildren(ctx context.Context, location *ent.Location) (int, error) {
	if children, err := location.Edges.ChildrenOrErr(); !ent.IsNotLoaded(err) {
		return len(children), err
	}
	return location.QueryChildren().Count(ctx)
}

func (locationResolver) Equipments(ctx context.Context, location *ent.Location) ([]*ent.Equipment, error) {
	if eqs, err := location.Edges.EquipmentOrErr(); !ent.IsNotLoaded(err) {
		return eqs, err
	}
	return location.QueryEquipment().All(ctx)
}

func (locationResolver) Properties(ctx context.Context, location *ent.Location) ([]*ent.Property, error) {
	if properties, err := location.Edges.PropertiesOrErr(); !ent.IsNotLoaded(err) {
		return properties, err
	}
	return location.QueryProperties().All(ctx)
}

func (locationResolver) filesOfType(ctx context.Context, location *ent.Location, typ string) ([]*ent.File, error) {
	fds, err := location.Edges.FilesOrErr()
	if ent.IsNotLoaded(err) {
		return location.QueryFiles().
			Where(file.Type(typ)).
			All(ctx)
	}
	files := make([]*ent.File, 0, len(fds))
	for _, f := range fds {
		if f.Type == typ {
			files = append(files, f)
		}
	}
	return files, nil
}

func (r locationResolver) Images(ctx context.Context, location *ent.Location) ([]*ent.File, error) {
	return r.filesOfType(ctx, location, models.FileTypeImage.String())
}

func (r locationResolver) Files(ctx context.Context, location *ent.Location) ([]*ent.File, error) {
	return r.filesOfType(ctx, location, models.FileTypeFile.String())
}

func (locationResolver) Hyperlinks(ctx context.Context, location *ent.Location) ([]*ent.Hyperlink, error) {
	if hls, err := location.Edges.HyperlinksOrErr(); !ent.IsNotLoaded(err) {
		return hls, err
	}
	return location.QueryHyperlinks().All(ctx)
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

func (*topologist) hkey(id1, id2 int) string {
	if id2 > id1 {
		id1, id2 = id2, id1
	}
	return strconv.Itoa(id1) + ":" + strconv.Itoa(id2)
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

func (locationResolver) LocationHierarchy(ctx context.Context, location *ent.Location) ([]*ent.Location, error) {
	var locations []*ent.Location
	for {
		parent, err := location.QueryParent().Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return locations, nil
			}
			return nil, fmt.Errorf("querying parent location: %w", err)
		}
		locations = append([]*ent.Location{parent}, locations...)
		location = parent
	}
}

func (locationResolver) DistanceKm(_ context.Context, location *ent.Location, latitude, longitude float64) (float64, error) {
	const (
		radian        = math.Pi / 180
		earthRadiusKm = 6371
	)
	locLat, locLong := location.Latitude, location.Longitude
	a := 0.5 - math.Cos((latitude-locLat)*radian)/2 +
		math.Cos(locLat*radian)*math.Cos(latitude*radian)*
			(1-math.Cos((longitude-locLong)*radian))/2
	return earthRadiusKm * 2 * math.Asin(math.Sqrt(a)), nil
}
