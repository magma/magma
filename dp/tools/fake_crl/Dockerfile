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

RUN virtualenv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY requirements.txt /requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

FROM python:2.7.18-alpine3.11

RUN apk add --no-cache \
	tcpdump=4.9.3-r1 \
	net-tools=1.60_git20140218-r2 \
	bash=5.0.11-r1 \
	openssl=1.1.1l-r0
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY main.py /opt/server/main.py
COPY openssl.cnf /opt/server/openssl.cnf
COPY revoke_and_generate_crl.sh /opt/server/revoke_and_generate_crl.sh

EXPOSE 9007
WORKDIR /opt/server
ENTRYPOINT ["python"]
CMD ["main.py"]
