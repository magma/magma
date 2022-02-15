# Mini Magma Orchestrator

Install dependant collections:
```bash
ansible-galaxy collection install community.crypto
ansible-galaxy collection install community.docker
ansible-galaxy collection install kubernetes.core
```

Copy your public SSH key to the host:
```bash
ssh-keygen -R 192.168.5.70
ssh-copy-id ubuntu@192.168.5.70
```

**Update your values in `hosts.yml` file before running the playbook.**

Deploy Magma orchestrator:
```bash
ansible-playbook deploy-orc8r.yml
```
> Note: After deployment is done it takes around 10 minutes to start all the magma services.

Create new user:
```bash
ORC_POD=$(kubectl -n orc8r get pod -l app.kubernetes.io/component=orchestrator -o jsonpath='{.items[0].metadata.name}')
kubectl -n orc8r exec -it ${ORC_POD} -- envdir /var/opt/magma/envdir /var/opt/magma/bin/accessc \
  add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator

NMS_POD=$(kubectl -n orc8r get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}')
kubectl -n orc8r exec -it ${NMS_POD} -- yarn setAdminPassword magma-test admin admin
kubectl -n orc8r exec -it ${NMS_POD} -- yarn setAdminPassword master admin admin
```

### Ansible Setup

Install Ansible - Ubuntu 20.04 LTS:
```bash
sudo apt remove ansible
sudo apt update
sudo apt install software-properties-common
sudo add-apt-repository --yes --update ppa:ansible/ansible
sudo apt install ansible -y
```

Install Ansible - macOS:
```bash
brew install ansible
```
