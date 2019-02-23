/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

// RequestOperator relies on x-magma-client-cert-serial HTTP request header,
// the header string is redefined here to avoid sharing it with magma GRPC
// Identity middleware & to comply with specific to Go's net/http header
// capitalization: https://golang.org/pkg/net/http/#Request
const (
	// Client Certificate CN Header (for logging only)
	CLIENT_CERT_CN_KEY = "X-Magma-Client-Cert-Cn"
	// Client Certificate Serial Number Header
	CLIENT_CERT_SN_KEY = "X-Magma-Client-Cert-Serial"
)
