// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/cenkalti/backoff/v4"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type client struct {
	client     *graphql.Client
	log        *zap.Logger
	tenant     string
	user       string
	automation bool
}

func TestMain(m *testing.M) {
	if err := waitFor("graph"); err != nil {
		fmt.Printf("FAIL\n%v\n", err)
		os.Exit(2)
	}
	os.Exit(m.Run())
}

func waitFor(services ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	g := ctxgroup.WithContext(ctx)
	for _, service := range services {
		service := service
		target := fmt.Sprintf("http://%s/healthz", service)
		g.Go(func(ctx context.Context) error {
			return backoff.Retry(func() error {
				rsp, err := ctxhttp.Get(ctx, nil, target)
				if err != nil {
					return err
				}
				rsp.Body.Close()
				if rsp.StatusCode != http.StatusOK {
					return fmt.Errorf("service %q not ready: status=%q", service, rsp.Status)
				}
				return nil
			}, backoff.WithContext(
				backoff.NewConstantBackOff(200*time.Millisecond), ctx),
			)
		})
	}
	return g.Wait()
}

type option func(*client)

func withAutomation() option {
	return func(c *client) {
		c.automation = true
	}
}

func newClient(t *testing.T, tenant, user string, opts ...option) *client {
	c := client{
		log: zaptest.NewLogger(t).With(
			zap.String("tenant", tenant),
		),
		tenant: tenant,
		user:   user,
	}
	c.client = graphql.NewClient(
		"http://graph/query",
		&http.Client{Transport: &c},
	)
	for _, opt := range opts {
		opt(&c)
	}
	require.NoError(t, c.createTenant())
	require.NoError(t, c.createOwnerUser())
	return &c
}

func (c *client) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-auth-organization", c.tenant)
	if !c.automation {
		req.Header.Set("x-auth-user-email", c.user)
	} else {
		req.Header.Set("x-auth-automation-name", c.user)
	}
	req.Header.Set("x-auth-user-role", "OWNER")
	return http.DefaultTransport.RoundTrip(req)
}

func (c *client) createTenant() error {
	conn, err := grpc.Dial("graph:443", grpc.WithInsecure())
	if err != nil {
		return err
	}
	_, err = graphgrpc.NewTenantServiceClient(conn).
		Create(context.Background(), &wrappers.StringValue{Value: c.tenant})
	switch st, _ := status.FromError(err); st.Code() {
	case codes.OK, codes.AlreadyExists:
	default:
		return st.Err()
	}
	return nil
}

func (c *client) createOwnerUser() error {
	conn, err := grpc.Dial("graph:443", grpc.WithInsecure())
	if err != nil {
		return err
	}
	_, err = graphgrpc.NewUserServiceClient(conn).
		Create(context.Background(), &graphgrpc.AddUserInput{Tenant: c.tenant, Id: c.user, IsOwner: true})
	switch st, _ := status.FromError(err); st.Code() {
	case codes.OK:
	default:
		return st.Err()
	}
	return nil
}

type addLocationTypeResponse struct {
	ID   graphql.ID
	Name graphql.String
}

func (c *client) addLocationType(name string, properties ...*models.PropertyTypeInput) (*addLocationTypeResponse, error) {
	var m struct {
		Response addLocationTypeResponse `graphql:"addLocationType(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": models.AddLocationTypeInput{
			Name:       name,
			Properties: properties,
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

type addLocationResponse struct {
	ID   graphql.ID
	Name graphql.String
}

func IDToInt(id graphql.ID) int {
	i, err := strconv.Atoi(id.(string))
	if err != nil {
		panic(err)
	}
	return i
}

func IDToIntOrNil(id graphql.ID) *int {
	if id == nil {
		return nil
	}
	i := IDToInt(id)
	return &i
}

func (c *client) addLocation(name string, parent graphql.ID) (*addLocationResponse, error) {
	typ, err := c.addLocationType("location_type_" + uuid.New().String())
	if err != nil {
		return nil, err
	}
	var m struct {
		Response addLocationResponse `graphql:"addLocation(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": models.AddLocationInput{
			Name:      name,
			Type:      IDToInt(typ.ID),
			Latitude:  pointer.ToFloat64(14.45),
			Longitude: pointer.ToFloat64(45.14),
			Parent:    IDToIntOrNil(parent),
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

type queryLocationResponse struct {
	ID       graphql.ID
	Name     graphql.String
	Children []struct {
		ID   graphql.ID
		Name graphql.String
		Type struct {
			ID   graphql.ID
			Name graphql.String
		} `graphql:"locationType"`
	}
}

func (c *client) queryLocation(id graphql.ID) (*queryLocationResponse, error) {
	var q struct {
		Node struct {
			Response queryLocationResponse `graphql:"... on Location"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}
	switch err := c.client.Query(context.Background(), &q, vars); {
	case err != nil:
		return nil, err
	case q.Node.Response.ID == nil:
		return nil, errors.New("location not found")
	}
	return &q.Node.Response, nil
}

type queryLocationsResponse struct {
	Edges []struct {
		Node struct {
			Name   graphql.String
			Parent struct {
				Name graphql.String
			} `graphql:"parentLocation"`
			Properties []struct {
				Type struct {
					Name graphql.String
				} `graphql:"propertyType"`
				Value graphql.String `graphql:"stringValue"`
			}
			Children []struct {
				Name graphql.String
			}
			ExternalID graphql.String
			Longitude  graphql.Float
			Latitude   graphql.Float
		}
	}
}

func (c *client) QueryLocations() (*queryLocationsResponse, error) {
	var q struct {
		Response queryLocationsResponse `graphql:"locations(first: null)"`
	}
	if err := c.client.Query(context.Background(), &q, nil); err != nil {
		return nil, err
	}
	return &q.Response, nil
}

type addEquipmentTypeResponse struct {
	ID         graphql.ID
	Name       graphql.String
	Properties []struct {
		ID   graphql.ID
		Name graphql.String
		Kind models.PropertyKind `graphql:"type"`
	} `graphql:"propertyTypes"`
}

func (c *client) addEquipmentType(name string, properties ...*models.PropertyTypeInput) (*addEquipmentTypeResponse, error) {
	var m struct {
		Response addEquipmentTypeResponse `graphql:"addEquipmentType(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": models.AddEquipmentTypeInput{
			Name:       name,
			Properties: properties,
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

type addEquipmentResponse struct {
	ID graphql.ID
}

func (c *client) addEquipment(name string, typ, location, workOrder graphql.ID) (*addEquipmentResponse, error) {
	var m struct {
		Response addEquipmentResponse `graphql:"addEquipment(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": models.AddEquipmentInput{
			Name:      name,
			Type:      IDToInt(typ),
			Location:  IDToIntOrNil(location),
			WorkOrder: IDToIntOrNil(workOrder),
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

func (c *client) removeEquipment(id, workOrder graphql.ID) error {
	var m struct {
		ID graphql.ID `graphql:"removeEquipment(id: $id, workOrderId: $workOrder)"`
	}
	vars := map[string]interface{}{
		"id":        id,
		"workOrder": workOrder,
	}
	return c.client.Mutate(context.Background(), &m, vars)
}

type queryEquipmentResponse struct {
	ID    graphql.ID
	Name  graphql.String
	State models.FutureState `graphql:"futureState"`
}

func (c *client) queryEquipment(id graphql.ID) (*queryEquipmentResponse, error) {
	var q struct {
		Node struct {
			Response queryEquipmentResponse `graphql:"... on Equipment"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}
	switch err := c.client.Query(context.Background(), &q, vars); {
	case err != nil:
		return nil, err
	case q.Node.Response.ID == nil:
		return nil, errors.New("equipment not found")
	}
	return &q.Node.Response, nil
}

type addWorkOrderTypeResponse struct {
	ID         graphql.ID
	Name       graphql.String
	Properties []struct {
		ID   graphql.ID
		Name graphql.String
		Kind models.PropertyKind `graphql:"type"`
	} `graphql:"propertyTypes"`
}

func (c *client) addWorkOrderType(name string, properties ...*models.PropertyTypeInput) (*addWorkOrderTypeResponse, error) {
	var m struct {
		Response addWorkOrderTypeResponse `graphql:"addWorkOrderType(input: $input)"`
	}
	if properties == nil {
		properties = []*models.PropertyTypeInput{}
	}
	vars := map[string]interface{}{
		"input": models.AddWorkOrderTypeInput{
			Name:       name,
			Properties: properties,
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

type User struct {
	ID    graphql.ID
	Email graphql.String
}

type addWorkOrderResponse struct {
	ID    graphql.ID
	Name  graphql.String
	Owner User
}

func (c *client) addWorkOrder(name string, typ graphql.ID) (*addWorkOrderResponse, error) {
	var m struct {
		Response addWorkOrderResponse `graphql:"addWorkOrder(input: $input)"`
	}
	vars := map[string]interface{}{
		"input": models.AddWorkOrderInput{
			Name:            name,
			WorkOrderTypeID: IDToInt(typ),
		},
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return nil, err
	}
	return &m.Response, nil
}

func (c *client) executeWorkOrder(workOrder *addWorkOrderResponse) error {
	var em struct {
		Response struct {
			ID graphql.ID
		} `graphql:"editWorkOrder(input: $input)"`
	}
	ownerID := IDToInt(workOrder.Owner.ID)
	vars := map[string]interface{}{
		"input": models.EditWorkOrderInput{
			ID:       IDToInt(workOrder.ID),
			Name:     string(workOrder.Name),
			OwnerID:  &ownerID,
			Status:   models.WorkOrderStatusDone,
			Priority: models.WorkOrderPriorityNone,
		},
	}
	if err := c.client.Mutate(context.Background(), &em, vars); err != nil {
		return xerrors.Errorf("editing work order: %w", err)
	}

	var m struct {
		Response struct {
			ID graphql.ID
		} `graphql:"executeWorkOrder(id: $id)"`
	}
	vars = map[string]interface{}{
		"id": workOrder.ID,
	}
	if err := c.client.Mutate(context.Background(), &m, vars); err != nil {
		return xerrors.Errorf("executing work order: %w", err)
	}
	return nil
}

const (
	testTenant = "integration-test"
	testUser   = "user@test.com"
)

func TestAddLocation(t *testing.T) {
	c := newClient(t, testTenant, testUser)
	name := "location_" + uuid.New().String()
	rsp, err := c.addLocation(name, nil)
	require.NoError(t, err)
	assert.NotNil(t, rsp.ID)
	assert.EqualValues(t, name, rsp.Name)
}

func TestAddLocationType(t *testing.T) {
	c := newClient(t, testTenant, testUser)
	name := "location_type_" + uuid.New().String()
	typ, err := c.addLocationType(name)
	require.NoError(t, err)
	assert.NotNil(t, typ.ID)
	assert.EqualValues(t, name, typ.Name)
}

func TestAddLocationWithAutomation(t *testing.T) {
	c := newClient(t, testTenant, testUser, withAutomation())
	name := "location_type_" + uuid.New().String()
	typ, err := c.addLocationType(name)
	require.NoError(t, err)
	assert.NotNil(t, typ.ID)
	assert.EqualValues(t, name, typ.Name)
}

func TestAddLocationsDifferentTenants(t *testing.T) {
	c1 := newClient(t, "integration-test-1", "user@test-1.com")
	c2 := newClient(t, "integration-test-2", "user@test-2.com")

	// makes sure tenant-2 does not have access to tenant-1 locations.
	name := "location_" + uuid.New().String()
	rsp, err := c1.addLocation(name, nil)
	require.NoError(t, err)
	l1, err := c1.queryLocation(rsp.ID)
	require.NoError(t, err)
	assert.EqualValues(t, name, l1.Name)
	_, err = c2.queryLocation(rsp.ID)
	assert.Error(t, err)

	name = "location_" + uuid.New().String()
	_, err = c2.addLocation(name, nil)
	require.NoError(t, err)
	locations, err := c1.QueryLocations()
	require.NoError(t, err)
	// make sure tenant-2 location does not exist in tenant-1.
	for i := range locations.Edges {
		assert.NotEqual(t, name, string(locations.Edges[i].Node.Name))
	}
}

func TestAddLocationWithChildren(t *testing.T) {
	c := newClient(t, testTenant, testUser)

	parentName := "parent_location_" + uuid.New().String()
	parent, err := c.addLocation(parentName, nil)
	require.NoError(t, err)
	childName := "child_location_" + uuid.New().String()
	child, err := c.addLocation(childName, parent.ID)
	require.NoError(t, err)

	location, err := c.queryLocation(parent.ID)
	require.NoError(t, err)
	assert.EqualValues(t, parentName, location.Name)
	require.Len(t, location.Children, 1)
	assert.Equal(t, child.ID, location.Children[0].ID)
	assert.EqualValues(t, childName, location.Children[0].Name)
}

func TestExecuteWorkOrder(t *testing.T) {
	c := newClient(t, testTenant, testUser)

	typ, err := c.addWorkOrderType("work_order_type_" + uuid.New().String())
	require.NoError(t, err)
	name := "work_order_" + uuid.New().String()
	workorder, err := c.addWorkOrder(name, typ.ID)
	require.NoError(t, err)
	assert.EqualValues(t, testUser, workorder.Owner.Email)

	et, err := c.addEquipmentType("router_type_" + uuid.New().String())
	require.NoError(t, err)
	l, err := c.addLocation("location_"+uuid.New().String(), nil)
	require.NoError(t, err)
	e, err := c.addEquipment("router_"+uuid.New().String(), et.ID, l.ID, workorder.ID)
	require.NoError(t, err)

	eq, err := c.queryEquipment(e.ID)
	require.NoError(t, err)
	assert.Equal(t, models.FutureStateInstall, eq.State)

	err = c.executeWorkOrder(workorder)
	require.NoError(t, err)

	eq, err = c.queryEquipment(e.ID)
	require.NoError(t, err)
	assert.Empty(t, eq.State)

	workorder, err = c.addWorkOrder(name, typ.ID)
	require.NoError(t, err)
	err = c.removeEquipment(eq.ID, workorder.ID)
	require.NoError(t, err)

	eq, err = c.queryEquipment(e.ID)
	require.NoError(t, err)
	assert.EqualValues(t, models.FutureStateRemove, eq.State)
	err = c.executeWorkOrder(workorder)
	require.NoError(t, err)
}

func TestPossibleProperties(t *testing.T) {
	c := newClient(t, testTenant, testUser)

	_, err := c.addEquipmentType(
		"router_type_"+uuid.New().String(),
		&models.PropertyTypeInput{
			Name: "Width",
			Type: models.PropertyKindInt,
		},
		&models.PropertyTypeInput{
			Name: "Manufacturer",
			Type: models.PropertyKindString,
		},
	)
	require.NoError(t, err)

	_, err = c.addEquipmentType(
		"router_type_"+uuid.New().String(),
		&models.PropertyTypeInput{
			Name: "Width",
			Type: models.PropertyKindInt,
		},
	)
	require.NoError(t, err)

	var q struct {
		Properties []struct {
			ID graphql.ID
		} `graphql:"possibleProperties(entityType: $entityType)"`
	}

	vars := map[string]interface{}{
		"entityType": models.PropertyEntityEquipment,
	}
	err = c.client.Query(context.Background(), &q, vars)
	require.NoError(t, err)
	assert.Len(t, q.Properties, 2)
}

func TestViewer(t *testing.T) {
	c := newClient(t, testTenant, testUser)
	var q struct {
		Viewer struct {
			Tenant graphql.String
			Email  graphql.String
		} `graphql:"me"`
	}
	err := c.client.Query(context.Background(), &q, nil)
	require.NoError(t, err)
	assert.EqualValues(t, testTenant, q.Viewer.Tenant)
	assert.EqualValues(t, testUser, q.Viewer.Email)
}
