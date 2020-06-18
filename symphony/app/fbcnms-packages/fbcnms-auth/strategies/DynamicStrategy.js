/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FBCNMSMiddleWareRequest} from '@fbcnms/express-middleware';

import {Strategy} from 'passport-strategy';

type StrategyBuilder = (req: FBCNMSMiddleWareRequest) => Promise<Strategy>;
type StrategyIDBuilder = (req: FBCNMSMiddleWareRequest) => Promise<string>;

export default class DynamicStrategy extends Strategy {
  _strategies: {[string]: Strategy} = {};
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

  async _getStrategy(req: FBCNMSMiddleWareRequest, name: string): Strategy {
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

  authenticate(req: FBCNMSMiddleWareRequest, options: any) {
    (async () => {
      const strategyID = await this._strategyIDBuilder(req);
      const strategy = await this._getStrategy(req, strategyID);
      strategy.authenticate(req, options);
    })().catch(_error => {
      this.error();
    });
  }
}
