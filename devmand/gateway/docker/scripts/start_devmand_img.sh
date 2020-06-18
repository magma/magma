#!/bin/bash
dirname="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${dirname}

DEVMANDDOCKERNET="devmanddevelnet"
DEVMANDIP="172.8.0.85"
DEVMANDSUBNET="172.8.0.0/16"

docker network inspect ${DEVMANDDOCKERNET}

if [ $? -eq 0 ]
then
    echo "docker network exists"
else
    docker network create --subnet=${DEVMANDSUBNET} ${DEVMANDDOCKERNET}
fi

docker run -d -h devmanddevel --net ${DEVMANDDOCKERNET}  \
      --ip ${DEVMANDIP} \
      --name 85 \
      -v "$(realpath ../../):/cache/devmand/repo:rw" \
      -v "$(realpath ~/cache_devmand_build):/cache/devmand/build:rw" \
      --entrypoint /bin/bash \
      "facebookconnectivity-southpoll-dev-docker.jfrog.io/devmand" \
      -c 'mkdir -p /run/sshd && /usr/sbin/sshd && bash -c "sleep infinity && ls"'
