#  Copyright 2018 SAS Project Authors. All Rights Reserved.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

"""A simple Certificate Revocation List (CRL) server implementation."""

from __future__ import absolute_import, division, print_function

import inspect
import logging
import os
import re
import socket
import subprocess
import sys
import threading

import six
from six.moves import input, zip
from six.moves.BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
from six.moves.urllib import parse as urllib

CRLDP_BASE = 'http://fake-crl-service:9007/'
DEFAULT_CRL_SERVER_PORT = 80
MENU_ACTIONS = {}
MENU_ID_MAIN_MENU = '-1'
MENU_ID_QUIT = '0'
MENU_ID_REVOKE_CERT = '1'
MENU_ID_CRL_UPDATE_REGENERATE = '2'
BANNER_TITLE = r'''   _____ _____  _         _____
  / ____|  __ \| |       / ____|
 | |    | |__) | |      | (___   ___ _ ____   _____ _ __
 | |    |  _  /| |       \___ \ / _ \ '__\ \ / / _ \ '__|
 | |____| | \ \| |____   ____) |  __/ |   \ V /  __/ |
  \_____|_|  \_\______| |_____/ \___|_|    \_/ \___|_|
 '''
MENU_OPTIONS = '''
Please select the options:
[1] Revoke Certificate
[2] Update CRL URL and Re-generate certificates
[0] Quit'''
CLI_PROMPT_TEXT = '''CRL Server> '''


class SimpleCrlServer(threading.Thread):
    """Implements a simple CRL server by creating a HTTP server.
    A simple HTTP server that responds with the ca.crl file present in the 'harness/certs/crl/' directory serving CRL files.
    The CRL server runs as a separate thread.
    """

    def __init__(self, crl_url, crl_directory):
        """Constructor for simple CRL Server.
        Args:
          crl_url: Prefix of the CRL distribution point URL listed in certificates.
              Must end with a slash. We will serve *.crl files under this virtual
              directory.
          crl_directory: Absolute path of the local directory containing *.crl
              files to serve.
        """

        threading.Thread.__init__(self)
        url_components = urllib.urlparse(crl_url)
        self.port = url_components.port if url_components.port is not None else \
            DEFAULT_CRL_SERVER_PORT
        self.setDaemon(True)

        # Setting the parameters for handler class.
        self.server = CrlHttpServer(
            url_components.path, crl_directory,
            ('', self.port),
            CrlServerHttpHandler,
        )

    def run(self):
        """Starts the HTTPServer as background thread.
        The run() method is an overridden method from thread class and it gets invoked
        when the thread is started using thread.start().
        """

        logging.info("CRL Server listening on port %d", self.port)
        self.server.serve_forever()

    def stopServer(self):
        """This method is used to stop HTTPServer Socket"""

        self.server.shutdown()
        logging.info("Stopped CRL Server on port %d", self.port)


class CrlHttpServer(HTTPServer):
    """The class CrlHttpServer overrides built-in HTTPServer.
    It takes the parameter base_path used to serve the files.
    This is needed to have CRL Server to serve files from different directories.
    """

    def __init__(
        self, crl_url_path, crl_directory, server_address,
        RequestHandlerClass,
    ):
        if not crl_url_path.endswith('/'):
            raise ValueError("crl_url_path %r must end with a slash" % crl_url_path)
        self.crl_url_path = crl_url_path
        self.crl_directory = crl_directory
        HTTPServer.__init__(self, server_address, RequestHandlerClass)


class CrlServerHttpHandler(BaseHTTPRequestHandler):
    """This class implements the HTTP handlers for the CRL Server.
    CrlServerHttpHandler class inherits with BaseHTTPRequestHandler
    to serve HTTP Response.
    """

    def do_GET(self):
        """Handles Pull/GET Request and returns path of the request to callback method."""

        req_path = urllib.urlparse(self.path).path
        self.log_message('Received GET Request: {}'.format(req_path))

        # Check the requested prefix.
        if not req_path.startswith(self.server.crl_url_path):
            self.log_message(
                'Requested path %r does not match URL prefix %r ' %
                (req_path, self.server.crl_url_path),
            )
            self.send_response(404)
            return

        # Restrict the filename, to prevent arbitrary filesystem access.
        req_filename = req_path[len(self.server.crl_url_path):]
        if not isCrlFilename(req_filename):
            self.log_message('Requested file %r is malformed' % req_filename)
            self.send_response(404)
            return

        # Try to read the requested CRL file.
        abs_filename = os.path.join(self.server.crl_directory, req_filename)
        try:
            with open(abs_filename, 'rb') as handle:
                file_data = handle.read()
        except IOError:
            self.log_message('Failed to read CRL file: %r' % abs_filename)
            self.send_response(404)
            return
        else:
            self.send_response(200)
            self.send_header("Content-type", "application/pkix-crl")
            self.end_headers()
            self.wfile.write(file_data)
            return


def getCertsDirectoryAbsolutePath():
    """This function will return absolute path of 'certs' directory."""

    harness_dir = os.path.dirname(os.path.abspath(inspect.getfile(inspect.currentframe())))
    path = os.path.join(harness_dir, 'certs')
    return path


def isCrlFilename(filename):
    """Returns true if the input resembles a well-formed CRL filename."""
    return re.match(r"^[a-z_]+[.]crl$", filename)


def getCertificateNameToBlacklist():
    """Allows the user to select the certificate to blacklist from a displayed list."""

    # Gets the all certificate files that are ends with .cert present in the
    # default certs directory.
    cert_files = [
        cert_file_name for cert_file_name in os.listdir(
            getCertsDirectoryAbsolutePath(),
        ) if cert_file_name.endswith('.cert')
    ]

    # Creates the mapping menu ID for each certificate file to displayed list.
    cert_files_with_options = dict(zip(range(len(cert_files)), sorted(cert_files)))
    print("Select the certificate to blacklist:")
    print(
        "\n".join(
            "[%s] %s" % (cert_file_option_id, cert_file_name)
            for (cert_file_option_id, cert_file_name) in
            sorted(six.iteritems(cert_files_with_options))
        ),
    )
    option_id = int(input(CLI_PROMPT_TEXT))
    if option_id not in cert_files_with_options:
        raise Exception(
            'RunTimeError:Invalid input:{}. Please select a valid '
            'certificate'.format(option_id),
        )
    return cert_files_with_options[option_id]


def revokeCertificate():
    """Revokes a leaf certificate and updates the CRL chain file."""

    # Retrives the all certificates from certs directory.
    cert_name = getCertificateNameToBlacklist()
    logging.info('Certificate selected to blacklist is:%s', cert_name)
    revoke_cert_command = "bash ./revoke_and_generate_crl.sh " \
                          "-r {0}/{1}".format(
                              getCertsDirectoryAbsolutePath(),
                              cert_name,
                          )
    command_exit_status = subprocess.call(revoke_cert_command, shell=True)
    if command_exit_status:
        raise Exception(
            "RunTimeError:Certificate Revocation is failed:"
            "exit_code:{}".format(command_exit_status),
        )
    logging.info('%s is blacklisted successfully', cert_name)


def updateCrlUrlAndRegenerateCertificates():
    """Updates the CRLDP field and re-generates certificates.
    Rewrites include_crldp_base.sh, which is consumed by generate_fake_certs.sh
    and related scripts.  Invokes the generate_fake_certs.sh script to regenerate
    all certificates.
    """

    # Populate the crldp file.
    certs_path = getCertsDirectoryAbsolutePath()
    with open(os.path.join(certs_path, 'include_crldp_base.sh'), 'w') as handle:
        handle.write('CRLDP_BASE="{}"\n'.format(CRLDP_BASE))

    # Re-generate certificates.
    script_name = './generate_fake_certs.sh'
    command = 'cd {0} && {1}'.format(certs_path, script_name)
    command_exit_status = subprocess.call(command, shell=True)
    if command_exit_status:
        raise Exception(
            "RunTimeError:Regeneration of certificates failed:"
            "exit_code:{}".format(command_exit_status),
        )
    logging.info("The certs are regenerated successfully")

    # Generate CRL Chain.
    generateCrlChain()


def generateCrlChain():
    """Generates CRL files for each CA.
    DER-encoded *.crl files are placed in the 'harness/certs/crl' directory.
    """
    create_crl_chain_command = "bash ./revoke_and_generate_crl.sh -u"
    command_exit_status = subprocess.call(create_crl_chain_command, shell=True)
    if command_exit_status:
        raise Exception(
            "RunTimeError: Failed to generate CRL files "
            "exit_code:{}".format(command_exit_status),
        )

    crl_files = []
    try:
        for filename in os.listdir(
            os.path.join(getCertsDirectoryAbsolutePath(), "../crl"),
        ):
            if isCrlFilename(filename):
                crl_files.append(filename)
    except EnvironmentError as e:
        raise Exception("RunTimeError: Failed to list crl directory: %s" % e)
    if not crl_files:
        raise Exception("RunTimeError: No CRL files found")
    logging.info("Serving CRL files: %r", sorted(crl_files))


def crlServerStart():
    """Start CRL server."""

    # Create CRL chain containing CRLs of sub CAs and root CA.
    generateCrlChain()

    try:
        crl_server = SimpleCrlServer(
            CRLDP_BASE,
            os.path.join(getCertsDirectoryAbsolutePath(), '../crl'),
        )
        crl_server.start()
        logging.info("CRL Server has been started")
        logging.info('URL pattern: %s<ca_name>.crl' % CRLDP_BASE)
    except socket.error as err:
        raise Exception(
            "RunTimeError:There is an error starting CRL Server:"
            "exit_reason:{}".format(err.strerror),
        )

    return crl_server


def crlServerStop(crl_server):
    """Stop CRL Server.
    Args:
      crl_server: A CRL Server object instance.
    """

    crl_server.stopServer()


def readInput():
    """Display the CRL server menu and executes the selected option."""

    print(BANNER_TITLE + MENU_OPTIONS)
    choice = input(CLI_PROMPT_TEXT)
    executeSelectedMenu(choice)
    return


def executeSelectedMenu(choice):
    """Performs the action based on the selected menu item.
    Args:
      choice: The selected index value from the main menu of CRL server.
    """

    if choice == '':
        MENU_ACTIONS[MENU_ID_MAIN_MENU]()
    else:
        try:
            MENU_ACTIONS[choice]()
        except KeyError:
            print("Invalid selection, please try again")
        except Exception as err:
            logging.error(str(err))
        MENU_ACTIONS[MENU_ID_MAIN_MENU]()
    return


# Mapping menu items to handler functions.
MENU_ACTIONS = {
    MENU_ID_MAIN_MENU: readInput,
    MENU_ID_REVOKE_CERT: revokeCertificate,
    MENU_ID_CRL_UPDATE_REGENERATE: updateCrlUrlAndRegenerateCertificates,
    MENU_ID_QUIT: exit,
}

# Set the logger for CRL Server.
logger = logging.getLogger()
handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(logging.Formatter('[%(levelname)s] %(asctime)s %(message)s'))
logger.addHandler(handler)
logger.setLevel(logging.INFO)

if __name__ == '__main__':
    # Start Simple CRL Server and waits for user input.
    try:
        #crl_server_instance = crlServerStart()
        readInput()
        crlServerStop(crl_server_instance)
    except Exception as err:
        logging.error(str(err))
