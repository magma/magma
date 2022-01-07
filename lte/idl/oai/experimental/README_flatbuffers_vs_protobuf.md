---
id: flat_buffers
title: Benchmarking FlatBuffers versus protobuf
hide_title: true
---
# Benchmarking Flatbuffers versus protobuf

The goal is to quantify and obtain measurements of the memory and time spent on serialization using protobuf vs [FlatBuffers](http://google.github.io/flatbuffers/index.html#flatbuffers_overview) library, using a minimal test object.

The goal of benchmarking FlatBuffers is to take advantage of zero-copy, meaning no serialization cost (CPU) from the handled contexts in tasks MME_APP, S1AP, etc. The cost of serialization would be mainly the cost of the Redis client call to send the context to Redis server. and could outperform the actual protobuf implementation.

## Scope of benchmarking

The purpose of the benchmark is to compare the serialization of an AGW LTE UE MM context into redis DB with protobuf versus [FlatBuffers](http://google.github.io/flatbuffers/index.html#flatbuffers_overview)

This AGW LTE UE MM context test object has been retrieved running an attach detach S1AP test and obtaining it by running state_cli.py: https://github.com/magma/magma/blob/master/docs/readmes/lte/dev_notes.md#checking-redis-entries-for-stateless-services)

### Use of FlatBuffers

The goal of using  FlatBuffers is to take advantage of zero-copy, meaning no serialization cost (CPU) from the handled contexts in tasks MME_APP, S1AP, etc.

So this first trial is done with mutables FlatBuffers. In order to have no serialization phase, we work with pre-order construction of data (pre-allocated arrays of structs and pre-allocated arrays of bytes for strings), meaning we do not use [object based API](http://google.github.io/flatbuffers/flatbuffers_guide_use_cpp.html#flatbuffers_cpp_object_based_api).

#### Limitations

One limitation is that we cannot modelize bitfields, the smallest information in an IDL file that can be modeled is a byte.

We encountered an issue after trying to handle most of the attributes of the UE MM context in the [IDL schemas](http://google.github.io/flatbuffers/flatbuffers_guide_writing_schema.html) files (lte/idl/oai/experimental/mme_ue_state.fbs) . This issue is that we already got a root UE_MM_CONTEXT object size slightly greater than allowed (65536 bytes). We had to reduce the number of allowed PDN contexts (5 instead of 10 in current implementation), and reduce the biggest size of ESM message we could receive.

This has to be considered versus using the FlatBuffers Object API.

## Benchmarking

For time benchmarking both PB and FB, std::chrono is used and measurements are flushed to stdout.

For memory benchmarking, [rusage](https://linux.die.net/man/2/getrusage) is used.

Benchmarking for AGW is done on the magma dev VM. You just need to bring only this VM up.

### Prerequisites

First, [install](http://google.github.io/flatbuffers/flatbuffers_guide_building.html#autotoc_md7) flatbuffer compiler.

### Run benchmarks

Inside magma VM, build mme:

```bash
[VM] cd $MAGMA_ROOT/lte/gateway
[VM] make build_oai
```

Run benchmark script for Protobuf and Flatbuffers

```bash
# Prepare mapped temp forder for results
[VM] TEST_RESULT_DIR=$MAGMA_ROOT/tmp/RUN1 &&
     mkdir -p -m 0777 $TEST_RESULT_DIR &&
     cd $TEST_RESULT_DIR &&
     # Do the benchmarking for Protobuf and Flatbuffers
     $MAGMA_ROOT/lte/gateway/c/scripts/mme_ue_context_serialization_benchmarking
    # You should see several log files in RUN1 folder
```

### Results

To display the timing measurements, on your host, go in the same temp folder:

```bash
[HOST] # MAGMA_ROOT should point to your magma folder on your host
[HOST] TEST_RESULT_DIR=$MAGMA_ROOT/tmp/RUN1 &&
       cd $TEST_RESULT_DIR &&
       $MAGMA_ROOT/lte/gateway/c/scripts/mme_ue_context_serialization_plotting
```
