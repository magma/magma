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

import json
import threading
import time
from datetime import datetime

import pyrebase


class FirebaseClient:
    def __init__(self):
        # Read db config
        f = open("FirebaseConfig.json")
        self.config = json.load(f)
        # Initialize pyrebase app
        self.firebase = pyrebase.initialize_app(self.config)
        # Get a reference to the auth service
        self.auth = self.firebase.auth()
        # Log the user in
        self.user = self.auth.sign_in_with_email_and_password(
            self.config["auth_email"], self.config["auth_password"],
        )
        # Get a reference to the database service
        self.db = self.firebase.database()
        # token refresh timer
        self.RefreshTokenTime = 1800
        threading.Timer(self.RefreshTokenTime, self.do_token_refresh).start()

    def push_test_report(self, workload, verdict, report):
        build_id = workload.val().get("build_id")
        data = {
            "timestamp": int(time.time()),
            "verdict": verdict,
            "report": report,
        }
        # Use build_id as key for report. Overwrites any existing report
        self.db.child("workers").child(self.config["agent_id"]).child("reports").child(build_id).set(
            data, self.user["idToken"],
        )

    def update_worker_state(self, num_of_testers, testers_state):
        data = {
            "timestamp": int(time.time()),
            "num_of_testers": num_of_testers,
            "testers_state": testers_state,
        }
        self.db.child("workers").child(self.config["agent_id"]).child("state").update(
            data, self.user["idToken"],
        )

    def print_reports(self):
        reports = (
            self.db.child("workers")
            .child(self.config["agent_id"])
            .child("reports")
            .get(self.user["idToken"])
        )
        for r in reports.each():
            print(r.val())

    def pop_next_workload(self):
        workloads = (
            self.db.child("workers")
            .child(self.config["agent_id"])
            .child("workloads")
            .get(self.user["idToken"])
        )
        bestWorkload = None
        for workload in workloads.each():
            print(workload.val())
            # skip if not valid
            if workload.val().get("state") != "queued":
                continue

            # check if first workload
            if bestWorkload is None:
                bestWorkload = workload
            # check if higher priority or newer
            elif (
                workload.val().get("priority") > bestWorkload.val().get("priority")
            ) or (
                (workload.val().get("priority") == bestWorkload.val().get("priority"))
                and (
                    workload.val().get("timestamp")
                    > bestWorkload.val().get("timestamp")
                )
            ):
                bestWorkload = workload

        if bestWorkload:
            # remove from workloads
            self.db.child("workers").child(self.config["agent_id"]).child(
                "workloads",
            ).child(bestWorkload.key()).update(
                {"state": "in_progress"}, self.user["idToken"],
            )
            print("Best workload: ", bestWorkload.key(), "==>", bestWorkload.val())
            return bestWorkload
        else:
            return None

    def get_build(self, workload):
        build_id = workload.val().get("build_id")
        build = self.db.child("builds").child(build_id).get(self.user["idToken"])
        if build:
            return build
        else:
            return None

    def do_token_refresh(self):
        """Run in its own thread and refresh token every 30 min."""

        print("Refreshing Firebase Token:", datetime.now())
        self.user = self.auth.refresh(self.user["refreshToken"])
        threading.Timer(self.RefreshTokenTime, self.do_token_refresh).start()

    def mark_workload_done(self, workload):
        print("Workload execution completed clearing it")
        self.db.child("workers").child(self.config["agent_id"]).child(
            "workloads",
        ).child(workload.key()).set({}, self.user["idToken"])
