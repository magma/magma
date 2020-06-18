// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func websocket(client *client.Client, query string) *client.Subscription {
	sub := client.Websocket(query)
	next := sub.Next
	sub.Next = func(rsp interface{}) error {
		for {
			if err := next(rsp); err == nil ||
				!strings.HasPrefix(err.Error(), "expected data message, got") ||
				!strings.Contains(err.Error(), "ka") {
				return err
			}
		}
	}
	return sub
}

func TestSubscriptionWorkOrder(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.Close()
	c := resolver.GraphClient()

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
		sub = websocket(c, `subscription { workOrderAdded { id name workOrderType { name } } }`)
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
		input := models.AddWorkOrderInput{Name: "clean"}
		input.WorkOrderTypeID, _ = strconv.Atoi(typ)
		var rsp struct{ AddWorkOrder struct{ ID string } }
		c.MustPost(
			`mutation($input: AddWorkOrderInput!) { addWorkOrder(input: $input) { id } }`,
			&rsp,
			client.Var("input", input),
		)
		id = rsp.AddWorkOrder.ID
	}
	wg.Wait()
	assert.Equal(t, id, sid)

	sub = websocket(c, `subscription { workOrderDone { id } }`)
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
		input := models.EditWorkOrderInput{
			Name:     "foo",
			Status:   models.WorkOrderStatusDone,
			Priority: models.WorkOrderPriorityNone,
		}
		input.ID, _ = strconv.Atoi(id)
		c.MustPost(`mutation($input: EditWorkOrderInput!) { editWorkOrder(input: $input) { id } }`,
			&rsp,
			client.Var("input", input),
		)
		id = rsp.EditWorkOrder.ID
	}
	wg.Wait()
	assert.Equal(t, id, sid)
}
