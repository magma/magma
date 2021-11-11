"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import time

from FirebaseClient import FirebaseClient
from Tester import Tester


class WorkerState:
    READY = 0
    BUSY = 1


num_of_testers = 1
# Initialize firebase client as global
db_client = FirebaseClient()


def test_done_callback(tester_id, workload, verdict, report):
    db_client.user = db_client.auth.refresh(db_client.user["refreshToken"])
    db_client.mark_workload_done(workload)
    db_client.push_test_report(workload, verdict, report)
    return


def main():
    state = WorkerState.READY
    testers = []
    for i in range(num_of_testers):
        testers.append(Tester(i))
    while True:
        testers_state = []
        for tester in testers:
            if tester.is_ready():
                print("tester {} is ready".format(tester.get_id()))
                new_workload = db_client.pop_next_workload()
                print("new workload is", new_workload)
                if new_workload:
                    build = db_client.get_build(new_workload)
                    if build:
                        tester.start_test(new_workload, build, test_done_callback)
                else:
                    print("Waiting for new workload")
            current_workload = tester.get_current_workload()
            if current_workload:
                testers_state.append(
                    {
                        "state": tester.get_state(),
                        "current_workload": current_workload.val(),
                    }
                )
            else:
                testers_state.append(
                    {"state": tester.get_state(), "current_workload": None}
                )

        time.sleep(1)
        db_client.user = db_client.auth.refresh(db_client.user["refreshToken"])
        db_client.update_worker_state(num_of_testers, testers_state)
        time.sleep(14)


if __name__ == "__main__":
    main()
