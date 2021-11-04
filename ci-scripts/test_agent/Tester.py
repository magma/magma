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


class TesterState:
    READY = 'READY'
    BUSY = 'BUSY'
    OFFLINE = 'OFFLINE'


class Tester:
    def __init__(self, id):
        self.id = id
        self.current_workload = None
        self.current_build = None
        self.state = TesterState.READY

    def test_ended(self):
        # Get test results and prepare report

        # Callback
        # self.callback(self.id, self.current_workload, valid, report)

        self.current_workload = None
        self.current_build = None
        self.state = TesterState.READY

    def start_test(self, workload, build, test_done_callback):
        self.current_workload = workload
        self.current_build = build
        self.callback = test_done_callback

        # start test on current_workload
        print('Tester {} Starting test on workload'.format(self.id))
        print(workload.key(), "==>", workload.val())

        # register callback to call_ended()

        self.state = TesterState.BUSY
        return

    def is_ready(self):
        return self.state == TesterState.READY

    def get_state(self):
        return self.state

    def get_current_workload(self):
        return self.current_workload

    def get_current_build(self):
        return self.current_build
