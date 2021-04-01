# Prepare vars
export CLUSTER_NAME="cluster.local"
export KUBESPRAY_RELEASE="release-2.14"
declare -a IPS=(10.22.85.6 10.22.85.7 10.22.85.8 10.22.85.9 10.22.85.10 10.22.85.11)

# Prepare KubeSpray
git clone https://github.com/kubernetes-sigs/kubespray.git || :
cd kubespray
git fetch --all
git checkout $KUBESPRAY_RELEASE

sudo apt install -y python3-venv || :
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Prepare KeepaliveD role
ansible-galaxy install evrardjp.keepalived
# Prepare HAProxy role
ansible-galaxy install git+https://github.com/uoi-io/ansible-haproxy,fe397a380ad733be7d17b567b626301e1ee90089

# Copy ``inventory/sample`` as ``inventory/$CLUSTER_NAME``
cp -rfp inventory/sample inventory/$CLUSTER_NAME

# Update Ansible inventory file with inventory builder
HOST_PREFIX=compute KUBE_MASTERS_MASTERS=3 CONFIG_FILE=inventory/$CLUSTER_NAME/hosts.yml python3 contrib/inventory_builder/inventory.py ${IPS[@]}

ansible-playbook -b -i inventory/$CLUSTER_NAME/hosts.yml \
   ../setup_ha_proxy.yml \
   -e @../deploy_ansible_vars.yaml
# evrardjp.keepalived does not start service automatically
sudo service keepalived start

ansible-playbook -b -i inventory/$CLUSTER_NAME/hosts.yml \
   cluster.yml \
   -e @../deploy_ansible_vars.yaml

 cd ..
mkdir -p ~/.kube
sudo cat /root/.kube/config > ~/.kube/config

# Create PV for NFS server provisioner
kubectl apply -f pv

#Install NFS server provisioner
helm repo add stable https://charts.helm.sh/stable/
helm repo update

mkdir -p /mnt/persistentvols/nfs
dpkg -l nfs-common || sudo apt install -y nfs-common
helm upgrade --install nfs-server-provisioner \
  -n kube-system \
  stable/nfs-server-provisioner \
  -f charts/nfs-server-provisioner.yaml

# Install kubevirt (takes ~5 minutes after apply to be ready)
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/v0.28.0/kubevirt-operator.yaml
sleep 5
kubectl apply -f https://github.com/kubevirt/kubevirt/releases/download/v0.28.0/kubevirt-cr.yaml
