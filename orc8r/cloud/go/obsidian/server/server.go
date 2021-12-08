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

// Server's main package, run with obsidian -h to see all available options
package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/cloud/go/obsidian/reverse_proxy"
	"magma/orc8r/cloud/go/obsidian/swagger/handlers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/service/config"
)

const (
	reverseProxyRefreshPeriod = 1 * time.Minute
)

func Start() {
	e := echo.New()
	e.HideBanner = true

	obsidian.AttachAll(e)
	// Metrics middleware is used before all other middlewares
	e.Use(CollectStats)
	e.Use(middleware.Recover())
	err := handlers.RegisterSwaggerHandlers(e)

	if err != nil {
		// Swallow RegisterHandlerError because the obsidian service should
		// continue to run even if Swagger handlers aren't registered.
		glog.Errorf("Error registering Swagger handlers %+v", err)
	}

	// Serve static assets for the Swagger UI
	e.Static(obsidian.StaticURLPrefix+"/static/swagger-ui/dist", obsidian.StaticFolder+"/swagger-ui/dist")
	e.Static(obsidian.StaticURLPrefix+"/v1/static", obsidian.StaticFolder+"/swagger/v1/static")

	portStr := fmt.Sprintf(":%d", obsidian.Port)
	log.Printf("Starting %s on %s", obsidian.Product, portStr)

	if obsidian.TLS {
		var caCerts []byte
		caCerts, err = ioutil.ReadFile(obsidian.ClientCAPoolPath)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCerts)
		if ok {
			log.Printf("Loaded %d Client CA Certificate[s] from '%s'", len(caCertPool.Subjects()), obsidian.ClientCAPoolPath)
		} else {
			log.Printf(
				"ERROR: No Certificates found in '%s'", obsidian.ClientCAPoolPath)
		}
		// Possible clientCertVerification values:
		// 	NoClientCert
		// 	RequestClientCert
		// 	RequireAnyClientCert
		// 	VerifyClientCertIfGiven
		// 	RequireAndVerifyClientCert
		clientCertVerification := tls.RequireAndVerifyClientCert
		if obsidian.AllowAnyClientCert {
			clientCertVerification = tls.RequireAnyClientCert
		}
		s := e.TLSServer
		s.TLSConfig = &tls.Config{
			Certificates: make([]tls.Certificate, 1),
			ClientCAs:    caCertPool,
			ClientAuth:   clientCertVerification,
			// Limit versions & Ciphers to our preferred list
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // 4 HTTP2 support
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}
		s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(obsidian.ServerCertPemPath, obsidian.ServerKeyPemPath)
		if err != nil {
			log.Fatalf(
				"ERROR loading server certificate ('%s') and/or key ('%s'): %s",
				obsidian.ServerCertPemPath, obsidian.ServerKeyPemPath, err,
			)
		}
		s.TLSConfig.BuildNameToCertificate()
		s.Addr = portStr
		if !e.DisableHTTP2 {
			s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
		}
	} else {
		e.Use(access.CertificateMiddleware)
	}
	var serviceConfig certifier.Config
	_, _, err = config.GetStructuredServiceConfig(orc8r.ModuleName, certifier.ServiceName, &serviceConfig)
	if err != nil {
		glog.Infof("Failed unmarshalling service config %v", err)
	}
	reverseProxyHandler := reverse_proxy.NewReverseProxyHandler(&serviceConfig)
	pathPrefixesByAddr, err := reverse_proxy.GetEchoServerAddressToPathPrefixes()
	if err != nil {
		log.Fatalf("Error querying service registry for reverse proxy paths: %s", err)
	}
	e, err = reverseProxyHandler.AddReverseProxyPaths(e, pathPrefixesByAddr)
	if err != nil {
		log.Fatalf("Error adding reverse proxy paths: %s", err)
	}
	go func() {
		for {
			<-time.After(reverseProxyRefreshPeriod)
			pathPrefixesByAddr, err := reverse_proxy.GetEchoServerAddressToPathPrefixes()
			if err != nil {
				log.Printf("Error querying service registry for reverse proxy paths: %s", err)
				continue
			}
			e, err = reverseProxyHandler.AddReverseProxyPaths(e, pathPrefixesByAddr)
			if err != nil {
				log.Printf("An error occurred while updating reverse proxy paths: %s", err)
			}
		}
	}()

	if obsidian.TLS {
		err = e.StartServer(e.TLSServer)
	} else {
		err = e.Start(portStr)
	}
	if err != nil {
		log.Println(err)
	}
}
