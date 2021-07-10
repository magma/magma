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

HEADER_TEMPLATE = 'ci-scripts/html-templates/file-header.htm'
FOOTER_TEMPLATE = 'ci-scripts/html-templates/file-footer.htm'
PUSH_BUILD_SUMMARY = 'ci-scripts/html-templates/git-summary-push.htm'
PULL_REQUEST_BUILD_SUMMARY = 'ci-scripts/html-templates/git-summary-pr.htm'
BUILD_HEADER = 'ci-scripts/html-templates/build-summary-header.htm'
BUILD_FOOTER = 'ci-scripts/html-templates/build-summary-footer.htm'
MAGMA_COMMON_ROWS = 'ci-scripts/html-templates/build-rows.htm'
TARGET_ROW = 'ci-scripts/html-templates/image-row.htm'
DETAILS_HEADER = 'ci-scripts/html-templates/details-header.htm'
DETAILS_FOOTER = 'ci-scripts/html-templates/details-footer.htm'
DETAILS_T_HEADER = 'ci-scripts/html-templates/details-t-header.htm'
DETAILS_T_FOOTER = 'ci-scripts/html-templates/details-t-footer.htm'
DETAILS_T_ROW = 'ci-scripts/html-templates/details-t-row.htm'

MAX_ALLOWED_WARNINGS = 20
COMMON_TYPE = 'MAGMA-COMMON'
MME_TYPE = 'MAGMA-OAI-MME'
SCTPD_TYPE = 'MAGMA-SCTPD'
U18_BUILD_LOG_FILE = 'build_magma_mme.log'
RHEL8_BUILD_LOG_FILE = 'build_magma_mme_rhel8.log'

COMMON_SECTION_START = "Creating directories for 'MagmaCommon'"
COMMON_SECTION_STOP = "Completed 'MagmaCommon'"
MME_SECTION_START = "Creating directories for 'MagmaCore'"
MME_SECTION_STOP = "Completed 'MagmaCore'"
SCTPD_SECTION_START = "Creating directories for 'Sctpd'"
SCTPD_SECTION_STOP = "Completed 'Sctpd'"


def generate_header(args):
    """
    Append HTML header to file

    Args:
        args: results from argument parser

    Returns:
        a string with formatted HTML header.
    """
    cwd = os.getcwd()
    header = ''
    with open(os.path.join(cwd, HEADER_TEMPLATE), 'r') as temp:
        header = temp.read()
        header = re.sub('JOB_NAME', args.job_name, header)
        header = re.sub('BUILD_ID', args.job_id, header)
        header = re.sub('BUILD_URL', args.job_url, header)
    return header


def generate_footer():
    """
    Append HTML footer to file

    Returns:
        a string with formatted HTML footer.
    """
    cwd = os.getcwd()
    footer = ''
    with open(os.path.join(cwd, FOOTER_TEMPLATE), 'r') as temp:
        footer = temp.read()
    return footer


def generate_git_summary(args):
    """
    Append HTML build summary to file

    Args:
        args: results from argument parser

    Returns:
        a string with formatted HTML build summary.
    """
    if args.git_merge_request:
        template = PULL_REQUEST_BUILD_SUMMARY
    else:
        template = PUSH_BUILD_SUMMARY
    cwd = os.getcwd()
    summary = ''
    with open(os.path.join(cwd, template), 'r') as temp:
        summary = temp.read()
        summary = re.sub('GIT_REPO_URL', args.git_url, summary)
        summary = re.sub('GIT_SRC_BRANCH', args.git_src_branch, summary)
        summary = re.sub('GIT_SRC_COMMIT', args.git_src_commit, summary)
        if args.git_merge_request:
            summary = re.sub('GIT_TGT_BRANCH', args.git_target_branch, summary)
            summary = re.sub('GIT_TGT_COMMIT', args.git_target_commit, summary)
    return summary


def generate_build_header():
    """
    Append HTML Build Summary Header to file

    Returns:
        a string with formatted HTML Build Summary Header.
    """
    cwd = os.getcwd()
    header = ''
    with open(os.path.join(cwd, BUILD_HEADER), 'r') as temp:
        header = temp.read()
    return header


def generate_build_footer():
    """
    Append HTML Build Summary footer to file

    Returns:
        a string with formatted HTML Build Summary footer.
    """
    cwd = os.getcwd()
    footer = ''
    with open(os.path.join(cwd, BUILD_FOOTER), 'r') as temp:
        footer = temp.read()
    return footer


def add_compile_rows():
    """
    Append HTML Build Summary rows

    Returns:
        a string with formatted HTML Build rows.
    """
    rows = add_compile_row(COMMON_TYPE)
    rows += add_compile_row(MME_TYPE)
    rows += add_compile_row(SCTPD_TYPE)
    return rows


def add_compile_row(nf_type):
    """
    Append HTML Build Summary row

    Args:
        nf_type: part of the build

    Returns:
        a string with formatted HTML Build row.
    """
    cwd = os.getcwd()
    rows = ''
    with open(os.path.join(cwd, MAGMA_COMMON_ROWS), 'r') as temp:
        rows = temp.read()

    log_file_names = [
        (U18_BUILD_LOG_FILE, 'U18_'),
        (RHEL8_BUILD_LOG_FILE, 'RHEL8_'),
    ]
    for log_file_name, prefix in log_file_names:
        status = False
        error_cnt = 0
        warning_cnt = 0
        if nf_type == COMMON_TYPE:
            section_start_pattern = COMMON_SECTION_START
            section_end_pattern = COMMON_SECTION_STOP
        if nf_type == MME_TYPE:
            section_start_pattern = MME_SECTION_START
            section_end_pattern = MME_SECTION_STOP
        if nf_type == SCTPD_TYPE:
            section_start_pattern = SCTPD_SECTION_START
            section_end_pattern = SCTPD_SECTION_STOP
        section_status = False
        with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
            for line in logfile:
                my_res = re.search(section_start_pattern, line)
                if my_res is not None:
                    section_status = True
                my_res = re.search(section_end_pattern, line)
                if my_res is not None and section_status:
                    section_status = False
                    status = True
                if section_status:
                    my_res = re.search('error:', line)
                    if my_res is not None:
                        error_cnt += 1
                    my_res = re.search('warning:', line)
                    if my_res is not None:
                        warning_cnt += 1
        if nf_type == COMMON_TYPE:
            rows = re.sub('BUILD_PART_TITLE', 'magma-common', rows)
            rows = re.sub('BUILD_COMMAND', 'magma_common', rows)
        if nf_type == MME_TYPE:
            rows = re.sub('BUILD_PART_TITLE', 'magma-oai-mme', rows)
            rows = re.sub('BUILD_COMMAND', 'core', rows)
        if nf_type == SCTPD_TYPE:
            rows = re.sub('BUILD_PART_TITLE', 'magma-sctpd', rows)
            rows = re.sub('BUILD_COMMAND', 'sctpd', rows)
        if status:
            rows = re.sub(prefix + 'STATUS_COLOR', 'LimeGreen', rows)
            rows = re.sub(prefix + 'STATUS', 'OK', rows)
        else:
            rows = re.sub(prefix + 'STATUS_COLOR', 'Tomato', rows)
            rows = re.sub(prefix + 'STATUS', 'KO', rows)
        if error_cnt == 0 and warning_cnt == 0:
            rows = re.sub(prefix + 'COUNT_COLOR', 'LimeGreen', rows)
        elif error_cnt == 0 and warning_cnt < MAX_ALLOWED_WARNINGS:
            rows = re.sub(prefix + 'COUNT_COLOR', 'Orange', rows)
        else:
            rows = re.sub(prefix + 'COUNT_COLOR', 'Tomato', rows)
        rows = re.sub(prefix + 'WARNINGS', str(warning_cnt), rows)
        rows = re.sub(prefix + 'ERRORS', str(error_cnt), rows)

    return rows


def add_target_image_gen_row():
    """
    Append HTML Build Summary row

    Returns:
        a string with formatted HTML Build row.
    """
    cwd = os.getcwd()
    rows = ''
    with open(os.path.join(cwd, TARGET_ROW), 'r') as temp:
        rows = temp.read()
        rows = re.sub('ROW_TITLE', 'Image Creation Status', rows)
        rows = re.sub('IMAGE_KIND', 'Target Image', rows)

    log_file_names = [
        (U18_BUILD_LOG_FILE, 'U18_'),
        (RHEL8_BUILD_LOG_FILE, 'RHEL8_'),
    ]
    for log_file_name, prefix in log_file_names:
        status = False
        start_pattern = '[aA][sS] magma-mme$'
        end_pattern = 'Successfully tagged magma-mme:|COMMIT magma-mme:'
        section_status = False
        with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
            for line in logfile:
                my_res = re.search(start_pattern, line)
                if my_res is not None:
                    section_status = True
                my_res = re.search(end_pattern, line)
                if my_res is not None and section_status:
                    section_status = False
                    status = True
        if status:
            rows = re.sub(prefix + 'COLOR', 'LimeGreen', rows)
            rows = re.sub(prefix + 'STATUS', 'OK', rows)
        else:
            rows = re.sub(prefix + 'COLOR', 'Tomato', rows)
            rows = re.sub(prefix + 'STATUS', 'KO', rows)

    return rows


def add_target_image_size_row(args):
    """
    Append HTML Build Summary row

    Args:
        args: results from argument parser

    Returns:
        a string with formatted HTML Build row.
    """
    cwd = os.getcwd()
    rows = ''
    with open(os.path.join(cwd, TARGET_ROW), 'r') as temp:
        rows = temp.read()
        rows = re.sub('ROW_TITLE', 'Image Size', rows)
        rows = re.sub('IMAGE_KIND', 'Target Image', rows)

    if args.git_merge_request:
        image_tag = 'magma-mme *ci-tmp .*ago *([0-9 A-Z]+)'
    else:
        image_tag = 'magma-mme *master .*ago *([0-9 A-Z]+)'

    log_file_names = [
        (U18_BUILD_LOG_FILE, 'U18_'),
        (RHEL8_BUILD_LOG_FILE, 'RHEL8_'),
    ]
    for log_file_name, prefix in log_file_names:
        status = False
        start_pattern = 'Successfully tagged magma-mme:|COMMIT magma-mme:'
        end_pattern = 'MAGMA-OAI-MME '
        section_status = False
        with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
            for line in logfile:
                my_res = re.search(start_pattern, line)
                if my_res is not None:
                    section_status = True
                my_res = re.search(end_pattern, line)
                if my_res is not None and section_status:
                    section_status = False
                my_res = re.search(image_tag, line)
                if my_res is not None and section_status:
                    size = my_res.group(1)
                    status = True

        if status:
            rows = re.sub(prefix + 'COLOR', 'LimeGreen', rows)
            stat_txt = 'OK\n  ' + size
            rows = re.sub(prefix + 'STATUS', stat_txt, rows)
        else:
            rows = re.sub(prefix + 'COLOR', 'Tomato', rows)
            rows = re.sub(prefix + 'STATUS', 'KO', rows)

    return rows


def add_compilation_details():
    """
    Append HTML compilation warnings/errors details

    Returns:
        a string with formatted HTML details section.
    """
    cwd = os.getcwd()
    details = ''
    with open(os.path.join(cwd, DETAILS_HEADER), 'r') as temp:
        details = temp.read()

    details += add_compilation_details_table(COMMON_TYPE)
    details += add_compilation_details_table(MME_TYPE)
    details += add_compilation_details_table(SCTPD_TYPE)

    with open(os.path.join(cwd, DETAILS_FOOTER), 'r') as temp:
        details += temp.read()
    return details


def add_compilation_details_table(nf_type):
    """
    Append HTML Build Summary row

    Args:
        nf_type: part of the build

    Returns:
        a string with formatted HTML Build row.
    """
    cwd = os.getcwd()
    with open(os.path.join(cwd, DETAILS_T_ROW), 'r') as temp:
        def_row = temp.read()

    table = ''

    log_file_names = [
        (U18_BUILD_LOG_FILE, 'UBUNTU 18'),
        (RHEL8_BUILD_LOG_FILE, 'RHEL8'),
    ]
    for log_file_name, variant in log_file_names:
        error_cnt = 0
        warning_cnt = 0
        errors_and_warnings = []
        if nf_type == COMMON_TYPE:
            section_start_pattern = COMMON_SECTION_START
            section_end_pattern = COMMON_SECTION_STOP
            build_variant = 'MAGMA-COMMON ' + variant
        if nf_type == MME_TYPE:
            section_start_pattern = MME_SECTION_START
            section_end_pattern = MME_SECTION_STOP
            build_variant = 'MAGMA-OAI-MME ' + variant
        if nf_type == SCTPD_TYPE:
            section_start_pattern = SCTPD_SECTION_START
            section_end_pattern = SCTPD_SECTION_STOP
            build_variant = 'MAGMA-SCTPD ' + variant
        section_status = False
        with open(cwd + '/archives/' + log_file_name, 'r') as logfile:
            for line in logfile:
                my_res = re.search(section_start_pattern, line)
                if my_res is not None:
                    section_status = True
                my_res = re.search(section_end_pattern, line)
                if my_res is not None and section_status:
                    section_status = False
                if section_status:
                    my_res = re.search('error:', line)
                    if my_res is not None:
                        error_cnt += 1
                        l_details = retrieve_details(line.strip(), 'error')
                        errors_and_warnings.append(l_details)
                    my_res = re.search('warning:', line)
                    if my_res is not None:
                        warning_cnt += 1
                        l_details = retrieve_details(line.strip(), 'warning')
                        errors_and_warnings.append(l_details)

        if error_cnt == 0 and warning_cnt == 0:
            continue
        with open(os.path.join(cwd, DETAILS_T_HEADER), 'r') as temp:
            table += temp.read()
            table = re.sub('BUILD_VARIANT', build_variant, table)
        for detail in errors_and_warnings:
            new_row = re.sub('FILENAME', detail['file_name'], def_row)
            new_row = re.sub('LINE_NB', detail['line_nb'], new_row)
            new_row = re.sub('KIND', detail['kind'], new_row)
            new_row = re.sub('WAR_TYPE', detail['warning_type'], new_row)
            new_row = re.sub('MSG', detail['msg'], new_row)
            table += new_row
        with open(os.path.join(cwd, DETAILS_T_FOOTER), 'r') as temp:
            table += temp.read()

    return table


def retrieve_details(line, kind):
    """
    Retrieve details

    Args:
        line: string from compilation log
        kind: error or warning

    Returns:
        an object with details info
    """
    info = {}
    file_name = re.sub(':.*$', '', line)
    file_name = re.sub('^/magma/', '', file_name)

    msg = re.sub('^.*' + kind + ': ', '', line)
    if kind == 'warning':
        msg = msg.replace(' [-W', ' -W')
        msg = re.sub(' -W.*$', '', msg)
        warning_type = re.sub('^.*-W', '', line)
        warning_type = re.sub('].*$', '', warning_type)
        if re.search('proto but not used', line) is not None:
            warning_type = 'proto is not used'

    line_nb = '0'
    infos = re.search(':(?P<linenb>[0-9]+):', line)
    if infos is not None:
        line_nb = infos.group('linenb')

    if kind == 'error':
        info['kind'] = 'Error'
        info['warning_type'] = 'fatal'
    else:
        info['kind'] = 'Warning'
        info['warning_type'] = warning_type
    info['file_name'] = file_name
    info['line_nb'] = line_nb
    info['msg'] = msg

    return info
