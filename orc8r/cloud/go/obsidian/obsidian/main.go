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

package main

import (
	"flag"
	"log"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/server"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/definitions"
)

func main() {
	flag.IntVar(&obsidian.Port, "port", -1, "HTTP (REST) Server Port")
	flag.IntVar(&obsidian.Port, "p", -1, "HTTP (REST) Server Port (shorthand)")

	// HTTPS settings
	flag.BoolVar(&obsidian.TLS, "tls", false, "HTTPS only access")
	flag.StringVar(
		&obsidian.ServerCertPemPath, "cert",
		definitions.GetEnvWithDefault("REST_CERT", obsidian.DefaultServerCert),
		"Server's certificate PEM file",
	)
	flag.StringVar(
		&obsidian.ServerKeyPemPath, "cert_key",
		definitions.GetEnvWithDefault("REST_CERT_KEY", obsidian.DefaultServerCertKey),
		"Server's certificate private key PEM file",
	)
	flag.StringVar(
		&obsidian.ClientCAPoolPath, "client_ca",
		definitions.GetEnvWithDefault("REST_CLIENT_CERT", obsidian.DefaultClientCAs),
		"Client certificate CA pool PEM file",
	)
	flag.BoolVar(
		&obsidian.AllowAnyClientCert, "client_cert_any", false,
		"Accept Any Client Certificate (Do not verify with given client CAs)",
	)
	flag.StringVar(
		&obsidian.StaticFolder, "static_folder", obsidian.DefaultStaticFolder,
		"Folder containing the static files served",
	)

	srv, err := service.NewOrchestratorService(orc8r.ModuleName, obsidian.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	obsidian.EnableDynamicSwaggerSpecs = srv.Config.MustGetBool(obsidian.EnableDynamicSwaggerSpecsKey)

	if obsidian.Port == -1 {
		obsidian.Port = obsidian.DefaultPort
		if obsidian.TLS {
			obsidian.Port = obsidian.DefaultHttpsPort
		}
	}

	go srv.Run()
	server.Start()
}
