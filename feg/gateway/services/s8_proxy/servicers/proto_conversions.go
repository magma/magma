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
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"

	"magma/feg/cloud/go/protos"
)

// parseCreateSessionResponse parses a gtp message into a CreateSessionResponsePgw. In case
// the message is proper it also returns the session. In case it there is an error it returns
// the cause of error
func parseCreateSessionResponse(msg message.Message) (csRes *protos.CreateSessionResponsePgw, err error) {
	//glog.V(4).Infof("Received Create Session Response (gtp):\n%s", message.Prettify(msg))
	csResGtp := msg.(*message.CreateSessionResponse)
	glog.V(2).Infof("Received Create Session Response (gtp):\n%s", csResGtp.String())
	csRes = &protos.CreateSessionResponsePgw{}
	// check Cause value first.
	if causeIE := csResGtp.Cause; causeIE != nil {
		csRes.GtpError, err = handleCause(causeIE, msg)
		if err != nil || csRes.GtpError != nil {
			// return either GtpError or err
			return csRes, err
		}
		// If we get here, the message will be processed
	} else {
		csRes.GtpError = errorIeMissing(ie.Cause)
		return csRes, nil
	}

	// get C AGW Teid (is the same used on Create Session Request by MME)
	csRes.CAgwTeid = msg.TEID()

	// get PDN Allocation from PGW
	if paaIE := csResGtp.PAA; paaIE != nil {
		csRes.Paa, csRes.PdnType, err = handlePDNAddressAllocation(paaIE)
		if err != nil {
			return nil, err
		}
	} else {
		csRes.GtpError = errorIeMissing(ie.PDNAddressAllocation)
		// error is in GtpError
		return csRes, nil
	}

	// Pgw control plane fteid
	if pgwCFteidIE := csResGtp.PGWS5S8FTEIDC; pgwCFteidIE != nil {
		csRes.CPgwFteid, _, err = handleFTEID(pgwCFteidIE)
		if err != nil {
			err = fmt.Errorf("Couldn't get PGW control plane FTEID: %s ", err)
			return nil, err
		}
	} else {
		csRes.GtpError = errorIeMissing(ie.FullyQualifiedTEID)
		// error is in GtpError
		return csRes, nil
	}

	// Protocol Configuration Options (PCO) optional
	if pgwPcoIE := csResGtp.PCO; pgwPcoIE != nil {
		csRes.ProtocolConfigurationOptions, err = handlePCO(pgwPcoIE)
		if err != nil {
			err = fmt.Errorf("Couldn't get Protocol Configuration Options: %s ", err)
			return nil, err
		}
	}

	// TODO: handle more than one bearer
	if brCtxIE := csResGtp.BearerContextsCreated; brCtxIE != nil {
		csRes.BearerContext, csRes.GtpError, err = handleBearerCtx(brCtxIE)
		if err != nil {
			return nil, err
		}
		if csRes.GtpError != nil {
			return csRes, nil
		}
	} else {
		csRes.GtpError = errorIeMissing(ie.BearerContext)
		// error is in GtpError
		return csRes, nil
	}
	return csRes, nil
}

// parseDeleteSessionResponse parses a gtp message into a DeleteSessionResponsePgw. In case
// the message is proper it also returns the session. In case it there is an error it returns
// the cause of error
func parseDeleteSessionResponse(msg message.Message) (
	*protos.DeleteSessionResponsePgw, error) {
	//glog.V(4).Infof("Received Delete Session Response (gtp):\n%s", (msg))
	cdResGtp := msg.(*message.DeleteSessionResponse)
	glog.V(2).Infof("Received Delete Session Response (gtp):\n%s", cdResGtp.String())

	dsRes := &protos.DeleteSessionResponsePgw{}
	var err error
	// check Cause value first.
	if causeIE := cdResGtp.Cause; causeIE != nil {

		dsRes.GtpError, err = handleCause(causeIE, msg)
		if err != nil || dsRes.GtpError != nil {
			// return either GtpError or err
			return dsRes, err
		}
		// If we get here, the message will be processed
	} else {
		dsRes.GtpError = errorIeMissing(ie.Cause)
		// error is in GtpError
		return dsRes, nil
	}
	return dsRes, nil
}

func parseCreateBearerRequest(msg message.Message) (*protos.CreateBearerRequestPgw, *protos.GtpError, error) {
	cbReqGtp := msg.(*message.CreateBearerRequest)
	glog.V(2).Infof("Received Create Bearer Request (gtp):\n%s", cbReqGtp.String())

	cbRes := &protos.CreateBearerRequestPgw{}

	// cgw control plane teid
	if !cbReqGtp.HasTEID() {
		return nil, errorIeMissing(ie.FullyQualifiedTEID), nil
	}
	cbRes.CAgwTeid = cbReqGtp.TEID()

	if linkedEBI := cbReqGtp.LinkedEBI; linkedEBI != nil {
		cbRes.LinkedBearerId = uint32(linkedEBI.MustEPSBearerID())
	} else {
		return nil, errorIeMissing(ie.EPSBearerID), nil
	}

	// TODO: handle more than one bearer
	if brCtxIE := cbReqGtp.BearerContexts; brCtxIE != nil {
		bearerContext, gtpError, err := handleBearerCtx(brCtxIE)
		if err != nil || gtpError != nil {
			return nil, gtpError, err
		}
		cbRes.BearerContext = bearerContext
	} else {
		return nil, errorIeMissing(ie.BearerContext), nil
	}

	// Protocol Configuration Options (PCO) optional
	if pgwPcoIE := cbReqGtp.PCO; pgwPcoIE != nil {
		pco, err := handlePCO(pgwPcoIE)
		if err != nil {
			err = fmt.Errorf("Couldn't get Protocol Configuration Options: %s ", err)
			return nil, nil, err
		}
		cbRes.ProtocolConfigurationOptions = pco
	}

	return cbRes, nil, nil
}

func handleCause(causeIE *ie.IE, msg message.Message) (*protos.GtpError, error) {
	cause, err := causeIE.Cause()
	if err != nil {
		return nil, fmt.Errorf("Couldn't check cause of %s: %s", msg.MessageTypeName(), err)
	}

	switch cause {
	case gtpv2.CauseRequestAccepted:
		return nil, nil
	default:
		gtpErrorString := fmt.Sprintf("%s with sequence # %d not accepted. Cause: %d", msg.MessageTypeName(), msg.Sequence(), cause)
		offendingIE, _ := causeIE.OffendingIE()
		if offendingIE != nil {
			gtpErrorString = fmt.Sprintf("%s %s: %d %s", gtpErrorString, " With Offending IE", offendingIE.Type, offendingIE)
		}
		glog.Warning(gtpErrorString)
		return &protos.GtpError{
			Cause: uint32(cause),
			Msg:   gtpErrorString,
		}, nil
	}
}

// handleFTEID converts FTEID IE format into Proto format returning also the type of interface
func handleFTEID(fteidIE *ie.IE) (*protos.Fteid, uint8, error) {
	interfaceType, err := fteidIE.InterfaceType()
	if err != nil {
		return nil, interfaceType, err
	}
	teid, err := fteidIE.TEID()
	if err != nil {
		return nil, interfaceType, err
	}

	fteid := &protos.Fteid{Teid: teid}
	if !fteidIE.HasIPv4() && !fteidIE.HasIPv6() {
		return nil, interfaceType, fmt.Errorf("Error: fteid %+v has no ips", fteidIE.String())
	}
	if fteidIE.HasIPv4() {
		ipv4, err := fteidIE.IPv4()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv4Address = ipv4.String()
	}
	if fteidIE.HasIPv6() {
		ipv6, err := fteidIE.IPv6()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv6Address = ipv6.String()
	}
	return fteid, interfaceType, nil
}

func handlePDNAddressAllocation(paaIE *ie.IE) (*protos.PdnAddressAllocation, protos.PDNType, error) {
	pdnTypeIE, err := paaIE.PDNType()
	if err != nil {
		return nil, 0, err
	}
	var pdnType protos.PDNType
	var paa protos.PdnAddressAllocation
	switch pdnTypeIE {
	case gtpv2.PDNTypeIPv4:
		pdnType = protos.PDNType_IPV4
		paa.Ipv4Address = paaIE.MustIPv4().String()
	case gtpv2.PDNTypeIPv6:
		pdnType = protos.PDNType_IPV6
		paa.Ipv6Address = paaIE.MustIPv6().String()
	case gtpv2.PDNTypeIPv4v6:
		pdnType = protos.PDNType_IPV4V6
		paa.Ipv6Address = paaIE.MustIPv6().String()
		paa.Ipv6Address = paaIE.MustIPv6().String()
	case gtpv2.PDNTypeNonIP:
		pdnType = protos.PDNType_NonIP
	}
	return &paa, pdnType, nil
}

func handleQOStoProto(qosIE *ie.IE) (*protos.QosInformation, error) {
	qos := &protos.QosInformation{}

	// priority level
	pl, err := qosIE.PriorityLevel()
	if err != nil {
		return nil, err
	}
	qos.PriorityLevel = uint32(pl)

	// qci label
	qci, err := qosIE.QCILabel()
	if err != nil {
		return nil, err
	}
	qos.Qci = uint32(qci)

	// Preemption Capability
	if qosIE.HasPCI() {
		qos.PreemptionCapability = 1
	}

	// Preemption Vulnerability
	if qosIE.HasPVI() {
		qos.PreemptionVulnerability = 1
	}

	// maximum bitrate
	mAmbr := &protos.Ambr{}
	mAmbr.BrUl, err = qosIE.MBRForUplink()
	if err != nil {
		return nil, err
	}
	mAmbr.BrDl, err = qosIE.MBRForDownlink()
	if err != nil {
		return nil, err
	}
	qos.Mbr = mAmbr

	// granted bitrate
	gAmbr := &protos.Ambr{}
	gAmbr.BrUl, err = qosIE.GBRForUplink()
	if err != nil {
		return nil, err
	}
	gAmbr.BrDl, err = qosIE.GBRForDownlink()
	if err != nil {
		return nil, err
	}
	qos.Gbr = gAmbr

	return qos, nil
}

func handleBearerCtx(brCtxIE *ie.IE) (*protos.BearerContext, *protos.GtpError, error) {
	bearerCtx := &protos.BearerContext{}
	for _, childIE := range brCtxIE.ChildIEs {
		switch childIE.Type {
		case ie.Cause:
			cause, err := childIE.Cause()
			if err != nil {
				return nil, nil, err
			}
			if cause != gtpv2.CauseRequestAccepted {
				gtpError := &protos.GtpError{
					Cause: uint32(cause),
					Msg:   "Bearer could not be created",
				}
				// error is in GtpError
				return bearerCtx, gtpError, nil
			}

		case ie.EPSBearerID:
			ebi, err := childIE.EPSBearerID()
			if err != nil {
				return nil, nil, err
			}
			bearerCtx.Id = uint32(ebi)

		case ie.FullyQualifiedTEID:
			userPlaneFteid, _, err := handleFTEID(childIE)
			if err != nil {
				return nil, nil, err
			}
			bearerCtx.UserPlaneFteid = userPlaneFteid

		case ie.ChargingID:
			chargingId, err := childIE.ChargingID()
			if err != nil {
				return nil, nil, err
			}
			bearerCtx.ChargingId = chargingId

		case ie.BearerQoS:
			// save for testing purposes
			qos, err := handleQOStoProto(childIE)
			if err != nil {
				return nil, nil, err
			}
			bearerCtx.Qos = qos

		case ie.BearerTFT:
			bearerTFT, err := handleTFT(childIE)
			if err != nil {
				return nil, nil, err
			}
			bearerCtx.Tft = bearerTFT
		}
	}
	return bearerCtx, nil, nil
}

func handlePCO(pcoIE *ie.IE) (*protos.ProtocolConfigurationOptions, error) {
	pgwPcoField, err := pcoIE.ProtocolConfigurationOptions()
	if err != nil {
		err = fmt.Errorf("Couldn't get PGW  Protocol Configuration Options: %s ", err)
		return nil, err
	}
	var containers []*protos.PcoProtocolOrContainerId
	for _, containerField := range pgwPcoField.ProtocolOrContainers {
		containers = append(containers,
			&protos.PcoProtocolOrContainerId{
				Id:       uint32(containerField.ID),
				Contents: containerField.Contents,
			})
	}
	return &protos.ProtocolConfigurationOptions{
		ConfigProtocol:     uint32(pgwPcoField.ConfigurationProtocol),
		ProtoOrContainerId: containers,
	}, nil
}

func errorIeMissing(missingIE uint8) *protos.GtpError {
	errMsg := gtpv2.RequiredIEMissingError{Type: missingIE}
	return &protos.GtpError{
		Cause: uint32(gtpv2.CauseMandatoryIEMissing),
		Msg:   errMsg.Error(),
	}
}

func ip2Long(ip string) uint32 {
	var long uint32
	addrs := net.ParseIP(ip).To4()

	binary.Read(bytes.NewBuffer(addrs), binary.BigEndian, &long)
	return long
}

func uintToIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}
