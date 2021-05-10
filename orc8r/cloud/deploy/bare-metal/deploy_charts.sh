#!/bin/bash
. orc8r_settings
for chart in orc8r fluentd elasticsearch mariadb elasticsearch metallb; do
  envsubst < charts/$chart.yaml.tpl > charts/$chart.yaml
done

kubectl create ns mariadb || :
kubectl create ns infra || :
kubectl create ns magma || :

echo "Setting up secrets and certs..."
if [ -e ../scripts/self_sign_certs.sh ]; then
  scripts_dir=../scripts
else
  test -d magma || git clone https://github.com/magma/magma.git
  scripts_dir=magma/orc8r/cloud/deploy/scripts
fi
mkdir -p certs
cd certs
if ! [ -f admin_operator.pem ]; then
  ../$scripts_dir/self_sign_certs.sh $dns_domain
  ../$scripts_dir/create_application_certs.sh $dns_domain
fi

kubectl -n $namespace create secret generic orc8r-certs \
  --from-file rootCA.pem \
  --from-file controller.key \
  --from-file controller.crt \
  --from-file certifier.key \
  --from-file certifier.pem \
  --from-file bootstrapper.key \
  --from-file admin_operator.pem \
  --dry-run=client  \
  -oyaml > ../secrets/orc8r-certs.yaml
kubectl -n $namespace create secret generic fluentd-certs \
  --from-file fluentd.key \
  --from-file fluentd.pem \
  --from-file certifier.pem \
  --dry-run=client -oyaml > ../secrets/fluentd-certs.yaml
kubectl -n $namespace create secret generic nms-certs \
  --from-file admin_operator.pem \
  --from-file admin_operator.key.pem \
  --from-file controller.key \
  --from-file controller.crt \
  --dry-run=client -oyaml > ../secrets/nms-certs.yaml
cd ..
kubectl -n $namespace apply -f secrets/

echo "Checking kube-dns service..."
if ! kubectl -n kube-system get svc kube-dns &>/dev/null; then
  echo "kube-dns service not found. Copying coredns service as coredns..."
  kubectl -n kube-system get svc coredns -oyaml | sed -e 's/name: coredns/name: kube-dns/' -e '/clusterIP:/d' | kubectl -n kube-system create -f -
fi

echo "Setting up infra helm charts..."
helm repo add stable https://charts.helm.sh/stable/
helm repo add jetstack https://charts.jetstack.io
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add elastic https://helm.elastic.co

# This should point to a repo where built charts are located
#helm repo add github-repo

kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml
sleep 3
helm -n mariadb upgrade --install mariadb bitnami/mariadb --version 7.3.14 -f charts/mariadb.yaml --wait

helm -n infra upgrade --install cert-manager jetstack/cert-manager
helm -n infra upgrade --install metallb stable/metallb -f charts/metallb.yaml --wait

echo "Setting up DB..."
envsubst < db_setup.sql.tpl > db_setup.sql
kubectl -n mariadb exec -it mariadb-master-0 -- mysql -u root --password=$db_root_password < db_setup.sql

echo "Setting up Magma helm charts..."
helm -n $namespace upgrade --install fluentd stable/fluentd -f charts/fluentd.yaml
helm -n $namespace upgrade --install elasticsearch-master elastic/elasticsearch -f charts/elasticsearch-master.yaml
helm -n $namespace upgrade --install elasticsearch-data elastic/elasticsearch -f charts/elasticsearch-data.yaml
helm -n $namespace upgrade --install elasticsearch-data2 elastic/elasticsearch -f charts/elasticsearch-data2.yaml
helm -n $namespace upgrade --install elasticsearch-curator stable/elasticsearch-curator -f charts/elasticsearch-curator.yaml
helm -n $namespace upgrade --install kibana stable/kibana -f charts/kibana.yaml

helm -n $namespace upgrade --install orc8r github-repo/orc8r -f charts/orc8r.yaml --wait

export ORC_POD=$(kubectl -n $namespace get pod -l app.kubernetes.io/component=orchestrator -o jsonpath='{.items[0].metadata.name}')
export NMS_POD=$(kubectl -n $namespace get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}')

kubectl -n $namespace exec -it ${ORC_POD} -- envdir /var/opt/magma/envdir /var/opt/magma/bin/accessc add-existing -admin -cert /var/opt/magma/certs/admin_operator.pem admin_operator || :
kubectl -n $namespace exec -it ${NMS_POD} -- yarn setAdminPassword master $admin_email $admin_password

NMS_ADDR="master.nms.$dns_domain"
echo "Magma UI: https://$NMS_ADDR"
echo "Email: $admin_email"
echo "Password: $admin_password"
