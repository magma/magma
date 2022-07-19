/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package obsidian

const (
	Product                      = "Obsidian Server"
	Version                      = "0.1"
	DefaultPort                  = 9081
	DefaultHttpsPort             = 9443
	DefaultServerCert            = "server_cert.pem"
	DefaultServerCertKey         = "server_cert.key.pem"
	DefaultClientCAs             = "ca_cert.pem"
	DefaultStaticFolder          = "/var/opt/magma/static"
	StaticURLPrefix              = "/swagger"
	StaticURLPrefixLegacy        = "/apidocs"
	ServiceName                  = "OBSIDIAN"
	EnableDynamicSwaggerSpecsKey = "enable_dynamic_swagger_specs"
)

// configs
var (
	TLS                bool
	Port               int
	ServerCertPemPath  string
	ServerKeyPemPath   string
	ClientCAPoolPath   string
	AllowAnyClientCert bool
	StaticFolder       string
	// EnableDynamicSwaggerSpecs is a config in the obsidian
	// service config.
	// When true, poll and combine Swagger specs at runtime
	// When false, fall back to serving the static Swagger spec asset
	EnableDynamicSwaggerSpecs bool
)
