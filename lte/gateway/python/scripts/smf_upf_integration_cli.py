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

from lte.protos.session_manager_pb2 import (
    CommonSessionContext,
    M5GSMSessionContext,
    PduSessionType,
    RatSpecificContext,
    RatSpecificNotification,
    RATType,
    RequestType,
    SetSmNotificationContext,
    SetSMSessionContext,
    SMSessionFSMState,
    SscMode,
    TeidSet,
)
from lte.protos.session_manager_pb2_grpc import AmfPduSessionSmContextStub
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.common.rpc_utils import grpc_wrapper


class CreateAmfSession(object):

    def __init__(self):
        self._set_session = SetSMSessionContext(
            common_context=CommonSessionContext(
                sid=SubscriberID(id="IMSI12345"),
                ue_ipv4="192.168.128.11",
                apn=bytes("BLR", 'utf-8'),
                rat_type=RATType.Name(2),
                sm_session_state=SMSessionFSMState.Name(0),
                sm_session_version=0,
            ),
            rat_specific_context=RatSpecificContext(
                m5gsm_session_context=M5GSMSessionContext(
                    pdu_session_id=1,
                    request_type=RequestType.Name(
                        0,
                    ),
                    gnode_endpoint=TeidSet(
                        teid=10000,
                        end_ipv4_addr="192.168.60.141",
                    ),
                    pdu_session_type=PduSessionType.Name(0),
                    ssc_mode=SscMode.Name(2),
                ),
            ),
        )


class CreateAmfSessionIPv6(object):
    """
    Creates IPv6 PDU session
    """

    def __init__(self):
        self._set_session = SetSMSessionContext(
            common_context=CommonSessionContext(
                sid=SubscriberID(id="IMSI12345"),
                apn=bytes("BLR", 'utf-8'),
                ue_ipv6="2001:db8::1",
                rat_type=RATType.Name(2),
                sm_session_state=SMSessionFSMState.Name(0),
                sm_session_version=0,
            ),
            rat_specific_context=RatSpecificContext(
                m5gsm_session_context=M5GSMSessionContext(
                    pdu_session_id=1,
                    request_type=RequestType.Name(
                        0,
                    ),
                    gnode_endpoint=TeidSet(
                        teid=10000,
                        end_ipv4_addr="192.168.60.141",
                    ),
                    pdu_session_type=PduSessionType.Name(1),
                    ssc_mode=SscMode.Name(2),
                ),
            ),
        )


class CreateAmfSessionIPv4v6(object):
    """
    Creates IPv4v6 PDU session
    """

    def __init__(self):
        self._set_session = SetSMSessionContext(
            common_context=CommonSessionContext(
                sid=SubscriberID(id="IMSI12345"),
                apn=bytes("BLR", 'utf-8'),
                ue_ipv4="192.168.128.11",
                ue_ipv6="2001:db8::1",
                rat_type=RATType.Name(2),
                sm_session_state=SMSessionFSMState.Name(0),
                sm_session_version=0,
            ),
            rat_specific_context=RatSpecificContext(
                m5gsm_session_context=M5GSMSessionContext(
                    pdu_session_id=1,
                    request_type=RequestType.Name(
                        0,
                    ),
                    gnode_endpoint=TeidSet(
                        teid=10000,
                        end_ipv4_addr="192.168.60.141",
                    ),
                    pdu_session_type=PduSessionType.Name(2),
                    ssc_mode=SscMode.Name(2),
                ),
            ),
        )


class CreateAmfMultiSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    ue_ipv4="192.168.128.12",
                    apn=bytes("BLR", 'utf-8'),
                    rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(0),
                    sm_session_version=0,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=2,
                        request_type=RequestType.Name(
                            0,
                        ),
                        gnode_endpoint=TeidSet(
                            teid=10001,
                            end_ipv4_addr="192.168.60.141",
                        ),
                        pdu_session_type=PduSessionType.Name(0),
                        ssc_mode=SscMode.Name(2),
                    ),
                ),
            )


class ReleaseAmfSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    ue_ipv4="192.168.128.11",
                    apn=bytes("BLR", 'utf-8'),
                    rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=1,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV4,
                    ),
                ),
            )


class ReleaseAmfSessionIPv6(object):
    """
    Releases IPv6 PDU session.
    """

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    ue_ipv6="2001:db8::1",
                    apn=bytes("BLR", 'utf-8'),
                    rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=1,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV6,
                    ),
                ),
            )


class ReleaseAmfSessionIPv4v6(object):
    """
    Releases IPv4v6 PDU session.
    """

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    ue_ipv4="192.168.128.11",
                    ue_ipv6="2001:db8::1",
                    apn=bytes("BLR", 'utf-8'),
                    rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=1,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV4IPV6,
                    ),
                ),
            )


class CleanAmfSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    apn=bytes("BLR", 'utf-8'),
                    ue_ipv4="192.168.128.12",
                    rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=1,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV4,
                    ),
                ),
            )


class CreateAmfSecondSubSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI987654"),
                    ue_ipv4="192.168.128.110",
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(0),
                    sm_session_version=0,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=2,
                        request_type=RequestType.Name(
                            0,
                        ),
                        gnode_endpoint=TeidSet(
                            teid=5000,
                            end_ipv4_addr="192.168.60.141",
                        ),
                        pdu_session_type=PduSessionType.Name(0),
                        ssc_mode=SscMode.Name(2),
                    ),
                ),
            )


class CreateAmfSecondSubSecondSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI987654"),
                    ue_ipv4="192.168.128.111",
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(0),
                    sm_session_version=0,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=2,
                        request_type=RequestType.Name(
                            0,
                        ),
                        gnode_endpoint=TeidSet(
                            teid=300,
                            end_ipv4_addr="192.168.60.141",
                        ),
                        pdu_session_type=PduSessionType.Name(0),
                        ssc_mode=SscMode.Name(2),
                    ),
                ),
            )


class ReleaseAmfSecondSubSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI987654"),
                    ue_ipv4="192.168.128.110",
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=2,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV4,
                    ),
                ),
            )


class ReleaseAmfSecondSubSecondSession(object):

    def __init__(self):
        self._set_session =\
            SetSMSessionContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI987654"),
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    ue_ipv4="192.168.128.111",
                    sm_session_state=SMSessionFSMState.Name(4),
                    sm_session_version=6,
                ),
                rat_specific_context=RatSpecificContext(
                    m5gsm_session_context=M5GSMSessionContext(
                        pdu_session_id=2,
                        request_type=RequestType.Name(
                            1,
                        ),
                        pdu_session_type=PduSessionType.IPV4,
                    ),
                ),
            )


class CreateAmfIdleModeSession(object):

    def __init__(self):
        self._set_session =\
            SetSmNotificationContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(3),
                    sm_session_version=2,
                ),
                rat_specific_notification=RatSpecificNotification(
                    pdu_session_id=1,
                    request_type=RequestType.Name(1), notify_ue_event=1,
                ),
            )


class CreateAmfActiveModeSession(object):

    def __init__(self):
        self._set_session =\
            SetSmNotificationContext(
                common_context=CommonSessionContext(
                    sid=SubscriberID(id="IMSI12345"),
                    apn=bytes("BLR", 'utf-8'), rat_type=RATType.Name(2),
                    sm_session_state=SMSessionFSMState.Name(2),
                    sm_session_version=6,
                ),
                rat_specific_notification=RatSpecificNotification(
                    pdu_session_id=1,
                    request_type=RequestType.Name(1), notify_ue_event=5,
                ),
            )


@grpc_wrapper
def set_amf_session_tc1(client, args):
    print("=========TEST CASE-1 PDU SESSION ESTABLISHMENT===========")
    cls_sess = CreateAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc2(client, args):
    print("=========TEST CASE-2 PDU SESSION RELEASE============")
    cls_sess = ReleaseAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc3(client, args):
    print("====TEST CASE-5 MULTIPLE PDU SESSION IN A SUBSCRIBER=====")
    cls_sess = CreateAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfMultiSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc4(client, args):
    print("TEST CASE-6 SINGLE PDU SESSION RELEASE IN A MULTIPLE PDU SESSION")
    cls_sess = ReleaseAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc100(client, args):
    print("==================CLEAN PDU SESSIONS=====================")
    cls_sess = CleanAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc5(client, args):
    print("=====TEST CASE-5 MULTIPLE SUBSCRIBERS SESSION ESTABLISHMENT=====")
    cls_sess = CreateAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfSecondSubSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc101(client, args):
    print("=========MULTIPLE PDU SESSIONS RELEASE============")
    cls_sess = ReleaseAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = ReleaseAmfSecondSubSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc6(client, args):
    print("TEST CASE-5 MULTIPLE SUBSCRIBERS MULTIPLE SESSIONS ESTABLISHMENT==")
    cls_sess = CreateAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfMultiSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfSecondSubSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfSecondSubSecondSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc102(client, args):
    print("=========MULTIPLE SUBSCRIBERS MULTIPLE SESSIONS RELEASE==========")
    cls_sess = ReleaseAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CleanAmfSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = ReleaseAmfSecondSubSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)

    cls_sess = CreateAmfSecondSubSecondSession()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc7(client, args):
    print("=========TEST CASE-7 PDU SESSION IDLE MODE===========")
    cls_sess = CreateAmfIdleModeSession()
    print(cls_sess._set_session)
    response = client.SetSmfNotification(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc8(client, args):
    print("=====TEST CASE-8 PDU SESSION IDLE MODE TO ACTIVE MODE======")
    cls_sess = CreateAmfActiveModeSession()
    print(cls_sess._set_session)
    response = client.SetSmfNotification(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc11(client, args):
    """
    Create an object for CreateAmfSessionIPv6 and fill CommonSessionContext,
    RatSpecificContext and call grpc SetAmfSessionContext towards sessiond

    Args:
        client: self
        args: command line arguments, sid and apn
    """
    print("=========TEST CASE-11 IPv6 PDU SESSION ESTABLISHMENT===========")
    cls_sess = CreateAmfSessionIPv6()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc12(client, args):
    """
    Create an object for ReleaseAmfSessionIPv6 and fill CommonSessionContext,
    RatSpecificContext and call grpc SetAmfSessionContext towards sessiond

    Args:
        client: self
        args: command line arguments, sid and apn
    """
    print("=========TEST CASE-12 IPv6 PDU SESSION RELEASE============")
    cls_sess = ReleaseAmfSessionIPv6()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc13(client, args):
    """
    Create an object for CreateAmfSessionIPv4v6 and fill CommonSessionContext,
    RatSpecificContext and call grpc SetAmfSessionContext towards sessiond

    Args:
        client: self
        args: command line arguments, sid and apn
    """
    print("=========TEST CASE-13 IPv4IPV6 PDU SESSION ESTABLISHMENT===========")
    cls_sess = CreateAmfSessionIPv4v6()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


@grpc_wrapper
def set_amf_session_tc14(client, args):
    """
    Create an object for ReleaseAmfSessionIPv4v6 and fill CommonSessionContext,
    RatSpecificContext and call grpc SetAmfSessionContext towards sessiond

    Args:
        client: self
        args: command line arguments, sid and apn
    """
    print("=========TEST CASE-14 IPV4IPV6 PDU SESSION RELEASE============")
    cls_sess = ReleaseAmfSessionIPv4v6()
    print(cls_sess._set_session)
    response = client.SetAmfSessionContext(cls_sess._set_session)
    print(response)


def create_amf_parser(apps):
    """
    Create the argparse subparser for the ng_services app
    """

    app = apps.add_parser('amf_context')
    subparsers = app.add_subparsers(title='subcommands', dest='cmd')

    subcmd = subparsers.add_parser(
        'set_amf_session_tc1',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc1)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc2',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc2)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc3',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc3)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc4',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc4)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc100',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc100)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc5',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc5)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc101',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc101)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc6',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc6)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc102',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc102)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc7',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc7)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc8',
        help='AMF Set Session',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc8)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc11',
        help='AMF Set Session for UE IPv6',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc11)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc12',
        help='AMF Release Session for UE IPv6',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc12)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc13',
        help='AMF Set Session for UE IPv4v6',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc13)

    subcmd = subparsers.add_parser(
        'set_amf_session_tc14',
        help='AMF Release Session for UE IPv4v6',
    )
    subcmd.add_argument(
        '--sid',
        help='Subscriber_ID',
        default="imsi00000000001",
    )
    subcmd.add_argument('--apn', help='APN', default='12345')
    subcmd.set_defaults(func=set_amf_session_tc14)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='Management CLI for sessiond',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
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
    args.func(args, AmfPduSessionSmContextStub, 'sessiond')


if __name__ == "__main__":
    main()
