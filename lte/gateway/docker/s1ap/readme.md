# Build and publish s1aptester images

1. Move into AGW docker directory in the repo and run build script

```
cd lte/gateway/docker
s1ap/build-s1ap.sh
```

2. Publish images to your registry

```
s1ap/publish.sh http://yourregistry.com/yourrepo/
```

# Run s1aptester

1. Move into AGW docker directory on the host and run start script. Make sure that your `.env` file points to your registry.
```
cd /var/opt/magma/docker
s1ap/start-s1ap.sh
```

2. This will drop you into a shell that you can start to run tests from, or run the full suite of tests.
```
root@472f8708ec12:/magma/lte/gateway/python/integ_tests#
# Run individual test(s)
make integ_test TESTS=test_attach_detach.py
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
