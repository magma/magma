/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package exporters

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/golang/glog"
	"github.com/marpaia/graphite-golang"

	dto "github.com/prometheus/client_model/go"
)

const (
	exportInterval = time.Second * 30

	NetworkTagName = "networkID"
	GatewayTagName = "gatewayID"
	defaultNetwork = "defaultNetwork"
	defaultGateway = "defaultGateway"
)

// GraphiteExporter handles registering and updating graphite metrics
type GraphiteExporter struct {
	graphite          *graphite.Graphite
	registeredMetrics map[string]GraphiteMetric
	connected         bool
	host              string
	port              int
	sync.Mutex
}

// NewGraphiteExporter create a new GraphiteExporter with own registry
func NewGraphiteExporter(graphiteAddress string, graphiteReceivePort int) exporters.Exporter {
	var graphiteObj *graphite.Graphite

	if graphiteAddress == "" {
		graphiteObj = graphite.NewGraphiteNop(graphiteAddress, graphiteReceivePort)
		glog.Error("Created No-Op graphite exporter because of empty graphite address\n")
		return &GraphiteExporter{
			registeredMetrics: make(map[string]GraphiteMetric),
			graphite:          graphite.NewGraphiteNop(graphiteAddress, graphiteReceivePort),
		}
	}
	var connected bool
	graphiteObj, err := graphite.NewGraphite(graphiteAddress, graphiteReceivePort)
	if err != nil {
		connected = false
		glog.Errorf("Could not connect to graphite address %s on start. Retrying on export", graphiteAddress)
	} else {
		connected = true
	}
	return &GraphiteExporter{
		registeredMetrics: make(map[string]GraphiteMetric),
		graphite:          graphiteObj,
		connected:         connected,
		host:              graphiteAddress,
		port:              graphiteReceivePort,
	}
}

// Submit takes in a metric and either registers it to or updates the metric if
// it is already registered
func (e *GraphiteExporter) Submit(metrics []exporters.MetricAndContext) error {
	e.Lock()
	defer e.Unlock()
	for _, metric := range metrics {
		err := e.submitSingleFamilyUnsafe(metric.Family, metric.Context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *GraphiteExporter) submitSingleFamilyUnsafe(family *dto.MetricFamily, ctx exporters.MetricsContext) error {
	for _, metric := range family.GetMetric() {
		registeredName := makeGraphiteName(metric, ctx)
		if registeredMetric, ok := e.registeredMetrics[registeredName]; ok {
			registeredMetric.Update(metric)
			return nil
		}
		e.registerMetric(metric, family, registeredName)
	}
	return nil
}

func (e *GraphiteExporter) registerMetric(metric *dto.Metric,
	family *dto.MetricFamily,
	name string,
) {
	var newMetric GraphiteMetric
	switch family.GetType() {
	case dto.MetricType_COUNTER:
		newMetric = NewGraphiteCounter()
	case dto.MetricType_GAUGE:
		newMetric = NewGraphiteGauge()
	case dto.MetricType_SUMMARY:
		newMetric = NewGraphiteSummary()
	case dto.MetricType_HISTOGRAM:
		newMetric = NewGraphiteHistogram()
	default:
		glog.Errorf("Cannot register unsupported Type: %v", family.GetType())
	}
	newMetric.Register(metric, name, e)
}

// Run ExportEvery() in a goroutine to avoid blocking
func (e *GraphiteExporter) Start() {
	go e.exportEvery()
}

func (e *GraphiteExporter) exportEvery() {
	for range time.Tick(exportInterval) {
		err := e.Export()
		if err != nil {
			glog.Errorf("Error submitting to graphite: %v\n", err)
		}
	}
}

func (e *GraphiteExporter) Export() error {
	if !e.connected {
		err := e.reconnect()
		if err != nil {
			return err
		}
	}
	e.Lock()
	defer e.Unlock()
	for _, metric := range e.registeredMetrics {
		err := metric.Export(e)
		if err != nil {
			return err
		}
	}
	e.clearRegistry()
	return nil
}

// clearRegistry erases the stored metrics map because we don't need to keep
// old metrics around forever if they aren't updated
func (e *GraphiteExporter) clearRegistry() {
	e.registeredMetrics = make(map[string]GraphiteMetric)
}

// reconnect attempts to connect the graphite client to the graphite server
func (e *GraphiteExporter) reconnect() error {
	newGraphite, err := graphite.NewGraphite(e.host, e.port)
	if err != nil {
		return fmt.Errorf("Could not connect to graphite address %s:%d on export. Retrying on next export", e.host, e.port)
	}
	e.graphite = newGraphite
	e.connected = true
	glog.Infof("Successfully created graphite connection on %s:%d. Exporting now.", e.host, e.port)
	return nil
}

func makeGraphiteName(metric *dto.Metric, ctx exporters.MetricsContext) string {
	name := ctx.MetricName
	labels := protos.GetDecodedLabel(metric)

	networkID := ctx.NetworkID
	gatewayID := ctx.GatewayID
	if networkID == "" {
		networkID = defaultNetwork
	}
	if gatewayID == "" {
		gatewayID = defaultGateway
	}

	tags := make(TagSet)
	tags.Insert(NetworkTagName, networkID)
	tags.Insert(GatewayTagName, gatewayID)

	for _, labelPair := range labels {
		tags.Insert(labelPair.GetName(), labelPair.GetValue())
	}
	return name + tags.String()
}

type TagSet map[string]string

func (s TagSet) Insert(name, value string) {
	if _, ok := s[name]; !ok {
		s[name] = value
	}
}

// String prints the tagSet sorted by key in a format that can be appended to
// a graphite metric name
func (s TagSet) String() string {
	var sortedKeys []string

	for key := range s {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var str strings.Builder
	for _, key := range sortedKeys {
		str.WriteString(fmt.Sprintf(";%s=%s", key, s[key]))
	}
	return str.String()
}
