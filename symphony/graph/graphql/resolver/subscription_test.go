// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"sync"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionWorkOrder(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()
	c := newGraphClient(t, resolver)

	var typ string
	{
		var rsp struct{ AddWorkOrderType struct{ ID string } }
		c.MustPost(
			`mutation($input: AddWorkOrderTypeInput!) { addWorkOrderType(input: $input) { id } }`,
			&rsp,
			client.Var("input", models.AddWorkOrderTypeInput{Name: "chore"}),
		)
		typ = rsp.AddWorkOrderType.ID
	}

	var (
		sub = c.Websocket(`subscription { workOrderAdded { id name workOrderType { name } } }`)
		wg  sync.WaitGroup
		sid string
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var rsp struct {
			WorkOrderAdded struct {
				ID            string
				Name          string
				WorkOrderType struct {
					Name string
				}
			}
		}
		err := sub.Next(&rsp)
		require.NoError(t, err)
		sid = rsp.WorkOrderAdded.ID
		require.NotEmpty(t, sid)
		assert.Equal(t, "clean", rsp.WorkOrderAdded.Name)
		assert.Equal(t, "chore", rsp.WorkOrderAdded.WorkOrderType.Name)
		err = sub.Close()
		assert.NoError(t, err)
	}()

	var id string
	{
		var rsp struct{ AddWorkOrder struct{ ID string } }
		c.MustPost(
			`mutation($input: AddWorkOrderInput!) { addWorkOrder(input: $input) { id } }`,
			&rsp,
			client.Var("input", models.AddWorkOrderInput{Name: "clean", WorkOrderTypeID: typ}),
		)
		id = rsp.AddWorkOrder.ID
	}
	wg.Wait()
	assert.Equal(t, id, sid)

	sub = c.Websocket(`subscription { workOrderDone { id } }`)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var rsp struct{ WorkOrderDone struct{ ID string } }
		err := sub.Next(&rsp)
		require.NoError(t, err)
		sid = rsp.WorkOrderDone.ID
		require.NotEmpty(t, sid)
		err = sub.Close()
		assert.NoError(t, err)
	}()

	{
		var rsp struct{ EditWorkOrder struct{ ID string } }
		c.MustPost(`mutation($input: EditWorkOrderInput!) { editWorkOrder(input: $input) { id } }`,
			&rsp,
			client.Var("input", models.EditWorkOrderInput{
				ID:       id,
				Name:     "foo",
				Status:   models.WorkOrderStatusDone,
				Priority: models.WorkOrderPriorityNone,
			}),
		)
		id = rsp.EditWorkOrder.ID
	}
	wg.Wait()
	assert.Equal(t, id, sid)
}
