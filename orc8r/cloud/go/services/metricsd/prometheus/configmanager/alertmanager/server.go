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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort                   = "9101"
	defaultAlertmanagerURL        = "alertmanager:9093"
	defaultAlertmanagerConfigPath = "./alertmanager.yml"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	alertmanagerConfPath := flag.String("alertmanager-conf", defaultAlertmanagerConfigPath, fmt.Sprintf("Path to alertmanager configuration file. Default is %s", defaultAlertmanagerConfigPath))
	alertmanagerURL := flag.String("alertmanagerURL", defaultAlertmanagerURL, fmt.Sprintf("URL of the alertmanager instance that is being used. Default is %s", defaultAlertmanagerURL))
	flag.Parse()

	e := echo.New()

	e.GET("/", statusHandler)

	receiverClient := receivers.NewClient(*alertmanagerConfPath, *alertmanagerURL, fsclient.NewFSClient())
	e.POST(ReceiverPath, GetReceiverPostHandler(receiverClient))
	e.GET(ReceiverPath, GetGetReceiversHandler(receiverClient))
	e.DELETE(ReceiverPath, GetDeleteReceiverHandler(receiverClient))
	e.PUT(ReceiverPath+"/:"+ReceiverNamePathParam, GetUpdateReceiverHandler(receiverClient))

	e.POST(RoutePath, GetUpdateRouteHandler(receiverClient))
	e.GET(RoutePath, GetGetRouteHandler(receiverClient))

	glog.Infof("Alertmanager Config server listening on port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}

func statusHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Alertmanager Config server")
}
