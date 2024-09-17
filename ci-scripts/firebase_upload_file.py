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
import json
import os

from firebase_admin import credentials, initialize_app, storage

# Initiate the parser
parser = argparse.ArgumentParser()

# Add arguments
parser.add_argument("--filename", "-f", required=True, help="file to upload")
parser.add_argument("--output", "-o", help="file to print url")

# Read arguments from the command line
args = parser.parse_args()
print("Uploading %s to firebase storage" % args.filename)

# Read Firebase service account config from envirenment
firebase_config = os.environ["FIREBASE_SERVICE_CONFIG"]
config = json.loads(firebase_config)

cred = credentials.Certificate(config)
initialize_app(cred, {'storageBucket': 'magma-ci.appspot.com'})

# Upload file
bucket = storage.bucket()
blob = bucket.blob(args.filename)
blob.upload_from_filename(args.filename)

# Make public access from the URL
blob.make_public()

print("File URL is: ", blob.public_url)
if args.output:
    with open(args.output, 'w') as f:
        f.write(blob.public_url)
        f.close()
