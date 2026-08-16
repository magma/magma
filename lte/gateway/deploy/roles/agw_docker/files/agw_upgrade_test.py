"""
Copyright 2026 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os
import subprocess  # noqa: S404
import tempfile
from pathlib import Path
from unittest import TestCase, main

BASH = '/bin/bash'
SCRIPT = Path(__file__).with_name('agw_upgrade.sh')


class AgwUpgradeTest(TestCase):
    """Verify that AGW upgrades remain scoped to their Compose project."""

    def test_upgrade_only_mutates_compose_project(self) -> None:
        """Keep unrelated containers and referenced images intact."""
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            docker_dir = root / 'docker'
            fake_bin = root / 'bin'
            docker_dir.mkdir()
            fake_bin.mkdir()

            env_execution_marker = root / 'env-was-executed'
            (docker_dir / '.env').write_text(
                'IMAGE_VERSION=new\n'
                'DOCKER_REGISTRY=registry.example/\n'
                f'UNUSED=$(touch {env_execution_marker})\n',
                encoding='utf-8',
            )
            (docker_dir / 'docker-compose.yaml').touch()

            docker = fake_bin / 'docker'
            docker.write_text(
                '#!/bin/bash\n'
                'echo "$*" >> "$DOCKER_LOG"\n'
                'case "$*" in\n'
                '  "ps --filter name=magmad --format {{.Image}}")\n'
                '    echo "agw_gateway_python:old" ;;\n'
                '  "compose --compatibility -f docker-compose.yaml config")\n'
                '    echo "services: valid" ;;\n'
                '  "compose --compatibility -f docker-compose.yaml images -q")\n'
                '    echo "sha256:old-agw-image" ;;\n'
                '  "ps -a -q --filter ancestor=sha256:old-agw-image")\n'
                '    test -z "$ANCESTOR_IN_USE" || echo "unrelated" ;;\n'
                'esac\n',
                encoding='utf-8',
            )
            docker.chmod(0o755)

            pidof = fake_bin / 'pidof'
            pidof.write_text('#!/bin/bash\nexit 1\n', encoding='utf-8')
            pidof.chmod(0o755)

            docker_log = root / 'docker.log'
            env = os.environ.copy()
            env.update({
                'DOCKER_LOG': str(docker_log),
                'PATH': f'{fake_bin}:{env["PATH"]}',
            })
            command = [
                BASH,
                '-c',
                'source "$1"; upgrade_agw "$2"',
                'agw-upgrade-test',
                str(SCRIPT),
                str(docker_dir),
            ]

            subprocess.run(  # noqa: S603
                command,
                check=True,
                env=env,
            )
            calls = docker_log.read_text(encoding='utf-8').splitlines()

            self.assertIn(
                'compose --compatibility -f docker-compose.yaml '
                'down --remove-orphans',
                calls,
            )
            self.assertIn(
                'compose --compatibility -f docker-compose.yaml up -d',
                calls,
            )
            self.assertIn('image rm sha256:old-agw-image', calls)
            self.assertFalse(any(call.startswith('stop ') for call in calls))
            self.assertFalse(any('system prune' in call for call in calls))
            self.assertFalse(env_execution_marker.exists())

            docker_log.unlink()
            env['ANCESTOR_IN_USE'] = '1'
            subprocess.run(  # noqa: S603
                command,
                check=True,
                env=env,
            )
            calls = docker_log.read_text(encoding='utf-8').splitlines()
            self.assertNotIn('image rm sha256:old-agw-image', calls)


if __name__ == '__main__':
    main()
