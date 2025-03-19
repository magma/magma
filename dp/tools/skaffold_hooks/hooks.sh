#!/bin/bash

# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MAGMA_ROOT=$(realpath .)
CERTS_DIR=${MAGMA_ROOT}/orc8r/cloud/helm/orc8r/charts/secrets/certs
export MAGMA_ROOT
export CERTS_DIR


build_controller() {
    cd "$MAGMA_ROOT/orc8r/cloud/docker" || exit 1
    python3 build.py -b controller &&
    docker tag orc8r_controller:latest "$IMAGE"
    if $PUSH_IMAGE; then
        docker push "$IMAGE"
    fi
}

build_nginx() {
    cd "$MAGMA_ROOT/orc8r/cloud/docker" || exit 1
    python3 build.py -b nginx &&
    docker tag orc8r_nginx:latest "$IMAGE"
    if $PUSH_IMAGE; then
        docker push "$IMAGE"
    fi
}

build_magmalte() {
    cd "$MAGMA_ROOT/nms" || exit 1
    COMPOSE_PROJECT_NAME=magmalte docker compose --compatibility build magmalte &&
    docker tag magmalte_magmalte:latest "$IMAGE"
    if $PUSH_IMAGE; then
        docker push "$IMAGE"
    fi
}

create_secrets() {

  if [ ! -d "$CERTS_DIR" ]; then
      mkdir -p "$CERTS_DIR"
      gen_certs
      mv -- *.key *.pem *crt "$CERTS_DIR"
  fi
  cd "$MAGMA_ROOT/orc8r/cloud/deploy/scripts" || exit 1

  if [ ! -v "$REMOTE_NAMESPACE" ]; then
      apply_secrets
  fi

}

gen_certs() {

  cd "$MAGMA_ROOT/orc8r/cloud/deploy/scripts" || exit 1
  ./self_sign_certs.sh localhost
  ./create_application_certs.sh localhost
}

apply_secrets() {

  cd "$MAGMA_ROOT/orc8r/cloud/helm/orc8r" || exit 1
  helm template orc8r charts/secrets \
    --namespace default \
    -f "${MAGMA_ROOT}/dp/cloud/helm/dp/charts/domain-proxy/examples/orc8r_secrets_values.yaml" \
    --set-string secret.certs.enabled=true \
    --set-file secret.certs.files."rootCA\.pem=${CERTS_DIR}/rootCA.pem" \
    --set-file secret.certs.files."bootstrapper\.key=${CERTS_DIR}/bootstrapper.key" \
    --set-file secret.certs.files."controller\.crt=${CERTS_DIR}/controller.crt" \
    --set-file secret.certs.files."controller\.key=${CERTS_DIR}/controller.key" \
    --set-file secret.certs.files."admin_operator\.pem=${CERTS_DIR}/admin_operator.pem" \
    --set-file secret.certs.files."admin_operator\.key\.pem=${CERTS_DIR}/admin_operator.key.pem" \
    --set-file secret.certs.files."certifier\.pem=${CERTS_DIR}/certifier.pem" \
    --set-file secret.certs.files."certifier\.key=${CERTS_DIR}/certifier.key" \
    --set-file secret.certs.files."fluentd\.pem=${CERTS_DIR}/fluentd.pem" \
    --set-file secret.certs.files."fluentd\.key=${CERTS_DIR}/fluentd.key" \
    --set=docker.registry="$DOCKER_REGISTRY" \
    --set=docker.username="$DOCKER_USERNAME" \
    --set=docker.password="$DOCKER_PASSWORD" |
    kubectl apply -f -
    #--set-file secret.certs.files."nms_nginx\.pem"=${CERTS_DIR}/nms_nginx.pem \
    #--set-file secret.certs.files."nms_nginx\.key\.pem"=${CERTS_DIR}/nms_nginx.key \

  cd "$MAGMA_ROOT/dp/cloud/helm/dp/charts" || exit 1
  helm template domain-proxy -s templates/fluentd/secrets.yaml \
    -f domain-proxy/examples/minikube_values.yaml \
    --namespace default \
    --release-name domain-proxy \
    --set-string dp.fluentd.secret.certs.enabled=true \
    --set-file dp.fluentd.secret.certs.files."ca\.pem=${CERTS_DIR}/certifier.pem" \
    --set-file dp.fluentd.secret.certs.files."fluentd\.pem=${CERTS_DIR}/dp_fluentd.pem" \
    --set-file dp.fluentd.secret.certs.files."fluentd\.key=${CERTS_DIR}/dp_fluentd.key" |
  kubectl apply -f -

}

deploy_fluentd_forwarder() {
  cd "$MAGMA_ROOT/orc8r/cloud/helm/orc8r" || exit 1
  helm dependency build
  helm template . \
    -s charts/logging/templates/fluentd-forward.deployment.yaml \
    -s charts/logging/templates/fluentd-forward.service.yaml \
    -s charts/logging/templates/fluentd-forward.configmap.yaml \
    --namespace default --name-template orc8r \
    -f ../../../../dp/cloud/helm/dp/charts/domain-proxy/examples/orc8r_minikube_values.yaml |
  kubectl apply -f -

}

cleanup() {
  rm -rf "$CERTS_DIR"
}

"$@"
