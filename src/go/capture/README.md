#Capture Service

The capture service allows remote introspection into RPC calls of any GRPC service.
Targeting is config driven

##Prerequisites
- Running the capture service requires that the capture service dial/server options are attached to the GRPC client or GRPC server.
- /src/go/middleware.go has two helper functions```GetDialOptions```and```GetServerOptions```that handle this for you. The capture service expects identical config and capture.Buffer pointers to be passed into the middleware.

##Configuration Instructions
- The capture service holds a pointer to the config option, so changing the targeted MatchSpecs can be done without restarting the service.
- Multiple MatchSpecs can be set and it also supports wildcarding.
### Examples:
- MatchSpec
```
&configpb.CaptureConfig_MatchSpec{
	Service: "magma.sctpd.SctpdDownlink",
	Method:  "SendDl",
}
```
- Wildcard Service
```
&configpb.CaptureConfig_MatchSpec{
	Service: "*",
	Method:  "SendDl",
}
```
- Wildcard Method
```
&configpb.CaptureConfig_MatchSpec{
	Service: "magma.sctpd.SctpdDownlink",
	Method:  "*",
}
```
4. Record all RPCs.
```
&configpb.CaptureConfig_MatchSpec{
	Service: "*",
	Method:  "*",
}
```

#Command Line Use
Using prototool is an easy way to make grpc calls.
https://github.com/uber/prototool

##Flush Capture Service Buffer

```prototool grpc --address 192.168.60.142:6001 --method magma.capture.Capture/Flush --data '{}'```

##Get Config from Config Service

```prototool grpc --address 192.168.60.142:6000 --method magma.config.Config/GetConfig --data '{}'```

##Replace config with wild card matchspec

```prototool grpc --address 192.168.60.142:6000 --method magma.config.Config/ReplaceConfig --data '{"config":{"logLevel":"INFO","sctpdDownstreamServiceTarget":"unix:///tmp/sctpd_downstream.sock","sctpdUpstreamServiceTarget":"unix:///tmp/sctpd_upstream.sock","mmeSctpdDownstreamServiceTarget":"unix:///tmp/mme_sctpd_downstream.sock","mmeSctpdUpstreamServiceTarget":"unix:///tmp/mme_sctpd_upstream.sock","configServicePort":"6000","vagrantPrivateNetworkIp":"192.168.60.142","captureServicePort":"6001","captureConfig":{"matchSpecs":[{"service":"*","method":"*"}]}}}'```



#Golden File Generation
1. Follow prerequisites listed above and ensure that magma vm and trfgen vm are both on and AGWD is running on magma VM and proxying sctpd uplink and downlink traffic.
2. Build and run the capture service from the magma_test vm
```
bazel build //src/go/capture/gen:gen
bazel run //src/go/capture/gen:gen
```
3. Capture will parse tests from the integ_test makefile and iterate over them, running each test, flushing the capture buffer and wrting them to a file in /resources/s1aptests/**testname**.py.golden
4. **TODO** update the test parser to be configurable.
