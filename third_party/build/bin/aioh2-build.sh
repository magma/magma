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



#build package from aioh2 master branch.

PKGNAME=aioh2
WORK_DIR=/tmp/build-"$PKGNAME"

PWD1=$(pwd)

mkdir "$WORK_DIR"

cd "$WORK_DIR" || exit
git clone https://github.com/URenko/aioh2.git

cd aioh2 || exit
sed -i 's/0.2.2/0.2.3/g' setup.py

python3 -m pip install --upgrade build
sudo apt-get install python3-venv

pip3 install stem stdeb

pip3 install  debhelper

sudo apt install python3-all

sudo apt install debhelper

sudo apt install dh-python

python3 setup.py --command-packages=stdeb.command bdist_deb

cp deb_dist/python3-aioh2*.deb "$PWD1"

