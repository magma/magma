// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

var https = require('https');
var fs = require('fs');
var jp = require('jsonpath');

var network = process.env['symphony_network'];

function doHttp(path, printError, handler, method, body) {
    var options = {
        hostname: 'api.magma.etagecom.io',
        port: 443,
        path: path,
        method: method,
        json: 'true',
        key: fs.readFileSync('/run/secrets/orc8r_api_key'),
        cert: fs.readFileSync('/run/secrets/orc8r_api_cert'),
        headers: {
            accept: 'application/json',
            'Content-Type': 'application/json',
            'Content-Length': body.length
        },
        body: body
    };

    var req = https.request(options, function(res) {
        var body = '';

        res.on('data', function (chunk) {
            body = body + chunk;
        });

        res.on('end',function(){
            try {
                handler(body);
            } catch (err) {
                if (printError) {
                    process.stdout.write("error [" + err + "]\n");
                }
            }
        });
    });

    if (body.length != 0) {
        req.write(body);
    }

    req.end();
}

function doGet(path, printError, handler) {
    doHttp(path, printError, handler, 'GET', "");
}

function doPost(path, printError, handler, body) {
    doHttp(path, printError, handler, 'POST', body);
}

function doPut(path, printError, handler, body) {
    doHttp(path, printError, handler, 'PUT', body);
}

function p(header, obj, path) {
    process.stdout.write("\t\t " + header + " = ");
    var result = jp.query(obj, path);
    if (result.length > 1) {
        process.stdout.write(JSON.stringify(result) + "\n");
    } else {
        process.stdout.write(JSON.stringify(result[0]) + "\n");
    }
}

function printInterfaceModels(deviceState) {
    var interfaces = deviceState["openconfig-interfaces:interfaces"];
    if (interfaces) {
        p("s/oper-status", interfaces, "$..['interface']..[*].state..['oper-status']");
        p("s/admin-status", interfaces, "$..['interface']..[*].state..['admin-status']");
        p("c/enabled", interfaces, "$..['interface']..[*].config..['enabled']");
    }
}

function printWifiModels(deviceState) {
    var joinedAps = deviceState["openconfig-ap-manager:joined-aps"];
    if (joinedAps) {
        p("mac", joinedAps, "$..['joined-ap']..[0].state.mac");
        p("opstate", joinedAps, "$..['joined-ap']..[0].state.opstate");
        p("model", joinedAps, "$..['joined-ap']..[0].state.model");
        p("serial", joinedAps, "$..['joined-ap']..[0].state.serial");
        p("hostname", joinedAps, "$..['joined-ap']..[0].state.hostname");
        p("ipv4", joinedAps, "$..['joined-ap']..[0].state.ipv4");
        p("software-version",
            joinedAps, "$..['joined-ap']..[0].state..['software-version']");

        p("uptime", joinedAps, "$..['joined-ap']..[0].state.uptime");
    }

    var aps = deviceState["openconfig-access-points:access-points"];
    if (aps) {
        p("operating-frequency", aps,
            "$..['access-point']..[0].ssids.ssid[*]." +
            "state..['operating-frequency']");
        p("ssids", aps, "$..['access-point']..[0].ssids.ssid[*].name");
        p("bssids", aps,
            "$..['access-point']..[0].ssids.ssid[*].bssids.bssid[*].bssid");
    }

    var system = deviceState["fbc-symphony-device:system"];
    if (system) {
        p("longitude", system, "$..['geo-location']..longitude");
        p("latitude", system, "$..['geo-location']..latitude");
        p("latencies", system, "$.latencies.latency[*].rtt");
        p("hotspot", system, "$.venue");
    }
}

function processAgent(agentId, agent) {
    if (agent && agent.status && agent.status.meta && agent.status.meta.devmand) {
        process.stdout.write("Agent " + agentId + "\n");
        var managedDevices = JSON.parse(agent.status.meta.devmand);
        for (var managedDevice in managedDevices) {
            process.stdout.write("\t" + managedDevice + "\n");
            printInterfaceModels(managedDevices[managedDevice]);
            printWifiModels(managedDevices[managedDevice]);

            doGet('/magma/v1/symphony/' +  network + '/devices/' + managedDevice,
                true,
                function(body) {
                    var device = JSON.parse(body);
                    console.log(managedDevices[managedDevice]
                        ["openconfig-interfaces:interfaces"].interface[1]
                        .config.enabled);

                    device.config.device_config = JSON.stringify(managedDevices[managedDevice]);

                    //doPut('/magma/v1/symphony/' +  network + '/devices/' + managedDevice, true,
                    //    function(body) {
                    //        process.stdout.write("P".repeat(80) + "\n");
                    //    }, JSON.stringify(device));
                });
        }
    }
}

function loop() {
    doGet('/magma/v1/symphony/' +  network + '/agents', true,
        function(body) {
            var agents = JSON.parse(body);
            process.stdout.write("#".repeat(80) + "\n");
            for (var agentId in agents) {
                processAgent(agentId, agents[agentId]);
            }
        });
}

let timerId = setInterval(loop, 10000);
setTimeout(() => { clearInterval(timerId);}, 108000);
