/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package http2 contains a minimal implementation of non-TLS http/2 server
// and client
package http2

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// HTTPGrpcError is an error wraps a error message, a grpc status code and
// a http status code. It implements Error interface.
type HTTPGrpcError struct {
	Msg            string
	GrpcStatusCode int
	HttpStatusCode int
}

// NewHTTPGrpcError creates a new HTTPGrpcError.
func NewHTTPGrpcError(
	msg string,
	grpcStatusCode int,
	httpStatusCode int,
) *HTTPGrpcError {
	return &HTTPGrpcError{Msg: msg, GrpcStatusCode: grpcStatusCode, HttpStatusCode: httpStatusCode}
}

// Error returns HTTPGrpcError as a string
func (err *HTTPGrpcError) Error() string {
	return fmt.Sprintf("HTTPGrpcError %v: httpStatus: %d, grpcStatus: %d\n",
		err.Msg, err.HttpStatusCode, err.GrpcStatusCode)
}

// WriteErrResponse writes the HTTPGrpcError err into the http.ResponseWriter
// w. This is used by http2 server, and the server's handler can return
// after calling this function.
func WriteErrResponse(w http.ResponseWriter, err *HTTPGrpcError) {
	httpCode := err.HttpStatusCode
	// http status code has to exist in the statusText map.
	if statusText := http.StatusText(err.HttpStatusCode); len(statusText) == 0 {
		glog.Errorf("Received unrecognized httpStatusCode: %v\n", err.HttpStatusCode)
		httpCode = 400
	}
	w.Header().Set("content-type", "application/grpc")
	w.Header().Set("trailer", "Grpc-Status")
	w.Header().Add("trailer", "Grpc-Message")
	w.Header().Set("grpc-status", strconv.Itoa(err.GrpcStatusCode))
	w.Header().Set("grpc-message", PercentEncode(err.Msg))
	w.WriteHeader(httpCode)
}

// grpc-message has to be percent-encoded, though space can be kept.
// If the grpc-message is not properly encoded, the whole message will show up
// as an empty string in the application layer on the client side.
// I.e., without encoding, errMsg "Failed gwId: 123" will be ""
// because ':' has to be replaced with '%3A'
// read here for more info on how grpc encodes on http/2:
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
// read here on percent-encoding:
// https://developer.mozilla.org/en-US/docs/Glossary/percent-encoding
func PercentEncode(str string) string {
	// percent-encode the message, then replace + back to " "
	// for better readability.
	urlEncoded := url.QueryEscape(str)
	return strings.Replace(urlEncoded, "+", " ", -1)
}
