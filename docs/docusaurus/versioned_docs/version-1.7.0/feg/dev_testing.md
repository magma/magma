---
id: version-1.7.0-dev_testing
title: Test Federation Gateway
hide_title: true
original_id: dev_testing
---

# Test FeG

This guide covers tips for quickly validating Federation Gateway changes.

## Run tests on the host

The normal way to run FeG unit tests is `$MAGMA_ROOT/feg/gateway/docker/build.py --tests`. This
builds and brings up a test container, then runs the full set of unit
tests.

The script also includes additional commands, such as linters and formatters.
When committing changes in FeG, use `$MAGMA_ROOT/feg/gateway/docker/build.py --precommit` to run a set of recommended commands.
