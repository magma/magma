# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM linuxserver/openssh-server:8.6_p1-r2-ls56
RUN apk add --no-cache gettext=0.21-r0

ENV TPL="config.template"

COPY "$TPL" "/$TPL"
COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
