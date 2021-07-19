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

package mock_pgw

import (
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"

	"magma/feg/cloud/go/protos"
)

func (mPgw *MockPgw) getHandleCreateSessionRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a CreateSessionRequest")

		session := gtpv2.NewSession(sgwAddr, &gtpv2.Subscriber{Location: &gtpv2.Location{}})
		bearer := session.GetDefaultBearer()

		var err error
		csReqFromSGW := msg.(*message.CreateSessionRequest)
		if imsiIE := csReqFromSGW.IMSI; imsiIE != nil {
			imsi, err2 := imsiIE.IMSI()
			if err2 != nil {
				return err2
			}
			session.IMSI = imsi

			// remove previous session for the same subscriber if exists.
			sess, err2 := c.GetSessionByIMSI(imsi)
			if err2 != nil {
				switch err2.(type) {
				case *gtpv2.UnknownIMSIError:
					// whole new session. just ignore.
				default:
					return errors.Wrap(err2, "got something unexpected")
				}
			} else {
				fmt.Printf("Existing IMSI during Create Session Request on PGW (%s). Deleting previous session\n", imsi)
				c.RemoveSession(sess)
			}
		} else {
			fmt.Println("Missing IE (IMSI) on Create Session Request that PGW received")
			return &gtpv2.RequiredIEMissingError{Type: ie.IMSI}
		}
		if uliIE := csReqFromSGW.ULI; uliIE != nil {
			mPgw.LastULI, err = uliIE.UserLocationInformation()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.UserLocationInformation}
		}
		if msisdnIE := csReqFromSGW.MSISDN; msisdnIE != nil {
			session.MSISDN, err = msisdnIE.MSISDN()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.MSISDN}
		}
		if meiIE := csReqFromSGW.MEI; meiIE != nil {
			session.IMEI, err = meiIE.MobileEquipmentIdentity()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.MobileEquipmentIdentity}
		}
		if apnIE := csReqFromSGW.APN; apnIE != nil {
			bearer.APN, err = apnIE.AccessPointName()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.AccessPointName}
		}
		if netIE := csReqFromSGW.ServingNetwork; netIE != nil {
			session.MNC, err = netIE.MNC()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.ServingNetwork}
		}
		if ratIE := csReqFromSGW.RATType; ratIE != nil {
			session.RATType, err = ratIE.RATType()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.RATType}
		}
		if sgwTEID := csReqFromSGW.SenderFTEIDC; sgwTEID != nil {
			teid, err := sgwTEID.TEID()
			if err != nil {
				return err
			}
			session.AddTEID(gtpv2.IFTypeS5S8SGWGTPC, teid)
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		if brCtxIE := csReqFromSGW.BearerContextsToBeCreated; brCtxIE != nil {
			for _, childIE := range brCtxIE.ChildIEs {
				switch childIE.Type {
				case ie.EPSBearerID:
					bearer.EBI, err = childIE.EPSBearerID()
					if err != nil {
						return err
					}
				case ie.FullyQualifiedTEID:
					it, err := childIE.InterfaceType()
					if err != nil {
						return err
					}
					// only used for user plane
					teidOut, err := childIE.TEID()
					if err != nil {
						return err
					}
					session.AddTEID(it, teidOut)
				case ie.BearerQoS:
					err = handleQOStoBearer(childIE, bearer)
					if err != nil {
						return err
					}
					// save for testing purposes
					mPgw.LastQos, err = handleQOStoProto(childIE)
					if err != nil {
						return err
					}

				}
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
		}

		if paaIE := csReqFromSGW.PAA; paaIE != nil {
			bearer.SubscriberIP = paaIE.MustIP().String()
		}

		// FTEIDS and TEIDS
		// create PGW control plane FTeids
		cIP := strings.Split(c.LocalAddr().String(), ":")[0]
		var pgwFTEIDc *ie.IE
		if mPgw.CreateSessionOptions.PgwTEIDc != 0 {
			// use passed options value
			pgwFTEIDc = ie.NewFullyQualifiedTEID(
				gtpv2.IFTypeS5S8PGWGTPC, mPgw.CreateSessionOptions.PgwTEIDc, cIP, "").WithInstance(1)
		} else {
			pgwFTEIDc = c.NewSenderFTEID(cIP, "").WithInstance(1)
		}

		// create PGW user plane FTeids
		uIP := strings.Split(dummyUserPlanePgwIP, ":")[0]
		var pgwUteid uint32
		if mPgw.CreateSessionOptions.PgwTEIDu != 0 {
			// use passed options value
			pgwUteid = mPgw.CreateSessionOptions.PgwTEIDu
		} else {
			mPgw.randGenMux.Lock()
			pgwUteid = (rand.Uint32() / 1000) * 1000 // for easy identification, this teid will always end in 000
			mPgw.randGenMux.Unlock()
		}
		pgwFTEIDu := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8PGWGTPU, pgwUteid, uIP, "").WithInstance(2)

		// get SGW user plane Teid
		var sgwTEIDc uint32
		if mPgw.CreateSessionOptions.SgwTEIDc != 0 {
			// used passed options value
			sgwTEIDc = mPgw.CreateSessionOptions.SgwTEIDc
		} else {
			// get the teid received by the request and stored in the seession previously
			sgwTEIDc, err = session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
			if err != nil {
				return err
			}
		}

		qosPCI := uint8(0)
		if bearer.QoSProfile.PCI {
			qosPCI = 1
		}
		qosPVI := uint8(0)
		if bearer.QoSProfile.PVI {
			qosPVI = 1
		}

		// Protocol Configuration Options (nil if not existent)
		pco := csReqFromSGW.PCO

		// send
		csRspFromPGW := message.NewCreateSessionResponse(
			sgwTEIDc, msg.Sequence(),
			ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
			pgwFTEIDc,
			ie.NewPDNAddressAllocation(bearer.SubscriberIP),
			ie.NewAPNRestriction(gtpv2.APNRestrictionPublic2),
			pco,
			ie.NewBearerContext(
				ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
				ie.NewEPSBearerID(bearer.EBI),
				pgwFTEIDu,
				ie.NewChargingID(bearer.ChargingID),
				ie.NewBearerQoS(qosPCI, bearer.QoSProfile.PL, qosPVI,
					bearer.QoSProfile.QCI, bearer.QoSProfile.MBRUL, bearer.QoSProfile.MBRDL,
					bearer.QoSProfile.GBRUL, bearer.QoSProfile.GBRDL),
			))

		session.AddTEID(gtpv2.IFTypeS5S8PGWGTPC, pgwFTEIDc.MustTEID())
		session.AddTEID(gtpv2.IFTypeS5S8PGWGTPU, pgwFTEIDu.MustTEID())

		s5pgwTEID, err := session.GetTEID(gtpv2.IFTypeS5S8PGWGTPC)
		if err != nil {
			return err
		}
		c.RegisterSession(s5pgwTEID, session)
		if err := session.Activate(); err != nil {
			return err
		}

		// save values given for testing purposes
		mPgw.LastTEIDc, err = pgwFTEIDc.TEID()
		if err != nil {
			return err
		}
		mPgw.LastTEIDu, err = pgwFTEIDu.TEID()
		if err != nil {
			return err
		}
		if err := c.RespondTo(sgwAddr, csReqFromSGW, csRspFromPGW); err != nil {
			fmt.Printf("mock PGW couldnt create a session for %s\n", session.IMSI)
			c.RemoveSession(session)
			return err
		}
		fmt.Printf("mock PGW created a session for: %s\n", session.IMSI)
		return nil
	}
}

func getRandomIp() string {
	return fmt.Sprintf("192.168.1.%d", (1 + rand.Intn(250)))
}

func handleQOStoBearer(qosIE *ie.IE, br *gtpv2.Bearer) error {
	var err error
	br.PL, err = qosIE.PriorityLevel()
	if err != nil {
		return err
	}
	br.QCI, err = qosIE.QCILabel()
	if err != nil {
		return err
	}
	br.PCI = qosIE.HasPCI()
	br.PVI = qosIE.HasPVI()

	br.MBRUL, err = qosIE.MBRForUplink()
	if err != nil {
		return err
	}
	br.MBRDL, err = qosIE.MBRForDownlink()
	if err != nil {
		return err
	}
	br.GBRUL, err = qosIE.GBRForUplink()
	if err != nil {
		return err
	}
	br.GBRDL, err = qosIE.GBRForDownlink()
	if err != nil {
		return err
	}
	return nil
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
	if qosIE.PreemptionCapability() {
		qos.PreemptionCapability = 1
	}

	// Preemption Vulnerability
	if qosIE.PreemptionVulnerability() {
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

func createQosIE(qp *gtpv2.QoSProfile) *ie.IE {
	pci, pvi := 0, 0
	if qp.PCI {
		pci = 1
	}
	if qp.PVI {
		pvi = 1
	}
	qosIE := ie.NewBearerQoS(uint8(pci), qp.PL, uint8(pvi),
		qp.QCI, qp.MBRUL, qp.MBRDL, qp.GBRUL, qp.GBRDL)
	return qosIE

}

// getHandleCreateSessionRequestWithErrorCause Responds with an arbitrary error cause
func (mPgw *MockPgw) getHandleCreateSessionRequestWithErrorCause(errorCause uint8) gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a CreateSessionRequest, but returning ERROR")
		csReqFromSGW := msg.(*message.CreateSessionRequest)
		sgwTEID := csReqFromSGW.SenderFTEIDC
		if sgwTEID != nil {
			_, err := sgwTEID.TEID()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		// send
		csRspFromPGW := message.NewCreateSessionResponse(
			sgwTEID.MustTEID(), msg.Sequence(),
			ie.NewCause(errorCause, 0, 0, 0, nil),
		)

		if err := c.RespondTo(sgwAddr, csReqFromSGW, csRspFromPGW); err != nil {
			return err
		}

		return nil
	}
}

// getHandleCreateSessionResponseWithMissingIE responds with a CreateSessionResponse that has a missing field
func (mPgw *MockPgw) getHandleCreateSessionResponseWithMissingIE() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a CreateSessionRequest, but returning ERROR")
		csReqFromSGW := msg.(*message.CreateSessionRequest)
		sgwTEID := csReqFromSGW.SenderFTEIDC
		if sgwTEID != nil {
			_, err := sgwTEID.TEID()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		// Mising pgwFTEID and bearer pgwFTEIDu
		csRspFromPGW := message.NewCreateSessionResponse(
			sgwTEID.MustTEID(), msg.Sequence(),
			ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
			//pgwFTEIDc,
			ie.NewPDNAddressAllocation("10.1.2.3"),
			ie.NewAPNRestriction(gtpv2.APNRestrictionPublic2),
			ie.NewBearerContext(
				ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
				ie.NewEPSBearerID(5),
			))

		if err := c.RespondTo(sgwAddr, csReqFromSGW, csRspFromPGW); err != nil {
			return err
		}
		return nil
	}
}

// getHandleCreateSessionRequestWithMissingIE responds with a CauseMandatoryIEMissing
func (mPgw *MockPgw) getHandleCreateSessionRequestWithMissingIE(missingIE *ie.IE) gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a CreateSessionRequest, but returning ERROR")
		csReqFromSGW := msg.(*message.CreateSessionRequest)
		sgwTEID := csReqFromSGW.SenderFTEIDC
		if sgwTEID != nil {
			_, err := sgwTEID.TEID()
			if err != nil {
				return err
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		// Mising pgwFTEID and bearer pgwFTEIDu
		csRspFromPGW := message.NewCreateSessionResponse(
			sgwTEID.MustTEID(), msg.Sequence(),
			ie.NewCause(gtpv2.CauseMandatoryIEMissing, 0, 0, 0, missingIE),
		)

		if err := c.RespondTo(sgwAddr, csReqFromSGW, csRspFromPGW); err != nil {
			return err
		}
		return nil
	}
}
