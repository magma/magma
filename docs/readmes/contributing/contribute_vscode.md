---
id: contribute_vscode
title: Develop Magma With VSCode
hide_title: true
---

# Develop Magma With VSCode

This document walks a new developer through the process of setting up VSCode. The primary focus of this guide is for AGW development.

There are two types of workspaces supported by the Magma team

- Remote SSH workspace with Magma VM: A workspace for working with files in the Magma VM from your local machine. [Jump to Setup Remote SSH workspace with Magma VM](#setup-remote-ssh-workspace-with-magma-vm)
- Devcontainer workspace: A Docker based development flow that is best suited for changes that only require unit testing. It is natively supported with GitHub Codespaces. However, it is currently only supported for AGW C/C++ code. [Setup Devcontainer workspace with GitHub Codespaces](#setup-devcontainer-workspace-with-github-codespaces)

Visual Studio Code is available for downloads at [their official website](https://code.visualstudio.com).

## Setup Remote SSH workspace with Magma VM

The steps below need to only be done once. After your workspace is setup once, it is easily accessible via **File->Open Recent**.

### Setup default extensions for SSH workspace

Open VSCode and use **Command+Shift+P** to open the editor preferences and select **Preferences: Open Settings (JSON)**. This will open the user settings file for VSCode. Insert the following into the file. This configuration will take effect the next time VSCode connects to a remote host.

```json
    "remote.SSH.defaultExtensions": [
        "llvm-vs-code-extensions.vscode-clangd",
        "stackbuild.bazel-stack-vscode",
        "coolchyni.beyond-debug",
        "stackbuild.bazel-stack-vscode-cc",
        "augustocdias.tasks-shell-input",
        "ryuta46.multi-command",
        "ms-python.python",
        "njpwerner.autodocstring",
        "ms-python.vscode-pylance",
        "pucelle.run-on-save",
    ],
```

Unlike with Devcontainer settings, there is no way to configure default extensions ouside of user settings. It may be good to periodically check this section for any updated extensions that should be installed.

### Add the VM SSH config to ~/.ssh/config

**Prerequisite**: Have the Magma VM running.

Go to `$MAGMA_ROOT/lte/gateway` and run `vagrant ssh-config`. Copy the output into `~/.ssh/config`.
You SSH config should now look like this:

```text
...

Host magma
  HostName 127.0.0.1
  User vagrant
  Port 2222
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /Users/<user_name>/.vagrant.d/insecure_private_key
  IdentitiesOnly yes
  LogLevel FATAL
  ForwardAgent yes
```

### Create a remote SSH workspace

A more detailed documentation is available on [the official VSCode doc on remote SSH](https://code.visualstudio.com/docs/remote/ssh).

Use **Command+Shift+X** to open the extensions tab and install the **[Remote SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh)** extension. Once installed, open the command palette with **Command+Shift+P** and select **Remote-SSH: Connect to Host**. Select *magma* to start a remote SSH workspace.

Inside the newly created VSCode window, go to **File → Open Workspace** and select  `/home/vagrant/magma/vscode-workspaces/workspace.magma-vm-workspace.code-workspace`.

## Setup Devcontainer workspace with GitHub Codespaces

**Prerequisite**: Have access to [GitHub Codespaces](https://github.com/features/codespaces). If you do not have access, follow the instructions at `.devcontainer/README.md` to open it using Docker.

If you open a GitHub Codespaces from a Magma repo Pull Request, the devcontainer configuration should be loaded up automatically. You can also open a workspace by going to `https://github.com/codespaces` and selecting `New Codespaces`. If you intend to make changes and submit a Pull Request, always select your fork of Magma.

Any local changes made inside the devcontainer will be wiped when it is upgraded. This should not happen often and you will always be given a prompt before the devcontainer upgrades. To have persisting personal configs, refer to the GitHub's doc on [Personalizing Codespaces for your account](https://docs.github.com/en/codespaces/customizing-your-codespace/personalizing-codespaces-for-your-account).

## Language and tool specific descriptions of what is enabled in VSCode

We've enabled several extensions and custom configurations to improve developer experience, refer to the sections below for language and tool specific commands and features.

### C / C++ code completion and navigation

> Supported in both Remote SSH + Magma VM and Devcontainer

We utilize [clangd](https://clangd.llvm.org) to enable smart code insights. Clangd searches for a `compile_commands.json` file that serves as a compilation database at $MAGMA_ROOT. Most often, clangd will need to be restarted when the compile_commands.json is modified.

The compilation database can be gernated with both CMake and Bazel. The Bazel method will be described more in detail below. With CMake, it is not possible to generate a single compilation database for all C/C++ targets, so you will have generate one at a time and restart clangd.

To generate the comilation database, simply run `make build_oai`, or any other C/C++ target. Then run a task to symlink the root compilation database file to the newly generated one with the following steps.

  1. Open the command palette (**Command+Shift+P**)
  2. Select **Tasks: Run Task** and choose **Set compile_commands.json for IntelliSense**
  3. Select a `compile_commands.json` to that was generated
  4. After the task completes, restart clangd

The `compile_commands.json` at project root should contain a compilation instruction for each file in the target you specified. If you only see partial results, it might help to run the same make command a few more times.

Once clangd is retarted, you should see some note about indexing in the bottom of your editor. Code completion and navigation should work after the indexing completes.

#### Dealing with clangd extension errors

If you see errors about clangd, try the following:

  1. For errors about `clangd` not being found, try running **clangd: Download language server** from the command palette
  2. For errors about the extension commands not being found, try running **clangd: Manually activate extension** to start the extension

### Python code completion and navigation

> Supported only in Remote SSH + Magma VM

First SSH into your VM and run `magitvate`. This will setup a virtual Python environment we will need.

Follow the following steps to enable code completion and navigation

  1. Open the command palette (**Command+Shift+P**)
  2. Select **Python: Select interpreter** and choose `magma` and the default Python interpreter (`~/build/python/bin/python`)
  3. Run **Python: Build Workplace Symbols** via the command palette

Additional features such as formatting on save is enabled for the Python source files.

### Building and testing with Bazel

> Supported in both Remote SSH + Magma VM and Devcontainer

#### Build specific targets and unit tests via codelens

The **bazel-stack-vscode** extension adds [codelens](https://code.visualstudio.com/blogs/2017/02/12/code-lens-roundup) directly into `BUILD.bazel` files. Utilizing this makes building and testing as easy as clicking a button.
For example, to run a single unit test for SessionD, open `lte/gateway/c/session_manager/test/BUILD.bazel` and click the `test` codelens. Similarly, click the `build` codelens to build only.

![SessionD Unit Test Codelens](assets/contributing/sessiond-unit-test-codelens.png)

At the top of each `BUILD.bazel` file, there is a codelens to build and test all targets in the file.

#### Generate compilation database with Bazel

To generate the compilation database with Bazel, run **Command+Shift+P** to open the command palette and select **Multi command: Execute multi command**. Select the command **sentry_generateCcWithBazelAndRestartClangderror**. This is a wrapper command that runs two extension commands: `bsv.cc.compdb.generate` (**Bzl: Bazel/C++: Generate Compilation Database**) and then `clangd.restart` (**clangd: Restart language server**).

This compilation database will contain necessary information for all C / C++ targets available to be built with Bazel.

Refer to [Dealing with clangd extension errors](#dealing-with-clangd-extension-errors) for debugging extension issues.

#### Run unit tests with GDB

> Currently only available for SessionD unit tests under `lte/gateway/c/session_manager/tests`, but it is easy to add configurations to enable it for any `cc_test`. Run `bazel query 'kind(cc_test,//...)'` to get the full list of available targets. Modify `.vscode/tasks.json` and `.vscode/launch.json` to enable GDB debugging for other targets.

Run **Command+Shift+D** to open the debug tab. In the drop down menu at the top of the tab, select **(Remote SSH) Run SessionD test with GDB** and press the gree arrow. This will open up a new drop down menu with all SessionD unit test targets. Once a test is selected, VSCode will build the target in debug mode and launch the test with GDB.

![SessionD Start Debug](assets/contributing/sessiond-start-debug.png)

Once the task is launched, the test will start execution immeditately. It is recommended to add a breakpoint before triggering the debugger to halt the execution.

To add a breakpoint, simply click on the left most edge of the code to add a red circle.

![SessionD Breakpoint Code](assets/contributing/sessiond-breakpoint-code.png)

With a breakpoint added, the debug console will show when the breakpoint is hit. Finally, use the debug console like a normal GDB console to aid your testing!

![SessionD GDB List](assets/contributing/sessiond-gdb-list.png)
