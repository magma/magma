# Multi-tenancy for conductor

Proxy between clients and conductor that trasforms requests
and responses so that tenants are restricted to their own
domain.

So far it does not do much.

## Building
```$sh
./gradlew build
```

## Running the server
```shell script
java -jar server/build/libs/server-1.0.jar
```
This starts http server on `localhost:8081`.

## Testing
Make sure conductor is listening on `localhost:8080`.
Start the proxy and issue a HTTP request:
```sh
curl -v "localhost:8081/api/workflow/search?query="
```
