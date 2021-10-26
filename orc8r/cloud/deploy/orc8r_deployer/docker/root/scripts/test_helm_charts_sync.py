#!/usr/bin/env python
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


import logging
import re
import sys

import yaml

logging.basicConfig(stream=sys.stdout, level=logging.INFO)

charts_fn_map = {
    'orc8r': 'orc8r/cloud/helm/orc8r/Chart.yaml',
    'fbinternal-orc8r': 'fbinternal/cloud/helm/fbinternal-orc8r/Chart.yaml',
    'feg-orc8r': 'feg/cloud/helm/feg-orc8r/Chart.yaml',
    'lte-orc8r': 'lte/cloud/helm/lte-orc8r/Chart.yaml',
    'wifi-orc8r': 'wifi/cloud/helm/wifi-orc8r/Chart.yaml',
    'cwf-orc8r': 'cwf/cloud/helm/cwf-orc8r/Chart.yaml',
}


def read_all_chart_versions(constants: dict) -> dict:
    """
    Read the version information from helm charts and return chart version map
    Args:
        constants: config read from config.yml
    Returns:
        map of chart name to versions
    """
    magma_root = constants['magma_root']
    chart_versions = {}
    for chart_name, chart_fn in charts_fn_map.items():
        with open(f'{magma_root}/{chart_fn}') as chart_f:
            chart_info = yaml.load(chart_f, Loader=yaml.FullLoader)
            chart_name = chart_name.replace('-', '_')
            chart_versions[chart_name] = chart_info['version']
    return chart_versions


def read_tf_chart_versions(constants: dict) -> dict:
    """Read the version information in variables.tf and return chart version map
    Args:
        constants: config read from config.yml
    Returns:
        map of chart name to versions
    """
    magma_root = constants['magma_root']
    tf_root = f'{magma_root}/orc8r/cloud/deploy/terraform'

    # parse input variables
    orc8r_app_tf = f'{tf_root}/orc8r-helm-aws/variables.tf'
    tf_chart_versions = {}
    version_search_re = r'variable\s\"(?P<chart_name>\w+?)_chart_version\"\s\{.*?default\s*=\s*\"(?P<chart_version>\d+\.\d+\.\d+)\".*?\}'
    with open(orc8r_app_tf) as f:
        # parse variables
        vars_content = f.read()
        for m in re.finditer(
            version_search_re, vars_content,
            re.MULTILINE | re.DOTALL,
        ):
            chart_name = (m.group('chart_name'))
            chart_version = (m.group('chart_version'))
            tf_chart_versions[chart_name] = chart_version
    return tf_chart_versions


def read_orcl_chart_versions(constants) -> dict:
    """Read the version information in vars.yml and return chart version map
    Args:
        constants: config read from config.yml
    Returns:
        map of chart name to versions
    """
    vars_fn = constants['vars_definition']
    try:
        with open(vars_fn) as f:
            config_vars = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        print(f"Failed opening vars file {vars_fn}")

    orcl_chart_versions = {}
    for component in constants['components']:
        for k, v in config_vars[component].items():
            if k.endswith('_chart_version'):
                chart_name = k.split('_chart_version')[0]
                orcl_chart_versions[chart_name] = v['Default']
    return orcl_chart_versions


if __name__ == '__main__':
    try:
        with open("/root/config.yml") as f:
            constants = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        print("Failed opening config.yml file")
        sys.exit(1)

    chart_versions = read_all_chart_versions(constants)
    tf_chart_versions = read_tf_chart_versions(constants)
    orcl_chart_versions = read_orcl_chart_versions(constants)

    logging.debug('Actual Chart Versions %s', repr(chart_versions))
    logging.debug(
        'Chart versions in tf variables file %s',
        repr(tf_chart_versions),
    )
    logging.debug('Chart versions in orcl %s', repr(orcl_chart_versions))

    tf_charts_in_sync = True
    orcl_charts_in_sync = True
    for k, v in chart_versions.items():
        tf_chart_version = tf_chart_versions.get(k)
        if tf_chart_version and tf_chart_version != v:
            print(
                f'Actual chart version for {k} = {v}, tf variables chart version is {tf_chart_version}',
            )
            tf_charts_in_sync = False

        orcl_chart_version = orcl_chart_versions.get(k)
        if orcl_chart_version and orcl_chart_version != v:
            print(
                f'Actual chart version for {k} = {v}, orcl variables chart version is {orcl_chart_version}',
            )
            orcl_charts_in_sync = False

    assert tf_charts_in_sync and orcl_charts_in_sync
    print("helm charts are in sync")
