FROM python:3.8-alpine as builder

RUN addgroup -S linter && adduser -S -G linter linter
RUN apk add --no-cache git gcc musl-dev # temp

COPY lte/gateway/docker/python-precommit/requirements.txt /
RUN pip wheel -r /requirements.txt --wheel-dir=/wheelhouse/


FROM python:3.8-alpine as runner

RUN addgroup -S linter && adduser -S -G linter linter

RUN apk add --no-cache jq

COPY --from=builder /wheelhouse/ /wheelhouse/
RUN pip install /wheelhouse/*

USER linter
WORKDIR /code/
