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
import AppContext from '../../../app/components/context/AppContext';
import Grid from '@material-ui/core/Grid';
import LteMetrics from '../lte/LteMetrics';
import Paper from '@material-ui/core/Paper';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  tab: {
    backgroundColor: colors.primary.white,
    borderRadius: '4px 4px 0 0',
    boxShadow: `inset 0 -2px 0 0 ${colors.primary.concrete}`,
    '& + &': {
      marginLeft: '4px',
    },
  },
  emptyTable: {
    backgroundColor: colors.primary.white,
    padding: theme.spacing(4),
    minHeight: '96px',
  },
  emptyTableContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    color: colors.primary.comet,
  },
}));

export default function () {
  const classes = useStyles();
  const grafanaEnabled =
    useContext(AppContext).isFeatureEnabled('grafana_metrics') &&
    useContext(AppContext).user.isSuperUser;

  return (
    <>
      {grafanaEnabled ? (
        <LteMetrics />
      ) : (
        <Paper elevation={0}>
          <Grid
            container
            alignItems="center"
            justifyContent="center"
            className={classes.emptyTable}>
            <Grid item xs={12} className={classes.emptyTableContent}>
              <Text variant="body2">
                Metrics is only enabled when grafana feature is enabled and user
                is a superuser
              </Text>
            </Grid>
          </Grid>
        </Paper>
      )}
    </>
  );
}
