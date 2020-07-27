/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package eap (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
//
//go:generate protoc -I. -I ../aaa/protos --go_out=plugins=grpc,paths=source_relative:. protos/eap_auth.proto
//
package eap

import "magma/feg/gateway/services/aaa/protos"

const (
	// EAP Message Payload Offsets
	EapMsgCode int = iota
	EapMsgIdentifier
	EapMsgLenHigh
	EapMsgLenLow
	EapMsgMethodType
	EapMsgData
	EapReserved1
	EapReserved2
	EapFirstAttribute
	EapFirstAttributeLen
)

const (
	// EapSubtype - pseudonym for EapMsgData
	EapSubtype   = EapMsgData
	EapHeaderLen = EapMsgMethodType
	// EapMaxLen maximum possible length of EAP Packet
	EapMaxLen uint = 1<<16 - 1

	UndefinedCode = 0
	RequestCode   = 1
	ResponseCode  = 2
	SuccessCode   = 3
	FailureCode   = 4
)

const (
	// EAP Related Consts
	MethodIdentity = uint8(protos.EapType_Identity)
	MethodNak      = uint8(protos.EapType_Legacy_Nak)
	CodeResponse   = uint8(protos.EapCode_Response)
)
