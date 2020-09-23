# Prepare vars
export CLUSTER_NAME="cluster.local"
export KUBESPRAY_RELEASE="release-2.14"
declare -a IPS=(192.168.1.10 192.168.1.11 192.168.1.12 192.168.1.13 192.168.1.14)

# Prepare KubeSpray
git clone https://github.com/kubernetes-sigs/kubespray.git || :
cd kubespray
git fetch --all
git checkout $KUBESPRAY_RELEASE

virtualenv -p python3 .venv
source .venv/bin/activate
pip install -r requirements.txt

# Prepare KeepaliveD role
ansible-galaxy install evrardjp.keepalived
# Prepare HAProxy role
ansible-galaxy install uoi-io.haproxy

# Copy ``inventory/sample`` as ``inventory/$CLUSTER_NAME``
cp -rfp inventory/sample inventory/$CLUSTER_NAME

# Update Ansible inventory file with inventory builder
KUBE_MASTERS_MASTERS=3 CONFIG_FILE=inventory/$CLUSTER_NAME/hosts.yml python3 contrib/inventory_builder/inventory.py ${IPS[@]}

ansible-playbook -b -i inventory/$CLUSTER_NAME/hosts.yml \
   ../setup_ha_proxy.yml \
   -e @../deploy_ansible_vars.yaml

ansible-playbook -b -i inventory/$CLUSTER_NAME/hosts.yml \
   cluster.yml \
   -e @../deploy_ansible_vars.yaml

cd ..

# Create PV for NFS server provisioner
kubectl create -f pv

#Install NFS server provisioner
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm repo update
helm upgrade --install nfs-server-provisioner \
  -n kube-system \
  stable/nfs-server-provisioner \
  -f charts/nfs-server-provisioner.yaml

# Install kubevirt (takes ~5 minutes after apply to be ready)
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/v0.28.0/kubevirt-operator.yaml
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/v0.28.0/kubevirt-cr.yaml
