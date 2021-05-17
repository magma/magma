---
id: healthchecker
title: Health Checker
sidebar_label: Health Checker
hide_title: true
---
# Health Checker
Health checker reports:
* Gateway - Controller connectivity
* Status for all the running services
* Number of restarts per each service
* Number of errors per each service
* Internet and DNS status
* Kernel version
* Magma version

# Usage
```bash
docker-compose exec magmad bash
health_cli.py
```