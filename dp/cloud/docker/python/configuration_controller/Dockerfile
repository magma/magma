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

ARG ENV=standard
FROM python:3.9.2-slim-buster AS standard-version

ENV CC_DIRECTORY=dp/cloud/python/magma/configuration_controller
ENV DB_DIRECTORY=dp/cloud/python/magma/db_service
ENV FC_DIRECTORY=dp/cloud/python/magma/fluentd_client
ENV MC_DIRECTORY=dp/cloud/python/magma/metricsd_client
ENV GW_COMMON_DIRECTORY=orc8r/gateway/python/magma/common
COPY $CC_DIRECTORY/requirements.txt \
     /$CC_DIRECTORY/requirements.txt
WORKDIR /$CC_DIRECTORY
RUN pip3 install --upgrade pip --no-cache-dir -r requirements.txt

#FROM standard-version AS tests-version
#COPY $CC_DIRECTORY/tests/requirements.txt \
#     /$CC_DIRECTORY/tests/requirements.txt
#RUN pip3 install --upgrade pip --no-cache-dir -r tests/requirements.txt

# hadolint ignore=DL3006
FROM ${ENV}-version AS final
COPY $CC_DIRECTORY /$CC_DIRECTORY/
COPY $DB_DIRECTORY /$DB_DIRECTORY/
COPY $FC_DIRECTORY /$FC_DIRECTORY/
COPY $MC_DIRECTORY /$MC_DIRECTORY/
COPY --from=protos-generator /magma/build/gen /magma/build/gen
COPY $GW_COMMON_DIRECTORY/metrics_export.py /$GW_COMMON_DIRECTORY/metrics_export.py
COPY $GW_COMMON_DIRECTORY/__init__.py /$GW_COMMON_DIRECTORY/__init__.py
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/
ENV PYTHONPATH=/:/dp/cloud/python:/orc8r/gateway/python:/magma/build/gen
ENTRYPOINT ["python"]
CMD ["run.py"]
