// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/pkg/errors"
)

const maxFormSize = 5 << 20

func findIndex(a []string, x string) int {
	for i, n := range a {
		if strings.EqualFold(x, n) {
			return i
		}
	}
	return -1
}

func findStringContainsIndex(a []string, x string) int {
	for i, n := range a {
		if strings.Contains(n, x) {
			return i
		}
	}
	return -1
}

func sortSlice(a []int, acs bool) []int {
	sort.Slice(a, func(i, j int) bool {
		if acs {
			return a[i] < a[j]
		}
		return a[i] > a[j]
	})
	return a
}

func findIndexForSimilar(a []string, x string) int {
	i := findIndex(a, x)
	if i != -1 {
		return i
	}
	newX := strings.ReplaceAll(x, " ", "_")
	i = findIndex(a, newX)
	if i != -1 {
		return i
	}
	newX = strings.ReplaceAll(x, "_", " ")
	i = findIndex(a, newX)
	return i
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func getPropInput(propertyType ent.PropertyType, value string) (*models.PropertyInput, error) {
	typ := propertyType.Type
	switch typ {
	case "date", "email", "string", "enum":
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			StringValue:    &value,
		}, nil
	case "int":
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			IntValue:       &intVal,
		}, nil
	case "float":
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			FloatValue:     &floatVal,
		}, nil
	case "gps_location": // 45.6 , 67.89
		split := strings.Split(value, ",")
		if len(split) != 2 {
			return nil, errors.Errorf("gps location data isn't of form '<LAT>,<LONG>' %s", value)
		}
		lat, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil {
			return nil, err
		}
		long, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil {
			return nil, err
		}
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			LatitudeValue:  &lat,
			LongitudeValue: &long,
		}, nil
	case "range":
		split := strings.Split(value, "-")
		if len(split) != 2 {
			return nil, errors.Errorf("range data isn't of form '<FROM>-<TO>' %s", value)
		}

		from, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
		if err != nil {
			return nil, err
		}
		to, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
		if err != nil {
			return nil, err
		}
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			RangeFromValue: &from,
			RangeToValue:   &to,
		}, nil
	case "bool":
		var b bool
		b, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed parsing bool property %s", value)
		}
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			BooleanValue:   &b,
		}, nil
	case "location":
		if value != "" {
			return &models.PropertyInput{
				PropertyTypeID:  propertyType.ID,
				LocationIDValue: &value,
			}, nil
		} else {
			return &models.PropertyInput{
				PropertyTypeID:  propertyType.ID,
				LocationIDValue: nil,
			}, nil
		}
	case "equipment":
		if value != "" {
			return &models.PropertyInput{
				PropertyTypeID:   propertyType.ID,
				EquipmentIDValue: &value,
			}, nil
		} else {
			return &models.PropertyInput{
				PropertyTypeID:   propertyType.ID,
				EquipmentIDValue: nil,
			}, nil
		}
	case "service":
		if value != "" {
			return &models.PropertyInput{
				PropertyTypeID: propertyType.ID,
				ServiceIDValue: &value,
			}, nil
		} else {
			return &models.PropertyInput{
				PropertyTypeID: propertyType.ID,
				ServiceIDValue: nil,
			}, nil
		}
	default:
		return &models.PropertyInput{
			PropertyTypeID: propertyType.ID,
			StringValue:    &value,
		}, nil
	}
}

// Supports "pyramid shape CSV", in order to add properties/equipment to location level other than the smallest
func getLowestLocationHierarchyIdxForRow(ctx context.Context, line []string) int {
	for idx, column := range line {
		if idx > getImportContext(ctx).lowestHierarchyIndex || column == "" {
			return idx - 1
		}
	}
	return len(line) - 1
}

// Prepare map of locationTypes to indexes and the "lowestHirarchyIndex"
func (m *importer) populateIndexToLocationTypeMap(ctx context.Context, firstLine []string, populateLocationProperties bool) {
	indexToLocationTypeID := getImportContext(ctx).indexToLocationTypeID
	lowestHierarchyIndex := &getImportContext(ctx).lowestHierarchyIndex
	qr, ltr := m.r.Query(), m.r.LocationType()
	var idx int
	locTypes, _ := qr.LocationTypes(ctx, nil, nil, nil, nil)
	for _, typeEdge := range locTypes.Edges {
		typeName := typeEdge.Node.Name
		typeID := typeEdge.Node.ID
		idx = findIndex(firstLine, typeName)
		if idx != -1 {
			indexToLocationTypeID[idx] = typeID
			if idx > *lowestHierarchyIndex {
				*lowestHierarchyIndex = idx
			}
			if populateLocationProperties {
				typeIDsToProperties := getImportContext(ctx).typeIDsToProperties
				propNameToIndex := getImportContext(ctx).propNameToIndex
				properties, _ := ltr.PropertyTypes(ctx, typeEdge.Node)
				for _, prop := range properties {
					typeIDsToProperties[typeID] = append(typeIDsToProperties[typeID], prop.Name)
					propIdx := findIndex(firstLine, prop.Name)
					if propIdx != -1 {
						propNameToIndex[prop.Name] = propIdx
					}
				}
			}
		}
	}
}

func (m *importer) populateEquipmentTypeNameToIDMapGeneral(ctx context.Context, firstLine []string, populateEquipProperties bool) error {
	header, err := NewImportHeader(firstLine, ImportEntityEquipment)
	if err != nil {
		return err
	}
	return m.populateEquipmentTypeNameToIDMap(ctx, header, populateEquipProperties)
}

func (m *importer) populateEquipmentTypeNameToIDMap(ctx context.Context, firstLine ImportHeader, populateEquipProperties bool) error {
	equipmentTypeNameToID := getImportContext(ctx).equipmentTypeNameToID
	qr := m.r.Query()
	equipTypes, err := qr.EquipmentTypes(ctx, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, equipTypeEdge := range equipTypes.Edges {
		equipTypeNode := equipTypeEdge.Node
		equipmentTypeNameToID[equipTypeNode.Name] = equipTypeNode.ID
		if populateEquipProperties {
			propNameToIndex := getImportContext(ctx).propNameToIndex
			equipmentTypeIDToProperties := getImportContext(ctx).equipmentTypeIDToProperties
			eqr := m.r.EquipmentType()
			properties, err := eqr.PropertyTypes(ctx, equipTypeNode)
			if err != nil {
				return err
			}
			for _, prop := range properties {
				equipmentTypeIDToProperties[equipTypeNode.ID] = append(equipmentTypeIDToProperties[equipTypeNode.ID], prop.Name)
				propIdx := firstLine.Find(prop.Name)
				if propIdx != -1 {
					propNameToIndex[prop.Name] = propIdx
				}
			}
		}
	}
	return nil
}

type reader struct {
	*csv.Reader
	io.Closer
}

func (importer) charsetReader(r io.Reader, headers ...textproto.MIMEHeader) (io.Reader, error) {
	for _, header := range headers {
		if label := header.Get(textproto.CanonicalMIMEHeaderKey("X-MIME-Charset")); label != "" {
			if rd, err := charset.NewReaderLabel(label, r); err == nil {
				return rd, nil
			}
		}
		if ct := header.Get(textproto.CanonicalMIMEHeaderKey("Content-Type")); ct != "" {
			if _, _, ok := charset.DetermineEncoding(nil, ct); ok {
				return charset.NewReader(r, ct)
			}
		}
	}
	return nil, errors.New("no charset in headers")
}

func (m *importer) newReader(key string, req *http.Request) ([]string, *reader, error) {
	f, hdr, err := req.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	r, err := m.charsetReader(f, hdr.Header, textproto.MIMEHeader(req.Header))
	if err != nil {
		m.logger.For(req.Context()).Warn("cannot detect mime charset", zap.Error(err))
		r, err = charset.NewReader(f, req.Header.Get("Content-Type"))
	}
	if err != nil {
		return nil, nil, err
	}
	rd := csv.NewReader(transform.NewReader(r, unicode.BOMOverride(transform.Nop)))
	line, err := rd.Read()
	if err != nil {
		return nil, nil, err
	}
	return line, &reader{Reader: rd, Closer: f}, nil
}

func (m *importer) ClientFrom(ctx context.Context) *ent.Client {
	client := ent.FromContext(ctx)
	if client == nil {
		panic("no client attached to context")
	}
	return client
}

func (m *importer) getOrCreatePropTypeForLocation(ctx context.Context, lTypeID string, pname string) (*ent.PropertyType, error) {
	ptype, err := m.ClientFrom(ctx).LocationType.Query().
		Where(locationtype.ID(lTypeID)).
		QueryPropertyTypes().
		Where(propertytype.Name(pname)).
		Only(ctx)
	if ent.IsNotFound(err) {
		ptype, err = m.ClientFrom(ctx).PropertyType.Create().
			SetName(pname).
			// TODO T40408163. get "Type" from model
			SetType("string").
			SetLocationTypeID(lTypeID).
			Save(ctx)
	}
	return ptype, err
}

func (m *importer) getOrCreatePropTypeForEquipment(ctx context.Context, eTypeID string, pname string) (*ent.PropertyType, error) {
	ptype, err := m.ClientFrom(ctx).EquipmentType.Query().
		Where(equipmenttype.ID(eTypeID)).
		QueryPropertyTypes().
		Where(propertytype.Name(pname)).
		Only(ctx)
	if ent.IsNotFound(err) {
		ptype, err = m.ClientFrom(ctx).PropertyType.Create().
			SetName(pname).
			// TODO T40408163. get "Type" from model
			SetType("string").
			SetEquipmentTypeID(eTypeID).
			Save(ctx)
	}
	return ptype, err
}

func (m *importer) trimLine(line []string) []string {
	for i, value := range line {
		line[i] = strings.Trim(value, " ")
	}
	return line
}

func errorReturn(w http.ResponseWriter, msg string, log *zap.Logger, err error) {
	log.Warn(msg, zap.Error(err))
	if err == nil {
		http.Error(w, msg, http.StatusBadRequest)
	} else {
		http.Error(w, fmt.Sprintf("%s %q", msg, err), http.StatusBadRequest)
	}
}
