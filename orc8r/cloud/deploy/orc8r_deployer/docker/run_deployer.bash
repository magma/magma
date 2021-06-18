#!/usr/bin/env bash
set -euf -o pipefail

print_help()
{
    echo "
./run_deployer runs the orc8r deployer container
orc8r deployer contains scripts which enable user to configure, run prechecks,
install, upgrade, verify and cleanup an orc8r deployment

Usage: run_deployer [-deploy-dir|-root-dir|-build|-h]
options:
-h           Print this Help
--deploy-dir  deployment dir containing configs and secrets (mandatory)
--root-dir    magma root directory
--build       build the deployer container
--build-testframework build the deployer container with go library and other testframework utilities
--test        'check_all' or any specific test function[run_unit_tests,check_helmcharts_insync, check_tfvars_insync ]
example: ./run_deployer.bash --deploy-dir ~/orc8r_15_deployment"
}

run_unit_tests()
{
    echo "Running orcl container unit tests"
    docker run -it \
        --entrypoint /root/scripts/cli/configlib_test.py \
        -v "${DEPLOY_WORKDIR}":/root/project \
        -v "${MAGMA_ROOT}":/root/magma \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --rm orc8r_deployer:latest
}

check_helmcharts_insync()
{
    echo "Checking if helm charts are in sync"
    docker run -it \
        --entrypoint /root/scripts/test_helm_charts_sync.py \
        -v "${DEPLOY_WORKDIR}":/root/project \
        -v "${MAGMA_ROOT}":/root/magma \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --rm orc8r_deployer:latest
}

check_tfvars_insync()
{
    echo "Checking tf vars are in sync"
    docker run -it \
        --entrypoint /root/scripts/test_vars_sync.py \
        -v "${DEPLOY_WORKDIR}":/root/project \
        -v "${MAGMA_ROOT}":/root/magma \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --rm orc8r_deployer:latest
}

check_all() {
    run_unit_tests
    check_helmcharts_insync
    check_tfvars_insync
}

if (( $# < 2 )); then
    print_help
    exit 1
fi

DOCKER_BUILD=false
DOCKER_BUILD_TESTFRAMEWORK=false
DEPLOY_WORKDIR=
MAGMA_ROOT=
TEST_TO_RUN=

while [ -n "${1-}" ]; do
	case "$1" in
	--deploy-dir)
        DEPLOY_WORKDIR="$2"
        shift
        ;;
	--root-dir)
        MAGMA_ROOT="$2"
        shift
        ;;
    --build)
        DOCKER_BUILD=true
        ;;
    --build-testframework)
        DOCKER_BUILD_TESTFRAMEWORK=true
        ;;        
    --test)
        TEST_TO_RUN="$2"
        shift
        ;;
    -h)
        print_help
        exit;;
	*) echo "Option $1 not recognized"
        print_help
        exit;;
	esac
    shift
done

echo "Build $DOCKER_BUILD"

echo "Deploy workddir $DEPLOY_WORKDIR"
if [[ ! -d $DEPLOY_WORKDIR ]]; then
    echo "${DEPLOY_WORKDIR} does not exist. Creating a new directory"
    mkdir -p $DEPLOY_WORKDIR
fi

if [ -z $MAGMA_ROOT ]; then
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
    PATTERN='\/orc8r\/cloud\/deploy\/orc8r_deployer\/docker'
    MAGMA_ROOT=${SCRIPT_DIR%%$PATTERN}
    echo "Warning: inferring magma root($MAGMA_ROOT) from run_deployer script directory"
fi

if $DOCKER_BUILD; then
    docker build -t orc8r_deployer:latest .
fi

if $DOCKER_BUILD_TESTFRAMEWORK; then
    docker build  --build-arg ENV=testframework -t orc8r_deployer:latest .
fi

if declare -F "$TEST_TO_RUN"; then
    $TEST_TO_RUN
else
    docker run -it \
        -v "${DEPLOY_WORKDIR}":/root/project \
        -v "${MAGMA_ROOT}":/root/magma \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --rm orc8r_deployer:latest bash
fi
