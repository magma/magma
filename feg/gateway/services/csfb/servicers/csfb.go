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

package servicers

import (
	"context"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/csfb/servicers/encode/message"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type PortNumber = int

type CsfbServer struct {
	Conn            ClientConnectionInterface
	ReceivingBuffer SafeBuffer
}

type ServerConnectionInterface interface {
	// Start a SCTP listener from server side
	StartListener(ipAddr string, port PortNumber) (PortNumber, error)
	// Close the active SCTP listener from server side
	CloseListener() error
	// Check if any connection is established
	ConnectionEstablished() bool
	// Accept connection from client
	AcceptConn() error
	// Close connection from client
	CloseConn() error
	// Receive data through the active listener in server side
	ReceiveThroughListener() ([]byte, error)
	// Send data through the established connection from server side
	SendFromServer([]byte) error
}

type ClientConnectionInterface interface {
	// Establish connection to a remote server through DialSCTP from client side
	EstablishConn() error
	// Close the established connection from client side
	CloseConn() error
	// Send data through the established connection from client side
	Send(message []byte) error
	// Receive data through the established connection in client side
	Receive() ([]byte, error)
}

func CreateVlrSCTPconnection(config *mconfig.CsfbConfig) (*SCTPClientConnection, error) {
	vlrSCTP, _ := convertIPAddressFromStrip(config.Client.ServerAddress)
	localSCTP, _ := convertIPAddressFromStrip(config.Client.LocalAddress)
	return NewSCTPClientConnection(vlrSCTP, localSCTP)
}

func NewCsfbServer(ConnectionInterface ClientConnectionInterface) (*CsfbServer, error) {
	return &CsfbServer{Conn: ConnectionInterface}, nil
}

// AlertAc sends SGsAP-ALERT-ACK to VLR
// to acknowledge a previous SGsAP-ALERT-REQUEST message
func (srv *CsfbServer) AlertAc(
	ctx context.Context,
	req *protos.AlertAck,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPAlertAck(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-ALERT-ACK: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// AlertRej sends SGsAP-ALERT-REJECT to VLR to indicate that the MME
// could not identify the IMSI indicated in the SGsAP-ALERT-REQUEST message
func (srv *CsfbServer) AlertRej(
	ctx context.Context,
	req *protos.AlertReject,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPAlertReject(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-ALERT-REJECT: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// EPSDetachInd sends SGsAP-EPS-DETACH-INDICATION to VLR
// to indicate an EPS detach performed from the UE or the MME
func (srv *CsfbServer) EPSDetachInd(
	ctx context.Context,
	req *protos.EPSDetachIndication,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPEPSDetachIndication(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-EPS-DETACH-INDICATION: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// IMSIDetachInd sends SGsAP-IMSI-DETACH-INDICATION to VLR
// to indicate an IMSI detach performed from the UE
func (srv *CsfbServer) IMSIDetachInd(
	ctx context.Context,
	req *protos.IMSIDetachIndication,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPIMSIDetachIndication(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-IMSI-DETACH-INDICATION: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// LocationUpdateReq sends SGsAP-LOCATION-UPDATE-REQUEST to VLR either
// to request update of its location file (normal update) or to request IMSI attach
func (srv *CsfbServer) LocationUpdateReq(
	ctx context.Context,
	req *protos.LocationUpdateRequest,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPLocationUpdateRequest(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-LOCATION-UPDATE-REQUEST: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// PagingRej sends SGsAP-PAGING-REJECT to VLR to indicate that
// the delivery of a previous SGsAP-PAGING-REQUEST message has failed
func (srv *CsfbServer) PagingRej(
	ctx context.Context,
	req *protos.PagingReject,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPPagingReject(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-PAGING-REJECT: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// ServiceReq sends SGsAP-SERVICE-REQUEST to VLR as a response
// to a previously received SGsAP-PAGING-REQUEST message
// to indicate the existence of a NAS signaling Connection
// between the UE and the MME or to indicate to the VLR that
// the NAS signaling Connection has been established after the paging procedure
func (srv *CsfbServer) ServiceReq(
	ctx context.Context,
	req *protos.ServiceRequest,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPServiceRequest(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-SERVICE-REQUEST: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// TMSIReallocationComp sends SGsAP-TMSI-REALLOCATION-COMPLETE to VLR
// to indicate that TMSI reallocation on the UE has been successfully completed
func (srv *CsfbServer) TMSIReallocationComp(
	ctx context.Context,
	req *protos.TMSIReallocationComplete,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPTMSIReallocationComplete(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-TMSI-REALLOCATION-COMPLETE: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// UEActivityInd sends SGsAP-UE-ACTIVITY-INDICATION to VLR
// to indicate that activity from a UE has been detected
func (srv *CsfbServer) UEActivityInd(
	ctx context.Context,
	req *protos.UEActivityIndication,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPUEActivityIndication(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-UE-ACTIVITY-INDICATION: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// UEUnreach sends SGsAP-UE-UNREACHABLE to VLR to indicate that,
// for example, paging could not be performed
// because the UE is marked as unreachable at the MME
func (srv *CsfbServer) UEUnreach(
	ctx context.Context,
	req *protos.UEUnreachable,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPUEUnreachable(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-UE-UNREACHABLE: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// Uplink sends SGsAP-UPLINK-UNITDATA to VLR
// to transparently convey a NAS message, from the UE, to the VLR
func (srv *CsfbServer) Uplink(
	ctx context.Context,
	req *protos.UplinkUnitdata,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPUplinkUnitdata(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-UPLINK-UNITDATA: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// MMEResetAck sends SGsAP-RESET-ACK to VLR to acknowledge
// a previous SGsAP-RESET-INDICATION message. This message indicates that
// all the SGs associations to the VLR or the MME have been marked as invalid.
func (srv *CsfbServer) MMEResetAck(
	ctx context.Context,
	req *protos.ResetAck,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPResetAck(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-RESET-ACK: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// MMEResetIndication sends SGsAP-RESET-INDICATION to VLR
// to indicate that a failure in the MME has occurred
// and all the SGs associations to the MME are be marked as invalid.
func (srv *CsfbServer) MMEResetIndication(
	ctx context.Context,
	req *protos.ResetIndication,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPResetIndication(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-RESET-INDICATION: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// MMEStatus sends SGsAP-STATUS to VLR to indicate an error
func (srv *CsfbServer) MMEStatus(
	ctx context.Context,
	req *protos.Status,
) (*orcprotos.Void, error) {
	encodedMsg, err := message.EncodeSGsAPStatus(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-STATUS: %s", err)
		return &orcprotos.Void{}, err
	}
	return &orcprotos.Void{}, srv.Conn.Send(encodedMsg)
}

// SendResetAck sends SGsAP-RESET-ACK to VLR
// Different from the MMEResetAck invoked by the gateway through GRPC,
// SendResetAck is invoked in the FeG as soon as the SGsAP-RESET-INDICATION
// is received and decoded.
func (srv *CsfbServer) SendResetAck() error {
	req, err := constructResetAck()
	if err != nil {
		glog.Errorf("Failed to construct SGsAP-RESET-ACK: %s", err)
		return err
	}
	encodedMsg, err := message.EncodeSGsAPResetAck(req)
	if err != nil {
		glog.Errorf("Failed to encode SGsAP-RESET-ACK: %s", err)
		return err
	}
	return srv.Conn.Send(encodedMsg)
}

func constructResetAck() (*protos.ResetAck, error) {
	mmeName, err := ConstructMMEName()
	if err != nil {
		glog.Errorf("Failed to construct MME name: %s", err)
		return nil, err
	}
	return &protos.ResetAck{MmeName: mmeName}, nil
}
