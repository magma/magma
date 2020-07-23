#!/bin/bash
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
# *      contact@openairinterface.org
# */

function usage {
    echo "OAI Coding / Formatting Guideline Check script"
    echo "   Original Author: Raphael Defosseux"
    echo ""
    echo "   Requirement: clang-format / git shall be installed"
    echo ""
    echo "   By default (no options) the complete repository will be checked"
    echo "   In case of merge/pull request, provided source and target branch,"
    echo "   the script will check only the modified files"
    echo ""
    echo "Usage:"
    echo "------"
    echo "    checkCodingFormattingRules.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "--------"
    echo "    --src-branch #### OR -sb ####"
    echo "    Specify the source branch of the merge request."
    echo ""
    echo "    --target-branch #### OR -tb ####"
    echo "    Specify the target branch of the merge request (usually develop)."
    echo ""
    echo "    --help OR -h"
    echo "    Print this help message."
    echo ""
}

if [ $# -ne 4 ] && [ $# -ne 1 ] && [ $# -ne 0 ]
then
    echo "Syntax Error: not the correct number of arguments"
    echo ""
    usage
    exit 1
fi

if [ $# -eq 0 ]
then
    echo " ---- Checking the whole repository ----"
    echo ""
    if [ -f oai_rules_result.txt ]
    then
        rm -f oai_rules_result.txt
    fi
    if [ -f oai_rules_result_list.txt ]
    then
        rm -f oai_rules_result_list.txt
    fi
    EXTENSION_LIST=("h" "hpp" "c" "cpp")
    NB_TO_FORMAT=0
    NB_TOTAL=0
    for EXTENSION in "${EXTENSION_LIST[@]}"
    do
        echo "Checking for all files with .${EXTENSION} extension"
        FILE_LIST=`tree -n --noreport -i -f -P *.${EXTENSION} | sed -e 's#^\./##' | grep -v test | grep "\.${EXTENSION}"`
        for FILE_TO_CHECK in "${FILE_LIST[@]}"
        do
            TO_FORMAT=`clang-format -output-replacements-xml ${FILE_TO_CHECK} 2>&1 | grep -v replacements | grep -c replacement`
            NB_TOTAL=$((NB_TOTAL + 1))
            if [ $TO_FORMAT -ne 0 ]
            then
                NB_TO_FORMAT=$((NB_TO_FORMAT + 1))
                # In case of full repo, being silent
                #echo "$FILE_TO_CHECK"
                echo "$FILE_TO_CHECK" >> ./oai_rules_result_list.txt
            fi
        done
    done
    echo "Nb Files that do NOT follow OAI rules: $NB_TO_FORMAT over $NB_TOTAL checked!"
    echo "NB_FILES_FAILING_CHECK=$NB_TO_FORMAT" > ./oai_rules_result.txt
    echo "NB_FILES_CHECKED=$NB_TOTAL" >> ./oai_rules_result.txt
    exit 0
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
    -tb|--target-branch)
    TARGET_BRANCH="$2"
    let "checker|=0x2"
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


if [ $checker -ne 3 ]
then
    echo "Source Branch is    : $SOURCE_BRANCH"
    echo "Target Branch is    : $TARGET_BRANCH"
    echo ""
    echo "Syntax Error: missing option"
    echo ""
    usage
    exit 1
fi

# Merge request scenario

MERGE_COMMMIT=`git log -n1 --pretty=format:%H`
if [ -f .git/refs/remotes/origin/$TARGET_BRANCH ]
then
    TARGET_INIT_COMMIT=`cat .git/refs/remotes/origin/$TARGET_BRANCH`
else
    TARGET_INIT_COMMIT=`git log -n1 --pretty=format:%H origin/$TARGET_BRANCH`
fi

echo " ---- Checking the modified files by the merge request ----"
echo ""
echo "Source Branch is    : $SOURCE_BRANCH"
echo "Target Branch is    : $TARGET_BRANCH"
echo "Merged Commit is    : $MERGE_COMMMIT"
echo "Target Init   is    : $TARGET_INIT_COMMIT"
echo ""
echo " ----------------------------------------------------------"
echo ""

# Retrieve the list of modified files since the latest develop commit
MODIFIED_FILES=`git log $TARGET_INIT_COMMIT..$MERGE_COMMMIT --oneline --name-status | egrep "^M|^A" | sed -e "s@^M\t*@@" -e "s@^A\t*@@" | sort | uniq | grep -v test`
NB_TO_FORMAT=0
NB_TOTAL=0

if [ -f oai_rules_result.txt ]
then
    rm -f oai_rules_result.txt
fi
if [ -f oai_rules_result_list.txt ]
then
    rm -f oai_rules_result_list.txt
fi
for FULLFILE in $MODIFIED_FILES
do
    filename=$(basename -- "$FULLFILE")
    EXT="${filename##*.}"
    if [ $EXT = "c" ] || [ $EXT = "h" ] || [ $EXT = "cpp" ] || [ $EXT = "hpp" ]
    then
        SRC_FILE=`echo $FULLFILE | sed -e "s#src/##"`
        TO_FORMAT=`clang-format -output-replacements-xml ${SRC_FILE} 2>&1 | grep -v replacements | grep -c replacement`
        NB_TOTAL=$((NB_TOTAL + 1))
        if [ $TO_FORMAT -ne 0 ]
        then
            NB_TO_FORMAT=$((NB_TO_FORMAT + 1))
            echo $FULLFILE
            echo $FULLFILE >> ./oai_rules_result_list.txt
        fi
    fi
done
echo ""
echo " ----------------------------------------------------------"
echo "Nb Files that do NOT follow OAI rules: $NB_TO_FORMAT over $NB_TOTAL checked!"
echo "NB_FILES_FAILING_CHECK=$NB_TO_FORMAT" > ./oai_rules_result.txt
echo "NB_FILES_CHECKED=$NB_TOTAL" >> ./oai_rules_result.txt

exit 0
