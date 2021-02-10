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
	Product              = "Obsidian Server"
	Version              = "0.1"
	DefaultPort          = 9081
	DefaultHttpsPort     = 9443
	DefaultServerCert    = "server_cert.pem"
	DefaultServerCertKey = "server_cert.key.pem"
	DefaultClientCAs     = "ca_cert.pem"
	DefaultStaticFolder  = "/var/opt/magma/static"
	StaticURLPrefix      = "/apidocs"
	ServiceName          = "OBSIDIAN"
	// EnableRunTimeSpecs is a parameter name in the obsidian service config.
	// When true, poll and merge Swagger specs at runtime
	// When false, do not poll or merge Swagger specs at runtime
	EnableRunTimeSpecs = "enable_runtime_specs"
)

// configs
var (
	TLS                  bool
	Port                 int
	ServerCertPemPath    string
	ServerKeyPemPath     string
	ClientCAPoolPath     string
	AllowAnyClientCert   bool
	StaticFolder         string
	CombineSpecAtRuntime bool
)
