FROM python:3.9-slim-buster
ARG ENV=prod
ARG HELM_VERSION="3.5.1"
ARG TERRAFORM_VERSION="0.15.5"
ARG KUBECTL_VERSION="1.20.2"
ARG ANSIBLE_VERSION="3.0.0"

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
            git \
            wget \
            unzip \
            curl \
            vim \
            jq \
            gcc \
            procps && \
    apt-get clean -y && \
    apt-get autoclean -y && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /var/cache/apt/archives/*

RUN pip3 install --no-cache-dir \
                        boto \
                        boto3 \
                        kubernetes \
                        ansible==${ANSIBLE_VERSION} \
                        prettytable \
                        requests \
                        docker \
                        pyOpenSSL \
                        unittest2 \
                        colorama \
                        dataclasses-json \
                        pytest

WORKDIR /root/download

# Install aws cli
RUN wget "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -O "/root/download/awscli2.zip" && \
    unzip awscli2.zip && \
    ./aws/install

# Install aws iam authenticator
RUN wget https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.9/2020-11-02/bin/linux/amd64/aws-iam-authenticator -O /usr/local/bin/aws-iam-authenticator && \
    chmod +x /usr/local/bin/aws-iam-authenticator

# Install helm
RUN mkdir helm3 && \
    curl -SsL --retry 5 "https://get.helm.sh/helm-v$HELM_VERSION-linux-amd64.tar.gz" | tar xz -C ./helm3 && \
    cp /root/download/helm3/linux-amd64/helm /usr/local/bin/helm

# Install terraform
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform\_${TERRAFORM_VERSION}\_linux_amd64.zip && \
    unzip ./terraform\_${TERRAFORM_VERSION}\_linux_amd64.zip -d terraform14_cli && \
    cp /root/download/terraform14_cli/terraform /usr/local/bin/terraform && \
    wget https://storage.googleapis.com/kubernetes-release/release/v$KUBECTL_VERSION/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    rm -rf /root/download/*


# Install go if we are building testframework image
WORKDIR /usr/local
RUN if [ "$ENV" = "testframework" ] ; then curl https://facebookconnectivity.jfrog.io/artifactory/generic/go1.13.4.linux-amd64.tar.gz -O && \
    tar xf go1.13.4.linux-amd64.tar.gz && \
    cp -r go/bin/* /usr/local/bin/ ; fi

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
ENV ANSIBLE_CONFIG=/root/.ansible.cfg

# Uncomment to install avalanche
# WORKDIR /root/download
# RUN if [ "$ENV" = "testframework" ] ; then go get -v github.com/open-fresh/avalanche/cmd/... && \
#     mv /root/go/bin/cmd /usr/local/bin/avalanche && \
#     chmod +x /usr/local/bin/avalanche ; fi

COPY root/ /root/

# Install the orc8r cli (orcl)
WORKDIR /root/scripts
RUN pip3 install .

# Set these environment variables to ensure aws configuration remains
# in the deployment directory even when container is removed
ENV AWS_CONFIG_FILE=/root/project/.aws/config
ENV AWS_SHARED_CREDENTIALS_FILE=/root/project/.aws/credentials
ENV AWS_DEFAULT_OUTPUT=json

WORKDIR /root/project
CMD ["/bin/bash"]