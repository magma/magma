/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

import (
	"errors"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/orc8r"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectHandlerFuncCalled         bool
	expectMigratedHandlerFuncCalled bool
	multiplexHandlers               bool

	handlerFuncError         string
	migratedHandlerFuncError string

	expectedError string
}

func TestRegister(t *testing.T) {
	oldRegistry := registries

	// Unmigrated, env var not set
	runCase(t, testCase{
		expectHandlerFuncCalled:         true,
		expectMigratedHandlerFuncCalled: false,
		multiplexHandlers:               true,
	})
	runCase(t, testCase{expectHandlerFuncCalled: true})
	runCase(t, testCase{
		expectHandlerFuncCalled: true,
		handlerFuncError:        "foo",
		expectedError:           "foo",
	})

	err := os.Setenv(orc8r.UseConfiguratorEnv, "1")
	assert.NoError(t, err)
	runCase(t, testCase{expectMigratedHandlerFuncCalled: true})
	runCase(t, testCase{
		expectMigratedHandlerFuncCalled: true,
		migratedHandlerFuncError:        "foo",
		expectedError:                   "foo",
	})

	runCase(t, testCase{
		expectMigratedHandlerFuncCalled: true,
		expectHandlerFuncCalled:         true,
		multiplexHandlers:               true,
		handlerFuncError:                "foo",
		expectedError:                   "foo",
	})

	runCase(t, testCase{
		expectHandlerFuncCalled:         false,
		expectMigratedHandlerFuncCalled: true,
		multiplexHandlers:               true,
		migratedHandlerFuncError:        "foo",
		expectedError:                   "foo",
	})

	registries = oldRegistry
}

func runCase(t *testing.T, tc testCase) {
	registries = map[HttpMethod]handlerRegistry{
		GET:    {},
		POST:   {},
		PUT:    {},
		DELETE: {},
	}
	handlerFuncCalled, migratedHandlerFuncCalled := false, false
	mockHandler := Handler{
		Methods: GET,
		Path:    "/foo/",
		HandlerFunc: func(c echo.Context) error {
			handlerFuncCalled = true
			if tc.handlerFuncError != "" {
				return errors.New(tc.handlerFuncError)
			}
			return nil
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			migratedHandlerFuncCalled = true
			if tc.migratedHandlerFuncError != "" {
				return errors.New(tc.migratedHandlerFuncError)
			}
			return nil
		},
		MultiplexAfterMigration: tc.multiplexHandlers,
	}
	err := Register(mockHandler)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/foo/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = registries[GET]["/foo/"](c)
	assert.Equal(t, 200, rec.Code)
	if tc.expectedError != "" {
		assert.EqualError(t, err, tc.expectedError)
	} else {
		assert.NoError(t, err)
	}
	assert.Equal(t, tc.expectHandlerFuncCalled, handlerFuncCalled)
	assert.Equal(t, tc.expectMigratedHandlerFuncCalled, migratedHandlerFuncCalled)
}
