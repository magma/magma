# Handle opts
if [ -n "$1" ]; then
  playbook=deploy-$1.yml
  if [ ! -e "$playbook" ]; then
    echo "Invalid selection. Playbook file $playbook not found."
    exit
  fi
else
  playbook=deploy-all.yml
fi
# Prepare vars
export CLUSTER_NAME="cluster.local"
export KUBESPRAY_RELEASE="release-2.14"
# Optionally set IPS variable externally to define the IPs for your k8s nodes
IPS=${IPS:-92.168.0.10 192.168.0.11}

# Prepare Kubespray
git clone https://github.com/kubernetes-sigs/kubespray.git || :
cd kubespray
git fetch --all
git checkout $KUBESPRAY_RELEASE
# Copy ``inventory/sample`` as ``inventory/$CLUSTER_NAME``
mkdir -p ../inventory/$CLUSTER_NAME
cp -rnp inventory/sample ../inventory/$CLUSTER_NAME
cd ..

sudo apt install -y python3-venv || :
python3 -m venv .venv
source .venv/bin/activate
pip install -U pip
pip install -r kubespray/requirements.txt

# Prepare Keepalived role
ansible-galaxy install evrardjp.keepalived
# Prepare HAProxy role
ansible-galaxy install git+https://github.com/uoi-io/ansible-haproxy,fe397a380ad733be7d17b567b626301e1ee90089


# Update Ansible inventory file with inventory builder
HOST_PREFIX=compute KUBE_MASTERS_MASTERS=3 CONFIG_FILE=inventory/$CLUSTER_NAME/hosts.yml python3 kubespray/contrib/inventory_builder/inventory.py ${IPS[@]}

ansible-playbook -b -i inventory/$CLUSTER_NAME/hosts.yml  $playbook -e @ansible_vars.yaml

echo "If you did not specify passwords, you can find generated passwords here:"
find inventory/$CLUSTER_NAME/credentials
