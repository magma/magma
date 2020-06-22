# Lambda workers (js & python) executed using WASM engine

Run lambda tasks inside a web assembly engine [wasmer](https://wasmer.io/).
Every execution spawns a new short lived process.

## Usage

Currently supports two task types: `GLOBAL___js` and `GLOBAL___py`.

* `lambdaValue` - convention for storing task inputs. E.g. workflow input: `${workflow.input.enter_your_name}`,
result of previous task: `${some_task_ref.output.result}`
* `outputIsJson` - if set to `true`, output is interpreted as JSON and
task will be marked as failed if parsing fails. If `false`, output is interpreted as plaintext.
Any other value, including empty one, means output will be parsed as JSON, will fallback
to plaintext on parsing failure.
* `scriptExpression` - script to be executed

## Javascript engine
Task `GLOBAL___js` uses [QuickJs](https://bellard.org/quickjs/) engine, compiled to wasm [(demo)](https://wapm.io/package/quickjs).

### APIs
Task result is written using `console.log` or by `return`ing the value (preferred).

Log messages are written using `log` or `console.error`.

Input data is available in `$` global variable.
Use `$.lambdaValue` to get task input.
This is backwards compatibile with
[Lambda tasks](https://netflix.github.io/conductor/configuration/systask/#lambda-task)

## Python interpreter
Task `GLOBAL___py` uses CPython 3.6 compiled to wasm [(demo)](https://wapm.io/package/python).

### APIs
Task result is written using `print` or by `return`ing the value (preferred).

Log messages are written using `log` or `eprint`.

Input data is available in `inputData` global variable.
Use `inputData["lambdaValue"]` to get task input.

## Example workflow

This example asks user for name, then executes:
* python task:
lambdaValue: `${workflow.input.enter_your_name}`
```python
log('logging from python')
name = inputData['lambdaValue']
return {'name': name}
```
* javascript task:
lambdaValue: `${create_json_ref.output.result}`
```javascript
log('logging from js');
var result = $.lambdaValue;
result.name_length = (result.name||'').length;
return result;
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
                "lambdaValue": "${workflow.input.enter_your_name}",
                "outputIsJson": "true",
                "scriptExpression": "log('logging from python')\nname = inputData['lambdaValue']\nreturn {'name': name}\n"
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
                "lambdaValue": "${create_json_ref.output.result}",
                "outputIsJson": "true",
                "scriptExpression": "log('logging from js');\nvar result = $.lambdaValue;\nresult.name_length = (result.name||'').length;\nreturn result;\n"
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
* Compared to QuickJs this approach introduces 5-200x worse latency for small scripts: ~30ms for QuickJs, ~.5s for Python
* Python needs writable lib directory, thus a temp directory needs to be created/deleted for each execution
