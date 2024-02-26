# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

PYTHON_VERSION=3.8

.test: .tests .sudo_tests

RESULTS_DIR := /var/tmp/test_results
CODECOV_DIR := /var/tmp/codecovs

.tests:
ifdef TESTS
ifndef SKIP_NON_SUDO_TESTS
	$(eval NAME ?= $(shell $(BIN)/python$(PYTHON_VERSION) setup.py --name))
	. $(PYTHON_BUILD)/bin/activate; $(BIN)/pytest --junit-xml=$(RESULTS_DIR)/tests_$(NAME).xml --cov=magma --cov-branch --cov-report xml:$(CODECOV_DIR)/cover_$(NAME).xml $(TESTS)
endif
endif

.sudo_tests:
ifdef SUDO_TESTS
ifndef SKIP_SUDO_TESTS
	$(eval NAME ?= $(shell $(BIN)/python$(PYTHON_VERSION) setup.py --name))
	. $(PYTHON_BUILD)/bin/activate; sudo $(BIN)/pytest --junit-xml=$(RESULTS_DIR)/sudo_$(NAME).xml --cov=magma --cov-branch --cov-report xml:$(CODECOV_DIR)/cover_sudo_$(NAME).xml $(SUDO_TESTS)
endif
endif
