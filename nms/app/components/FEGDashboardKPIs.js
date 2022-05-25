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
 *
 * @flow strict-local
 * @format
 */

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from './layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ServicingAccessGatewayKPIs from './FEGServicingAccessGatewayKPIs';

import {GpsFixed} from '@material-ui/icons';

/**
 * Returns the KPI's in the federation network dashboard.
 * It currently supports KPI for the count of serviced access gateways
 * by the federation network.
 */
export default function () {
  return (
    <>
      <CardTitleRow icon={GpsFixed} label="Events" />
      <Grid container item zeroMinWidth alignItems="center" spacing={4}>
        <Grid item xs={12} lg={6}>
          <ServicingAccessGatewayKPIs />
        </Grid>
      </Grid>
    </>
  );
}
