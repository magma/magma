#!/bin/bash

# This script regenerates cpp files from gRPC proto file,
# assuming docker/scripts/start_devmand_image.sh was executed and container is running.
# Passing argument 'clean' will clean the destination directory.

dirname="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
protoDir="/cache/devmand/repo/src/devmand/channels/cli/plugin/proto"
#protoFile="${protoDir}/ReaderPlugin.proto"
#protoFile2="${protoDir}/PluginRegistration.proto"
protocppRelativeDir="/../protocpp"
protocppDir="${protoDir}/${protocppRelativeDir}"
hostProtocppDir="${dirname}/${protocppRelativeDir}"
protoFileArray=(${dirname}/*.proto)


if [ "$1" == "clean" ]; then
    # rm -rf
    docker exec --user $(id -u):$(id -g) 85 rm -rf "${protocppDir}"
    shift
fi
docker exec --user $(id -u):$(id -g) 85 mkdir -p "${protocppDir}"


for ((i=0; i<${#protoFileArray[@]}; i++)); do
    fileName=$(basename -- "${protoFileArray[$i]}")
    protoFile="${protoDir}/${fileName}"
    docker exec --user $(id -u):$(id -g) 85 protoc -I "${protoDir}" --cpp_out="${protocppDir}" "${protoFile}"
    docker exec --user $(id -u):$(id -g) 85 bash -c \
      "protoc -I \"${protoDir}\" --grpc_out=\"${protocppDir}\" --plugin=protoc-gen-grpc=\`which grpc_cpp_plugin\` \"${protoFile}\""
done


# rename cc to cpp
rename 's/\.cc$/\.cpp/' ${hostProtocppDir}/*.cc

# fix warnings:
# unused parameter options
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(StubOptions& options.*\)/\1\n(void)options;/g"
# unused parameter deterministic
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(bool deterministic.*\)/\1\n(void)deterministic;/g"
# unused definition
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION.*\)/\1\n#ifdef INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION\n#endif/g"
# unused parameter output in SerializeWithCachedSizes()
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(::google::protobuf::io::CodedOutputStream\* output.*) const {\)/\1\n(void)output;/g"
# uint->int in for loop of CapabilitiesResponse::SerializeWithCachedSizes in PluginRegistration.pb.cpp:1277:61
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(for (unsigned int i = 0, n =\)/for (int i = 0, n =/g"

# format using clang-format
find ${hostProtocppDir} \( -name "*.cpp" -or -name "*.h" \) -exec clang-format -i --style=file {} \;

# casting to int
find ${hostProtocppDir} -name '*.cpp' | xargs sed -i "s/\(.*length(),$\)/(int)\1/g"
