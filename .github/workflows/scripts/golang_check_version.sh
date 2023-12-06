#!/usr/bin/env bash
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

expected_version="1.20.1"
all_versions_good=true

while IFS= read -r -d '' file
do
  while read -r line;
  do
    version=$(echo "$line" | awk -F '[]["'"'"']' '{print $2}')
    if [[ -n $version && $version != "$expected_version" ]]
    then
      echo "Found unexpected Go version $version in file $(realpath --relative-to="$MAGMA_ROOT" "$file"):"
      echo "$line"
      all_versions_good=false
    fi
  done < <(grep -i 'go[-_]version:' "$file" )

  while read -r line;
  do
    version=$(echo "$line" | awk '{ print $NF }')
    if [[ -n $version && $version != "go$expected_version.linux-amd64.tar.gz" ]]
    then
      echo "Found unexpected Go version $version in file $(realpath --relative-to="$MAGMA_ROOT" "$file"):"
      echo "$line"
      all_versions_good=false
    fi
  done < <(grep -i 'golang_tar:' "$file" )
done < <(find "$MAGMA_ROOT" -regex '.*\.ya?ml' -print0)

while IFS= read -r -d '' file
do
  while read -r line;
  do
    version=$(echo "$line" | awk -F '"' '{print $2}')
    if [[ -n $version && $version != "$expected_version" && $version != "1.21.0" ]]
    then
      echo "Found unexpected Go version $version in file $(realpath --relative-to="$MAGMA_ROOT" "$file"):"
      echo "$line"
      all_versions_good=false
    fi
  done < <(grep -i '^ARG GOLANG_VERSION' "$file" )
done < <(find "$MAGMA_ROOT" -name Dockerfile -print0)

while IFS= read -r -d '' file
do
  while read -r version;
  do
    if [[ -n $version && $version != "go$expected_version" ]]
    then
      echo "Found unexpected Go version $version in file $(realpath --relative-to="$MAGMA_ROOT" "$file"):"
      all_versions_good=false
    fi
  done < <(grep -o -P 'go[\d\.]+\d' "$file" )
done < <(find "$MAGMA_ROOT/docs/readmes/basics/" -regex '.*\.md' -print0)

if [ $all_versions_good = true ]
then
  exit 0
fi
exit 1
