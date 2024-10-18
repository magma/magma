---
id: version-1.0.0-building_the_symphony_agent
title: Building the Symphony Agent
hide_title: true
original_id: building_the_symphony_agent
---
# Building the Symphony Agent

## Building the Base
If you've already built this base image and still have it cached and hanging around, skip down to Building the Agent. If you need to refresh the base image or this is your first time building the Symphony agent, read on.

First, start up your Docker daemon. Then, `cd` to where you've cloned Magma, and in a shell run:

```bash
cd magma/devmand/gateway/docker
./scripts/build_cached
```

This will build a base image that the Symphony Agent will subsequently build on top of. This image won't change much, so you probably won't need to rebuild it in the future.

## Building the Agent
Now, in the same directory, you can run:

```bash
./scripts/build
```

This will build the actual Symphony Agent Docker image that you will be able to run.
