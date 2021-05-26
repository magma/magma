---
id: contribute_onboarding
title: Developer Onboarding
hide_title: true
---

# Developer Onboarding

This document walks a new developer through the process of onboarding to the Magma project.

## Project overview

Magma is an open-source software platform that gives network operators an open, flexible, and extendable mobile core network solution. Our mission is to connect the world to a faster network by enabling service providers to build cost-effective and extensible carrier-grade networks.

In more approachable terms, Magma is a collection of software that makes running things like a cellular network affordable and customizable. With Magma's current rate of growth, increasing stability and reliability of the platform are all huge wins.

There are three main components in the Magma software platform: Access Gateway (AGW), Orchestrator (Orc8r), and Federation Gateway (FeG). If you want to read more, see the [Magma intro docs](https://magma.github.io/magma/docs/next/basics/introduction.html).

### Access Gateway

The AGW provides network services and policy enforcement. In an LTE network, the AGW implements an evolved packet core (EPC), and a combination of an AAA and a PGW. It works with existing, unmodified commercial radio hardware.

More generally, the AGW defines datapath rules for connecting subscribers through to the Internet. It pulls configuration from Orc8r, sets datapath rules, manages charging and accounting, reports state and metrics to Orc8r, and more.

### Orchestrator

The Orc8r is a centralized controller for a set of networks. In [SDN](https://en.wikipedia.org/wiki/Software-defined_networking) terms, Orc8r handles the control plane. This means one Orc8r serves many gateways -- pushing configuration to the gateways and pulling state and metrics from the gateways.

One of the functions of the Orc8r is to expose a management [REST API](https://restfulapi.net/). The REST API is defined as an [OpenAPI](https://swagger.io/solutions/getting-started-with-oas/) [specification](https://swagger.io/specification/) (aka [Swagger](https://swagger.io/blog/api-strategy/difference-between-swagger-and-openapi/) specification), and made available to operators over [mutually-authenticated](https://comodosslstore.com/blog/what-is-ssl-tls-client-authentication-how-does-it-work.html) HTTPS. The Orc8r also exposes a series of [gRPC](https://grpc.io/) services which gateways can call to report state, receive config updates, etc.

### Federation Gateway

The FeG serves as the interface to existing operator cores, affording a modern, flexible interface on top of existing protocols and architectures.

## Points of contact

Check out the [Community](https://www.magmacore.org/community/) page for community entry points. The principal resources are

- [GitHub Discussions page](https://github.com/magma/magma/discussions). Best place to ask questions if you're stuck or want more context.
- Community Slack channel. Say hello in `#general` and let us know you're onboarding to the project.

## Ramp-up resources

Briefly check out the resources below

- [Magma documentation](https://magma.github.io/magma/docs/next/basics/introduction.html) (aka Docusaurus)
- [Magma GitHub repo](https://github.com/magma/magma)
- External guides
    - Go
        - [Effective Go](https://golang.org/doc/effective_go.html)
        - [Video tutorial](https://www.youtube.com/watch?v=YS4e4q9oBaU&ab_channel=freeCodeCamp.org)
    - IntelliJ
        - [Go IDE overview](https://www.youtube.com/watch?v=o3igXAE9eDo&ab_channel=JetBrainsTV)
    - [Protocol buffers](https://developers.google.com/protocol-buffers)
    - [gRPC](https://grpc.io/)
    - [Swagger ecosystem](https://swagger.io/)

## Getting started

### Set up development environment

Part of setting up the development environment involves updating your [shell rc files](https://superuser.com/questions/183870/difference-between-bashrc-and-bash-profile#183980). When we say "add to your shell rc file", if you don't know what that means, then [add it to your `~/.bash_profile`](https://joshstaiger.org/archives/2005/07/bash_profile_vs.html) (or `~/.zshrc` if you're using zsh).

Note: this guide assumes you have access to an IntelliJ Ultimate Edition license. The majority of the functionality should still work without the license. Also, plenty of Magma developers use VS Code, and a minority also use other IDEs.

**Howto**

1. [Import](https://support.google.com/chrome/answer/96816?hl=en) @hcgatewood's Google Chrome bookmarks: [bookmarks-hcgatewood.html](https://www.dropbox.com/s/rvhcofsrkpvkbfm/bookmarks-hcgatewood.html?dl=0). These bookmarks provide a starting point for accessing resources across the Magma ecosystem.
2. Add the following to your shell rc file, then restart your terminal

```
# MAGMA_ROOT denotes the root of the Magma repo
export MAGMA_ROOT=~/magma

# noti sends notification based on previous command success/fail
# Default sound is "default"
function noti() {
    if [[ $? = 0 ]] ; then
        title="✅" ; sound="ping"
    else
        title="❌" ; sound="basso"
    fi
    terminal-notifier -title "$title" -message "Finished at $(date +%X)" -sound "$sound"
}
```

3. Install the [Homebrew package manager](https://brew.sh/)
4. Install [terminal-notifier](https://github.com/julienXX/terminal-notifier) via `brew install terminal-notifier`
5. Install [IntelliJ Ultimate Edition](https://www.jetbrains.com/idea/) via `brew install --cask intellij-idea`
6. [Import](https://www.jetbrains.com/help/idea/sharing-your-ide-settings.html#import-export-settings) @hcgatewood's IntelliJ settings: [intellij-hcgatewood.zip](https://www.dropbox.com/s/2i38wrfrfhjyicz/intellij-hcgatewood.zip?dl=0)
7. Clone the [Magma repository](https://github.com/magma/magma) via `git clone git@github.com:magma/magma.git ${MAGMA_ROOT}`

### Install Magma locally

Install Magma locally and get everything running.

**Howto**

1. Follow the [prerequisites guide](https://magma.github.io/magma/docs/next/basics/prerequisites) and install all development tools, up to but not including the "Build/Deploy Tooling" section
2. (Optional) If you opt to use IntelliJ IDEA as your local IDE, follow the instructions in the [Set up IntelliJ](#set-up-intellij) section below before you proceed
3. Run all Orc8r tests
    1. Via Docker build script: `cd ${MAGMA_ROOT}/orc8r/cloud/docker && ./build.py -t ; noti`
    2. [Via IntelliJ](https://magma.github.io/magma/docs/orc8r/dev_testing#testing-tips)
4. Follow the [quick start guide](https://magma.github.io/magma/docs/next/basics/quick_start_guide) to get an AGW and Orc8r instance running on your dev machine
5. Visit the local [Swagger UI](https://swagger.io/tools/swagger-ui/) view of our REST API (URL is in @hcgatewood's Google Chrome bookmarks) and [list the set of managed networks](https://localhost:9443/apidocs/v1/#/Networks/get_networks) -- there should be one named "test"
    - You will need to toggle a Google Chrome preference to [allow insecure localhost](https://superuser.com/questions/772762/how-can-i-disable-security-checks-for-localhost)

Note: remember to periodically call `docker system prune` to clear outdated Docker artifacts from your dev machine.

### Set Up IntelliJ
We recommend using [IntelliJ IDEA](https://www.jetbrains.com/idea/) for general Magma development, or [Visual Studio Code](https://code.visualstudio.com/) for a free alternative.

For IntelliJ IDEA, we provide a set of run configurations that support rapidly testing Magma code. See [Testing Tips](https://magma.github.io/magma/docs/orc8r/dev_testing#testing-tips) for more details.

To set up your local IntelliJ environment, perform the following
1. After cloning the Magma repo, open the directory in IntelliJ
2. Ensure the [Go plugin](https://plugins.jetbrains.com/plugin/9568-go) has been installed by going to `Preferences > Plugins > search for the plugin "Go"`
3. [Specify the location of the Go SDK](https://www.jetbrains.com/help/idea/quick-start-guide-goland.html#step-1-open-or-create-a-project) by going to `Preferences > Languages & Frameworks > Go > GOROOT` and selecting the relevant location
4. Create a Go module for the project by going to `Files > Project Structure > Project Settings > Modules > Click on "+" sign in the toolbar > New Module > Next`. When you reach the new module creation page, enter the following information:
    - Module name: `magma`
    - Content root, module file location: full path to your local Magma clone, e.g. `/Users/your_username/magma`

After completing the above steps, restart your IDE and ensure the environment is properly set up

1. Open "Project" on the left toolbar, and display "Project Files". All the files in the root `magma` directory should be displayed *without* a yellow background. This indicates IntelliJ recognizes the files as part of the module.
2. At the top-right corner of your IDE, you should see a drop-down menu showing a list of run configurations for the Magma test suites, with a green triangular button that allows you to run the selected test. Alternatively, when you open your run configurations (`Run > Edit Configurations`), you should see something like the below

![intellij_initial_run_configs](assets/intellij_initial_run_configs.png)

You can now run all (Orchestrator) tests in one click.

NOTE: a minority of tests require a running Postgres instance, and will otherwise fail with connection errors. You can use `orc8r/cloud/docker/run.py` to spin up the required test DB.

### Complete a starter task

If you haven't already, join our community Slack channel via the link from the [Community](https://www.magmacore.org/community/) page and say hello in `#general`. We can point you toward a good starter task. You can also use the [`good first issue` tag on our GitHub repo to search for good starter tasks](https://github.com/magma/magma/labels/good%20first%20issue).

As you're working on your starter task, refer to [Development Workflow](./contribute_workflow.md) for guidance on how to author a pull request and [Contributing Conventions](./contribute_conventions.md) for Magma conventions. Once your pull request is approved and merged, you're now an official contributor to the Magma project!

### Next steps

The community Slack workspace is a great place to get connected. From there, we can coordinate choosing appropriate tasks given your skill set and the Magma roadmap.

As you get acquainted with the codebase, consider the following sources of documentation

- Magma Docusaurus. This site! First stop for documentation.
- `doc.go` files. Many Go packages have a `doc.go` file with a summary of the package's functionality.
- Tests. Tests provide both testing and documentation of expected functionality.

