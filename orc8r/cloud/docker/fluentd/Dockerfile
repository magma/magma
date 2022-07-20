FROM fluent/fluentd:v1.14.6-debian-1.0
USER root
RUN gem install \
    elasticsearch:7.13.0 \
    fluent-plugin-elasticsearch:5.2.1 \
    fluent-plugin-multi-format-parser:1.0.0 \
    --no-document
USER fluent
