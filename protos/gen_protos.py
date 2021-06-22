"""
Copyright 2020 The Magma Authors.

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
import sys
from typing import Iterable, List

from grpc.tools import protoc


def find_all_proto_files_in_dir(input_dir: str) -> List[str]:
    """
    Returns a list of filenames of .proto files in the given directory

    Args:
        input_dir: Directory to search in

    Returns:
        list(str): List of .proto filenames in the directory
    """
    proto_files = []
    for root, _, names in os.walk(input_dir):
        for name in names:
            full_filename = os.path.join(root, name)
            extn = os.path.splitext(name)[1]

            # Recurse into subdirectories
            if os.path.isdir(full_filename):
                proto_files += find_all_proto_files_in_dir(full_filename)
            if os.path.isfile(full_filename) and extn == '.proto':
                # TODO: find a better way to exclude the prometheus proto
                if not full_filename.endswith('prometheus/metrics.proto'):
                    proto_files.append(full_filename)
    return proto_files


def gen_bindings(
        input_dir: str,
        include_paths: Iterable[str],
        proto_path: str,
        output_dir: str,
) -> None:
    """
    Generates python and Go bindings for all .proto files in input dir
       @input_dir - input directory with .proto files to generate code for
       @include_paths - a list of include paths to resolve relative imports in .protos
       @output_dir - output directory to put generated code in
    """
    protofiles = find_all_proto_files_in_dir(input_dir)

    inouts = [
        '-I' + proto_path,
        '--python_out=' + output_dir,
        '--grpc_python_out=' + output_dir,
    ]
    # Only run mypy (dev dependency) when the protoc-consumed executable exists
    if shutil.which('protoc-gen-mypy') is not None:
        inouts.append('--mypy_out=' + output_dir)

    protoc.main(
        ('',)
        + tuple('-I' + path for path in include_paths)
        + tuple(inouts)
        + tuple(f for f in protofiles),
    )


def main():
    """
    Default main module. Generates .py code for all proto files
    specified by the arguments
    """
    if len(sys.argv) != 5:
        print(
            "Usage: ./gen_protos.py <dir containing .proto's> <include paths CSV> <proto_path for imports> <output dir>",
        )
        exit(1)
    input_dir = sys.argv[1]
    include_paths = sys.argv[2].split(',')
    # The deprecated vagrant box image amarpad/magma_dev has grpc installed
    # from source, with header files located at /usr/local/include. In the new
    # box image amarpad/debian_jessie, grpc is installed from deb package, with
    # headers located at /usr/include. Currently, only the magma_test vm uses
    # amarpad/magma_dev.
    #
    # TODO: Migrate magma_test to amarpad/debian_jessie to to remove the
    # '/usr/local/include'
    include_paths.append('/usr/include')
    include_paths.append('/usr/local/include')
    proto_path = sys.argv[3]
    output_dir = sys.argv[4]
    gen_bindings(input_dir, include_paths, proto_path, output_dir)


if __name__ == "__main__":
    main()
