/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import Tooltip from '@material-ui/core/Tooltip';
import {groupBy} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks/index';

import type {ApiUtil} from './AlarmsApi';
import type {InputChangeFunc} from './PrometheusEditor';
import type {prometheus_labelset} from '../../fbcnms-magma-api';

const useStyles = makeStyles(theme => ({
  button: {
    marginRight: theme.spacing(1),
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  helpButton: {
    color: 'black',
  },
  labeledToggleSwitch: {
    paddingBottom: 0,
  },
}));

export type ThresholdExpression = {
  metricName: string,
  comparator: ?Comparator,
  // TODO: Add support for filters
  filter: {name: string, value: string},
  value: number,
};

export function thresholdToPromQL(
  thresholdExpression: ThresholdExpression,
): string {
  if (!thresholdExpression.comparator || !thresholdExpression.metricName) {
    return '';
  }
  let filteredMetric = thresholdExpression.metricName;
  if (thresholdExpression.filter.name && thresholdExpression.filter.value) {
    const filter = thresholdExpression.filter;
    filteredMetric += `{${filter.name}=~"^${filter.value}"$}`;
  }
  return (
    filteredMetric + thresholdExpression.comparator + thresholdExpression.value
  );
}

const COMPARATORS = {
  '<': '<',
  '<=': '<=',
  '=': '==',
  '==': '==',
  '>=': '>=',
  '>': '>',
};
type Comparator = $Keys<typeof COMPARATORS>;

function getComparator(newComparator: string): Comparator {
  return COMPARATORS[newComparator];
}

export default function ToggleableExpressionEditor(props: {
  onChange: InputChangeFunc,
  onThresholdExpressionChange: (expresion: ThresholdExpression) => void,
  expression: ThresholdExpression,
  stringExpression: string,
  apiUtil: ApiUtil,
}) {
  const [toggleOn, setToggleOn] = React.useState(false);

  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {response, error} = props.apiUtil.useAlarmsApi(
    props.apiUtil.getMetricSeries,
    {networkId: match.params.networkId},
  );
  if (error) {
    enqueueSnackbar('Error retrieving metrics: ' + error);
  }
  const metricsByName = groupBy(response, '__name__');

  return (
    <Grid container item xs={12}>
      <Grid item xs={12}>
        <ToggleSwitch
          toggleOn={toggleOn}
          onChange={({target}) => setToggleOn(target.checked)}
        />
      </Grid>
      {toggleOn ? (
        <Grid item xs={4}>
          <AdvancedExpressionEditor
            onChange={props.onChange}
            expression={props.stringExpression}
          />
        </Grid>
      ) : (
        <ThresholdExpressionEditor
          onChange={props.onThresholdExpressionChange}
          expression={props.expression}
          metricsByName={metricsByName}
        />
      )}
    </Grid>
  );
}

export function AdvancedExpressionEditor(props: {
  onChange: InputChangeFunc,
  expression: string,
}) {
  return (
    <TextField
      required
      label="Expression"
      placeholder="Ex: up == 0"
      value={props.expression}
      onChange={props.onChange(value => ({expression: value}))}
      fullWidth
    />
  );
}

function ThresholdExpressionEditor(props: {
  onChange: (expression: ThresholdExpression) => void,
  expression: ThresholdExpression,
  metricsByName: {[string]: Array<prometheus_labelset>},
}) {
  const metricSelector = (
    <Select
      displayEmpty
      value={props.expression.metricName}
      onChange={({target}) => {
        props.onChange({...props.expression, metricName: target.value});
      }}>
      <MenuItem disabled value="">
        <em>Metric Name</em>
      </MenuItem>
      {Object.keys(props.metricsByName).map(item => (
        <MenuItem key={item} value={item}>
          {item}
        </MenuItem>
      ))}
    </Select>
  );

  return (
    <Grid container spacing={2} alignItems="center">
      <Grid item>
        <Text>IF</Text>
      </Grid>
      <Grid item>{metricSelector}</Grid>
      <Grid item>
        <Text>IS</Text>
      </Grid>
      <Grid item>
        <ToggleButtonGroup
          exclusive={true}
          value={props.expression.comparator}
          onChange={(event, val) => {
            props.onChange({
              ...props.expression,
              comparator: getComparator(val),
            });
          }}>
          <ToggleButton value="<">{'<'}</ToggleButton>
          <ToggleButton value="<=">{'≤'}</ToggleButton>
          <ToggleButton value="==">{'='}</ToggleButton>
          <ToggleButton value=">=">{'≥'}</ToggleButton>
          <ToggleButton value=">">{'>'}</ToggleButton>
        </ToggleButtonGroup>
      </Grid>
      <Grid item xs={2}>
        <TextField
          value={props.expression.value}
          type="number"
          InputLabelProps={{
            shrink: true,
          }}
          onChange={({target}) => {
            props.onChange({
              ...props.expression,
              value: parseFloat(target.value),
            });
          }}
          margin="normal"
        />
      </Grid>
    </Grid>
  );
}

function ToggleSwitch(props: {
  toggleOn: boolean,
  onChange: (event: any) => void,
}) {
  const classes = useStyles();
  return (
    <>
      <FormLabel>Advanced Expression</FormLabel>
      <Switch onChange={props.onChange} checked={props.toggleOn} />
      <Tooltip
        title={
          'Switch the toggle on to write an arbitrary alerting expression' +
          'in PromQL.\n' +
          'To learn more about how to write alert expressions, click ' +
          'on the help icon to open the prometheus querying basics guide.'
        }
        placement="right">
        <IconButton
          className={classes.helpButton}
          href="https://prometheus.io/docs/prometheus/latest/querying/basics/"
          target="_blank"
          size="small">
          <HelpIcon />
        </IconButton>
      </Tooltip>
    </>
  );
}
