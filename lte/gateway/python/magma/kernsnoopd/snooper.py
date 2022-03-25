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
import logging
import os
from typing import List

from bcc import BPF  # pylint:disable=import-error
from jinja2 import Template
from magma.common.job import Job
from magma.kernsnoopd.handlers import ByteCounter, ebpf_handlers

EBPF_SRC_DIR = "/var/opt/magma/ebpf/kernsnoopd/"
if not os.path.isdir(EBPF_SRC_DIR):
    EBPF_SRC_DIR = os.path.join(
        os.path.dirname(os.path.abspath(__file__)),
        'ebpf',
    )
EBPF_COMMON_FILE = 'common.bpf.h'


def _get_ebpf_source(filename, context) -> str:
    """
    _get_ebpf_source reads template source from file and renders it with
    context parameters

    Args:
        filename: absolute path of file from which to read template source
        context: dict containing parameter values

    Returns:
        Rendered source contents
    """
    with open(filename, 'r', encoding="utf-8") as src_f:
        src = src_f.read()
    template = Template(src)
    return template.render(context)


class NoSourcesFoundError(Exception):
    """
    NoSourcesFoundError is thrown when Snooper does not find any eBPF programs
    or source files to load into the kernel
    """


class Snooper(Job):
    """
    Snooper is a Job that compiles and loads eBPF programs, registered relevant
    front-end programs as handlers, and periodically calls their handle methods
    """

    def __init__(
        self, programs: list, collect_interval: int,
        service_registry, service_loop,
    ):

        super().__init__(interval=collect_interval, loop=service_loop)
        self._bpf = None
        self._handlers: List[ByteCounter] = []
        self._loop = service_loop
        self._ebpf_programs = programs
        self._service_registry = service_registry
        self._context = {
            'PROXY_PORT': service_registry.get_proxy_config().get(
                'local_port',
            ),
        }
        try:
            self._load_ebpf_programs()
            self.start()
        except NoSourcesFoundError:
            logging.error('Fatal: no eBPF sources loaded')

    def _load_ebpf_programs(self) -> None:
        """
        _load_ebpf_programs reads eBPF templates from _ebpf_programs, renders
        them with context, compiles and loads them into kernel, and registers
        corresponding front-end handlers

        Raises:
            NoSourcesFoundError: self._ebpf_programs was empty or no source in
            self._ebpf_programs could be loaded
        """
        if not self._ebpf_programs:
            raise NoSourcesFoundError()

        sources = []
        for basename in self._ebpf_programs:
            filename = os.path.join(EBPF_SRC_DIR, f'{basename}.bpf.c')
            try:
                sources.append(_get_ebpf_source(filename, self._context))
                handler = ebpf_handlers[basename](self._service_registry)
                self._handlers.append(handler)
            except FileNotFoundError:
                logging.error('Could not open eBPF source file %s', filename)
            except KeyError:
                logging.error('Fatal: did not find handler for %s', basename)

        # found eBPF sources to load into kernel
        if sources:
            # find and prepend header
            header = os.path.join(EBPF_SRC_DIR, EBPF_COMMON_FILE)
            try:
                sources.insert(0, _get_ebpf_source(header, self._context))
                self._bpf = BPF(text='\n'.join(sources))
                logging.info('Loaded sources into kernel')
            except FileNotFoundError:
                logging.error('Fatal: Could not open header file %s', header)
        else:
            raise NoSourcesFoundError()

    async def _run(self) -> None:
        if self._bpf is not None:
            for handler in self._handlers:
                handler.handle(self._bpf)
