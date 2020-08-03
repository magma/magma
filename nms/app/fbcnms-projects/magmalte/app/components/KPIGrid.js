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
 *
 * @flow strict-local
 * @format
 */

import type {ComponentType} from 'react';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DeviceStatusCircle from '../theme/design-system/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import React from 'react';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  kpiHeaderBlock: {
    display: 'flex',
    alignItems: 'center',
    padding: 0,
  },
  kpiHeaderContent: {
    display: 'flex',
    alignItems: 'center',
  },
  kpiHeaderIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  kpiBlock: {
    boxShadow: `0 0 0 1px ${colors.primary.concrete}`,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: props =>
      props.hasStatus
        ? 'calc(100% - 16px)'
        : props.hasIcon
        ? 'calc(100% - 32px)'
        : '100%',
  },
  kpiObscuredValue: {
    color: colors.primary.brightGray,
    width: '100%',

    '& input': {
      whiteSpace: 'nowrap',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
    },
  },
  kpiBox: {
    width: '100%',
    '& > div': {
      width: '100%',
    },
  },
  kpiIcon: {
    display: 'flex',
    alignItems: 'center',

    '& svg': {
      fill: colors.primary.comet,
      marginRight: theme.spacing(1),
    },
  },
}));

// Status Indicator displays a small text with an DeviceStatusCircle icon
// disabled indicates if the status color is to be grayed out
// up/down indicates if we have to display status to be in green or in red
function StatusIndicator(disabled: boolean, up: boolean, val: string) {
  const props = {hasStatus: true};
  const classes = useStyles(props);
  return (
    <Grid container alignItems="center">
      <Grid item>
        <DeviceStatusCircle isGrey={disabled} isActive={up} isFilled={true} />
      </Grid>
      <Grid item className={classes.kpiValue}>
        {val}
      </Grid>
    </Grid>
  );
}

// KPI Icon adds an icon to the left of the value
function KpiIcon(icon: ComponentType<SvgIconExports>, val: string) {
  const props = {hasIcon: true};
  const classes = useStyles(props);
  const Icon = icon;
  return (
    <Grid container alignItems="center">
      <Grid item className={classes.kpiIcon}>
        <Icon />
      </Grid>
      <Grid item className={classes.kpiValue}>
        {val}
      </Grid>
    </Grid>
  );
}

// KPI Obscure makes the field into a password type filed with a visibility toggle for more sensitive fields.
function KpiObscure(val: number | string) {
  const [showPassword, setShowPassword] = React.useState(false);
  return (
    <Input
      type={showPassword ? 'text' : 'password'}
      fullWidth={true}
      value={val}
      disableUnderline={true}
      readOnly={true}
      data-testid={'epcPassword'} // TODO: TEMPORARY fix for yarn - test
      endAdornment={
        <InputAdornment position="end">
          <IconButton
            aria-label="toggle password visibility"
            onClick={() => setShowPassword(!showPassword)}
            onMouseDown={event => event.preventDefault()}>
            {showPassword ? <Visibility /> : <VisibilityOff />}
          </IconButton>
        </InputAdornment>
      }
    />
  );
}

type KPIData = {
  icon?: ComponentType<SvgIconExports>,
  category?: string,
  value: number | string,
  obscure?: boolean,
  unit?: string,
  statusCircle?: boolean,
  statusInactive?: boolean,
  status?: boolean,
};

export type KPIRows = KPIData[];

type Props = {data: KPIRows[]};

export default function KPIGrid(props: Props) {
  const classes = useStyles();
  const kpiGrid = props.data.map((row, i) => (
    <Grid key={i} container direction="row">
      {row.map((kpi, j) => (
        <Grid
          item
          xs={12}
          md
          key={`data-${i}-${j}`}
          zeroMinWidth
          className={classes.kpiBlock}>
          <Grid container direction="row" alignItems="center">
            <Grid item xs={12}>
              <CardHeader
                data-testid={kpi.category}
                className={classes.kpiBox}
                title={kpi.category}
                titleTypographyProps={{
                  variant: 'body3',
                  className: classes.kpiLabel,
                  title: kpi.category,
                }}
                subheaderTypographyProps={{
                  variant: 'body1',
                  className:
                    kpi.obscure === true
                      ? classes.kpiObscuredValue
                      : classes.kpiValue,
                  title: kpi.value + (kpi.unit ?? ''),
                }}
                subheader={
                  kpi.statusCircle === true
                    ? StatusIndicator(
                        kpi.statusInactive || false,
                        kpi.status || false,
                        kpi.value + (kpi.unit ?? ''),
                      )
                    : kpi.icon
                    ? KpiIcon(kpi.icon, kpi.value + (kpi.unit ?? ''))
                    : kpi.obscure === true
                    ? KpiObscure(kpi.value)
                    : kpi.value + (kpi.unit ?? '')
                }
              />
            </Grid>
          </Grid>
        </Grid>
      ))}
    </Grid>
  ));
  return (
    <Card elevation={0}>
      <Grid container alignItems="center" justify="center">
        {kpiGrid}
      </Grid>
    </Card>
  );
}
