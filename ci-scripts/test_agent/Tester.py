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
import pickle
import subprocess
import threading


class TesterState:
    READY = "READY"
    BUSY = "BUSY"
    OFFLINE = "OFFLINE"


class Tester:
    def __init__(self, id):
        self.id = id
        self.current_workload = None
        self.current_build = None
        self.state = TesterState.READY

    def test_ended(self):
        # Get test results and prepare report
        verdict, report = self.load_test_results(dbfile="/tmp/test_res_pickle")
        # Callback
        print("Test ended. Invoking callback function")
        self.callback(self.id, self.current_workload, verdict, report)
        self.current_workload = None
        self.current_build = None
        self.state = TesterState.READY

    def load_test_results(self, dbfile):
        """test run can dump their results (dict) in pickle
        which can be loaded here and used to push back to DB
        data format = {'verdict' = True/False, 'report': 'html file text'}
        """
        try:
            with open(dbfile, "rb") as dbfile:
                db = pickle.load(dbfile)
        except (IOError, OSError, pickle.PickleError, pickle.UnpicklingError):
            print("Not a valid returning null")
            return False, None
        else:
            with open(db["report"], "r") as rfile:
                report = rfile.read()
            verdict = "pass" if db["verdict"] else "fail"
            return verdict, report

    def start_test(self, workload, build, test_done_callback):
        self.current_workload = workload
        self.current_build = build
        self.callback = test_done_callback

        def run_hil_thread(call_ended, popen_args):
            proc = subprocess.Popen(*popen_args, stdout=subprocess.PIPE, shell=True)
            proc.wait()
            call_ended()
            return

        if self.current_build.val().get("agw", {}).get("valid", False):
            magma_build = next(
                (
                    pkg
                    for pkg in self.current_build.val().get("agw", {}).get("packages", [])
                    if "magma_" in pkg
                ),
                None,
            )
            if magma_build:
                # start test on current_workload
                print(f"Tester {self.id} Starting test on workload")
                print(workload.key(), "==>", workload.val())

                # register callback to call_ended()
                # TODO pass pickle file from here to test run.

                thread = threading.Thread(
                    target=run_hil_thread,
                    args=(self.test_ended, ["./run_test.sh " + magma_build]),
                )
                thread.start()

                self.state = TesterState.BUSY
                print(f"test {self.id} started on workload")
            else:
                print("No Magma package found in packages list")
                self.callback(self.id, self.current_workload, "fail", "NA")
                self.current_workload = None
                self.current_build = None
                self.state = TesterState.READY

        else:
            print(f"Invalid AGW build {self.id}. Terminate test")
            self.callback(self.id, self.current_workload, "fail", "NA")
            self.current_workload = None
            self.current_build = None
            self.state = TesterState.READY

        return

    def is_ready(self):
        return self.state == TesterState.READY

    def get_id(self):
        return self.id

    def get_state(self):
        return self.state

    def get_current_workload(self):
        return self.current_workload

    def get_current_build(self):
        return self.current_build
