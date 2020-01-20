// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func getInput(name, email string) models.TechnicianInput {
	return models.TechnicianInput{Name: name, Email: email}
}

func TestAddTechnician(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	inp := getInput("name_1", "email_1@mail.com")
	tech1, err := mr.AddTechnician(ctx, inp)
	require.NoError(t, err)

	inp = getInput("name_2", "email_2@wow.com")
	tech2, err := mr.AddTechnician(ctx, inp)
	require.NoError(t, err)

	client := ent.FromContext(ctx)
	techs := client.Technician.Query().AllX(ctx)
	require.Len(t, techs, 2)
	require.NotEqual(t, tech1.Name, tech2.Name)
}

func TestAddTechniciansSameName(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	inp := getInput("name_1", "email_1@mail.com")
	_, err = mr.AddTechnician(ctx, inp)
	require.NoError(t, err)

	inp = getInput("name_2", "email_1@mail.com")
	_, err = mr.AddTechnician(ctx, inp)
	require.Error(t, err, "same email")

	inp = getInput("name_3", "@invalid.com")
	_, err = mr.AddTechnician(ctx, inp)
	require.Error(t, err, "invalid email")

	inp = getInput("name_3", "invalid.com")
	_, err = mr.AddTechnician(ctx, inp)
	require.Error(t, err, "invalid email")

	client := ent.FromContext(ctx)
	techs := client.Technician.Query().AllX(ctx)
	require.Len(t, techs, 1)
}
