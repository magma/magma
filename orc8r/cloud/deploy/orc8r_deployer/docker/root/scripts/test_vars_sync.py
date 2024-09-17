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
import re
import sys

import yaml

# Following script reads all the variables in vars.yml, variables.tf and output.tf
# identifies variables missing in vars.yml and asserts. This script is to be used
# for flagging variables file being out of sync with orcl vars file


def read_all_vars(constants: dict) -> set:
    """Read all variables present in vars.yml and return the entire set
    Args:
        constants: config read from config.yml
    Returns:
        all variables defined in vars.yml
    """
    vars_fn = constants['vars_definition']
    try:
        with open(vars_fn) as f:
            config_vars = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        print(f"Failed opening vars file {vars_fn}")

    existing_vars = set()
    for component in constants['components']:
        for k in config_vars[component]:
            existing_vars.add(k)

    return existing_vars


def parse_tf(constants: dict) -> set:
    """Read user configured variables in variables tf and return the entire set
    Args:
        constants: config read from config.yml
    Returns:
        all variables defined in variables.tf
    """
    magma_root = constants['magma_root']
    tf_root = f'{magma_root}/orc8r/cloud/deploy/terraform'

    # parse input variables
    orc8r_var_fn_list = [
        f'{tf_root}/orc8r-aws/variables.tf',
        f'{tf_root}/orc8r-helm-aws/variables.tf',
    ]
    var_search_re = re.compile(r'variable\s\"(?P<variable_name>\w+?)\"\s\{')
    actual_vars = set()
    for fn in orc8r_var_fn_list:
        with open(fn) as f:
            # parse variables
            for line in f.readlines():
                m = var_search_re.search(line)
                if m and m.group('variable_name'):
                    actual_vars.add(m.group('variable_name'))

    # parse output variables
    orc8r_var_fn_list = [
        f'{tf_root}/orc8r-aws/outputs.tf',
    ]

    # remove variables which are set through outputs
    output_search_re = re.compile(r'output\s\"(?P<variable_name>\w+?)\"\s\{')
    for fn in orc8r_var_fn_list:
        with open(fn) as f:
            # parse variables
            for line in f.readlines():
                m = output_search_re.search(line)
                if m and m.group('variable_name'):
                    output_var = m.group('variable_name')
                    if output_var in actual_vars:
                        actual_vars.remove(output_var)
    return actual_vars


if __name__ == '__main__':
    try:
        with open("/root/config.yml") as f:
            constants = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        print("Failed opening config.yml file")
        sys.exit(1)

    existing_vars = read_all_vars(constants)
    actual_vars = parse_tf(constants)
    if not actual_vars.issubset(existing_vars):
        diff = actual_vars.difference(existing_vars)
        print(f"missing variable definitions {diff!r}")
        assert(not diff)

    print("vars.yml found to be in sync with terraform modules")
