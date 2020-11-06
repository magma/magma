sudo apt install curl -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt update -y && sudo apt install docker-ce -y
sudo hostnamectl set-hostname kubecluster
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
sudo apt-get update -y && sudo apt-get install -y apt-transport-https && curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list && sudo apt-get update
sudo apt install -y kubeadm  kubelet kubernetes-cni
sudo kubeadm init --pod-network-cidr=10.244.20.0/24 --apiserver-advertise-address=192.168.1.111
mkdir $HOME/.k8s
sudo cp /etc/kubernetes/admin.conf $HOME/.k8s/
sudo chown $(id -u):$(id -g) $HOME/.k8s/admin.conf
export KUBECONFIG=$HOME/.k8s/admin.conf
echo "export KUBECONFIG=$HOME/.k8s/admin.conf" | tee -a ~/.bashrc
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/k8s-manifests/kube-flannel-rbac.yml
kubectl taint nodes --all node-role.kubernetes.io/master-
curl -L https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo update
