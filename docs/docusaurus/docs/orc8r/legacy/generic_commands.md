---
id: version-1.0.0-generic_commands
sidebar_label: Generic commands framework
title: Using the generic command framework
hide_title: true
original_id: generic_commands
---
# Using the generic command framework
The generic command framework is a way to define commands in the gateway without having to implement all aspects of a command (such as the cloud implementation and handler). Instead, we can define the commands on the gateway, and call the REST endpoint `/networks/{network_id}/gateways/{gateway_id}/command/generic` with the command name and parameters to execute the command.

## Explanation

In `magmad.yml`, the generic command framework with the field `generic_command_config`:

```
# magmad.yml

generic_command_config:
  module: ...
  class: ...
```

When the gateway starts up, it looks for the class inside that module and sets it as the command executor. The purpose of the command executor is to hold a table that maps command names to async functions, execute those functions, and return the output. When the RPC method `GenericCommand` is called, it will attempt to execute the command, returning the response if successful or an error if a timeout or exception occurs. We can configure the timeout by providing a `timeout_secs` field within `generic_command_config` to determine the number of seconds before a command times out.

On the cloud side, we can send a POST request to `/networks/{network_id}/gateways/{gateway_id}/command/generic` with the following request body:

```
{
    "command": "name_of_command",
    "params": {
        ...
    }
}
```

## Using the Shell Command Executor

The default provided command executor is `ShellCommandExecutor`, which resides in the module `magma.magmad.generic_command.shell_command_executor`. This command executor will run shell commands located in `generic_command_config`. For each shell command, we can decide if we want to allow parameters with the field `allow_params`, which will treat the command as a format string. Parameters are read from the field `shell_params` ( a list of parameters) within the `params` field of the request. For example:

```
# magmad.yml

generic_command_config:
  module: magma.magmad.generic_command.shell_command_executor
  class: ShellCommandExecutor
  shell_commands:
    - name: tail_syslog
      command: "sudo tail /var/log/syslog -n 20"
    - name: echo
      command: "echo {}"
      allow_params: True
```

We can then use the API endpoint to execute the command. For example:

```
POST `/networks/{network_id}/gateways/{gateway_id}/command/generic

{
    "command": "echo",
    "params": {
        "shell_params": ["Hello world!"]
    }
}`
```

We then get a response with the return code, stderr, and stdout:

```
{
  "response": {
    "returncode": 0,
    "stderr": "",
    "stdout": "Hello world!\n"
  }
}
```

## Creating a custom command executor

The generic command framework is designed to be extensible. If we want more complex functionality, we can define our own command executor and configure `magmad.yml` to use it instead.

All command executors must be an instance of `CommandExecutor`, which is the abstract base class for all command executors.

```
# magma/magmad/generic_command/command_executor.py

class CommandExecutor(ABC):
    """
    Abstract class for command executors
    """
    def __init__(
            self,
            config: Dict[str, Any],
            loop: asyncio.AbstractEventLoop,
    ) -> None:
        self._config = config
        self._loop = loop

    async def execute_command(
            self,
            command: str,
            params: Dict[str, ParamValueT],
    ) -> Dict[str, Any]:
        """
        Run the command from the dispatch table with params
        """
        result = await self.get_command_dispatch()[command](params)
        return result

    @abstractmethod
    def get_command_dispatch(self) -> Dict[str, ExecutorFuncT]:
        """
        Returns the command dispatch table for this command executor
        """
        pass
```

Command executors must provide a method `get_command_dispatch`, which returns the dispatch table of commands. We can then add our own commands to the dispatch table. Command functions are coroutines that take in a dictionary of parameters as arguments and return a dictionary as a response.

```
class CustomCommandExecutor(CommandExecutor):
    def __init__(
            self,
            config: Dict[str, Any],
            loop: asyncio.AbstractEventLoop,
    ) -> None:
        ...
        self._dispatch_table = {
            "hello_world": self._handle_hello_world
        }

    def get_command_dispatch(self) -> Dict[str, ExecutorFuncT]:
        return self._dispatch_table
        
    async def _handle_hello_world(
            self,
            params: Dict[str, ParamValueT]
    ) -> Dict[str, Any]:
        return {
            "hello": "world!"
        }
```


