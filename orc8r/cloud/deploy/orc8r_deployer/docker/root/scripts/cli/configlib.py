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

import json
import os

import click
import yaml
from cli.style import print_error_msg, print_success_msg, print_title
from jinja2 import Environment, FileSystemLoader
from prettytable import PrettyTable
from utils.awslib import set_aws_configs
from utils.common import get_json, put_json


def get_input(text, default_val):
    if default_val:
        resp = click.prompt(text, default=default_val)
    else:
        resp = click.prompt(text, default=default_val, show_default=False)

    # strip any quotes
    resp = resp.strip("\'").strip("\"")
    return resp


def input_to_type(input_str: str, input_type: str):
    ''' convert input string to its appropriate type '''
    if input_type == 'map':
        return json.loads(input_str)
    elif input_type == 'bool':
        return input_str.lower() in ["true", "True"]
    return input_str


def add_pretty_table(fields, items):
    table = PrettyTable()
    table.field_names = fields
    for field in fields:
        table.align[field] = 'l'

    for item in items:
        table.add_row(item)
    return table


def render_j2_template(src_dir, dst_fn, cfg):
    env = Environment(loader=FileSystemLoader(searchpath=src_dir))
    fn = os.path.basename(dst_fn)
    template = env.get_template(f'{fn}.j2')
    try:
        with open(dst_fn, "w") as f:
            f.write(template.render(cfg=cfg))
    except Exception as err:
        click.echo(f"Error: {fn} rendering {err!r} file ")


class ConfigManager(object):
    '''
    ConfigManager manages the orcl configuration. It reads the variable
    definitions during init and uses the information parsed to configure
    the various component attributes.
    Currently the variables are used to configure aws cli and
    terraform.
    '''

    def _get_config_fn(self, component: str):
        config_dir = self.constants['config_dir']
        return f"{config_dir}/{component}.tfvars.json"

    def set(self, component: str, key: str, value: str):
        click.echo(f"Setting key {key} value {value} "
                   f"for component {component}")
        config_vars = self.config_vars[component]
        config_info = config_vars.get(key)
        if not config_info:
            print_error_msg(f"{key} not a valid attribute in {component}")
            return
        self.configs[component][key] = input_to_type(value, config_info.get('Type'))

    def commit(self, component):
        set_aws_configs(self.aws_vars, self.configs[component])
        self._configure_tf(component, self.configs[component])
        put_json(self._get_config_fn(component), self.configs[component])

    def configure(self, component: str):
        print_title(f"\nConfiguring {component} deployment variables")
        cfgs = self.configs[component]

        self.initialize_defaults(component)

        # TODO: use a different yaml loader to ensure we load inorder
        # sort the variables to group the inputs together
        config_vars = self.config_vars[component].items()
        for config_key, config_info in sorted(config_vars, key=lambda s: s[0]):
            if not config_info['Required']:
                continue

            config_desc = config_info.get('Description', config_key).strip('.')
            v = cfgs.get(config_key)
            if v:
                inp = get_input(f"{config_key}({config_desc})", v)
            else:
                inp = get_input(f"{config_key}({config_desc})", v)

            if not inp and v is None:
                if not click.confirm("press 'y' to set empty string and "
                                     "'n' to skip", prompt_suffix=': '):
                    continue
                inp = ""

            cfgs[config_key] = input_to_type(inp, config_info.get('Type'))

    def _configure_tf(self, component: str, cfgs: dict):
        ''' updates the terraform auto configuration and main.tf '''
        auto_tf = self.constants['auto_tf']
        auto_cfgs = get_json(auto_tf)
        for k, v in cfgs.items():
            if k in self.tf_vars:
                auto_cfgs[k] = v
        put_json(auto_tf, auto_cfgs)

        # render main.tf with terraform variables alone
        tf_cfgs = {}
        for component in self.configs:
            tf_cfgs[component] = {}
            for k, v in self.configs[component].items():
                if k in self.tf_vars:
                    tf_cfgs[component][k] = v

        render_j2_template(
            self.constants['tf_dir'],
            self.constants['main_tf'],
            tf_cfgs)
        render_j2_template(
            self.constants['tf_dir'],
            self.constants['vars_tf'],
            self.tf_vars)

    def check(self, component: str) -> bool:
        ''' check if all mandatory options of a specific component is set '''
        cfgs = self.configs[component]
        valid = True
        missing_cfgs = []
        for k, v in self.config_vars[component].items():
            if v['Required'] and cfgs.get(k) is None:
                missing_cfgs.append(k)
                valid = False

        if missing_cfgs:
            print_error_msg(
                f"Missing {missing_cfgs!r} configs for {component} component")
        else:
            print_success_msg(
                f"All mandatory configs for {component} has been configured")
        return valid

    def info(self, component: str):
        ''' pretty click.echo vars yml '''
        print_title(f"{component} Configuration Options")
        fields = ["Name", "Description", "Type", "Required", "Used By"]
        items = []
        for k, v in self.config_vars[component].items():
            items.append([
                k,
                v["Description"],
                v["Type"],
                v["Required"],
                v["ConfigApps"]
            ])
        click.echo(add_pretty_table(fields, items))

    def show(self, component: str):
        ''' pretty click.echo existing configuration '''
        print_title(f"{component} Configuration")
        items = [[k, v] for k, v in self.configs[component].items()]
        fields = ["Name", "Configuration"]

        click.echo(add_pretty_table(fields, items))

    def __init__(self, constants: dict):
        self.config_vars = {}
        self.configs = {}
        self.tf_vars = set()
        self.aws_vars = set()
        self.constants = constants
        vars_fn = constants['vars_definition']
        try:
            with open(vars_fn) as f:
                self.config_vars = yaml.load(f, Loader=yaml.FullLoader)
        except OSError:
            click.echo(f"Failed opening vars file {vars_fn}")

        # read configs
        for component in constants['components']:
            self.configs[component] = {}
            try:
                fn = self._get_config_fn(component)
                self.configs[component] = get_json(fn)
            except OSError:
                pass

            for config_key, config_info in self.config_vars[component].items():
                if 'tf' in config_info["ConfigApps"]:
                    self.tf_vars.add(config_key)

                if 'awscli' in config_info["ConfigApps"]:
                    self.aws_vars.add(config_key)

    def initialize_defaults(self, component):
        cfgs = self.configs[component]
        for config_key, config_info in self.config_vars[component].items():
            # add defaults to configs inorder to run prechecks
            default = config_info.get('Default')
            typ = config_info.get('Type')
            if default is not None:
                curr_val = cfgs.get(config_key)

                # if defaults are overriden already, explicitly confirm
                # to reset defaults
                if (not curr_val or (curr_val and default != curr_val and
                    click.confirm(f"Override {config_key} "
                    f"current val {curr_val} with default {default}"))):
                    cfgs[config_key] = default
