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
# Build context for Docker files for AMF, SMF, UPF
BUILD_CONTEXT_PARENT_DIR = '/tmp/orc8r_docker'
BUILD_CONTEXT_AMF = os.path.join(BUILD_CONTEXT_PARENT_DIR, 'AMF')
BUILD_CONTEXT_SMF = os.path.join(BUILD_CONTEXT_PARENT_DIR, 'SMF')
BUILD_CONTEXT_UPF = os.path.join(BUILD_CONTEXT_PARENT_DIR, 'UPF')
# Get the Current Working Directory
CWD = os.getcwd()
# Path for MAGMA dockerfile
MAGMA_DOCKER = os.path.join(CWD, '../docker')
# Path for Orc8r temperary files
ORC8R_TEMP_DIR = '/tmp/orc8r_temp'
INFRA_SOFTWARE_VER = os.path.join(ORC8R_TEMP_DIR, 'infra_software_version.txt')
K8S_GET_DEP = os.path.join(ORC8R_TEMP_DIR, 'k8s_get_deployment.txt')
K8S_GET_SVC = os.path.join(ORC8R_TEMP_DIR, 'k8s_get_service.txt')
K8S_DEL_OBJ = os.path.join(ORC8R_TEMP_DIR, 'k8s_del_objects.txt')


class Error(Exception):
    """Base class for other exceptions"""
    pass


class NotInstalled(Error):
    """Raised when Installation not done"""
    pass


def Code(type1):
    switcher = {
        'WARNING': 93,
        'FAIL': 91,
        'GREEN': 92,
        'BLUE': 94,
        'ULINE': 4,
        'BLD': 1,
        'HDR': 95,
    }
    return switcher.get(type1)

# Print messages with colours on console


def myprint(type1, msg):
    code = Code(type1)
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
    myprint("BLUE", "        installing Orc8r monitoring stack")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")

# Creating docker images for stubbed AMF, SMF, UPF 5G core components


def create_docker_images(pwd):
    myprint("BLUE", "Creating docker images for stubbed AMF, SMF, UPF core components")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker pull nginx' /dev/null" % pwd
    execute_cmd(cmd)
    if os.path.exists(BUILD_CONTEXT_PARENT_DIR):
        shutil.rmtree(BUILD_CONTEXT_PARENT_DIR)
    os.mkdir(BUILD_CONTEXT_PARENT_DIR)
    if os.path.exists(BUILD_CONTEXT_AMF):
        shutil.rmtree(BUILD_CONTEXT_AMF)
    os.mkdir(BUILD_CONTEXT_AMF)
    shutil.copyfile(os.path.join(MAGMA_DOCKER, 'amf_dockerfile'), os.path.join(BUILD_CONTEXT_AMF, 'Dockerfile'))
    os.chdir(BUILD_CONTEXT_AMF)
    cmd = 'echo " This is AMF !!" > index.html'
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker build -t amf-container:latest %s' /dev/null" % (pwd, BUILD_CONTEXT_AMF)
    execute_cmd(cmd)
    if os.path.exists(BUILD_CONTEXT_SMF):
        shutil.rmtree(BUILD_CONTEXT_SMF)
    os.mkdir(BUILD_CONTEXT_SMF)
    shutil.copyfile(os.path.join(MAGMA_DOCKER, 'smf_dockerfile'), os.path.join(BUILD_CONTEXT_SMF, 'Dockerfile'))
    os.chdir(BUILD_CONTEXT_SMF)
    cmd = 'echo " This is SMF !!" > index.html'
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker build -t smf-container:latest %s' /dev/null" % (pwd, BUILD_CONTEXT_SMF)
    execute_cmd(cmd)
    if os.path.exists(BUILD_CONTEXT_UPF):
        shutil.rmtree(BUILD_CONTEXT_UPF)
    os.mkdir(BUILD_CONTEXT_UPF)
    shutil.copyfile(os.path.join(MAGMA_DOCKER, 'upf_dockerfile'), os.path.join(BUILD_CONTEXT_UPF, 'Dockerfile'))
    os.chdir(BUILD_CONTEXT_UPF)
    cmd = 'echo " This is UPF !!" > index.html'
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker build -t upf-container:latest %s' /dev/null" % (pwd, BUILD_CONTEXT_UPF)
    execute_cmd(cmd)
    os.chdir(CWD)

# thread1 : Getting the status of k8s objects like deployment and updating the k8s_obj_dict dictionary


def get_status(lock):
    while True:
        if os.path.exists(K8S_GET_DEP):
            if os.stat(K8S_GET_DEP).st_size == 0:
                break
        for _ in k8s_obj_dict.values():
            # Get the deployment which are not in Running state
            cmd = "kubectl get deployment -n default | awk " + "'{{if ($2 ~ " + '!"1"' + " || $2 ~ " + '!"READY"' + ") print $1,$2};}' > " + K8S_GET_DEP
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
                            cmd = "kubectl get deployment -n default | awk " + "'{{if (NR != 1) print $1,$2};}' > " + K8S_GET_SVC
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

# From the k8s services updating k8s_obj_dict dictionary and creating get_status, get_ports threads


def start_to_run():
    cmd = "kubectl get services -n default | awk " + "'{{if ($5 ~ " + '"TCP"' + " || $5 ~ " + '"UDP"' + ") print $1, $5};}' > " + K8S_GET_SVC
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
    execute_cmd("helm repo add prometheus-community https://prometheus-community.github.io/helm-charts")
    execute_cmd("helm repo add stable https://charts.helm.sh/stable")
    execute_cmd("helm repo update")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_alertmanagers.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_podmonitors.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_prometheuses.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_prometheusrules.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_servicemonitors.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/monitoring.coreos.com_thanosrulers.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/amfchart.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/smfchart.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/upfchart.yaml -n default")
    execute_cmd("kubectl apply -f $PWD/../helm/templates/efkchart.yaml -n default")
    time.sleep(3)
    execute_cmd("helm install prometheus-default stable/prometheus-operator --namespace default")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-default-prometh-alertmanager -n default")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-default-grafana -n default")
    myprint("FAIL", "change type(key) value from 'ClusterIP' to 'NodePort' and save it")
    time.sleep(3)
    execute_cmd("kubectl edit service/prometheus-default-prometh-prometheus -n default")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Orc8r monitoring stack installed successfully")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("HDR", "-------------------------------------------------")
    myprint("WARNING", "         Installing 5G-core apps")
    myprint("HDR", "-------------------------------------------------")
    start_to_run()

# Deleting the k8s objects like pods, services, deployemnts..etc


def del_objects(file, type1):
    k8s_obj = "kubectl get all -n default| grep %s > %s" % (type1, file)
    execute_cmd(k8s_obj)
    with open(file) as fop:
        while True:
            del_obj = fop.readline()
            if not del_obj:
                break
            k8s_obj = "kubectl delete %s -n default" % del_obj.split(' ')[0]
            execute_cmd(k8s_obj)

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
    execute_cmd("helm uninstall prometheus-default stable/prometheus-operator --namespace default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_thanosrulers.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_servicemonitors.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_prometheusrules.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_prometheuses.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_podmonitors.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/monitoring.coreos.com_alertmanagers.yaml -n default")
    myprint("BLUE", "*****Trying to Un-install AMF/SMF/UPF containers*****")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/amfchart.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/smfchart.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/upfchart.yaml -n default")
    execute_cmd("kubectl delete -f $PWD/../helm/templates/efkchart.yaml -n default")
    myprint("BLUE", "*****Trying to remove AMF/SMF/UPF containers*****")
    cmd = '''{ sleep 0.1; echo '%s'; } | script -q -c "sudo docker ps -aqf 'name=upf_proto|smf_proto|amf_proto' | sudo xargs docker container rm -f" /dev/null''' % pwd
    execute_cmd(cmd)
    myprint("BLUE", "*****Trying to remove AMF/SMF/UPF docker images*****")
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker rmi upf-container' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker rmi smf-container' /dev/null" % pwd
    execute_cmd(cmd)
    cmd = "{ sleep 0.1; echo '%s'; } | script -q -c 'sudo docker rmi amf-container' /dev/null" % pwd
    execute_cmd(cmd)
    myprint("BLUE", "*****Trying to Cleanup the temporay files & Directories created as part of installation*****")
    del_objects(K8S_DEL_OBJ, "service")
    del_objects(K8S_DEL_OBJ, "deployment")
    del_objects(K8S_DEL_OBJ, "daemonset")
    del_files(INFRA_SOFTWARE_VER)
    del_files(K8S_GET_DEP)
    del_files(K8S_GET_SVC)
    del_files(K8S_DEL_OBJ)
    del_files(os.path.join(BUILD_CONTEXT_AMF, 'Dockerfile'))
    del_files(os.path.join(BUILD_CONTEXT_AMF, 'index.html'))
    del_files(os.path.join(BUILD_CONTEXT_SMF, 'Dockerfile'))
    del_files(os.path.join(BUILD_CONTEXT_SMF, 'index.html'))
    del_files(os.path.join(BUILD_CONTEXT_UPF, 'Dockerfile'))
    del_files(os.path.join(BUILD_CONTEXT_UPF, 'index.html'))
    if os.path.exists(BUILD_CONTEXT_UPF):
        shutil.rmtree(BUILD_CONTEXT_UPF)
    if os.path.exists(BUILD_CONTEXT_SMF):
        shutil.rmtree(BUILD_CONTEXT_SMF)
    if os.path.exists(BUILD_CONTEXT_AMF):
        shutil.rmtree(BUILD_CONTEXT_AMF)
    if os.path.exists(BUILD_CONTEXT_PARENT_DIR):
        shutil.rmtree(BUILD_CONTEXT_PARENT_DIR)
    if os.path.exists(ORC8R_TEMP_DIR):
        shutil.rmtree(ORC8R_TEMP_DIR)
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")
    myprint("BLUE", "  Orc8r monitoring stack Uninstalled successfully")
    myprint("GREEN", "+++++++++++++++++++++++++++++++++++++++++++++++")


def get_help(color):
    myprint(color, './mvc1_5g_orc8r_deployment_script.py -p <sudo-password> -i')
    myprint(color, './mvc1_5g_orc8r_deployment_script.py -p <sudo-password> -u')
    myprint(color, '    (OR)   ')
    myprint(color, './mvc1_5g_orc8r_deployment_script.py --password <sudo-password> --install')
    myprint(color, './mvc1_5g_orc8r_deployment_script.py --password <sudo-password> --uninstall')


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
            create_docker_images(password)
            run_services()
            myprint("HDR", "-------------------------------------------------")
            myprint("WARNING", "    5G-core apps are Running Successfully")
            myprint("HDR", "-------------------------------------------------")
        elif (opt == "-u" or opt == "--uninstall"):
            un_install(password)


if __name__ == "__main__":
    main(sys.argv[1:])
