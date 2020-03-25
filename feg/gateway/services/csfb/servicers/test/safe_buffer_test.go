/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/test_utils"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulReadAndWrite(t *testing.T) {
	safeBuffer, err := servicers.NewSafeBuffer()
	assert.NoError(t, err)

	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)

	safeBuffer.WriteChunk(chunk)

	messageType, message, err := safeBuffer.GetNextMessage(1)
	assert.NoError(t, err)
	assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
	expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, message)
}

func TestReadSuccessWithWait(t *testing.T) {
	safeBuffer, err := servicers.NewSafeBuffer()
	assert.NoError(t, err)

	go func() {
		messageType, message, err := safeBuffer.GetNextMessage(1)
		assert.NoError(t, err)
		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
			Imsi: "111111",
		})
		assert.Equal(t, expectedMsg, message)
	}()
	time.Sleep(time.Millisecond * 500)
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)
	safeBuffer.WriteChunk(chunk)
}

func TestReadFailWithWaitTimeout(t *testing.T) {
	safeBuffer, err := servicers.NewSafeBuffer()
	assert.NoError(t, err)

	go func() {
		messageType, message, err := safeBuffer.GetNextMessage(1)
		assert.EqualError(t, err, "buffer read timeout")
		assert.Equal(t, decode.SGsMessageType(0x00), messageType)
		assert.Equal(t, &any.Any{}, message)
	}()
	time.Sleep(time.Second * 2)
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)
	safeBuffer.WriteChunk(chunk)
}

func TestSuccessfulConsecutiveReadWrite(t *testing.T) {
	safeBuffer, err := servicers.NewSafeBuffer()
	assert.NoError(t, err)

	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)

	safeBuffer.WriteChunk(chunk)
	safeBuffer.WriteChunk(chunk)
	safeBuffer.WriteChunk(chunk)

	messageType, message, err := safeBuffer.GetNextMessage(1)
	assert.NoError(t, err)
	assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
	expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, message)

	messageType, message, err = safeBuffer.GetNextMessage(1)
	assert.NoError(t, err)
	assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
	expectedMsg, _ = ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, message)
	messageType, message, err = safeBuffer.GetNextMessage(1)
	assert.NoError(t, err)
	assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
	expectedMsg, _ = ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, message)
}

func TestSuccessfulConsecutiveReadWriteWithWait(t *testing.T) {
	safeBuffer, err := servicers.NewSafeBuffer()
	assert.NoError(t, err)

	go func() {
		messageType, message, err := safeBuffer.GetNextMessage(1)
		assert.NoError(t, err)
		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
			Imsi: "111111",
		})
		assert.Equal(t, expectedMsg, message)
	}()
	go func() {
		messageType, message, err := safeBuffer.GetNextMessage(1)
		assert.NoError(t, err)
		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
			Imsi: "111111",
		})
		assert.Equal(t, expectedMsg, message)
	}()
	go func() {
		messageType, message, err := safeBuffer.GetNextMessage(1)
		assert.NoError(t, err)
		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
			Imsi: "111111",
		})
		assert.Equal(t, expectedMsg, message)
	}()
	time.Sleep(time.Millisecond * 500)
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)
	safeBuffer.WriteChunk(chunk)
	safeBuffer.WriteChunk(chunk)
	safeBuffer.WriteChunk(chunk)
}

/*
	TODO: Re-enable this test.
	This test has temporary been removed for 2 reasons:
		- Flackiness in CI
		- csfb implementation needs rework
*/

// func TestReadSuccessAndFailWithTimeout(t *testing.T) {
// 	safeBuffer, err := servicers.NewSafeBuffer()
// 	assert.NoError(t, err)
// 	signal := make(chan bool)
//
// 	go func() {
// 		messageType, message, err := safeBuffer.GetNextMessage(1)
// 		assert.NoError(t, err)
// 		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
// 		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
// 			Imsi: "111111",
// 		})
// 		assert.Equal(t, expectedMsg, message)
// 		signal <- true
// 	}()
// 	go func() {
// 		messageType, message, err := safeBuffer.GetNextMessage(1)
// 		assert.NoError(t, err)
// 		assert.Equal(t, decode.SGsAPIMSIDetachAck, messageType)
// 		expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
// 			Imsi: "111111",
// 		})
// 		assert.Equal(t, expectedMsg, message)
// 		signal <- true
// 	}()
// 	go func() {
// 		<-signal
// 		<-signal
// 		messageType, message, err := safeBuffer.GetNextMessage(1)
// 		assert.EqualError(t, err, "buffer read timeout")
// 		assert.Equal(t, decode.SGsMessageType(0x00), messageType)
// 		assert.Equal(t, &any.Any{}, message)
// 	}()
// 	imsi, _ := test_utils.ConstructIMSI("111111")
// 	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)
// 	safeBuffer.WriteChunk(chunk)
// 	safeBuffer.WriteChunk(chunk)
// }
