FROM fluent/fluentd:v1.14.6-debian-1.0

USER root
WORKDIR /home/fluent
ENV PATH /fluentd/vendor/bundle/ruby/2.7.0/bin:$PATH
ENV GEM_PATH /fluentd/vendor/bundle/ruby/2.7.0
ENV GEM_HOME /fluentd/vendor/bundle/ruby/2.7.0
# skip runtime bundler installation
ENV FLUENTD_DISABLE_BUNDLER_INJECTION 1

COPY Gemfile* /fluentd/
RUN buildDeps="sudo make gcc g++ libc-dev libffi-dev" \
  && apt-get update \
  && apt-get upgrade -y \
  && apt-get install \
    -y --no-install-recommends \
    $buildDeps net-tools \
  && gem install bundler --version 1.16.2 \
  && bundle config silence_root_warning true \
  && bundle install --gemfile=/fluentd/Gemfile --path=/fluentd/vendor/bundle \
  && SUDO_FORCE_REMOVE=yes \
    apt-get purge -y --auto-remove \
                -o APT::AutoRemove::RecommendsImportant=false \
                $buildDeps \
  && rm -rf /var/lib/apt/lists/* \
  && gem sources --clear-all \
  && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

# Copy configuration files
RUN touch /fluentd/etc/disable.conf

# Copy plugins
COPY plugins /fluentd/plugins/
COPY entrypoint.sh /fluentd/entrypoint.sh

# Environment variables
ENV FLUENTD_OPT=""
ENV FLUENTD_CONF="fluent.conf"

# Overwrite ENTRYPOINT to run fluentd as root for /var/log / /var/lib
ENTRYPOINT ["tini", "--", "/fluentd/entrypoint.sh"]
