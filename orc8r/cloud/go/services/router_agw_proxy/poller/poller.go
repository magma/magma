package poller

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"magma/orc8r/cloud/go/services/router_agw_proxy/auth"
	"magma/orc8r/cloud/go/services/router_agw_proxy/metrics"
	"magma/orc8r/cloud/go/services/router_agw_proxy/types"
)

func Start(apiURL, tokenURL, username, password, clientID, clientSecret string, interval time.Duration) {
	go func() {
		for {
			networksList, err := CallOrc8rNetworks()
			if err != nil {
				log.Println("Error fetching network list:", err)
				time.Sleep(interval)
				continue
			}

			token, err := auth.GetToken(tokenURL, username, password, clientID, clientSecret)
			if err != nil {
				log.Println("Error getting token:", err)
				time.Sleep(interval)
				continue
			}

			req, _ := http.NewRequest("GET", apiURL, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error calling API:", err)
				time.Sleep(interval)
				continue
			}
			defer resp.Body.Close()

			var apiResp types.APIResponse
			if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
				log.Println("Error decoding JSON:", err)
				time.Sleep(interval)
				continue
			}

			for _, network := range networksList {
				for _, item := range apiResp.Items {
					if item.Data.Imei == "" {
						continue
					}

					labels := prometheus.Labels{
						"name":      item.Name,
						"type":      item.Type,
						"imei":      item.Data.Imei,
						"networkID": network,
					}

					if item.Data.Latitude != nil {
						metrics.DeviceLatitude.With(labels).Set(*item.Data.Latitude)
					}
					if item.Data.Longitude != nil {
						metrics.DeviceLongitude.With(labels).Set(*item.Data.Longitude)
					}

					if item.Data.BytesSent != nil && item.Data.BytesReceived != nil {
						sent, err1 := strconv.Atoi(*item.Data.BytesSent)
						recv, err2 := strconv.Atoi(*item.Data.BytesReceived)
						if err1 == nil && err2 == nil && sent > 0 && recv > 0 {
							metrics.DevicePrivate5GActive.With(labels).Set(1)
						} else {
							metrics.DevicePrivate5GActive.With(labels).Set(0)
						}
					} else {
						metrics.DevicePrivate5GActive.With(labels).Set(0)
					}

					if item.Data.NetworkServiceType == "5G" {
						metrics.DeviceNetworkService.With(labels).Set(2)
					} else if item.Data.NetworkServiceType == "4G" {
						metrics.DeviceNetworkService.With(labels).Set(1)
					} else {
						metrics.DeviceNetworkService.With(labels).Set(0)
					}

					if item.Data.Rsrp != nil {
						val, err := strconv.ParseFloat(*item.Data.Rsrp, 64)
						if err == nil {
							metrics.DeviceRsrp.With(labels).Set(val)
						}
					}
					if item.Data.Rsrq != nil {
						val, err := strconv.ParseFloat(*item.Data.Rsrq, 64)
						if err == nil {
							metrics.DeviceRsrq.With(labels).Set(val)
						}
					}
					if item.Data.Rssi != nil {
						val, err := strconv.ParseFloat(*item.Data.Rssi, 64)
						if err == nil {
							metrics.DeviceRssi.With(labels).Set(val)
						}
					}
					if item.Data.Snr != nil {
						val, err := strconv.ParseFloat(*item.Data.Snr, 64)
						if err == nil {
							metrics.DeviceSnr.With(labels).Set(val)
						}
					}
					if item.Data.Speed != nil {
						val, err := strconv.ParseFloat(*item.Data.Speed, 64)
						if err == nil {
							metrics.DeviceSpeed.With(labels).Set(val)
						}
					}
					if item.Data.RadioModuleTemp != nil {
						val, err := strconv.ParseFloat(*item.Data.RadioModuleTemp, 64)
						if err == nil {
							metrics.DeviceRadioModuleTemp.With(labels).Set(val)
						}
					}
					if item.Data.BoardTemp != nil {
						val, err := strconv.ParseFloat(*item.Data.BoardTemp, 64)
						if err == nil {
							metrics.DeviceBoardTemp.With(labels).Set(val)
						}
					}
					if item.CommStatus == "OK" {
						metrics.DeviceCommStatus.With(labels).Set(1)
					} else {
						metrics.DeviceCommStatus.With(labels).Set(0)
					}
					if item.Data.BytesSent != nil {
						val, err := strconv.ParseFloat(*item.Data.BytesSent, 64)
						if err == nil {
							metrics.DeviceBytesSent.With(labels).Set(val)
						}
					}
					if item.Data.BytesReceived != nil {
						val, err := strconv.ParseFloat(*item.Data.BytesReceived, 64)
						if err == nil {
							metrics.DeviceBytesReceived.With(labels).Set(val)
						}
					}
				}
			}

			log.Printf("Updated metrics for %d devices across %d networks\n", len(apiResp.Items), len(networksList))
			time.Sleep(interval)
		}
	}()
}

func CallOrc8rNetworks() ([]string, error) {
	certFile := "/var/opt/test_certs/admin_operator.pem"
	keyFile := "/var/opt/test_certs/admin_operator.key.pem"
	caFile := "/var/opt/test_certs/rootCA.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("Error loading keypair: %v\n", err)
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Printf("Error loading CA: %v\n", err)
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // opcionalmente false se CN estiver correto
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
		Timeout:   10 * time.Second,
	}

	url := "https://orc8r_nginx_1:9443/magma/v1/networks"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling endpoint: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error on response (%s): %s\n", resp.Status, string(body))
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var networks []string
	if err := json.Unmarshal(body, &networks); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	log.Printf("Found %d networks: %v", len(networks), networks)
	return networks, nil
}