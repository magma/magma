ARG ENV=standard
FROM python:3.9.2-slim-buster as protos-generator

RUN apt update && apt install -y curl zip make
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.18.0/protoc-3.18.0-linux-x86_64.zip \
    -o protoc3.zip
RUN unzip protoc3.zip -d protoc3 &&\
    mv protoc3/bin/protoc /bin/protoc &&\
    chmod a+rx /bin/protoc &&\
    mkdir -p /magma &&\
    mv ./protoc3/include/google /magma/google &&\
    rm -rf protoc3.zip protoc3
RUN pip3 install protobuf setuptools==49.6.0 grpcio==1.37.1 grpcio-tools==1.37.1
ENV MAGMA_ROOT=/magma
ENV PYTHONPATH=$MAGMA_ROOT:$MAGMA_ROOT/build/gen
ENV PYTHON_BUILD=$MAGMA_ROOT/build
COPY . /magma/
WORKDIR /magma/dp
RUN mkdir -p $PYTHON_BUILD && make protos

FROM python:3.9.2-slim-buster as standard-version

ENV RC_DIRECTORY=dp/cloud/python/magma/radio_controller
ENV DB_DIRECTORY=dp/cloud/python/magma/db_service
COPY $RC_DIRECTORY/requirements.txt \
     /$RC_DIRECTORY/requirements.txt
RUN pip3 install --upgrade pip --no-cache-dir -r /$RC_DIRECTORY/requirements.txt

FROM standard-version as tests-version

COPY $RC_DIRECTORY/tests/requirements.txt \
     /$RC_DIRECTORY/tests/requirements.txt
WORKDIR /$RC_DIRECTORY
RUN pip3 install --upgrade pip --no-cache-dir -r tests/requirements.txt

FROM ${ENV}-version as final

ENV PYTHONPATH=/magma/build/gen:/dp/cloud/python
COPY $RC_DIRECTORY /$RC_DIRECTORY/
COPY $DB_DIRECTORY /$DB_DIRECTORY/
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/
COPY --from=protos-generator /magma/build/gen /magma/build/gen
WORKDIR /$RC_DIRECTORY
EXPOSE 50053
ENTRYPOINT ["python"]
CMD ["run.py"]
