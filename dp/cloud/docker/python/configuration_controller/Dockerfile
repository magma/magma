ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

COPY dp/cloud/python/configuration_controller dp/cloud/python/configuration_controller
COPY dp/cloud/python/db_service/config.py \
    dp/cloud/python/db_service/models.py \
    dp/cloud/python/db_service/session_manager.py \
    dp/cloud/python/db_service/session_manager.py \
    dp/cloud/python/db_service/
COPY dp/cloud/python/mappings dp/cloud/python/mappings
WORKDIR dp/cloud/python/configuration_controller
RUN pip3 install --upgrade pip --no-cache-dir -r requirements.txt

FROM standard-version as tests-version
RUN pip3 install --upgrade pip --no-cache-dir -r tests/requirements.txt

FROM ${ENV}-version as final
ENV PYTHONPATH=$PYTHONPATH:/
ENTRYPOINT ["python"]
CMD ["run.py"]
