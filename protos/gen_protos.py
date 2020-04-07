"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import sys

import os
from grpc.tools import protoc


def find_all_proto_files_in_dir(input_dir):
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


def gen_bindings(input_dir, include_paths, proto_path, output_dir):
    """
    Generates python and Go bindings for all .proto files in input dir
       @input_dir - input directory with .proto files to generate code for
       @include_paths - a list of include paths to resolve relative imports in .protos
       @output_dir - output directory to put generated code in
    """
    protofiles = find_all_proto_files_in_dir(input_dir)
    protoc.main(
        ('',) +
        tuple('-I' + path for path in include_paths) +
        ('--proto_path=' + proto_path,
         '--python_out=' + output_dir,
         '--grpc_python_out=' + output_dir) +
        tuple(f for f in protofiles),
    )


def main():
    """
    Default main module. Generates .py code for all proto files
    specified by the arguments
    """
    if len(sys.argv) != 5:
        print(
            "Usage: ./gen_protos.py <dir containing .proto's> <include paths CSV> <proto_path for imports> <output dir>")
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
