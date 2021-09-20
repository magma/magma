ARG ENV=standard
FROM python:3.9.2-slim-buster as protos-generator

ENV MAGMA_ROOT=/magma
RUN apt update && apt install -y curl zip make
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.18.0/protoc-3.18.0-linux-x86_64.zip \
    -o protoc3.zip
RUN unzip protoc3.zip -d protoc3 &&\
    mv protoc3/bin/protoc /bin/protoc &&\
    chmod a+rx /bin/protoc &&\
    mkdir -p $MAGMA_ROOT &&\
    mv ./protoc3/include/google $MAGMA_ROOT/google &&\
    rm -rf protoc3.zip protoc3
RUN pip3 install protobuf setuptools==49.6.0 grpcio==1.37.1 grpcio-tools==1.37.1
ENV PYTHONPATH=$MAGMA_ROOT:$MAGMA_ROOT/build/gen
ENV PYTHON_BUILD=$MAGMA_ROOT/build
COPY . $MAGMA_ROOT/
WORKDIR $MAGMA_ROOT/dp
RUN mkdir -p $PYTHON_BUILD && make protos

FROM python:3.9.2-slim-buster as standard-version

ENV PC_DIRECTORY=dp/cloud/python/magma/protocol_controller
COPY $PC_DIRECTORY/requirements.txt \
     /$PC_DIRECTORY/requirements.txt
RUN pip3 install --upgrade pip --no-cache-dir -r /$PC_DIRECTORY/requirements.txt

FROM standard-version as tests-version

COPY $PC_DIRECTORY/plugins/cbsd_sas/tests/requirements.txt \
     /$PC_DIRECTORY/plugins/cbsd_sas/tests/requirements.txt
WORKDIR /$PC_DIRECTORY
RUN pip3 install --upgrade pip --no-cache-dir -r plugins/cbsd_sas/tests/requirements.txt

FROM ${ENV}-version as final

ENV PYTHONPATH=/magma/build/gen:/dp/cloud/python
COPY $PC_DIRECTORY /$PC_DIRECTORY/
COPY --from=protos-generator /magma/build/gen /magma/build/gen
WORKDIR /$PC_DIRECTORY
EXPOSE 8000
CMD ["python", "service.py"]
