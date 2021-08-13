/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import type {wifi_gateway} from '@fbcnms/magma-api';

export const RAW_GATEWAY: wifi_gateway = {
  description: '',
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 15,
    checkin_timeout: 12,
  },
  tier: 'default',
  wifi: {
    additional_props: {example_prop: 'lol', expected_is_gateway: '1'},
    client_channel: '1',
    info: 'binney lab, top shelf back wall',
    is_production: true,
    latitude: 83,
    longitude: -70,
    mesh_id: 'shared_d',
    mesh_rssi_threshold: -80,
  },
  id: 'shared_d_id_5ce28cf1aeb6',
  name: '5ce28cf1aeb6',
  device: {
    hardware_id: 'faceb00c-face-b00c-face-5ce28cf1aeb6',
    key: {key_type: 'ECHO'},
  },
  status: {
    checkin_time: 1561156214384,
    hardware_id: 'faceb00c-face-b00c-face-5ce28cf1aeb6',
    kernel_version: '4.14.104',
    machine_info: {
      cpu_info: {
        architecture: 'armv7l',
        core_count: 2,
        model_name: 'Krait',
        threads_per_core: 1,
      },
      network_info: {
        network_interfaces: [],
        routing_table: [],
      },
    },
    meta: {
      version:
        'Facebook Wi-Fi soma-image-1.0 (nbg6817) Release 52461bc4d24d+dirty (yerv@devvm354 Wed Jun 19 18:30:23 UTC 2019) (fbpkg:none) (cfg:none)',
      is_gateway: 'true',
      openr_inet_monitor: 'success',
      default_route: 'eth0,10.1.0.1',
      default_route_v6:
        'eth0,fe80::7e25:8602:ccfb:3819;eth0,fe80::c242:d002:cc83:52d8;eth0,fe80::200:5eff:fe00:201',
    },
    platform_info: {
      config_info: {mconfig_created_at: 1561156211},
      kernel_version: '4.14.104',
      packages: [{name: 'magma', version: '0.0.0'}],
      vpn_ip: 'N/A',
    },
    system_status: {
      cpu_idle: 179422550,
      cpu_system: 14347810,
      cpu_user: 14490730,
      disk_partitions: [
        {
          device: '/dev/loop0',
          mount_point: '/boot/root/mnt/squashroot',
          total: 79167488,
          used: 79167488,
        },
      ],
      mem_available: 170078208,
      mem_free: 36532224,
      mem_total: 487931904,
      mem_used: 283856896,
      time: 1561156213,
      uptime_secs: 109884,
    },
    version: '0.0.0',
    vpn_ip: 'N/A',
  },
};
