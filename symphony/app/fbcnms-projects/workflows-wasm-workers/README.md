### Set conductor proxy as CONDUCTOR_API variable
```shell script
CONDUCTOR_API="http://localhost:8088/proxy/api"
```

### Create new workflow wasm-example
POST to `/metadata/workflow`

```shell script
curl -v \
-H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
-H 'Content-Type: application/json' \
${CONDUCTOR_API}/metadata/workflow -d '
{
    "name": "js-example",
    "description": "javascript lambdas in running in wasm",
    "ownerEmail": "example@example.com",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
            "taskReferenceName": "create_json_ref",
            "name": "GLOBAL___js",
            "inputParameters": {
                "args": "${workflow.input.enter_your_name}",
                "outputIsJson": "true",
                "script": "console.log(JSON.stringify({name: argv[1]}));"
            },
            "type": "SIMPLE",
            "startDelay": 0,
            "optional": false,
            "asyncComplete": false
        },
        {
            "taskReferenceName": "calculate_name_length_ref",
            "name": "GLOBAL___js",
            "inputParameters": {
                "args": "${create_json_ref.output.result}",
                "outputIsJson": "true",
                "script": "let json=JSON.parse(argv[1]); json.name_length = json.name.length; console.log(JSON.stringify(json));"
            },
            "type": "SIMPLE",
            "startDelay": 0,
            "optional": false,
            "asyncComplete": false
        }
    ]
}
'
```

### Execute the workflow
POST to `/workflow`

```shell script
WORKFLOW_ID=$(curl -v \
  -H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
  -H 'Content-Type: application/json' \
  $CONDUCTOR_API/workflow \
  -H 'Content-Type: application/json' \
  -d '
{
  "name": "js-example",
  "version": 1,
  "input": {
    "enter_your_name": "John"
  }
}
')
```

Check result:
```shell script
curl -v \
  -H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
  "${CONDUCTOR_API}/workflow/${WORKFLOW_ID}"
```

Output of the workflow execution should contain:
```json
{
   "result": {
      "name": "John",
      "name_length": 4
   }
}
```
### QuickJs bugs, limitations:
* Syntax errors are printed to stdout

## Python
### Create new workflow wasm-example
POST to `/metadata/workflow`

```shell script
curl -v \
-H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
-H 'Content-Type: application/json' \
${CONDUCTOR_API}/metadata/workflow -d @- << 'EOF'
{
    "name": "py-example",
    "description": "python lambdas in running in wasm",
    "ownerEmail": "example@example.com",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
            "taskReferenceName": "create_json_ref2",
            "name": "GLOBAL___py",
            "inputParameters": {
                "args": "${workflow.input.enter_your_name}",
                "outputIsJson": "true",
                "script": "import json;print(json.dumps({'name': argv[1]}));"
            },
            "type": "SIMPLE",
            "startDelay": 0,
            "optional": false,
            "asyncComplete": false
        }
    ]
}
EOF
```


### Execute the workflow
POST to `/workflow`

```shell script
WORKFLOW_ID=$(curl -v \
  -H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
  -H 'Content-Type: application/json' \
  $CONDUCTOR_API/workflow \
  -H 'Content-Type: application/json' \
  -d '
{
  "name": "py-example",
  "version": 1,
  "input": {
    "enter_your_name": "John"
  }
}
')
```

Check result:
```shell script
curl -v \
  -H "x-auth-organization: fb-test" -H "x-auth-user-role: OWNER" -H "x-auth-user-email: foo" \
  "${CONDUCTOR_API}/workflow/${WORKFLOW_ID}"
```

Output of the workflow execution should contain:
```json
{
   "result": {
      "name": "John"
   }
}
```

### Python bugs, limitations:
* Syntax errors end up having status COMPLETED instead of FAILED as status code is always 0.
* Compared to QuickJs this approach introduces 5-20x worse latency for small scripts.
* FIXME: Could not find platform dependent libraries <exec_prefix>\nConsider setting $PYTHONHOME to <prefix>[:<exec_prefix>]
* FIXME: stderr
