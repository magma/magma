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
import Button from '@material-ui/core/Button';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import RemoveCircleIcon from '@material-ui/icons/RemoveCircle';
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
import type {prometheus_labelset} from '@fbcnms/magma-api';

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
  metricFilterItem: {
    marginRight: theme.spacing(1),
  },
}));

export type ThresholdExpression = {
  metricName: string,
  comparator: ?Comparator,
  filters: Array<{name: string, value: string}>,
  value: number,
};

export function thresholdToPromQL(
  thresholdExpression: ThresholdExpression,
): string {
  if (!thresholdExpression.comparator || !thresholdExpression.metricName) {
    return '';
  }
  let filteredMetric = thresholdExpression.metricName;
  if (thresholdExpression.filters.length > 0) {
    filteredMetric += '{';
    thresholdExpression.filters.forEach(filter => {
      if (filter.name != '' && filter.value != '') {
        filteredMetric += `${filter.name}=~"^${filter.value}$",`;
      }
    });
    filteredMetric += '}';
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
    <>
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
      <Grid item>
        {props.expression.filters.length > 0 ? (
          <FormLabel>For metrics matching:</FormLabel>
        ) : (
          <></>
        )}
        <MetricFilters
          metricSeries={props.metricsByName[props.expression.metricName]}
          expression={props.expression}
          onChange={props.onChange}
        />
      </Grid>
    </>
  );
}

function MetricFilters(props: {
  metricSeries: Array<prometheus_labelset>,
  expression: ThresholdExpression,
  onChange: (expression: ThresholdExpression) => void,
}) {
  const classes = useStyles();
  return (
    <Grid container>
      {props.expression.filters.map((filter, idx) => (
        <Grid item xs={12}>
          <MetricFilter
            key={idx}
            metricSeries={props.metricSeries}
            onChange={props.onChange}
            onRemove={filterIdx => {
              const filters = props.expression.filters;
              filters.splice(filterIdx, 1);
              props.onChange({...props.expression, filters: filters});
            }}
            expression={props.expression}
            filterIdx={idx}
            selectedLabel={filter.name}
            selectedValue={filter.value}
          />
        </Grid>
      ))}
      <Grid item>
        <Button
          disabled={props.expression.metricName == ''}
          variant="contained"
          className={classes.button}
          onClick={() => {
            props.onChange({
              ...props.expression,
              filters: [...props.expression.filters, {name: '', value: ''}],
            });
          }}>
          Add Filter
        </Button>
      </Grid>
    </Grid>
  );
}

function MetricFilter(props: {
  metricSeries: Array<prometheus_labelset>,
  onChange: (expression: ThresholdExpression) => void,
  onRemove: (filerIdx: number) => void,
  expression: ThresholdExpression,
  filterIdx: number,
  selectedLabel: string,
  selectedValue: string,
}) {
  return (
    <>
      <FilterSelector
        values={getFilteredListOfLabelNames([
          ...new Set(...props.metricSeries.map(Object.keys)),
        ])}
        defaultVal="Label"
        onChange={({target}) => {
          const filters = props.expression.filters;
          filters[props.filterIdx].name = target.value;
          props.onChange({...props.expression, filters: filters});
        }}
        selectedValue={props.selectedLabel}
      />
      <FilterSelector
        values={[
          ...new Set(props.metricSeries.map(item => item[props.selectedLabel])),
        ]}
        disabled={props.selectedLabel == ''}
        defaultVal="Value"
        onChange={({target}) => {
          const filters = props.expression.filters;
          filters[props.filterIdx].value = target.value;
          props.onChange({...props.expression, filters: filters});
        }}
        updateExpression={props.onChange}
        selectedValue={props.selectedValue}
      />
      <IconButton onClick={() => props.onRemove(props.filterIdx)}>
        <RemoveCircleIcon />
      </IconButton>
    </>
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

function FilterSelector(props: {
  values: Array<string>,
  defaultVal: string,
  onChange: (event: SyntheticInputEvent<HTMLElement>) => void,
  selectedValue: string,
  disabled?: boolean,
}) {
  const classes = useStyles();
  const menuItems = props.values.map(val => (
    <MenuItem value={val} key={val}>
      {val}
    </MenuItem>
  ));

  return (
    <Select
      disabled={props.disabled}
      displayEmpty
      className={classes.metricFilterItem}
      value={props.selectedValue}
      onChange={props.onChange}>
      <MenuItem disabled value="">
        <em>{props.defaultVal}</em>
      </MenuItem>
      {menuItems}
    </Select>
  );
}

// Labels we don't want to show during metric filtering since they are useless
const forbiddenLabels = new Set(['networkID', '__name__']);
function getFilteredListOfLabelNames(labelNames: Array<string>): Array<string> {
  return labelNames.filter(label => !forbiddenLabels.has(label));
}
