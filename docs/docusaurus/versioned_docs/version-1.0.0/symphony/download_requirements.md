---
id: version-1.0.0-symphony_download_requirements
title: Download Requirements
hide_title: true
original_id: symphony_download_requirements
---
# Download Requirements

If you are going to run or build the Symphony Agent you will need to install Docker and Docker Compose on your machine. Here are some guides for how to do it.

## Docker
### MacOS (OSX)
**WARNING:**
There are some known limitations to running Docker on MacOS, the biggest being that you are unable to run ping over IPv6 from inside these Docker containers. But if you don't need your agent to do these things, this should be fine.


With MacOS you actually **don't want to install Docker via Homebrew**, but if you go to the docker website it will ask you to create an account to login before installing. To avoid having to create an account, though, you can use this [direct link](https://download.docker.com/mac/stable/Docker.dmg) to download a dmg image file instead. (You can read more background on this issue [here](https://github.com/docker/docker.github.io/issues/6910))

Once you have installed Docker via the dmg file you will want to increase the Docker disk (~50GB) and memory (at least 5GB) in the Docker desktop UI. You can do this by going to the Disk and Advanced tabs in Docker's preferences window.

### Linux
For a Fedora based system, this is as simple as:

```bash
sudo dnf install docker
```

## Docker Compose
### MacOS
This should already be included in the Docker desktop install.

### Linux
For a Fedora based system, run:
```bash
sudo dnf install docker-compose
```
