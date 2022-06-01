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
 * @flow
 * @format
 *
 * Base container for showing details of different types of alerts. To show
 * a custom component for an alert type, 2 interfaces must be implemented:
 *  Implement the getAlertType in AlarmContext. This function should
 *  inspect the labels/annotations of an alert and determine which rule type
 *  generated it.
 *
 *  Implement the AlertViewer interface for the rule type. By default, the
 *  MetricAlertViewer will be shown.
 */

import * as React from 'react';
import CloseIcon from '@material-ui/icons/Close';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Link from '@material-ui/core/Link';
import MetricAlertViewer from './MetricAlertViewer';
import Paper from '@material-ui/core/Paper';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SeverityIndicator from '../../severity/SeverityIndicator';
import Typography from '@material-ui/core/Typography';
import moment from 'moment';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../../AlarmContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useSnackbars} from '../../../../../hooks/useSnackbar';

import type {
  AlertViewerProps,
  RuleInterfaceMap,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from '../../rules/RuleInterface';
// $FlowFixMe migrated to typescript
import type {FiringAlarm} from '../../AlarmAPIType';
// $FlowFixMe migrated to typescript
import type {Labels} from '../../AlarmAPIType';
import type {SvgIconProps} from '@material-ui/core/@@SvgIcon';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(3),
  },
  capitalize: {
    textTransform: 'capitalize',
  },
  // annotations can potentially contain json so it should wrap properly
  objectViewerValue: {
    wordBreak: 'break-word',
  },
  objectViewerItem: {
    marginBottom: '0',
    justifyContent: 'space-between',
  },
}));

type Props = {|
  alert: FiringAlarm,
  onClose: () => void,
|};

export default function AlertDetailsPane({alert, onClose}: Props) {
  const classes = useStyles();
  const {getAlertType, ruleMap} = useAlarmContext();
  const alertType = getAlertType ? getAlertType(alert) : '';
  const {startsAt, labels} = alert || {};
  const {alertname, severity} = labels || {};
  const AlertViewer = getAlertViewer(ruleMap, alertType);

  return (
    <Paper elevation={1} data-testid="alert-details-pane">
      <Grid container direction="column" spacing={2} className={classes.root}>
        <Grid item container direction="column" spacing={1}>
          <Grid item container justifyContent="space-between">
            <Grid item>
              <SeverityIndicator severity={severity} chip={true} />
            </Grid>
            <Grid item>
              <IconButton
                size="small"
                edge="end"
                onClick={onClose}
                data-testid="alert-details-close">
                <CloseIcon />
              </IconButton>
            </Grid>
          </Grid>
          <Grid item>
            <Typography variant="h5">{alertname}</Typography>
          </Grid>
          <Grid item>
            <AlertDate date={startsAt} />
          </Grid>
        </Grid>
        <Grid item>
          <AlertViewer alert={alert} />
        </Grid>
      </Grid>
      <AlertTroubleshootingLink alertName={alertname} />
    </Paper>
  );
}

/**
 * Get the AlertViewer for this alert's rule type or fallback to the default.
 */
function getAlertViewer(
  ruleMap: RuleInterfaceMap<mixed>,
  alertType: string,
): React.ComponentType<AlertViewerProps> {
  const ruleInterface = ruleMap[alertType];
  if (!(ruleInterface && ruleInterface.AlertViewer)) {
    return MetricAlertViewer;
  }
  return ruleInterface.AlertViewer;
}

function AlertDate({date}: {date: string}) {
  const classes = useStyles();
  const fromNow = React.useMemo(() => moment(date).local().fromNow(), [date]);
  const startDate = React.useMemo(
    () => moment(date).local().format('MMM Do YYYY, h:mm:ss a'),
    [date],
  );
  return (
    <Typography variant="body2" color="textSecondary">
      <span className={classes.capitalize}>{fromNow}</span> | {startDate}
    </Typography>
  );
}

/**
 * Link to troubleshooting documentation or display nothing if no link provided
 */
function AlertTroubleshootingLink({alertName}: {alertName: string}) {
  const classes = useStyles();
  const snackbars = useSnackbars();
  const {apiUtil} = useAlarmContext();
  const {error, response: troubleshootingLink} = apiUtil.useAlarmsApi(
    apiUtil.getTroubleshootingLink,
    {
      alertName,
    },
  );
  React.useEffect(() => {
    if (error) {
      snackbars.error(
        `Unable to load troubleshooting link. ${
          error.response ? error.response.data.message : error.message || ''
        }`,
      );
    }
  }, [error, snackbars]);
  return (
    <>
      {(troubleshootingLink?.link || '').length > 0 && (
        <>
          <Divider variant="fullWidth" />
          <Grid
            container
            direction="column"
            spacing={2}
            className={classes.root}>
            <Grid item>
              <Link
                variant="subtitle1"
                href={troubleshootingLink?.link}
                target="_blank"
                rel="noopener">
                {troubleshootingLink?.title}
              </Link>
            </Grid>
          </Grid>
        </>
      )}
    </>
  );
}

/**
 * Shows the key-value pairs of an object such as annotations or labels.
 */
export function ObjectViewer({object}: {object: Labels}) {
  const labelKeys = Object.keys(object);
  const classes = useStyles();
  return (
    <Grid container item>
      {labelKeys.length < 1 && (
        <Grid item>
          <Typography color="textSecondary">None</Typography>
        </Grid>
      )}
      {labelKeys.map(key => (
        <Grid container item spacing={4} className={classes.objectViewerItem}>
          <Grid item>
            <Typography variant="subtitle1">{key}:</Typography>
          </Grid>
          <Grid item>
            <Typography
              className={classes.objectViewerValue}
              color="textSecondary"
              variant="subtitle1">
              {object[key]}
            </Typography>
          </Grid>
        </Grid>
      ))}
    </Grid>
  );
}

export function Section({
  title,
  children,
  divider,
}: {
  title: React.Node,
  children: React.Node,
  /**
   * we shouldn't show a divider for the last section. Only hide if false is
   * passed
   */
  divider?: boolean,
}) {
  return (
    <Grid item container direction="column" spacing={2}>
      <Grid item>
        <Typography variant="h5">{title}</Typography>
      </Grid>
      <Grid item container spacing={2}>
        {children}
      </Grid>
      {divider !== false && (
        <Grid item>
          <Divider />
        </Grid>
      )}
    </Grid>
  );
}

const useDetailIconStyles = makeStyles(_theme => ({
  root: {
    verticalAlign: 'middle',
    fontSize: '1rem',
  },
}));
// layout for items in the Details section
export function Detail({
  icon: Icon,
  title,
  children,
}: {
  icon: React.ComponentType<SvgIconProps>,
  title: string,
  children: React.Node,
}) {
  const iconStyles: $Shape<$Values<SvgIconProps>> = useDetailIconStyles();
  return (
    <Grid item container wrap="nowrap" spacing={1}>
      <Grid item>
        <Icon classes={iconStyles} fontSize="small" />
      </Grid>
      <Grid container direction="column" item>
        <Grid item>
          <Typography variant="body1">{title}</Typography>
        </Grid>
        <Grid item>
          <Typography color="textSecondary">{children}</Typography>
        </Grid>
      </Grid>
    </Grid>
  );
}
