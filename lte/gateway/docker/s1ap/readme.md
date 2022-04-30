# Build and publish s1aptester images

1. Clone the magma repository on an ubuntu 20.04 host and move into AGW docker directory in the repo and run build script.

```
cd lte/gateway/docker
s1ap/build-s1ap.sh
```

2. Publish images to your registry

```
s1ap/publish.sh yourregistry.com/yourrepo/
```

# Run s1aptester

1. Create an ubuntu 20.04 host that has 3 interfaces. eth0 which is your SGi interface, eth1 which is your S1 interface, and a third interface which is your ssh management interface to connect to the instance while doing the tests. You could skip the third interface if you have some kind of serial console access. Your host should have at least 20 GB of HD.

2. [Deploy a containerized AGW](https://github.com/magma/magma/tree/master/lte/gateway/docker), move into AGW docker directory `/var/opt/magma/docker` on the host and run start script. Make sure that your `.env` file points to your registry with your AGW and s1aptester images.
```
cd /var/opt/magma/docker
s1ap/start-s1ap.sh
```

2. This will drop you into a shell that you can start to run tests from, or run the full suite of tests.
```
root@472f8708ec12:/magma/lte/gateway/python/integ_tests#
# Run individual test(s)
make integ_test TESTS=s1aptests/test_attach_detach.py
# Run full suite
make integ_test
```

# Stop s1aptester

Move into AGW docker directory on the host and run stop script.
```
cd /var/opt/magma/docker
s1ap/stop-s1ap.sh
```

If inside of container, CTRL+d or exit from container and run stop script
```
root@472f8708ec12:/magma/lte/gateway/python/integ_tests# exit
s1ap/stop-s1ap.sh
```

# Simple example

```
sudo -i
wget https://raw.githubusercontent.com/magma/magma/master/lte/gateway/deploy/agw_install_docker.sh
mkdir -p /var/opt/magma/certs
echo "-----BEGIN CERTIFICATE-----
MIIDXzCCAkegAwIBAgIUakfCUNf7JMKbLDqHnuiG1QNhCQ8wDQYJKoZIhvcNAQEL
BQAwPzELMAkGA1UEBhMCVVMxMDAuBgNVBAMMJ3Jvb3RjYS5tYWdtYS0xNi1zeWRu
ZXkuZmFpbGVkd2l6YXJkLmRldjAeFw0yMTA3MTQyMjI5MjNaFw0zMTA3MTIyMjI5
MjNaMD8xCzAJBgNVBAYTAlVTMTAwLgYDVQQDDCdyb290Y2EubWFnbWEtMTYtc3lk
bmV5LmZhaWxlZHdpemFyZC5kZXYwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDFVmNFaAOkVD3c8W28FkGUVmBKDyj/T7N8C7PE43WvbBZJmO5TO1c887Dt
yiX8Ua2mpCQ2SF8DZtXojkLGKOFM85uxzTV1YI656u5BDSejRkm1UDeMT5R+tQJK
fyHYTt5ZNprX/dUrxYnp+h2zEl0PlzO5ijrktuZgM4KZjtQVaC1VirSC//2ZKQEo
2aX3L81ALrjVzsmH4TePKEY8StjDHC2Mg6LOaYR/+Gu272P39/heULrm147g1k0k
haeKv8qrI0dfvBcZveTzYf77iA6/OeVzYtWwM3zr1Z1cFALZrcuS6R6DrsAInseH
qeiMh4kLfyoh0vQCNpEAJQgt5PmVAgMBAAGjUzBRMB0GA1UdDgQWBBTpk41oSZDv
hSlsCLboVWzT5w414TAfBgNVHSMEGDAWgBTpk41oSZDvhSlsCLboVWzT5w414TAP
BgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBWeu3+kB2WmxXbIDPU
1JGw4rvj/+u4mvN1maYkibqCZyKMXLuqwOy9wMhtniHgKp/RIxFI+W3FTq4Tik++
kwDemaYq3nwbHMwBXwFh/T9I9ExtWBCogj+LLFUsrPDJNmUwYnnEMRh6beF8oT1E
Da3oNVZ70Tyv0DnWozW+4TQZ8bTOQ/bpjoFNZPVB3Jr7tjVLfPez8m/clM8+War+
gjTiiUsJkJP2uhKmWkb58CCiH+k2EH3fw2IUmc0fgTMGZ5vv8g1OjCBrXspnGSpk
iJf9ryw/jIH/9RGxSUO7tiQxe/IShf65clsyxlAjrSr7JvbYwOyIXAbgNA7vk0lc
nmmv
-----END CERTIFICATE-----" > /var/opt/magma/certs/rootCA.pem
bash agw_install_docker.sh
sed -i 's/DOCKER_REGISTRY=/DOCKER_REGISTRY=public.ecr.aws\/yourrepo\//' /var/opt/magma/docker/.env
cd /var/opt/magma/docker
s1ap/start-s1ap.sh
make integ_test TESTS=s1aptests/test_attach_detach.py
```
