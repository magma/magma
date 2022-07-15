/**
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

import {Request} from 'express';
import {Strategy} from 'passport-strategy';

type StrategyBuilder = (req: Request) => Promise<Strategy>;
type StrategyIDBuilder = (req: Request) => Promise<string>;

export default class DynamicStrategy extends Strategy {
  _strategies: Record<string, Strategy> = {};
  _strategyBuilder: StrategyBuilder;
  _strategyIDBuilder: StrategyIDBuilder;

  constructor(
    strategyIDBuilder: StrategyIDBuilder,
    strategyBuilder: StrategyBuilder,
  ) {
    super();
    this._strategyIDBuilder = strategyIDBuilder;
    this._strategyBuilder = strategyBuilder;
  }

  async _getStrategy(req: Request, name: string): Promise<Strategy> {
    let strategy = this._strategies[name];

    if (!strategy) {
      strategy = this._strategies[name] = await this._strategyBuilder(req);
    }

    strategy.error = this.error;
    strategy.redirect = this.redirect;
    strategy.success = this.success;
    strategy.fail = this.fail;
    strategy.pass = this.pass;

    return strategy;
  }

  authenticate(req: Request, options: any) {
    (async () => {
      const strategyID = await this._strategyIDBuilder(req);
      const strategy = await this._getStrategy(req, strategyID);
      strategy.authenticate(req, options);
    })().catch(error => {
      this.error(error as Error);
    });
  }
}
