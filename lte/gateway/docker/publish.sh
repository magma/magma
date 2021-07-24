

publish (){
# aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/z2g3r6f7
# docker tag control_proxy:latest public.ecr.aws/z2g3r6f7/control_proxy:latest
# docker push public.ecr.aws/z2g3r6f7/control_proxy:latest
  docker tag $1:latest public.ecr.aws/z2g3r6f7/$1:latest
  docker push public.ecr.aws/z2g3r6f7/$1:latest
}




publish mobilityd
# build enodebd
# build health
# build policydb
# build smsd
# build subscriberdb
# build ctraced
# build magmad
# build state
# build directoryd
# build pipelined
# build eventd
# build control_proxy
#
# # build container for C services
# cd ../../../
# docker build . -f lte/gateway/docker/mme/Dockerfile.ubuntu20.04 -t mme_builder:latest
# cd lte/gateway/docker
# docker build . -f services/build/Dockerfile.c -t cbuilder:latest
#
# # C services
# build mme
# build sctpd
# build sessiond
# build envoy_controller
