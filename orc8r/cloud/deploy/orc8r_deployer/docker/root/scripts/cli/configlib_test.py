#!/usr/bin/env python3

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

import os
import shutil
import subprocess
from tempfile import mkdtemp
from unittest import TestCase, main
from unittest.mock import patch

import configlib
import utils.awslib
import yaml
from utils.common import get_json


class TestConfigManager(TestCase):
    def setUp(self):
        self.root_dir = mkdtemp()
        self.constants = {
            'project_dir': self.root_dir,
            'config_dir': '%s/configs' % self.root_dir,
            'vars_definition': '%s/vars.yml' % self.root_dir,
            'components': ['infra', 'platform', 'service'],
            'auto_tf': self.root_dir + '/terraform.tfvars.json',
            'tf_dir': self.root_dir,
            'main_tf': self.root_dir + '/main.tf',
            'vars_tf': self.root_dir + '/vars.tf'
        }
        # add config directory
        os.makedirs(self.constants["config_dir"])
        test_vars = {
            "infra": {
                "aws_access_key_id": {
                    "Required": True,
                    "ConfigApps": ["awscli"]
                },
                "aws_secret_access_key": {
                    "Required": True,
                    "ConfigApps": ["awscli"]
                },
                "cluster_name": {
                    "Required": False,
                    "ConfigApps": ["tf"]
                },
                "secretsmanager_orc8r_secret": {
                    "Required": True,
                    "ConfigApps": ["tf"]
                }
            },
            "platform": {
                "deploy_elastic": {
                    "Required": False,
                    "ConfigApps": ["tf"]
                },
                "nms_db_password": {
                    "Required": True,
                    "ConfigApps": ["tf"]
                }
            },
            "service": {
                "lte_orc8r_chart_version": {
                    "Required": False,
                    "Default": "0.2.4",
                    "ConfigApps": ["tf"]
                },
            }
        }
        with open(self.constants['vars_definition'], 'w') as f:
            yaml.dump(test_vars, f)

        # write a simple jinja template file
        jinja_template = (
            '''
{% for k in cfg['infra'] %}
{{k}}=var.{{k}}{% endfor %}
''')
        with open(self.constants['tf_dir'] + '/main.tf.j2', 'w') as f:
            f.write(jinja_template)

        jinja_template = (
            '''
{% for k in cfg %}
variable "{{k}}" {}{% endfor %}
''')
        with open(self.constants['tf_dir'] + '/vars.tf.j2', 'w') as f:
            f.write(jinja_template)

    def tearDown(self):
        shutil.rmtree(self.root_dir, ignore_errors=True)

    @patch("configlib.get_input")
    @patch("utils.awslib.run_command")
    def test_configure_sanity(self, mock_run_command, mock_get_input):
        mock_vals = {
            'aws_access_key_id': 'foo',
            'aws_secret_access_key': 'bar',
            'secretsmanager_orc8r_secret': 'jar',
        }
        mock_get_input.side_effect = [
            mock_vals['aws_access_key_id'],
            mock_vals['aws_secret_access_key'],
            mock_vals['secretsmanager_orc8r_secret'],
        ]
        mock_run_command.return_value = subprocess.CompletedProcess(
            args=[], returncode=0)

        # verify if components tfvars json is created
        mgr = configlib.ConfigManager(self.constants)
        mgr.configure('infra')
        mgr.commit('infra')

        # verify if configs are set in infra tfvars json
        fn = "%s/infra.tfvars.json" % self.constants['config_dir']
        cfg = get_json(fn)
        self.assertEqual(len(cfg.keys()), 3)
        self.assertEqual(cfg['aws_access_key_id'], "foo")
        self.assertEqual(cfg['aws_secret_access_key'], "bar")
        self.assertEqual(cfg['secretsmanager_orc8r_secret'], "jar")

        # check if aws configs are set
        aws_config_cmd = ['aws', 'configure', 'set']
        mock_run_command.assert_any_call(
            aws_config_cmd + ['aws_access_key_id', 'foo'])
        mock_run_command.assert_any_call(
            aws_config_cmd + ['aws_secret_access_key', 'bar'])

        # verify that platform tfvars json file isn't present
        fn = "%s/platform.tfvars.json" % self.constants['config_dir']
        self.assertEqual(os.path.isfile(fn), False)

        # reset mocks
        mock_get_input.reset_mock()
        mock_run_command.reset_mock()

        mock_vals = {
            'nms_db_password': 'foo',
        }
        mock_get_input.side_effect = [
            mock_vals['nms_db_password'],
        ]

        # configure platform
        mgr.configure('platform')
        mgr.commit('platform')

        # verify that no aws call was invoked
        self.assertEqual(mock_run_command.call_count, 0)

        # check if we only invoked input for nms_db_password
        self.assertEqual(mock_get_input.call_count, 1)
        fn = "%s/platform.tfvars.json" % self.constants['config_dir']
        cfg = get_json(fn)
        self.assertEqual(len(cfg.keys()), 1)
        self.assertEqual(cfg['nms_db_password'], "foo")

        # verify that service tfvars json file isn't present
        fn = "%s/service.tfvars.json" % self.constants['config_dir']
        self.assertEqual(os.path.isfile(fn), False)

        # reset mocks
        mock_get_input.reset_mock()
        mock_run_command.reset_mock()

        # configure service
        mgr.configure('service')
        mgr.commit('service')

        # verify that no input or aws call was invoked
        self.assertEqual(mock_run_command.call_count, 0)
        self.assertEqual(mock_get_input.call_count, 0)

        fn = "%s/service.tfvars.json" % self.constants['config_dir']
        cfg = get_json(fn)

        # verify that default value was set
        self.assertEqual(len(cfg.keys()), 1)
        self.assertEqual(cfg['lte_orc8r_chart_version'], "0.2.4")

        # finally verify if all configs required by tf is present
        cfg = get_json(self.constants['auto_tf'])
        self.assertEqual(len(cfg.keys()), 3)
        self.assertEqual(cfg['secretsmanager_orc8r_secret'], "jar")
        self.assertEqual(cfg['nms_db_password'], "foo")
        self.assertEqual(cfg['lte_orc8r_chart_version'], "0.2.4")

        # verify if jinja template has been rendered accordingly
        with open(self.constants['main_tf']) as f:
            jinja_cfg = dict(ln.split('=')
                             for ln in f.readlines() if ln.strip())

        # all infra terraform keys should be present in the jinja template
        self.assertEqual(
            set(jinja_cfg.keys()),
            set(['secretsmanager_orc8r_secret']))

        # variable tf is of the form variable "var_name" {}
        # get the middle element and remove the quotes
        with open(self.constants['vars_tf']) as f:
            jinja_cfg = [ln.split()[1][1:-1]
                         for ln in f.readlines() if ln.strip()]

        # all infra terraform keys should be present in the jinja template
        self.assertEqual(set(jinja_cfg), set(mgr.tf_vars))

    @patch("configlib.get_input")
    @patch("utils.awslib.run_command")
    def test_configure_set(self, mock_run_command, mock_get_input):
        mock_vals = {
            'nms_db_password': 'foo',
        }
        mock_get_input.side_effect = [
            mock_vals['nms_db_password'],
        ]
        mock_run_command.return_value = 0

        # configure platform
        mgr = configlib.ConfigManager(self.constants)
        mgr.configure('platform')
        mgr.commit('platform')

        # check if we only invoked input for nms_db_password
        self.assertEqual(mock_get_input.call_count, 1)
        fn = "%s/platform.tfvars.json" % self.constants['config_dir']
        cfg = get_json(fn)
        self.assertEqual(len(cfg.keys()), 1)
        self.assertEqual(cfg['nms_db_password'], "foo")

        # set a specific variable
        mgr.set('platform', 'deploy_elastic', 'true')
        mgr.commit('platform')

        cfg = get_json(fn)
        self.assertEqual(len(cfg.keys()), 2)
        self.assertEqual(cfg['deploy_elastic'], "true")

        # finally verify if all configs required by tf is present
        cfg = get_json(self.constants['auto_tf'])
        self.assertEqual(len(cfg.keys()), 2)
        self.assertEqual(cfg['nms_db_password'], "foo")
        self.assertEqual(cfg['deploy_elastic'], "true")


if __name__ == '__main__':
    main()
