ARG ENV=standard
FROM python:3.9.2-slim-buster as builder

COPY dp/cloud/python/magma/db_service/requirements.txt /dp/cloud/python/magma/db_service/requirements.txt
WORKDIR /dp/cloud/python/magma/db_service/migrations
RUN pip3 install --upgrade pip --no-cache-dir -r ../requirements.txt

COPY dp/cloud/python/magma/db_service /dp/cloud/python/magma/db_service/
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/

FROM builder as final
ENV PYTHONPATH=/:/dp/cloud/python
ENV ALEMBIC_CONFIG=./alembic.ini

ENTRYPOINT ["python"]
CMD ["../db_initialize.py"]
