# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM fluent/fluentd:v1.14.6-debian-1.0
USER root
  RUN gem install rubygems-update -v 3.2.33 --no-document && \
  update_rubygems && \
  gem install multi_json -v 1.15.0 --no-document && \
  gem install elasticsearch -v 7.13.0 --no-document && \
  gem install fluent-plugin-elasticsearch -v 5.2.1 --no-document && \
  gem install fluent-plugin-multi-format-parser -v 1.0.0 --no-document
USER fluent
