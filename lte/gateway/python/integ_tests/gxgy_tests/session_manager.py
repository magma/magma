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
from lte.protos import session_manager_pb2, session_manager_pb2_grpc
from lte.protos.session_manager_pb2 import (
    ChargingCredit,
    CreditUnit,
    CreditUpdateResponse,
    GrantedUnits,
    UsageMonitoringCredit,
    UsageMonitoringUpdateResponse,
)
from ryu.lib import hub


class MockSessionManager(session_manager_pb2_grpc.CentralSessionControllerServicer):
    def __init__(self, *args, **kwargs):
        super(MockSessionManager, self).__init__(*args, **kwargs)
        self.mock_create_session = None
        self.mock_terminate_session = None
        self.mock_update_session = None

    def CreateSession(self, req, ctx):
        return self.mock_create_session(req, ctx)

    def TerminateSession(self, req, ctx):
        return self.mock_terminate_session(req, ctx)

    def UpdateSession(self, req, ctx):
        return self.mock_update_session(req, ctx)


def get_from_queue(q, retries=10, sleep_time=0.5):
    """
    get an object from a hub.Queue by polling
    Args:
        q (hub.queue): queue to wait on
        retries (int): number of times to try getting from the queue
        sleep_time (float): amount of seconds to wait between retries
    Returns:
        The object, or None if it wasn't retrieved in the number of retries
    """
    for _ in range(retries):
        try:
            return q.get(block=False)
        except hub.QueueEmpty:
            hub.sleep(0.5)
            continue
    return None


def get_standard_update_response(
    update_complete, monitor_complete, quota,
    is_final=False,
    success=True,
    monitor_action=UsageMonitoringCredit.CONTINUE,
):
    """
    Create a CreditUpdateResponse with some useful defaults
    Args:
        update_complete (hub.Queue): eventlet queue to wait for update responses on
        monitor_complete (hub.Queue): eventlet queue to wait for monitor responses on
        quota (int): number of bytes to return
        is_final (bool): True if these are the last credits to return
        success (bool): True if the update was successful
        monitor_action (UsageMonitoringCredit.Action): action to take with response,
            defaults to CONTINUE
    """
    def update_response(*args, **kwargs):
        charging_responses = []
        monitor_responses = []
        for update in args[0].updates:
            charging_responses.append(
                create_update_response(
                    update.sid, update.usage.charging_key, quota,
                    is_final=is_final, success=success,
                ),
            )
            update_complete.put(update)
        for monitor in args[0].usage_monitors:
            monitor_responses.append(
                create_monitor_response(
                    monitor.sid, monitor.update.monitoring_key, quota,
                    monitor.update.level, action=monitor_action,
                    success=success,
                ),
            )
            monitor_complete.put(monitor)
        return session_manager_pb2.UpdateSessionResponse(
            responses=charging_responses,
            usage_monitor_responses=monitor_responses,
        )
    return update_response


def create_update_response(
    imsi, charging_key, total_quota,
    is_final=False,
    success=True,
):
    """
    Create a CreditUpdateResponse with some useful defaults
    Args:
        imsi (string): subscriber id
        charging_key (int): rating group
        quota (int): number of bytes to return
        is_final (bool): True if these are the last credits to return
        success (bool): True if the update was successful
    """
    return CreditUpdateResponse(
        success=success,
        sid=imsi,
        charging_key=charging_key,
        credit=ChargingCredit(
            granted_units=GrantedUnits(
                total=CreditUnit(is_valid=True, volume=total_quota),
            ),
            is_final=is_final,
        ),
    )


def create_monitor_response(
    imsi, m_key, total_quota, level,
    action=UsageMonitoringCredit.CONTINUE,
    success=True,
):
    """
    Create a UsageMonitoringUpdateResponse with some useful defaults
    Args:
        imsi (string): subscriber id
        m_key (string): monitoring key
        quota (int): number of bytes to return
        level (MonitoringLevel): session level or rule level
        action (UsageMonitoringCredit.Action): action to take with response,
            defaults to CONTINUE
        success (bool): True if the update was successful
    """
    return UsageMonitoringUpdateResponse(
        success=success,
        sid=imsi,
        credit=UsageMonitoringCredit(
            monitoring_key=m_key,
            granted_units=GrantedUnits(
                total=CreditUnit(is_valid=True, volume=total_quota),
            ),
            level=level,
            action=action,
        ),
    )
