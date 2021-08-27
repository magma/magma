REGISTRY=$1
VERSION=latest
publish (){
  docker tag agw_gateway_$1:latest ${REGISTRY}gateway_$1:${VERSION}
  docker push ${REGISTRY}gateway_$1:${VERSION}
}

publish python
publish c
