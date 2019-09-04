/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package tests

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Method  string
	URL     string
	Payload encoding.BinaryMarshaler
	Handler echo.HandlerFunc

	ParamNames  []string
	ParamValues []string

	ExpectedStatus int
	ExpectedResult encoding.BinaryMarshaler

	ExpectedError string
}

// RunUnitTest runs a test case using the given Echo instance. This function
// does not start an obsidian server..
func RunUnitTest(t *testing.T, e *echo.Echo, test Test) {
	var req *http.Request
	if test.Payload != nil {
		payloadBytes, err := test.Payload.MarshalBinary()
		if !assert.NoError(t, err) {
			return
		}
		req = httptest.NewRequest(test.Method, test.URL, bytes.NewReader(payloadBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(test.Method, test.URL, bytes.NewReader([]byte{}))
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames(test.ParamNames...)
	c.SetParamValues(test.ParamValues...)
	err := test.Handler(c)
	if err != nil {
		// in prod CollectStats middleware does this
		c.Error(err)
	}
	assert.Equal(t, test.ExpectedStatus, rec.Code)

	if test.ExpectedError != "" {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			assert.Equal(t, test.ExpectedError, httpErr.Message)
		} else {
			assert.EqualError(t, err, test.ExpectedError)
		}
	} else if assert.NoError(t, err) {
		if test.ExpectedResult != nil {
			expectedBytes, err := test.ExpectedResult.MarshalBinary()
			if assert.NoError(t, err) {
				// convert to string for more readable assert failure messages
				assert.Equal(t, string(expectedBytes), string(rec.Body.Bytes()))
			}
		}
	}
}

// GetHandlerByPathAndMethod fetches the first obsidian.Handler that matches the
// given path and method from a list of handlers. If no such handler exists, it
// will fail.
func GetHandlerByPathAndMethod(t *testing.T, handlers []obsidian.Handler, path string, method obsidian.HttpMethod) obsidian.Handler {
	for _, handler := range handlers {
		if handler.Path == path && handler.Methods == method {
			return handler
		}
	}
	assert.Fail(t, fmt.Sprintf("No handler registered for path %s", path))
	return obsidian.Handler{}
}

func JSONMarshaler(v interface{}) encoding.BinaryMarshaler {
	return &jsonMarshaler{v: v}
}

type jsonMarshaler struct {
	v interface{}
}

func (j *jsonMarshaler) MarshalBinary() (data []byte, err error) {
	return json.Marshal(j.v)
}
