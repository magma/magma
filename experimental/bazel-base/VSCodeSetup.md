# Setting Up VSCode For Magma VM

This document describes configrations needed to enable development workflow with Bazel and Magma VM.

These steps need to only be done once. After your workspace is setup once, it is easily accessible via **File->Open Recent**.

## Setup default extensions for SSH workspace

Open VSCode and use **Command+Shift+P** to open the editor preferences and select **Preferences: Open Settings (JSON)**. This will open the user settings file for VSCode. Insert the following into the file. This configuration will take effect the next time VSCode connects to a remote host.

```
    "remote.SSH.defaultExtensions": [
        "llvm-vs-code-extensions.vscode-clangd",
        "mitaki28.vscode-clang",
        "stackbuild.bazel-stack-vscode",
        "coolchyni.beyond-debug",
        "stackbuild.bazel-stack-vscode-cc",
        "augustocdias.tasks-shell-input",
    ],
```

It may be good to periodically check this section for any updated extensions that should be installed.

## Add the VM SSH config to ~/.ssh/config

**Prerequisite**: Have the Magma VM running.

Go to `$MAGMA_ROOT/lte/gateway` and run `vagrant ssh-config`. Copy the output into `~/.ssh/config`.
You SSH config should now look like this:

```
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

## Create a remote SSH workspace

Official VSCode doc on remote SSH: https://code.visualstudio.com/docs/remote/ssh

Use **Command+Shift+X** to open the extensions tab and install the **[Remote SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh)** extension. Once installed, open the command palette with **Command+Shift+P** and select **Remote-SSH: Connect to Host**. Select *magma* to start a remote SSH workspace. 

Inside the newly created VSCode window, go to **File → Open Workspace** and select  `/home/vagrant/magma/vscode-workspaces/workspace.magma-vm-workspace.code-workspace` .
