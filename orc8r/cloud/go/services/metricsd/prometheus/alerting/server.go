package main

import (
	"flag"
	"fmt"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/handlers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort          = "9093"
	defaultPrometheusURL = "localhost:9090"
	rootPath             = "/:network_id"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	rulesDir := flag.String("rules-dir", ".", "Directory to write rules files. Default is '.'")
	prometheusURL := flag.String("prometheusURL", "localhost:9090", fmt.Sprintf("URL of the prometheus instance that is reading these rules. Default is %s", defaultPrometheusURL))
	flag.Parse()

	client, err := alert.NewClient(*rulesDir)
	if err != nil {
		glog.Errorf("error creating alert client: %v", err)
		return
	}

	e := echo.New()

	e.POST(rootPath, handlers.GetPostHandler(client, *prometheusURL))
	e.GET(rootPath, handlers.GetGetHandler(client))

	glog.Infof("Alertconfig server listening on Port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
