FROM alpine:3.11 as build

COPY docker/user-grafana/orc8r-datasource /orc8r-datasource

# Build datasource plugin
WORKDIR /orc8r-datasource
RUN apk add npm
RUN npm install
RUN npm run-script build

FROM grafana/grafana

ARG PLUGIN_PATH="/grafana-plugins"
ENV GF_PATHS_PLUGINS=${PLUGIN_PATH}

COPY --from=build /orc8r-datasource ${PLUGIN_PATH}/orc8r-datasource

ADD docker/user-grafana/grafana.ini /etc/grafana/grafana.ini
