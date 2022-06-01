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
 *
 */

import * as React from 'react';
import Collapse from '@material-ui/core/Collapse';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import {ObjectViewer} from './AlertDetailsPane';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAlarmContext} from '../../AlarmContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {AlertViewerProps} from '../../rules/RuleInterface';

export default function MetricAlertViewer({alert}: AlertViewerProps) {
  const {filterLabels} = useAlarmContext();
  const {labels, annotations} = alert || {};
  const {alertname: _a, severity: _s, ...extraLabels} = labels || {};
  const {description, ...extraAnnotations} = annotations || {};
  const [showDetails, setShowDetails] = React.useState(false);
  return (
    <Grid container data-testid="metric-alert-viewer" spacing={5}>
      <Grid item>
        <Typography variant="body1">{description}</Typography>
      </Grid>
      <Grid item>
        <ObjectViewer
          object={filterLabels ? filterLabels(extraLabels) : extraLabels}
        />
      </Grid>
      <Grid item xs={12}>
        <Link
          variant="subtitle1"
          component="button"
          onClick={() => setShowDetails(!showDetails)}>
          {'Show More Details'}
        </Link>
      </Grid>
      <Grid item>
        <Collapse in={showDetails}>
          <ObjectViewer object={extraAnnotations} />
        </Collapse>
      </Grid>
    </Grid>
  );
}
