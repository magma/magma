# /*
# * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# * contributor license agreements.  See the NOTICE file distributed with
# * this work for additional information regarding copyright ownership.
# * The OpenAirInterface Software Alliance licenses this file to You under
# * the terms found in the LICENSE file in the root of this
# * source tree.
# *
# * Unless required by applicable law or agreed to in writing, software
# * distributed under the License is distributed on an "AS IS" BASIS,
# * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# * See the License for the specific language governing permissions and
# * limitations under the License.
# *-------------------------------------------------------------------------------
# * For more information about the OpenAirInterface (OAI) Software Alliance:
# *   contact@openairinterface.org
# */
# ---------------------------------------------------------------------

import os
import re
import subprocess
import sys


class HtmlReport():
  def __init__(self):
    self.job_name = ''
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

  def generate(self):
    cwd = os.getcwd()
    self.file = open(cwd + '/test_results_magma_converged_mme.html', 'w')
    self.generateHeader()

    self.coding_formatting_log()

    self.analyze_sca_log()

    self.addPagination()

    # no-S11 part
    self.buildSummaryHeader('agw1-no-s11')
    self.makeRunRow('agw1-no-s11')
    self.statusCheckRow('agw1-no-s11')
    self.buildSummaryFooter()

    self.testSummaryHeader('agw1-no-s11')
    self.s1apTesterTable('agw1-no-s11')
    self.testSummaryFooter()

    # with-S11 part
    self.buildSummaryHeader('agw1-with-s11')
    self.makeRunRow('agw1-with-s11')
    self.statusCheckRow('agw1-with-s11')
    self.buildSummaryFooter()

    self.testSummaryHeader('agw1-with-s11')
    self.s1apTesterTable('agw1-with-s11')
    self.testSummaryFooter()

    self.generateFooter()
    self.file.close()

  def generateHeader(self):
    # HTML Header
    self.file.write('<!DOCTYPE html>\n')
    self.file.write('<html class="no-js" lang="en-US">\n')
    self.file.write('<head>\n')
    self.file.write('  <meta name="viewport" content="width=device-width, initial-scale=1">\n')
    self.file.write('  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">\n')
    self.file.write('  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>\n')
    self.file.write('  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>\n')
    self.file.write('  <title>MAGMA/OAI Core Network Test Results for ' + self.job_name + ' job build #' + self.job_id + '</title>\n')
    self.file.write('</head>\n')
    self.file.write('<body><div class="container">\n')
    self.file.write('  <table width = "100%" style="border-collapse: collapse; border: none;">\n')
    self.file.write('   <tr style="border-collapse: collapse; border: none;">\n')
    self.file.write('     <td style="border-collapse: collapse; border: none;">\n')
    self.file.write('       <a href="http://www.openairinterface.org/">\n')
    self.file.write('          <img src="http://www.openairinterface.org/wp-content/uploads/2016/03/cropped-oai_final_logo2.png" alt="" border="none" height=50 width=150>\n')
    self.file.write('          </img>\n')
    self.file.write('       </a>\n')
    self.file.write('     </td>\n')
    self.file.write('     <td style="border-collapse: collapse; border: none; vertical-align: center;">\n')
    self.file.write('       <b><font size = "6">Job Summary -- Job: ' + self.job_name + ' -- Build-ID: <a href="' + self.job_url + '">' + self.job_id + '</a></font></b>\n')
    self.file.write('     </td>\n')
    self.file.write('   </tr>\n')
    self.file.write('  </table>\n')
    self.file.write('  <br>\n')

    # Build Info Summary
    self.file.write('  <table class="table-bordered" width = "80%" align = "center" border = "1">\n')
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-time"></span> Build Start Time</td>\n')
    #date_formatted = re.sub('\..*', '', self.created)
    self.file.write('      <td>' + self.job_start_time + '</td>\n')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-wrench"></span> Build Trigger</td>\n')
    if self.git_merge_request:
      self.file.write('      <td>Merge Request</td>\n')
    else:
      self.file.write('      <td>Push Event</td>\n')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-cloud-upload"></span> GIT Repository</td>\n')
    self.file.write('      <td><a href="' + self.git_url + '">' + self.git_url + '</a></td>\n')
    self.file.write('    </tr>\n')
    if self.git_merge_request:
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-out"></span> Source Branch</td>\n')
      self.file.write('      <td>' + self.git_src_branch + '</td>\n')
      self.file.write('    </tr>\n')
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Source Commit ID</td>\n')
      self.file.write('      <td>' + self.git_src_commit + '</td>\n')
      self.file.write('    </tr>\n')
      if (self.git_src_commit_msg is not None):
        self.file.write('    <tr>\n')
        self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Source Commit Message</td>\n')
        self.file.write('      <td>' + self.git_src_commit_msg + '</td>\n')
        self.file.write('    </tr>\n')
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-in"></span> Target Branch</td>\n')
      self.file.write('      <td>' + self.git_target_branch + '</td>\n')
      self.file.write('    </tr>\n')
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Target Commit ID</td>\n')
      self.file.write('      <td>' + self.git_target_commit + '</td>\n')
      self.file.write('    </tr>\n')
    else:
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tree-deciduous"></span> Branch</td>\n')
      self.file.write('      <td>' + self.git_src_branch + '</td>\n')
      self.file.write('    </tr>\n')
      self.file.write('    <tr>\n')
      self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Commit ID</td>\n')
      self.file.write('      <td>' + self.git_src_commit + '</td>\n')
      self.file.write('    </tr>\n')
      if (self.git_src_commit_msg is not None):
        self.file.write('    <tr>\n')
        self.file.write('      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Commit Message</td>\n')
        self.file.write('      <td>' + self.git_src_commit_msg + '</td>\n')
        self.file.write('    </tr>\n')
    self.file.write('  </table>\n')
    self.file.write('  <br>\n')

  def generateFooter(self):
    self.file.write('  </nav>\n')
    self.file.write('  <div class="well well-lg">End of Build Report -- Copyright <span class="glyphicon glyphicon-copyright-mark"></span> 2020 <a href="http://www.openairinterface.org/">OpenAirInterface</a>. All Rights Reserved.</div>\n')
    self.file.write('</div></body>\n')
    self.file.write('</html>\n')

  def coding_formatting_log(self):
    cwd = os.getcwd()
    self.file.write('  <h2>OAI Coding / Formatting Guidelines Check</h2>\n')
    if os.path.isfile(cwd + '/oai_rules_result.txt'):
      cmd = 'grep NB_FILES_FAILING_CHECK ' + cwd + '/oai_rules_result.txt | sed -e "s#NB_FILES_FAILING_CHECK=##"'
      nb_fail = subprocess.check_output(cmd, shell=True, universal_newlines=True)
      cmd = 'grep NB_FILES_CHECKED ' + cwd + '/oai_rules_result.txt | sed -e "s#NB_FILES_CHECKED=##"'
      nb_total = subprocess.check_output(cmd, shell=True, universal_newlines=True)
      if int(nb_fail.strip()) == 0:
        self.file.write('  <div class="alert alert-success">\n')
        if self.git_merge_request:
          self.file.write('    <strong>All modified files in Merge-Request follow OAI rules. <span class="glyphicon glyphicon-ok-circle"></span> -> (' + nb_total.strip() + ' were checked)</strong>\n')
        else:
          self.file.write('    <strong>All files in repository follow OAI rules. <span class="glyphicon glyphicon-ok-circle"></span> -> (' + nb_total.strip() + ' were checked)</strong>\n')
        self.file.write('  </div>\n')
      else:
        self.file.write('  <div class="alert alert-warning">\n')
        if self.git_merge_request:
          self.file.write('    <strong>' + nb_fail.strip() + ' modified files in Merge-Request DO NOT follow OAI rules. <span class="glyphicon glyphicon-warning-sign"></span> -> (' + nb_total.strip() + ' were checked)</strong>\n')
        else:
          self.file.write('    <strong>' + nb_fail.strip() + ' files in repository DO NOT follow OAI rules. <span class="glyphicon glyphicon-warning-sign"></span> -> (' + nb_total.strip() + ' were checked)</strong>\n')
        self.file.write('  </div>\n')

        if os.path.isfile(cwd + '/oai_rules_result_list.txt'):
          self.file.write('  <button data-toggle="collapse" data-target="#oai-formatting-details">More details on formatting check</button>\n')
          self.file.write('  <div id="oai-formatting-details" class="collapse">\n')
          self.file.write('  <p>Please apply the following command to this(ese) file(s): </p>\n')
          self.file.write('  <p style="margin-left: 30px"><strong><code>clang-format -i filename(s)</code></strong></p>\n')
          self.file.write('  <table class="table-bordered" width = "60%" align = "center" border = 1>\n')
          self.file.write('    <tr><th bgcolor = "lightcyan" >Filename</th></tr>\n')
          with open(cwd + '/oai_rules_result_list.txt', 'r') as filelist:
            for line in filelist:
              self.file.write('    <tr><td>' + line.strip() + '</td></tr>\n')
            filelist.close()
          self.file.write('  </table>\n')
          self.file.write('  </div>\n')
    else:
      self.file.write('  <div class="alert alert-danger">\n')
      self.file.write('   <strong>Was NOT performed (with CLANG-FORMAT tool). <span class="glyphicon glyphicon-ban-circle"></span></strong>\n')
      self.file.write('  </div>\n')

    self.file.write('  <br>\n')

  def analyze_sca_log(self):
    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/cppcheck_build.log'):
      self.file.write('  <h2>Static Code Analysis</h2>\n')
    if os.path.isfile(cwd + '/archives/cppcheck.xml'):
      nb_errors = 0
      nb_warnings = 0
      nb_uninitvar = 0
      nb_uninitStructMember = 0
      nb_memleak = 0
      nb_doubleFree = 0
      nb_resourceLeak = 0
      nb_nullPointer = 0
      nb_arrayIndexOutOfBounds = 0
      nb_bufferAccessOutOfBounds = 0
      nb_unknownEvaluationOrder = 0
      with open(cwd + '/archives/cppcheck.xml', 'r') as xmlfile:
        for line in xmlfile:
          result = re.search('severity="warning"', line)
          if result is not None:
            nb_warnings += 1
          result = re.search('severity="error"', line)
          if result is not None:
            nb_errors += 1
            result = re.search('uninitvar', line)
            if result is not None:
              nb_uninitvar += 1
            result = re.search('uninitStructMember', line)
            if result is not None:
              nb_uninitStructMember += 1
            result = re.search('memleak', line)
            if result is not None:
              nb_memleak += 1
            result = re.search('doubleFree', line)
            if result is not None:
              nb_doubleFree += 1
            result = re.search('resourceLeak', line)
            if result is not None:
              nb_resourceLeak += 1
            result = re.search('nullPointer', line)
            if result is not None:
              nb_nullPointer += 1
            result = re.search('arrayIndexOutOfBounds', line)
            if result is not None:
              nb_arrayIndexOutOfBounds += 1
            result = re.search('bufferAccessOutOfBounds', line)
            if result is not None:
              nb_bufferAccessOutOfBounds += 1
            result = re.search('unknownEvaluationOrder', line)
            if result is not None:
              nb_unknownEvaluationOrder += 1
        xmlfile.close()
      if (nb_errors == 0) and (nb_warnings == 0):
        self.file.write('   <div class="alert alert-success">\n')
        self.file.write('   <strong>CPPCHECK found NO error and NO warning <span class="glyphicon glyphicon-ok-circle"></span></strong>\n')
        self.file.write('   </div>\n')
      elif (nb_errors == 0):
        self.file.write('   <div class="alert alert-warning">\n')
        self.file.write('   <strong>CPPCHECK found NO error and ' + str(nb_warnings) + ' warnings <span class="glyphicon glyphicon-warning-sign"></span></strong>\n')
        self.file.write('   </div>\n')
      else:
        self.file.write('   <div class="alert alert-danger">\n')
        self.file.write('   <strong>CPPCHECK found ' + str(nb_errors) + ' errors and ' + str(nb_warnings) + ' warnings <span class="glyphicon glyphicon-ban-circle"></span></strong>\n')
        self.file.write('   </div>\n')
      if (nb_errors > 0) or (nb_warnings > 0):
        self.file.write('   <button data-toggle="collapse" data-target="#oai-cppcheck-details">More details on CPPCHECK results</button>\n')
        self.file.write('   <div id="oai-cppcheck-details" class="collapse">\n')
        self.file.write('   <br>\n')
        self.file.write('   <table class="table-bordered" width = "80%" align = "center" border = "1">\n')
        self.file.write('   <tr bgcolor = "#33CCFF" >\n')
        self.file.write('   <th>Error / Warning Type</th>\n')
        self.file.write('   <th>Nb Errors</th>\n')
        self.file.write('   <th>Nb Warnings</th>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Uninitialized variable</td>\n')
        self.file.write('   <td>' + str(nb_uninitvar) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Uninitialized struct member</td>\n')
        self.file.write('   <td>' + str(nb_uninitStructMember) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Memory leak</td>\n')
        self.file.write('   <td>' + str(nb_memleak) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Memory is freed twice</td>\n')
        self.file.write('   <td>' + str(nb_doubleFree) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Resource leak</td>\n')
        self.file.write('   <td>' + str(nb_resourceLeak) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Possible null pointer dereference</td>\n')
        self.file.write('   <td>' + str(nb_nullPointer) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Array access  out of bounds</td>\n')
        self.file.write('   <td>' + str(nb_arrayIndexOutOfBounds) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Buffer is accessed out of bounds</td>\n')
        self.file.write('   <td>' + str(nb_bufferAccessOutOfBounds) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Expression depends on order of evaluation of side effects</td>\n')
        self.file.write('   <td>' + str(nb_unknownEvaluationOrder) + '</td>\n')
        self.file.write('   <td>N/A</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr>\n')
        self.file.write('   <td>Others</td>\n')
        nb_others = nb_uninitvar + nb_uninitStructMember + nb_memleak + nb_doubleFree + nb_resourceLeak + nb_nullPointer + nb_arrayIndexOutOfBounds + nb_arrayIndexOutOfBounds + nb_bufferAccessOutOfBounds + nb_unknownEvaluationOrder
        nb_others = nb_errors - nb_others
        self.file.write('   <td>' + str(nb_others) + '</td>\n')
        self.file.write('   <td>' + str(nb_warnings) + '</td>\n')
        self.file.write('   </tr>\n')
        self.file.write('   <tr bgcolor = "#33CCFF" >\n')
        self.file.write('   <th>Total</th>\n')
        self.file.write('   <th>' + str(nb_errors) + '</th>\n')
        self.file.write('   <th>' + str(nb_warnings) + '</th>\n')
        self.file.write('   </tr>\n')
        self.file.write('   </table>\n')
        self.file.write('   <br>\n')
        self.file.write('   <p>Full details in artifact (cppcheck.xml) </p>\n')
        self.file.write('   <p style="margin-left: 30px">Graphical Interface tool : <strong><code>cppcheck-gui -l cppcheck.xml</code></strong></p>\n')
        self.file.write('   <br>\n')
        self.file.write('   </div>\n')
    else:
      self.file.write('  <div class="alert alert-danger">\n')
      self.file.write('   <strong>Was NOT performed (with CPPCHECK tool). <span class="glyphicon glyphicon-ban-circle"></span></strong>\n')
      self.file.write('  </div>\n')

  def addPagination(self):
    self.file.write('\n')
    self.file.write('  <nav aria-label="Page nav">\n')
    self.file.write('  <ul class="pagination pagination-lg">\n')
    # Check status for AGW1-noS11
    agw1_noS11_status = False
    cwd = os.getcwd()
    s1ap_logfile = 'magma_run_s1ap_tester.log'
    if os.path.isfile(cwd + '/archives/' + s1ap_logfile):
      cmd = 'egrep -c "AGW-VM-S1AP-TESTS: OK" ' + cwd + '/archives/' + s1ap_logfile
      try:
        found = subprocess.check_output(cmd, shell=True, universal_newlines=True)
      except:
        found = '0'
      if int(found) > 0:
        agw1_noS11_status = True
    if agw1_noS11_status:
      self.file.write('     <li class="page-item"><a class="page-link" href="#agw1-no-s11-details">AGW1-no-S11 <font color = "Green"><span class="glyphicon glyphicon-ok-sign"></span></font></a></li>\n')
    else:
      self.file.write('     <li class="page-item"><a class="page-link" href="#agw1-no-s11-details">AGW1-no-S11 <font color = "Red"><span class="glyphicon glyphicon-remove-sign"></span></font></a></li>\n')

    # Check status for AGW1-with-S11
    agw1_wiS11_status = False
    s1ap_logfile = 'magma_run_s1ap_tester_s11.log'
    if os.path.isfile(cwd + '/archives/' + s1ap_logfile):
      cmd = 'egrep -c "AGW-VM-S1AP-TESTS: OK" ' + cwd + '/archives/' + s1ap_logfile
      try:
        found = subprocess.check_output(cmd, shell=True, universal_newlines=True)
      except:
        found = '0'
      if int(found) > 0:
        agw1_wiS11_status = True
    if agw1_wiS11_status:
      self.file.write('     <li class="page-item"><a class="page-link" href="#agw1-wi-s11-details">AGW1-with-S11 <font color = "Green"><span class="glyphicon glyphicon-ok-sign"></span></font></a></li>\n')
    else:
      self.file.write('     <li class="page-item"><a class="page-link" href="#agw1-wi-s11-details">AGW1-with-S11 <font color = "Red"><span class="glyphicon glyphicon-remove-sign"></span></font></a></li>\n')
    self.file.write('  </ul>\n')
    self.file.write('\n')

  def buildSummaryHeader(self, kind):
    if kind == 'agw1-no-s11':
      self.file.write('  <div id="agw1-no-s11-details" class="container">\n')
      self.file.write('  <h1>AGW1 NO S11</h1>\n')
      self.file.write('  <h2>AGW1 no S11 -- Build Summary</h2>\n')
    if kind == 'agw1-with-s11':
      self.file.write('  <div id="agw1-wi-s11-details" class="container">\n')
      self.file.write('  <h1>AGW1 WITH S11</h1>\n')
      self.file.write('  <h2>AGW1 with S11 -- Build Summary</h2>\n')

    self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
    self.file.write('     <tr bgcolor="#33CCFF" >\n')
    self.file.write('       <th>Stage Name</th>\n')
    self.file.write('       <th>Magma-dev VM</th>\n')
    self.file.write('       <th>Magma-test VM</th>\n')
    self.file.write('       <th>Magma-traffic VM</th>\n')

  def buildSummaryFooter(self):
    self.file.write('  </table>\n')
    self.file.write('  <br>\n')

  def analyze_vagrant_up_log(self, vmType):
    logFileName = vmType.lower().replace('magma', 'magma_vagrant') + '_up.log'
    if vmType == 'magma':
      vmName = 'magma_dev'
    elif vmType == 'magma_trfserver':
      vmName = 'magma-trfserver'
    else:
      vmName = 'magma_test'
    pattern = 'Machine booted and ready!'
    mount_pattern = 'Mounting shared folders'

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      status = False
      mount = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(pattern, line)
          if result is not None:
            status = True
          result = re.search(mount_pattern, line)
          if result is not None:
            mount = True
        logfile.close()
      if status and mount:
        cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
        cell_msg += 'OK: VM ' + vmName + ':\n'
        cell_msg += ' -- woken-up         successfully\n'
        cell_msg += ' -- folders mounted  successfully</b></pre></td>\n'
      elif not status:
        cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
        cell_msg += 'KO: VM ' + vmName + ':\n'
        cell_msg += ' -- did not start properly?</b></pre></td>\n'
      else:
        cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
        cell_msg += 'KO: VM ' + vmName + ':\n'
        cell_msg += ' -- woken-up         successfully\n'
        cell_msg += ' -- folders DID NOT MOUNT</b></pre></td>\n'
    else:
      cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def makeRunRow(self, kind):
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" >Make Run Results</td>\n')
    self.analyze_vagrant_make_run_log('magma', kind)
    self.analyze_vagrant_make_run_log('magma_test', kind)
    self.analyze_vagrant_make_run_log('magma_trfserver', kind)
    self.file.write('    </tr>\n')

  def analyze_vagrant_make_run_log(self, vmType, kind):
    if vmType == 'magma_test' or vmType == 'magma_trfserver':
      cell_msg = '      <td bgcolor="LightGray"><pre style="border:none; background-color:LightGray"><b>'
      cell_msg += 'N/A</b></pre></td>\n'
      self.file.write(cell_msg)
      return

    if kind == 'agw1-no-s11':
      logFileName = vmType.lower().replace('magma', 'magma_vagrant') + '_make_run.log'
    if kind == 'agw1-with-s11':
      logFileName = vmType.lower().replace('magma', 'magma_vagrant') + '_make_run2.log'
    module_pattern = 'ninja -C  '
    end_pattern = 'sudo service magma@magmad start'

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      status = False
      firstModulePassed = False
      moduleList = ''
      module_name = ''
      firstLineIsOptions = True
      optionLine = ''
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          if firstLineIsOptions:
            firstLineIsOptions = False
            optionLine = line
          result = re.search(end_pattern, line)
          if result is not None:
            status = True
          result = re.search(module_pattern, line)
          if result is not None:
            if firstModulePassed:
              moduleList += '  -- ' + module_name + ': OK\n'
            firstModulePassed = True
            module_name = line.replace('\n', '').replace('ninja -C  /home/vagrant/build/c//', '')
        logfile.close()
      if status:
        moduleList += '  -- ' + module_name + ': OK\n'
        cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
        cell_msg += 'Build Options are: \n'
        cell_msg += optionLine + '\n'
        cell_msg += 'OK: \n'
      else:
        moduleList += '  -- ' + module_name + ': KO\n'
        cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
        cell_msg += 'Build Options are: \n'
        cell_msg += optionLine + '\n'
        cell_msg += 'KO: \n'
      cell_msg += moduleList
      cell_msg += '</b></pre></td>\n'
    else:
      cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def statusCheckRow(self, kind):
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" >Magma Services Status</td>\n')
    self.analyze_status_check_log('magma', kind)
    self.analyze_status_check_log('magma_test', kind)
    self.analyze_status_check_log('magma_trfserver', kind)
    self.file.write('    </tr>\n')

  def analyze_status_check_log(self, vmType, kind):
    if vmType == 'magma_test' or vmType == 'magma_trfserver':
      cell_msg = '      <td bgcolor="LightGray"><pre style="border:none; background-color:LightGray"><b>'
      cell_msg += 'N/A</b></pre></td>\n'
      self.file.write(cell_msg)
      return

    if kind == 'agw1-no-s11':
      logFileName = 'magma_status.log'
    if kind == 'agw1-with-s11':
      logFileName = 'magma_status2.log'
    pattern_module = '● magma@'
    active_pattern = 'Active: active'

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      status = True
      service_status = False
      firstService = False
      listServices = ''
      service_name = ''
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(pattern_module, line)
          if result is not None:
            if firstService and not service_status:
              listServices += ' -- ' + service_name + '  NOT ACTIVE\n'
              status = False
            service_name = line.replace('\n', '').replace(pattern_module, '')
            service_name = service_name.replace('.service', '')
            firstService = True
            service_status = False
          result = re.search(active_pattern, line)
          if result is not None:
            if not service_status:
              listServices += ' -- ' + service_name + '  ACTIVE\n'
              service_status = True
        logfile.close()
      if status:
        cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
        cell_msg += 'OK: \n'
      else:
        cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
        cell_msg += 'KO: \n'
      cell_msg += listServices
      cell_msg += '</b></pre></td>\n'
    else:
      cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def testSummaryHeader(self, kind):
    if kind == 'agw1-no-s11':
      self.file.write('  <h2>AGW1 no S11 -- Test Summary</h2>\n')
    if kind == 'agw1-with-s11':
      self.file.write('  <h2>AGW1 with S11 -- Test Summary</h2>\n')

  def testSummaryFooter(self):
    self.file.write('  <br>\n')
    self.file.write('  </div>\n')

  def s1apTesterTable(self, kind):
    cwd = os.getcwd()
    if kind == 'agw1-no-s11':
      s1ap_logfile = 'magma_run_s1ap_tester.log'
      button_ref = 's1ap-tester'
    if kind == 'agw1-with-s11':
      s1ap_logfile = 'magma_run_s1ap_tester_s11.log'
      button_ref = 's1ap-tester-s11'

    if os.path.isfile(cwd + '/archives/' + s1ap_logfile):
      cmd = 'egrep -c "AGW-VM-S1AP-TESTS: OK" ' + cwd + '/archives/' + s1ap_logfile
      try:
        found = subprocess.check_output(cmd, shell=True, universal_newlines=True)
      except:
        found = '0'
      if int(found) > 0:
        self.file.write('  <h3>S1AP Tester Summary -- STATUS is <font color = "Green"><span class="glyphicon glyphicon-ok-sign"></span></font></h3>\n')
      else:
        self.file.write('  <h3>S1AP Tester Summary -- STATUS is <font color = "Red"><span class="glyphicon glyphicon-remove-sign"></span></font></h3>\n')
    else:
      self.file.write('  <h3>S1AP Tester Summary -- STATUS is <font color = "Gray"><span class="glyphicon glyphicon-question-sign"></span></font></h3>\n')

    self.file.write('  <br>\n')
    self.file.write('  <button data-toggle="collapse" data-target="#' + button_ref + '">More details on S1AP Tester checks</button>\n')
    self.file.write('  <div id="' + button_ref + '" class="collapse">\n')
    self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
    self.file.write('     <tr bgcolor="#33CCFF" >\n')
    self.file.write('       <th>Test Name</th>\n')
    self.file.write('       <th>Test Status</th>\n')
    self.file.write('       <th>Test Duration</th>\n')
    self.file.write('     </tr>\n')

    if os.path.isfile(cwd + '/archives/' + s1ap_logfile):
      pattern = 'echo "Running test: '
      status = True
      global_test_duration = 0.0
      with open(cwd + '/archives/' + s1ap_logfile, 'r') as logfile:
        for line in logfile:
          result = re.search(pattern + '(.+)$', line)
          if result is not None:
            test_name = result.group(1).replace('"', '')
            test_duration = 'unknown'
          result = re.search('Ran 1 test in (.+)$', line)
          if result is not None:
            test_duration = result.group(1)
            global_test_duration += float(test_duration.replace('s', ''))
          result = re.search('^OK$|^FAILED|^Killed by signal 15', line)
          if result is not None:
            listTests = '    <tr>\n'
            listTests += '      <td bgcolor = "LightGray" >' + test_name + '</td>\n'
            result = re.search('^OK$', line)
            if result is not None:
              listTests += '      <td bgcolor = "Green" ><font color = "white"><b>OK</b></font></td>\n'
            else:
              listTests += '      <td bgcolor = "Red"><font color = "white"><b>FAIL</b></font></td>\n'
              status = False
            listTests += '      <td bgcolor = "LightGray" >' + test_duration + '</td>\n'
            listTests += '    </tr>\n'
            self.file.write(listTests)
      logfile.close()
      self.file.write('     <tr>\n')
      self.file.write('       <td bgcolor = "LightBlue" ><b>FINAL STATUS</b></td>\n')
      if status:
        self.file.write('       <td align = "left" bgcolor = "Green"><b><font color = "white">OK</font></b></td>\n')
        self.file.write('       <td align = "left" bgcolor = "Green"><b><font color = "white">' + str(global_test_duration) + 's</font></b></td>\n')
      else:
        self.file.write('       <td align = "left" bgcolor = "Red"><b><font color = "white">FAIL</font></b></td>\n')
        self.file.write('       <td align = "left" bgcolor = "Red"><b><font color = "white">' + str(global_test_duration) + 's</font></b></td>\n')
      self.file.write('     </tr>\n')
    else:
      self.file.write('     <tr bgcolor="Tomato" >\n')
      self.file.write('       <td colspan=3><b>KO: logfile (' + s1ap_logfile + ') not found</b></td>\n')
      self.file.write('     </tr>\n')

    if kind == 'agw1-no-s11' and os.path.isfile(cwd + '/archives/magma_vagrant_make_coverage_oai.log'):
      with open(cwd + '/archives/magma_vagrant_make_coverage_oai.log', 'r') as logfile:
        cov_global_rate = False
        for line in logfile:
          result = re.search('Overall coverage rate:', line)
          if result is not None:
            cov_global_rate = True
            self.file.write('     <tr>\n')
            self.file.write('       <td colspan = 3 bgcolor = "DarkSlateGrey"><b><font color = "white">Overall coverage rate</font></b></td>\n')
            self.file.write('     </tr>\n')
          if cov_global_rate:
            result = re.search(' lines.*: (.+)% \(([0-9]+) of ([0-9]+) lines', line)
            if result is not None:
              percentage = result.group(1)
              reached_line_nb = result.group(2)
              total_line_nb = result.group(3)
              self.file.write('     <tr>\n')
              self.file.write('       <td><b>Lines</b></td>\n')
              self.file.write('       <td bgcolor = "DarkSlateGrey"><b><font color = "white">' + percentage + '%</font></b></td>\n')
              self.file.write('       <td bgcolor = "DarkSlateGrey"><b><font color = "white">' + reached_line_nb + ' over ' + total_line_nb + '</font></b></td>\n')
              self.file.write('     </tr>\n')
            result = re.search(' functions.*: (.+)% \(([0-9]+) of ([0-9]+) functions', line)
            if result is not None:
              percentage = result.group(1)
              reached_function_nb = result.group(2)
              total_function_nb = result.group(3)
              self.file.write('     <tr>\n')
              self.file.write('       <td><b>Functions</font></b></td>\n')
              self.file.write('       <td bgcolor = "DarkSlateGrey"><b><font color = "white">' + percentage + '%</font></b></td>\n')
              self.file.write('       <td bgcolor = "DarkSlateGrey"><b><font color = "white">' + reached_function_nb + ' over ' + total_function_nb + '</font></b></td>\n')
              self.file.write('     </tr>\n')
            result = re.search('Generated coverage output', line)
            if result is not None:
              self.file.write('     <tr>\n')
              self.file.write('       <td colspan = 3>More details in artifact magma_logs.zip (code_coverage.zip)</td>\n')
              self.file.write('     </tr>\n')
      logfile.close()

    self.file.write('  </table>\n')
    self.file.write('  </div>\n')
    self.file.write('  <br>\n')


# --------------------------------------------------------------------------------------------------------
#
# Start of main
#
# --------------------------------------------------------------------------------------------------------

argvs = sys.argv
argc = len(argvs)

HTML = HtmlReport()

while len(argvs) > 1:
  myArgv = argvs.pop(1)
  if re.match('^\-\-help$', myArgv, re.IGNORECASE):
    sys.exit(0)
  elif re.match('^\-\-job_name=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-job_name=(.+)$', myArgv, re.IGNORECASE)
    HTML.job_name = matchReg.group(1)
  elif re.match('^\-\-job_id=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-job_id=(.+)$', myArgv, re.IGNORECASE)
    HTML.job_id = matchReg.group(1)
  elif re.match('^\-\-job_url=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-job_url=(.+)$', myArgv, re.IGNORECASE)
    HTML.job_url = matchReg.group(1)
  elif re.match('^\-\-git_url=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-git_url=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_url = matchReg.group(1)
  elif re.match('^\-\-git_src_branch=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-git_src_branch=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_src_branch = matchReg.group(1)
  elif re.match('^\-\-git_src_commit=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-git_src_commit=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_src_commit = matchReg.group(1)
  elif re.match('^\-\-git_src_commit_msg=(.+)$', myArgv, re.IGNORECASE):
    # Not Mandatory
    matchReg = re.match('^\-\-git_src_commit_msg=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_src_commit_msg = matchReg.group(1)
  elif re.match('^\-\-git_merge_request=(.+)$', myArgv, re.IGNORECASE):
    # Can be silent: would be false!
    matchReg = re.match('^\-\-git_merge_request=(.+)$', myArgv, re.IGNORECASE)
    if matchReg.group(1) == 'true' or matchReg.group(1) == 'True':
      HTML.git_merge_request = True
  elif re.match('^\-\-git_target_branch=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-git_target_branch=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_target_branch = matchReg.group(1)
  elif re.match('^\-\-git_target_commit=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-git_target_commit=(.+)$', myArgv, re.IGNORECASE)
    HTML.git_target_commit = matchReg.group(1)
  else:
    sys.exit('Invalid Parameter: ' + myArgv)

if HTML.job_name == '' or HTML.job_id == '' or HTML.job_url == '':
  sys.exit('Missing Parameter in job description')

if HTML.git_url == '' or HTML.git_src_branch == '' or HTML.git_src_commit == '':
  sys.exit('Missing Parameter in Git Repository description')

if HTML.git_merge_request:
  if HTML.git_target_commit == '' or HTML.git_target_branch == '':
     sys.exit('Missing Parameter in Git Pull Request Repository description')

HTML.generate()
