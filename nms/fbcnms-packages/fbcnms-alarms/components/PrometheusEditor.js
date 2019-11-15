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
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import {makeStyles} from '@material-ui/styles';

//TODO move to more shared location
import {SEVERITY} from './SimpleTable';

import type {AlertConfig} from './AlarmAPIType';

type Props = {
  rule: ?AlertConfig,
  updateAlertConfig: AlertConfig => void,
  saveAlertRule: () => Promise<void>,
  onExit: () => void,
  isNew: boolean,
};

type MenuItemProps = {key: string, value: string, children: string};

type TimeUnit = {value: string, label: string};

const timeUnits: Array<TimeUnit> = [
  {
    value: '',
    label: '',
  },
  {
    value: 's',
    label: 'seconds',
  },
  {
    value: 'm',
    label: 'minutes',
  },
  {
    value: 'h',
    label: 'hours',
  },
];

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
}));

/**
 * An easier to edit representation of the form's state, then convert
 * to and from the AlertConfig type for posting to the api.
 */
type FormState = {
  ruleName: string,
  expression: string,
  severity: string,
  timeNumber: string,
  timeUnit: string,
};

type inputChangeFunc = (
  formUpdate: (val: string) => $Shape<FormState>,
) => (event: SyntheticInputEvent<HTMLElement>) => void;

function RuleNameEditor(props: {onChange: inputChangeFunc, ruleName: string}) {
  return (
    <TextField
      required
      label="Rule Name"
      placeholder="Ex: Service Down"
      value={props.ruleName}
      onChange={props.onChange(value => ({ruleName: value}))}
      fullWidth
    />
  );
}

function AdvancedExpressionEditor(props: {
  onChange: inputChangeFunc,
  expression: string,
}) {
  const classes = useStyles();
  return (
    <TextField
      required
      label="Expression"
      placeholder="Ex: up == 0"
      value={props.expression}
      onChange={props.onChange(value => ({expression: value}))}
      fullWidth
      InputProps={{
        endAdornment: (
          <InputAdornment position="end">
            <Tooltip
              title={
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
          </InputAdornment>
        ),
      }}
    />
  );
}

function SeverityEditor(props: {
  onChange: inputChangeFunc,
  severity: string,
  options: Array<MenuItemProps>,
}) {
  return (
    <TextField
      required
      label="Severity"
      select
      fullWidth
      value={props.severity}
      onChange={props.onChange(value => ({severity: value}))}>
      {props.options.map(opt => (
        <MenuItem {...opt} />
      ))}
    </TextField>
  );
}

function TimeEditor(props: {
  onChange: inputChangeFunc,
  timeNumber: string,
  timeUnit: string,
}) {
  return (
    <>
      <Grid item xs={6}>
        <TimeNumberEditor
          onChange={props.onChange}
          timeNumber={props.timeNumber}
        />
      </Grid>
      <Grid item xs={5}>
        <TimeUnitEditor
          onChange={props.onChange}
          timeUnit={props.timeUnit}
          timeUnits={timeUnits}
        />
      </Grid>
      <Grid item xs={1}>
        <Tooltip
          title={
            'Enter the amount of time the alert expression needs to be ' +
            'true for before the alert fires.'
          }
          placement="right">
          <HelpIcon />
        </Tooltip>
      </Grid>
    </>
  );
}

function TimeNumberEditor(props: {
  onChange: inputChangeFunc,
  timeNumber: string,
}) {
  return (
    <TextField
      type="number"
      value={props.timeNumber}
      onChange={props.onChange(val => ({timeNumber: val}))}
      label="Duration"
      fullWidth
    />
  );
}

function TimeUnitEditor(props: {
  onChange: inputChangeFunc,
  timeUnit: string,
  timeUnits: Array<TimeUnit>,
}) {
  return (
    <TextField
      select
      value={props.timeUnit}
      onChange={props.onChange(val => ({timeUnit: val}))}
      label="Unit"
      fullWidth>
      {props.timeUnits.map(option => (
        <MenuItem key={option.value} value={option.value}>
          {option.label}
        </MenuItem>
      ))}
    </TextField>
  );
}

export default function PrometheusEditor(props: Props) {
  const {updateAlertConfig, onExit, saveAlertRule, isNew, rule} = props;

  const classes = useStyles();
  const [form, setFormState] = React.useState<FormState>(fromAlertConfig(rule));

  /**
   * Passes the event value to an updater function which returns an update
   * object to be merged into the form. After the internal form state is
   * updated, the parent component is notified of the updated AlertConfig
   */
  const handleInputChange = React.useCallback(
    (formUpdate: (val: string) => $Shape<FormState>) => (
      event: SyntheticInputEvent<HTMLElement>,
    ) => {
      const value = event.target.value;
      const updatedForm = {
        ...form,
        ...formUpdate(value),
      };
      setFormState(updatedForm);
      updateAlertConfig(toAlertConfig(updatedForm));
    },
    [form],
  );

  const severityOptions = React.useMemo<Array<MenuItemProps>>(
    () =>
      Object.keys(SEVERITY).map(key => ({
        key: key,
        value: key,
        children: key.toUpperCase(),
      })),
    [],
  );

  return (
    <Grid container spacing={3}>
      <Grid container item direction="column" spacing={2} wrap="nowrap">
        <Grid item xs={12} sm={3}>
          <RuleNameEditor
            onChange={handleInputChange}
            ruleName={form.ruleName}
          />
        </Grid>
        <Grid item xs={12} sm={3}>
          <AdvancedExpressionEditor
            onChange={handleInputChange}
            expression={form.expression}
          />
        </Grid>
        <Grid item xs={12} sm={3}>
          <SeverityEditor
            onChange={handleInputChange}
            options={severityOptions}
            severity={form.severity}
          />
        </Grid>
        <Grid container item xs={12} sm={3} spacing={1} alignItems="flex-end">
          <TimeEditor
            onChange={handleInputChange}
            timeNumber={form.timeNumber}
            timeUnit={form.timeUnit}
          />
        </Grid>
      </Grid>

      <Grid container item>
        <Button
          variant="outlined"
          onClick={() => onExit()}
          className={classes.button}>
          Close
        </Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => saveAlertRule()}
          className={classes.button}>
          {isNew ? 'Add' : 'Edit'}
        </Button>
      </Grid>
    </Grid>
  );
}

function fromAlertConfig(rule: ?AlertConfig): FormState {
  if (!rule) {
    return {
      ruleName: '',
      expression: '',
      severity: '',
      timeNumber: '',
      timeUnit: '',
    };
  }
  const timeString = rule.for ?? '';
  const {timeNumber, timeUnit} = parseTimeString(timeString);
  return {
    ruleName: rule.alert,
    expression: rule.expr,
    severity: rule.labels?.severity || '',
    timeNumber,
    timeUnit,
  };
}

function toAlertConfig(form: FormState): AlertConfig {
  return {
    alert: form.ruleName,
    expr: form.expression,
    labels: {
      severity: form.severity,
    },
    for: `${form.timeNumber}${form.timeUnit}`,
  };
}

/***
 * When editing a rule with a duration like 1h, the api will return a duration
 * string like 1h0m0s instead of just 1h. Since the editor only allows for
 * one duration and time unit pair, take the most significant pair and return
 * only that. For example: 1h0m0s we'll just return
 * { timeNumber: 1, timeUnit: h}
 */
function parseTimeString(
  timeStamp: string,
): {timeNumber: string, timeUnit: string} {
  const units = new Set(['h', 'm', 's']);
  let duration = '';
  let unit = '';
  for (const char of timeStamp) {
    if (units.has(char)) {
      unit = char;
      break;
    }
    duration += char;
  }
  return {
    timeNumber: duration,
    timeUnit: unit,
  };
}
