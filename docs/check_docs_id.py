#!/usr/bin/env python3
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
from typing import Dict

from common import extract_doc_id, get_readme_path, get_version_prefix


def get_all_pages(version: str = "latest") -> Dict[str, str]:
    """
    Scrape the relevant docs folder to get all relevant page ids.

    Args:
        version (str): The version of the docs to check. Defaults to "latest".

    Returns:
        Set[str]: A set of all pages that are available for this version.
    """
    readme_path = get_readme_path(version)
    all_pages = {}
    for root, _, filenames in os.walk(readme_path):
        for filename in filenames:
            if filename.endswith('.md'):
                doc_id = extract_doc_id(filename, root, get_version_prefix(version))
                filename_key = os.path.join(root.replace(f'{readme_path}/', ''), filename.replace('.md', ''))
                all_pages[filename_key] = root.replace(f'{readme_path}/', '') + '/' + doc_id
    return all_pages


def main():
    """Check the doc ids against the filenames."""
    versions = [
        'latest', '1.8.0', '1.7.0', '1.6.X', '1.5.X',
        '1.4.X', '1.3.X', '1.2.X', '1.1.X', '1.0.X', '1.0.0',
    ]
    flag = False
    print("Checking docs id against docs filenames...")
    for version in versions:
        print('-------------')
        print(f"Version: {version}")
        pages = get_all_pages(version=version)
        for k, v in pages.items():
            if k != v:
                flag = True
                print(f"Error: Expected filename {v}.md, but got {k}.md")

    if flag:
        print('-------------')
        print("Error: Some docs id are not matching with docs filenames.")
        exit(1)


if __name__ == "__main__":
    main()
