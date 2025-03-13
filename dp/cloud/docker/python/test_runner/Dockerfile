# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG ENV=standard
FROM python:3.9.2-slim-buster AS protos-generator

RUN apt-get update && apt-get install -y --no-install-recommends curl zip make unzip
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip \
    -o protoc3.zip
RUN unzip protoc3.zip -d protoc3 &&\
    mv protoc3/bin/protoc /bin/protoc &&\
    chmod a+rx /bin/protoc &&\
    mkdir -p /magma &&\
    mv ./protoc3/include/google /magma/google &&\
    rm -rf protoc3.zip protoc3
RUN pip3 install --no-cache-dir protobuf==3.19.1 setuptools==49.6.0 grpcio==1.37.1 grpcio-tools==1.37.1
ENV MAGMA_ROOT=/magma
ENV PYTHONPATH=$MAGMA_ROOT:$MAGMA_ROOT/build/gen
ENV PYTHON_BUILD=$MAGMA_ROOT/build
COPY ./protos $MAGMA_ROOT/protos/
COPY ./orc8r/protos $MAGMA_ROOT/orc8r/protos/
COPY ./lte/protos $MAGMA_ROOT/lte/protos/
COPY ./feg/protos $MAGMA_ROOT/feg/protos/
COPY ./cwf/protos $MAGMA_ROOT/cwf/protos/
COPY ./dp $MAGMA_ROOT/dp/
WORKDIR /magma/dp
RUN mkdir -p $PYTHON_BUILD && make protos

FROM python:3.9.2-slim-buster AS standard-version

COPY dp/cloud/python/magma/test_runner/requirements.txt /dp/cloud/python/magma/test_runner/requirements.txt
RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/magma/test_runner/requirements.txt

COPY dp/cloud/python/magma/db_service /dp/cloud/python/magma/db_service/
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/
COPY dp/cloud/python/magma/test_runner /dp/cloud/python/magma/test_runner/
COPY dp/cloud/python/magma/fixtures /dp/cloud/python/magma/fixtures/
COPY --chmod=755 dp/cloud/docker/python/test_runner/entrypoint.sh /entrypoint.sh
COPY --from=protos-generator /magma/build/gen /magma/build/gen

WORKDIR /dp/cloud/python/magma/test_runner
ENV PYTHONPATH=/magma/build/gen:/dp/cloud/python
ENTRYPOINT ["/entrypoint.sh"]
