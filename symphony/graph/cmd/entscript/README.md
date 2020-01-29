# Entscript

This binary should be used to run custom ent "scripts" on the database.
This is safer than performing manual sql queries on database.

### How to test

First add your lines in the relevant function in the tool

Build and re-create graph
```
$ docker-compose build graph
$ docker-compose up -d

```

Connect to graph docker instance
```
$ docker-compose exec graph /bin/sh
```

From docker instance
```
$ /bin/enscript --tenant=fb-test --user=fbuser@fb.com
```

### How to run in production

Find the github revision that is currently in production. It can be found in the output of this command
```
$ kubectl describe pod {graph_pod_name}
```

Checkout the symphony github repository to the correct revision
```
$ git reset --hard {github_revison}
```

Now add your lines in the relevant function in the tool

Compile and upload the tool to the relevant kubernetes container
```
$ GOOS=linux go build -o ./build ./cmd/enscript
$ kubectl cp build/enscript {relavent_graph_pod_name}:/bin
```

Connect to graph kubernetes instance
```
$ kubectl exec -it svc/inventory-graph --container graph sh
```

From kubernetes instance
```
$ /bin/enscript --tenant=fb-test --user=fbuser@fb.com
```