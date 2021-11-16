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

import json
import os
import time

import pyrebase


def main():
    """Publish builds and workloads to Firebase realtime database"""
    # Read db config
    firebase_config = os.environ["FIREBASE_CONFIG"]
    config = json.loads(firebase_config)

    # Initialize pyrebase app
    firebase = pyrebase.initialize_app(config)

    # Get a reference to the auth service
    auth = firebase.auth()

    # Log the user in
    user = auth.sign_in_with_email_and_password(
        config["auth_email"], config["auth_password"],
    )

    # Get a reference to the database service
    db = firebase.database()

    # Grab environment variables
    workers_env = os.environ["WORKERS"]
    build_id = os.environ["BUILD_ID"]
    build_metadata_env = os.environ["BUILD_METADATA"]
    agw_artifacts_env = os.environ["AGW_ARTIFACTS"]
    feg_artifacts_env = os.environ["FEG_ARTIFACTS"]

    # Prepare list of registered test workers
    workers = [x.strip() for x in workers_env.split(",")]

    # Prepare build metadata
    build_metadata = {}
    try:
        build_metadata = json.loads(build_metadata_env)
        build_metadata["timestamp"] = int(time.time())
    except ValueError:
        print("Decoding build_metadata_env JSON failed: ", build_metadata_env)
    build_info = {"metadata": build_metadata}

    # Add AGW artifacts
    agw_artifacts = {}
    try:
        agw_artifacts = json.loads(agw_artifacts_env)
    except ValueError:
        print("Decoding agw artifacts JSON has failed: ", agw_artifacts_env)
        agw_artifacts = {"packages": [], "valid": False}

    # TODO: Remove this backward compatibility code
    for package in agw_artifacts["packages"]:
        if "magma_" in package:
            agw_artifacts["artifacts"] = {"downloadUri": package}
            break
    build_info["agw"] = agw_artifacts

    # Add FEG artifacts
    feg_artifacts = {}
    try:
        feg_artifacts = json.loads(feg_artifacts_env)
    except ValueError:
        print("Decoding feg artifacts JSON has failed: ", feg_artifacts_env)
        feg_artifacts = {"packages": [], "valid": False}
    build_info["feg"] = feg_artifacts

    # Prepare workload
    workload = {
        "build_id": build_id,
        "priority": 2,
        "state": "queued",
        "timestamp": int(time.time()),
    }

    # Publish build to Firebase realtime database
    print("Publishing build to database: builds/", build_id)
    print(build_info)
    db.child("builds").child(build_id).set(build_info, user["idToken"])

    # Publish workloads to Firebase realtime database
    print("Pushing the following workload to workers: ", workers)
    print(workload)
    for worker in workers:
        db.child("workers").child(worker).child("workloads").push(
            workload, user["idToken"],
        )


if __name__ == "__main__":
    main()
