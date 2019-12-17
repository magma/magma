# Magma manipulator
This is standalone tool that helps automatically register [Magma](https://github.com/facebookincubator/magma) gateways.
It works in Kubernetes cluster with virtlet only.
For example if some pod with gateway was recreated by Kubernetes this tool
will re-register the new gateway automatically in [Magma](https://github.com/facebookincubator/magma) orc8r.

## Installation
```
git clone git@github.com:vladiskuz/magma-manipulator.git
cd magma-manipulator
virtualenv -p python3 .venv
source .venv/bin/activate
python3 setup.py develop
```

## Configuration
* change *config.yml* for your purposes
* change *kconfig* regarding your k8s cluster
* run the tool magma-manipulator
* delete some pod and wait until the pod will recreate and this tool will re-register them in Magma orc8r

