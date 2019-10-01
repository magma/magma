var https = require('https');
var fs = require('fs');

function loop() {
    var symphony_agent = process.env['symphony_agent'];
    var options = {
        hostname: 'api.magma.etagecom.io',
        port: 443,
        path: '/magma/networks/southpoll_dev/gateways/' + symphony_agent + '/status',
        method: 'GET',
        json: 'true',
        key: fs.readFileSync('/run/secrets/orc8r_api_key'),
        cert: fs.readFileSync('/run/secrets/orc8r_api_cert'),
        headers: { accept: 'application/json' }
    };

    var req = https.request(options, function(res) {
        var body = '';

        res.on('data', function (chunk) {
            body = body + chunk;
        });

        res.on('end',function(){
            try {
                var agent_data = JSON.parse(body);
                process.stdout.write(body);
                var managed_devices = JSON.parse(agent_data.meta.devmand);
                process.stdout.write("\n");
                process.stdout.write("\n");
                for (var managed_device in managed_devices) {
                    try {
                        var ifindex = process.env['symphony_ifindex'];
                        process.stdout.write(
                            managed_devices[managed_device]
                            ["openconfig-interfaces:interfaces"]
                            ["interface"][ifindex]["state"]
                            ["counters"]["out-unicast-pkts"]);
                        process.stdout.write("\n");
                    } catch (err) {
                        process.stdout.write(managed_device);
                        process.stdout.write("\n");
                    }
                }
            } catch (err) {
            }
        });
    });

    req.end();
}

let timerId = setInterval(loop, 2000);

setTimeout(() => { clearInterval(timerId);}, 2000000);
