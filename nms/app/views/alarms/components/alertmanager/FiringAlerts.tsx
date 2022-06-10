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
import AddAlertTwoToneIcon from '@material-ui/icons/AddAlertTwoTone';
import AlertDetailsPane from './AlertDetails/AlertDetailsPane';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SeverityIndicator from '../severity/SeverityIndicator';
import SimpleTable from '../table/SimpleTable';
import Slide from '@material-ui/core/Slide';
import Typography from '@material-ui/core/Typography';
import moment from 'moment';
import {Link, useResolvedPath} from 'react-router-dom';
import {colors} from '../../../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useEffect, useState} from 'react';
import {useNetworkId} from '../hooks';
import {useSnackbars} from '../../../../hooks/useSnackbar';

import {Theme} from '@material-ui/core/styles';
import {getErrorMessage} from '../../../../util/ErrorUtils';
import type {PromFiringAlert} from '../../../../../generated-ts';

const useStyles = makeStyles<Theme>(theme => ({
  root: {
    paddingTop: theme.spacing(4),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
  addAlertIcon: {
    fontSize: '200px',
    margin: theme.spacing(1),
  },
  helperText: {
    color: colors.primary.brightGray,
    fontSize: theme.typography.pxToRem(20),
  },
}));

type Props = {
  emptyAlerts?: React.ReactNode;
};

export default function FiringAlerts(props: Props) {
  const resolvedPath = useResolvedPath('');
  const {apiUtil, filterLabels} = useAlarmContext();
  const [selectedRow, setSelectedRow] = useState<PromFiringAlert | null>(null);
  const [lastRefreshTime] = useState<string>(new Date().toLocaleString());
  const [alertData, setAlertData] = useState<Array<PromFiringAlert> | null>(
    null,
  );
  const classes = useStyles();
  const snackbars = useSnackbars();
  const networkId = useNetworkId();
  const {error, isLoading, response} = apiUtil.useAlarmsApi(
    apiUtil.viewFiringAlerts,
    {networkId},
    lastRefreshTime,
  );

  useEffect(() => {
    if (!isLoading) {
      const alertData = response
        ? response.map(alert => {
            let labels = alert.labels;
            if (labels && filterLabels) {
              labels = filterLabels(labels);
            }
            return {
              ...alert,
              labels,
            };
          })
        : [];
      setAlertData(alertData);
    }
    return () => {
      setAlertData(null);
    };
  }, [filterLabels, isLoading, response, setAlertData]);

  const showRowDetailsPane = React.useCallback(
    (row: PromFiringAlert) => {
      setSelectedRow(row);
    },
    [setSelectedRow],
  );
  const hideDetailsPane = React.useCallback(() => {
    setSelectedRow(null);
  }, [setSelectedRow]);

  React.useEffect(() => {
    if (error) {
      snackbars.error(
        `Unable to load firing alerts. ${getErrorMessage(error)}`,
      );
    }
  }, [error, snackbars]);

  if (!isLoading && alertData?.length === 0) {
    return (
      <Grid
        container
        spacing={2}
        direction="column"
        alignItems="center"
        justifyContent="center"
        data-testid="no-alerts-icon"
        style={{minHeight: '60vh'}}>
        {!(props.emptyAlerts ?? false) ? (
          <>
            <Grid item>
              <AddAlertTwoToneIcon
                color="primary"
                className={classes.addAlertIcon}
              />
            </Grid>
            <Grid item>
              <span className={classes.helperText}>
                Start creating alert rules
              </span>
            </Grid>
            <Grid item>
              <Button
                color="primary"
                size="small"
                variant="contained"
                component={Link}
                to={`${resolvedPath.pathname.slice(
                  0,
                  resolvedPath.pathname.lastIndexOf('/'),
                )}/rules`}>
                Add Alert Rule
              </Button>
            </Grid>
          </>
        ) : (
          <>{props.emptyAlerts}</>
        )}
      </Grid>
    );
  }
  return (
    <Grid className={classes.root} container spacing={2}>
      <Grid item xs={selectedRow ? 8 : 12}>
        <SimpleTable
          onRowClick={showRowDetailsPane}
          columnStruct={[
            {
              title: 'Name',
              field: 'labels.alertname',
              render: currRow => (
                <Typography noWrap>
                  <span>{currRow.labels?.alertname}</span>
                </Typography>
              ),
            },
            {
              title: 'Severity',
              field: 'labels.severity',
              render: currRow => (
                <SeverityIndicator severity={currRow.labels?.severity} />
              ),
            },
            {
              title: 'Date',
              field: 'startsAt',
              render: currRow => {
                const date = moment(new Date(currRow.startsAt));
                return (
                  <>
                    <Typography variant="body1">{date.fromNow()}</Typography>
                    <div>{date.format('dddd, MMMM Do YYYY')}</div>
                  </>
                );
              },
            },
          ]}
          tableData={alertData || []}
          dataTestId="firing-alerts"
        />
        {isLoading && (
          <div className={classes.loading}>
            <CircularProgress />
          </div>
        )}
      </Grid>
      <Slide direction="left" in={!!selectedRow}>
        <Grid item xs={4}>
          {selectedRow && (
            <AlertDetailsPane alert={selectedRow} onClose={hideDetailsPane} />
          )}
        </Grid>
      </Slide>
    </Grid>
  );
}
