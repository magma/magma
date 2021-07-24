REGISTRY=$1
VERSION=latest
publish (){
  docker tag $1:latest ${REGISTRY}$1:${VERSION}
  docker push ${REGISTRY}$1:${VERSION}
}

publish mobilityd
publish enodebd
publish health
publish policydb
publish smsd
publish subscriberdb
publish ctraced
publish magmad
publish state
publish directoryd
publish pipelined
publish eventd
publish control_proxy

# C services
publish mme
publish sctpd
publish sessiond
publish envoy_controller
