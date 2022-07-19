# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG GOLANG_VERSION="1.18.3"
FROM golang:${GOLANG_VERSION}-alpine3.16 AS build
COPY dp/cloud/go/active_mode_controller /active_mode_controller
WORKDIR /active_mode_controller/cmd
RUN go build

FROM alpine:3.14.3 as final
COPY --from=build /active_mode_controller/cmd/cmd /active_mode_controller/cmd
WORKDIR /active_mode_controller
CMD ["./cmd"]
