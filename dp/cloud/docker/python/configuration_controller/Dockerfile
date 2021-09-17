ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

COPY dp/cloud/python/magma/configuration_controller dp/cloud/python/magma/configuration_controller
COPY dp/cloud/python/magma/db_service/config.py \
    dp/cloud/python/magma/db_service/models.py \
    dp/cloud/python/magma/db_service/session_manager.py \
    dp/cloud/python/magma/db_service/session_manager.py \
    dp/cloud/python/magma/db_service/
COPY dp/cloud/python/magma/mappings dp/cloud/python/magma/mappings
WORKDIR dp/cloud/python/magma/configuration_controller
RUN pip3 install --upgrade pip --no-cache-dir -r requirements.txt

FROM standard-version as tests-version
RUN pip3 install --upgrade pip --no-cache-dir -r tests/requirements.txt

FROM ${ENV}-version as final
ENV PYTHONPATH=/:/dp/cloud/python
ENTRYPOINT ["python"]
CMD ["run.py"]
