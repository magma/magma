ARG ENV=standard
FROM python:3.9.2-slim-buster as protos-generator

ENV MAGMA_ROOT=/magma
RUN apt-get update && apt-get install -y --no-install-recommends curl=7.64.0-4+deb10u2 zip=3.0-11+b1 make=4.2.1-1.2 unzip=6.0-23+deb10u2
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.18.0/protoc-3.18.0-linux-x86_64.zip \
    -o protoc3.zip
RUN unzip protoc3.zip -d protoc3 &&\
    mv protoc3/bin/protoc /bin/protoc &&\
    chmod a+rx /bin/protoc &&\
    mkdir -p $MAGMA_ROOT &&\
    mv ./protoc3/include/google $MAGMA_ROOT/google &&\
    rm -rf protoc3.zip protoc3
RUN pip3 install --no-cache-dir protobuf==3.19.1 setuptools==49.6.0 grpcio==1.37.1 grpcio-tools==1.37.1
ENV PYTHONPATH=$MAGMA_ROOT:$MAGMA_ROOT/build/gen
ENV PYTHON_BUILD=$MAGMA_ROOT/build
COPY ./protos $MAGMA_ROOT/protos/
COPY ./orc8r/protos $MAGMA_ROOT/orc8r/protos/
COPY ./lte/protos $MAGMA_ROOT/lte/protos/
COPY ./feg/protos $MAGMA_ROOT/feg/protos/
COPY ./cwf/protos $MAGMA_ROOT/cwf/protos/
COPY ./dp $MAGMA_ROOT/dp/
WORKDIR $MAGMA_ROOT/dp
RUN mkdir -p $PYTHON_BUILD && make protos

FROM python:3.9.2-slim-buster as standard-version

ENV PC_DIRECTORY=dp/cloud/python/magma/protocol_controller
COPY $PC_DIRECTORY/requirements.txt \
     /$PC_DIRECTORY/requirements.txt
RUN pip3 install --upgrade pip --no-cache-dir -r /$PC_DIRECTORY/requirements.txt

#FROM standard-version as tests-version
#
#COPY $PC_DIRECTORY/plugins/cbsd_sas/tests/requirements.txt \
#     /$PC_DIRECTORY/plugins/cbsd_sas/tests/requirements.txt
#WORKDIR /$PC_DIRECTORY
#RUN pip3 install --upgrade pip --no-cache-dir -r plugins/cbsd_sas/tests/requirements.txt

# hadolint ignore=DL3006
FROM ${ENV}-version as final

ENV PYTHONPATH=/magma/build/gen:/dp/cloud/python
COPY $PC_DIRECTORY /$PC_DIRECTORY/
COPY --from=protos-generator /magma/build/gen /magma/build/gen
WORKDIR /$PC_DIRECTORY
EXPOSE 8000
CMD ["python", "service.py"]
