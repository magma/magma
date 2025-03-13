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
	&& pip install --no-cache-dir virtualenv==20.10.0

RUN virtualenv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY requirements.txt /requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

FROM python:2.7.18-alpine3.11

COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN mkdir /opt/server  && chown -R nobody /opt/server
COPY fake_sas.py /opt/server
COPY sas_interface.py /opt/server
COPY sas.cfg /opt/server
EXPOSE 9000
WORKDIR /opt/server
USER nobody
ENTRYPOINT ["python"]
CMD ["fake_sas.py"]

