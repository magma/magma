---
id: version-1.7.0-dev_testing
title: Test NMS
hide_title: true
original_id: dev_testing
---

# Test NMS

This guide covers tips for quickly validating NMS changes.

## Run tests on the host

There are some simple tests that you should run after making changes. Before running the tests make sure to `cd ${MAGMA_ROOT}/nms` and run the tests from this directory.

### Pre-commit tests

These are the set of tests that you should run before opening a pull request.

#### Eslint

Run `yarn run eslint` to test that there are no linter errors. You can also run `yarn run eslint --fix` to fix automatically fixable linter errors.

#### Flow

Run `yarn run flow` to test that there are no type errors.

#### Unit Test

Run `yarn run test` to test that your changes have not broken existing individual units / components. This command also runs any new unit tests that you created.

### End to end test

This is not a mandatory test. You can run `yarn run test:e2e` to test that your
changes have passed all the end to end tests.
