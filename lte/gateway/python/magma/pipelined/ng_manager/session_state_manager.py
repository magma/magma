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

from typing import (
    NamedTuple,
    Optional)

from enum import Enum
from magma.pipelined.set_interface_client import (
    send_periodic_session_update)

from magma.pipelined.ng_manager.session_state_manager_util import (
    pdr_create_rule_entry)

from lte.protos.session_manager_pb2 import (
    UPFSessionConfigState,
    UPFSessionState)

from lte.protos.pipelined_pb2 import (
    PdrState,
    UPFSessionContextState,
    OffendingIE,
    CauseIE)

# Help to build failure report
MsgParseOutput = NamedTuple(
                   'MsgParseOutput',
                   [('offending_ie', OffendingIE),
                    ('cause_info', int)])

class SessionMessageType(Enum):
    MSG_TYPE_CONTEXT_STATE = 1  # In response to session configuraiton from SMF
    MSG_TYPE_CONFIG_STATE  = 2  # Periodic messages to update session information to SMF
    MSG_TYPE_CONFIG_REPORT = 3  # Event or Periodic Report from session to SMF

class SessionStateManager:
    send_message_offset = 0
    periodic_config_msg_count = 0

    """
    This controller manages session state information
    and reports session config to SMF.
    """
    def __init__(self, loop, logger):
        """
        Launch the SessionStateManager under ng_services
        """
        self._loop = loop
        self.logger = logger

    # Creating the dict entries for the far group
    @staticmethod
    def _pdr_create_rule_group(new_session, pdr_rules) -> Optional[MsgParseOutput]:

        for pdr_entry in new_session.set_gr_pdr:
            # PDR Validation
            if pdr_entry.HasField('pdi') == False or pdr_entry.pdr_id == 0:
                offending_ie = OffendingIE(identifier=pdr_entry.pdr_id,
                                           version=pdr_entry.pdr_version)
                return MsgParseOutput(offending_ie, CauseIE.MANDATORY_IE_INCORRECT)

            # If session is creted or activiated FAR_IDs cann't be 0
            if  len(pdr_entry.set_gr_far.ListFields()) == 0 and \
                     pdr_entry.pdr_state == PdrState.Value('INSTALL'):
                offending_ie = OffendingIE(identifier=pdr_entry.pdr_id,
                                           version=pdr_entry.pdr_version)
                return MsgParseOutput(offending_ie, CauseIE.INVALID_FORWARDING_POLICY)

            pdr_rules.update({pdr_entry.pdr_id: pdr_create_rule_entry(pdr_entry)})

        return None

    @staticmethod
    def validate_session_msg(new_session):
        """
        Initial session validation. Check subscriber_id and version no.
        If existing session with same id is found return the
        existing session
        """

        #if SEID is not found or version is 0
        if len(new_session.subscriber_id) == 0 or\
           new_session.session_version == 0:
            return CauseIE.SESSION_CONTEXT_NOT_FOUND

        # Check if the new session operation is without any pdr group
        if not new_session.set_gr_pdr:
            return CauseIE.MANDATORY_IE_MISSING

        return CauseIE.REQUEST_ACCEPTED

    def process_session_message(self, new_session, process_pdr_rules):
        """
        Process the messages recevied from session. Return True
        if parsing is successfull.
        """

        # Assume things are green
        context_response =\
             UPFSessionContextState(cause_info=CauseIE(cause_ie=CauseIE.REQUEST_ACCEPTED),
                                    session_snapshot=UPFSessionState(subscriber_id=new_session.subscriber_id,
                                                                     local_f_teid=new_session.local_f_teid,
                                                                     session_version=new_session.session_version))

        context_response.cause_info.cause_ie = \
                  SessionStateManager.validate_session_msg(new_session)
        if context_response.cause_info.cause_ie != CauseIE.REQUEST_ACCEPTED:
            self.logger.error("Error : Parsing Error in SetInterface Message %d",
                              context_response.cause_info.cause_ie)
            return context_response

        #Create PDR rules
        pdr_validator = SessionStateManager._pdr_create_rule_group(new_session, process_pdr_rules)
        if pdr_validator:
            context_response.failure_rule_id.pdr.extend([pdr_validator.offending_ie])
            context_response.cause_info.cause_ie = pdr_validator.cause_info

        return context_response

    @classmethod
    def report_session_config_state(cls, session_config_dict, sessiond_stub):

        SessionStateManager.send_message_offset +=1

        # Send session config messages every 10 seconds
        if SessionStateManager.send_message_offset % 5:
            return

        session_config_list = []
        for index in session_config_dict:
            session_config_list.append(session_config_dict[index])

        session_config_msg = UPFSessionConfigState(upf_session_state=session_config_list)
        if send_periodic_session_update(session_config_msg, sessiond_stub) == True:
            SessionStateManager.periodic_config_msg_count += 1
