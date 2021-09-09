ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

COPY dp/cloud/python/db_service/config.py dp/cloud/python/db_service/models.py dp/cloud/python/db_service/session_manager.py dp/cloud/python/db_service/db_initialize.py dp/cloud/python/db_service/
COPY dp/cloud/python/db_service/tests dp/cloud/python/db_service/tests
COPY dp/cloud/python/mappings dp/cloud/python/mappings
COPY dp/protos dp/protos
COPY dp/cloud/python/test_runner dp/cloud/python/test_runner
COPY dp/cloud/python/fixtures dp/cloud/python/fixtures

RUN pip3 install --upgrade pip --no-cache-dir -r dp/cloud/python/test_runner/requirements.txt
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/requests.proto
RUN python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. dp/protos/active_mode.proto

WORKDIR dp/cloud/python/test_runner

FROM ${ENV}-version as final
ENV PYTHONPATH=$PYTHONPATH:/

ENTRYPOINT ["python"]