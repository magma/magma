package main

import (
	"flag"
	"fmt"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/handlers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort                   = "9093"
	defaultPrometheusURL          = "localhost:9090"
	defaultAlertmanagerURL        = "localhost:9092"
	defaultAlertmanagerConfigPath = "./alertmanager.yml"

	rootPath     = "/:network_id"
	alertPath    = rootPath + "/alert"
	receiverPath = rootPath + "/receiver"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	rulesDir := flag.String("rules-dir", ".", "Directory to write rules files. Default is '.'")
	prometheusURL := flag.String("prometheusURL", "localhost:9090", fmt.Sprintf("URL of the prometheus instance that is reading these rules. Default is %s", defaultPrometheusURL))
	alertmanagerConfPath := flag.String("alertmanager-conf", "./alertmanager.yml", fmt.Sprintf("Path to alertmanager configuration file. Default is %s", defaultAlertmanagerConfigPath))
	alertmanagerURL := flag.String("alertmanagerURL", "localhost:9092", fmt.Sprintf("URL of the alertmanager instance that is being used. Default is %s", defaultAlertmanagerURL))
	flag.Parse()

	e := echo.New()

	alertClient, err := alert.NewClient(*rulesDir)
	if err != nil {
		glog.Errorf("error creating alert client: %v", err)
		return
	}
	e.POST(alertPath, handlers.GetPostHandler(alertClient, *prometheusURL))
	e.GET(alertPath, handlers.GetGetHandler(alertClient))
	e.DELETE(alertPath, handlers.GetDeleteHandler(alertClient, *prometheusURL))

	receiverClient := receivers.NewClient(*alertmanagerConfPath)
	e.POST(receiverPath, handlers.GetReceiverPostHandler(receiverClient, *alertmanagerURL))
	e.GET(receiverPath, handlers.GetGetReceiversHandler(receiverClient))

	e.POST(receiverPath+"/route", handlers.GetUpdateRouteHandler(receiverClient))
	e.GET(receiverPath+"/route", handlers.GetGetRouteHandler(receiverClient))

	glog.Infof("Alertconfig server listening on Port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
