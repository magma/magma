#!/bin/bash
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

function lint {
    LINTER=$(which golangci-lint)
    if [ -z "$LINTER" ]; then
        echo "-> installing golangci " && \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sudo sh -s -- -b /usr/sbin/ v1.51.2;
    fi

    golangci-lint run -c "$MAGMA_ROOT"/.golangci.yml

    if [ $? -ne 0 ]; then
        echo Lint failed.
        exit 1
    fi
}

function build {
    go build .
}

function gen {
    go generate ./...
}

function clean {
    rm ./radius
}

function start {
    build
    ./radius
}

function pretty {
    build
    ./radius 2>&1 |  zap-pretty
}

function test {
    # Run tests sequentially to avoid radius server cross-talk between tests.
    gotestsum -- -p 1 ./...
}

function e2e {
    pushd ../integration/lb/sim || exit
    echo Building E2E test containers for LB with simulator
    docker compose --compatibility build
    echo Starting E2E test containers for LB with simulator
    echo Automatically terminating after 60 secs...
    docker compose --compatibility up &
    sleep 60s
    docker compose down
    popd || exit
}

case $1 in
build*)
	build
	;;
clean*)
    clean
    ;;
start*)
	start
	;;
pretty*)
        pretty
        ;;
test*)
	test
	;;
e2e*)
	e2e
	;;
lint*)
    lint
    ;;
gen*)
    gen
    ;;
*)
	echo "usage: ./run.sh {build | clean | start | test | e2e | lint | gen}"
	;;
esac
