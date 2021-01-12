# Here is script.sh file.

For installing kubernetes(Version 1.9.1), Docker and Helm(Version 3.3.1) I have created script.sh file.

Note :
    1) Before running script.sh file we need to change below line according to our VM ip.
        sudo kubeadm init --pod-network-cidr=10.244.20.0/24 --apiserver-advertise-address=<VM_ip>
    2) And also make sure our hosts file entry as same in hostname accordingly in script.sh file.
        sudo hostnamectl set-hostname <kubecluster>