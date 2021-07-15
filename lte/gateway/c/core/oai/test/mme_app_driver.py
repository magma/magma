#!/usr/bin/env python3
#
# Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The OpenAirInterface Software Alliance licenses this file to You under
# the terms found in the LICENSE file in the root of this source tree.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -------------------------------------------------------------------------------
# For more information about the OpenAirInterface (OAI) Software Alliance:
#      contact@openairinterface.org
#

import argparse
import logging
import os
import re
import shutil
import signal
import subprocess
import sys
import tempfile
import time

MAGMA_ROOT = os.getenv('MAGMA_ROOT')


class MMEAppDriver(object):
    """Driver class for MME app end-to-end testing"""

    def __init__(self):
        """Creates a new MMEAppDriver.

        Sets up the general framework by caching the general run parameters:
          - self._mme_binary is the path to the MME binary to run
          - self._timeout is the max amount of time to run each test
          - self._sample_period is the polling frequency for each test
        """

        # Instantiate variables
        self._elapsed_time = 0  # Make sure that report() doesn't crash
        self._failed_conditions = None  # Indicate that nothing has run yet
        self._mme_bin = None
        self._mme_conf = None
        self._mme_log = None
        self._mme_pid = None
        self._sampling_range = None

        self._setup()
        self._report_settings()

    @classmethod
    def get_parser(cls):
        """Class method for obtaining an argument parser.

        Returns:
            an ArgumentParser object. This allows for extensions built on top
            to make the super call and continue to build on top of the same
            parser.

        """
        parser = argparse.ArgumentParser(
            description='End-to-end test to ensure MME startup is nonblocking',
        )

        parser.add_argument(
            '-b', '--binary', metavar='PATH', default='~/build/mme_app/mme',
            help='The location of the mme binary to test',
        )
        parser.add_argument(
            '-r', '--rate', metavar='SECONDS', default=1, type=int,
            help='Sampling rate for success conditions, in seconds',
        )
        parser.add_argument(
            '-t', '--timeout', metavar='SECONDS', default=30, type=int,
            help='Maximum duration of the test before timing out, in seconds',
        )

        return parser

    def report(self):
        """Report run results. Logs and returns a tuple of failed conditions

        If no report is cached (i.e. no tests have completed yet), then this
        logs an error and returns None. Otherwise, success is logged in info,
        while failure is logged in error, along with the failed condition
        expressions.

        Returns:
            a tuple() of the failed expressions (an empty tuple on success), or
            None if there is nothing to report.

        """
        if self._failed_conditions is None:  # Empty tuple means something else
            logging.error(
                'No results to report yet. Try calling run() first.',
            )
            return None

        if self._failed_conditions:
            logging.error(
                'Test failed after %.6f seconds!', self._elapsed_time,
            )
            for condition in self._failed_conditions:
                logging.error('..Failed expression: \'%s\'', condition)
        else:
            logging.info(
                'Test succeeded after %.6f seconds!', self._elapsed_time,
            )

        # Copy the tuple to avoid representation exposure
        return tuple(self._failed_conditions)

    def run(self, log_conditions=tuple()):  # noqa: B008 T25377293 Grandfathered in
        """Run the driver, verifying the given log regexes.

        Sets up a temporary directory, which contains all the files generated
        during the test. Conf files, certs, and logs all go in here. After
        generating the files, spins off a process to run the MME binary, while
        the main thread proceeds to verify the conditions in log_conditions
        using the log. Upon success or timeout, the MME process gets killed.

        Args:
            log_conditions: a tuple or list of the regex strings to search for
                in the logs. The ones that fail will appear in the report.
        """
        # Use a temp dir for our custom conf and cert files
        with tempfile.TemporaryDirectory(prefix='/tmp/test') as temp_dir:
            oai_dir = os.path.join(temp_dir, 'oai')
            self._setup_oai(oai_dir)
            try:
                self.start_mme()
                self._verify(log_conditions)
            finally:  # Make sure that we kill our MME instance before we exit
                self.stop_mme()

        self._mme_conf = None
        self._mme_log = None

        self.report()

    def start_mme(self):
        """Starts a new MME instance as a subprocess under the current user.

        Spins up a new MME instance by running the MME binary at the cached
        location. If a pid_basepath is passed in, then it is forwarded to the
        MME binary as a command line argument.

        Note that this instance is run under the current user, and thus will
        have only the permissions of the current user.
        """
        cmd = (self._mme_bin, '-c', self._mme_conf)
        mme_proc = subprocess.Popen(cmd)
        self._mme_pid = mme_proc.pid

    def stop_mme(self):
        """Terminates the MME instance subprocess.

        Sends a SIGTERM to the subprocess. It is possible that the subprocess
        may ignore the SIGTERM signal, in which case manual intervention may be
        necessary to kill the zombie process.
        """
        os.kill(self._mme_pid, signal.SIGTERM)

    def _report_settings(self):
        """Logs the cached framework information.

        Reports the presumed location of the MME binary, the max duration, and
        the sample period.
        """
        logging.info('Assuming that %s is the MME binary.', self._mme_bin)
        logging.info(
            'Test will run for a max of %ds, sampling every %ds.',
            self._sampling_range.stop, self._sampling_range.step,
        )

    def _setup(self):
        """Set driver instance variables.

        Configures the logger, parses the command line arguments.
        """
        logging.basicConfig(level=logging.DEBUG)

        opts = MMEAppDriver.get_parser().parse_args()

        timeout = opts.timeout

        # Sample period should be no longer than the max duration
        sample_period = min(timeout, opts.rate)

        # Generate sampling range now instead of once for each test
        self._sampling_range = range(0, timeout, sample_period)

        # Absolute path to the MME binary
        self._mme_bin = os.path.normpath(os.path.expanduser(opts.binary))
        if not os.path.isfile(self._mme_bin):
            logging.error(
                'No binary found at declared MME bin path: %s', self._mme_bin,
            )
            sys.exit(1)

    def _setup_conf_and_certs(self, oai_dir):
        """Set up conf and cert files.

        Copies the conf files from MAGMA_ROOT/config/oai to the target
        directory, making path corrections as well, for the locations of the
        configuration files and the log. Also caches the location of the
        mme.conf file.

        Generates certificates by making a subprocess call to the certificate
        generation script, placing the generated certificates in the
        freeDiameter subdirectory of the directory path in the argument.

        Args:
            oai_dir: the root of the OAI configuration directory to place the
                certificates in. Note that the certificates will actually be in
                the freeDiameter subdirectory of this directory.
        """

        # Conf file setup
        shutil.copytree(os.path.join(MAGMA_ROOT, 'config/oai'), oai_dir)

        # Replace template /usr/local/etc/oai with the conf dir
        for path, _, files in os.walk(oai_dir):
            for fname in files:
                fname = os.path.join(path, fname)
                with open(fname, 'r') as fin:
                    data = fin.read()

                with open(fname, 'w') as fout:
                    fout.write(
                        data
                        .replace('/usr/local/etc/oai', oai_dir)
                        .replace('/var/run', oai_dir)
                        .replace('/tmp/mme.log', self._mme_log),
                    )

        self._mme_conf = os.path.join(oai_dir, 'mme.conf')

        # Certs setup
        fd_dir = os.path.join(oai_dir, 'freeDiameter')
        os.makedirs(fd_dir, exist_ok=True)
        subprocess.call((
            '%s/config/create_certs.py' % MAGMA_ROOT,
            '-c', fd_dir,
        ))

    def _setup_oai(self, oai_dir):
        """Set up OAI conf and cert files in the specified directory.

        Generates confs, certs, and sets up the mme_log instance variable,
        which contains the location of the log file to inspect.

        Args:
            oai_dir: the temporary directory that will contain the conf and
                cert files.
        """
        self._mme_log = os.path.join(oai_dir, 'log')
        self._setup_conf_and_certs(oai_dir)

    def _verify(self, conditions):
        """Verify the log conditions.

        Checks the log for the conditions, using regular expression search.
        Runs the search query every self._sample_period seconds by putting the
        main thread to sleep for that length of time.

        After either success or failure by timeout, caches the elapsed time and
        a tuple of the failed condition strings.

        Args:
            conditions: the conditions to use to verify. An empty set results
                in no verifications, but will only return upon seeing that the
                log matches 'all' the conditions.
        """
        conditions = set(conditions)  # Optimize out any duplicate conditions
        log = ''  # Default value in case log does not get created
        start_time = time.time()
        for _ in self._sampling_range:
            time.sleep(self._sampling_range.step)
            if os.path.exists(self._mme_log):
                with open(self._mme_log, 'r') as log_file:
                    log = log_file.read()
                    if all(re.search(cond, log) for cond in conditions):
                        break
        self._elapsed_time = time.time() - start_time

        self._failed_conditions = tuple(
            cond for cond in conditions if re.search(cond, log) is None
        )
