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

package scale

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"magma/orc8r/cloud/api/v1/go/client"
	"magma/orc8r/cloud/api/v1/go/client/a_p_ns"
	"magma/orc8r/cloud/api/v1/go/client/lte_gateways"
	"magma/orc8r/cloud/api/v1/go/client/lte_networks"
	"magma/orc8r/cloud/api/v1/go/client/subscribers"
	"magma/orc8r/cloud/api/v1/go/client/upgrades"
	"magma/orc8r/cloud/api/v1/go/models"

	"magma/orc8r/cloud/test/testlib"

	oclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/api"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/pkcs12"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// TODO(hcgatewood): clean up this scale test then add it to daily CI job

func init() {
	_ = flag.Set("logtostderr", "true")
}

var NSubs int64 = 20_000
var NGateways = 198

// var HardwareID = "b639e30b-1fd4-4f45-934a-9ce44a1cf880"

// var KconfigPath = "/Users/hcgatewood/.kube/config.minikube"
var KconfigPath = "/Users/hcgatewood/Desktop/tmp/orc8r_deployment/kubeconfig_orc8r"

var PostgresClearCmdTemplate = `apt-get update && \
apt-get install -y lsb-release curl ca-certificates gnupg && \
curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' && \
apt-get update && \
apt-get install -y postgresql-client-9.6 && \
psql -d '%s' -c 'DROP SCHEMA public CASCADE ; CREATE SCHEMA public' magma
`
var AdminOperatorCmd = "/var/opt/magma/bin/accessc add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator"

var CloudCertURL = "staging.testminster.com"

var (
	RootCertFilepath  = "/Users/hcgatewood/Desktop/tmp/orc8r_deployment/secrets/rootCA.pem"
	ClientPFXFilepath = "/Users/hcgatewood/Desktop/tmp/orc8r_deployment/secrets/admin_operator.pfx"
	ClientPFXPassword = ""
)

// var CloudURL = "localhost:9443"
var CloudURL = CloudCertURL

var (
	shouldRegenClusterInfo = true
	shouldClean            = true
	shouldPortForward      = true
	shouldCreate           = true
	shouldResetGateways    = true
)

func TestBasicScale(t *testing.T) {
	glog.Info("Get cluster info")
	cluster, err := testlib.GetClusterFromFile()
	assert.NoError(t, err)
	if shouldRegenClusterInfo {
		cluster, err = testlib.GetClusterInfo()
		assert.NoError(t, err)
		err = testlib.SetClusterInfo(cluster)
		assert.NoError(t, err)
	}

	glog.Info("Read server cert")
	rootCAs, err := x509.SystemCertPool()
	require.NoError(t, err)
	pemBytes, err := ioutil.ReadFile(RootCertFilepath)
	require.NoError(t, err)
	ok := rootCAs.AppendCertsFromPEM(pemBytes)
	require.True(t, ok)

	glog.Info("Read client cert")
	pfxBytes, err := ioutil.ReadFile(ClientPFXFilepath)
	require.NoError(t, err)
	key, cert, err := pkcs12.Decode(pfxBytes, ClientPFXPassword)
	require.NoError(t, err)

	glog.Info("Get REST API client")
	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		ServerName:   fmt.Sprintf("*.%s", CloudCertURL),
		Certificates: []tls.Certificate{{Certificate: [][]byte{cert.Raw}, PrivateKey: key}},
	}
	tlsConfig.BuildNameToCertificate()
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	openAPIClient := oclient.NewWithClient(fmt.Sprintf("api.%s", CloudURL), client.DefaultBasePath, []string{"https"}, httpClient)
	c := client.New(openAPIClient, nil)

	glog.Info("Get K8s client")
	kconfigBytes, err := ioutil.ReadFile(KconfigPath)
	require.NoError(t, err)
	kRestConfig, err := clientcmd.RESTConfigFromKubeConfig(kconfigBytes)
	require.NoError(t, err)
	k, err := kubernetes.NewForConfig(kRestConfig)
	require.NoError(t, err)

	promClient, err := api.NewClient(api.Config{Address: "http://localhost:9090"})
	require.NoError(t, err)
	p := prom.NewAPI(promClient)

	glog.Info("Get DB connection string")
	connStrSecret, err := k.CoreV1().Secrets("orc8r").Get("orc8r-controller", metav1.GetOptions{})
	require.NoError(t, err)
	connStr := string(connStrSecret.Data["postgres.connstr"])
	require.NotEmpty(t, connStr)

	if shouldClean {
		glog.Info("Clear DB")
		pod := getOrchestratorPod(t, k)
		postgresClearCmd := fmt.Sprintf(PostgresClearCmdTemplate, connStr)
		kubectlExec(t, k, kRestConfig, pod, postgresClearCmd)

		glog.Info("Bounce application pods")
		err = k.CoreV1().Pods("orc8r").DeleteCollection(nil, metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=orc8r"})
		require.NoError(t, err)
		err = k.CoreV1().Pods("orc8r").DeleteCollection(nil, metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=lte-orc8r"})
		require.NoError(t, err)

		glog.Info("Wait for pods to come up")
		time.Sleep(1 * time.Minute)

		glog.Info("Create admin operator")
		pod = getOrchestratorPod(t, k)
		kubectlExec(t, k, kRestConfig, pod, AdminOperatorCmd)

		time.Sleep(2 * time.Second)
	}

	if shouldPortForward {
		glog.Info("Port-forward")
		// cancelNginx := portForward("svc/orc8r-nginx-proxy", "7443:8443", "7444:8444", "9443:443")
		cancelProm := portForward("svc/orc8r-prometheus", "9090:9090")
		// defer cancelNginx()
		defer cancelProm()

		time.Sleep(2 * time.Second)
	}

	// Assert (instead of require) from here to gracefully handle calling the defer

	glog.Info("Test REST API is up")
	_, err = c.LTENetworks.GetLTE(&lte_networks.GetLTEParams{Context: context.Background()})
	assert.NoError(t, err)

	if shouldCreate {
		glog.Info("Create test network")
		_, err = c.LTENetworks.PostLTE(&lte_networks.PostLTEParams{
			Context: context.Background(),
			LTENetwork: &models.LTENetwork{
				ID:          "test",
				Name:        "test lte network",
				Description: "for testing",
				DNS: &models.NetworkDNSConfig{
					EnableCaching: swag.Bool(false),
					LocalTTL:      swag.Uint32(60),
				},
				Cellular: &models.NetworkCellularConfigs{
					Epc: &models.NetworkEpcConfigs{
						CloudSubscriberdbEnabled: false,
						DefaultRuleID:            "",
						GxGyRelayEnabled:         swag.Bool(false),
						HssRelayEnabled:          swag.Bool(false),
						LTEAuthAmf:               []byte("gA"),
						LTEAuthOp:                []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
						Mcc:                      "001",
						Mnc:                      "01",
						Tac:                      1,
					},
					Ran: &models.NetworkRanConfigs{
						BandwidthMhz: 20,
						TddConfig: &models.NetworkRanConfigsTddConfig{
							Earfcndl:               44590,
							SubframeAssignment:     2,
							SpecialSubframePattern: 7,
						},
					},
				},
			}})
		assert.NoError(t, err)
		nRes, err := c.LTENetworks.GetLTE(&lte_networks.GetLTEParams{Context: context.Background()})
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"test"}, nRes.Payload)

		glog.Info("Create tier")
		_, err = c.Upgrades.PostNetworksNetworkIDTiers(&upgrades.PostNetworksNetworkIDTiersParams{
			Context:   context.Background(),
			NetworkID: "test",
			Tier: &models.Tier{
				Gateways: models.TierGateways{},
				ID:       "default",
				Images:   models.TierImages{},
				Version:  "0.0.0-0",
			},
		})
		assert.NoError(t, err)
		tRes, err := c.Upgrades.GetNetworksNetworkIDTiers(&upgrades.GetNetworksNetworkIDTiersParams{
			Context:   context.Background(),
			NetworkID: "test",
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, []models.TierID{"default"}, tRes.Payload)

		glog.Info("Create APN")
		_, err = c.ApNs.PostLTENetworkIDAPNS(&a_p_ns.PostLTENetworkIDAPNSParams{
			Context:   context.Background(),
			NetworkID: "test",
			APN: &models.APN{
				APNName: "apn0",
				APNConfiguration: &models.APNConfiguration{
					Ambr: &models.AggregatedMaximumBitrate{
						MaxBandwidthDl: swag.Uint32(100),
						MaxBandwidthUl: swag.Uint32(100),
					},
					QosProfile: &models.QosProfile{
						ClassID:                 swag.Int32(9),
						PreemptionCapability:    swag.Bool(true),
						PreemptionVulnerability: swag.Bool(false),
						PriorityLevel:           swag.Uint32(15),
					},
				},
			},
		})
		aRes, err := c.ApNs.GetLTENetworkIDAPNS(&a_p_ns.GetLTENetworkIDAPNSParams{
			Context:   context.Background(),
			NetworkID: "test",
		})
		assert.NoError(t, err)
		assert.Len(t, aRes.Payload, 1)
		assert.NotNil(t, aRes.Payload["apn0"])

		glog.Info("Create subscribers")
		subs := make(models.MutableSubscribers, 0, NSubs)
		for i := int64(0); i < NSubs; i++ {
			sub := &models.MutableSubscriber{
				ID: models.SubscriberID(fmt.Sprintf("IMSI%015d", i)),
				LTE: &models.LTESubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				ActiveAPNS: models.APNList{"apn0"},
			}
			subs = append(subs, sub)
		}
		// Chunk into separate calls to evade timeouts
		glog.Info("Create subscribers a")
		_, err = c.Subscribers.PostLTENetworkIDSubscribersV2(&subscribers.PostLTENetworkIDSubscribersV2Params{
			Context:     context.Background(),
			NetworkID:   "test",
			Subscribers: subs[:5000],
		})
		require.NoError(t, err)
		glog.Info("Create subscribers b")
		_, err = c.Subscribers.PostLTENetworkIDSubscribersV2(&subscribers.PostLTENetworkIDSubscribersV2Params{
			Context:     context.Background(),
			NetworkID:   "test",
			Subscribers: subs[5_000:10_000],
		})
		require.NoError(t, err)
		glog.Info("Create subscribers c")
		_, err = c.Subscribers.PostLTENetworkIDSubscribersV2(&subscribers.PostLTENetworkIDSubscribersV2Params{
			Context:     context.Background(),
			NetworkID:   "test",
			Subscribers: subs[10_000:15_000],
		})
		require.NoError(t, err)
		glog.Info("Create subscribers d")
		_, err = c.Subscribers.PostLTENetworkIDSubscribersV2(&subscribers.PostLTENetworkIDSubscribersV2Params{
			Context:     context.Background(),
			NetworkID:   "test",
			Subscribers: subs[15_000:20_000],
		})
		require.NoError(t, err)
		sRes, err := c.Subscribers.GetLTENetworkIDSubscribersV2(&subscribers.GetLTENetworkIDSubscribersV2Params{
			Context:   context.Background(),
			NetworkID: "test",
			PageSize:  swag.Uint32(1000),
		})
		assert.NoError(t, err)
		assert.Len(t, sRes.Payload.Subscribers, 1000)
		assert.NotEmpty(t, sRes.Payload.NextPageToken)
		assert.Equal(t, NSubs, sRes.Payload.TotalCount)

		glog.Info("Register gateways")
		for _, gateway := range cluster.Gateways {
			_, err = c.LTEGateways.PostLTENetworkIDGateways(&lte_gateways.PostLTENetworkIDGatewaysParams{
				Context:   context.Background(),
				NetworkID: "test",
				Gateway:   testlib.GetDefaultLteGateway(gateway.ID, gateway.HardwareID),
			})
			assert.NoError(t, err)
		}
		gRes, err := c.LTEGateways.GetLTENetworkIDGateways(&lte_gateways.GetLTENetworkIDGatewaysParams{
			Context:   context.Background(),
			NetworkID: "test",
		})
		assert.NoError(t, err)
		assert.Len(t, gRes.Payload, NGateways)
	}

	if shouldResetGateways {
		glog.Info("Clear each gateway's session certs")
		_, err = cluster.RunCmdOnGateways("sudo rm /var/opt/magma/certs/gateway.* && sudo service magma@* stop && sudo service magma@magmad start")
		require.NoError(t, err)

		glog.Info("Wait for gateways to report some data")
		time.Sleep(2 * time.Minute)
	}

	glog.Info("Ensure gateways checked in")
	for _, gw := range cluster.Gateways {
		assertGatewayCheckedIn(t, p, gw.ID)
	}

	val, warn, err := p.Query(context.Background(), `grpc_server_handling_seconds_count{grpc_method="ListSubscribers",grpc_service="magma.lte.SubscriberDBCloud"}`, time.Now())
	assert.NoError(t, err)
	assert.Empty(t, warn)
	glog.Infof("total ListSubscribers requests: %v", val)

	val, warn, err = p.Query(context.Background(), `grpc_server_handling_seconds_bucket{grpc_method="ListSubscribers",grpc_service="magma.lte.SubscriberDBCloud"}`, time.Now())
	assert.NoError(t, err)
	assert.Empty(t, warn)
	glog.Infof("ListSubscribers request buckets: %v", val)
}

func kubectlExec(t *testing.T, k *kubernetes.Clientset, kRestConfig *rest.Config, pod, cmd string) {
	options := &v1.PodExecOptions{
		Command: []string{"sh", "-c", cmd},
		Stdin:   false,
		Stdout:  true,
		Stderr:  false,
		TTY:     false,
	}
	req := k.CoreV1().RESTClient().Post().Resource("pods").
		Name(pod).
		Namespace("orc8r").
		SubResource("exec").
		VersionedParams(options, scheme.ParameterCodec)
	exc, err := remotecommand.NewSPDYExecutor(kRestConfig, "POST", req.URL())
	require.NoError(t, err)
	stdout := &bytes.Buffer{}
	err = exc.Stream(remotecommand.StreamOptions{Stdout: stdout})
	require.NoError(t, err, stdout.String())
}

func getOrchestratorPod(t *testing.T, k *kubernetes.Clientset) string {
	oPods, err := k.CoreV1().Pods("orc8r").List(metav1.ListOptions{LabelSelector: "app.kubernetes.io/component=orchestrator"})
	require.NoError(t, err)
	require.NotEmpty(t, oPods.Items)
	podName := oPods.Items[0].Name
	return podName
}

func portForward(obj string, portPairs ...string) func() {
	args := append([]string{"--namespace", "orc8r", "port-forward", obj}, portPairs...)
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", KconfigPath))
	go func() {
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}()
	return cancel
}

func assertGatewayCheckedIn(t *testing.T, p prom.API, gatewayID string) {
	val, warn, err := p.Query(context.Background(), fmt.Sprintf(`gateway_checkin_status{gatewayID="%s"}`, gatewayID), time.Now())
	assert.NoError(t, err)
	assert.Empty(t, warn)
	assert.Equal(t, model.ValVector, val.Type())

	checkedIn, ok := val.(model.Vector)
	if !assert.True(t, ok) {
		return
	}
	assert.NotEmpty(t, checkedIn)
	assert.Equal(t, model.SampleValue(1), checkedIn[0].Value)
}
