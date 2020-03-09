#/*
# * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# * contributor license agreements.  See the NOTICE file distributed with
# * this work for additional information regarding copyright ownership.
# * The OpenAirInterface Software Alliance licenses this file to You under
# * the OAI Public License, Version 1.1  (the "License"); you may not use this file
# * except in compliance with the License.
# * You may obtain a copy of the License at
# *
# *   http://www.openairinterface.org/?page_id=698
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

		#self.analyze_sca_log()

		self.buildSummaryHeader()
		self.vmWakeUpRow()
		self.makeRunRow()
		self.statusCheckRow()
		self.buildSummaryFooter()

		self.testSummaryHeader()
		self.s1apTesterTable()
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
		self.file.write('  <div class="well well-lg">End of Build Report -- Copyright <span class="glyphicon glyphicon-copyright-mark"></span> 2020 <a href="http://www.openairinterface.org/">OpenAirInterface</a>. All Rights Reserved.</div>\n')
		self.file.write('</div></body>\n')
		self.file.write('</html>\n')

	def buildSummaryHeader(self):
		self.file.write('  <h2>Build Summary</h2>\n')
		self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
		self.file.write('     <tr bgcolor="#33CCFF" >\n')
		self.file.write('       <th>Stage Name</th>\n')
		self.file.write('       <th>Magma-dev VM</th>\n')
		self.file.write('       <th>Magma-test VM</th>\n')
		self.file.write('       <th>Magma-traffic VM</th>\n')

	def buildSummaryFooter(self):
		self.file.write('  </table>\n')
		self.file.write('  <br>\n')

	def vmWakeUpRow(self):
		self.file.write('    <tr>\n')
		self.file.write('      <td bgcolor="lightcyan" >Waking-Up Vagrant VMs</td>\n')
		self.analyze_vagrant_up_log('magma')
		self.analyze_vagrant_up_log('magma_test')
		self.analyze_vagrant_up_log('magma_trfserver')
		self.file.write('    </tr>\n')

	def analyze_vagrant_up_log(self, vmType):
		logFileName = vmType.lower().replace('magma','magma_vagrant') + '_up.log'
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

	def makeRunRow(self):
		self.file.write('    <tr>\n')
		self.file.write('      <td bgcolor="lightcyan" >Make Run Results</td>\n')
		self.analyze_vagrant_make_run_log('magma')
		self.analyze_vagrant_make_run_log('magma_test')
		self.analyze_vagrant_make_run_log('magma_trfserver')
		self.file.write('    </tr>\n')

	def analyze_vagrant_make_run_log(self, vmType):
		if vmType == 'magma_test' or vmType == 'magma_trfserver':
			cell_msg = '      <td bgcolor="LightGray"><pre style="border:none; background-color:LightGray"><b>'
			cell_msg += 'N/A</b></pre></td>\n'
			self.file.write(cell_msg)
			return

		logFileName = vmType.lower().replace('magma','magma_vagrant') + '_make_run.log'
		module_pattern = 'ninja -C  '
		end_pattern = 'sudo service magma@magmad start'

		cwd = os.getcwd()
		if os.path.isfile(cwd + '/archives/' + logFileName):
			status = False
			firstModulePassed = False
			moduleList = ''
			module_name = ''
			with open(cwd + '/archives/' + logFileName, 'r') as logfile:
				for line in logfile:
					result = re.search(end_pattern, line)
					if result is not None:
						status = True
					result = re.search(module_pattern, line)
					if result is not None:
						if firstModulePassed:
							moduleList += '  -- ' + module_name + ': OK\n'
						firstModulePassed = True
						module_name = line.replace('\n','').replace('ninja -C  /home/vagrant/build/c//','')
				logfile.close()
			if status:
				moduleList += '  -- ' + module_name + ': OK\n'
				cell_msg = '      <td bgcolor="LimeGreen"><pre style="border:none; background-color:LimeGreen"><b>'
				cell_msg += 'OK: \n'
			else:
				moduleList += '  -- ' + module_name + ': KO\n'
				cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
				cell_msg += 'KO: \n'
			cell_msg += moduleList
			cell_msg += '</b></pre></td>\n'
		else:
			cell_msg = '      <td bgcolor="Tomato"><pre style="border:none; background-color:Tomato"><b>'
			cell_msg += 'KO: logfile (' + logFileName + ') not found</b></pre></td>\n'

		self.file.write(cell_msg)

	def statusCheckRow(self):
		self.file.write('    <tr>\n')
		self.file.write('      <td bgcolor="lightcyan" >Magma Services Status</td>\n')
		self.analyze_status_check_log('magma')
		self.analyze_status_check_log('magma_test')
		self.analyze_status_check_log('magma_trfserver')
		self.file.write('    </tr>\n')

	def analyze_status_check_log(self, vmType):
		if vmType == 'magma_test' or vmType == 'magma_trfserver':
			cell_msg = '      <td bgcolor="LightGray"><pre style="border:none; background-color:LightGray"><b>'
			cell_msg += 'N/A</b></pre></td>\n'
			self.file.write(cell_msg)
			return

		logFileName = 'magma_status.log'
		pattern_module = '● magma@'
		active_pattern = 'Active: active'

		cwd = os.getcwd()
		if os.path.isfile(cwd + '/archives/' + logFileName):
			status = True
			service_status = False
			firstService = False
			listServices = ''
			with open(cwd + '/archives/' + logFileName, 'r') as logfile:
				for line in logfile:
					result = re.search(pattern_module, line)
					if result is not None:
						if firstService and not service_status:
							listServices += ' -- ' + service_name + '  NOT ACTIVE\n'
							status = False
						service_name = line.replace('\n','').replace(pattern_module,'')
						service_name = service_name.replace('.service','')
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

	def testSummaryHeader(self):
		self.file.write('  <h2>Test Summary</h2>\n')

	def testSummaryFooter(self):
		self.file.write('  <br>\n')

	def s1apTesterTable(self):
		cwd = os.getcwd()
		if os.path.isfile(cwd + '/archives/magma_run_s1ap_tester.log'):
			cmd = 'egrep -c "AGW-VM-S1AP-TESTS: OK" ' + cwd + '/archives/magma_run_s1ap_tester.log'
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
		self.file.write('  <button data-toggle="collapse" data-target="#s1ap-tester">More details on S1AP Tester checks</button>\n')
		self.file.write('  <div id="s1ap-tester" class="collapse">\n')
		self.file.write('  <table class="table-bordered" width = "100%" align = "center" border = "1">\n')
		self.file.write('     <tr bgcolor="#33CCFF" >\n')
		self.file.write('       <th>Test Name</th>\n')
		self.file.write('       <th>Test Status</th>\n')
		self.file.write('       <th>Test Duration</th>\n')
		self.file.write('     </tr>\n')

		if os.path.isfile(cwd + '/archives/magma_run_s1ap_tester.log'):
			pattern = 'echo "Running test: '
			status = True
			with open(cwd + '/archives/magma_run_s1ap_tester.log', 'r') as logfile:
				for line in logfile:
					result = re.search(pattern + '(.+)$', line)
					if result is not None:
						test_name = result.group(1).replace('"','')
					result = re.search('Ran 1 test in (.+)$', line)
					if result is not None:
						test_duration = result.group(1)
					result = re.search('^OK$|^FAILED', line)
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
				self.file.write('       <td align = "right" colspan = 2 bgcolor = "Green"><b><font color = "white">OK</font></b></td>\n')
			else:
				self.file.write('       <td align = "right" colspan = 2 bgcolor = "Red"><b><font color = "white">FAIL</font></b></td>\n')
			self.file.write('     </tr>\n')
		else:
			self.file.write('     <tr bgcolor="Tomato" >\n')
			self.file.write('       <td colspan=3><b>KO: logfile (magma_run_s1ap_tester.log) not found</b></td>\n')
			self.file.write('     </tr>\n')

		self.file.write('  </table>\n')
		self.file.write('  </div>\n')
		self.file.write('  <br>\n')


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
