#!/bin/bash

OUTPUT_DIR=bin
WHICH_FBGO=$(command -v fbgo)
WHICH_GO=$(command -v go)
if [ -z "${WHICH_FBGO}" ]; then
    GO=${WHICH_GO}
else
    GO=${WHICH_FBGO}
fi

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
    gen_eap_client
    gen_authorization_client
    ${GO} build  .
}

function gen_eap_client {
    AAA_PROTOS_DIR=../../gateway/services/aaa/protos
    GRPC_GEN_OUTPUT_DIR=./modules/eap/methods/akamagma/protos
    mkdir -p ${GRPC_GEN_OUTPUT_DIR}
    protoc -I${AAA_PROTOS_DIR} --go_out=plugins=grpc,paths=source_relative:${GRPC_GEN_OUTPUT_DIR} ${AAA_PROTOS_DIR}/context.proto ${AAA_PROTOS_DIR}/eap.proto
}

function gen_authorization_client {
    AAA_PROTOS_DIR=../../gateway/services/aaa/protos
    GRPC_GEN_OUTPUT_DIR=./modules/coa/protos
    mkdir -p ${GRPC_GEN_OUTPUT_DIR}
    protoc -I${AAA_PROTOS_DIR} --go_out=plugins=grpc,paths=source_relative:${GRPC_GEN_OUTPUT_DIR} ${AAA_PROTOS_DIR}/context.proto ${AAA_PROTOS_DIR}/authorization.proto
}

function clean {
    rm -rf ${OUTPUT_DIR}
}

function start {
    build
    ./radius
}

function test {
    find . | grep _test\.go | sed 's/\(.*\)\/.*/\1/' | xargs -L1 "${GO}" "test"
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
test*)
	test
	;;
e2e*)
	e2e
	;;
gen_eap_client*)
    gen_eap_client
    ;;
lint*)
    lint
    ;;
gen_authorization_client*)
    gen_authorization_client
    ;;
gen*)
    gen_eap_client
    gen_authorization_client
    ;;
*)
	echo "usage: ./run.sh {build | clean | start | test | e2e | lint | gen_eap_client | gen_authorization_client}"
	;;
esac
