#!/bin/bash 

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

docker run -h devmanddevel --net ${DEVMANDDOCKERNET} --ip ${DEVMANDIP} -it "facebookconnectivity-southpoll-dev-docker.jfrog.io/devmand" /bin/bash -c "/usr/sbin/service ssh start && bash"
