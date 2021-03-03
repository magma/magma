# magmad service

[![magma](https://circleci.com/gh/magma/magma.svg?style=shield)](https://circleci.com/gh/magma/magma)

**magmad** service is a core part of [magma](https://magma.github.io/magma) gateway platform, it is required on every gateway and facilitates the following gateway functionality:

* Maintaining and displaying unique gateway identity information necessary for the gateway register on and join a network
* Bootstraping gateway on the network and maintaining gateway's client certificates
* Checking in with the cloud and providing the cloud with up to date gateway status
* Monitoring gateway services, collecting and publishing the services' and system metrics to the cloud
* Creating and maintaining secure cloud to gateway RPC channel for real time cloud requests fulfilment
* Delivering and maintaining of cloud controlled gateway services' configuration
* Providing gateway control APIs to the cloud, such as:
  * starting/stopping gateway services
  * rebooting gateway
  * running network diagnostics
  * remote execution of permitted commands
  * etc.
* Managing Gateway software & facilitating Gateway software upgrades **

## Implementation

This is Go native implementation of magmad service, it is intended as an lightweight alternative to Magma's main 
python magmad service implementation (magma/orc8r/gateway/python/magmad), primary targeting embedded systems or systems 
with limited resources and/or lack of full python support (systems with less then 256MB of RAM or storage).

This implementation is not guranteed to support all features supported by current python magmad (such as Gateway software
upgrades), it provides required core functionality for a Gateway to be a part of Magma network.

It'll work with all existing *magmad.yml*, *service_registry.yml* and *control_proxy.yml* configurations located in
*/etc/magma/configs* or legacy */etc/magma* directories.

Current implementation requires at least 16MB of RAM and 3MB of storage (compressed) on most supported systems/architectures.


## Building magmad

magmad binary can be built on any system with [installed Go tools](https://golang.org/doc/install#install). Go1.12 or later is recommended.

To build:
* Clone magma [github repository](https://github.com/magma/magma)
* (optional) Create a target directory for magmad binary (example: *mkdir -p ~/bin/arm; mkdir -p ~/bin/amd64*)
* *cd magma/orc8r/gateway/go*
* Use Go tool to build magmad for your target platform. Examples:
  * *GOOS=linux GOARCH=arm go build -o ~/bin/arm/ magma/gateway/services/magmad*
  * *GOOS=linux GOARCH=amd64 go build -o ~/bin/amd64/ magma/gateway/services/magmad*
* To reduce magmad binary size, you may use the following options:
  * *GOOS=linux GOARCH=arm go build -o ~/bin/arm/magmad.prod -ldflags="-s -w" magma/gateway/services/magmad*
  * You can also use [upx](https://upx.github.io/) to compress the binary:
    * *upx --brute -o ~/bin/arm/magmad ~/bin/arm/magmad.prod*

## Supported Architectures & Systems

magmad is a Go native, statically-linked application. It should be compatible with all supported by Go distributions.
As of go1.14 the supported distributions are:
* android/386
* android/amd64
* android/arm
* android/arm64
* darwin/386
* darwin/amd64
* darwin/arm
* darwin/arm64
* dragonfly/amd64
* freebsd/386
* freebsd/amd64
* freebsd/arm
* freebsd/arm64
* illumos/amd64
* js/wasm
* linux/386
* linux/amd64
* linux/arm
* linux/arm64
* linux/mips
* linux/mips64
* linux/mips64le
* linux/mipsle
* linux/ppc64
* linux/ppc64le
* linux/riscv64
* linux/s390x
* netbsd/386
* netbsd/amd64
* netbsd/arm
* netbsd/arm64
* openbsd/386
* openbsd/amd64
* openbsd/arm
* openbsd/arm64
* plan9/386
* plan9/amd64
* plan9/arm
* solaris/amd64
* windows/386
* windows/amd64
* windows/arm

## Join the Magma Community

- Mailing lists:
  - Join [magma-dev](https://groups.google.com/forum/#!forum/magma-dev) for technical discussions
  - Join [magma-announce](https://groups.google.com/forum/#!forum/magma-announce) for announcements
- Discord:
  - [magma\_dev](https://discord.gg/WDBpebF) channel

See the [CONTRIBUTING](../../../../CONTRIBUTING.md) file for how to help out.

## License

Magma is BSD License licensed, as found in the [LICENSE](../../../../LICENSE) file.
The EPC is OAI is offered under the OAI Apache 2.0 license, as found in the LICENSE file in the OAI directory.

** - Go magmad does not currently provide Gateway software upgrade capabilities