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

// package servce implements EAP-SIM GRPC service
package sim

import (
	"fmt"

	"magma/feg/gateway/services/aaa/protos"
	aaa "magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim/metrics"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewSIMNotificationReq(identifier uint8, code uint16) eap.Packet {
	metrics.FailureNotifications.Inc()
	return []byte{
		eap.RequestCode,
		identifier,
		0, 12, // EAP Len
		TYPE,
		byte(SubtypeNotification),
		0, 0,
		byte(AT_NOTIFICATION),
		1, // EAP SIM Attr Len
		uint8(code >> 8), uint8(code)}
}

func EapErrorResPacket(id uint8, code uint16, rpcCode codes.Code, f string, a ...interface{}) (eap.Packet, error) {
	Errorf(rpcCode, f, a...) // log only
	return NewSIMNotificationReq(id, code), nil
}

func EapErrorResPacketWithMac(id uint8, code uint16, K_aut []byte, rpcCode codes.Code, f string, a ...interface{}) (eap.Packet, error) {
	p := NewSIMNotificationReq(id, code)
	p, err := AppendMac(p, K_aut)
	if err != nil {
		panic(err) // should never happen
	}
	Errorf(rpcCode, f, a...) // log only
	return p, nil
}

func EapErrorRes(
	id uint8, code uint16,
	rpcCode codes.Code,
	ctx *aaa.Context,
	f string, a ...interface{}) (*protos.Eap, error) {

	Errorf(rpcCode, f, a...) // log only
	return &protos.Eap{Payload: NewSIMNotificationReq(id, code), Ctx: ctx}, nil
}

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	glog.Errorf("SIM RPC [%s] %s", code, msg)
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	glog.Errorf("SIM RPC [%s] %s", code, err)
	return status.Error(code, err.Error())
}
