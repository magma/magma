# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM python:2.7.18-alpine3.11 AS builder

RUN apk add --no-cache \
        gcc=9.3.0-r0 \
        musl-dev=1.1.24-r3 \
        libffi-dev=3.2.1-r6 \
        openssl-dev=1.1.1l-r0 \
        git==2.24.4-r0 \
        && pip install --no-cache-dir virtualenv==20.10.0

WORKDIR /cbrs
ENV CLONE_URL="https://github.com/Wireless-Innovation-Forum/Citizens-Broadband-Radio-Service-Device.git"
RUN git clone \
    --branch v1.0.0.3 \
        --depth 1 \
        --filter=blob:none \
    ${CLONE_URL}

WORKDIR /cbrs/Citizens-Broadband-Radio-Service-Device
RUN virtualenv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --no-cache-dir -r requirements.txt

FROM python:2.7.18-alpine3.11
RUN apk add --no-cache \
		tcpdump==4.9.3-r1 \
		net-tools==1.60_git20140218-r2
COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /cbrs /cbrs
COPY Configuration /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/Configuration
RUN ln -s /opt/server/certs/server.cert /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/certificates/server.cert \
&& ln -s /opt/server/certs/server.key /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/certificates/server.key \
&& ln -s /opt/server/certs/ca.cert /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/certificates/ca.cert \
&& ln -s /opt/server/certs/server.cert /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/certificates/cpi.cert
ENV PATH="/opt/venv/bin:$PATH"
EXPOSE 9000

WORKDIR /cbrs/Citizens-Broadband-Radio-Service-Device/CBRS_Mock_SAS_Python_Master/cbrsPython-master/cbrsPython/controllers
ENTRYPOINT [ "sh" ]

