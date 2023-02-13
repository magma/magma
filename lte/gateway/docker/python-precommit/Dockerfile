FROM python:3.9.10-alpine

RUN addgroup -S linter && adduser -S -G linter linter

# Installing wemake-python-styleguide==0.17.0 is not
# working with a requirements.in file.
# hadolint ignore=DL3013
RUN pip install --no-cache-dir wemake-python-styleguide==0.17.0 \
                add-trailing-comma \
                isort \
                autopep8

USER linter
WORKDIR /code/
