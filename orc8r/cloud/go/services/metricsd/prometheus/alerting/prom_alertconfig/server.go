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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/files"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/handlers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	defaultPort          = "9100"
	defaultPrometheusURL = "prometheus:9090"
)

func main() {
	port := flag.String("port", defaultPort, fmt.Sprintf("Port to listen for requests. Default is %s", defaultPort))
	rulesDir := flag.String("rules-dir", ".", "Directory to write rules files. Default is '.'")
	prometheusURL := flag.String("prometheusURL", defaultPrometheusURL, fmt.Sprintf("URL of the prometheus instance that is reading these rules. Default is %s", defaultPrometheusURL))
	flag.Parse()

	e := echo.New()

	fileLocks, err := alert.NewFileLocker(alert.NewDirectoryClient(*rulesDir))
	alertClient := alert.NewClient(fileLocks, *rulesDir, files.NewFSClient())
	if err != nil {
		glog.Errorf("error creating alert client: %v", err)
		return
	}
	e.GET("/", statusHandler)

	e.POST(handlers.AlertPath, handlers.GetConfigureAlertHandler(alertClient, *prometheusURL))
	e.GET(handlers.AlertPath, handlers.GetRetrieveAlertHandler(alertClient))
	e.DELETE(handlers.AlertPath, handlers.GetDeleteAlertHandler(alertClient, *prometheusURL))
	e.PUT(handlers.AlertUpdatePath, handlers.GetUpdateAlertHandler(alertClient, *prometheusURL))

	e.PUT(handlers.AlertBulkPath, handlers.GetBulkAlertUpdateHandler(alertClient, *prometheusURL))

	glog.Infof("Prometheus Config server listening on port: %s\n", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}

func statusHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Prometheus Config server")
}
