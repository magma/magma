---
id: prerequisites
title: Prerequisites
hide_title: true
---
# Prerequisites

We support MacOS and Linux host operating systems, and developing on Windows
should be possible but has not been tested.

* Orc8r and federated gateway development is done using 
[Docker](https://www.docker.com/get-started), and Docker Compose.

* We develop the access gateway and its components on virtual machines managed by
[Vagrant](https://www.vagrantup.com/). This helps us ensure that every
developer has a consistent development environment for the gateway.
First, install [Virtualbox](https://www.virtualbox.org/wiki/Downloads) and 
[Vagrant](http://www.vagrantup.com/downloads.html).

* Then, install some additional prereqs (replace `brew` with your OS-appropriate
package manager as necessary):

```console
$ brew install python3
$ pip3 install ansible fabric3 requests PyYAML
$ vagrant plugin install vagrant-vbguest
```

## Recommended versions:

The following are the versions of the various tools that we have used, and hence recommend, to develop Magma:

* vagrant: 2.2.3
* virtualbox: 5.2.22
* ansible: 2.7.8
* docker: 18.09.2
* docker-compose: 1.23.2
