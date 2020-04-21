// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"

	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/gqlerror"
	"go.uber.org/zap"
)

type mutationResolver struct{ resolver }

func (mutationResolver) Me(ctx context.Context) *viewer.Viewer {
	return viewer.FromContext(ctx)
}

const badID = -1

func (mutationResolver) isEmptyProp(ptype *ent.PropertyType, input interface{}) (bool, error) {
	var (
		typ                           models.PropertyKind
		strVal                        *string
		boolVal                       *bool
		lat, long, rangeTo, rangeFrom *float64
	)
	switch v := input.(type) {
	case *models.PropertyInput:
		typ = models.PropertyKind(ptype.Type)
		strVal = v.StringValue
		boolVal = v.BooleanValue
		lat, long = v.LatitudeValue, v.LongitudeValue
		rangeTo, rangeFrom = v.RangeToValue, v.RangeFromValue
	case *models.PropertyTypeInput:
		typ = v.Type
		strVal = v.StringValue
		boolVal = v.BooleanValue
		lat, long = v.LatitudeValue, v.LongitudeValue
		rangeTo, rangeFrom = v.RangeToValue, v.RangeFromValue
	default:
		return false, errors.New("input not of type property or property type")
	}
	switch typ {
	case models.PropertyKindDate,
		models.PropertyKindEmail,
		models.PropertyKindString,
		models.PropertyKindEnum,
		models.PropertyKindDatetimeLocal:
		return strVal == nil || *strVal == "", nil
	case models.PropertyKindInt:
		// TODO detect int no-value
		return false, nil
	case models.PropertyKindGpsLocation:
		if lat == nil || long == nil {
			return false, errors.New("gps_location type, with no lat/long provided")
		}
		return *lat == 0 && *long == 0, nil
	case models.PropertyKindRange:
		if rangeTo == nil || rangeFrom == nil {
			return false, gqlerror.Errorf("range type, with no to/from provided")
		}
		return *rangeTo == 0 && *rangeFrom == 0, nil
	case models.PropertyKindBool:
		return boolVal == nil, nil
	default:
		return false, nil
	}
}

func (r mutationResolver) setNodePropertyCreate(ctx context.Context, setter *ent.PropertyCreate, nodeID int) error {
	client := r.ClientFrom(ctx)
	value, err := client.Node(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("invalid node id: %d %w", nodeID, err)
	}
	switch value.Type {
	case ent.TypeEquipment:
		setter.SetEquipmentValueID(value.ID)
	case ent.TypeLocation:
		setter.SetLocationValueID(value.ID)
	case ent.TypeService:
		setter.SetServiceValueID(value.ID)
	case ent.TypeWorkOrder:
		setter.SetWorkOrderValueID(value.ID)
	default:
		return fmt.Errorf("invalid node type: %d %s", value.ID, value.Type)
	}
	return nil
}

func (r mutationResolver) setNodePropertyUpdate(ctx context.Context, setter *ent.PropertyUpdate, nodeID int) error {
	client := r.ClientFrom(ctx)
	value, err := client.Node(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("invalid node id: %d %w", nodeID, err)
	}
	switch value.Type {
	case ent.TypeEquipment:
		setter.SetEquipmentValueID(value.ID)
	case ent.TypeLocation:
		setter.SetLocationValueID(value.ID)
	case ent.TypeService:
		setter.SetServiceValueID(value.ID)
	case ent.TypeWorkOrder:
		setter.SetWorkOrderValueID(value.ID)
	default:
		return fmt.Errorf("invalid node type: %d %s", value.ID, value.Type)
	}
	return nil
}

func (r mutationResolver) AddProperty(
	input *models.PropertyInput,
	args resolverutil.AddPropertyArgs,
) (*ent.Property, error) {
	ctx := args.Context
	client := r.ClientFrom(ctx)
	propType, err := client.PropertyType.Get(ctx, input.PropertyTypeID)
	if err != nil {
		return nil, err
	}
	isTemplate := args.IsTemplate != nil && *args.IsTemplate
	if !isTemplate && !propType.IsInstanceProperty {
		return nil, nil
	}
	query := client.Property.Create()
	if args.EntSetter != nil {
		args.EntSetter(query)
	}

	query = query.
		SetTypeID(input.PropertyTypeID).
		SetNillableStringVal(input.StringValue).
		SetNillableIntVal(input.IntValue).
		SetNillableBoolVal(input.BooleanValue).
		SetNillableFloatVal(input.FloatValue).
		SetNillableLatitudeVal(input.LatitudeValue).
		SetNillableLongitudeVal(input.LongitudeValue).
		SetNillableRangeFromVal(input.RangeFromValue).
		SetNillableRangeToVal(input.RangeToValue)

	if input.NodeIDValue != nil {
		if err = r.setNodePropertyCreate(ctx, query, *input.NodeIDValue); err != nil {
			return nil, err
		}
	}

	p, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating property: %w", err)
	}
	return p, nil
}

func (r mutationResolver) AddProperties(inputs []*models.PropertyInput, args resolverutil.AddPropertyArgs) ([]*ent.Property, error) {
	properties := make([]*ent.Property, 0, len(inputs))
	for _, input := range inputs {
		p, err := r.AddProperty(input, args)
		if err != nil {
			return nil, err
		}
		if p != nil {
			properties = append(properties, p)
		}
	}
	return properties, nil
}

func (r mutationResolver) AddPropertyTypes(
	ctx context.Context, inputs ...*models.PropertyTypeInput,
) ([]*ent.PropertyType, error) {
	var (
		client = r.ClientFrom(ctx).PropertyType
		types  = make([]*ent.PropertyType, len(inputs))
		err    error
	)
	for i, input := range inputs {
		if types[i], err = client.Create().
			SetName(input.Name).
			SetType(input.Type.String()).
			SetNillableNodeType(input.NodeType).
			SetNillableExternalID(input.ExternalID).
			SetNillableIndex(input.Index).
			SetNillableCategory(input.Category).
			SetNillableStringVal(input.StringValue).
			SetNillableIntVal(input.IntValue).
			SetNillableBoolVal(input.BooleanValue).
			SetNillableFloatVal(input.FloatValue).
			SetNillableLatitudeVal(input.LatitudeValue).
			SetNillableLongitudeVal(input.LongitudeValue).
			SetNillableIsInstanceProperty(input.IsInstanceProperty).
			SetNillableRangeFromVal(input.RangeFromValue).
			SetNillableRangeToVal(input.RangeToValue).
			SetNillableEditable(input.IsEditable).
			SetNillableMandatory(input.IsMandatory).
			SetNillableDeleted(input.IsDeleted).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating property type: %w", err)
		}
	}
	return types, nil
}

func (r mutationResolver) AddSurveyTemplateCategories(
	ctx context.Context, inputs ...*models.SurveyTemplateCategoryInput,
) ([]*ent.SurveyTemplateCategory, error) {
	var (
		client     = r.ClientFrom(ctx).SurveyTemplateCategory
		categories = make([]*ent.SurveyTemplateCategory, len(inputs))
	)
	for i, input := range inputs {
		questions, err := r.AddSurveyTemplateQuestions(ctx, input.SurveyTemplateQuestions...)
		if err != nil {
			return nil, err
		}
		if categories[i], err = client.Create().
			SetCategoryTitle(input.CategoryTitle).
			SetCategoryDescription(input.CategoryDescription).
			AddSurveyTemplateQuestions(questions...).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating survey template categories: %w", err)
		}
	}
	return categories, nil
}

func (r mutationResolver) AddSurveyTemplateQuestions(ctx context.Context, inputs ...*models.SurveyTemplateQuestionInput) ([]*ent.SurveyTemplateQuestion, error) {
	var (
		client    = r.ClientFrom(ctx).SurveyTemplateQuestion
		questions = make([]*ent.SurveyTemplateQuestion, len(inputs))
		err       error
	)
	for i, input := range inputs {
		if questions[i], err = client.Create().
			SetQuestionTitle(input.QuestionTitle).
			SetQuestionDescription(input.QuestionDescription).
			SetQuestionType(input.QuestionType.String()).
			SetIndex(input.Index).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating survey template questions: %w", err)
		}
	}
	return questions, nil
}

type ScanParentIDs struct {
	qid             *int
	locationID      *int
	checklistItemID *int
}

func validateScanParentIDs(ids ScanParentIDs) error {
	if ids.qid == nil && ids.locationID == nil && ids.checklistItemID == nil {
		return errors.New("must specify either question, location or checklist item")
	}

	return nil
}

func (r mutationResolver) AddWiFiScans(ctx context.Context, data []*models.SurveyWiFiScanData, locationID int) ([]*ent.SurveyWiFiScan, error) {
	return r.CreateWiFiScans(ctx, data, ScanParentIDs{locationID: &locationID})
}

func (r mutationResolver) CreateWiFiScans(ctx context.Context, inputs []*models.SurveyWiFiScanData, parentIDs ScanParentIDs) ([]*ent.SurveyWiFiScan, error) {
	if err := validateScanParentIDs(parentIDs); err != nil {
		return nil, err
	}
	var (
		client = r.ClientFrom(ctx).SurveyWiFiScan
		scans  = make([]*ent.SurveyWiFiScan, len(inputs))
		err    error
	)
	for i, input := range inputs {
		if scans[i], err = client.Create().
			SetTimestamp(time.Unix(int64(input.Timestamp), 0)).
			SetFrequency(input.Frequency).
			SetChannel(input.Channel).
			SetBssid(input.Bssid).
			SetStrength(input.Strength).
			SetNillableSsid(input.Ssid).
			SetNillableBand(input.Band).
			SetNillableChannelWidth(input.ChannelWidth).
			SetNillableCapabilities(input.Capabilities).
			SetNillableLatitude(input.Latitude).
			SetNillableLongitude(input.Longitude).
			SetNillableSurveyQuestionID(parentIDs.qid).
			SetNillableLocationID(parentIDs.locationID).
			SetNillableChecklistItemID(parentIDs.checklistItemID).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating survey wifi scan: %w", err)
		}
	}
	return scans, nil
}

func (r mutationResolver) AddCellScans(ctx context.Context, data []*models.SurveyCellScanData, locationID int) ([]*ent.SurveyCellScan, error) {
	return r.CreateCellScans(ctx, data, ScanParentIDs{locationID: &locationID})
}

func (r mutationResolver) CreateCellScans(ctx context.Context, inputs []*models.SurveyCellScanData, parentIDs ScanParentIDs) ([]*ent.SurveyCellScan, error) {
	if err := validateScanParentIDs(parentIDs); err != nil {
		return nil, err
	}
	var (
		client = r.ClientFrom(ctx).SurveyCellScan
		scans  = make([]*ent.SurveyCellScan, len(inputs))
		err    error
	)
	for i, input := range inputs {
		var timestamp *time.Time
		if input.Timestamp != nil {
			inputTime := time.Unix(int64(*input.Timestamp), 0)
			timestamp = &inputTime
		}
		if scans[i], err = client.Create().
			SetNetworkType(input.NetworkType.String()).
			SetSignalStrength(input.SignalStrength).
			SetNillableTimestamp(timestamp).
			SetNillableBaseStationID(input.BaseStationID).
			SetNillableNetworkID(input.NetworkID).
			SetNillableSystemID(input.SystemID).
			SetNillableCellID(input.CellID).
			SetNillableLocationAreaCode(input.LocationAreaCode).
			SetNillableMobileCountryCode(input.MobileCountryCode).
			SetNillableMobileNetworkCode(input.MobileNetworkCode).
			SetNillablePrimaryScramblingCode(input.PrimaryScramblingCode).
			SetNillableOperator(input.Operator).
			SetNillableArfcn(input.Arfcn).
			SetNillablePhysicalCellID(input.PhysicalCellID).
			SetNillableTrackingAreaCode(input.TrackingAreaCode).
			SetNillableTimingAdvance(input.TimingAdvance).
			SetNillableEarfcn(input.Earfcn).
			SetNillableUarfcn(input.Uarfcn).
			SetNillableLatitude(input.Latitude).
			SetNillableLongitude(input.Longitude).
			SetNillableSurveyQuestionID(parentIDs.qid).
			SetNillableLocationID(parentIDs.locationID).
			SetNillableChecklistItemID(parentIDs.checklistItemID).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating survey cell scan: %w", err)
		}
	}
	return scans, nil
}

func (r mutationResolver) CreateSurvey(ctx context.Context, data models.SurveyCreateData) (int, error) {

	client := r.ClientFrom(ctx)
	query := client.Survey.
		Create().
		SetLocationID(data.LocationID).
		SetCompletionTimestamp(time.Unix(int64(data.CompletionTimestamp), 0)).
		SetName(data.Name).
		SetOwnerName(r.Me(ctx).User)
	if data.CreationTimestamp != nil {
		query.SetCreationTimestamp(time.Unix(int64(*data.CreationTimestamp), 0))
	}
	srv, err := query.Save(ctx)
	if err != nil {
		return badID, fmt.Errorf("creating survey: %w", err)
	}

	for _, sr := range data.SurveyResponses {
		query := r.ClientFrom(ctx).SurveyQuestion.
			Create().
			SetFormIndex(sr.FormIndex).
			SetNillableFormName(sr.FormName).
			SetNillableFormDescription(sr.FormDescription).
			SetQuestionIndex(sr.QuestionIndex).
			SetQuestionText(sr.QuestionText).
			SetNillableBoolData(sr.BoolData).
			SetNillableEmailData(sr.EmailData).
			SetNillableLatitude(sr.Latitude).
			SetNillableLongitude(sr.Longitude).
			SetNillableLocationAccuracy(sr.LocationAccuracy).
			SetNillablePhoneData(sr.PhoneData).
			SetNillableTextData(sr.TextData).
			SetNillableFloatData(sr.FloatData).
			SetNillableIntData(sr.IntData).
			SetSurvey(srv)
		if sr.QuestionFormat != nil {
			query.SetQuestionFormat(sr.QuestionFormat.String())
		}
		if sr.DateData != nil {
			query.SetDateData(time.Unix(int64(*sr.DateData), 0))
		}

		if sr.PhotoData != nil {
			f, err :=
				r.createImage(
					ctx,
					&models.AddImageInput{
						ImgKey:   sr.PhotoData.StoreKey,
						FileName: sr.PhotoData.FileName,
						FileSize: func() int {
							if sr.PhotoData.SizeInBytes != nil {
								return *sr.PhotoData.SizeInBytes
							}
							return 0
						}(),
						Modified: time.Now(),
						ContentType: func() string {
							if sr.PhotoData.MimeType != nil {
								return *sr.PhotoData.MimeType
							}
							return "image/jpeg"
						}(),
					},
				)
			if err != nil {
				return badID, err
			}
			query.AddPhotoData(f)
		}

		if sr.ImagesData != nil {
			err = r.addSurveyQuestionImagesData(ctx, sr, query)
			if err != nil {
				return badID, err
			}
		}

		question, err := query.Save(ctx)
		if err != nil {
			return badID, fmt.Errorf("creating survey question: %w", err)
		}

		if sr.QuestionFormat != nil {
			switch *sr.QuestionFormat {
			case models.SurveyQuestionTypeWifi:
				_, err = r.CreateWiFiScans(ctx, sr.WifiData, ScanParentIDs{qid: &question.ID})
			case models.SurveyQuestionTypeCellular:
				_, err = r.CreateCellScans(ctx, sr.CellData, ScanParentIDs{qid: &question.ID})
			}
		}
		if err != nil {
			return badID, err
		}
	}
	return srv.ID, nil
}

func (r mutationResolver) addSurveyQuestionImagesData(ctx context.Context, sr *models.SurveyQuestionResponse, query *ent.SurveyQuestionCreate) error {
	for _, imageData := range sr.ImagesData {
		image, err :=
			r.createImage(
				ctx,
				&models.AddImageInput{
					ImgKey:   imageData.StoreKey,
					FileName: imageData.FileName,
					FileSize: func() int {
						if imageData.SizeInBytes != nil {
							return *imageData.SizeInBytes
						}
						return 0
					}(),
					Modified: time.Now(),
					ContentType: func() string {
						if imageData.MimeType != nil {
							return *imageData.MimeType
						}
						return "image/jpeg"
					}(),
				},
			)
		if err != nil {
			return fmt.Errorf("creating and saving images while creating survey question: %w", err)
		}
		query.AddImages(image)
	}

	return nil
}

func (r mutationResolver) validateRootLocationUniqueness(ctx context.Context, typeid int, name string) error {
	switch exist, err := r.ClientFrom(ctx).
		Location.Query().
		Where(location.Name(name), location.Not(location.HasParent())).
		QueryType().
		Where(locationtype.ID(typeid)).
		Exist(ctx); {
	case err != nil:
		return fmt.Errorf("querying location name existence: %w", err)
	case exist:
		return gqlerror.Errorf("A root location with the name %s already exist", name)
	}
	return nil
}

func (r mutationResolver) verifyLocationParent(ctx context.Context, typeID, parentID int) error {
	typ, err := r.ClientFrom(ctx).
		LocationType.Query().
		Where(locationtype.ID(typeID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("querying location type by id %q: %w", typeID, err)
	}
	ptype, err := r.ClientFrom(ctx).
		Location.Query().
		Where(location.ID(parentID)).
		QueryType().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("querying parent location type by parent id %q: %w", parentID, err)
	}
	if ptype.Index > typ.Index {
		return gqlerror.Errorf("Can't link child to parent with bigger index (%d, %d)", ptype.Index, typ.Index)
	}
	return nil
}

func (r mutationResolver) AddLocation(
	ctx context.Context, input models.AddLocationInput,
) (*ent.Location, error) {
	if input.Parent == nil {
		// ent index enforces uniqueness on (name, type, parent) tuple however
		// no enforcement occurs when parent is not set as NULL is not indexed
		if err := r.validateRootLocationUniqueness(ctx, input.Type, input.Name); err != nil {
			return nil, err
		}
	} else {
		if err := r.verifyLocationParent(ctx, input.Type, *input.Parent); err != nil {
			return nil, err
		}
	}
	var ei *string
	if input.ExternalID != nil && *input.ExternalID != "" {
		ei = input.ExternalID
	}
	l, err := r.ClientFrom(ctx).
		Location.Create().
		SetName(input.Name).
		SetNillableLatitude(input.Latitude).
		SetNillableLongitude(input.Longitude).
		SetTypeID(input.Type).
		SetNillableParentID(input.Parent).
		SetNillableExternalID(ei).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating location: %w", err)
	}
	funcSetLocation := func(b *ent.PropertyCreate) { b.SetLocation(l) }
	if _, err := r.AddProperties(input.Properties, resolverutil.AddPropertyArgs{Context: ctx, EntSetter: funcSetLocation}); err != nil {
		return nil, fmt.Errorf("creating location properties: %w", err)
	}
	return l, nil
}

func (r mutationResolver) AddLocationType(
	ctx context.Context, input models.AddLocationTypeInput,
) (*ent.LocationType, error) {
	props, err := r.AddPropertyTypes(ctx, input.Properties...)
	if err != nil {
		return nil, err
	}
	categories, err := r.AddSurveyTemplateCategories(ctx, input.SurveyTemplateCategories...)
	if err != nil {
		return nil, err
	}
	index, err := r.ClientFrom(ctx).LocationType.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	typ, err := r.ClientFrom(ctx).LocationType.
		Create().
		SetName(input.Name).
		SetNillableMapType(input.MapType).
		SetNillableMapZoomLevel(input.MapZoomLevel).
		SetNillableSite(input.IsSite).
		SetIndex(index).
		AddPropertyTypes(props...).
		AddSurveyTemplateCategories(categories...).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A location type with the name %s already exists", input.Name)
		}
		return nil, fmt.Errorf("creating location type: %w", err)
	}
	return typ, nil
}

func (r mutationResolver) AddEquipmentPorts(ctx context.Context, et *ent.EquipmentType, e *ent.Equipment) ([]*ent.EquipmentPort, error) {
	ids, err := et.QueryPortDefinitions().IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying port definitions %q: %w", et.ID, err)
	}
	var (
		client = r.ClientFrom(ctx).EquipmentPort
		ports  = make([]*ent.EquipmentPort, len(ids))
	)
	for i, id := range ids {
		if ports[i], err = client.Create().
			SetDefinitionID(id).
			SetParent(e).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating equipment port: %w", err)
		}
	}
	return ports, nil
}

func (r mutationResolver) AddEquipmentPositions(ctx context.Context, et *ent.EquipmentType, e *ent.Equipment) ([]*ent.EquipmentPosition, error) {
	ids, err := et.QueryPositionDefinitions().IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying position definitions %q: %w", et.ID, err)
	}
	var (
		client    = r.ClientFrom(ctx).EquipmentPosition
		positions = make([]*ent.EquipmentPosition, len(ids))
	)
	for i, id := range ids {
		if positions[i], err = client.Create().
			SetDefinitionID(id).
			SetParent(e).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating equipment position: %w", err)
		}
	}
	return positions, nil
}

func (r mutationResolver) getOrCreatePort(ctx context.Context, side *models.LinkSide) (*ent.EquipmentPort, error) {
	client := r.ClientFrom(ctx)
	port, err := client.Equipment.Query().
		Where(equipment.ID(side.Equipment)).
		QueryPorts().
		Where(equipmentport.HasDefinitionWith(
			equipmentportdefinition.ID(side.Port),
		)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying port %d: %w", side.Port, err)
	}
	if port != nil {
		return port, nil
	}
	if port, err = client.EquipmentPort.Create().
		SetDefinitionID(side.Port).
		SetParentID(side.Equipment).
		Save(ctx); err != nil {
		return nil, fmt.Errorf("creating equipment port: %w", err)
	}
	return port, nil
}

func (r mutationResolver) addEquipment(
	ctx context.Context, typ *ent.EquipmentType,
	input models.AddEquipmentInput,
) (*ent.Equipment, error) {
	ep, err := resolverutil.GetOrCreatePosition(ctx, r.ClientFrom(ctx), input.Parent, input.PositionDefinition, true)
	if err != nil {
		return nil, err
	}
	var positionID *int
	if ep != nil {
		switch exist, err := ep.QueryParent().QueryPositions().
			Where(equipmentposition.HasAttachmentWith(
				equipment.Name(input.Name),
				equipment.HasTypeWith(equipmenttype.ID(typ.ID)),
			)).
			Exist(ctx); {
		case err != nil:
			return nil, fmt.Errorf("querying parent position: %w", err)
		case exist:
			return nil, errors.New("equipment already exist under parent")
		}
		positionID = &ep.ID
	}
	if err := r.validateEquipmentNameIsUnique(
		ctx, input.Name, input.Location,
		positionID, nil,
	); err != nil {
		return nil, err
	}

	var ei *string
	if input.ExternalID != nil && *input.ExternalID != "" {
		ei = input.ExternalID
	}
	e, err := r.ClientFrom(ctx).
		Equipment.Create().
		SetName(input.Name).
		SetType(typ).
		SetNillableExternalID(ei).
		SetNillableParentPositionID(positionID).
		SetNillableLocationID(input.Location).
		SetNillableWorkOrderID(input.WorkOrder).
		SetNillableFutureState(func() *string {
			if input.WorkOrder != nil {
				state := models.FutureStateInstall.String()
				return &state
			}
			return nil
		}()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating equipment: %w", err)
	}
	addPropertyArgs := resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetEquipment(e) },
	}
	if _, err := r.AddProperties(input.Properties, addPropertyArgs); err != nil {
		return nil, fmt.Errorf("creating equipment properties: %w", err)
	}
	if _, err := r.AddEquipmentPorts(ctx, typ, e); err != nil {
		return nil, fmt.Errorf("creating equipment ports: %w", err)
	}
	if _, err := r.AddEquipmentPositions(ctx, typ, e); err != nil {
		return nil, fmt.Errorf("creating equipment positions: %w", err)
	}
	return e, nil
}

func (r mutationResolver) AddEquipment(
	ctx context.Context, input models.AddEquipmentInput,
) (*ent.Equipment, error) {
	if input.Location == nil && (input.Parent == nil || input.PositionDefinition == nil) {
		return nil, errors.New("location or position data is required")
	}
	typ, err := r.ClientFrom(ctx).EquipmentType.Get(ctx, input.Type)
	if err != nil {
		return nil, fmt.Errorf("querying equipment type %q: %w", input.Type, err)
	}
	return r.addEquipment(ctx, typ, input)
}

func (r mutationResolver) AddEquipmentPositionDefinitions(
	ctx context.Context, inputs []*models.EquipmentPositionInput, equipmentTypeID *int,
) ([]*ent.EquipmentPositionDefinition, error) {
	if equipmentTypeID != nil {
		query := r.ClientFrom(ctx).
			EquipmentType.Query().
			Where(equipmenttype.ID(*equipmentTypeID)).
			QueryPositionDefinitions()
		for _, input := range inputs {
			def, err := query.Clone().
				Where(equipmentpositiondefinition.Name(input.Name)).
				First(ctx)
			switch {
			case err != nil && !ent.IsNotFound(err):
				return nil, fmt.Errorf("querying position definition name existence: %w", err)
			case def != nil:
				r.logger.For(ctx).Error("duplicate position definition name for equipment type",
					zap.String("name", input.Name),
					zap.Intp("type", equipmentTypeID),
				)
				return nil, gqlerror.Errorf(
					"A position definition with the name %v already exists under %v",
					input.Name, equipmentTypeID,
				)
			}
		}
	}
	var (
		client = r.ClientFrom(ctx).EquipmentPositionDefinition
		defs   = make([]*ent.EquipmentPositionDefinition, len(inputs))
		err    error
	)
	for i, input := range inputs {
		if defs[i], err = client.Create().
			SetName(input.Name).
			SetNillableIndex(input.Index).
			SetNillableVisibilityLabel(input.VisibleLabel).
			SetNillableEquipmentTypeID(equipmentTypeID).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating equipment position definition: %w", err)
		}
	}
	return defs, nil
}

func (r mutationResolver) AddEquipmentPortDefinitions(
	ctx context.Context, inputs []*models.EquipmentPortInput, equipmentTypeID *int,
) ([]*ent.EquipmentPortDefinition, error) {
	if equipmentTypeID != nil {
		query := r.ClientFrom(ctx).
			EquipmentType.Query().
			Where(equipmenttype.ID(*equipmentTypeID)).
			QueryPortDefinitions()
		for _, input := range inputs {
			pd, err := query.Clone().
				Where(equipmentportdefinition.Name(input.Name)).
				First(ctx)
			switch {
			case err != nil && !ent.IsNotFound(err):
				return nil, fmt.Errorf("querying port definition name existence: %w", err)
			case pd != nil:
				r.logger.For(ctx).Error("duplicate port definition name for equipment type ",
					zap.String("name", input.Name),
					zap.Intp("type", equipmentTypeID),
				)
				return nil, gqlerror.Errorf(
					"A port definition with the name %v already exists under %v",
					input.Name, *equipmentTypeID,
				)
			}
		}
	}
	var (
		client = r.ClientFrom(ctx).EquipmentPortDefinition
		defs   = make([]*ent.EquipmentPortDefinition, len(inputs))
		err    error
	)
	for i, input := range inputs {
		if defs[i], err = client.Create().
			SetName(input.Name).
			SetNillableIndex(input.Index).
			SetNillableBandwidth(input.Bandwidth).
			SetNillableVisibilityLabel(input.VisibleLabel).
			SetNillableEquipmentPortTypeID(input.PortTypeID).
			SetNillableEquipmentTypeID(equipmentTypeID).
			Save(ctx); err != nil {
			return nil, fmt.Errorf("creating equipment port definition: %w", err)
		}
	}
	return defs, nil
}

func (r mutationResolver) AddEquipmentPortType(
	ctx context.Context, input models.AddEquipmentPortTypeInput,
) (*ent.EquipmentPortType, error) {
	props, err := r.AddPropertyTypes(ctx, input.Properties...)
	if err != nil {
		return nil, err
	}
	linkProps, err := r.AddPropertyTypes(ctx, input.LinkProperties...)
	if err != nil {
		return nil, err
	}
	et, err := r.ClientFrom(ctx).
		EquipmentPortType.
		Create().
		SetName(input.Name).
		AddPropertyTypes(props...).
		AddLinkPropertyTypes(linkProps...).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("An equipment port type with the name %s already exists", input.Name)
		}
		return nil, fmt.Errorf("creating equipment type: %w", err)
	}
	return et, nil
}

func (r mutationResolver) AddEquipmentType(
	ctx context.Context, input models.AddEquipmentTypeInput,
) (*ent.EquipmentType, error) {
	positions, err := r.AddEquipmentPositionDefinitions(ctx, input.Positions, nil)
	if err != nil {
		return nil, err
	}
	ports, err := r.AddEquipmentPortDefinitions(ctx, input.Ports, nil)
	if err != nil {
		return nil, err
	}
	props, err := r.AddPropertyTypes(ctx, input.Properties...)
	if err != nil {
		return nil, err
	}
	client := r.ClientFrom(ctx)
	typ, err := client.
		EquipmentType.Create().
		SetName(input.Name).
		AddPositionDefinitions(positions...).
		AddPortDefinitions(ports...).
		AddPropertyTypes(props...).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("An equipment type with the name %v already exists", input.Name)
		}
		return nil, fmt.Errorf("creating equipment type: %w", err)
	}
	if input.Category != nil {
		if typ, err = r.updateEquipmentTypeCategory(
			ctx, client, typ, *input.Category,
		); err != nil {
			return nil, fmt.Errorf("updating equipment category: %w", err)
		}
	}
	return typ, nil
}

func (r mutationResolver) EditLocation(
	ctx context.Context, input models.EditLocationInput,
) (*ent.Location, error) {
	client := r.ClientFrom(ctx)
	l, err := client.Location.Get(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("querying location %q: %w", input.ID, err)
	}
	lt, err := l.QueryType().OnlyID(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying location type of %q: %w", l.ID, err)
	}
	if l.Name != input.Name {
		switch exist, err := l.QueryParent().Exist(ctx); {
		case err != nil:
			return nil, fmt.Errorf("querying location parent existence: %w", err)
		case !exist:
			// root location requires validation, see comment in AddLocation
			if err := r.validateRootLocationUniqueness(ctx, lt, input.Name); err != nil {
				return nil, err
			}
		}
	}

	upd := client.Location.
		UpdateOne(l).
		SetName(input.Name).
		SetLatitude(input.Latitude).
		SetLongitude(input.Longitude)
	if input.ExternalID != nil && *input.ExternalID != "" {
		upd.SetNillableExternalID(input.ExternalID)
	} else {
		upd.ClearExternalID()
	}
	if l, err = upd.Save(ctx); err != nil {
		return nil, fmt.Errorf("updating location %q: %w", input.ID, err)
	}
	var added, edited []*models.PropertyInput
	directPropertiesTypes, err := l.QueryProperties().QueryType().IDs(ctx)
	if err != nil {
		return nil, err
	}
	for _, input := range input.Properties {
		if r.isNewProp(directPropertiesTypes, input.ID, input.PropertyTypeID) {
			added = append(added, input)
		} else {
			if input.ID == nil {
				propID, err := l.QueryProperties().Where(property.HasTypeWith(propertytype.ID(input.PropertyTypeID))).OnlyID(ctx)
				if err != nil {
					return nil, err
				}
				input.ID = &propID
			}
			edited = append(edited, input)
		}
	}
	if _, err := r.AddProperties(added, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetLocation(l) },
	}); err != nil {
		return nil, err
	}
	for _, input := range edited {
		typ, err := client.LocationType.Query().
			Where(locationtype.ID(lt)).
			QueryPropertyTypes().
			Where(propertytype.ID(input.PropertyTypeID)).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("querying location property type %q: %w", input.PropertyTypeID, err)
		}
		if typ.Editable && typ.IsInstanceProperty {
			query := client.Property.
				Update().
				Where(
					property.HasLocationWith(location.ID(l.ID)),
					property.ID(*input.ID),
				)
			if r.updatePropValues(ctx, input, query) != nil {
				return nil, fmt.Errorf("updating property values: %w", err)
			}
		}
	}
	return l, nil
}

func (r mutationResolver) RemoveEquipmentFromPosition(ctx context.Context, positionID int, workOrderID *int) (*ent.EquipmentPosition, error) {
	client := r.ClientFrom(ctx)
	ep, err := client.EquipmentPosition.Get(ctx, positionID)
	if err != nil {
		return nil, fmt.Errorf("querying equipment position: %w", err)
	}

	e, err := ep.QueryAttachment().First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("querying position attachment: %w", err)
	}
	if e == nil {
		return ep, nil
	}
	if workOrderID != nil {
		exist, err := client.WorkOrder.Query().
			Where(workorder.ID(*workOrderID)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("querying work order %q from equipment %q: %w", e.ID, *workOrderID, err)
		}
		if exist {
			switch exist, err := e.QueryWorkOrder().Where(workorder.ID(*workOrderID)).Exist(ctx); {
			case err != nil:
				return nil, fmt.Errorf("querying work order %q: %w", e.ID, err)
			case exist:
				return ep, r.removeEquipment(ctx, e)
			}
			if err := client.Equipment.
				UpdateOne(e).
				ClearWorkOrder().
				SetWorkOrderID(*workOrderID).
				SetFutureState(models.FutureStateRemove.String()).
				Exec(ctx); err != nil {
				return nil, fmt.Errorf("updating attached equipment %q: %w", e.ID, err)
			}
			return ep, nil
		}
	} else if err := r.removeEquipment(ctx, e); err != nil {
		return nil, err
	}
	return ep, nil
}

func (r mutationResolver) hasPositionCycle(ctx context.Context, parent, child int) bool {
	current := r.ClientFrom(ctx).Equipment.GetX(ctx, parent)
	seen := map[int]struct{}{child: {}}
	for current != nil {
		if _, ok := seen[current.ID]; ok {
			r.logger.For(ctx).Warn("equipment position cycle",
				zap.Int("current", current.ID),
				zap.Reflect("seen", seen),
			)
			return true
		}
		seen[current.ID] = struct{}{}
		current = current.QueryParentPosition().QueryParent().FirstX(ctx)
	}
	return false
}

func (r mutationResolver) MoveEquipmentToPosition(
	ctx context.Context, parentEquipmentID, positionDefinitionID *int, equipmentID int,
) (*ent.EquipmentPosition, error) {
	ep, err := resolverutil.GetOrCreatePosition(
		ctx, r.ClientFrom(ctx), parentEquipmentID, positionDefinitionID, true,
	)
	if err != nil {
		return nil, err
	}
	var (
		client = r.ClientFrom(ctx)
		e      *ent.Equipment
	)
	if e, err = client.Equipment.Get(ctx, equipmentID); err != nil {
		return nil, fmt.Errorf("querying equipment %d: %w", equipmentID, err)
	}
	if parentEquipmentID != nil && r.hasPositionCycle(ctx, *parentEquipmentID, equipmentID) {
		return nil, fmt.Errorf("equipment position %d cycle, parent %d", equipmentID, e.ID)
	}
	if err := client.Equipment.
		UpdateOne(e).
		SetParentPosition(ep).
		ClearLocation().
		Exec(ctx); err != nil {
		return nil, fmt.Errorf("moving equipment %d to position %d: %w", equipmentID, ep.ID, err)
	}
	return ep, nil
}

// NOTE: Be aware that this method is used to create both images and files. Will be renamed in another Diff.
func (r mutationResolver) createImage(ctx context.Context, input *models.AddImageInput) (*ent.File, error) {
	img, err := r.ClientFrom(ctx).
		File.Create().
		SetStoreKey(input.ImgKey).
		SetName(input.FileName).
		SetSize(input.FileSize).
		SetModifiedAt(input.Modified).
		SetUploadedAt(time.Now()).
		SetType(func() string {
			if strings.HasPrefix(input.ContentType, "image/") {
				return models.FileTypeImage.String()
			}
			return models.FileTypeFile.String()
		}()).
		SetContentType(input.ContentType).
		SetNillableCategory(input.Category).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating image for key %q: %w", input.ImgKey, err)
	}
	return img, nil
}

type execer interface{ Exec(context.Context) error }

func (r mutationResolver) AddImage(ctx context.Context, input models.AddImageInput) (*ent.File, error) {
	image, err := r.createImage(ctx, &input)
	if err != nil {
		return nil, err
	}
	var (
		client = r.ClientFrom(ctx)
		execer execer
	)
	switch input.EntityType {
	case models.ImageEntityLocation:
		execer = client.Location.
			UpdateOneID(input.EntityID).
			AddFiles(image)
	case models.ImageEntitySiteSurvey:
		execer = client.Survey.
			UpdateOneID(input.EntityID).
			SetSourceFile(image)
	case models.ImageEntityWorkOrder:
		execer = client.WorkOrder.
			UpdateOneID(input.EntityID).
			AddFiles(image)
	case models.ImageEntityEquipment:
		execer = client.Equipment.
			UpdateOneID(input.EntityID).
			AddFiles(image)
	case models.ImageEntityUser:
		execer = client.User.
			UpdateOneID(input.EntityID).
			SetProfilePhoto(image)
	default:
		return nil, fmt.Errorf("unknown image owner type: %s", input.EntityType)
	}
	if err := execer.Exec(ctx); err != nil {
		return nil, fmt.Errorf("adding image to type %s: %w", input.EntityType, err)
	}
	return image, nil
}

func (r mutationResolver) AddHyperlink(ctx context.Context, input models.AddHyperlinkInput) (*ent.Hyperlink, error) {
	client := r.ClientFrom(ctx)
	hyperlink, err := client.Hyperlink.
		Create().
		SetURL(input.URL).
		SetNillableName(input.DisplayName).
		SetNillableCategory(input.Category).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating hyperlink for url %q: %w", input.URL, err)
	}
	var execer execer
	switch input.EntityType {
	case models.ImageEntityLocation:
		execer = client.Location.
			UpdateOneID(input.EntityID).
			AddHyperlinks(hyperlink)
	case models.ImageEntityWorkOrder:
		execer = client.WorkOrder.
			UpdateOneID(input.EntityID).
			AddHyperlinks(hyperlink)
	case models.ImageEntityEquipment:
		execer = client.Equipment.
			UpdateOneID(input.EntityID).
			AddHyperlinks(hyperlink)
	default:
		return nil, fmt.Errorf("unknown hyperlink owner type: %s", input.EntityType)
	}
	if err := execer.Exec(ctx); err != nil {
		return nil, fmt.Errorf("adding hyperlink to type %s: %w", input.EntityType, err)
	}
	return hyperlink, nil
}

func (r mutationResolver) DeleteHyperlink(ctx context.Context, id int) (*ent.Hyperlink, error) {
	client := r.ClientFrom(ctx).Hyperlink
	h, err := client.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("querying hyperlink %d: %w", id, err)
	}
	if err := client.DeleteOne(h).Exec(ctx); err != nil {
		return nil, fmt.Errorf("deleting hyperlink %d: %w", id, err)
	}
	return h, nil
}

func (r mutationResolver) DeleteImage(ctx context.Context, _ models.ImageEntity, _, id int) (*ent.File, error) {
	client := r.ClientFrom(ctx).File
	f, err := client.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("querying file id %d: %w", id, err)
	}
	if err := client.DeleteOne(f).Exec(ctx); err != nil {
		return nil, fmt.Errorf("deleting file id %d: %w", id, err)
	}
	return f, nil
}

func (r mutationResolver) AddComment(ctx context.Context, input models.CommentInput) (*ent.Comment, error) {
	client := r.ClientFrom(ctx)
	u, err := viewer.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}
	c, err := client.Comment.Create().
		SetAuthor(u).
		SetText(input.Text).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}
	var execer execer
	switch input.EntityType {
	case models.CommentEntityWorkOrder:
		execer = client.WorkOrder.
			UpdateOneID(input.ID).
			AddComments(c)
	case models.CommentEntityProject:
		execer = client.Project.
			UpdateOneID(input.ID).
			AddComments(c)
	default:
		return nil, fmt.Errorf("unknown comment owner type: %s", input.EntityType)
	}
	if err := execer.Exec(ctx); err != nil {
		return nil, fmt.Errorf("adding comment to type %s: %w", input.EntityType, err)
	}
	return c, nil
}

func (r mutationResolver) AddLink(
	ctx context.Context, input models.AddLinkInput,
) (*ent.Link, error) {
	ids := make([]int, len(input.Sides))
	for i, side := range input.Sides {
		port, err := r.getOrCreatePort(ctx, side)
		if err != nil {
			return nil, err
		}
		switch linked, err := port.QueryLink().Exist(ctx); {
		case err != nil:
			return nil, fmt.Errorf("querying link existence: %w", err)
		case linked:
			return nil, fmt.Errorf("port %d already has link: %w", port.ID, err)
		}
		ids[i] = port.ID
	}
	if count, err := r.ClientFrom(ctx).EquipmentPort.Query().
		Where(
			equipmentport.IDIn(ids...),
			equipmentport.Not(equipmentport.HasLink()),
		).
		Count(ctx); err != nil || count != 2 {
		return nil, fmt.Errorf("querying ports %v: %w", ids, err)
	}
	l, err := r.ClientFrom(ctx).Link.Create().
		AddPortIDs(ids...).
		SetNillableWorkOrderID(input.WorkOrder).
		SetNillableFutureState(func() *string {
			if input.WorkOrder != nil {
				state := models.FutureStateInstall.String()
				return &state
			}
			return nil
		}()).
		AddServiceIDs(input.ServiceIds...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating link between ports %v: %w", ids, err)
	}
	if _, err := r.AddProperties(input.Properties, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetLink(l) },
	}); err != nil {
		return nil, fmt.Errorf("creating link properties: %w", err)
	}
	return l, err
}

func (r mutationResolver) EditLink(
	ctx context.Context, input models.EditLinkInput,
) (*ent.Link, error) {
	client := r.ClientFrom(ctx)
	l, err := client.Link.Get(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("querying link %d: %w", input.ID, err)
	}

	var added, edited []*models.PropertyInput
	directPropertiesTypes, err := l.QueryProperties().QueryType().IDs(ctx)
	if err != nil {
		return nil, err
	}
	for _, input := range input.Properties {
		if r.isNewProp(directPropertiesTypes, input.ID, input.PropertyTypeID) {
			added = append(added, input)
		} else {
			if input.ID == nil {
				propID, err := l.QueryProperties().Where(property.HasTypeWith(propertytype.ID(input.PropertyTypeID))).OnlyID(ctx)
				if err != nil {
					return nil, err
				}
				input.ID = &propID
			}
			edited = append(edited, input)
		}
	}
	if _, err := r.AddProperties(added, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetLinkID(l.ID) },
	}); err != nil {
		return nil, err
	}
	for _, input := range edited {
		typ, err := l.QueryPorts().
			QueryDefinition().
			QueryEquipmentPortType().
			QueryLinkPropertyTypes().
			Where(propertytype.ID(input.PropertyTypeID)).
			First(ctx)
		if err != nil {
			return nil, fmt.Errorf("querying link property type %d: %w", input.PropertyTypeID, err)
		}
		if typ.Editable && typ.IsInstanceProperty {
			query := client.Property.
				Update().
				Where(
					property.HasLinkWith(link.ID(l.ID)),
					property.ID(*input.ID),
				)
			if r.updatePropValues(ctx, input, query) != nil {
				return nil, fmt.Errorf("updating property values: %w", err)
			}
		}
	}

	currentServiceIds, err := l.QueryService().IDs(ctx)
	if err != nil {
		return nil, err
	}
	addedServiceIds, deletedServiceIds := resolverutil.GetDifferenceBetweenSlices(currentServiceIds, input.ServiceIds)
	for _, serviceID := range addedServiceIds {
		if _, err := r.AddServiceLink(ctx, serviceID, l.ID); err != nil {
			return nil, err
		}
	}
	for _, serviceID := range deletedServiceIds {
		if _, err := r.RemoveServiceLink(ctx, serviceID, l.ID); err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (r mutationResolver) removeLink(ctx context.Context, link *ent.Link) error {
	if err := r.ClientFrom(ctx).Link.
		DeleteOne(link).
		Exec(ctx); err != nil {
		return fmt.Errorf("removing link %d: %w", link.ID, err)
	}
	return nil
}

func (r mutationResolver) RemoveLink(ctx context.Context, id int, workOrderID *int) (*ent.Link, error) {
	client := r.ClientFrom(ctx)
	l, err := client.Link.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting link %d: %w", id, err)
	}
	if workOrderID != nil {
		switch exist, err := client.WorkOrder.
			Query().
			Where(workorder.ID(*workOrderID)).
			Exist(ctx); {
		case err != nil:
			return nil, errors.Wrapf(err, "querying work order from link: l=%q, wo=%q", l.ID, *workOrderID)
		case exist:
			if err := client.Link.
				UpdateOne(l).
				ClearWorkOrder().
				SetWorkOrderID(*workOrderID).
				SetFutureState(models.FutureStateRemove.String()).
				Exec(ctx); err != nil {
				return nil, err
			}
			return l, nil
		}
	} else if err := r.removeLink(ctx, l); err != nil {
		return nil, errors.Wrapf(err, "removing link: id=%q", id)
	}
	return l, nil
}

func (r mutationResolver) removeSurveyQuestion(ctx context.Context, question *ent.SurveyQuestion) error {
	client := r.ClientFrom(ctx)
	if _, err := client.SurveyCellScan.Delete().
		Where(surveycellscan.HasSurveyQuestionWith(surveyquestion.ID(question.ID))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "deleting survey cell scan: id=%q", question.ID)
	}
	if _, err := r.ClientFrom(ctx).SurveyWiFiScan.Delete().
		Where(surveywifiscan.HasSurveyQuestionWith(surveyquestion.ID(question.ID))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "deleting survey wifi scan: id=%q", question.ID)
	}
	ids, err := question.QueryPhotoData().IDs(ctx)
	if err != nil {
		return errors.Wrapf(err, "querying question photos ids: id=%q", question.ID)
	}
	// TODO(T47446957): Delete S3 files of sitesurvey on sitesurvey graphql deletion
	if _, err := client.File.Delete().Where(file.IDIn(ids...)).Exec(ctx); err != nil {
		return errors.Wrapf(err, "deleting question photos: id=%q, count=%d", question.ID, len(ids))
	}
	if err := client.SurveyQuestion.DeleteOne(question).Exec(ctx); err != nil {
		return errors.Wrapf(err, "deleting survey question: id=%q", question.ID)
	}
	return nil
}

func (r mutationResolver) RemoveSiteSurvey(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	questions, err := client.SurveyQuestion.Query().
		Where(surveyquestion.HasSurveyWith(survey.ID(id))).
		All(ctx)
	if err != nil {
		return id, errors.Wrapf(err, "querying survey questions")
	}
	for _, question := range questions {
		if err := r.removeSurveyQuestion(ctx, question); err != nil {
			return id, err
		}
	}
	if err := client.Survey.DeleteOneID(id).Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting survey")
	}
	return id, nil
}

func (r mutationResolver) RemoveLocation(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	l, err := client.Location.Query().
		Where(
			location.ID(id),
			location.Not(location.HasChildren()),
			location.Not(location.HasFiles()),
			location.Not(location.HasEquipment()),
			location.Not(location.HasSurvey()),
		).
		Only(ctx)
	if err != nil {
		return id, errors.Wrapf(err, "querying location: id=%q", id)
	}
	if _, err := client.Property.Delete().Where(property.HasLocationWith(location.ID(id))).Exec(ctx); err != nil {
		return id, errors.Wrapf(err, "deleting location properties: id=%q", id)
	}
	if err := client.Location.DeleteOne(l).Exec(ctx); err != nil {
		return id, errors.Wrapf(err, "deleting location: id=%q", id)
	}
	return id, nil
}

func (r mutationResolver) RemoveWorkOrder(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	wo, err := client.WorkOrder.Get(ctx, id)
	if err != nil {
		return id, errors.Wrapf(err, "querying work order: id=%q", id)
	}

	equipments, err := wo.QueryEquipment().All(ctx)
	if err != nil {
		return id, errors.Wrapf(err, "query work order equipment: id=%q", id)
	}
	for _, e := range equipments {
		e := e
		if e.FutureState == models.FutureStateInstall.String() {
			if err := r.removeEquipment(ctx, e); err != nil {
				return id, errors.Wrapf(err, "deleting to be installed equipment in work order e=%q, wo=%q", e.ID, id)
			}
		} else {
			err := client.Equipment.
				UpdateOne(e).
				ClearWorkOrder().
				SetFutureState("").
				Exec(ctx)
			if err != nil {
				return id, errors.Wrapf(err, "deleting future remove state from to be removed equipment in work order e=%q, wo=%q", e.ID, id)
			}
		}
	}

	links, err := wo.QueryLinks().All(ctx)
	if err != nil {
		return id, errors.Wrapf(err, "query work order links: id=%q", id)
	}
	for _, l := range links {
		if l.FutureState == models.FutureStateInstall.String() {
			if _, err := r.RemoveLink(ctx, l.ID, nil); err != nil {
				return id, errors.Wrapf(err, "deleting to be installed link in work order l=%q, wo=%q", l.ID, id)
			}
		} else {
			if err := client.Link.
				UpdateOne(l).
				ClearWorkOrder().
				SetFutureState("").
				Exec(ctx); err != nil {
				return id, errors.Wrapf(err, "deleting future remove state from to be removed link in work order l=%q, wo=%q", l.ID, id)
			}
		}
	}

	if err := client.WorkOrder.DeleteOne(wo).Exec(ctx); err != nil {
		return id, errors.Wrapf(err, "deleting work order wo=%q", id)
	}
	return id, nil
}

func (r mutationResolver) removeEquipment(ctx context.Context, e *ent.Equipment) error {
	client := r.ClientFrom(ctx)
	if _, err := r.ClientFrom(ctx).Property.Delete().
		Where(property.HasEquipmentWith(equipment.ID(e.ID))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "deleting equipment properties e=%q", e.ID)
	}

	ids, err := e.QueryPositions().IDs(ctx)
	if err != nil {
		return errors.Wrapf(err, "querying equipment positions: id=%q", e.ID)
	}
	if len(ids) > 0 {
		for _, id := range ids {
			if _, err := r.RemoveEquipmentFromPosition(ctx, id, nil); err != nil {
				return errors.Wrapf(err, "remove equipment from position e=%q, id=%q", e.ID, id)
			}
		}
		if _, err := client.EquipmentPosition.Delete().
			Where(equipmentposition.IDIn(ids...)).
			Exec(ctx); err != nil {
			return errors.Wrapf(err, "remove equipment positions e=%q", e.ID)
		}
	}

	if _, err := client.Link.Delete().
		Where(link.HasPortsWith(equipmentport.HasParentWith(equipment.ID(e.ID)))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "delete links of equipment e=%q", e.ID)
	}
	if _, err := client.ServiceEndpoint.Delete().
		Where(serviceendpoint.HasPortWith(equipmentport.HasParentWith(equipment.ID(e.ID)))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "delete service endpoints of equipment e=%q", e.ID)
	}
	if _, err := client.EquipmentPort.Delete().
		Where(equipmentport.HasParentWith(equipment.ID(e.ID))).
		Exec(ctx); err != nil {
		return errors.Wrapf(err, "delete ports of equipment e=%q", e.ID)
	}

	if err := client.Equipment.DeleteOne(e).Exec(ctx); err != nil && !ent.IsNotFound(err) {
		return errors.Wrapf(err, "delete equipment e=%q", e.ID)
	}
	return nil
}

func (r mutationResolver) RemoveEquipment(ctx context.Context, id int, workOrderID *int) (int, error) {
	client := r.ClientFrom(ctx)
	e, err := client.Equipment.Get(ctx, id)
	if err != nil {
		return id, errors.Wrapf(err, "query equipment: id=%q", id)
	}
	if workOrderID != nil {
		exist, err := client.WorkOrder.Query().
			Where(workorder.ID(*workOrderID)).
			Exist(ctx)
		if err != nil || !exist {
			return id, errors.Wrapf(err, "querying work order from equipment: e=%q, wo=%q", e.ID, *workOrderID)
		}
		if err := client.Link.Update().
			Where(link.HasPortsWith(equipmentport.HasParentWith(equipment.ID(e.ID)))).
			ClearWorkOrder().
			SetWorkOrderID(*workOrderID).
			SetFutureState(models.FutureStateRemove.String()).
			Exec(ctx); err != nil {
			return id, errors.Wrapf(err, "delete links of equipment e=%q, wo=%q", e.ID, *workOrderID)
		}

		ids, err := e.QueryPositions().IDs(ctx)
		if err != nil {
			return id, errors.Wrapf(err, "querying positions of equipment: e=%q", e.ID)
		}
		for _, id := range ids {
			if _, err := r.RemoveEquipmentFromPosition(ctx, id, workOrderID); err != nil {
				return id, errors.WithMessagef(err, "removing equipment from position: e=%q, id=%q, wo=%q", e.ID, id, *workOrderID)
			}
		}
		if err := client.Equipment.UpdateOne(e).
			ClearWorkOrder().
			SetWorkOrderID(*workOrderID).
			SetFutureState(models.FutureStateRemove.String()).
			Exec(ctx); err != nil {
			return id, errors.Wrapf(err, "attaching equipment to work order: e=%q, wo=%q", id, *workOrderID)
		}
		return id, nil
	}
	return id, r.removeEquipment(ctx, e)
}

func (r mutationResolver) RemoveEquipmentPortType(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	pt, err := client.EquipmentPortType.Get(ctx, id)
	if err != nil {
		return id, errors.Wrapf(err, "equipment port type does not exist: id=%q", id)
	}
	switch exist, err := pt.QueryPortDefinitions().Exist(ctx); {
	case err != nil:
		return id, errors.Wrapf(err, "querying locations for type: id=%q", pt.ID)
	case exist:
		return id, errors.Errorf("cannot delete location type with existing locations")
	}
	if _, err := client.PropertyType.Delete().
		Where(propertytype.HasEquipmentPortTypeWith(equipmentporttype.ID(id))).
		Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting property type")
	}
	if err := client.EquipmentPortType.DeleteOne(pt).Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting equipment port type")
	}
	return id, nil
}

func (r mutationResolver) RemoveEquipmentType(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	t, err := client.EquipmentType.Query().
		Where(
			equipmenttype.ID(id),
			equipmenttype.Not(equipmenttype.HasEquipment()),
		).
		Only(ctx)
	if err != nil {
		return id, errors.Wrapf(err, "querying equipment type: id=%q", id)
	}
	if _, err := client.EquipmentPortDefinition.Delete().
		Where(equipmentportdefinition.HasEquipmentTypeWith(equipmenttype.ID(id))).
		Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting equipment port definition")
	}
	if _, err := client.EquipmentPositionDefinition.Delete().
		Where(equipmentpositiondefinition.HasEquipmentTypeWith(equipmenttype.ID(id))).
		Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting equipment position definition")
	}
	if _, err := client.PropertyType.Delete().
		Where(propertytype.HasEquipmentTypeWith(equipmenttype.ID(id))).
		Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting property type")
	}
	if err := client.EquipmentType.DeleteOne(t).Exec(ctx); err != nil {
		return id, errors.Wrap(err, "deleting equipment type")
	}
	return id, nil
}

func (r mutationResolver) ExecuteWorkOrder(ctx context.Context, id int) (*models.WorkOrderExecutionResult, error) {
	wo, err := r.ClientFrom(ctx).WorkOrder.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot find work order with id=%q", id)
	}

	var (
		equipments []*ent.Equipment
		links      []*ent.Link
	)
	if equipments, err = wo.QueryEquipment().All(ctx); err != nil {
		return nil, errors.Wrapf(err, "query work order equipments wo=%q", id)
	}
	if links, err = wo.QueryLinks().All(ctx); err != nil {
		return nil, errors.Wrapf(err, "query work order links wo=%q", id)
	}

	result := models.WorkOrderExecutionResult{ID: wo.ID, Name: wo.Name}
	for _, l := range links {
		if l.FutureState == models.FutureStateRemove.String() {
			if err := r.removeLink(ctx, l); err != nil {
				return nil, errors.Wrapf(err, "remove work order link l=%q, wo=%q", l.ID, id)
			}
			result.LinkRemoved = append(result.LinkRemoved, l.ID)
		}
	}

	for _, e := range equipments {
		if e.FutureState == models.FutureStateRemove.String() {
			if err := r.removeEquipment(ctx, e); err != nil {
				return nil, errors.Wrapf(err, "remove work order equipment e=%q, wo=%q", e.ID, id)
			}
			result.EquipmentRemoved = append(result.EquipmentRemoved, e.ID)
		}
	}

	for _, e := range equipments {
		if e.FutureState == models.FutureStateInstall.String() {
			eid := e.ID
			switch exist, err := e.QueryParentPosition().Exist(ctx); {
			case err != nil:
				return nil, errors.Wrapf(err, "checking existence of equipment parent position")
			case exist:
				switch parent, err := e.QueryParentPosition().QueryParent().QueryWorkOrder().Only(ctx); {
				case err != nil && !ent.IsNotFound(err):
					return nil, errors.Wrapf(err, "checking existence of equipment parent equipment work order")
				case parent != nil && parent.ID != wo.ID:
					return nil, errors.New("work order depend on another work order")
				}
			}
			e, err := r.ClientFrom(ctx).Equipment.
				UpdateOne(e).
				ClearWorkOrder().
				SetFutureState("").
				Save(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "install work order equipment e=%q, wo=%q", eid, id)
			}
			result.EquipmentAdded = append(result.EquipmentAdded, e)
		}
	}

	for _, l := range links {
		if l.FutureState == models.FutureStateInstall.String() {
			lid := l.ID
			switch exist, err := l.QueryPorts().
				QueryParent().
				Where(equipment.FutureState(models.FutureStateInstall.String())).
				Exist(ctx); {
			case err != nil:
				return nil, errors.Wrapf(err, "querying equipment link existence")
			case exist:
				return nil, errors.Errorf("installing link on equipment to be installed wo=%q", id)
			}
			l, err := r.ClientFrom(ctx).Link.
				UpdateOne(l).
				ClearWorkOrder().
				SetFutureState("").
				Save(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "install work order link l=%q, wo=%q", lid, id)
			}
			result.LinkAdded = append(result.LinkAdded, l)
		}
	}

	if err := r.ClientFrom(ctx).WorkOrder.
		UpdateOne(wo).
		SetStatus(models.WorkOrderStatusDone.String()).
		Exec(ctx); err != nil {
		return nil, errors.Wrapf(err, "Installing and removing work order items wo=%q", id)
	}
	return &result, nil
}

func (r mutationResolver) RemoveLocationType(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	lt, err := client.LocationType.Get(ctx, id)
	if err != nil {
		return id, errors.Wrapf(err, "location type does not exist: id=%q", id)
	}
	switch exist, err := lt.QueryLocations().Exist(ctx); {
	case err != nil:
		return id, errors.Wrapf(err, "querying locations for type: id=%q", id)
	case exist:
		return id, errors.Errorf("cannot delete location type with existing locations: id=%q", id)
	}
	if _, err := client.PropertyType.Delete().
		Where(propertytype.HasLocationTypeWith(locationtype.ID(id))).
		Exec(ctx); err != nil {
		return id, errors.Wrapf(err, "deleting property type: id=%q", id)
	}
	if err := client.LocationType.DeleteOne(lt).Exec(ctx); err != nil {
		return id, errors.Wrapf(err, "deleting location type: id=%q", id)
	}
	return id, nil
}

func (r mutationResolver) MarkSiteSurveyNeeded(ctx context.Context, locationID int, needed bool) (*ent.Location, error) {
	l, err := r.ClientFrom(ctx).
		Location.UpdateOneID(locationID).
		SetSiteSurveyNeeded(needed).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot set site survey requested: id=%q", locationID)
	}
	return l, nil
}

func (r mutationResolver) AddService(ctx context.Context, data models.ServiceCreateData) (*ent.Service, error) {
	if data.Status == nil {
		return nil, errors.New("status is a mandatory param")
	}
	client := r.ClientFrom(ctx)
	err := resolverutil.CheckServiceNameNotExist(ctx, client, data.Name)
	if err != nil {
		return nil, err
	}
	if data.ExternalID != nil {
		err := resolverutil.CheckServiceExternalIDNotExist(ctx, client, *data.ExternalID)
		if err != nil {
			return nil, err
		}
	}
	mutation := client.Service.Create().
		SetName(data.Name).
		SetStatus(data.Status.String()).
		SetNillableExternalID(data.ExternalID).
		SetTypeID(data.ServiceTypeID).
		AddUpstreamIDs(data.UpstreamServiceIds...)
	if data.CustomerID != nil {
		mutation.AddCustomerIDs(*data.CustomerID)
	}

	s, err := mutation.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating service: %w", err)
	}
	if _, err := r.AddProperties(data.Properties, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetService(s) },
	}); err != nil {
		return nil, fmt.Errorf("creating service properties: %w", err)
	}
	return s, nil
}

// nolint: funlen
func (r mutationResolver) EditService(ctx context.Context, data models.ServiceEditData) (*ent.Service, error) {
	client := r.ClientFrom(ctx)
	s, err := client.Service.Get(ctx, data.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying service: id=%q", data.ID)
	}

	st, err := s.QueryType().OnlyID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying service type id")
	}

	query := client.Service.UpdateOne(s)

	if data.ExternalID != nil && (s.ExternalID == nil || *s.ExternalID != *data.ExternalID) {
		err := resolverutil.CheckServiceExternalIDNotExist(ctx, client, *data.ExternalID)
		if err != nil {
			return nil, err
		}
		query.SetExternalID(*data.ExternalID)
	}

	if data.Name != nil && s.Name != *data.Name {
		err := resolverutil.CheckServiceNameNotExist(ctx, client, *data.Name)
		if err != nil {
			return nil, err
		}
		query.SetName(*data.Name)
	}

	if data.Status != nil {
		query.SetStatus(data.Status.String())
	}

	if data.UpstreamServiceIds != nil {
		oldUpstreamIds := s.QueryDownstream().IDsX(ctx)
		addedUpstreamIds, deletedUpstreamIds := resolverutil.GetDifferenceBetweenSlices(oldUpstreamIds, data.UpstreamServiceIds)
		query.RemoveUpstreamIDs(deletedUpstreamIds...).AddUpstreamIDs(addedUpstreamIds...)
	}

	if data.CustomerID != nil {
		oldCustomerIds := s.QueryCustomer().IDsX(ctx)
		addedCustomerIds, deletedCustomerIds := resolverutil.GetDifferenceBetweenSlices(oldCustomerIds, []int{*data.CustomerID})
		query.RemoveCustomerIDs(deletedCustomerIds...).AddCustomerIDs(addedCustomerIds...)
	}

	if s, err = query.Save(ctx); err != nil {
		return nil, fmt.Errorf("updating service %d: %w", data.ID, err)
	}

	var added, edited []*models.PropertyInput
	directPropertiesTypes, err := s.QueryProperties().QueryType().IDs(ctx)
	if err != nil {
		return nil, err
	}
	for _, input := range data.Properties {
		if r.isNewProp(directPropertiesTypes, input.ID, input.PropertyTypeID) {
			added = append(added, input)
		} else {
			if input.ID == nil {
				propID, err := s.QueryProperties().Where(property.HasTypeWith(propertytype.ID(input.PropertyTypeID))).OnlyID(ctx)
				if err != nil {
					return nil, err
				}
				input.ID = &propID
			}
			edited = append(edited, input)
		}
	}
	if _, err := r.AddProperties(added, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetService(s) },
	}); err != nil {
		return nil, err
	}
	for _, input := range edited {
		typ, err := client.ServiceType.Query().
			Where(servicetype.ID(st)).
			QueryPropertyTypes().
			Where(propertytype.ID(input.PropertyTypeID)).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("querying service property type %d: %w", input.PropertyTypeID, err)
		}
		if typ.Editable && typ.IsInstanceProperty {
			query := client.Property.
				Update().
				Where(
					property.HasServiceWith(service.ID(s.ID)),
					property.ID(*input.ID),
				)
			if r.updatePropValues(ctx, input, query) != nil {
				return nil, fmt.Errorf("updating property values: %w", err)
			}
		}
	}
	return s, nil
}

func (r mutationResolver) AddServiceLink(ctx context.Context, id, linkID int) (*ent.Service, error) {
	svc, err := r.ClientFrom(ctx).
		Service.
		UpdateOneID(id).
		AddLinkIDs(linkID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("adding link %q to service %q: %w", linkID, id, err)
	}
	return svc, nil
}

func (r mutationResolver) RemoveServiceLink(ctx context.Context, id, linkID int) (*ent.Service, error) {
	svc, err := r.ClientFrom(ctx).
		Service.
		UpdateOneID(id).
		RemoveLinkIDs(linkID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("removing link %q from service %q: %w", linkID, id, err)
	}
	return svc, nil
}

func (r mutationResolver) AddServiceType(ctx context.Context, data models.ServiceTypeCreateData) (*ent.ServiceType, error) {
	types, err := r.AddPropertyTypes(ctx, data.Properties...)
	if err != nil {
		return nil, errors.WithMessage(err, "creating service type properties")
	}

	epTypes, err := r.addServiceEndpointDefinitions(ctx, data.Endpoints...)
	if err != nil {
		return nil, errors.WithMessage(err, "creating service endpoint definition")
	}

	st, err := r.ClientFrom(ctx).
		ServiceType.Create().
		SetName(data.Name).
		SetHasCustomer(data.HasCustomer).
		AddPropertyTypes(types...).
		AddEndpointDefinitions(epTypes...).
		Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("creating service type. error=%v", err.Error())
	}
	return st, nil
}

func (r mutationResolver) EditServiceType(ctx context.Context, data models.ServiceTypeEditData) (*ent.ServiceType, error) {
	typ, err := r.ClientFrom(ctx).
		ServiceType.UpdateOneID(data.ID).
		SetName(data.Name).
		SetHasCustomer(data.HasCustomer).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A service type with the name %v already exists", data.Name)
		}
		return nil, errors.Wrapf(err, "updating service type: id=%q", data.ID)
	}
	for _, input := range data.Properties {
		if input.ID == nil {
			err = r.validateAndAddNewPropertyType(
				ctx, input, func(b *ent.PropertyTypeUpdateOne) {
					b.SetServiceType(typ)
				},
			)
		} else {
			err = r.updatePropType(ctx, input)
		}
		if err != nil {
			return nil, err
		}
	}

	// update indexes to be invalid number before re-assigning (field is unique+mandatory)
	for _, input := range data.Endpoints {
		if input.ID != nil {
			ept, err := r.ClientFrom(ctx).ServiceEndpointDefinition.Get(ctx, *input.ID)
			if err != nil {
				return nil, err
			}
			_, err = ept.Update().SetIndex(-1*ept.Index - 1).Save(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, input := range data.Endpoints {
		if input.ID == nil {
			err = r.validateAndAddEndpointDefinition(ctx, input, data.ID)
		} else {
			err = r.updateEndpointDefinition(ctx, input, data.ID)
		}
		if err != nil {
			return nil, err
		}
	}

	return typ, nil

}

func (r mutationResolver) RemoveServiceType(ctx context.Context, id int) (int, error) {
	client := r.ClientFrom(ctx)
	st, err := client.ServiceType.Get(ctx, id)
	if err != nil {
		return id, fmt.Errorf("getting service type %d: %w", id, err)
	}
	switch exist, err := st.QueryServices().Exist(ctx); {
	case err != nil:
		return id, fmt.Errorf("querying services for type %d: %w", id, err)
	case exist:
		return id, fmt.Errorf("cannot delete service type %d with existing services", id)
	}
	if _, err := client.Property.Delete().
		Where(property.HasServiceWith(service.HasTypeWith(servicetype.ID(st.ID)))).
		Exec(ctx); err != nil {
		return id, fmt.Errorf("deleting service type %d properties: %w", id, err)
	}
	if err := client.ServiceType.DeleteOne(st).Exec(ctx); err != nil {
		return id, fmt.Errorf("deleting service type %q: %w", id, err)
	}
	return id, nil
}

func (r mutationResolver) EditEquipment(
	ctx context.Context, input models.EditEquipmentInput,
) (*ent.Equipment, error) {
	client := r.ClientFrom(ctx)
	e, err := client.Equipment.Get(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying equipment: id=%q", input.ID)
	}

	var added, edited []*models.PropertyInput
	for _, input := range input.Properties {
		if input.ID == nil {
			added = append(added, input)
		} else {
			edited = append(edited, input)
		}
	}
	if _, err := r.AddProperties(added, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetEquipment(e) },
	}); err != nil {
		return nil, err
	}

	if e.Name != input.Name {
		var lid, pid *int
		l, err := e.QueryLocation().FirstID(ctx)
		if err == nil {
			lid = &l
		}
		p, err := e.QueryParentPosition().FirstID(ctx)
		if err == nil {
			pid = &p
		}
		if err := r.validateEquipmentNameIsUnique(ctx, input.Name, lid, pid, &e.ID); err != nil {
			return nil, err
		}
	}

	if e.Name != input.Name || input.DeviceID != nil && e.DeviceID != *input.DeviceID {
		upd := client.Equipment.
			UpdateOne(e).
			SetName(input.Name).
			SetNillableDeviceID(input.DeviceID)
		if input.ExternalID != nil && *input.ExternalID != "" {
			upd.SetNillableExternalID(input.ExternalID)
		} else {
			upd.ClearExternalID()
		}
		if e, err = upd.Save(ctx); err != nil {
			return nil, errors.Wrapf(err, "updating equipment: id=%q", input.ID)
		}
	}

	for _, input := range edited {
		et, err := e.QueryType().OnlyID(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment type: id=%q", e.ID)
		}
		typ, err := client.
			EquipmentType.Query().
			Where(equipmenttype.ID(et)).
			QueryPropertyTypes().
			Where(propertytype.ID(input.PropertyTypeID)).
			Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment property type %q", input.PropertyTypeID)
		}
		if typ.Editable && typ.IsInstanceProperty {
			updater := client.Property.Update().
				Where(
					property.HasEquipmentWith(equipment.ID(e.ID)),
					property.ID(*input.ID),
				)
			if r.updatePropValues(ctx, input, updater) != nil {
				return nil, errors.Wrap(err, "updating property values")
			}
		}
	}
	return e, nil
}

// TODO T58981969 Add isNewProp to all edit mutations
func (r mutationResolver) isNewProp(directPropertiesTypes []int, propertyID *int, propertyTypeID int) bool {
	if propertyID != nil {
		return false
	}
	for _, id := range directPropertiesTypes {
		if id == propertyTypeID {
			return false
		}
	}
	return true
}

func (r mutationResolver) EditEquipmentPort(
	ctx context.Context, input models.EditEquipmentPortInput,
) (*ent.EquipmentPort, error) {
	client := r.ClientFrom(ctx)
	p, err := r.getOrCreatePort(ctx, input.Side)
	if err != nil || p == nil {
		return nil, err
	}

	var added, edited []*models.PropertyInput
	directPropertiesTypes, err := p.QueryProperties().QueryType().IDs(ctx)
	if err != nil {
		return nil, err
	}
	for _, input := range input.Properties {
		if r.isNewProp(directPropertiesTypes, input.ID, input.PropertyTypeID) {
			added = append(added, input)
		} else {
			if input.ID == nil {
				propID, err := p.QueryProperties().Where(property.HasTypeWith(propertytype.ID(input.PropertyTypeID))).OnlyID(ctx)
				if err != nil {
					return nil, err
				}
				input.ID = &propID
			}
			edited = append(edited, input)
		}
	}
	if _, err = r.AddProperties(added, resolverutil.AddPropertyArgs{
		Context:   ctx,
		EntSetter: func(b *ent.PropertyCreate) { b.SetEquipmentPort(p) },
	}); err != nil {
		return nil, err
	}

	for _, input := range edited {
		def, err := p.QueryDefinition().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment port definition type %q", input.PropertyTypeID)
		}
		id, err := def.QueryEquipmentPortType().OnlyID(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment port type type %q", input.PropertyTypeID)
		}
		typ, err := client.
			EquipmentPortType.Query().
			Where(equipmentporttype.ID(id)).
			QueryPropertyTypes().
			Where(propertytype.ID(input.PropertyTypeID)).
			Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying equipment port property type %q", input.PropertyTypeID)
		}
		if typ.Editable && typ.IsInstanceProperty {
			updater := client.Property.
				Update().
				Where(
					property.HasEquipmentPortWith(equipmentport.ID(p.ID)),
					property.ID(*input.ID),
				)
			if r.updatePropValues(ctx, input, updater) != nil {
				return nil, errors.Wrap(err, "updating property values")
			}
		}
	}
	return p, nil
}

func (r mutationResolver) validateEquipmentNameIsUnique(ctx context.Context, name string, locationID, positionID, equipID *int) error {
	query := r.ClientFrom(ctx).Equipment.Query().Where(equipment.Name(name))
	if equipID != nil {
		query = query.Where(equipment.IDNEQ(*equipID))
	}
	if positionID != nil {
		query = query.Where(equipment.HasParentPositionWith(equipmentposition.ID(*positionID)))
	} else if locationID != nil {
		query = query.Where(equipment.HasLocationWith(location.ID(*locationID)))
	}
	exist, err := query.Exist(ctx)
	if err != nil {
		return errors.Wrapf(err, "error querying equipment existence for %q", name)
	}
	if exist {
		var parentName interface{}
		if locationID != nil {
			parent, err := r.ClientFrom(ctx).Location.Get(ctx, *locationID)
			if err != nil {
				return errors.Wrapf(err, "error querying equipment location for %q", *locationID)
			}
			parentName = parent.Name
		} else if positionID != nil {
			parent, err := r.ClientFrom(ctx).EquipmentPosition.Get(ctx, *positionID)
			if err != nil {
				return errors.Wrapf(err, "error querying equipment position for position id %q", *positionID)
			}
			parentName = parent.ID
		}
		r.logger.For(ctx).Error(
			"duplicate equipment name",
			zap.String("name", name),
			zap.Any("parent", parentName))
		return gqlerror.Errorf("An equipment with the name %v already exists under %v", name, parentName)
	}
	return nil
}

func (r mutationResolver) validateAndAddNewPropertyType(ctx context.Context, input *models.PropertyTypeInput, entSetter func(*ent.PropertyTypeUpdateOne)) error {
	isEmpty, err := r.isEmptyProp(nil, input)
	if err != nil {
		return err
	}
	if isEmpty {
		return gqlerror.Errorf("The new property %v must have a default value", input.Name)
	}
	types, err := r.AddPropertyTypes(ctx, input)
	if err != nil || len(types) == 0 {
		return err
	}
	query := r.ClientFrom(ctx).
		PropertyType.
		UpdateOne(types[0])
	entSetter(query)
	if _, err =
		query.
			Save(ctx); ent.IsConstraintError(err) {
		return gqlerror.Errorf("A property type with the name %v already exists under in the selected object", input.Name)
	}
	return err
}

func (r mutationResolver) validateAndAddEndpointDefinition(ctx context.Context, input *models.ServiceEndpointDefinitionInput, serviceTypeID int) error {
	if input == nil {
		return nil
	}
	_, err := r.addServiceEndpointDefinition(ctx, *input, serviceTypeID)
	return err
}

func (r mutationResolver) EditLocationTypesIndex(ctx context.Context, locationTypesIndex []*models.LocationTypeIndex) ([]*ent.LocationType, error) {
	var updated []*ent.LocationType
	client := r.ClientFrom(ctx)
	for _, obj := range locationTypesIndex {
		lt, err := client.LocationType.Get(ctx, obj.LocationTypeID)
		if err != nil {
			r.logger.For(ctx).Error("couldn't fetch location type",
				zap.Int("id", obj.LocationTypeID),
			)
			return nil, gqlerror.Errorf("couldn't fetch location type. id=%q", obj.LocationTypeID)
		}
		saved, err := lt.Update().SetIndex(obj.Index).Save(ctx)
		if err != nil {
			r.logger.For(ctx).Error("couldn't update location type",
				zap.Int("id", obj.LocationTypeID),
				zap.Int("index", obj.Index),
			)
			return nil, gqlerror.Errorf("couldn't update location type. id=%q, index=%q", obj.LocationTypeID, obj.Index)
		}
		updated = append(updated, saved)
	}
	return updated, nil
}

func (r mutationResolver) EditLocationType(
	ctx context.Context, input models.EditLocationTypeInput,
) (*ent.LocationType, error) {
	typ, err := r.ClientFrom(ctx).
		LocationType.UpdateOneID(input.ID).
		SetName(input.Name).
		SetNillableMapType(input.MapType).
		SetNillableMapZoomLevel(input.MapZoomLevel).
		SetNillableSite(input.IsSite).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, gqlerror.Errorf("A location type with the name %v already exists", input.Name)
		}
		return nil, errors.Wrapf(err, "updating location type: id=%q", input.ID)
	}
	for _, input := range input.Properties {
		if input.ID == nil {
			err = r.validateAndAddNewPropertyType(
				ctx, input, func(b *ent.PropertyTypeUpdateOne) {
					b.SetLocationType(typ)
				},
			)
		} else {
			err = r.updatePropType(ctx, input)
		}
		if err != nil {
			return nil, err
		}
	}
	return typ, nil
}

func (r mutationResolver) EditLocationTypeSurveyTemplateCategories(
	ctx context.Context, id int, surveyTemplateCategories []*models.SurveyTemplateCategoryInput,
) ([]*ent.SurveyTemplateCategory, error) {
	var (
		categories = make([]*ent.SurveyTemplateCategory, len(surveyTemplateCategories))
		keepIDs    = make(map[int]bool)
		added      []*ent.SurveyTemplateCategory
		err        error
	)
	for i, input := range surveyTemplateCategories {
		if input.ID == nil {
			category, err := r.AddSurveyTemplateCategories(ctx, input)
			if err != nil {
				return nil, err
			}
			categories[i] = category[0]
			added = append(added, category[0])
		} else {
			keepIDs[*input.ID] = true
			if categories[i], err = r.updateSurveyTemplateCategory(ctx, input); err != nil {
				return nil, err
			}
		}
	}

	lt, err := r.ClientFrom(ctx).LocationType.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch location type: id=%q", id)
	}

	existingCategories, err := r.ClientFrom(ctx).LocationType.QuerySurveyTemplateCategories(lt).All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch survey template categories for location type: id=%q", id)
	}

	var deleteIDs []int
	for _, existingCategory := range existingCategories {
		if _, ok := keepIDs[existingCategory.ID]; !ok {
			deleteIDs = append(deleteIDs, existingCategory.ID)
		}
	}

	if err := r.ClientFrom(ctx).
		LocationType.
		UpdateOneID(id).
		AddSurveyTemplateCategories(added...).
		RemoveSurveyTemplateCategoryIDs(deleteIDs...).
		Exec(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to update survey categories for location type")
	}

	return categories, nil
}

func (r mutationResolver) updateEquipmentTypeCategory(ctx context.Context, client *ent.Client, et *ent.EquipmentType, category string) (*ent.EquipmentType, error) {
	c, err := client.EquipmentCategory.Query().Where(equipmentcategory.Name(category)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, errors.Wrapf(err, "querying equipment category with name %v", category)
		}
		c, err = client.EquipmentCategory.Create().SetName(category).Save(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "creating equipment category with name %v", category)
		}
	}
	et, err = client.EquipmentType.UpdateOne(et).SetCategory(c).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "updating equipment category")
	}
	return et, nil
}

func (r mutationResolver) EditEquipmentType(
	ctx context.Context, input models.EditEquipmentTypeInput,
) (et *ent.EquipmentType, err error) {
	client := r.ClientFrom(ctx)
	if et, err = client.EquipmentType.Get(ctx, input.ID); err != nil {
		return nil, errors.Wrapf(err, "querying equipment type: id=%q", input.ID)
	}
	if input.Name != et.Name {
		if et, err = client.EquipmentType.
			UpdateOne(et).
			SetName(input.Name).
			Save(ctx); err != nil {
			if ent.IsConstraintError(err) {
				return nil, gqlerror.Errorf("An equipment type with the name %v already exists", input.Name)
			}
			return nil, errors.Wrap(err, "updating equipment type name")
		}
	}

	if input.Category == nil {
		switch exist, err := et.QueryCategory().Exist(ctx); {
		case err != nil:
			return nil, errors.Wrap(err, "querying category existence")
		case exist:
			et, err = client.EquipmentType.UpdateOne(et).ClearCategory().Save(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "clearing equipment category")
			}
		}
	} else if et, err = r.updateEquipmentTypeCategory(ctx, client, et, *input.Category); err != nil {
		return nil, errors.Wrap(err, "updating equipment category")
	}

	for _, input := range input.Properties {
		if input.ID == nil {
			err = r.validateAndAddNewPropertyType(
				ctx, input, func(b *ent.PropertyTypeUpdateOne) {
					b.SetEquipmentTypeID(et.ID)
				},
			)
		} else {
			err = r.updatePropType(ctx, input)
		}
		if err != nil {
			return nil, err
		}
	}

	{
		var added, edited []*models.EquipmentPortInput
		for _, input := range input.Ports {
			if input.ID == nil {
				added = append(added, input)
			} else {
				edited = append(edited, input)
			}
		}
		if _, err := r.AddEquipmentPortDefinitions(ctx, added, &et.ID); err != nil {
			return nil, err
		}
		for _, input := range edited {
			if err := client.EquipmentPortDefinition.
				UpdateOneID(*input.ID).
				SetName(input.Name).
				SetNillableIndex(input.Index).
				SetNillableBandwidth(input.Bandwidth).
				SetNillableVisibilityLabel(input.VisibleLabel).
				Exec(ctx); err != nil {
				return nil, errors.Wrapf(err, "updating equipment port definition: id=%q", *input.ID)
			}
		}
	}

	{
		var added, edited []*models.EquipmentPositionInput
		for _, input := range input.Positions {
			if input.ID == nil {
				added = append(added, input)
			} else {
				edited = append(edited, input)
			}
		}
		if _, err := r.AddEquipmentPositionDefinitions(ctx, added, &et.ID); err != nil {
			return nil, err
		}
		for _, input := range edited {
			if err := client.EquipmentPositionDefinition.
				UpdateOneID(*input.ID).
				SetName(input.Name).
				SetNillableIndex(input.Index).
				SetNillableVisibilityLabel(input.VisibleLabel).
				SetEquipmentType(et).
				Exec(ctx); err != nil {
				return nil, errors.Wrapf(err, "updating equipment position definition: id=%q", *input.ID)
			}
		}
	}
	return et, nil
}

func (r mutationResolver) EditEquipmentPortType(
	ctx context.Context, input models.EditEquipmentPortTypeInput,
) (*ent.EquipmentPortType, error) {
	client := r.ClientFrom(ctx)
	pt, err := client.EquipmentPortType.Get(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying equipment port type: id=%q", input.ID)
	}
	if input.Name != pt.Name {
		if pt, err = client.EquipmentPortType.
			UpdateOne(pt).
			SetName(input.Name).
			Save(ctx); err != nil {
			if ent.IsConstraintError(err) {
				return nil, gqlerror.Errorf("An equipment port type with the name %s already exists", input.Name)
			}
			return nil, errors.Wrap(err, "updating equipment port type")
		}
	}

	for _, input := range input.Properties {
		if input.ID == nil {
			if err := r.validateAndAddNewPropertyType(ctx, input,
				func(b *ent.PropertyTypeUpdateOne) {
					b.SetEquipmentPortTypeID(pt.ID)
				}); err != nil {
				return nil, err
			}
		} else {
			if _, err := client.
				EquipmentPortType.Query().
				QueryPropertyTypes().
				Where(propertytype.ID(*input.ID)).
				Only(ctx); err != nil {
				return nil, gqlerror.Errorf("%v error querying property type %v (id: %v)", err, input.Name, input.ID)
			}
			if err := r.updatePropType(ctx, input); err != nil {
				return nil, err
			}
		}
	}
	for _, input := range input.LinkProperties {
		if input.ID == nil {
			if err := r.validateAndAddNewPropertyType(ctx, input,
				func(b *ent.PropertyTypeUpdateOne) {
					b.SetLinkEquipmentPortTypeID(pt.ID)
				}); err != nil {
				return nil, err
			}
		} else {
			if _, err := client.
				EquipmentPortType.Query().
				QueryLinkPropertyTypes().
				Where(propertytype.ID(*input.ID)).
				Only(ctx); err != nil {
				return nil, gqlerror.Errorf("%v error querying link property type %v (id: %v)", err, input.Name, input.ID)
			}
			if err := r.updatePropType(ctx, input); err != nil {
				return nil, err
			}
		}
	}
	return pt, nil
}

func (r mutationResolver) DeleteLocationTypeEquipments(ctx context.Context, locationTypeID int, blacklistedLocationIds []int, limit int) (int, error) {
	equipments, err := r.ClientFrom(ctx).
		EquipmentType.Query().
		QueryEquipment().
		Where(
			equipment.HasLocationWith(
				location.IDNotIn(blacklistedLocationIds...),
				location.HasTypeWith(
					locationtype.ID(locationTypeID),
				),
			),
		).
		Limit(limit).
		All(ctx)
	if err != nil {
		return 0, errors.Wrapf(err, "querying equipments")
	}

	for _, e := range equipments {
		if err := r.removeEquipment(ctx, e); err != nil {
			return 0, err
		}
	}
	return len(equipments), nil
}

func (r mutationResolver) updatePropValues(ctx context.Context, input *models.PropertyInput, pu *ent.PropertyUpdate) error {
	pu = pu.SetNillableStringVal(input.StringValue).
		SetNillableIntVal(input.IntValue).
		SetNillableBoolVal(input.BooleanValue).
		SetNillableFloatVal(input.FloatValue).
		SetNillableLatitudeVal(input.LatitudeValue).
		SetNillableLongitudeVal(input.LongitudeValue).
		SetNillableRangeFromVal(input.RangeFromValue).
		SetNillableRangeToVal(input.RangeToValue)

	if input.NodeIDValue != nil {
		if err := r.setNodePropertyUpdate(ctx, pu, *input.NodeIDValue); err != nil {
			return err
		}
	} else {
		pu = pu.ClearEquipmentValue()
		pu = pu.ClearLocationValue()
		pu = pu.ClearServiceValue()
	}

	return pu.Exec(ctx)
}

func (r mutationResolver) updatePropType(ctx context.Context, input *models.PropertyTypeInput) error {
	if err := r.ClientFrom(ctx).PropertyType.
		UpdateOneID(*input.ID).
		SetName(input.Name).
		SetType(input.Type.String()).
		SetNillableNodeType(input.NodeType).
		SetNillableIndex(input.Index).
		SetNillableExternalID(input.ExternalID).
		SetNillableStringVal(input.StringValue).
		SetNillableIntVal(input.IntValue).
		SetNillableBoolVal(input.BooleanValue).
		SetNillableFloatVal(input.FloatValue).
		SetNillableLatitudeVal(input.LatitudeValue).
		SetNillableLongitudeVal(input.LongitudeValue).
		SetNillableRangeFromVal(input.RangeFromValue).
		SetNillableRangeToVal(input.RangeToValue).
		SetNillableIsInstanceProperty(input.IsInstanceProperty).
		SetNillableEditable(input.IsEditable).
		SetNillableMandatory(input.IsMandatory).
		SetNillableDeleted(input.IsDeleted).
		Exec(ctx); err != nil {
		return errors.Wrap(err, "updating property type")
	}
	return nil
}

func (r mutationResolver) updateEndpointDefinition(ctx context.Context, input *models.ServiceEndpointDefinitionInput, serviceTypeID int) error {
	if err := r.ClientFrom(ctx).ServiceEndpointDefinition.
		UpdateOneID(*input.ID).
		SetName(input.Name).
		SetNillableRole(input.Role).
		SetIndex(input.Index).
		SetEquipmentTypeID(input.EquipmentTypeID).
		SetServiceTypeID(serviceTypeID).
		Exec(ctx); err != nil {
		return errors.Wrap(err, "updating service endpoint definition")
	}
	return nil
}

func (r mutationResolver) updateSurveyTemplateCategory(ctx context.Context, input *models.SurveyTemplateCategoryInput) (*ent.SurveyTemplateCategory, error) {
	updater := r.ClientFrom(ctx).SurveyTemplateCategory.UpdateOneID(*input.ID)
	keepIDs := make(map[int]bool)
	for _, questionInput := range input.SurveyTemplateQuestions {
		if questionInput.ID == nil {
			question, err := r.AddSurveyTemplateQuestions(ctx, questionInput)
			if err != nil {
				return nil, err
			}
			updater.AddSurveyTemplateQuestions(question...)
		} else {
			if err := r.updateSurveyTemplateQuestion(ctx, questionInput); err != nil {
				return nil, err
			}
			keepIDs[*questionInput.ID] = true
		}
	}

	category, err := r.ClientFrom(ctx).SurveyTemplateCategory.Get(ctx, *input.ID)
	if err != nil {
		return nil, fmt.Errorf("querying survey template category %d: %w", *input.ID, err)
	}

	existingQuestions, err := category.QuerySurveyTemplateQuestions().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying survey template questions for category %d: %w", *input.ID, err)
	}

	var deleteIDs []int
	for _, existingQuestion := range existingQuestions {
		if _, ok := keepIDs[existingQuestion.ID]; !ok {
			deleteIDs = append(deleteIDs, existingQuestion.ID)
		}
	}

	category, err = updater.
		RemoveSurveyTemplateQuestionIDs(deleteIDs...).
		SetCategoryTitle(input.CategoryTitle).
		SetCategoryDescription(input.CategoryDescription).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating survey template category: %w", err)
	}
	return category, nil
}

func (r mutationResolver) updateSurveyTemplateQuestion(ctx context.Context, input *models.SurveyTemplateQuestionInput) error {
	if err := r.ClientFrom(ctx).SurveyTemplateQuestion.
		UpdateOneID(*input.ID).
		SetQuestionTitle(input.QuestionTitle).
		SetQuestionDescription(input.QuestionDescription).
		SetQuestionType(input.QuestionType.String()).
		Exec(ctx); err != nil {
		return fmt.Errorf("updating survey template question: %w", err)
	}
	return nil
}

func (r mutationResolver) MarkLocationPropertyAsExternalID(ctx context.Context, name string) (string, error) {
	client := r.ClientFrom(ctx).Location
	sites, err := client.Query().
		Where(
			location.HasPropertiesWith(
				property.HasTypeWith(
					propertytype.Name(name),
				),
			),
		).
		All(ctx)
	if err != nil {
		return "", errors.Wrap(err, "querying locations with property")
	}

	for _, site := range sites {
		p, err := site.QueryProperties().
			Where(
				property.HasTypeWith(
					propertytype.Name(name),
				),
			).
			Only(ctx)
		if err != nil {
			return "", fmt.Errorf("querying property type: %w", err)
		}
		if err := client.UpdateOne(site).
			SetExternalID(p.StringVal).
			Exec(ctx); err != nil {
			return "", fmt.Errorf("updating external id: %w", err)
		}
	}
	return name, nil
}

func (r mutationResolver) deleteLocationHierarchy(ctx context.Context, l *ent.Location) error {
	children, err := l.QueryChildren().All(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed querying child locations l=%v", l.ID)
	}
	for _, child := range children {
		if err := r.deleteLocationHierarchy(ctx, child); err != nil {
			return err
		}
	}
	err = r.ClientFrom(ctx).Location.DeleteOne(l).Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete location l=%v", l.ID)
	}
	return nil
}

func (r mutationResolver) DeleteLocationHierarchy(ctx context.Context, id int) (int, error) {
	l, err := r.ClientFrom(ctx).Location.Get(ctx, id)
	if err != nil {
		return id, errors.Wrapf(err, "can't query location l=%v", id)
	}
	return id, r.deleteLocationHierarchy(ctx, l)
}

func (r mutationResolver) MoveLocation(ctx context.Context, locationID int, parentLocationID *int) (*ent.Location, error) {
	client := r.ClientFrom(ctx)
	l, err := client.Location.Get(ctx, locationID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying location: id=%q", locationID)
	}
	if parentLocationID == nil {
		// location becoming root which requires validation, see comment in AddLocation
		if err := r.validateRootLocationUniqueness(ctx, l.QueryType().OnlyXID(ctx), l.Name); err != nil {
			return nil, err
		}
		return client.Location.
			UpdateOne(l).
			ClearParent().
			Save(ctx)
	}
	newParentID := *parentLocationID
	newParent, err := client.Location.Get(ctx, newParentID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying parent location: parent id=%q", *parentLocationID)
	}
	parentAncestors, err := r.Location().LocationHierarchy(ctx, newParent)
	if err != nil {
		return nil, errors.Wrapf(err, "querying parent ancestors: parent id=%q", *parentLocationID)
	}
	for _, parentAncestor := range parentAncestors {
		if parentAncestor.ID == l.ID {
			return nil, errors.Errorf("new parent (%q)is a descendant of the location (%q)", *parentLocationID, locationID)
		}
	}
	if err := r.verifyLocationParent(ctx, l.QueryType().OnlyX(ctx).ID, newParentID); err != nil {
		return nil, err
	}
	if l, err = client.Location.
		UpdateOne(l).
		SetParentID(newParentID).
		Save(ctx); err != nil {
		return nil, errors.Wrapf(err, "updating location parent: id=%q, parent id=%q", locationID, *parentLocationID)
	}
	return l, nil
}

func (r mutationResolver) AddTechnician(
	ctx context.Context, input models.TechnicianInput,
) (*ent.Technician, error) {
	t, err := r.ClientFrom(ctx).
		Technician.Create().
		SetName(input.Name).
		SetEmail(input.Email).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating technician")
	}
	return t, nil
}

func (r mutationResolver) AddCustomer(ctx context.Context, input models.AddCustomerInput) (*ent.Customer, error) {
	exist, _ := r.ClientFrom(ctx).Customer.Query().Where(customer.Name(input.Name)).Exist(ctx)
	if exist {
		return nil, gqlerror.Errorf("A customer with the name %v already exists", input.Name)
	}

	if input.ExternalID != nil {
		exist, _ = r.ClientFrom(ctx).Customer.Query().Where(customer.ExternalID(*input.ExternalID)).Exist(ctx)
		if exist {
			return nil, gqlerror.Errorf("A customer with the external id %v already exists", *input.ExternalID)
		}
	}

	t, err := r.ClientFrom(ctx).
		Customer.Create().
		SetName(input.Name).
		SetNillableExternalID(input.ExternalID).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating customer")
	}
	return t, nil
}

func (r mutationResolver) RemoveCustomer(ctx context.Context, id int) (int, error) {
	if err := r.ClientFrom(ctx).Customer.DeleteOneID(id).Exec(ctx); err != nil {
		return id, errors.Wrap(err, "removing customer")
	}
	return id, nil
}

func actionsInputToSchema(ctx context.Context, inputActions []*models.ActionsRuleActionInput) ([]*core.ActionsRuleAction, error) {
	ac := actions.FromContext(ctx)
	ruleActions := make([]*core.ActionsRuleAction, 0, len(inputActions))
	for _, ruleAction := range inputActions {
		_, err := ac.ActionForID(ruleAction.ActionID)
		if err != nil {
			return nil, errors.Wrap(err, "validating action")
		}

		ruleActions = append(ruleActions, &core.ActionsRuleAction{
			ActionID: ruleAction.ActionID,
			Data:     ruleAction.Data,
		})
	}
	return ruleActions, nil
}

func filtersInputToSchema(inputFilters []*models.ActionsRuleFilterInput) []*core.ActionsRuleFilter {
	ruleFilters := make([]*core.ActionsRuleFilter, 0, len(inputFilters))
	for _, ruleFilter := range inputFilters {
		ruleFilters = append(ruleFilters, &core.ActionsRuleFilter{
			FilterID:   ruleFilter.FilterID,
			OperatorID: ruleFilter.OperatorID,
			Data:       ruleFilter.Data,
		})
	}
	return ruleFilters
}

func (r mutationResolver) AddActionsRule(ctx context.Context, input models.AddActionsRuleInput) (*ent.ActionsRule, error) {
	ac := actions.FromContext(ctx)

	_, err := ac.TriggerForID(input.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "validating trigger")
	}

	ruleActions, err := actionsInputToSchema(ctx, input.RuleActions)
	if err != nil {
		return nil, errors.Wrap(err, "validating action")
	}

	ruleFilters := filtersInputToSchema(input.RuleFilters)

	actionsRule, err := r.ClientFrom(ctx).
		ActionsRule.Create().
		SetName(input.Name).
		SetTriggerID(string(input.TriggerID)).
		SetRuleActions(ruleActions).
		SetRuleFilters(ruleFilters).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating actionsrule")
	}
	return actionsRule, nil
}

func (r mutationResolver) AddFloorPlan(ctx context.Context, input models.AddFloorPlanInput) (*ent.FloorPlan, error) {
	img, err := r.createImage(ctx, input.Image)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create image")
	}

	client := r.ClientFrom(ctx)
	referencePoint, err := client.FloorPlanReferencePoint.Create().
		SetX(input.ReferenceX).
		SetY(input.ReferenceY).
		SetLatitude(input.Latitude).
		SetLongitude(input.Longitude).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create reference point")
	}

	scale, err := client.FloorPlanScale.Create().
		SetReferencePoint1X(input.ReferencePoint1x).
		SetReferencePoint1Y(input.ReferencePoint1y).
		SetReferencePoint2X(input.ReferencePoint2x).
		SetReferencePoint2Y(input.ReferencePoint2y).
		SetScaleInMeters(input.ScaleInMeters).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scale")
	}

	floorPlan, err := client.FloorPlan.Create().
		SetName(input.Name).
		SetLocationID(input.LocationID).
		SetImage(img).
		SetReferencePoint(referencePoint).
		SetScale(scale).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create floor plan")
	}

	return floorPlan, nil
}

func (r mutationResolver) EditActionsRule(ctx context.Context, id int, input models.AddActionsRuleInput) (*ent.ActionsRule, error) {
	ac := actions.FromContext(ctx)

	_, err := ac.TriggerForID(input.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "validating trigger")
	}

	ruleActions, err := actionsInputToSchema(ctx, input.RuleActions)
	if err != nil {
		return nil, errors.Wrap(err, "validating action")
	}

	ruleFilters := filtersInputToSchema(input.RuleFilters)

	actionsRule, err := r.ClientFrom(ctx).
		ActionsRule.UpdateOneID(id).
		SetName(input.Name).
		SetTriggerID(string(input.TriggerID)).
		SetRuleActions(ruleActions).
		SetRuleFilters(ruleFilters).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "updating actionsrule")
	}
	return actionsRule, nil
}

func (r mutationResolver) RemoveActionsRule(ctx context.Context, id int) (_ bool, err error) {
	if err = r.ClientFrom(ctx).ActionsRule.DeleteOneID(id).Exec(ctx); err != nil {
		err = fmt.Errorf("removing actions rule: %w", err)
	}
	return err == nil, err
}

func (r mutationResolver) DeleteFloorPlan(ctx context.Context, id int) (_ bool, err error) {
	if err = r.ClientFrom(ctx).FloorPlan.DeleteOneID(id).Exec(ctx); err != nil {
		err = fmt.Errorf("deleting floor plan %q: err %w", id, err)
	}
	return err == nil, err
}

func (r mutationResolver) TechnicianWorkOrderCheckIn(ctx context.Context, id int) (*ent.WorkOrder, error) {
	client := r.ClientFrom(ctx).WorkOrder
	wo, err := client.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting work order %q: %w", id, err)
	}
	if wo.Status != models.WorkOrderStatusPlanned.String() {
		return wo, nil
	}
	if wo, err = wo.Update().
		SetStatus(models.WorkOrderStatusPending.String()).
		Save(ctx); err != nil {
		return nil, fmt.Errorf("updating work order %q status to pending: %w", id, err)
	}
	if _, err = r.AddComment(ctx, models.CommentInput{
		EntityType: models.CommentEntityWorkOrder,
		ID:         id,
		Text:       r.Me(ctx).User + " checked-in",
	}); err != nil {
		return nil, fmt.Errorf("adding technician check-in comment: %w", err)
	}
	return wo, nil
}

func validateFilterTypeEntity(input models.ReportFilterInput) error {
	for _, f := range input.Filters {
		if f.Key == "" {
			return fmt.Errorf("filter key was not provided for %q", input.Entity)
		}
		var validator interface{ IsValid() bool }
		switch input.Entity {
		case models.FilterEntityEquipment:
			validator = models.EquipmentFilterType(f.FilterType)
		case models.FilterEntityLink:
			validator = models.LinkFilterType(f.FilterType)
		case models.FilterEntityLocation:
			validator = models.LocationFilterType(f.FilterType)
		case models.FilterEntityPort:
			validator = models.PortFilterType(f.FilterType)
		case models.FilterEntityService:
			validator = models.ServiceFilterType(f.FilterType)
		case models.FilterEntityWorkOrder:
			validator = models.WorkOrderFilterType(f.FilterType)
		}
		if validator == nil || !validator.IsValid() {
			return fmt.Errorf("entity %q and filter type %q does not match", input.Entity, f.FilterType)
		}
	}
	return nil
}

func (r mutationResolver) AddReportFilter(ctx context.Context, input models.ReportFilterInput) (*ent.ReportFilter, error) {
	if err := validateFilterTypeEntity(input); err != nil {
		return nil, err
	}
	filters, err := json.Marshal(input.Filters)
	if err != nil {
		return nil, err
	}
	rf, err := r.ClientFrom(ctx).
		ReportFilter.
		Create().
		SetName(input.Name).
		SetEntity(reportfilter.Entity(input.Entity)).
		SetFilters(string(filters)).
		Save(ctx)
	if err != nil && ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("a saved search with the name %s already exists", input.Name)
	}
	return rf, err
}

func (r mutationResolver) EditReportFilter(ctx context.Context, input models.EditReportFilterInput) (*ent.ReportFilter, error) {
	rf, err := r.ClientFrom(ctx).
		ReportFilter.
		UpdateOneID(input.ID).
		SetName(input.Name).
		Save(ctx)
	if err != nil && ent.IsConstraintError(err) {
		return nil, gqlerror.Errorf("a saved search with the name %s already exists", input.Name)
	}
	return rf, err
}

func (r mutationResolver) DeleteReportFilter(ctx context.Context, id int) (_ bool, err error) {
	if err = r.ClientFrom(ctx).ReportFilter.DeleteOneID(id).Exec(ctx); err != nil {
		err = fmt.Errorf("deleting report filter %q: %w", id, err)
	}
	return err == nil, err
}
