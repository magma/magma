ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

COPY dp/cloud/python/db_service/config.py \
    dp/cloud/python/db_service/models.py \
    dp/cloud/python/db_service/session_manager.py \
    dp/cloud/python/db_service/
COPY dp/cloud/python/mappings dp/cloud/python/mappings
COPY dp/cloud/python/radio_controller dp/cloud/python/radio_controller
COPY dp/protos dp/protos

RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/radio_controller/requirements.txt
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/requests.proto
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/active_mode.proto
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/common.proto

FROM standard-version as tests-version
COPY /dp/cloud/python/db_service/db_initialize.py dp/cloud/python/db_service/
COPY dp/cloud/python/db_service/tests dp/cloud/python/db_service/tests
RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/radio_controller/tests/requirements.txt
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/requests.proto
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/active_mode.proto
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/common.proto

FROM ${ENV}-version as final

WORKDIR dp/cloud/python/radio_controller

ENV PYTHONPATH=$PYTHONPATH:/

EXPOSE 50053

ENTRYPOINT ["python"]
CMD ["run.py"]
