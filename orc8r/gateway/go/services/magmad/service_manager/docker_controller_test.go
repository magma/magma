/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package service_manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInspectResult = []byte(`[
    {
        "Id": "407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba",
        "Created": "2020-02-26T00:03:31.39300804Z",
        "Path": "envdir",
        "Args": [
            "/var/opt/magma/envdir",
            "/var/opt/magma/bin/aaa_server",
            "-logtostderr=true",
            "-v=0"
        ],
        "State": {
            "Status": "%s",
            "Running": true,
            "Paused": false,
            "Restarting": false,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 13100,
            "ExitCode": 0,
            "Error": "",
            "StartedAt": "2020-03-04T19:04:47.252591285Z",
            "FinishedAt": "2020-03-04T19:04:41.307047441Z"
        },
        "Image": "sha256:e802dcd6d730aa59c3b5b8d1cb362e0b11978e77c6c0bb708d51d6d7a629d45f",
        "ResolvConfPath": "/var/lib/docker/containers/407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba/hostname",
        "HostsPath": "/var/lib/docker/containers/407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba/hosts",
        "LogPath": "/var/lib/docker/containers/407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba/407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba-json.log",
        "Name": "/aaa_server",
        "RestartCount": 0,
        "Driver": "overlay2",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "docker-default",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": [
                "/var/opt/magma/configs:/var/opt/magma/configs:rw",
                "/etc/magma/magmad.yml:/etc/magma/magmad.yml:rw",
                "/etc/magma/service_registry.yml:/etc/magma/service_registry.yml:rw",
                "/var/opt/magma/certs:/var/opt/magma/certs:rw",
                "/etc/magma/control_proxy.yml:/etc/magma/control_proxy.yml:rw",
                "/var/opt/magma/certs/rootCA.pem:/var/opt/magma/certs/rootCA.pem:rw",
                "/etc/magma/magmad_legacy.yml:/etc/magma/magmad_legacy.yml:rw"
            ],
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {
                    "max-file": "10",
                    "max-size": "10mb"
                }
            },
            "NetworkMode": "docker_magma",
            "PortBindings": {},
            "RestartPolicy": {
                "Name": "always",
                "MaximumRetryCount": 0
            },
            "AutoRemove": false,
            "VolumeDriver": "",
            "VolumesFrom": [],
            "CapAdd": null,
            "CapDrop": null,
            "Dns": null,
            "DnsOptions": null,
            "DnsSearch": null,
            "ExtraHosts": null,
            "GroupAdd": null,
            "IpcMode": "shareable",
            "Cgroup": "",
            "Links": null,
            "OomScoreAdj": 0,
            "PidMode": "",
            "Privileged": false,
            "PublishAllPorts": false,
            "ReadonlyRootfs": false,
            "SecurityOpt": null,
            "UTSMode": "",
            "UsernsMode": "",
            "ShmSize": 67108864,
            "Runtime": "runc",
            "ConsoleSize": [
                0,
                0
            ],
            "Isolation": "",
            "CpuShares": 0,
            "Memory": 0,
            "NanoCpus": 0,
            "CgroupParent": "",
            "BlkioWeight": 0,
            "BlkioWeightDevice": null,
            "BlkioDeviceReadBps": null,
            "BlkioDeviceWriteBps": null,
            "BlkioDeviceReadIOps": null,
            "BlkioDeviceWriteIOps": null,
            "CpuPeriod": 0,
            "CpuQuota": 0,
            "CpuRealtimePeriod": 0,
            "CpuRealtimeRuntime": 0,
            "CpusetCpus": "",
            "CpusetMems": "",
            "Devices": null,
            "DeviceCgroupRules": null,
            "DiskQuota": 0,
            "KernelMemory": 0,
            "MemoryReservation": 0,
            "MemorySwap": 0,
            "MemorySwappiness": null,
            "OomKillDisable": false,
            "PidsLimit": 0,
            "Ulimits": null,
            "CpuCount": 0,
            "CpuPercent": 0,
            "IOMaximumIOps": 0,
            "IOMaximumBandwidth": 0,
            "MaskedPaths": [
                "/proc/asound",
                "/proc/acpi",
                "/proc/kcore",
                "/proc/keys",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/proc/scsi",
                "/sys/firmware"
            ],
            "ReadonlyPaths": [
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger"
            ]
        },
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/f1b3a5fb9e019c4edd014e079bf4c45d7d12c9a46f416e6d0cffb0873dfa2e48-init/diff:/var/lib/docker/overlay2/b72b684d6128e998d1fd959b2d4679b827f9b3ce0dd6b5f67950cf768d5f226b/diff:/var/lib/docker/overlay2/b4ebf001a245395e584289c96842494162960e2a8973bd0bf81fd1d543b1d835/diff:/var/lib/docker/overlay2/e9fc8ba35e962ed553c3af31acd6430158984e57cbc012879be75cb109905823/diff:/var/lib/docker/overlay2/3b0a3308f074f70f7d76a85410decb03eb95c0ddd02cb2d1ea28e092368962ae/diff:/var/lib/docker/overlay2/e17dd793d9e8f2996c4cee00dbf718ea605dbc9f0d86cb97488b9646d9e58590/diff:/var/lib/docker/overlay2/e06f12d0f2a92e9a734026715f5458e86d319e3cd39f753e8f06503a0e6eeb46/diff:/var/lib/docker/overlay2/58ad390e4295bf0f82077dc8f07cc924369aee7595668555a69d8030d51956f0/diff:/var/lib/docker/overlay2/1ee7aafe15298d165bf12f722e16b3d1dab78f10df2510e1d619c5a99d919452/diff:/var/lib/docker/overlay2/e854926f26772143ef6809ac3669cda0df4b23e7bdaddaef13c54e109bfc8435/diff",
                "MergedDir": "/var/lib/docker/overlay2/f1b3a5fb9e019c4edd014e079bf4c45d7d12c9a46f416e6d0cffb0873dfa2e48/merged",
                "UpperDir": "/var/lib/docker/overlay2/f1b3a5fb9e019c4edd014e079bf4c45d7d12c9a46f416e6d0cffb0873dfa2e48/diff",
                "WorkDir": "/var/lib/docker/overlay2/f1b3a5fb9e019c4edd014e079bf4c45d7d12c9a46f416e6d0cffb0873dfa2e48/work"
            },
            "Name": "overlay2"
        },
        "Mounts": [
            {
                "Type": "bind",
                "Source": "/etc/magma/magmad_legacy.yml",
                "Destination": "/etc/magma/magmad_legacy.yml",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/var/opt/magma/configs",
                "Destination": "/var/opt/magma/configs",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/etc/magma/magmad.yml",
                "Destination": "/etc/magma/magmad.yml",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/etc/magma/service_registry.yml",
                "Destination": "/etc/magma/service_registry.yml",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/var/opt/magma/certs",
                "Destination": "/var/opt/magma/certs",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/etc/magma/control_proxy.yml",
                "Destination": "/etc/magma/control_proxy.yml",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "bind",
                "Source": "/var/opt/magma/certs/rootCA.pem",
                "Destination": "/var/opt/magma/certs/rootCA.pem",
                "Mode": "rw",
                "RW": true,
                "Propagation": "rprivate"
            }
        ],
        "Config": {
            "Hostname": "407f286d41e6",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": false,
            "AttachStderr": false,
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            ],
            "Cmd": [
                "envdir",
                "/var/opt/magma/envdir",
                "/var/opt/magma/bin/aaa_server",
                "-logtostderr=true",
                "-v=0"
            ],
            "Image": "facebookconnectivity-feg-docker.jfrog.io/gateway_go:276ce097",
            "Volumes": {
                "/etc/magma/control_proxy.yml": {},
                "/etc/magma/magmad.yml": {},
                "/etc/magma/magmad_legacy.yml": {},
                "/etc/magma/service_registry.yml": {},
                "/var/opt/magma/certs": {},
                "/var/opt/magma/certs/rootCA.pem": {},
                "/var/opt/magma/configs": {}
            },
            "WorkingDir": "",
            "Entrypoint": null,
            "OnBuild": null,
            "Labels": {
                "com.docker.compose.config-hash": "a2d50980cc3f1e558f5dd3f5791d6f159880dc14f20aca5d2e2ba72ace92fdc5",
                "com.docker.compose.container-number": "1",
                "com.docker.compose.oneoff": "False",
                "com.docker.compose.project": "docker",
                "com.docker.compose.service": "aaa_server",
                "com.docker.compose.version": "1.25.0-rc1"
            }
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "16dd6dc5e3e95ad5ad9d622f9a9e3595030c791aeb8dec4b8464ac7c24a27f6a",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {},
            "SandboxKey": "/var/run/docker/netns/16dd6dc5e3e9",
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "",
            "Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "",
            "IPPrefixLen": 0,
            "IPv6Gateway": "",
            "MacAddress": "",
            "Networks": {
                "docker_magma": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": [
                        "407f286d41e6",
                        "aaa_server"
                    ],
                    "NetworkID": "50e76c5ba19336a72f97080ec4030c4fe1d54ec010bde89e5b9ef4e51c45e7bf",
                    "EndpointID": "216f0dec6c82b92e19d8cd7fcbb57a0908445f0b60dfe352dacc1e3ecdbd6afb",
                    "Gateway": "172.20.0.1",
                    "IPAddress": "172.20.0.3",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:14:00:03",
                    "DriverOpts": null
                }
            }
        }
    }
]
`)

	emptyInspectResult  = []byte{}
	emptyInspectArr     = []byte("[]")
	missingInspectState = []byte(`[
    {
        "Id": "407f286d41e6c03e65bf5e049330efa061fe4780f50991354895bcedd6a93dba",
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "16dd6dc5e3e95ad5ad9d622f9a9e3595030c791aeb8dec4b8464ac7c24a27f6a",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {},
            "SandboxKey": "/var/run/docker/netns/16dd6dc5e3e9",
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "",
            "Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "",
            "IPPrefixLen": 0,
            "IPv6Gateway": "",
            "MacAddress": "",
            "Networks": {
                "docker_magma": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": [
                        "407f286d41e6",
                        "aaa_server"
                    ],
                    "NetworkID": "50e76c5ba19336a72f97080ec4030c4fe1d54ec010bde89e5b9ef4e51c45e7bf",
                    "EndpointID": "216f0dec6c82b92e19d8cd7fcbb57a0908445f0b60dfe352dacc1e3ecdbd6afb",
                    "Gateway": "172.20.0.1",
                    "IPAddress": "172.20.0.3",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:14:00:03",
                    "DriverOpts": null
                }
            }
        }
    }
	]`)
)

func TestDockerControllerParser(t *testing.T) {
	out := []byte(fmt.Sprintf(string(testInspectResult), "running"))
	s, e := parseDockerInspectResult(out)
	assert.NoError(t, e)
	assert.Equal(t, Active, s)

	out = []byte(fmt.Sprintf(string(testInspectResult), "dead"))
	s, e = parseDockerInspectResult(out)
	assert.NoError(t, e)
	assert.Equal(t, Failed, s)

	out = []byte(fmt.Sprintf(string(testInspectResult), "foobar"))
	s, e = parseDockerInspectResult(out)
	assert.NoError(t, e)
	assert.Equal(t, Unknown, s)

	s, e = parseDockerInspectResult(emptyInspectResult)
	assert.Error(t, e)
	assert.Equal(t, Error, s)
	s, e = parseDockerInspectResult(emptyInspectArr)
	assert.Error(t, e)
	assert.Equal(t, Error, s)
	s, e = parseDockerInspectResult(missingInspectState)
	assert.NoError(t, e)
	assert.Equal(t, Unknown, s)
}
