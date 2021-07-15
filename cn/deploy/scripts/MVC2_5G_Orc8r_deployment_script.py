#!/usr/bin/env python3

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import fileinput
import getopt
import logging
import os
import platform
import re
import shutil
import socket
import subprocess
import sys
import threading
import time
import webbrowser

# Initialize the lock
lock = threading.Lock()
# Dictionary to maintain k8s services whether they are Running or not
k8s_obj_dict = {}
# Get the Current Working Directory
CWD = os.getcwd()
# Path for Orc8r temperary files
ORC8R_TEMP_DIR = '/tmp/Orc8r_temp'
INFRA_SOFTWARE_VER = os.path.join(ORC8R_TEMP_DIR, 'infra_software_version.txt')
K8S_GET_DEP = os.path.join(ORC8R_TEMP_DIR, 'k8s_get_deployment.txt')
K8S_GET_SVC = os.path.join(ORC8R_TEMP_DIR, 'k8s_get_service.txt')
# Path for Orc8r VM temperary files
ORC8R_VM_DIR = '/tmp/Orc8r_vm'
K8S_GET_OBJ = os.path.join(ORC8R_VM_DIR, 'k8s_get_objects.txt')
# Path for Templates directory where all source yaml files present
TEMPLATES_DIR = os.path.join(CWD, '../helm/templates')
# Debian-9-openstack-amd64.qcow2 file
DEBIAN_QCOW2_FILE = os.path.join(TEMPLATES_DIR, 'debian-9-openstack-amd64.qcow2')
# Path for multus-cni home directory
MULTUS_DIR = os.path.join(TEMPLATES_DIR, 'multus-cni')


class Error(Exception):
    """Base class for other exceptions"""
    pass


class NotInstalled(Error):
    """Raised when Installation not done"""
    pass


def Code(type):
    switcher = {
        'WARNING': 93,
        'FAIL': 91,
        'GREEN': 92,
        'BLUE': 94,
        'ULINE': 4,
        'BLD': 1,
        'HDR': 95,
    }
    return switcher.get(type)

# Print messages with colours on console


def myprint(type, msg):
    code = Code(type)
    message = '\033[%sm \n %s \n \033[0m' % (code, msg)
    print(message)

# Executing shell commands via subprocess.Popen() method


def execute_cmd(cmd):
    process = subprocess.Popen(cmd, shell=True)
    os.waitpid(process.pid, 0)

# Checking pre-requisites like kubeadm, helm should be installed before we run this script


def check_pre_requisite():
    # Setting logging basic configurations like severity level=DEBUG, timestamp, function name, line numner
    logging.basicConfig(
            format='[%(asctime)s %(levelname)s %(name)s:%(funcName)s:%(lineno)d] %(message)s',
            level=logging.DEBUG,
    )
    uname = platform.uname()
    logging.debug('Operating System : %s' % uname[0])
    logging.debug('Host name : %s' % uname[1])
    if os.path.exists(ORC8R_TEMP_DIR):
        shutil.rmtree(ORC8R_TEMP_DIR)
    os.mkdir(ORC8R_TEMP_DIR)
    cmd = 'cat /etc/os-release > %s' % INFRA_SOFTWARE_VER
    execute_cmd(cmd)
    with open(INFRA_SOFTWARE_VER) as fop1:
        all_lines = fop1.readlines()
        for distro_name in all_lines:
            if "PRETTY_NAME" in distro_name:
                logging.debug("Distro name : %s" % distro_name.split('=')[1])
    logging.debug('Kernel version : %s' % uname[2])
    logging.debug('Architecture  : %s' % uname[4])
    logging.debug('python version is : %s' % sys.version)
    try:
        cmd = 'kubeadm version > %s' % INFRA_SOFTWARE_VER
        out = os.system(cmd)
        if out == 0:
            myprint("GREEN", "kubeadm installed : YES")
            with open(INFRA_SOFTWARE_VER) as fop2:
                kubeadm_version = fop2.readline().split(' ')
                logging.debug("%s %s %s" % (kubeadm_version[2].split('{')[1], kubeadm_version[3], kubeadm_version[4]))
        else:
            raise NotInstalled
    except NotInstalled:
        print("kudeadm is not installed")
        myprint("FAIL", "kubeadm installed : NO")
    try:
        cmd = 'helm version > %s' % INFRA_SOFTWARE_VER
        out = os.system(cmd)
        if out == 0:
            myprint("GREEN", "HELM installed : YES")
            with open(INFRA_SOFTWARE_VER) as fop3:
                helm_version = fop3.readline().split(',')
                logging.debug("%s" % helm_version[0].split('{')[1])
        else:
            raise NotInstalled
    except NotInstalled:
        print("Helm is not installed")
        myprint("FAIL", "HELM installed : NO")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "        installing kubevirt and cdi")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")

# Delete files if exits


def del_files(file):
    if os.path.exists(file):
        os.remove(file)

# Un-installing all the k8s objects and deleting the temperary files in the path /tmp/Orc8r_temp/


def un_install(pwd):
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Uninstalling Orc8r monitoring stack ")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "*****Trying to Un-install Helm Charts*****")
    execute_cmd("helm uninstall prometheus stable/prometheus-operator --namespace kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_thanosrulers.yaml -n kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_servicemonitors.yaml -n kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_prometheusrules.yaml -n kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_prometheuses.yaml -n kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_podmonitors.yaml -n kubevirt")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_alertmanagers.yaml -n kubevirt")
    myprint("BLUE", "*****Trying to Cleanup the temporay files & Directories created as part of installation*****")
    del_files(INFRA_SOFTWARE_VER)
    del_files(K8S_GET_DEP)
    del_files(K8S_GET_SVC)
    if os.path.exists(ORC8R_TEMP_DIR):
        shutil.rmtree(ORC8R_TEMP_DIR)
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Orc8r monitoring stack Uninstalled successfully")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")

# Get magmadev VM IP


def get_magmadev_vm_ip():
    cmd = "kubectl get vmi -n kubevirt | awk '{print $1, $4}'"
    data = subprocess.Popen([cmd], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, shell=True)
    stdout, stderr = data.communicate()
    vmi_list = stdout.strip().decode("utf-8").split("\n")
    for vmi in vmi_list:
        if "magmadev" in vmi:
            return vmi.split(" ")[1]

# Deleting route information


def del_route(pwd):
    myprint("WARNING", "*****Trying to Un-install all 3 magma Virtual Machines*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo route del -net 192.168.60.0 netmask 255.255.255.0 dev br0' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo route del -net 192.168.129.0 netmask 255.255.255.0 dev br1' /dev/null" % pwd
    execute_cmd(cmd)

# Deleting iptables rules


def del_iptables(pwd):
    myprint("BLUE", "*****Trying to delete iptables rules added as part of installation*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo iptables -D FORWARD -s 192.168.0.0/16 -j ACCEPT' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo iptables -D FORWARD -d 192.168.0.0/16 -j ACCEPT' /dev/null" % pwd
    execute_cmd(cmd)

# Deleting 3 VMs magmatraffic, magmatest, magmadev


def del_vms():
    myprint("BLUE", "*****Deleting Alertmanger configurations*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/endpoint.yml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/service.yml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/service_monitor.yml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/alert_rules.yml")
    myprint("BLUE", "*****Revert the changes like remove magmadev VM IP from endpoint.yml, service.yml files*****")
    MAGMA_DEV_VM_IP = get_magmadev_vm_ip()
    os.chdir(TEMPLATES_DIR)
    for line in fileinput.input("endpoint.yml", inplace=True):
         if "ip" in line:
               print(line.replace(MAGMA_DEV_VM_IP, "YOUR_MAGMA_DEV_VM_IP"))
         else:
              print(line)
    for line in fileinput.input("service.yml", inplace=True):
        if "externalName:" in line:
               print(line.replace(MAGMA_DEV_VM_IP, "YOUR_MAGMA_DEV_VM_IP"))
        else:
              print(line)
    os.chdir(CWD)
    myprint("BLUE", "*****Deleting 3 VMs magmatraffic, magmatest, magmadev*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_traffic.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_test.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_dev.yaml")

# Deleting DataVolumes which are created for upload the Debian Image


def del_dvs(pwd):
    myprint("BLUE", "*****Deleting DataVolumes which are created for upload the Debian Image*****")
    execute_cmd("kubectl delete dv magma-traffic -n kubevirt")
    execute_cmd("kubectl delete dv magma-test -n kubevirt")
    execute_cmd("kubectl delete dv magma-dev -n kubevirt")
    time.sleep(10)
    myprint("BLUE", "*****Deleting PersistantVolumes [PVs] which are created for upload the Debian Image*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_dev_pv.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_test_pv.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/magma_traffic_pv.yaml")
    myprint("BLUE", "*****Deleting disk.img and tmpimage under /mnt path*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_dev/' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_dev_scratch/' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_test/' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_test_scratch/' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_traffic/' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm -rf /mnt/magma_traffic_scratch/' /dev/null" % pwd
    execute_cmd(cmd)

# Deleting network-attachment-definitions


def del_network_attachment_definition():
    myprint("BLUE", "*****Deleting Network-attachment-definitions*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/net_attach_def.yml")

# Removing ssh public key


def remove_ssh_key():
    myprint("BLUE", "*****Removing the id_rsa ssh-key [ssh public key]*****")
    execute_cmd("rm ~/.ssh/id_rsa.pub")
    execute_cmd("rm ~/.ssh/id_rsa")

# Delete Brdiges created to communicate with VMs


def del_bridges(pwd):
    myprint("BLUE", "*****Deleting Bridges created to communicate with VMs*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo ifconfig br0 down' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo ifconfig br1 down' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo brctl delbr br0' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo brctl delbr br1' /dev/null" % pwd
    execute_cmd(cmd)

# Deleting virtctl


def del_virtctl(pwd):
    myprint("BLUE", "*****Deleting virtctl*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo rm /usr/local/bin/virtctl' /dev/null" % pwd
    execute_cmd(cmd)

# Deleting CDI


def del_cdi():
    myprint("BLUE", "*****Deleting EFK [ElasticSearch Fluentd Kibana] configurations*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/kubevirt_efkchart.yaml -n kubevirt")
    myprint("BLUE", "*****Deleting Containerized Data Import [CDI]*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/cdi-cr.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/cdi-operator.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/storage-setup.yml")

# Deleting kubevirt


def del_kubevirt():
    myprint("BLUE", "*****Deleting kubevirt*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/kubevirt-cr.yaml")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/kubevirt-operator.yaml")

# Un-installing multus-cni plugin


def un_install_multus():
    myprint("BLUE", "*****Un-installing multus-cni plugin*****")
    os.chdir(MULTUS_DIR)
    execute_cmd("cat ./images/multus-daemonset.yml | kubectl delete -f -")
    os.chdir(TEMPLATES_DIR)
    if os.path.exists(MULTUS_DIR):
        shutil.rmtree(MULTUS_DIR)
    os.chdir(CWD)

# Deleting Temperay files and directories created as part of VM installation


def del_temp_files():
    myprint("BLUE", "*****Deleting Temperary files and directories created as part of VM installation*****")
    del_files(INFRA_SOFTWARE_VER)
    del_files(K8S_GET_OBJ)
    if os.path.exists(ORC8R_VM_DIR):
        shutil.rmtree(ORC8R_VM_DIR)

# Uninstalling VMs


def un_install_vm(pwd):
    del_route(pwd)
    del_iptables(pwd)
    del_vms()
    del_dvs(pwd)
    del_network_attachment_definition()
    remove_ssh_key()
    del_bridges(pwd)
    del_virtctl(pwd)
    del_cdi()
    del_kubevirt()
    un_install_multus()
    del_temp_files()

# Checking for pods and deployment status whether they are Running or not


def check_status(obj, namespace):
    print("check_satus", obj, namespace)
    if os.path.exists(ORC8R_VM_DIR):
        shutil.rmtree(ORC8R_VM_DIR)
    os.mkdir(ORC8R_VM_DIR)
    if obj == "pod":
        cmd = "kubectl get pods -n " + namespace + " | awk " + "'{{if ($3 ~ " + '!"Running"' + " || $3 ~ " + '!"STATUS"' + ") print $1,$3};}' > " + K8S_GET_OBJ
    elif obj == "deployment":
        cmd = "kubectl get deployment -n " + namespace + " | awk " + "'{{if ($2 ~ " + '!"1"' + " || $2 ~ " + '!"READY"' + ") print $1,$2};}' > " + K8S_GET_OBJ
    execute_cmd(cmd)
    if os.stat(K8S_GET_OBJ) == 0:
        return
    with open(K8S_GET_OBJ) as fop:
        while True:
            line = fop.readline()
            if not line:
                break
            myprint("WARNING", obj + "is not yet Running, please wait for a while")
            time.sleep(5)
            check_status(obj, namespace)

# thread1 : Getting the status of k8s objects like deployment and updating the k8s_obj_dict dictionary


def get_status(lock):
    while True:
        if os.path.exists(K8S_GET_DEP):
            if os.stat(K8S_GET_DEP).st_size == 0:
                break
        for values in k8s_obj_dict.values():
            # Get the deployment which are not in Running state
            cmd = "kubectl get deployment -n kubevirt | awk " + "'{{if ($2 ~ " + '!"1"' + " || $2 ~ " + '!"READY"' + ") print $1,$2};}' > " + K8S_GET_DEP
            execute_cmd(cmd)
            with open(K8S_GET_DEP) as fop1:
                while True:
                    k8s_obj_file1_line = fop1.readline()
                    if not k8s_obj_file1_line:
                        break
                    k8s_obj_name_list1 = k8s_obj_file1_line.split(' ')
                    for key in k8s_obj_dict.keys():
                        # Checking whether any key matches with deployment which are not in Running state
                        if re.search(k8s_obj_name_list1[0], key):
                            myprint("WARNING", "Few k8s Objects not Running YET!! Be patient, Please wait for a while")
                            # Get the latest status of all the deployments
                            cmd = "kubectl get deployment -n kubevirt | awk " + "'{{if (NR != 1) print $1,$2};}' > " + K8S_GET_SVC
                            execute_cmd(cmd)
                            with open(K8S_GET_SVC) as fop2:
                                while True:
                                    k8s_obj_file2_line = fop2.readline()
                                    if not k8s_obj_file2_line:
                                        break
                                    k8s_obj_name_list2 = k8s_obj_file2_line.split(' ')
                                    # Update the latest status of deployment into the k8s_obj_dict dictionary
                                    if re.search(k8s_obj_name_list1[0], k8s_obj_name_list2[0]):
                                        lock.acquire()
                                        k8s_obj_dict[key][0] = k8s_obj_name_list2[1]
                                        lock.release()

# thread2 : Getting the ports from running services and printing URL


def get_ports(lock):
    # Get the hostip into host_ip local variable
    host_ip = socket.gethostbyname(socket.gethostname())
    for key, values in k8s_obj_dict.items():
        if values[1] == 0:
            if len(values) > 2:
                port = values[2]
                cmd = "http://" + host_ip + ":" + port
                print("URL for :%s -->> %s" % (key, cmd))
                webbrowser.open(cmd, new=2)
                lock.acquire()
                values[1] = 1
                lock.release()

# Configure alert manager to get alerts from magmadev VM where AGW was Running


def configure_alert_manager():
    myprint("BLUE", "*****Get the magmadev VM IP and update in service.yml, endpoint.yml to get the alerts from magmadev VM*****")
    MAGMA_DEV_VM_IP = get_magmadev_vm_ip()
    os.chdir(TEMPLATES_DIR)
    for line in fileinput.input("endpoint.yml", inplace=True):
         if "ip" in line:
               print(line.replace("YOUR_MAGMA_DEV_VM_IP", MAGMA_DEV_VM_IP))
         else:
              print(line)
    for line in fileinput.input("service.yml", inplace=True):
        if "externalName:" in line:
               print(line.replace("YOUR_MAGMA_DEV_VM_IP", MAGMA_DEV_VM_IP))
        else:
              print(line)
    os.chdir(CWD)
    myprint("BLUE", "*****Applying the yaml files required to get the alerts from magmadev VM*****")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/endpoint.yml")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/service.yml")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/service_monitor.yml")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/alert_rules.yml")

# From the k8s services updating k8s_obj_dict dictionary and creating get_status, get_ports threads


def start_to_run():
    cmd = "kubectl get services -n kubevirt | awk " + "'{{if ($5 ~ " + '"TCP"' + " || $5 ~ " + '"UDP"' + ") print $1, $5};}' > " + K8S_GET_SVC
    execute_cmd(cmd)
    # Initializing the k8s_obj_dict with default values list[0, 0] for each key:k8s_obj_name
    with open(K8S_GET_SVC) as fop:
        while True:
            k8s_obj_file_line = fop.readline()
            if not k8s_obj_file_line:
                break
            k8s_obj_name_list = k8s_obj_file_line.split(' ')
            k8s_obj_dict[k8s_obj_name_list[0]] = [0, 0]
            # Updating the k8s_obj_dict with ports as values for each key:k8s_obj_name
            ports_list = k8s_obj_name_list[1].split('/')
            if len(ports_list[0].split(':')) > 1:
                for key in k8s_obj_dict.keys():
                    if re.search(k8s_obj_name_list[0], key):
                        k8s_obj_dict.setdefault(key, []).append(ports_list[0].split(':')[1])

    t1 = threading.Thread(target=get_status, args=(lock,))
    t2 = threading.Thread(target=get_ports, args=(lock,))
    t1.start()
    t2.start()
    t1.join()
    t2.join()

# Applying all the yaml files to create all k8s objects


def run_services():
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Installing Orc8r monitoring stack")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    execute_cmd("helm repo add prometheus-community https://prometheus-community.github.io/helm-charts")
    execute_cmd("helm repo add stable https://charts.helm.sh/stable")
    execute_cmd("helm repo update")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_alertmanagers.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_podmonitors.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_prometheuses.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_prometheusrules.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_servicemonitors.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_thanosrulers.yaml -n kubevirt")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/kubevirt_efkchart.yaml -n kubevirt")
    time.sleep(3)
    execute_cmd("helm install prometheus stable/prometheus-operator --namespace kubevirt")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-prometheus-oper-alertmanager -n kubevirt")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-grafana -n kubevirt")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-prometheus-oper-prometheus -n kubevirt")
    configure_alert_manager()
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Orc8r monitoring stack installed successfully")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("HDR", "-------------------------------------------------")
    myprint("WARNING", "        Printing URL's for Dashboards")
    myprint("HDR", "-------------------------------------------------")
    start_to_run()

# Install multus plugin which will be used for creating multiple interfaces in VM in addition to the default interfaces


def install_multus_plugin():
    myprint("BLUE", "*****Installing multus plugin which is used for creating multiple interfaces in VM in addition to the default interfaces*****")
    os.chdir(TEMPLATES_DIR)
    execute_cmd("git clone https://github.com/intel/multus-cni.git")
    os.chdir(MULTUS_DIR)
    execute_cmd("cat ./images/multus-daemonset.yml | kubectl apply -f -")
    os.chdir(CWD)

# Install kubevirt which allows to run virtual machines alongside your containers on a k8s platform


def install_kubevirt():
    myprint("BLUE", '*****Installing KubeVirt which allows to run virtual machines along with containers in k8s platform*****')
    execute_cmd("kubectl apply -f $PWD/../helm/templates/kubevirt-operator.yaml")
    check_status("pod", "kubevirt")
    execute_cmd("kubectl create configmap kubevirt-config -n kubevirt --from-literal debug-useEmulation=true")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/kubevirt-cr.yaml")
    check_status("pod", "kubevirt")
    myprint("BLUE", "*****Wait until all KubeVirt components is up*****")
    execute_cmd("kubectl -n kubevirt wait kv kubevirt --for condition=Available")

# Install Containerized Data Importer [CDI] used to import VM images to crate and control PVC


def install_cdi():
    myprint("BLUE", "*****Installing COntainerized Data Importer[CDI] used to import VM images to create PVC*****")
    execute_cmd("kubectl create -f $PWD/../helm/templates/storage-setup.yml")
    execute_cmd("kubectl create -f $PWD/../helm/templates/cdi-operator.yaml")
    execute_cmd("kubectl create -f $PWD/../helm/templates/cdi-cr.yaml")
    check_status("pod", "cdi")

# Install virtctl which is used to create DV,PVC to upload disk.img also used to connect and control VM via CLI


def install_virtctl(pwd):
    myprint("BLUE", '*****Installing virtctl which is used to create DV,PVC to upload disk.img*****')
    os.chdir(TEMPLATES_DIR)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo install virtctl /usr/local/bin' /dev/null" % pwd
    execute_cmd(cmd)
    os.chdir(CWD)

# Create Bridges which are required to communicate between Host to VM and VM to VM


def create_bridges(pwd):
    myprint("BLUE", "*****Creating Bridges required to communicate between Host to VM and VM to VM*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo brctl addbr br0' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo brctl addbr br1' /dev/null" % pwd
    execute_cmd(cmd)

# Creating NetworkAttachmentDefinition to configure Network Attachment with a L2 Bridge and Vlan


def create_network_attachment_definition():
    myprint("BLUE", "*****Creating NetworkAttachmentDefinition to configure Network Attachment with a L2 Bridge*****")
    execute_cmd("kubectl create -f $PWD/../helm/templates/net_attach_def.yml")

# Generate ssh-key and inject to debian qcow2 image to make use of passwordless authentication via root user


def generate_ssh_public_key(pwd):
    os.chdir(TEMPLATES_DIR)
    if not os.path.exists(DEBIAN_QCOW2_FILE):
        myprint("WARNING", "*****debian-9-openstack-amd64.qcow2 image is not present under magma/cn/deploy/helm/templates/ directory script will download it, Please be patient!! it may take some time based on your bandwidth!!*****")
        execute_cmd("wget http://cdimage.debian.org/cdimage/openstack/current-9/debian-9-openstack-amd64.qcow2")
    else:
        myprint("BLUE", "*****debian-9-openstack-amd64.qcow2 image is already present under magma/cn/deploy/helm/templates/ directory so skipping download!!*****")
    myprint("BLUE", "*****Generating password-less SSH key and inject to debian qcow2 image*****")
    execute_cmd('ssh-keygen -f ~/.ssh/id_rsa -q -N "" 0>&-')
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo virt-sysprep -a debian-9-openstack-amd64.qcow2 --ssh-inject root:file:$HOME/.ssh/id_rsa.pub' /dev/null" % pwd
    execute_cmd(cmd)
    os.chdir(CWD)
    execute_cmd("kubectl -n kubevirt wait kv kubevirt --for condition=Available")
    time.sleep(10)

# Creating DataVolumes for magmadev, magmatest, magmatraffic VMs, These DataVolumes will mount corresponding PVC


def create_datavolume(pwd):
    myprint("BLUE", "*****Creating DataVolumes to mount debian qcow2 image *****")
    # Get the cdi_uplodproxy service IP address which is used to frame URL to image-upload
    cmd = "kubectl get svc -n cdi | grep 'cdi-uploadproxy' | awk '{print $3}'"
    data = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    stdout, stderr = data.communicate()
    cdi_uplaod_proxy_ip_add = stdout.strip().decode('utf-8')
    # Create directories under /mnt to store disk.img from the mounted PVC
    myprint("BLUE", "*****Create directories under /mnt to store disk.img under /mnt from the mounted PVC *****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_dev' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_dev_scratch' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_test' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_test_scratch' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_traffic' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo mkdir -p /mnt/magma_traffic_scratch' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo chmod 777 /mnt/*' /dev/null" % pwd
    execute_cmd(cmd)
    # Create PVs which are going to claim by PVCs
    myprint("BLUE", "*****Create PVs which are going to Claim by PVCs*****")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/magma_dev_pv.yaml")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/magma_test_pv.yaml")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/magma_traffic_pv.yaml")
    # Create DataVolume[dv] which will mount the debian qcow2 disk.img to corresponding mounted path under /mnt
    myprint("BLUE", "*****Create DataVolume[dv] which will mount the debian qcow2 disk.img to directory under /mnt*****")
    try:
        cmd = "virtctl image-upload dv magma-dev --namespace kubevirt --pvc-size=50Gi --image-path $PWD/../helm/templates/debian-9-openstack-amd64.qcow2 --uploadproxy-url=https://%s:443 --insecure" % cdi_uplaod_proxy_ip_add
        execute_cmd(cmd)
        cmd = "virtctl image-upload dv magma-test --namespace kubevirt --pvc-size=50Gi --image-path $PWD/../helm/templates/debian-9-openstack-amd64.qcow2 --uploadproxy-url=https://%s:443 --insecure" % cdi_uplaod_proxy_ip_add
        execute_cmd(cmd)
        cmd = "virtctl image-upload dv magma-traffic --namespace kubevirt --pvc-size=50Gi --image-path $PWD/../helm/templates/debian-9-openstack-amd64.qcow2 --uploadproxy-url=https://%s:443 --insecure" % cdi_uplaod_proxy_ip_add
        execute_cmd(cmd)
    except NotInstalled:
        print("Image upload not completed")
        myprint("FAIL", "Image upload not completed")

# Creating 3 VMs magmadev, magmatest, magmatraffic


def create_vm():
    myprint("BLUE", "*****Creating 3 VMs magmadev, magmatest, magmatraffic*****")
    execute_cmd("kubectl create -f $PWD/../helm/templates/magma_dev.yaml")
    execute_cmd("kubectl create -f $PWD/../helm/templates/magma_test.yaml")
    execute_cmd("kubectl create -f $PWD/../helm/templates/magma_traffic.yaml")
    myprint("BLUE", "*****Wait for some time to VM to wake up to Running state*****")
    time.sleep(10)

# Adding route information of Bridge


def add_route_info(pwd):
    myprint("BLUE", "*****Add route information of bridge*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo ifconfig br0 up' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo ifconfig br1 up' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo route add -net 192.168.60.0 netmask 255.255.255.0 dev br0' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo route add -net 192.168.129.0 netmask 255.255.255.0 dev br1' /dev/null" % pwd
    execute_cmd(cmd)

# Updating iptables to forward VM traffic


def add_iptables_rule(pwd):
    myprint("BLUE", "*****Update iptables to forward VM traffic*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo iptables -A FORWARD -s 192.168.0.0/16 -j ACCEPT' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo iptables -A FORWARD -d 192.168.0.0/16 -j ACCEPT' /dev/null" % pwd
    execute_cmd(cmd)

# Create magmadev, magmatest, magmatraffic [3 VMs] by step by step


def install_vm(pwd):
    install_multus_plugin()
    install_kubevirt()
    install_cdi()
    install_virtctl(pwd)
    create_bridges(pwd)
    create_network_attachment_definition()
    generate_ssh_public_key(pwd)
    create_datavolume(pwd)
    create_vm()
    add_route_info(pwd)
    add_iptables_rule(pwd)

# Displays the Usage of the script


def get_help(color):
    myprint(color, './MVC2_5G_Orc8r_deployment_script.py -p <sudo-password> -i')
    myprint(color, './MVC2_5G_Orc8r_deployment_script.py -p <sudo-password> -u')
    myprint(color, '    (OR)   ')
    myprint(color, './MVC2_5G_Orc8r_deployment_script.py --password <sudo-password> --install')
    myprint(color, './MVC2_5G_Orc8r_deployment_script.py --password <sudo-password> --uninstall')


def main(argv):
    password = ''
    try:
        opts, args = getopt.getopt(argv, "hiup:", ["help", "install", "uninstall", "password="])
    except getopt.GetoptError:
        get_help("FAIL")

    for opt, arg in opts:
        if (re.match("-h", opt) or re.match("--help", opt)):
            get_help("BLUE")
        elif (opt == "-p" or opt == "--password"):
            password = arg
        elif (opt == "-i" or opt == "--install"):
            myprint("HDR", "-------------------------------------------------")
            myprint("GREEN", "           Checking Pre-requisites: ")
            myprint("HDR", "-------------------------------------------------")
            check_pre_requisite()
            install_vm(password)
            run_services()
            myprint("HDR", "-------------------------------------------------")
            myprint("WARNING", "    URL's for Dashboards printed successfully")
            myprint("HDR", "-------------------------------------------------")
        elif (opt == "-u" or opt == "--uninstall"):
            un_install_vm(password)
            un_install(password)


if __name__ == "__main__":
    main(sys.argv[1:])
