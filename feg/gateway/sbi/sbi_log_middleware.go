/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sbi

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggingTransport implements the http.RoundTripper interface and logs the
// HTTP request and response information
type LoggingTransport struct {
	Transport http.RoundTripper
	Logger    SbiLogger
}

type responseGrabber struct {
	io.Writer
	http.ResponseWriter
}

func NewLoggingTransport() LoggingTransport {
	return LoggingTransport{
		Transport: http.DefaultTransport,
		Logger:    SbiLogger{},
	}
}

func NewLoggingHttpClient() *http.Client {
	return &http.Client{
		Transport: NewLoggingTransport(),
	}
}

// RoundTrip overrides the http.DefaultTransport RoundTrip behavior and logs the HTTP request and response.
func (lt LoggingTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// Before sending
	var reqBody []byte

	if req.Body != nil {
		reqBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
		// reset the req.Body back to original
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	}
	lt.Logger.LogRequest(req.Method, req.URL, reqBody, req.Header)

	// Make the actual request
	start := time.Now()
	res, err = lt.Transport.RoundTrip(req)
	took := time.Since(start)
	if err != nil {
		// error responses shall be logged by caller
		return
	}

	// After sending
	var resBody []byte
	if res.Body != nil {
		resBody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		// reset the res.Body back to original
		res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))
	}
	lt.Logger.LogResponse(req.URL, res.Status, resBody, res.Header, took)
	return
}

func (w *responseGrabber) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseGrabber) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *responseGrabber) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseGrabber) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// ServerLoggingMiddleware logs HTTP request and response whenever a request is received.
// When logger is nil, the messages are logged to default glog verbose level 2 logger.
func ServerLoggingMiddleware() echo.MiddlewareFunc {
	logger := SbiLogger{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// Request
			var reqBody []byte
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
			}
			logger.LogRequest(c.Request().Method, c.Request().URL, reqBody, c.Request().Header)

			// Setting up for grabbing response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &responseGrabber{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			start := time.Now()
			// invoke request handler
			if err = next(c); err != nil {
				c.Error(err)
			}
			latency := time.Since(start)
			status := fmt.Sprintf("%d %s", c.Response().Status, http.StatusText(c.Response().Status))
			logger.LogResponse(c.Request().URL, status, resBody.Bytes(), c.Response().Header(), latency)

			return
		}
	}
}
