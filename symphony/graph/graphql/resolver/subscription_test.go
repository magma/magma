// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"strings"
	"sync"
	"testing"

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
			`mutation { addWorkOrderType(input: { name: "chore" }) { id } }`,
			&rsp,
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
		var rsp struct{ AddWorkOrder struct{ ID string } }
		c.MustPost(
			`mutation($type: ID!) { addWorkOrder(input: { name: "clean", workOrderTypeId: $type }) { id } }`,
			&rsp,
			client.Var("type", typ),
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
		c.MustPost(`mutation($id: ID!) { editWorkOrder(input: { id: $id, name: "foo", status: DONE }) { id } }`,
			&rsp,
			client.Var("id", id),
		)
		id = rsp.EditWorkOrder.ID
	}
	wg.Wait()
	assert.Equal(t, id, sid)
}
