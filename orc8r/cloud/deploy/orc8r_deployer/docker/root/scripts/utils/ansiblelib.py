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
from dataclasses import dataclass, field

from ansible import context
from ansible.executor.playbook_executor import PlaybookExecutor
from ansible.inventory.manager import InventoryManager
from ansible.module_utils.common.collections import ImmutableDict
from ansible.parsing.dataloader import DataLoader
from ansible.vars.manager import VariableManager
from utils.common import execute_command


@dataclass
class AnsiblePlay:
    """ Model defining ansible play
    """
    playbook: str
    tags: list[str]
    extra_vars: dict[str, any]
    skip_tags: list[any] = field(default_factory=list)
    inventory: str = ''
    verbosity: int = 3


def run_playbook(play: AnsiblePlay) -> int:
    """Run ansible playbook

    Args:
        play (AnsiblePlay): object describing the current play

    Returns:
        int: return code
    """
    if play.inventory:
        env = {"ANSIBLE_HOST_KEY_CHECKING": "False"}
        return execute_command(
            [
            "ansible-playbook",
            "-i",
            play.inventory,
            "-e",
            json.dumps(play.extra_vars),
            "--tags",
            ",".join(play.tags),
            play.playbook,
            ], env=env,
        )

    context.CLIARGS = ImmutableDict(
        tags=play.tags,
        skip_tags=play.skip_tags,
        connection='smart',
        verbosity=play.verbosity,
        forks=10,
        become=None,
        become_method=None,
        become_user=None,
        check=False,
        syntax=None,
        start_at_task=None,
        diff=False,
    )
    loader = DataLoader()
    variable_manager = VariableManager(loader=loader)
    variable_manager.extra_vars.update(play.extra_vars)
    inventory = InventoryManager(loader=loader)
    variable_manager.set_inventory(inventory)
    passwords = {}
    pbex = PlaybookExecutor(
        playbooks=[play.playbook],
        inventory=inventory,
        variable_manager=variable_manager,
        loader=loader,
        passwords=passwords,
    )
    return pbex.run()
