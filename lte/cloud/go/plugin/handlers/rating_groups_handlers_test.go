/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/plugin/handlers"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// Basic API workflow tests
func TestRatingGroupHandlersBasic(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.LteNetworkType})
	assert.NoError(t, err)

	listRatingGroups := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/rating_groups", obsidian.GET).HandlerFunc
	createRatingGroup := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/rating_groups", obsidian.POST).HandlerFunc
	getRatingGroup := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/rating_groups/:rating_group_id", obsidian.GET).HandlerFunc
	updateRatingGroup := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/rating_groups/:rating_group_id", obsidian.PUT).HandlerFunc
	deleteRatingGroup := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/rating_groups/:rating_group_id", obsidian.DELETE).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/rating_groups",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listRatingGroups,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.RatingGroup{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/rating_groups"

	// Test add rating group
	testRatingGroup := &models.RatingGroup{
		ID:        models.RatingGroupID(uint32(1)),
		LimitType: swag.String("FINITE"),
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/rating_groups",
		Payload:        testRatingGroup,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createRatingGroup,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that rating group was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/rating_groups",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listRatingGroups,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.RatingGroup{
			"1": testRatingGroup,
		}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/rating_groups"

	// Test Read Rating Group Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/rating_groups/1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rating_group_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        getRatingGroup,
		ExpectedStatus: 200,
		ExpectedResult: testRatingGroup,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update rating group
	testRatingGroup.LimitType = swag.String("INFINITE_METERED")
	testMutableGroup := &models.MutableRatingGroup{
		LimitType: swag.String("INFINITE_METERED"),
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/rating_groups/1",
		Payload:        testMutableGroup,
		ParamNames:     []string{"network_id", "rating_group_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        updateRatingGroup,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update results
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/rating_groups/1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rating_group_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        getRatingGroup,
		ExpectedStatus: 200,
		ExpectedResult: testRatingGroup,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete a rating group
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/rating_groups/1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rating_group_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        deleteRatingGroup,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/rating_groups",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listRatingGroups,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.RatingGroup{}),
	}
	tests.RunUnitTest(t, e, tc)
}
