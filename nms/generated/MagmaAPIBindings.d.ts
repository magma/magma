/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

class MagmaAPIBindings {
  public static async getNetworksByNetworkIdType(parameters: {
    networkId: string;
  }): Promise<string>;

  public static async deleteLteByNetworkIdApnsByApnName(parameters: {
    networkId: string;
    apnName: string;
  }): Promise<'Success'>;

  public static async getNetworks(): Promise<Array<string>>;
}

export default MagmaAPIBindings;
