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
	graphiteClients   []*graphiteClient
	registeredMetrics map[string]GraphiteMetric
	sync.Mutex
}

type Address struct {
	Host string
	Port int
}

func NewAddress(host string, port int) Address {
	if strings.HasPrefix(host, "http") {
		host = host[strings.LastIndex(host, "/")+1:]
	}
	return Address{
		Host: host,
		Port: port,
	}
}

type graphiteClient struct {
	client    *graphite.Graphite
	connected bool
	address   Address
}

func NewGraphiteClient(address Address) *graphiteClient {
	var client *graphite.Graphite
	connected := true
	if address.Host == "" {
		glog.Info("Created No-Op graphite exporter because of empty graphite address\n")
		client = graphite.NewGraphiteNop(address.Host, address.Port)
	} else {
		var err error
		client, err = graphite.NewGraphite(address.Host, address.Port)
		if err != nil {
			connected = false
			glog.Errorf("Could not connect to graphite address %s on start. Retrying on export", address.Host)
		}
	}
	return &graphiteClient{
		client:    client,
		connected: connected,
		address:   address,
	}
}

// NewGraphiteExporter create a new GraphiteExporter with own registry
func NewGraphiteExporter(graphiteAddresses []Address) exporters.Exporter {
	var graphiteClients = make([]*graphiteClient, 0, len(graphiteAddresses))

	for _, address := range graphiteAddresses {
		graphiteClients = append(graphiteClients, NewGraphiteClient(address))
	}

	return &GraphiteExporter{
		registeredMetrics: make(map[string]GraphiteMetric),
		graphiteClients:   graphiteClients,
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
	e.Lock()
	defer e.Unlock()
	for _, graphiteClient := range e.graphiteClients {
		if !graphiteClient.connected {
			err := graphiteClient.reconnect()
			if err != nil {
				glog.Errorf("Failed to reconnect: %v", err)
			}
			continue
		}
		for _, metric := range e.registeredMetrics {
			err := metric.Export(graphiteClient.client)
			if err != nil {
				graphiteClient.connected = false
				glog.Errorf("Graphite client failed to send to %s:%d %v. Retrying on next export.", graphiteClient.address.Host, graphiteClient.address.Port, err)
				break
			}
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
func (c *graphiteClient) reconnect() error {
	newGraphite, err := graphite.NewGraphite(c.address.Host, c.address.Port)
	if err != nil {
		return fmt.Errorf("Could not connect to graphite address %s:%d on export. Retrying on next export", c.address.Host, c.address.Port)
	}
	c.client = newGraphite
	c.connected = true
	glog.Infof("Successfully created graphite connection on %s:%d. Exporting now.", c.address.Host, c.address.Port)
	return nil
}

func makeGraphiteName(metric *dto.Metric, ctx exporters.MetricsContext) string {
	name := ctx.MetricName
	labels := metric.Label

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
	s[name] = value
}

func (s TagSet) SortedTags() []string {
	var sortedKeys []string
	var tagList []string

	for key := range s {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		tagList = append(tagList, fmt.Sprintf("%s=%s", key, s[key]))
	}
	return tagList
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
