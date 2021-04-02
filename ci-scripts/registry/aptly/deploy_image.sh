#!/bin/bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

CMD=$(basename "$0")

if [ "${USE_ENV_HOME}" != "true" ]; then
    # run as
    # USE_ENV_HOME=true ./pull_images
    # to run manually as another user
    export HOME=/home/ubuntu
fi

function logwrapper() {
    # avoid feeding unbounded input to this in the foreground
    if command -v logger >/dev/null 2>&1; then
	logger -t "${CMD}"
    else
    	cat > /dev/null
    fi
}

exec &> >(tee >(logwrapper))

echo "checking for new packages"

export PATH=${HOME}/bin:${HOME}/.local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin

APTLY_SOCK=${HOME}/docker/run/aptly/aptly.sock

if [ -z "${MAGMA_VERSION}" ]; then
    export MAGMA_VERSION="test"
fi

if [ -z "${DOCKER_DIR}" ]; then
    export DOCKER_DIR=${HOME}/docker
fi

aws_ls=$(aws s3 ls s3://magma-images/gateway/ | sed 's/.* //g')

# allow to specify a build on the command line
override_build=$1

most_recent_build=""
complete_file_name=""
for line in $aws_ls
do
  if [[ $line = *".deps.tar.gz" ]]; then
    most_recent_build=${line%".deps.tar.gz"}
    complete_file_name=$line


    # allow to specify a build on the command line, empty build ignored
    if [[ "x$override_build" != "x" && "$most_recent_build" = "$override_build" ]]; then
        break
    fi
  fi
done

if [[ $most_recent_build = "" ]]; then
  exit
fi

# check if the latest build is already in the stretch-${MAGMA_VERSION} repo
stretch_test_ls=$(curl -s --unix-socket ${APTLY_SOCK} "http://localhost/api/repos/stretch-${MAGMA_VERSION}/packages?q=Name+%28magma%29&format=details' | jq -r '.[].Version")

for line in $stretch_test_ls
do
  if [[ $line == "$most_recent_build" ]]; then
    exit
  fi
done

# create directory and download image from s3
temp_dir=temp_ci_packages
temp_path=${HOME}/${temp_dir}
mkdir $temp_path
function cleanup {
  rm -r $temp_path
}
trap cleanup EXIT

aws s3 cp s3://magma-images/gateway/"$complete_file_name" "$temp_path"/
tar -xvzf "$temp_path"/"$complete_file_name" -C "$temp_path"/
rm -f "$temp_path"/"$complete_file_name"

cd ${DOCKER_DIR} || exit 1

dc_exec_user=aptly-user

function dc_exec() {
    # use paths from docker container in arg list for this function
    # not host paths
    pushd ${DOCKER_DIR} >/dev/null 2>&1 || exit 1
    docker-compose exec -T -u ${dc_exec_user} aptly "$@"
    popd >/dev/null 2>&1 || exit 1
}

destdir=/home/aptly-user/upload

dc_exec mkdir -p ${destdir}
docker cp ${temp_path} "$(docker-compose ps -q aptly)":${destdir}
dc_exec_user=root dc_exec chown -R aptly-user:aptly-user ${destdir}
dc_exec aptly repo create -architectures=amd64 -distribution=stretch-${MAGMA_VERSION} stretch-${MAGMA_VERSION}
dc_exec aptly repo add -remove-files stretch-${MAGMA_VERSION} "${destdir}/${temp_dir}"
dc_exec rm -rf "${destdir}/${temp_dir}"
dc_exec aptly publish repo stretch-${MAGMA_VERSION}
dc_exec aptly publish update stretch-${MAGMA_VERSION}
dc_exec aptly db cleanup
