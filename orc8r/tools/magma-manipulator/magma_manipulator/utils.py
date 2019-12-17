"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from io import StringIO
import json
import logging
import os

import paramiko

from magma_manipulator import exceptions


LOG = logging.getLogger(__name__)
logging.getLogger("paramiko").setLevel(logging.WARNING)

CLOUD_INIT_CHECK_CMD = 'cloud-init status'
CLOUD_INIT_DONE = 'done'
CLOUD_INIT_RUNNING = 'running'

GET_GW_UUID_CMD = 'cd /var/opt/magma/docker ; '\
                  'sudo docker-compose exec '\
                  '-T magmad /usr/local/bin/show_gateway_info.py'


def is_gw_reachable(gw_ip):
    response = os.system('ping -c 1 ' + gw_ip)
    return response == 0


def exec_ssh_command(server, username, rsa_private_key_path, command):
    client = None
    try:
        with open(rsa_private_key_path, 'r') as f:
            s = f.read()
        pkey = paramiko.RSAKey(file_obj=StringIO(s))

        client = paramiko.SSHClient()
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        LOG.debug('Connection to server {server} '
                  'to execute command "{cmd}"'.format(server=server,
                                                      cmd=command))
        client.connect(server, username=username, pkey=pkey)
        ssh_stdin, ssh_stdout, ssh_stderr = client.exec_command(command)

        return ssh_stdout.read().decode('ascii')
    except Exception as e:
        msg = 'Execution ssh command "{cmd}" on server {server}'\
              'returns {msg}'.format(cmd=command, server=server, msg=e)
        LOG.error(msg)
        raise exceptions.SshRemoteCommandException(msg)
    finally:
        if client:
            client.close()


def is_cloud_init_done(gw_ip, gw_username, rsa_private_key_path):
    LOG.info('Check cloud-init status on gatewat {gw_ip}'.format(gw_ip=gw_ip))
    result = exec_ssh_command(gw_ip,
                              gw_username,
                              rsa_private_key_path,
                              CLOUD_INIT_CHECK_CMD)
    LOG.info('Cloud-init status: {status} on gateway {gw_ip}'.format(
        status=result,
        gw_ip=gw_ip))
    if CLOUD_INIT_DONE in result:
        return True
    elif CLOUD_INIT_RUNNING in result:
        return False
    else:
        msg = 'Something goes wrong with cloud-init '\
              'on gateway: {error}'.format(error=result)
        LOG.error(msg)
        raise exceptions.CloudInitException(msg)


def get_gw_uuid_and_key(gw_ip, gw_username, rsa_private_key_path):
    ssh_output = exec_ssh_command(gw_ip,
                                  gw_username,
                                  rsa_private_key_path,
                                  GET_GW_UUID_CMD)
    gw_uuid = ssh_output.split('\n')[2]
    gw_key = ssh_output.split('\n')[6]
    return (gw_uuid, gw_key)


def save_gateway_config(gw_id, configs_dir, cfg):
    if not os.path.exists(configs_dir):
        os.makedirs(configs_dir)
        LOG.info('Create directory for gateways configs {dir}'.format(
            dir=configs_dir))
    cfg_name = str(gw_id) + '.json'
    cfg_path = os.path.join(configs_dir, cfg_name)
    with open(cfg_path, 'w', encoding='utf-8') as f:
        json.dump(cfg, f, ensure_ascii=False, indent=4)
    return cfg_path


def load_gateway_config(gw_id, config_path):
    with open(config_path, 'r', encoding='utf-8') as f:
        json_data = json.load(f)
    return json_data
