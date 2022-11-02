#!/bin/bash
#
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build package from aioh2 master branch.

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR"/../lib/util.sh
PKGNAME=aioh2
WORK_DIR=/tmp/build-"$PKGNAME"
PWD1=$(pwd)

function buildrequires() {
    echo \
        python3-all \
        debhelper \
        dh-python \
        python3-stem \
        fakeroot
}

if_subcommand_exec

# Create python virtual env to install python3 packages that do not work with apt
python3 -m venv env
VENV_DIR=$(pwd)/env
source env/bin/activate

# Install wheel, otherwise building wheel fails for stdeb
pip3 install wheel

pip3 install stdeb

pip3 install debhelper

# The resulting package is placed in $OUTPUT_DIR
# or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=$(pwd)
else
  OUTPUT_DIR="$1"
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# build from source
if [ -d "$WORK_DIR" ]; then
  rm -rf "$WORK_DIR"
fi
mkdir "$WORK_DIR"
cd "$WORK_DIR"

git clone https://github.com/URenko/aioh2.git

cd aioh2
sed -i 's/0.2.2/0.2.3/g' setup.py

python3 setup.py --command-packages=stdeb.command bdist_deb

cp deb_dist/python3-aioh2*.deb "$PWD1"

# Deactivate and remove python virtual env
deactivate
sudo rm -rf "$VENV_DIR"
