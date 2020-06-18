# Lambda workers (js & python) executed using WASM engine

Run lambda tasks inside a web assembly engine [wasmer](https://wasmer.io/).
Every execution spawns a new short lived process.

## Usage

Currently supports two task types: `GLOBAL___js` and `GLOBAL___py`.

* `args` - arguments to the script. To use workflow input: `${workflow.input.enter_your_name}`
To read from previous task: `${some_task_ref.output.result}`
* `outputIsJson` - if set to `true`, output will be interpreted as JSON. Otherwise interpreted as plaintext.
* `scriptExpression` - script to be executed

## Javascript engine
Task `GLOBAL___js` uses [QuickJs](https://bellard.org/quickjs/) engine, compiled to wasm [(demo)](https://wapm.io/package/quickjs).

### APIs
Task result is written using `console.log`.

Log messages are written using `console.error`.

Arguments are available in `argv` global variable.

## Python interpreter
Task `GLOBAL___py` uses CPython 3.6 compiled to wasm [(demo)](https://wapm.io/package/python).

### APIs
Task result is written using `print`.

Log messages are written using `eprint`.

Arguments are available in `argv` global variable.

## Example workflow

This example asks user for name, then executes:
* python task:
```python
import json
print(json.dumps({'name': argv[1]}))
eprint('logging from python')
```
* javascript task:
```javascript
let json=JSON.parse(argv[1]);
json.name_length = json.name.length;
console.log(JSON.stringify(json));
console.error('logging from js');
```

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
${CONDUCTOR_API}/metadata/workflow -d @- << 'EOF'
{
    "name": "wasm-example",
    "description": "python and javascript lambdas running in wasm",
    "ownerEmail": "example@example.com",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
            "taskReferenceName": "create_json_ref",
            "name": "GLOBAL___py",
            "inputParameters": {
                "args": "${workflow.input.enter_your_name}",
                "outputIsJson": "true",
                "scriptExpression": "import json\nprint(json.dumps({'name': argv[1]}))\neprint('logging from python')"
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
                "scriptExpression": "let json=JSON.parse(argv[1]);\njson.name_length = json.name.length;\nconsole.log(JSON.stringify(json));\nconsole.error('logging from js');"
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
  "name": "wasm-example",
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

### Python bugs, limitations:
* Compared to QuickJs this approach introduces 5-20x worse latency for small scripts.
