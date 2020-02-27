---
id: readme_cwag_integ_test
title: CWAG Integration Test 
hide_title: true
---
## Integration Test Setup 
This integration currently uses 3 separate VMs to run the tests: `cwag-test` 
and `cwag-dev` in `magma/cwf/gateway`, and `magma-trfserver` in `magma/lte/gateway`.
The fabfile in this directory can be used to automate the setup needed to run
the test.

###  `cwag-dev` 
This VM will build and run all cwag services and mock core services needed to 
run the test. See the various `docker-compose` files in `cwf/gateway/docker` 
to see the complete list of services. 

### `cwag-test`
This VM will be used to run the tests. We will also run a UE simulator service 
to simulate a UE device in the test. 

### `magma-trfserver`
This VM runs an iperf3 server.

## Current Tests

- Authentication
- Session Creation
- Basic Policy enforcement

## Running the test 
#### Requirements

- fabric3 
- see https://facebookincubator.github.io/magma/docs/basics/prerequisites for 
our prerequisites on running our VMs. (Ignore docker part)

To the run the test, run `fab integ_test` from `magma/cwf/gateway`.
This fabfile will
- Provision the 3 VMs
- Build and start docker containers on `cwag-dev`
- Start the UE simulator service on `cwag-test`
- Start the iperf3 server on `magma-trf`
- Run the integration test on `cwag-test`
- Clean up

## Development 

- To see the list of running services, run `docker ps` in the `cwag-dev` VM.
- To see per-service logs, run `docker-compose logs <container_name>`
- To go into a running container, run `docker-compose exec <container_name> bash`
- `/usr/local/bin/pipelined_cli.py` in pipelined service maybe useful for 
viewing installed flows for debugging.

## FAQ

#### Unit tests are not able to run: `/cwf/gateway: No such file or directory`

&rightarrow; The cwag-dev VM probably did not provision properly. Run 
`vagrant provision cwag` in `magma/cwf/gateway` to provision the VM again 
to see specific errors. 

#### Docker is failing to build due to lack of space

&rightarrow; Since docker does not garbage collect previously built images, we 
will have to manually prune them. Run `docker system df` to see memory usage 
and what can be deleted. To remove these images, run `docker system prune to docker image prune --filter until=24h`.