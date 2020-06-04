/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as PromQL from '../../prometheus/PromQL';
import * as React from 'react';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Button from '@material-ui/core/Button';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import RemoveCircleIcon from '@material-ui/icons/RemoveCircle';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import useRouter from '../../../hooks/useRouter';
import {groupBy} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../../AlarmContext';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';

import type {InputChangeFunc} from './PrometheusEditor';

type prometheus_labelset = {
  [string]: string,
};

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
  comparator: PromQL.BinaryComparator,
  filters: PromQL.Labels,
  value: number,
};

export function thresholdToPromQL(
  thresholdExpression: ThresholdExpression,
): string {
  if (!thresholdExpression.comparator || !thresholdExpression.metricName) {
    return '';
  }
  const {metricName, comparator, filters, value} = thresholdExpression;
  const metricSelector = new PromQL.InstantSelector(metricName, filters);
  const exp = new PromQL.BinaryOperation(
    metricSelector,
    new PromQL.Scalar(value),
    comparator,
  );
  return exp.toPromQL();
}

export default function ToggleableExpressionEditor(props: {
  onChange: InputChangeFunc,
  onThresholdExpressionChange: (expresion: ThresholdExpression) => void,
  expression: ThresholdExpression,
  stringExpression: string,
  toggleOn: boolean,
  onToggleChange: boolean => void,
}) {
  const {apiUtil} = useAlarmContext();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {response, error} = apiUtil.useAlarmsApi(apiUtil.getMetricSeries, {
    networkId: match.params.networkId,
  });
  if (error) {
    enqueueSnackbar('Error retrieving metrics: ' + error, {
      variant: 'error',
    });
  }
  const metricsByName = groupBy(response, '__name__');

  return (
    <Grid container item xs={12}>
      <Grid item xs={12}>
        <ToggleSwitch
          toggleOn={props.toggleOn}
          onChange={({target}) => props.onToggleChange(target.checked)}
        />
      </Grid>
      {props.toggleOn ? (
        <Grid item xs={12}>
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
      {props.metricsByName[props.expression.metricName] ? (
        ''
      ) : (
        <MenuItem
          key={props.expression.metricName}
          value={props.expression.metricName}>
          {props.expression.metricName}
        </MenuItem>
      )}
    </Select>
  );

  return (
    <>
      <Grid container spacing={2} alignItems="center">
        <Grid item>
          <Typography variant="body2">IF</Typography>
        </Grid>
        <Grid item>{metricSelector}</Grid>
        <Grid item>
          <Typography variant="body2">IS</Typography>
        </Grid>
        <Grid item>
          <ToggleButtonGroup
            exclusive={true}
            value={props.expression.comparator}
            onChange={(event, val) => {
              props.onChange({
                ...props.expression,
                comparator: new PromQL.BinaryComparator(val),
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
            onChange={({target}) => {
              props.onChange({
                ...props.expression,
                value: parseFloat(target.value),
              });
            }}
          />
        </Grid>
      </Grid>
      <Grid item xs={12}>
        {props.expression.filters.len() > 0 ? (
          <FormLabel>For metrics matching:</FormLabel>
        ) : (
          <></>
        )}
        <MetricFilters
          metricSeries={props.metricsByName[props.expression?.metricName] || []}
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
      {props.expression.filters.labels.map((filter, idx) => (
        <Grid item xs={12}>
          <MetricFilter
            key={idx}
            metricSeries={props.metricSeries}
            onChange={props.onChange}
            onRemove={filterIdx => {
              const filtersCopy = props.expression.filters.copy();
              filtersCopy.remove(filterIdx);
              props.onChange({...props.expression, filters: filtersCopy});
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
            const filtersCopy = props.expression.filters.copy();
            filtersCopy.addEqual('', '');
            props.onChange({
              ...props.expression,
              filters: filtersCopy,
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
  const labelNames: Array<string> = [];
  props.metricSeries.forEach(metric => {
    labelNames.push(...Object.keys(metric));
  });
  labelNames.push(props.selectedLabel);

  return (
    <Grid container xs={12} spacing={1} alignItems="center">
      <Grid item>
        <FilterSelector
          values={getFilteredListOfLabelNames([...new Set(labelNames)])}
          defaultVal="Label"
          onChange={({target}) => {
            const filtersCopy = props.expression.filters.copy();
            filtersCopy.setIndex(props.filterIdx, target.value, '');
            props.onChange({...props.expression, filters: filtersCopy});
          }}
          selectedValue={props.selectedLabel}
        />
      </Grid>
      <Grid item xs={3}>
        <FilterAutocomplete
          values={
            props.selectedLabel
              ? [
                  ...new Set(
                    props.metricSeries.map(item => item[props.selectedLabel]),
                  ),
                ]
              : []
          }
          disabled={props.selectedLabel == ''}
          defaultVal="Value"
          onChange={(event, value) => {
            // TODO: This is here because we have to pass the onChange function
            // to both the Autocomplete element and the TextInput element
            // T57876329
            if (!value) {
              value = event.target.value;
            }
            const filtersCopy = props.expression.filters.copy();
            const filterOperator = isRegexValue(value) ? '=~' : '=';
            filtersCopy.setIndex(
              props.filterIdx,
              filtersCopy.labels[props.filterIdx].name,
              value || '',
              filterOperator,
            );
            props.onChange({...props.expression, filters: filtersCopy});
          }}
          updateExpression={props.onChange}
          selectedValue={props.selectedValue}
        />
      </Grid>
      <Grid item>
        <IconButton onClick={() => props.onRemove(props.filterIdx)}>
          <RemoveCircleIcon />
        </IconButton>
      </Grid>
    </Grid>
  );
}

function ToggleSwitch(props: {
  toggleOn: boolean,
  onChange: (event: SyntheticInputEvent<HTMLInputElement>) => void,
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

function FilterAutocomplete(props: {
  values: Array<string>,
  defaultVal: string,
  onChange: (event: SyntheticInputEvent<HTMLElement>) => void,
  selectedValue: string,
  disabled?: boolean,
}) {
  return (
    <Autocomplete
      freeSolo
      options={props.values}
      onChange={props.onChange}
      value={props.selectedValue}
      renderInput={({inputProps, ...params}) => (
        <TextField
          {...params}
          inputProps={{
            ...inputProps,
            autoComplete: 'off',
            onChange: props.onChange,
          }}
          disabled={props.values.length === 0}
          label={props.defaultVal}
          margin="normal"
          variant="filled"
          fullWidth
        />
      )}
    />
  );
}

// Labels we don't want to show during metric filtering since they are useless
const forbiddenLabels = new Set(['networkID', '__name__']);
function getFilteredListOfLabelNames(labelNames: Array<string>): Array<string> {
  return labelNames.filter(label => !forbiddenLabels.has(label));
}

// Checks if a value has regex characters
function isRegexValue(value: string): boolean {
  const regexChars = '.+*|?()[]{}:=';
  for (const char of regexChars.split('')) {
    if (value.includes(char)) {
      return true;
    }
  }
  return false;
}
