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

import AddAlertIcon from '@mui/icons-material/AddAlert';
import AlarmIcon from '@mui/icons-material/Alarm';
import Alarms from './components/Alarms';
import AppContext from '../../context/AppContext';
import Button from '@mui/material/Button';
import ContactMailIcon from '@mui/icons-material/ContactMail';
import Grid from '@mui/material/Grid';
import React from 'react';
import Text from '../../theme/design-system/Text';
import TopBar from '../../components/TopBar';
import nullthrows from '../../../shared/util/nullthrows';
import {MagmaAlarmsApiUtil} from '../../util/AlarmsApiUtil';
import {Theme} from '@mui/material/styles';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import {triggerAlertSync} from '../../util/SyncAlerts';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles<Theme>(theme => ({
  root: {
    padding: theme.spacing(4),
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    float: 'right',
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  emptyAlerts: {
    color: colors.primary.comet,
    marginBottom: '15px',
    width: '50%',
    textAlign: 'center',
  },
}));

const tabs = [
  {
    to: 'alerts',
    label: 'Alerts',
    icon: AlarmIcon,
  },
  {
    to: 'rules',
    label: 'Alert Rules',
    icon: AddAlertIcon,
  },
  {
    to: 'teams',
    label: 'Receivers',
    icon: ContactMailIcon,
  },
];

function EmptyAlerts() {
  const classes = useStyles();
  const isSuperUser = useContext(AppContext).user.isSuperUser;
  const params = useParams();

  const enqueueSnackbar = useEnqueueSnackbar();
  const networkId = nullthrows(params.networkId);
  return (
    <Grid container direction="column" alignItems="center">
      <Text className={classes.emptyAlerts} variant="h5">
        No Alerts Added
      </Text>
      <Text className={classes.emptyAlerts} variant="subtitle1">
        Find out about possible issues in the network by easily enabling alerts
        or creating custom ones.
      </Text>
      {isSuperUser && (
        <Button
          size="large"
          variant="contained"
          className={classes.appBarBtn}
          onClick={() => {
            void triggerAlertSync(networkId, enqueueSnackbar);
          }}>
          {'Enable Alerts'}
        </Button>
      )}
    </Grid>
  );
}

function AlarmsDashboard() {
  const apiUtil = MagmaAlarmsApiUtil;
  const classes = useStyles();

  return (
    <>
      <TopBar header={'Alarms'} tabs={tabs} />
      <div className={classes.root}>
        <Alarms
          makeTabLink={makeTabLink}
          disabledTabs={['alerts', 'rules', 'suppressions', 'routes', 'teams']}
          apiUtil={apiUtil}
          emptyAlerts={<EmptyAlerts />}
        />
      </div>
    </>
  );
}

function makeTabLink(input: {networkId?: string; keyName: string}): string {
  return input.keyName;
}

export default AlarmsDashboard;
