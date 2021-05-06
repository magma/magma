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
import re
import sys

MAX_ALLOWED_WARNINGS = 20
COMMON_TYPE = 'MAGMA-COMMON'
MME_TYPE = 'MAGMA-OAI-MME'
SCTPD_TYPE = 'MAGMA-SCTPD'
U18_BUILD_LOG_FILE = 'build_magma_mme.log'
RHEL8_BUILD_LOG_FILE = 'build_magma_mme_rhel8.log'
REPORT_NAME = 'build_results_magma_oai_mme.html'


class HtmlReport():
    """Creates Executive Summary HTML reports."""

    def __init__(self):
        """Initialize obeject."""
        self.job_name = ''
        self.mode = ''
        self.job_id = ''
        self.job_url = ''
        self.job_start_time = 'TEMPLATE_TIME'
        self.git_url = ''
        self.git_src_branch = ''
        self.git_src_commit = ''
        self.git_src_commit_msg = None
        self.git_merge_request = False
        self.git_target_branch = ''
        self.git_target_commit = ''
        self.errorWarningInfo = []
        self.variant = []

    def generate_build_report(self):
        """Create the BUILD HTML report."""
        cwd = os.getcwd()
        try:
            self.file = open(os.path.join(cwd, REPORT_NAME), 'w')
        except IOError:
            sys.exit('Could not open write output file')
        self.generate_header()

        self.add_build_summary_header()
        self.add_compile_rows()
        self.add_copy_to_target_image_row()
        self.add_copy_conf_tools_to_target_mage_row()
        self.add_image_size_row()
        self.add_build_summary_footer()

        self.add_details()

        self.generate_footer()
        self.file.close()

    def generate_header(self):
        """Append HTML header to file."""
        # HTML Header
        header = '<!DOCTYPE html>\n'
        header += '<html class="no-js" lang="en-US">\n'
        header += '<head>\n'
        header += '  <meta name="viewport" content="width=device-width, initial-scale=1">\n'
        header += '  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">\n'
        header += '  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>\n'
        header += '  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>\n'
        header += '  <title>MAGMA/OAI Core Network Build Results for ' + self.job_name + ' job build #' + self.job_id + '</title>\n'
        header += '</head>\n'
        header += '<body><div class="container">\n'
        header += '  <br>\n'
        header += '  <table width = "100%" style="border-collapse: collapse; border: none;">\n'
        header += '   <tr style="border-collapse: collapse; border: none;">\n'
        # SVG has a invisible background color -- adding it.
        header += '     <td bgcolor="#5602a4" style="border-collapse: collapse; border: none;">\n'
        header += '       <a href="https://www.magmacore.org/">\n'
        header += '          <img src="https://www.magmacore.org/img/magma-logo.svg" alt="" border="none" height=50 width=150>\n'
        header += '          </img>\n'
        header += '       </a>\n'
        header += '     </td>\n'
        header += '     <td align = "center" style="border-collapse: collapse; border: none; vertical-align: center;">\n'
        header += '       <b><font size = "6">Job Summary -- Job: ' + self.job_name + ' -- Build-ID: <a href="' + self.job_url + '">' + self.job_id + '</a></font></b>\n'
        header += '     </td>\n'
        header += '     <td style="border-collapse: collapse; border: none;">\n'
        header += '       <a href="http://www.openairinterface.org/">\n'
        header += '          <img src="http://www.openairinterface.org/wp-content/uploads/2016/03/cropped-oai_final_logo2.png" alt="" border="none" height=50 width=150>\n'
        header += '          </img>\n'
        header += '       </a>\n'
        header += '     </td>\n'
        header += '   </tr>\n'
        header += '  </table>\n'
        header += '  <br>\n'
        self.file.write(header)
        summary = self.generate_build_summary()
        self.file.write(summary)

    def generate_build_summary(self):
        """
        Create build summary string.

        Returns:
            a string with build information.
        """
        summary = ''
        # Build Info Summary
        summary += '  <table class="table-bordered" width = "80%" align = "center" border = "1">\n'
        summary += '    <tr>\n'
        summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-time"></span> Build Start Time</td>\n'
        # date_formatted = re.sub('\..*', '', self.created
        summary += '      <td>' + self.job_start_time + '</td>\n'
        summary += '    </tr>\n'
        summary += '    <tr>\n'
        summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-wrench"></span> Build Trigger</td>\n'
        if self.git_merge_request:
            summary += '      <td>Pull Request</td>\n'
        else:
            summary += '      <td>Push Event</td>\n'
        summary += '    </tr>\n'
        summary += '    <tr>\n'
        summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-cloud-upload"></span> GIT Repository</td>\n'
        summary += '      <td><a href="' + self.git_url + '">' + self.git_url + '</a></td>\n'
        summary += '    </tr>\n'
        if self.git_merge_request:
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-link"></span> Pull Request Link</td>\n'
            summary += '      <td><a href="TEMPLATE_PULL_REQUEST_LINK">TEMPLATE_PULL_REQUEST_LINK</a></td>\n'
            summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-header"></span> Pull Request Title</td>\n'
            summary += '      <td>TEMPLATE_PULL_REQUEST_TEMPLATE</td>\n'
            summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-out"></span> Source Branch</td>\n'
            summary += '      <td>' + self.git_src_branch + '</td>\n'
            summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Source Commit ID</td>\n'
            summary += '      <td>' + self.git_src_commit + '</td>\n'
            summary += '    </tr>\n'
            if (self.git_src_commit_msg is not None):
                summary += '    <tr>\n'
                summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Source Commit Message</td>\n'
                summary += '      <td>' + self.git_src_commit_msg + '</td>\n'
                summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-in"></span> Target Branch</td>\n'
            summary += '      <td>' + self.git_target_branch + '</td>\n'
            summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Target Commit ID</td>\n'
            summary += '      <td>' + self.git_target_commit + '</td>\n'
            summary += '    </tr>\n'
        else:
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tree-deciduous"></span> Branch</td>\n'
            summary += '      <td>' + self.git_src_branch + '</td>\n'
            summary += '    </tr>\n'
            summary += '    <tr>\n'
            summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Commit ID</td>\n'
            summary += '      <td>' + self.git_src_commit + '</td>\n'
            summary += '    </tr>\n'
            if (self.git_src_commit_msg is not None):
                summary += '    <tr>\n'
                summary += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Commit Message</td>\n'
                summary += '      <td>' + self.git_src_commit_msg + '</td>\n'
                summary += '    </tr>\n'
        summary += '  </table>\n'
        summary += '  <br>\n'
        return summary

    def generate_footer(self):
        """Append the HTML footer to report."""
        self.file.write('  <div class="well well-lg">End of Build Report -- Copyright <span class="glyphicon glyphicon-copyright-mark"></span> 2020 <a href="http://www.openairinterface.org/">OpenAirInterface</a>. All Rights Reserved.</div>\n')
        self.file.write('</div></body>\n')
        self.file.write('</html>\n')

    def add_build_summary_header(self):
        """Append Build Information Summary (Header)."""
        self.file.write('  <h2>Docker/Podman Images Build Summary</h2>\n')
        self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
        self.file.write('     <tr bgcolor="#33CCFF" >\n')
        self.file.write('       <th>Stage Name</th>\n')
        self.file.write('       <th>Image Kind</th>\n')
        cwd = os.getcwd()
        if os.path.isfile(cwd + '/archives/' + U18_BUILD_LOG_FILE):
            self.file.write('       <th>MAGMA - OAI MME cNF (Ubuntu-18)</th>\n')
        if os.path.isfile(cwd + '/archives/' + RHEL8_BUILD_LOG_FILE):
            self.file.write('       <th>MAGMA - OAI MME cNF (RHEL-8)</th>\n')
        self.file.write('     </tr>\n')

    def add_build_summary_footer(self):
        """Append Build Information Summary (Footer)."""
        self.file.write('  </table>\n')
        self.file.write('  <br>\n')

    def add_compile_rows(self):
        """Add rows for the compilation."""
        self.file.write('    <tr>\n')
        self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-common</b> Compile / Build</td>\n')
        self.analyze_build_log(COMMON_TYPE)
        self.file.write('    </tr>\n')
        self.file.write('    <tr>\n')
        self.analyze_compile_log(COMMON_TYPE)
        self.file.write('    </tr>\n')
        self.file.write('    <tr>\n')
        self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-oai-mme</b> Compile / Build</td>\n')
        self.analyze_build_log(MME_TYPE)
        self.file.write('    </tr>\n')
        self.file.write('    <tr>\n')
        self.analyze_compile_log(MME_TYPE)
        self.file.write('    </tr>\n')
        self.file.write('    <tr>\n')
        self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-sctpd</b> Compile / Build</td>\n')
        self.analyze_build_log(SCTPD_TYPE)
        self.file.write('    </tr>\n')
        self.file.write('    <tr>\n')
        self.analyze_compile_log(SCTPD_TYPE)
        self.file.write('    </tr>\n')

    def analyze_build_log(self, nf_type):
        """
        Add the row about build status.

        Args:
            nf_type: which build part
        """
        self.file.write('      <td>Builder Image</td>\n')
        cwd = os.getcwd()

        log_file_names = [U18_BUILD_LOG_FILE, RHEL8_BUILD_LOG_FILE]
        for log_file_name in log_file_names:
            if os.path.isfile(cwd + '/archives/' + log_file_name):
                status = False
                if nf_type == COMMON_TYPE:
                    section_start_pattern = 'ninja -C  /build/c/magma_common'
                    section_end_pattern = 'cmake  /magma/lte/gateway/c/core/oai -DCMAKE_BUILD_TYPE=Debug  -DS6A_OVER_GRPC=False -GNinja'
                if nf_type == MME_TYPE:
                    section_start_pattern = 'ninja -C  /build/c/core/oai'
                    section_end_pattern = 'cmake  /magma/orc8r/gateway/c/common -DCMAKE_BUILD_TYPE=Debug   -GNinja'
                if nf_type == SCTPD_TYPE:
                    section_start_pattern = 'ninja -C  /build/c/sctpd'
                    section_end_pattern = 'FROM ubuntu:bionic as magma-mme'
                section_status = False
                with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
                    for line in logfile:
                        my_res = re.search(section_start_pattern, line)
                        if my_res is not None:
                            section_status = True
                        my_res = re.search(section_end_pattern, line)
                        if my_res is not None:
                            section_status = False
                        if section_status:
                            if nf_type == COMMON_TYPE:
                                my_res = re.search('Linking CXX static library eventd/libEVENTD.a', line)
                            if nf_type == MME_TYPE:
                                my_res = re.search('Linking CXX executable core/oai_mme/mme', line)
                            if nf_type == SCTPD_TYPE:
                                my_res = re.search('Linking CXX executable sctpd', line)
                            if my_res is not None:
                                status = True
                    logfile.close()
                if status:
                    cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
                    cell_msg += 'OK:\n'
                else:
                    cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                    cell_msg += 'KO:\n'
                if nf_type == COMMON_TYPE:
                    cell_msg += ' -- ninja -C  /build/c/magma_common</b></pre></td>\n'
                if nf_type == MME_TYPE:
                    cell_msg += ' -- ninja -C  /build/c/core/oai</b></pre></td>\n'
                if nf_type == SCTPD_TYPE:
                    cell_msg += ' -- ninja -C  /build/c/sctpd</b></pre></td>\n'
            else:
                cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                cell_msg += 'KO: logfile (' + log_file_name + ') not found</b></pre></td>\n'

            self.file.write(cell_msg)

    def analyze_compile_log(self, nf_type):
        """
        Add the row about compilation errors/warnings/notes.

        Args:
            nf_type: which build part
        """
        self.file.write('      <td>Builder Image</td>\n')
        cwd = os.getcwd()

        log_file_names = [U18_BUILD_LOG_FILE, RHEL8_BUILD_LOG_FILE]
        for log_file_name in log_file_names:
            nb_errors = 0
            nb_warnings = 0
            nb_notes = 0
            if log_file_name.count('_rhel8') > 0:
                variant = 'RHEL8'
            else:
                variant = 'UBUNTU 18'
            self.errorWarningInfo.append([])
            self.variant.append(nf_type + ' ' + variant)
            idx = len(self.errorWarningInfo) - 1

            if os.path.isfile(cwd + '/archives/' + log_file_name):
                if nf_type == COMMON_TYPE:
                    section_start_pattern = '/build/c/magma_common'
                    section_end_pattern = 'mkdir -p  /build/c/core/oai'
                if nf_type == MME_TYPE:
                    section_start_pattern = '/build/c/core/oai'
                    section_end_pattern = 'mkdir -p  /build/c/magma_common'
                if nf_type == SCTPD_TYPE:
                    section_start_pattern = '/build/c/sctpd'
                    section_end_pattern = 'FROM ubuntu:bionic as magma-mme'
                section_status = False
                section_done = False
                with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
                    for line in logfile:
                        my_res = re.search(section_start_pattern, line)
                        if (my_res is not None) and not section_done and (re.search('cmake', line) is not None):
                            section_status = True
                        my_res = re.search(section_end_pattern, line)
                        if (my_res is not None) and not section_done and section_status:
                            section_status = False
                            section_done = True
                        if section_status:
                            my_res = re.search('error:', line)
                            if my_res is not None:
                                nb_errors += 1
                                errorandwarnings = {}
                                file_name = re.sub(':.*$', '', line.strip())
                                file_name = re.sub('^/magma/', '', file_name)
                                line_nb = '0'
                                warning_msg = re.sub('^.*error: ', '', line.strip())
                                details = re.search(':(?P<linenb>[0-9]+):', line)
                                if details is not None:
                                    line_nb = details.group('linenb')
                                errorandwarnings['kind'] = 'Error'
                                errorandwarnings['file_name'] = file_name
                                errorandwarnings['line_nb'] = line_nb
                                errorandwarnings['warning_msg'] = warning_msg
                                errorandwarnings['warning_type'] = 'fatal'
                                self.errorWarningInfo[idx].append(errorandwarnings)
                            my_res = re.search('warning:', line)
                            if my_res is not None:
                                nb_warnings += 1
                                errorandwarnings = {}
                                file_name = re.sub(':.*$', '', line.strip())
                                file_name = re.sub('^/magma/', '', file_name)
                                line_nb = '0'
                                details = re.search(':(?P<linenb>[0-9]+):', line)
                                if details is not None:
                                    line_nb = details.group('linenb')
                                warning_msg = re.sub('^.*warning: ', '', line.strip())
                                warning_msg = re.sub(' \[-W.*$', '', warning_msg)
                                warning_type = re.sub('^.* \[-W', '', line.strip())
                                warning_type = re.sub('\].*$', '', warning_type)
                                errorandwarnings['kind'] = 'Warning'
                                errorandwarnings['file_name'] = file_name
                                errorandwarnings['line_nb'] = line_nb
                                errorandwarnings['warning_msg'] = warning_msg
                                errorandwarnings['warning_type'] = warning_type
                                self.errorWarningInfo[idx].append(errorandwarnings)
                            my_res = re.search('note:', line)
                            if my_res is not None:
                                nb_notes += 1
                    logfile.close()
                if nb_warnings == 0 and nb_errors == 0:
                    cell_msg = '       <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
                elif nb_warnings < MAX_ALLOWED_WARNINGS and nb_errors == 0:
                    cell_msg = '       <td bgcolor="Orange"><pre style="border:none; background-color:Orange"><b>'
                else:
                    cell_msg = '       <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                if nb_errors > 0:
                    cell_msg += str(nb_errors) + ' errors found in compile log\n'
                cell_msg += str(nb_warnings) + ' warnings found in compile log\n'
                if nb_notes > 0:
                    cell_msg += str(nb_notes) + ' notes found in compile log\n'

                cell_msg += '</b></pre></td>\n'
            else:
                cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                cell_msg += 'KO: logfile (' + log_file_name + ') not found</b></pre></td>\n'

            self.file.write(cell_msg)

    def add_copy_to_target_image_row(self):
        """Add the row about start of target image creation."""
        self.file.write('    <tr>\n')
        self.file.write('      <td bgcolor="lightcyan" >SW libs Installation / Copy from Builder</td>\n')
        self.analyze_copy_log('MME')
        self.file.write('    </tr>\n')

    def analyze_copy_log(self, nf_type):
        """
        Add the row about copy of executables/packages to target image.

        Args:
            nf_type: which build part
        """
        if nf_type != 'MME':
            self.file.write('      <td>N/A</td>\n')
            self.file.write('      <td>Wrong NF Type for this Report</td>\n')
            return

        self.file.write('      <td>Target Image</td>\n')
        cwd = os.getcwd()

        log_file_names = [U18_BUILD_LOG_FILE, RHEL8_BUILD_LOG_FILE]
        for log_file_name in log_file_names:
            if os.path.isfile(cwd + '/archives/' + log_file_name):
                if log_file_name == U18_BUILD_LOG_FILE:
                    section_start_pattern = 'FROM ubuntu:bionic as magma-mme$'
                if log_file_name == RHEL8_BUILD_LOG_FILE:
                    section_start_pattern = 'FROM registry.access.redhat.com/ubi8/ubi:latest AS magma-mme$'
                section_end_pattern = 'WORKDIR /magma-mme/bin$'
                section_status = False
                status = False
                with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
                    for line in logfile:
                        my_res = re.search(section_start_pattern, line)
                        if my_res is not None:
                            section_status = True
                        my_res = re.search(section_end_pattern, line)
                        if (my_res is not None) and section_status:
                            section_status = False
                            status = True
                    logfile.close()
                if status:
                    cell_msg = '       <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
                    cell_msg += 'OK:\n'
                else:
                    cell_msg = '       <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                    cell_msg += 'KO:\n'
                cell_msg += '</b></pre></td>\n'
            else:
                cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                cell_msg += 'KO: logfile (' + log_file_name + ') not found</b></pre></td>\n'

            self.file.write(cell_msg)

    def add_copy_conf_tools_to_target_mage_row(self):
        """Add the row about copy of configuration/tools."""
        self.file.write('    <tr>\n')
        self.file.write('      <td bgcolor="lightcyan" >Copy Template Conf / Tools from Builder</td>\n')
        self.analyze_copy_conf_tool_log('MME')
        self.file.write('    </tr>\n')

    def analyze_copy_conf_tool_log(self, nf_type):
        """
        Retrieve info from log for conf/tools copy.

        Args:
            nf_type: which build part
        """
        if nf_type != 'MME':
            self.file.write('      <td>N/A</td>\n')
            self.file.write('      <td>Wrong NF Type for this Report</td>\n')
            return

        self.file.write('      <td>Target Image</td>\n')
        cwd = os.getcwd()

        log_file_names = [U18_BUILD_LOG_FILE, RHEL8_BUILD_LOG_FILE]
        for log_file_name in log_file_names:
            if os.path.isfile(cwd + '/archives/' + log_file_name):
                section_start_pattern = 'WORKDIR /magma-mme/bin$'
                if log_file_name == U18_BUILD_LOG_FILE:
                    section_end_pattern = 'Successfully tagged magma-mme:'
                if log_file_name == RHEL8_BUILD_LOG_FILE:
                    section_end_pattern = 'COMMIT magma-mme:'
                section_status = False
                status = False
                with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
                    for line in logfile:
                        my_res = re.search(section_start_pattern, line)
                        if my_res is not None:
                            section_status = True
                        my_res = re.search(section_end_pattern, line)
                        if (my_res is not None) and section_status:
                            section_status = False
                            status = True
                    logfile.close()
                if status:
                    cell_msg = '       <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
                    cell_msg += 'OK:\n'
                else:
                    cell_msg = '       <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                    cell_msg += 'KO:\n'
                cell_msg += '</b></pre></td>\n'
            else:
                cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                cell_msg += 'KO: logfile (' + log_file_name + ') not found</b></pre></td>\n'

            self.file.write(cell_msg)

    def add_image_size_row(self):
        """Add the row about image size of target image."""
        self.file.write('    <tr>\n')
        self.file.write('      <td bgcolor="lightcyan" >Image Size</td>\n')
        self.analyze_image_size_log('MME')
        self.file.write('    </tr>\n')

    def analyze_image_size_log(self, nf_type):
        """
        Retrieve image size from log.

        Args:
            nf_type: which build part
        """
        if nf_type != 'MME':
            self.file.write('      <td>N/A</td>\n')
            self.file.write('      <td>Wrong NF Type for this Report</td>\n')
            return

        self.file.write('      <td>Target Image</td>\n')
        cwd = os.getcwd()

        log_file_names = [U18_BUILD_LOG_FILE, RHEL8_BUILD_LOG_FILE]
        for log_file_name in log_file_names:
            if os.path.isfile(cwd + '/archives/' + log_file_name):
                if log_file_name == U18_BUILD_LOG_FILE:
                    section_start_pattern = 'Successfully tagged magma-mme'
                    section_end_pattern = 'MAGMA-OAI-MME DOCKER IMAGE BUILD'
                if log_file_name == RHEL8_BUILD_LOG_FILE:
                    section_start_pattern = 'COMMIT magma-mme:'
                    section_end_pattern = 'MAGMA-OAI-MME RHEL8 PODMAN IMAGE BUILD'
                section_status = False
                status = False
                with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
                    for line in logfile:
                        my_res = re.search(section_start_pattern, line)
                        if my_res is not None:
                            section_status = True
                        my_res = re.search(section_end_pattern, line)
                        if (my_res is not None) and section_status:
                            section_status = False
                        if section_status:
                            if self.git_merge_request:
                                my_res = re.search('magma-mme *ci-tmp', line)
                            else:
                                my_res = re.search('magma-mme *master *', line)
                            if my_res is not None:
                                my_res = re.search('ago *([0-9 A-Z]+)', line)
                                if my_res is not None:
                                    size = my_res.group(1)
                                    status = True
                    logfile.close()
                if status:
                    cell_msg = '       <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
                    cell_msg += 'OK:  ' + size + '\n'
                else:
                    cell_msg = '       <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                    cell_msg += 'KO:\n'
                cell_msg += '</b></pre></td>\n'
            else:
                cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
                cell_msg += 'KO: logfile (' + log_file_name + ') not found</b></pre></td>\n'

            self.file.write(cell_msg)

    def add_details(self):
        """Add the compilation warnings/errors details"""
        idx = 0
        needed_details = False
        while (idx < len(self.errorWarningInfo)):
            if len(self.errorWarningInfo[idx]) > 0:
                needed_details = True
            idx += 1
        if not needed_details:
            return

        details = '  <h3>Details</h3>\n'
        details += '  <button data-toggle="collapse" data-target="#compilation-details">Details for Compilation Errors and Warnings </button>\n'
        details += '  <div id="compilation-details" class="collapse">\n'
        idx = 0
        while (idx < len(self.errorWarningInfo)):
            if len(self.errorWarningInfo[idx]) == 0:
                idx += 1
                continue
            details += '  <h4>Details for ' + self.variant[idx] + '</h4>\n'
            details += '   <table class="table-bordered" width = "100%" align = "center" border = "1">\n'
            details += '      <tr bgcolor = "#33CCFF" >\n'
            details += '        <th>File</th>\n'
            details += '        <th>Line Number</th>\n'
            details += '        <th>Status</th>\n'
            details += '        <th>Kind</th>\n'
            details += '        <th>Message</th>\n'
            details += '      </tr>\n'
            for info in self.errorWarningInfo[idx]:
                details += '      <tr>\n'
                details += '        <td>' + info['file_name'] + '</td>\n'
                details += '        <td>' + info['line_nb'] + '</td>\n'
                details += '        <td>' + info['kind'] + '</td>\n'
                details += '        <td>' + info['warning_type'] + '</td>\n'
                details += '        <td>' + info['warning_msg'] + '</td>\n'
                details += '      </tr>\n'
            details += '   </table>\n'
            idx += 1
        details += '  </div>\n'
        details += '  <br>\n'
        self.file.write(details)

    def append_build_summary(self, mode):
        """
        Append in test results a correct build info summary.

        Args:
            mode: which test mode
        """
        cwd = os.getcwd()
        if mode == 'dsTester':
            filename = 'test_results_magma_oai_epc.html'
        if os.path.isfile(cwd + '/' + filename):
            new_test_report = open(cwd + '/new_' + filename, 'w')
            build_summary_to_be_done = True
            with open(cwd + '/' + filename, 'r') as original_test_report:
                for line in original_test_report:
                    my_res = re.search('Deployment Summary', line)
                    if (my_res is not None) and build_summary_to_be_done:
                        summary = self.generate_build_summary()
                        new_test_report.write(summary)
                        build_summary_to_be_done = False
                    new_test_report.write(line)
                original_test_report.close()
            new_test_report.close()
            os.rename(cwd + '/new_' + filename, cwd + '/' + filename)

# --------------------------------------------------------------------------------------------------------
#
# Start of main
#
# --------------------------------------------------------------------------------------------------------


argvs = sys.argv
argc = len(argvs)

HTML = HtmlReport()

while len(argvs) > 1:
    my_argv = argvs.pop(1)
    if re.match('^--help$', my_argv, re.IGNORECASE):
        print('No help yet.')
        sys.exit(0)
    elif re.match('^--job_name=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--job_name=(.+)$', my_argv, re.IGNORECASE)
        HTML.job_name = match.group(1)
    elif re.match('^--job_id=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--job_id=(.+)$', my_argv, re.IGNORECASE)
        HTML.job_id = match.group(1)
    elif re.match('^--job_url=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--job_url=(.+)$', my_argv, re.IGNORECASE)
        HTML.job_url = match.group(1)
    elif re.match('^--git_url=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--git_url=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_url = match.group(1)
    elif re.match('^--git_src_branch=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--git_src_branch=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_src_branch = match.group(1)
    elif re.match('^--git_src_commit=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--git_src_commit=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_src_commit = match.group(1)
    elif re.match('^--git_src_commit_msg=(.+)$', my_argv, re.IGNORECASE):
        # Not Mandatory
        match = re.match('^--git_src_commit_msg=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_src_commit_msg = match.group(1)
    elif re.match('^--git_merge_request=(.+)$', my_argv, re.IGNORECASE):
        # Can be silent: would be false!
        match = re.match('^--git_merge_request=(.+)$', my_argv, re.IGNORECASE)
        if match.group(1) == 'true' or match.group(1) == 'True':
            HTML.git_merge_request = True
    elif re.match('^--git_target_branch=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--git_target_branch=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_target_branch = match.group(1)
    elif re.match('^--git_target_commit=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--git_target_commit=(.+)$', my_argv, re.IGNORECASE)
        HTML.git_target_commit = match.group(1)
    elif re.match('^--mode=(.+)$', my_argv, re.IGNORECASE):
        match = re.match('^--mode=(.+)$', my_argv, re.IGNORECASE)
        if match.group(1) == 'Build':
            HTML.mode = 'build'
        elif match.group(1) == 'TestWithDsTest':
            HTML.mode = 'dsTester'
        else:
            sys.exit('Invalid mode: ' + match.group(1))
    else:
        sys.exit('Invalid Parameter: ' + my_argv)

if HTML.job_name == '' or HTML.job_id == '' or HTML.job_url == '' or HTML.mode == '':
    sys.exit('Missing Parameter in job description')

if HTML.git_url == '' or HTML.git_src_branch == '' or HTML.git_src_commit == '':
    sys.exit('Missing Parameter in Git Repository description')

if HTML.git_merge_request:
    if HTML.git_target_commit == '' or HTML.git_target_branch == '':
        sys.exit('Missing Parameter in Git Pull Request Repository description')

if HTML.mode == 'build':
    HTML.generate_build_report()
elif HTML.mode == 'dsTester':
    HTML.append_build_summary(HTML.mode)
