/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package nghttpxlogger

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

type NghttpxParser interface {
	Parse(str string) (*NghttpxMessage, error)
}

type NghttpxParserImpl struct {
}

//${time_iso8601}@|@${remote_addr}@|@${http_host}@|@${server_port}@|@${request}
//@|@${status}@|@${body_bytes_sent}@|@${request_time}@|@${alpn}@|@
//${tls_client_serial}@|@${tls_client_subject_name}@|@${tls_session_reused}@|@
//${tls_sni}@|@${tls_protocol}@|@${tls_cipher}@|@${backend_host}
//@|@${backend_port}
const (
	DELIMITER     = "@|@"
	NUM_OF_FIELDS = 17 // num of variables in the original string to parse
)

func NewNghttpParser() (*NghttpxParserImpl, error) {
	return &NghttpxParserImpl{}, nil
}

func (parser *NghttpxParserImpl) Parse(str string) (*NghttpxMessage, error) {
	str = strings.Trim(str, "\x00")
	res := strings.Split(str, DELIMITER)
	if len(res) < NUM_OF_FIELDS {
		return nil, fmt.Errorf("Expected # of fields:%v, got: %v", NUM_OF_FIELDS, len(res))
	}

	builder := NewNghttpxScribeDataBuilder()
	intMsg, normalMsg, time, errs := builder.
		Time(res[0]).
		StringField("client_ip", res[1]).
		StringField("http_host", res[2]).
		IntField("server_port", res[3]).
		ClientRequest(res[4]).
		StringField("status", res[5]).
		IntField("body_bytes_sent", res[6]).
		RequestTime(res[7]).
		StringField("alpn", res[8]).
		StringField("tls_client_serial", res[9]).
		StringField("tls_client_subject_name", res[10]).
		StringField("tls_session_reused", res[11]).
		StringField("tls_sni", res[12]).
		StringField("tls_protocol", res[13]).
		StringField("tls_cipher", res[14]).
		StringField("backend_host", res[15]).
		IntField("backend_port", res[16]).
		Build()

	msg := NghttpxMessage{Normal: normalMsg, Int: intMsg, Time: time}
	if len(errs) != 0 {
		return nil, errs[0]
	}
	return &msg, nil
}

type NghttpxScribeDataBuilder struct {
	normalMsg map[string]string
	intMsg    map[string]int64
	time      int64
	errors    []error
}

func NewNghttpxScribeDataBuilder() *NghttpxScribeDataBuilder {
	return &NghttpxScribeDataBuilder{
		normalMsg: map[string]string{},
		intMsg:    map[string]int64{},
		errors:    []error{},
	}
}

func (builder *NghttpxScribeDataBuilder) Time(token string) *NghttpxScribeDataBuilder {
	t, err := time.Parse(time.RFC3339, token)
	if err != nil {
		builder.errors = append(builder.errors, fmt.Errorf("Error parsing time: %v", err))
	} else {
		builder.time = int64(t.Unix())
	}
	return builder
}

func (builder *NghttpxScribeDataBuilder) StringField(name string, value string) *NghttpxScribeDataBuilder {
	builder.normalMsg[name] = value
	return builder
}

func (builder *NghttpxScribeDataBuilder) IntField(name string, value string) *NghttpxScribeDataBuilder {
	if value == "-" {
		glog.V(2).Infof("Cannot parse %s field with value: %s, ignored", name, value)
		return builder
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		builder.errors = append(builder.errors, fmt.Errorf("Error parsing %s: %v", name, err))
	} else {
		builder.intMsg[name] = int64(intValue)
	}
	return builder
}

func (builder *NghttpxScribeDataBuilder) ClientRequest(token string) *NghttpxScribeDataBuilder {
	// parse request, setting RequestMethod and RequestUrl
	res := strings.Split(token, " ")
	if len(res) != 3 {
		glog.V(2).Infof("Cannot parse client_request field with value: %s, ignored", token)
	} else {
		builder.normalMsg["request_method"] = res[0]
		builder.normalMsg["request_url"] = res[1]
	}
	return builder
}

func (builder *NghttpxScribeDataBuilder) RequestTime(token string) *NghttpxScribeDataBuilder {
	reqT, err := strconv.ParseFloat(token, 64)
	if err != nil {
		glog.V(2).Infof("Cannot parse request_time field with value: %s, ignored", token)
	} else {
		builder.intMsg["request_time_micro_secs"] = int64(reqT * 1000)
	}
	return builder
}

func (builder *NghttpxScribeDataBuilder) Build() (map[string]int64, map[string]string, int64, []error) {
	return builder.intMsg, builder.normalMsg, builder.time, builder.errors
}
