#  Magma LTE Gateway eBPF
This folder contains the **eBPF subsystem** for the Magma LTE Gateway.  
It includes kernel programs, userspace helpers, Python integration, and systemd services to manage traffic control, GTP sessions, and telemetry.

## Folder Structure

    ebpf/
    ├── bpf/                
    ├── config/          
    ├── python/             
    ├── scripts/         
    ├── src/          
    ├── systemd/         
    ├── tests/        
    └── WORKSPACE
   
   ## Deploying eBPF Programs Manually    
   
   1.  **Generate `vmlinux.h`**    

    cd scripts
    ./gen_vmlinux.sh
    
   2. Build eBPF programs

    cd ../bpf
    make all
3. Load the TC program
`sudo ./scripts/load_tc.sh'
4.  **Verify BPF maps
    'sudo ./scripts/debug_maps.sh'

5.  **Unload TC program (if needed)**
    

`sudo ./scripts/unload_tc.sh`

----
## Deploying via Ansible

The `ansible/` folder automates installation, configuration, and service management for eBPF.

### 1. Update the inventory

Edit `ansible/inventory/hosts.ini` with your gateway host(s):

`[gateway] magma-gateway ansible_host=192.168.60.10 ansible_user=ubuntu ansible_become=true` 

### 2. Run the playbook

`cd ansible
ansible-playbook -i inventory/hosts.ini playbooks/deploy_ebpf.yml` 

### 3. Playbook actions

-   Installs dependencies (`clang`, `bpftool`, kernel headers, etc.)
    
-   Generates `vmlinux.h`
    
-   Builds eBPF programs
    
-   Copies binaries and configuration to `/var/opt/magma/ebpf`
    
-   Deploys systemd services:
    
    -   `magma-ebpf-loader.service`
        
    -   `magma-ebpf-manager.service`
        

### 4. Verify deployment

`sudo systemctl status magma-ebpf-loader.service

sudo systemctl status magma-ebpf-manager.service

sudo bpftool prog list`