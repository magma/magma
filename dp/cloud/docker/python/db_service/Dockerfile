ARG ENV=standard
FROM python:3.9.2-slim-buster as builder

COPY dp/cloud/python/db_service dp/cloud/python/db_service
COPY dp/cloud/python/mappings dp/cloud/python/mappings
WORKDIR dp/cloud/python/db_service/migrations
RUN pip3 install --upgrade pip --no-cache-dir -r ../requirements.txt

FROM builder as final
ENV PYTHONPATH=$PYTHONPATH:/
ENV ALEMBIC_CONFIG=./alembic.ini

ENTRYPOINT ["python"]
CMD ["../db_initialize.py"]
