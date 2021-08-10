################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ENV TZ=America/New_York
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        apt-transport-https \
        apt-utils \
        build-essential \
        ca-certificates \
        curl \
        gcc \
        git \
        gnupg2 \
        g++ \
        python-dev \
        zip \
        unzip \
        vim \
        wget

WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.6.1/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

WORKDIR /magma
