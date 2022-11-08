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
from typing import Dict

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

    # Prepare list of registered test workers
    workers = [x.strip() for x in workers_env.split(",")]

    # Prepare build metadata
    build_metadata = {}
    try:
        build_metadata = json.loads(build_metadata_env)
        build_metadata["timestamp"] = int(time.time())
    except ValueError:
        print(f"Decoding build_metadata_env JSON failed: {build_metadata_env}")
    build_info = {"metadata": build_metadata}

    builds = ["AGW", "FEG", "ORC8R", "CWAG", "NMS"]
    for build in builds:
        artifact_env = os.environ[f"{build}_ARTIFACTS"]
        artifacts = {}
        try:
            artifacts = json.loads(artifact_env)
        except ValueError:
            print(
                f"Decoding {build} artifacts JSON "
                f"has failed: {artifact_env}",
            )
            artifacts = {"packages": [], "valid": False}
        if build == "AGW":
            artifacts = _set_download_uri(artifacts)
        build_info[f"{build.lower()}"] = artifacts

    # Prepare workload
    workload = {
        "build_id": build_id,
        "priority": 2,
        "state": "queued",
        "timestamp": int(time.time()),
    }

    # Publish build to Firebase realtime database
    print(f"Publishing build to database: builds/{build_id}")
    print(build_info)
    db.child("builds").child(build_id).set(build_info, user["idToken"])

    # Publish workloads to Firebase realtime database
    print("Pushing the following workload to workers: ", workers)
    print(workload)
    for worker in workers:
        db.child("workers").child(worker).child("workloads").push(
            workload, user["idToken"],
        )


def _set_download_uri(artifacts: Dict) -> Dict:
    # TODO: Remove this backward compatibility code
    for package in artifacts["packages"]:
        if "magma_" in package:
            artifacts["artifacts"] = {"downloadUri": package}
            break
    return artifacts


if __name__ == "__main__":
    main()
