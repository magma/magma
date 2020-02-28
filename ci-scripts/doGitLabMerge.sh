#!/bin/bash
#/*
# * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# * contributor license agreements.  See the NOTICE file distributed with
# * this work for additional information regarding copyright ownership.
# * The OpenAirInterface Software Alliance licenses this file to You under
# * the OAI Public License, Version 1.1  (the "License"); you may not use this file
# * except in compliance with the License.
# * You may obtain a copy of the License at
# *
# *      http://www.openairinterface.org/?page_id=698
# *
# * Unless required by applicable law or agreed to in writing, software
# * distributed under the License is distributed on an "AS IS" BASIS,
# * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# * See the License for the specific language governing permissions and
# * limitations under the License.
# *-------------------------------------------------------------------------------
# * For more information about the OpenAirInterface (OAI) Software Alliance:
# *      contact@openairinterface.org
# */

function usage {
    echo "OAI GitLab merge request applying script"
    echo "   Original Author: Raphael Defosseux"
    echo ""
    echo "Usage:"
    echo "------"
    echo ""
    echo "    doGitLabMerge.sh [OPTIONS] [MANDATORY_OPTIONS]"
    echo ""
    echo "Mandatory Options:"
    echo "------------------"
    echo ""
    echo "    --src-branch #### OR -sb ####"
    echo "    Specify the source branch of the merge request."
    echo ""
    echo "    --src-commit #### OR -sc ####"
    echo "    Specify the source commit ID (SHA-1) of the merge request."
    echo ""
    echo "    --target-branch #### OR -tb ####"
    echo "    Specify the target branch of the merge request (usually develop)."
    echo ""
    echo "    --target-commit #### OR -tc ####"
    echo "    Specify the target commit ID (SHA-1) of the merge request."
    echo ""
    echo "Options:"
    echo "--------"
    echo "    --help OR -h"
    echo "    Print this help message."
    echo ""
}

if [ $# -ne 8 ] && [ $# -ne 1 ]
then
    echo "Syntax Error: not the correct number of arguments"
    echo ""
    usage
    exit 1
fi

checker=0
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -h|--help)
    shift
    usage
    exit 0
    ;;
    -sb|--src-branch)
    SOURCE_BRANCH="$2"
    let "checker|=0x1"
    shift
    shift
    ;;
    -sc|--src-commit)
    SOURCE_COMMIT_ID="$2"
    let "checker|=0x2"
    shift
    shift
    ;;
    -tb|--target-branch)
    TARGET_BRANCH="$2"
    let "checker|=0x4"
    shift
    shift
    ;;
    -tc|--target-commit)
    TARGET_COMMIT_ID="$2"
    let "checker|=0x8"
    shift
    shift
    ;;
    *)
    echo "Syntax Error: unknown option: $key"
    echo ""
    usage
    exit 1
esac

done

if [[ $TARGET_COMMIT_ID == "latest" ]]
then
    TARGET_COMMIT_ID=`git log -n1 --pretty=format:%H origin/$TARGET_BRANCH`
fi

echo "Source Branch is    : $SOURCE_BRANCH"
echo "Source Commit ID is : $SOURCE_COMMIT_ID"
echo "Target Branch is    : $TARGET_BRANCH"
echo "Target Commit ID is : $TARGET_COMMIT_ID"

if [ $checker -ne 15 ]
then
    echo ""
    echo "Syntax Error: missing option"
    echo ""
    usage
    exit 1
fi

git config user.email "jenkins@openairinterface.org"
git config user.name "OAI Jenkins"

git checkout -f $SOURCE_COMMIT_ID > checkout.txt 2>&1
STATUS=`egrep -c "fatal: reference is not a tree" checkout.txt`
rm -f checkout.txt
if [ $STATUS -ne 0 ]
then
    echo "fatal: reference is not a tree --> $SOURCE_COMMIT_ID"
    STATUS=-1
    exit $STATUS
fi

git log -n1 --pretty=format:\"%s\" > .git/CI_COMMIT_MSG

git merge --ff $TARGET_COMMIT_ID -m "Temporary merge for CI"

STATUS=`git status | egrep -c "You have unmerged paths.|fix conflicts"`
if [ $STATUS -ne 0 ]
then
    echo "There are merge conflicts.. Cannot perform further build tasks"
    STATUS=-1
fi
exit $STATUS
