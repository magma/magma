// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hooks_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/hooks"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/suite"
)

type propertiesTestSuite struct {
	suite.Suite
	ctx       context.Context
	client    *ent.Client
	user      *ent.User
	typ       *ent.WorkOrderType
	template  *ent.WorkOrderTemplate
	workOrder *ent.WorkOrder
}

func (s *propertiesTestSuite) SetupSuite() {
	s.client = viewertest.NewTestClient(s.T())
	s.ctx = viewertest.NewContext(
		context.Background(),
		s.client,
	)
	u, ok := viewer.FromContext(s.ctx).(*viewer.UserViewer)
	s.Require().True(ok)
	s.user = u.User()
	var err error
	s.typ, err = s.client.WorkOrderType.
		Create().
		SetName("type").
		Save(s.ctx)
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) BeforeTest(_, _ string) {
	var err error
	s.template, err = s.client.WorkOrderTemplate.
		Create().
		SetName("template").
		Save(s.ctx)
	s.Require().NoError(err)
	s.workOrder, err = s.client.WorkOrder.
		Create().
		SetName("instance").
		SetType(s.typ).
		SetTemplate(s.template).
		SetCreationDate(time.Now()).
		SetOwner(s.user).
		Save(s.ctx)
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) createPropertyType(
	typ propertytype.Type, name string, change func(create *ent.PropertyTypeCreate)) *ent.PropertyType {
	create := s.client.PropertyType.Create().
		SetName(name).
		SetType(typ).
		SetIsInstanceProperty(true).
		SetMandatory(true).
		SetWorkOrderTemplate(s.template)
	if change != nil {
		change(create)
	}
	return create.SaveX(s.ctx)
}

func (s *propertiesTestSuite) closeWorkOrder(ctx context.Context, client *ent.Client) {
	_, err := client.WorkOrder.UpdateOne(s.workOrder).
		SetStatus(workorder.StatusDONE).
		Save(ctx)
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) withTransaction(f func(context.Context, *ent.Client)) error {
	tx, err := ent.FromContext(s.ctx).Tx(s.ctx)
	s.Require().NoError(err)
	ctx := ent.NewTxContext(s.ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	ctx = ent.NewContext(ctx, tx.Client())
	f(ctx, tx.Client())
	return tx.Commit()
}

func (s *propertiesTestSuite) TestStringPropertyExists() {
	pType := s.createPropertyType(propertytype.TypeString, "string", nil)
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetStringVal("value").
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) TestStringPropertyNotExists() {
	_ = s.createPropertyType(propertytype.TypeString, "string", nil)
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestEmptyStringPropertyExists() {
	pType := s.createPropertyType(propertytype.TypeString, "string", nil)
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestStringPropertyExistsWithDefault() {
	pType := s.createPropertyType(propertytype.TypeString, "string", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("default")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetStringVal("value").
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) TestStringPropertyNotExistsWithDefault() {
	_ = s.createPropertyType(propertytype.TypeString, "string", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("default")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		s.closeWorkOrder(ctx, client)
	})
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) TestEmptyStringPropertyExistsWithDefault() {
	pType := s.createPropertyType(propertytype.TypeString, "string", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("default")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestEnumPropertyExists() {
	pType := s.createPropertyType(propertytype.TypeEnum, "enum", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("[\"option1\",\"option2\"]")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetStringVal("option1").
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) TestEnumPropertyNotExists() {
	_ = s.createPropertyType(propertytype.TypeEnum, "enum", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("[\"option1\",\"option2\"]")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestEmptyEnumPropertyExists() {
	pType := s.createPropertyType(propertytype.TypeEnum, "enum", func(create *ent.PropertyTypeCreate) {
		create.SetStringVal("[\"option1\",\"option2\"]")
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestNodePropertyExists() {
	pType := s.createPropertyType(propertytype.TypeNode, "node", func(create *ent.PropertyTypeCreate) {
		create.SetNodeType(hooks.NodeTypeLocation)
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		locationType := client.LocationType.Create().
			SetName("locationType").
			SaveX(ctx)
		location := client.Location.Create().
			SetName("Location").
			SetType(locationType).
			SaveX(ctx)
		_, err := client.Property.Create().
			SetType(pType).
			SetWorkOrder(s.workOrder).
			SetLocationValue(location).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().NoError(err)
}

func (s *propertiesTestSuite) TestNodePropertyNotExists() {
	_ = s.createPropertyType(propertytype.TypeNode, "node", func(create *ent.PropertyTypeCreate) {
		create.SetNodeType(hooks.NodeTypeLocation)
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func (s *propertiesTestSuite) TestEmptyNodePropertyExists() {
	pType := s.createPropertyType(propertytype.TypeNode, "node", func(create *ent.PropertyTypeCreate) {
		create.SetNodeType(hooks.NodeTypeLocation)
	})
	err := s.withTransaction(func(ctx context.Context, client *ent.Client) {
		_, err := client.Property.Create().
			SetType(pType).
			SetWorkOrder(s.workOrder).
			Save(ctx)
		s.Require().NoError(err)
		s.closeWorkOrder(ctx, client)
	})
	s.Require().Error(err)
}

func TestPropertiesHooks(t *testing.T) {
	suite.Run(t, &propertiesTestSuite{})
}
