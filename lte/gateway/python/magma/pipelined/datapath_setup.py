"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import subprocess
import os
import logging

irq_utility = '/usr/local/bin/set_irq_affinity'
ethtool_utility = '/usr/sbin/ethtool'

'''
Following function sets various tuning parameters related
interface queue.
1. RX queue size
2. TX queue size
3. queue CPU assignment
'''


def tune_datapath(config_dict):
    # TODO move this to mconfig
    if 'dp_irq' not in config_dict:
        logging.info("DP Tuning not enabled.")
        return

    if _check_util_failed(irq_utility):
        return
    if _check_util_failed(ethtool_utility):
        return

    tune_dp_irqs = config_dict['dp_irq']
    logging.info("set tuning params: %s", tune_dp_irqs)

    # stop irq-balance
    stop_irq_balance = ['service', 'irqbalance', 'stop']
    logging.debug("cmd: %s", stop_irq_balance)
    try:
        subprocess.check_call(stop_irq_balance)
    except subprocess.CalledProcessError as ex:
        logging.debug('%s failed with: %s', stop_irq_balance, ex)

    # set_irq_affinity -X 1-2 eth1
    s1_interface = config_dict['enodeb_iface']
    s1_cpu = tune_dp_irqs['S1_cpu']
    set_s1_cpu_command = [irq_utility, '-X', s1_cpu, s1_interface]
    logging.debug("cmd: %s", set_s1_cpu_command)
    try:
        subprocess.check_call(set_s1_cpu_command)
    except subprocess.CalledProcessError as ex:
        logging.debug('%s failed with: %s', set_s1_cpu_command, ex)

    sgi_interface = config_dict['nat_iface']
    sgi_cpu = tune_dp_irqs['SGi_cpu']
    set_sgi_cpu_command = [irq_utility, '-X', sgi_cpu, sgi_interface]
    logging.debug("cmd: %s", set_sgi_cpu_command)
    try:
        subprocess.check_call(set_sgi_cpu_command)
    except subprocess.CalledProcessError as ex:
        logging.debug('%s failed with: %s', set_sgi_cpu_command, ex)

    # ethtool -G eth1 rx 1024 tx 1024
    s1_queue_size = tune_dp_irqs['S1_queue_size']
    set_s1_queue_sz = [ethtool_utility, '-G', s1_interface,
                       'rx', str(s1_queue_size), 'tx', str(s1_queue_size)]
    logging.debug("cmd: %s", set_s1_queue_sz)
    try:
        subprocess.check_call(set_s1_queue_sz)
    except subprocess.CalledProcessError as ex:
        logging.debug('%s failed with: %s', set_s1_queue_sz, ex)

    sgi_queue_size = tune_dp_irqs['SGi_queue_size']
    set_sgi_queue_sz = [ethtool_utility, '-G', sgi_interface,
                        'rx', str(sgi_queue_size), 'tx', str(sgi_queue_size)]
    logging.debug("cmd: %s", set_sgi_queue_sz)
    try:
        subprocess.check_call(set_sgi_queue_sz)
    except subprocess.CalledProcessError as ex:
        logging.debug('%s failed with: %s', set_sgi_queue_sz, ex)


def _check_util_failed(path: str):
    if not os.path.isfile(path) or not os.access(path, os.X_OK):
        logging.info("missing %s: path: %s perm: %s", path,
                     os.path.isfile(path),
                     os.access(path, os.X_OK))
        return True
    return False
