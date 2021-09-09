ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

COPY dp/cloud/python/protocol_controller dp/cloud/python/protocol_controller
COPY dp/protos dp/protos

RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/protocol_controller/requirements.txt
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/requests.proto

FROM standard-version as tests-version
RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/protocol_controller/plugins/cbsd_sas/tests/requirements.txt
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/requests.proto

FROM ${ENV}-version as final

WORKDIR dp/cloud/python/protocol_controller

ENV PYTHONPATH=$PYTHONPATH:/

EXPOSE 8000

CMD ["python", "service.py"]
