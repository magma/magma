"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import subprocess  # noqa: S404
import sys


def main() -> None:
    """Provide command-line options to flatten MAGMA-MME OAI image"""
    args = _parse_args()
    status = perform_flattening(args.tag)
    sys.exit(status)


def _parse_args() -> argparse.Namespace:
    """Parse the command line args

    Returns:
        argparse.Namespace: the created parser
    """
    parser = argparse.ArgumentParser(description='Flattening Image')

    parser.add_argument(
        '--tag', '-t',
        action='store',
        required=True,
        help='Image Tag in image-name:image tag format',
    )
    return parser.parse_args()


def perform_flattening(tag):
    """Parse the command line args

    Args:
        tag: Image Tag in image-name:image tag format

    Returns:
        int: pass / fail status
    """
    # First detect which docker/podman command to use
    cli = ''
    image_prefix = ''
    cmd = 'which podman || true'
    podman_check = subprocess.check_output(cmd, shell=True, universal_newlines=True)  # noqa: S602
    if podman_check.strip():
        cli = 'sudo podman'
        image_prefix = 'localhost/'
        # No more need to flatten with --squash option
        return 0
    else:
        cmd = 'which docker || true'
        docker_check = subprocess.check_output(cmd, shell=True, universal_newlines=True)  # noqa: S602
        if docker_check.strip():
            cli = 'docker'
            image_prefix = ''
        else:
            print('No docker / podman installed: quitting')
            return -1

    print(f'Flattening {tag}')
    # Creating a container
    cmd = cli + ' run --name test-flatten --entrypoint /bin/true -d ' + tag
    print(cmd)
    subprocess.check_call(cmd, shell=True, universal_newlines=True)  # noqa: S602

    # Export / Import trick
    cmd = cli + ' export test-flatten | ' + cli + ' import '
    # Bizarro syntax issue with podman
    if cli == 'docker':
        cmd += ' --change "ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" '
    else:
        cmd += ' --change "ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" '
    cmd += ' --change "WORKDIR /magma-mme" '
    cmd += ' --change "EXPOSE 3870/tcp" '
    cmd += ' --change "EXPOSE 5870/tcp" '
    cmd += ' --change "EXPOSE 2123/udp" '
    cmd += ' --change "CMD [\\"sleep\\", \\"infinity\\"]" '  # noqa: WPS342
    cmd += ' - ' + image_prefix + tag
    print(cmd)
    subprocess.check_call(cmd, shell=True, universal_newlines=True)  # noqa: S602

    # Remove container
    cmd = cli + ' rm -f test-flatten'
    print(cmd)
    subprocess.check_call(cmd, shell=True, universal_newlines=True)  # noqa: S602

    # At this point the original image is a dangling image.
    # CI pipeline will clean up (`image prune --force`)
    return 0


if __name__ == '__main__':
    main()
