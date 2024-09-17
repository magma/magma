# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM python:3.9.10-alpine

RUN addgroup -S linter && adduser -S -G linter linter

# Installing wemake-python-styleguide==0.17.0 is not
# working with a requirements.in file.
# hadolint ignore=DL3013
RUN pip install --no-cache-dir wemake-python-styleguide==0.17.0 \
                add-trailing-comma \
                isort \
                autopep8

USER linter
WORKDIR /code/
