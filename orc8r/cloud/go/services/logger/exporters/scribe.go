/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type HttpClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type ScribeExporter struct {
	scribeUrl      string
	appId          string
	appSecret      string
	queue          []*ScribeLogEntry
	queueMutex     sync.RWMutex
	queueLen       int
	exportInterval time.Duration
}

func NewScribeExporter(
	baseUrl string,
	appId string,
	appSecret string,
	queueLen int,
	exportInterval time.Duration,
) *ScribeExporter {
	e := new(ScribeExporter)
	e.scribeUrl = baseUrl
	e.appId = appId
	e.appSecret = appSecret
	e.queueLen = queueLen
	e.exportInterval = exportInterval
	return e
}

func (e *ScribeExporter) Start() {
	go e.exportEvery()
}

func (e *ScribeExporter) exportEvery() {
	for _ = range time.Tick(e.exportInterval) {
		client := http.DefaultClient
		err := e.Export(client)
		if err != nil {
			glog.Errorf("Error in exporting to scribe: %v\n", err)
		}
	}
}

// Write to Scribe
func (e *ScribeExporter) Export(client HttpClient) error {
	e.queueMutex.RLock()
	logs := e.queue
	e.queueMutex.RUnlock()
	if len(logs) != 0 {
		err := e.write(client, logs)
		if err != nil {
			return fmt.Errorf("Failed to export to scribe: %v\n", err)
		}
		// write to ods successful, clear written logs from queue
		e.queueMutex.Lock()
		e.queue = e.queue[len(logs):]
		e.queueMutex.Unlock()
	}
	return nil
}

func (e *ScribeExporter) write(client HttpClient, logEntries []*ScribeLogEntry) error {
	logJson, err := json.Marshal(logEntries)
	if err != nil {
		return err
	}
	accessToken := fmt.Sprintf("%s|%s", e.appId, e.appSecret)
	resp, err := client.PostForm(e.scribeUrl,
		url.Values{"access_token": {accessToken}, "logs": {string(logJson)}})
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = fmt.Errorf("Scribe status code %d: %s", resp.StatusCode, errMsg)
	}
	return err
}

func (e *ScribeExporter) Submit(logEntries []*protos.LogEntry) error {
	e.queueMutex.Lock()
	defer e.queueMutex.Unlock()
	if (len(e.queue) + len(logEntries)) > e.queueLen {
		// queue is full, clear queue and log that queue was full
		e.queue = []*ScribeLogEntry{}
		glog.Warningf("Queue is full, clearing...")
		if len(logEntries) > e.queueLen {
			return fmt.Errorf("dropping %v logEntries as it exceeds max queue length", len(logEntries))
		}
	}
	scribeEntries, err := ConvertToScribeLogEntries(logEntries)
	if err != nil {
		return err
	}
	e.queue = append(e.queue, scribeEntries...)
	return nil
}
