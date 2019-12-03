/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Admin from './admin/Admin';
import Automation from './automation/Automation';
import IDToolMain from './id/IDToolMain';
import Inventory from './Inventory';
import MagmaMain from '@fbcnms/magmalte/app/components/Main';
import React from 'react';
import WorkOrdersMain from './work_orders/WorkOrdersMain';

import {Route, Switch} from 'react-router-dom';

export default () => (
  <Switch>
    <Route path="/nms" component={MagmaMain} />
    <Route path="/inventory" component={Inventory} />
    <Route path="/workorders" component={WorkOrdersMain} />
    <Route path="/admin" component={Admin} />
    <Route path="/automation" component={Automation} />
    <Route path="/id" component={IDToolMain} />
  </Switch>
);
