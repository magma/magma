/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/handlers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort                   = "9093"
	defaultAlertmanagerURL        = "localhost:9092"
	defaultAlertmanagerConfigPath = "./alertmanager.yml"

	rootPath     = "/:network_id"
	receiverPath = rootPath + "/receiver"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	alertmanagerConfPath := flag.String("alertmanager-conf", "./alertmanager.yml", fmt.Sprintf("Path to alertmanager configuration file. Default is %s", defaultAlertmanagerConfigPath))
	alertmanagerURL := flag.String("alertmanagerURL", "localhost:9092", fmt.Sprintf("URL of the alertmanager instance that is being used. Default is %s", defaultAlertmanagerURL))
	flag.Parse()

	e := echo.New()

	e.GET("/", statusHandler)

	receiverClient := receivers.NewClient(*alertmanagerConfPath)
	e.POST(receiverPath, handlers.GetReceiverPostHandler(receiverClient, *alertmanagerURL))
	e.GET(receiverPath, handlers.GetGetReceiversHandler(receiverClient))
	e.DELETE(receiverPath, handlers.GetDeleteReceiverHandler(receiverClient, *alertmanagerURL))
	e.PUT(receiverPath+"/:"+handlers.ReceiverNamePathParam, handlers.GetUpdateReceiverHandler(receiverClient, *alertmanagerURL))

	e.POST(receiverPath+"/route", handlers.GetUpdateRouteHandler(receiverClient, *alertmanagerURL))
	e.GET(receiverPath+"/route", handlers.GetGetRouteHandler(receiverClient))

	glog.Infof("Alertmanager Config server listening on port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}

func statusHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Alertmanager Config server")
}
