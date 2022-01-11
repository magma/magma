ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

ENV CC_DIRECTORY=dp/cloud/python/magma/configuration_controller
ENV DB_DIRECTORY=dp/cloud/python/magma/db_service

COPY $CC_DIRECTORY/requirements.txt \
     /$CC_DIRECTORY/requirements.txt
WORKDIR /$CC_DIRECTORY
RUN pip3 install --upgrade pip --no-cache-dir -r requirements.txt

#FROM standard-version as tests-version
#COPY $CC_DIRECTORY/tests/requirements.txt \
#     /$CC_DIRECTORY/tests/requirements.txt
#RUN pip3 install --upgrade pip --no-cache-dir -r tests/requirements.txt

# hadolint ignore=DL3006
FROM ${ENV}-version as final
COPY $CC_DIRECTORY /$CC_DIRECTORY/
COPY $DB_DIRECTORY /$DB_DIRECTORY/
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/
ENV PYTHONPATH=/:/dp/cloud/python
ENTRYPOINT ["python"]
CMD ["run.py"]
