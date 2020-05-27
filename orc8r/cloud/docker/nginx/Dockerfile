FROM nginx:1.17

RUN apt-get update && \
  apt-get install -y python3-pip daemontools
RUN pip3 install PyYAML jinja2

RUN mkdir -p /var/opt/magma/envdir
COPY configs /etc/magma/configs

WORKDIR /tmp

ARG NGINX_FILES=src/magma/orc8r/cloud/docker/nginx
ARG PROXY_FILES=src/magma/orc8r/cloud/docker/proxy

COPY ${NGINX_FILES}/templates /etc/magma/templates
COPY ${NGINX_FILES}/generate_nginx_configs.py /usr/local/bin/generate_nginx_configs.py
COPY ${PROXY_FILES}/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

CMD /usr/local/bin/create_test_proxy_certs && envdir /var/opt/magma/envdir /usr/local/bin/generate_nginx_configs.py && exec nginx -g 'daemon off;'
