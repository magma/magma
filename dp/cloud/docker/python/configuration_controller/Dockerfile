ARG ENV=standard
FROM python:3.9.2-slim-buster as standard-version

ENV CC_DIRECTORY=dp/cloud/python/magma/configuration_controller
ENV DB_DIRECTORY=dp/cloud/python/magma/db_service
ENV FC_DIRECTORY=dp/cloud/python/magma/fluentd_client
ENV MC_DIRECTORY=dp/cloud/python/magma/metricsd_client
ENV GW_COMMON_DIRECTORY=orc8r/gateway/python/magma/common
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
COPY $FC_DIRECTORY /$FC_DIRECTORY/
COPY $MC_DIRECTORY /$MC_DIRECTORY/
COPY build/gen /build/gen
COPY $GW_COMMON_DIRECTORY/metrics_export.py /$GW_COMMON_DIRECTORY/metrics_export.py
COPY $GW_COMMON_DIRECTORY/__init__.py /$GW_COMMON_DIRECTORY/__init__.py
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/
ENV PYTHONPATH=/:/dp/cloud/python:/orc8r/gateway/python:/build/gen
ENTRYPOINT ["python"]
CMD ["run.py"]
