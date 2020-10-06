# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM ubuntu:bionic as builder

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https

# Install the runtime deps.
RUN apt-get update && apt-get install -y bzr curl daemontools gcc

# Install Golang 1.13.
WORKDIR /usr/local
RUN curl https://facebookconnectivity.jfrog.io/artifactory/generic/go1.13.4.linux-amd64.tar.gz -O && \
    tar xf go1.13.4.linux-amd64.tar.gz && \
    cp -r go/bin/* /usr/local/bin/

ENV GOBIN /var/opt/magma/bin
ENV PATH ${PATH}:${GOBIN}
ENV GO111MODULE on
# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

# Copy just the go.mod and go.sum files to download the golang deps.
# This step allows us to cache the downloads, and prevents reaching out to
# the internet unless any of the go.mod or go.sum files are changed.
COPY cwf/cloud/go/go.* $MAGMA_ROOT/cwf/cloud/go/
COPY cwf/k8s/cwf_operator/go.*  $MAGMA_ROOT/cwf/k8s/cwf_operator/
COPY lte/cloud/go/go.* $MAGMA_ROOT/lte/cloud/go/
COPY feg/cloud/go/go.* $MAGMA_ROOT/feg/cloud/go/
COPY feg/cloud/go/protos/go.* $MAGMA_ROOT/feg/cloud/go/protos/
COPY orc8r/lib/go/go.* $MAGMA_ROOT/orc8r/lib/go/
COPY orc8r/lib/go/protos/go.* $MAGMA_ROOT/orc8r/lib/go/protos/
COPY orc8r/cloud/go/go.* $MAGMA_ROOT/orc8r/cloud/go/
COPY orc8r/gateway/go/go.* $MAGMA_ROOT/orc8r/gateway/go/
WORKDIR $MAGMA_ROOT/cwf/k8s/cwf_operator
RUN go mod download; exit 0

# Copy and build the go code.
COPY cwf/cloud/go $MAGMA_ROOT/cwf/cloud/go
COPY feg/cloud/go $MAGMA_ROOT/feg/cloud/go
COPY lte/cloud/go $MAGMA_ROOT/lte/cloud/go
COPY orc8r/lib/go $MAGMA_ROOT/orc8r/lib/go
COPY orc8r/cloud/go $MAGMA_ROOT/orc8r/cloud/go
COPY orc8r/gateway/go $MAGMA_ROOT/orc8r/gateway/go
COPY cwf/k8s/cwf_operator/cmd $MAGMA_ROOT/cwf/k8s/cwf_operator/cmd
COPY cwf/k8s/cwf_operator/pkg $MAGMA_ROOT/cwf/k8s/cwf_operator/pkg
COPY cwf/k8s/cwf_operator/version $MAGMA_ROOT/cwf/k8s/cwf_operator/version
RUN go install ./cmd/manager/

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:bionic as cwf_operator

ENV OPERATOR=/usr/local/bin/cwf-operator \
    USER_UID=1001 \
    USER_NAME=cwf-operator

# Install the runtime deps.
RUN apt-get update && apt-get install -y curl openssh-server fping

RUN curl -L "https://github.com/Mirantis/virtlet/releases/download/v1.5.1/virtletctl" -o /usr/bin/virtletctl
RUN chmod 755 /usr/bin/virtletctl

# install operator binary
COPY --from=builder /var/opt/magma/bin/manager ${OPERATOR}

COPY cwf/k8s/cwf_operator/docker/entrypoint /usr/local/bin
COPY cwf/k8s/cwf_operator/docker/user_setup /usr/local/bin
COPY cwf/gateway/helm/cwf-kubevirt /helm-charts/cwf
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}
