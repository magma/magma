"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import os
import sys

from grpc.tools import protoc


def gen_prometheus_proto_py(proto_file_dir, output_dir):
    # Function For fb-internal build tools - open source should use this file
    # as a script
    protoc.main(
        (
            '-I ' + proto_file_dir,
            '--proto_path=' + proto_file_dir,
            '--python_out=' + output_dir,
            '--grpc_python_out=' + output_dir,
            os.path.join(proto_file_dir, 'metrics.proto'),
        )
    )


if __name__ == '__main__':
    # ./gen_prometheus_proto.py <magma root> <output_dir>
    magma_root, out_dir = sys.argv[1], sys.argv[2]
    file_dir = os.path.join(magma_root, 'orc8r/protos/prometheus')
    gen_prometheus_proto_py(file_dir, out_dir)
