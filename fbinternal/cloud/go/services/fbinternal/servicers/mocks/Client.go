/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
)

type HTTPClient struct {
	mock.Mock
}

func (client *HTTPClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := client.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}
