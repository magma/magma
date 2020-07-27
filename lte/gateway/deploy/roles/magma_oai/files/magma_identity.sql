/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

INSERT INTO mmeidentity (`idmmeidentity`,`mmehost`,`mmerealm`,`UE-reachability`)
SELECT * FROM (SELECT '7','magma-oai.openair4G.eur','openair4G.eur','0') AS tmp
WHERE NOT EXISTS (
  SELECT * FROM mmeidentity
  WHERE mmehost = 'magma-oai.openair4G.eur'
) LIMIT 1;
