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
    agw_valid_env = os.environ["AGW_VALID"]
    feg_docker_env = os.environ["FEG_DOCKER"]
    feg_jfrog_env = os.environ["FEG_JFROG"]
    feg_valid_env = os.environ["FEG_VALID"]

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

    # Add agw arifacts
    agw_artifacts = {}
    agw_valid = False
    if agw_valid_env == "true":
        try:
            agw_artifacts = json.loads(agw_artifacts_env)
            agw_valid = True
        except ValueError:
            print("Decoding agw artifacts JSON has failed: ", agw_artifacts_env)
    build_info["agw"] = {
        "valid": agw_valid,
        "artifacts": agw_artifacts,
    }

    # Add feg arifacts
    feg_docker = {}
    feg_jfrog = {}
    feg_valid = False
    if feg_valid_env == "true":
        try:
            feg_docker = json.loads(feg_docker_env)
            feg_jfrog = json.loads(feg_jfrog_env)
            feg_valid = True
        except ValueError:
            print("Decoding feg artifacts JSON has failed: ", feg_docker_env, feg_jfrog_env)
    build_info["feg"] = {
        "valid": feg_valid,
        "docker": feg_docker,
        "jfrog": feg_jfrog,
    }

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
