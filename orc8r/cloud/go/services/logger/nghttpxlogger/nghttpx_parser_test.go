/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package nghttpxlogger_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/logger/nghttpxlogger"

	"github.com/stretchr/testify/assert"
)

func TestNghttpxParser_Parse(t *testing.T) {
	str := "2018-05-11T16:49:10.657Z@|@192.168.80.1@|@192.168.80.10:9443@|@" +
		"9443@|@GET / HTTP/2@|@200@|@45@|@0.005@|@h2@|@" +
		"fcd2f92750744ead7921876870c3c277@|@CN=test_operator,OU=,O=,C=US@|@r@|@" +
		"-@|@TLSv1.2@|@ECDHE-RSA-AES256-GCM-SHA384@|@127.0.0.1@|@9081"
	parser, err := nghttpxlogger.NewNghttpParser()
	assert.NoError(t, err)
	msg, err := parser.Parse(str)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(msg.Int))
	assert.Equal(t, 13, len(msg.Normal))

	assert.Equal(t, msg.Time, int64(1526057350))
	assert.Equal(t, msg.Int["server_port"], int64(9443))
	assert.Equal(t, msg.Int["backend_port"], int64(9081))
	assert.Equal(t, msg.Int["body_bytes_sent"], int64(45))
	assert.Equal(t, msg.Int["request_time_micro_secs"], int64(5))
	assert.Equal(t, msg.Normal["alpn"], "h2")
	assert.Equal(t, msg.Normal["backend_host"], "127.0.0.1")
	assert.Equal(t, msg.Normal["client_ip"], "192.168.80.1")
	assert.Equal(t, msg.Normal["http_host"], "192.168.80.10:9443")
	assert.Equal(t, msg.Normal["request_method"], "GET")
	assert.Equal(t, msg.Normal["request_url"], "/")
	assert.Equal(t, msg.Normal["status"], "200")
	assert.Equal(t, msg.Normal["tls_cipher"], "ECDHE-RSA-AES256-GCM-SHA384")
	assert.Equal(t, msg.Normal["tls_client_serial"], "fcd2f92750744ead7921876870c3c277")
	assert.Equal(t, msg.Normal["tls_client_subject_name"], "CN=test_operator,OU=,O=,C=US")
	assert.Equal(t, msg.Normal["tls_protocol"], "TLSv1.2")
	assert.Equal(t, msg.Normal["tls_session_reused"], "r")
	assert.Equal(t, msg.Normal["tls_sni"], "-")
}

func TestNghttpxParser_ParseSuccessWithWrongFormat(t *testing.T) {
	// request_time is populated with "-" instead of a valid number
	// backend_port is populated with "-" instead of a valid number
	// client_requst is populated with wrong format which does not have space as the delimiter
	// Those fields should be ignored and parsed without any issue
	str := "2018-05-11T16:49:10.657Z@|@192.168.80.1@|@192.168.80.10:9443@|@" +
		"9443@|@GET/HTTP/2@|@200@|@45@|@-@|@h2@|@" +
		"fcd2f92750744ead7921876870c3c277@|@CN=test_operator,OU=,O=,C=US@|@r@|@" +
		"-@|@TLSv1.2@|@ECDHE-RSA-AES256-GCM-SHA384@|@127.0.0.1@|@-"

	parser, err := nghttpxlogger.NewNghttpParser()
	assert.NoError(t, err)
	msg, err := parser.Parse(str)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(msg.Int))
	assert.Equal(t, 11, len(msg.Normal))

	assert.Equal(t, msg.Time, int64(1526057350))
	assert.Equal(t, msg.Int["server_port"], int64(9443))
	assert.Equal(t, msg.Int["body_bytes_sent"], int64(45))
	assert.Equal(t, msg.Normal["alpn"], "h2")
	assert.Equal(t, msg.Normal["backend_host"], "127.0.0.1")
	assert.Equal(t, msg.Normal["client_ip"], "192.168.80.1")
	assert.Equal(t, msg.Normal["http_host"], "192.168.80.10:9443")
	assert.Equal(t, msg.Normal["status"], "200")
	assert.Equal(t, msg.Normal["tls_cipher"], "ECDHE-RSA-AES256-GCM-SHA384")
	assert.Equal(t, msg.Normal["tls_client_serial"], "fcd2f92750744ead7921876870c3c277")
	assert.Equal(t, msg.Normal["tls_client_subject_name"], "CN=test_operator,OU=,O=,C=US")
	assert.Equal(t, msg.Normal["tls_protocol"], "TLSv1.2")
	assert.Equal(t, msg.Normal["tls_session_reused"], "r")
	assert.Equal(t, msg.Normal["tls_sni"], "-")
}

func TestNghttpxParser_ParseFail(t *testing.T) {
	// timestamp is not in the correct format
	str := "2018657Z@|@192.168.80.1@|@192.168.80.10:9443@|@" +
		"9443@|@GET / HTTP/2@|@200@|@45@|@0.005@|@h2@|@" +
		"fcd2f92750744ead7921876870c3c277@|@CN=test_operator,OU=,O=,C=US@|@r@|@" +
		"-@|@TLSv1.2@|@ECDHE-RSA-AES256-GCM-SHA384@|@127.0.0.1@|@9081"
	parser, err := nghttpxlogger.NewNghttpParser()
	assert.NoError(t, err)
	_, err = parser.Parse(str)
	assert.EqualError(t, err, "Error parsing time: parsing time \"2018657Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"657Z\" as \"-\"")
}
