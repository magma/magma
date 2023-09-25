# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM python:3.9-slim-buster
ARG ENV=prod
ARG TERRAFORM_VERSION="0.15.5"
ARG KUBECTL_VERSION="1.20.2"
ARG ANSIBLE_VERSION="3.0.0"

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
            git \
            wget \
            unzip \
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
ARG HELM_VERSION="3.5.1"
RUN HELM_TARBALL="helm-v${HELM_VERSION}-linux-amd64.tar.gz" \
 && wget --no-verbose https://get.helm.sh/${HELM_TARBALL} \
 && tar --strip-components=1 -C /usr/local/bin -xzf ${HELM_TARBALL} linux-amd64/helm \
 && rm ${HELM_TARBALL}

# Install terraform
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform\_${TERRAFORM_VERSION}\_linux_amd64.zip && \
    unzip ./terraform\_${TERRAFORM_VERSION}\_linux_amd64.zip -d terraform14_cli && \
    cp /root/download/terraform14_cli/terraform /usr/local/bin/terraform && \
    wget https://storage.googleapis.com/kubernetes-release/release/v$KUBECTL_VERSION/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    rm -rf /root/download/*

# Install go if we are building testframework image
ARG GOLANG_VERSION="1.20.1"
RUN if [ "$ENV" = "testframework" ]; \
    then \
        GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
     && wget --no-verbose https://linuxfoundation.jfrog.io/artifactory/magma-blob/${GO_TARBALL} \
     && tar --strip-components=2 -C /usr/local/bin -xzf ${GO_TARBALL} go/bin/go go/bin/gofmt \
     && rm ${GO_TARBALL}; \
    fi

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
