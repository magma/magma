# symphony project integration environment

### Commands

Build all containers
```
$ docker-compose build
```

Create and start all containers (`-d` runs them in the backgroud)
```
$ docker-compose up -d
```

Get logs of all containers (`-f` for follow)
```
$ docker-compose logs -f
```

If you want to see logs of a specific container, suffix the command above with the service name in `docker-compose.yaml`. For example:
```
$ docker-compose logs -f graph 
```

Connect to the database of the `auth` service
```
$ docker-compose exec mysql mysql -proot auth
```

Connect to the database of `fb-test` tenant
```
$ docker-compose exec mysql mysql -proot --database="tenant_fb-test"
```

See all running containers
```
$ docker-compose ps
```

### Writing integration tests

Integration tests are sitting in `integration/tests` and written in Go.
If you write a new test, that for example, needs to talk to `graph` service,
you can simply use the service name as the hostname. For example, `http://graph`.

To run all Go tests
```
$ docker-compose -f docker-compose.yaml -f docker-compose.override.yaml -f docker-compose.testing.yaml run --use-aliases test go test -v
```

To run a specific test, add the `-run` argument to `go test`. For example:
```
$ docker-compose -f docker-compose.yaml -f docker-compose.override.yaml -f docker-compose.testing.yaml run --use-aliases test go test -v -run=TestUser
```

In order to get a shell in the "test container" (mainly for ongoing test development), get a shell of the "test container":
```
$ docker-compose -f docker-compose.yaml -f docker-compose.override.yaml -f docker-compose.testing.yaml run --use-aliases test sh
```
You will get a shell in the "test container". Now, add your test case, run `go test`, repeat (the tests are mounted to the container).


### Debugging local environment

- Tracing - http://localhost:16686/
- Prometheus - http://localhost:9090/