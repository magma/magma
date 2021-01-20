package download_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/download/servicers"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
)

func TestDownloadServiceClientMethods(t *testing.T) {

	// Start the server for testing
	portNumber := initTestServer()

	// Give it time to start
	time.Sleep(5 * time.Second)

	// Make an HTTP/2 request of the server
	client := NewH2CClient()

	// Test non-s3 path
	url := fmt.Sprintf("http://localhost:%s/TestMeNow", portNumber)
	resp, err := client.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 400)

	defer resp.Body.Close()

	// Read the response body
	_, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Test unknown application
	url = fmt.Sprintf("http://localhost:%s/s3/blah/TestMeNow", portNumber)
	resp, err = client.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

// Implementation of a non-TLS HTTP/2 client copied from D7784190
func NewH2CClient() *http.Client {
	return &http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}
}

func initTestServer() string {
	// start the listener on port 0 to use one of available ports
	tcpListener, err := net.Listen("tcp", ":0")
	if err != nil {
		glog.Fatalf("net.Listen err: %v", err)
	}
	listenerAddr := tcpListener.Addr()
	glog.V(2).Infof("Starting listener on address %s", listenerAddr)

	go func() {
		config := servicers.InitServiceConfig()
		server := http2.Server{}
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				glog.Fatalf("tcpListener.Accept err: %v", err)
			}
			server.ServeConn(conn, &http2.ServeConnOpts{
				Handler: servicers.RootHandler(config),
			})
		}
	}()
	return getPortNumber(listenerAddr)
}

func getPortNumber(listenerAddr net.Addr) string {
	splitAddr := strings.Split(listenerAddr.String(), ":")
	return splitAddr[len(splitAddr)-1]
}
