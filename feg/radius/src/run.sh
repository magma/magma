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
        echo ERROR :: Could not find 'golangci-lint' executable
        echo To install, run 'brew install golangci/tap/golangci-lint'
        exit 1
    fi

    golangci-lint run -c ../../ci/.golangci.yml

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
    find . | grep _test\.go | sed 's/\(.*\)\/.*/\1/' | xargs -L1 go "test"
}

function e2e {
    pushd ../integration/lb/sim || exit
    echo Building E2E test containers for LB with simulator
    docker-compose build
    echo Starting E2E test containers for LB with simulator
    echo Automatically terminating after 60 secs...
    docker-compose up &
    sleep 60s
    docker-compose down
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
