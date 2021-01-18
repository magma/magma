#!/usr/bin/env python3

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse

from lte.protos import session_manager_pb2
from lte.protos.subscriberdb_pb2 import (
     SubscriberID,
)
from lte.protos.session_manager_pb2 import (
    SetSMSessionContext,
    M5GSMSessionContext,
    RedirectServer,
    RequestType,
    priorityaccess,
    AccessType,
    CommonSessionContext,
    SMSessionFSMState,
    RatSpecificContext,
    PduSessionType,
    SscMode,
    RATType,
)
from lte.protos.session_manager_pb2_grpc import (
    AmfPduSessionSmContextStub,
    SessionProxyResponderStub,

)
from magma.common.rpc_utils import grpc_wrapper
from magma.subscriberdb.sid import SIDUtils

class CreateAmfSession(object):

	def __init__(self,rat_type = RATType.Name(2),sm_session_state=SMSessionFSMState.Name(0),sm_session_version = [0,'uint32'],pdu_session_id = bytes([2]),rquest_type = RequestType.Name(0),access_type=AccessType.Name(0),pdu_address=RedirectServer(redirect_address_type=RedirectServer.IPV4,redirect_server_address="10.20.35.45"),pdu_session_type=PduSessionType.Name(0),ssc_mode=SscMode.Name(2)):
		self._set_session = \
		       SetSMSessionContext(common_context=CommonSessionContext(sid = SubscriberID(id="imsi00000000002"),apn = bytes("HYD",'utf-8'),rat_type = RATType.Name(2),
                       sm_session_state=SMSessionFSMState.Name(0)),\
                       rat_specific_context = RatSpecificContext(m5gsm_session_context = M5GSMSessionContext(pdu_session_id = bytes([2]),\
                       rquest_type = RequestType.Name(0),\
                       pdu_address=RedirectServer(redirect_address_type=RedirectServer.IPV4, redirect_server_address="10.20.35.45"),access_type=AccessType.Name(0),
                       pdu_session_type=PduSessionType.Name(0),ssc_mode=SscMode.Name(2))))
@grpc_wrapper
def set_amf_session (client,args):

     cls_sess = CreateAmfSession (args.rat_type,args.sm_session_state,args.sm_session_version,args.pdu_session_id,args.rquest_type,args.pdu_address,args.access_type,args.pdu_session_type,args.ssc_mode)
     print (cls_sess._set_session)
     response = client.SetAmfSessionContext(cls_sess._set_session)
     print (response)

def create_amf_parser(apps):
    """
    Creates the argparse subparser for the ng_services app
    """

    app = apps.add_parser('amf_context')
    subparsers = app.add_subparsers(title='subcommands',dest='cmd')
    subcmd = subparsers.add_parser('set_amf_session',help='AMF Set Session')
    subcmd.add_argument('--sid', help='Subscriber_ID', default = "imsi00000000002")
    subcmd.add_argument('--apn',help='HYD')
    subcmd.add_argument('--rat_type',help='0-tgpp_lte,1-tgpp_wlan,2-tgpp_nr',default='2')
    subcmd.add_argument('--sm_session_state',help='0-creating_0,1-creating_1',default='0')
    subcmd.add_argument('--sm_session_version',help='SM_Session_Version',default='0')
    subcmd.add_argument('--pdu_session_id',help='PDU session ID',default='0x2')
    subcmd.add_argument('--rquest_type',help='0 - initial_request, 1 - existing_pdu_session, 2 - initial_emergency_request')
    subcmd.add_argument('--pdu_address',help='0 - ipv4,1 - ipv6,2 - url,3 - sip_uri',default='1')
    subcmd.add_argument('--pdu_addresse',help='pdu_address_ip',default="10.20.35.45")
    subcmd.add_argument('--access_type',help="0-m_3gpp_access_3gpp,1-non_3gpp_acess",default='0')
    subcmd.add_argument('--pdu_session_type',help='0 - ipv4,1 - ipv6,2-ipv4ipv6,3-unstructred',default='1')
    subcmd.add_argument('--ssc_mode',help='0 - ssc_mode_1,1 - ssc_mode2, 2 - ssc_mode3',default='2')
    subcmd.set_defaults(func=set_amf_session)

def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for sessiond',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    apps = parser.add_subparsers(title='apps', dest='cmd')
    create_amf_parser(apps)
    return parser
def main():
    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)
    # Execute the subcommand function
    args.func(args,AmfPduSessionSmContextStub , 'sessiond')
if __name__ == "__main__":
    main()



