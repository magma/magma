################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

FROM golang:alpine as builder
RUN apk add git gcc musl-dev bash protobuf
COPY cwf/radius/ /src/cwf/radius
COPY lib/go/ /src/lib/go
WORKDIR /src/cwf/radius
RUN go mod download
RUN ./run.sh build 
COPY cwf/radius/docker-entrypoint.sh /src/cwf/radius/bin/docker-entrypoint.sh
RUN chmod 0755 /src/cwf/radius/bin/docker-entrypoint.sh

FROM alpine
RUN apk add gettext musl
COPY --from=builder /src/cwf/radius/radius /app/
COPY --from=builder /src/cwf/radius/*.config.json /app/
WORKDIR /app

CMD ["./radius", "-config", "lb.config.json"]
