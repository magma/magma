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

package header

// Request fields
const (
	Accept             = "Accept"
	AcceptCharset      = "Accept-Charset"
	AcceptEncoding     = "Accept-Encoding"
	AcceptLanguage     = "Accept-Language"
	Authorization      = "Authorization"
	Cookie             = "Cookie"
	Expect             = "Expect"
	Forwarded          = "Forwarded"
	From               = "From"
	Host               = "Host"
	IfMatch            = "If-Match"
	IfModifiedSince    = "If-Modified-Since"
	IfNoneMatch        = "If-None-Match"
	IfRange            = "If-Range"
	IfUnmodifiedSince  = "If-Unmodified-Since"
	MaxForwards        = "Max-Forwards"
	Origin             = "Origin"
	ProxyAuthorization = "Proxy-Authorization"
	Range              = "Range"
	Referer            = "Referer"
	TE                 = "TE"
	UserAgent          = "User-Agent"
)

// WebSocket fields
const (
	SecWebSocketKey = "Sec-WebSocket-Key"
)

// Common non-standard request fields
const (
	XRequestedWith      = "X-Requested-With"
	DNT                 = "DNT"
	XForwardedFor       = "X-Forwarded-For"
	XForwardedHost      = "X-Forwarded-Host"
	XForwardedProto     = "X-Forwarded-Proto"
	XRealIP             = "X-Real-IP"
	XHttpMethodOverride = "X-Http-Method-Override"
	// nolint:gas
	XCsrfToken = "X-Csrf-Token"
	XCSRFToken = "X-CSRFToken"
	XXSRFToken = "X-XSRF-TOKEN"
)

// Request and Response fields
const (
	CacheControl  = "Cache-Control"
	Connection    = "Connection"
	ContentLength = "Content-Length"
	ContentType   = "Content-Type"
	Date          = "Date"
	Pragma        = "Pragma"
	Upgrade       = "Upgrade"
	Via           = "Via"
	Warning       = "Warning"
)

// Response fields
const (
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AcceptPatch                   = "Accept-Patch"
	AcceptRanges                  = "Accept-Ranges"
	Age                           = "Age"
	Allow                         = "Allow"
	AltSvc                        = "Alt-Svc"
	ContentDisposition            = "Content-Disposition"
	ContentEncoding               = "Content-Encoding"
	ContentLanguage               = "Content-Language"
	ContentLocation               = "Content-Location"
	ContentRange                  = "Content-Range"
	ETag                          = "ETag"
	Expires                       = "Expires"
	LastModified                  = "Last-Modified"
	Link                          = "Link"
	Location                      = "Location"
	P3P                           = "P3P"
	ProxyAuthenticate             = "Proxy-Authenticate"
	PublicKeyPins                 = "Public-Key-Pins"
	ReferrerPolicy                = "Referrer-Policy"
	Refresh                       = "Refresh"
	RetryAfter                    = "Retry-After"
	Server                        = "Server"
	SetCookie                     = "Set-Cookie"
	StrictTransportSecurity       = "Strict-Transport-Security"
	Trailer                       = "Trailer"
	TransferEncoding              = "Transfer-Encoding"
	TSV                           = "TSV"
	Vary                          = "Vary"
	WWWAuthenticate               = "WWW-Authenticate"
	XFrameOptions                 = "X-Frame-Options"
)

// Common non-standard response fields
const (
	XXSSProtection          = "X-XSS-Protection"
	ContentSecurityPolicy   = "Content-Security-Policy"
	XContentSecurityPolicy  = "X-Content-Security-Policy"
	XWebKitCSP              = "X-WebKit-CSP"
	XContentTypeOptions     = "X-Content-Type-Options"
	XPoweredBy              = "X-Powered-By"
	XUACompatible           = "X-UA-Compatible"
	XContentDuration        = "X-Content-Duration"
	UpgradeInsecureRequests = "Upgrade-Insecure-Requests"
)

// Common non-standard request and response fields
const (
	XRequestID     = "X-Request-ID"
	XCorrelationID = "X-Correlation-ID"
)

// Encoding
const (
	EncodingChunked  = "chunked"
	EncodingCompress = "compress"
	EncodingDeflate  = "deflate"
	EncodingGzip     = "gzip"
	EncodingIdentity = "identity"
)

// Connection
const (
	ConnectionKeepAlive = "keep-alive"
	ConnectionUpgrade   = "upgrade"
)
