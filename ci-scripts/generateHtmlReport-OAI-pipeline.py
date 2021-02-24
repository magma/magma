#/*
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
#---------------------------------------------------------------------

import os
import re
import sys
import subprocess

class HtmlReport():
  def __init__(self):
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

  def generateBuild(self):
    cwd = os.getcwd()
    self.file = open(cwd + '/build_results_magma_oai_mme.html', 'w')
    self.generateHeader()

    self.buildSummaryHeader()
    self.buildCompileRows()
    self.copyToTargetImage()
    self.copyConfToolsToTargetImage()
    self.imageSizeRow()
    self.buildSummaryFooter()

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
    self.file.write('  <title>MAGMA/OAI Core Network Build Results for ' + self.job_name + ' job build #' + self.job_id + '</title>\n')
    self.file.write('</head>\n')
    self.file.write('<body><div class="container">\n')
    self.file.write('  <br>\n')
    self.file.write('  <table width = "100%" style="border-collapse: collapse; border: none;">\n')
    self.file.write('   <tr style="border-collapse: collapse; border: none;">\n')
    # SVG has a invisible background color -- adding it.
    self.file.write('     <td bgcolor="#5602a4" style="border-collapse: collapse; border: none;">\n')
    self.file.write('       <a href="https://www.magmacore.org/">\n')
    self.file.write('          <img src="https://www.magmacore.org/img/magma-logo.svg" alt="" border="none" height=50 width=150>\n')
    self.file.write('          </img>\n')
    self.file.write('       </a>\n')
    self.file.write('     </td>\n')
    self.file.write('     <td align = "center" style="border-collapse: collapse; border: none; vertical-align: center;">\n')
    self.file.write('       <b><font size = "6">Job Summary -- Job: ' + self.job_name + ' -- Build-ID: <a href="' + self.job_url + '">' + self.job_id + '</a></font></b>\n')
    self.file.write('     </td>\n')
    self.file.write('     <td style="border-collapse: collapse; border: none;">\n')
    self.file.write('       <a href="http://www.openairinterface.org/">\n')
    self.file.write('          <img src="http://www.openairinterface.org/wp-content/uploads/2016/03/cropped-oai_final_logo2.png" alt="" border="none" height=50 width=150>\n')
    self.file.write('          </img>\n')
    self.file.write('       </a>\n')
    self.file.write('     </td>\n')
    self.file.write('   </tr>\n')
    self.file.write('  </table>\n')
    self.file.write('  <br>\n')
    buildSummary = self.generateBuildSummary()
    self.file.write(buildSummary)

  def generateBuildSummary(self):
    returnString = ''
    # Build Info Summary
    returnString += '  <table class="table-bordered" width = "80%" align = "center" border = "1">\n'
    returnString += '    <tr>\n'
    returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-time"></span> Build Start Time</td>\n'
    #date_formatted = re.sub('\..*', '', self.created
    returnString += '      <td>' + self.job_start_time + '</td>\n'
    returnString += '    </tr>\n'
    returnString += '    <tr>\n'
    returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-wrench"></span> Build Trigger</td>\n'
    if self.git_merge_request:
      returnString += '      <td>Pull Request</td>\n'
    else:
      returnString += '      <td>Push Event</td>\n'
    returnString += '    </tr>\n'
    returnString += '    <tr>\n'
    returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-cloud-upload"></span> GIT Repository</td>\n'
    returnString += '      <td><a href="' + self.git_url + '">' + self.git_url + '</a></td>\n'
    returnString += '    </tr>\n'
    if self.git_merge_request:
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-link"></span> Pull Request Link</td>\n'
      returnString += '      <td><a href="TEMPLATE_PULL_REQUEST_LINK">TEMPLATE_PULL_REQUEST_LINK</a></td>\n'
      returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-header"></span> Pull Request Title</td>\n'
      returnString += '      <td>TEMPLATE_PULL_REQUEST_TEMPLATE</td>\n'
      returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-out"></span> Source Branch</td>\n'
      returnString += '      <td>' + self.git_src_branch + '</td>\n'
      returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Source Commit ID</td>\n'
      returnString += '      <td>' + self.git_src_commit + '</td>\n'
      returnString += '    </tr>\n'
      if (self.git_src_commit_msg is not None):
        returnString += '    <tr>\n'
        returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Source Commit Message</td>\n'
        returnString += '      <td>' + self.git_src_commit_msg + '</td>\n'
        returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-log-in"></span> Target Branch</td>\n'
      returnString += '      <td>' + self.git_target_branch + '</td>\n'
      returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Target Commit ID</td>\n'
      returnString += '      <td>' + self.git_target_commit + '</td>\n'
      returnString += '    </tr>\n'
    else:
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tree-deciduous"></span> Branch</td>\n'
      returnString += '      <td>' + self.git_src_branch + '</td>\n'
      returnString += '    </tr>\n'
      returnString += '    <tr>\n'
      returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-tag"></span> Commit ID</td>\n'
      returnString += '      <td>' + self.git_src_commit + '</td>\n'
      returnString += '    </tr>\n'
      if (self.git_src_commit_msg is not None):
        returnString += '    <tr>\n'
        returnString += '      <td bgcolor="lightcyan" > <span class="glyphicon glyphicon-comment"></span> Commit Message</td>\n'
        returnString += '      <td>' + self.git_src_commit_msg + '</td>\n'
        returnString += '    </tr>\n'
    returnString += '  </table>\n'
    returnString += '  <br>\n'
    return returnString

  def generateFooter(self):
    self.file.write('  </nav>\n')
    self.file.write('  <div class="well well-lg">End of Build Report -- Copyright <span class="glyphicon glyphicon-copyright-mark"></span> 2020 <a href="http://www.openairinterface.org/">OpenAirInterface</a>. All Rights Reserved.</div>\n')
    self.file.write('</div></body>\n')
    self.file.write('</html>\n')

  def buildSummaryHeader(self):
    self.file.write('  <h2>Docker Image Build Summary</h2>\n')
    self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
    self.file.write('     <tr bgcolor="#33CCFF" >\n')
    self.file.write('       <th>Stage Name</th>\n')
    self.file.write('       <th>Image Kind</th>\n')
    self.file.write('       <th>MAGMA - OAI MME cNF</th>\n')
    self.file.write('     </tr>\n')

  def buildSummaryFooter(self):
    self.file.write('  </table>\n')
    self.file.write('  <br>\n')

  def buildCompileRows(self):
    self.file.write('    <tr>\n')
    self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-common</b> Compile / Build</td>\n')
    self.analyze_build_log('OAI-COMMON')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.analyze_compile_log('OAI-COMMON')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-oai-mme</b> Compile / Build</td>\n')
    self.analyze_build_log('OAI-MME')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.analyze_compile_log('OAI-MME')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.file.write('      <td rowspan=2 bgcolor="lightcyan" ><b>magma-sctpd</b> Compile / Build</td>\n')
    self.analyze_build_log('SCTPD')
    self.file.write('    </tr>\n')
    self.file.write('    <tr>\n')
    self.analyze_compile_log('SCTPD')
    self.file.write('    </tr>\n')

  def analyze_build_log(self, nfType):
    if nfType != 'OAI-COMMON' and nfType != 'OAI-MME' and nfType != 'SCTPD':
      self.file.write('      <td>N/A</td>\n')
      self.file.write('      <td>Wrong NF Type for this Report</td>\n')
      return

    logFileName = 'build_magma_mme.log'
    self.file.write('      <td>Builder Image</td>\n')

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      status = False
      if nfType == 'OAI-COMMON':
        section_start_pattern = 'ninja -C  /build/c/magma_common'
        section_end_pattern = 'cmake  /magma/lte/gateway/c/oai -DCMAKE_BUILD_TYPE=Debug  -DS6A_OVER_GRPC=False -GNinja'
      if nfType == 'OAI-MME':
        section_start_pattern = 'ninja -C  /build/c/oai'
        section_end_pattern = 'cmake  /magma/orc8r/gateway/c/common -DCMAKE_BUILD_TYPE=Debug   -GNinja'
      if nfType == 'SCTPD':
        section_start_pattern = 'ninja -C  /build/c/sctpd'
        section_end_pattern = 'FROM ubuntu:bionic as magma-mme'
      section_status = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
                    result = re.search(section_start_pattern, line)
                    if result is not None:
                        section_status = True
                    result = re.search(section_end_pattern, line)
                    if result is not None:
                        section_status = False
                    if section_status:
                        if nfType == 'OAI-COMMON':
                          result = re.search('Linking CXX static library eventd/libEVENTD.a', line)
                        if nfType == 'OAI-MME':
                          result = re.search('Linking CXX executable oai_mme/mme', line)
                        if nfType == 'SCTPD':
                          result = re.search('Linking CXX executable sctpd', line)
                        if result is not None:
                            status = True
        logfile.close()
      if status:
        cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
        cell_msg += 'OK:\n'
      else:
        cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
        cell_msg += 'KO:\n'
      if nfType == 'OAI-COMMON':
        cell_msg += ' -- ninja -C  /build/c/magma_common</b></pre></td>\n'
      if nfType == 'OAI-MME':
        cell_msg += ' -- ninja -C  /build/c/oai</b></pre></td>\n'
      if nfType == 'SCTPD':
        cell_msg += ' -- ninja -C  /build/c/sctpd</b></pre></td>\n'
    else:
      cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def analyze_compile_log(self, nfType):
    if nfType != 'OAI-COMMON' and nfType != 'OAI-MME' and nfType != 'SCTPD':
      self.file.write('      <td>N/A</td>\n')
      self.file.write('      <td>Wrong NF Type for this Report</td>\n')
      return

    logFileName = 'build_magma_mme.log'
    self.file.write('      <td>Builder Image</td>\n')

    cwd = os.getcwd()
    nb_errors = 0
    nb_warnings = 0

    if os.path.isfile(cwd + '/archives/' + logFileName):
      if nfType == 'OAI-COMMON':
        section_start_pattern = 'ninja -C  /build/c/magma_common'
        section_end_pattern = 'cmake  /magma/lte/gateway/c/oai -DCMAKE_BUILD_TYPE=Debug  -DS6A_OVER_GRPC=False -GNinja'
      if nfType == 'OAI-MME':
        section_start_pattern = 'ninja -C  /build/c/oai'
        section_end_pattern = 'cmake  /magma/orc8r/gateway/c/common -DCMAKE_BUILD_TYPE=Debug   -GNinja'
      if nfType == 'SCTPD':
        section_start_pattern = 'ninja -C  /build/c/sctpd'
        section_end_pattern = 'FROM ubuntu:bionic as magma-mme'
      section_status = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(section_start_pattern, line)
          if result is not None:
            section_status = True
          result = re.search(section_end_pattern, line)
          if result is not None:
            section_status = False
          if section_status:
            result = re.search('error:', line)
            if result is not None:
              nb_errors += 1
              result = re.search('warning:', line)
              if result is not None:
                nb_warnings += 1
        logfile.close()
      if nb_warnings == 0 and nb_errors == 0:
        cell_msg = '       <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
      elif nb_warnings < 20 and nb_errors == 0:
        cell_msg = '       <td bgcolor="Orange"><pre style="border:none; background-color:Orange"><b>'
      else:
        cell_msg = '       <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      if nb_errors > 0:
        cell_msg += str(nb_errors) + ' errors found in compile log\n'
      cell_msg += str(nb_warnings) + ' warnings found in compile log</b></pre></td>\n'
    else:
      cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def copyToTargetImage(self):
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" >SW libs Installation / Copy from Builder</td>\n')
    self.analyze_copy_log('MME')
    self.file.write('    </tr>\n')

  def analyze_copy_log(self, nfType):
    if nfType != 'MME':
      self.file.write('      <td>N/A</td>\n')
      self.file.write('      <td>Wrong NF Type for this Report</td>\n')
      return

    logFileName = 'build_magma_mme.log'
    self.file.write('      <td>Target Image</td>\n')

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      section_start_pattern = 'FROM ubuntu:bionic as magma-mme$'
      section_end_pattern = 'WORKDIR /magma-mme/bin$'
      section_status = False
      status = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(section_start_pattern, line)
          if result is not None:
            section_status = True
          result = re.search(section_end_pattern, line)
          if result is not None:
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
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def copyConfToolsToTargetImage(self):
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" >Copy Template Conf / Tools from Builder</td>\n')
    self.analyze_copy_conf_tool_log('MME')
    self.file.write('    </tr>\n')

  def analyze_copy_conf_tool_log(self, nfType):
    if nfType != 'MME':
      self.file.write('      <td>N/A</td>\n')
      self.file.write('      <td>Wrong NF Type for this Report</td>\n')
      return

    logFileName = 'build_magma_mme.log'
    self.file.write('      <td>Target Image</td>\n')

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      section_start_pattern = 'WORKDIR /magma-mme/bin$'
      section_end_pattern = 'Successfully tagged magma-mme:'
      section_status = False
      status = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(section_start_pattern, line)
          if result is not None:
            section_status = True
          result = re.search(section_end_pattern, line)
          if result is not None:
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
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def imageSizeRow(self):
    self.file.write('    <tr>\n')
    self.file.write('      <td bgcolor="lightcyan" >Image Size</td>\n')
    self.analyze_image_size_log('MME')
    self.file.write('    </tr>\n')

  def analyze_image_size_log(self, nfType):
    if nfType != 'MME':
      self.file.write('      <td>N/A</td>\n')
      self.file.write('      <td>Wrong NF Type for this Report</td>\n')
      return

    logFileName = 'build_magma_mme.log'
    self.file.write('      <td>Target Image</td>\n')

    cwd = os.getcwd()
    if os.path.isfile(cwd + '/archives/' + logFileName):
      section_start_pattern = 'Successfully tagged magma-mme'
      section_end_pattern = 'MAGMA-OAI-MME DOCKER IMAGE BUILD'
      section_status = False
      status = False
      with open(cwd + '/archives/' + logFileName, 'r') as logfile:
        for line in logfile:
          result = re.search(section_start_pattern, line)
          if result is not None:
            section_status = True
          result = re.search(section_end_pattern, line)
          if result is not None:
            section_status = False
          if section_status:
            result = re.search('magma-mme *ci-tmp', line)
            if result is not None:
              result = re.search('ago *([0-9A-Z]+)', line)
              if result is not None:
                size = result.group(1)
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
      cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

    self.file.write(cell_msg)

  def appendBuildSummary(self, mode):
    cwd = os.getcwd()
    if mode == 'dsTester':
      filename = 'test_results_magma_oai_epc.html'
    if os.path.isfile(cwd + '/' + filename):
      newEpcReport = open(cwd + '/new_' + filename, 'w')
      buildSummaryToBeDone = True
      with open(cwd + '/' + filename, 'r') as originalEpcReport:
        for line in originalEpcReport:
          result = re.search('Deployment Summary', line)
          if (result is not None) and buildSummaryToBeDone:
            buildSummary = self.generateBuildSummary()
            newEpcReport.write(buildSummary)
            buildSummaryToBeDone = False
          newEpcReport.write(line)
        originalEpcReport.close()
      newEpcReport.close()
      os.rename(cwd + '/new_' + filename, cwd + '/' + filename)

#--------------------------------------------------------------------------------------------------------
#
# Start of main
#
#--------------------------------------------------------------------------------------------------------

argvs = sys.argv
argc = len(argvs)

HTML = HtmlReport()

while len(argvs) > 1:
  myArgv = argvs.pop(1)
  if re.match('^\-\-help$', myArgv, re.IGNORECASE):
    print('No help yet.')
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
  elif re.match('^\-\-mode=(.+)$', myArgv, re.IGNORECASE):
    matchReg = re.match('^\-\-mode=(.+)$', myArgv, re.IGNORECASE)
    if matchReg.group(1) == 'Build':
      HTML.mode = 'build'
    elif matchReg.group(1) == 'TestWithDsTest':
      HTML.mode = 'dsTester'
    else:
      sys.exit('Invalid mode: ' + matchReg.group(1))
  else:
    sys.exit('Invalid Parameter: ' + myArgv)

if HTML.job_name == '' or HTML.job_id == '' or HTML.job_url == '' or HTML.mode == '':
  sys.exit('Missing Parameter in job description')

if HTML.git_url == '' or HTML.git_src_branch == '' or HTML.git_src_commit == '':
  sys.exit('Missing Parameter in Git Repository description')

if HTML.git_merge_request:
  if HTML.git_target_commit == '' or HTML.git_target_branch == '':
     sys.exit('Missing Parameter in Git Pull Request Repository description')

if HTML.mode == 'build':
  HTML.generateBuild()
elif HTML.mode == 'dsTester':
  HTML.appendBuildSummary(HTML.mode)
