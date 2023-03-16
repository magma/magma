"""
Copyright 2023 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os


def get_readme_path(version):
    """Get the path to the readme folder for a given version."""
    if version == 'latest':
        return 'readmes'
    return f'docusaurus/versioned_docs/version-{version}'


def extract_doc_id(filename, root, version_prefix):
    """Extract the doc id from a given file."""
    path = os.path.join(root, filename)
    doc_id = ""
    with open(path) as f:
        lines = f.readlines()
        if lines and lines[0].startswith('---'):
            for line in lines:
                if line.startswith('id: '):
                    doc_id = line.replace(f'id: {version_prefix}', '').rstrip('\n')
                    break
        else:
            doc_id = filename.replace('.md', '')
    return doc_id


def get_version_prefix(version):
    """Get the version prefix for a given version."""
    if version == 'latest':
        return ''
    return f'version-{version}-'
