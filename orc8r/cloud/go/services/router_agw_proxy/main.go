package main

import (
	"log"
	"net/http"
	"time"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"magma/orc8r/cloud/go/services/router_agw_proxy/metrics"
	"magma/orc8r/cloud/go/services/router_agw_proxy/poller"
)

func main() {
	apiURL := os.Getenv("API_URL")
	tokenURL := os.Getenv("TOKEN_URL")

	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASSWORD")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
		
	metrics.Init()
	poller.Start(apiURL, tokenURL, username, password, clientID, clientSecret, 60*time.Second)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter running at ... at :9097/metrics")
	log.Fatal(http.ListenAndServe(":9097", nil))
}
