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
from jinja2 import Environment, FileSystemLoader
from prettytable import PrettyTable


def get_input(text, default_val):
    if default_val:
        resp = click.prompt(text, default=default_val)
    else:
        resp = click.prompt(text, default=default_val, show_default=False)

    # strip any quotes
    resp = resp.strip("\'").strip("\"")
    return resp


def get_json(fn: str) -> dict:
    try:
        with open(fn) as f:
            return json.load(f)
    except OSError:
        pass
    return {}


def put_json(fn: str, cfgs: dict):
    with open(fn, 'w') as outfile:
        json.dump(cfgs, outfile)


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
        click.echo(
            f"Setting key {key} value {value} for component {component}")
        cfgs = get_json(self._get_config_fn(component))
        config_vars = self.config_vars[component]

        if not config_vars.get(key):
            click.echo("not a valid key")
            return
        cfgs[key] = value

        self._configure_aws(component, cfgs)
        self._configure_tf(component, cfgs)
        put_json(self._get_config_fn(component), cfgs)

    def configure(self, component: str):
        click.echo(
            click.style(
                f"\nConfiguring {component} deployment variables ",
                underline=True))
        cfgs = self.configs[component]
        # TODO: use a different yaml loader to ensure we load inorder
        # sort the variables to group the inputs together
        config_vars = self.config_vars[component].items()
        for config_key, config_info in sorted(config_vars, key=lambda s: s[0]):
            config_description = config_info.get(
                'Description', config_key).strip('.')
            defaultValue = config_info.get('Default')
            # add defaults to the json configs to ensure we can run prechecks
            if defaultValue:
                cfgs[config_key] = defaultValue
                continue

            if not config_info['Required']:
                continue

            v = cfgs.get(config_key)
            if v:
                inp = get_input(f"{config_key}({config_description})", v)
            else:
                inp = get_input(f"{config_key}({config_description})", v)

            # strip quotes from input
            if inp:
                cfgs[config_key] = inp
            else:
                if v is None:
                    if click.confirm("press 'y' to set empty string and "
                                     "'n' to skip", prompt_suffix=': '):
                        cfgs[config_key] = ""

        self.configs[component] = cfgs
        self._configure_aws(component, cfgs)
        self._configure_tf(component, cfgs)
        put_json(self._get_config_fn(component), cfgs)

    def _configure_aws(self, component: str, cfgs: dict):
        ''' configures aws cli with configuration '''
        for k, v in cfgs.items():
            if k in self.aws_vars:
                check_call(["aws", "configure", "set", k, v])

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
            click.echo(
                f"Missing {missing_cfgs!r} configs for {component} component")
        else:
            click.echo(
                f"All mandatory configs for {component} has been configured")
        return valid

    def info(self, component: str):
        ''' pretty click.echo vars yml '''
        click.echo(f"{component} Configuration Options")
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
        click.echo(f"{component} Configuration")
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

        # read all configs
        for component in constants['components']:
            try:
                fn = self._get_config_fn(component)
                self.configs[component] = get_json(fn)
                for k, v in self.config_vars[component].items():
                    if 'tf' in v["ConfigApps"]:
                        self.tf_vars.add(k)
                    if 'awscli' in v["ConfigApps"]:
                        self.aws_vars.add(k)
            except OSError:
                pass
