################################################################################
#  Copyright 2020 The Magma Authors.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
################################################################################
MAGMA_ROOT = ~/magma
PY_LTE     = $(MAGMA_ROOT)/lte/gateway/python
PY_ORC8R   = $(MAGMA_ROOT)/orc8r/gateway/python

magma.lockfile: $(PY_LTE)/setup.py $(PY_ORC8R)/setup.py
	./pydep finddep --install-from-repo -l ./magma.lockfile.$(os_release)  $(PY_ORC8R)/setup.py $(PY_LTE)/setup.py
