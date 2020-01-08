/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as PromQL from './prometheus/PromQL';
import * as React from 'react';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
import ToggleableExpressionEditor, {
  AdvancedExpressionEditor,
  thresholdToPromQL,
} from './ToggleableExpressionEditor';
import Tooltip from '@material-ui/core/Tooltip';
import {BINARY_COMPARATORS} from './prometheus/PromQLTypes';
import {Labels} from './prometheus/PromQL';
import {Parse} from './prometheus/PromQLParser';
import {SEVERITY} from './Severity';

import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

import type {AlertConfig} from './AlarmAPIType';
import type {BinaryComparator} from './prometheus/PromQLTypes';
import type {GenericRule, RuleEditorProps} from './RuleInterface';
import type {ThresholdExpression} from './ToggleableExpressionEditor';

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
  timeNumber: number,
  timeUnit: string,
  description: string,
};

export type InputChangeFunc = (
  formUpdate: (val: string) => $Shape<FormState>,
) => (event: SyntheticInputEvent<HTMLElement>) => void;

function RuleNameEditor({
  ruleName,
  ...props
}: {
  onChange: InputChangeFunc,
  ruleName: string,
}) {
  return (
    <TextField
      required
      label="Rule Name"
      placeholder="Ex: Service Down"
      value={ruleName}
      onChange={props.onChange(value => ({ruleName: value}))}
      fullWidth
    />
  );
}

function SeverityEditor(props: {
  onChange: InputChangeFunc,
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
  onChange: InputChangeFunc,
  timeNumber: number,
  timeUnit: string,
}) {
  return (
    <Grid container spacing={1} alignItems="flex-end">
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
    </Grid>
  );
}

function TimeNumberEditor(props: {
  onChange: InputChangeFunc,
  timeNumber: number,
}) {
  return (
    <TextField
      type="number"
      value={isNaN(props.timeNumber) ? '' : props.timeNumber}
      onChange={props.onChange(val => ({timeNumber: parseInt(val, 10)}))}
      label="Duration"
      fullWidth
    />
  );
}

function TimeUnitEditor(props: {
  onChange: InputChangeFunc,
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

function DescriptionEditor(props: {
  onChange: InputChangeFunc,
  description: string,
}) {
  return (
    <TextField
      value={props.description}
      onChange={props.onChange(val => ({description: val}))}
      label="Description"
      fullWidth
    />
  );
}

type PrometheusEditorProps = {
  ...RuleEditorProps<AlertConfig>,
  thresholdEditorEnabled?: ?boolean,
};
export default function PrometheusEditor(props: PrometheusEditorProps) {
  const {apiUtil, isNew, onRuleUpdated, onExit, rule} = props;
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const classes = useStyles();
  const [formState, setFormState] = React.useState<FormState>(
    fromAlertConfig(rule ? rule.rawRule : null),
  );
  const [
    thresholdExpression,
    setThresholdExpression,
  ] = React.useState<ThresholdExpression>({
    metricName: '',
    comparator: '==',
    value: 0,
    filters: new Labels(),
  });

  const [advancedEditorMode, setAdvancedEditorMode] = React.useState<boolean>(
    !props.thresholdEditorEnabled,
  );

  const parsedExpression = React.useMemo(() => {
    try {
      return Parse(props.rule?.expression);
    } catch {
      return null;
    }
    // We only want to parse the expression on the first
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  React.useEffect(() => {
    if (!props.rule?.expression) {
      setAdvancedEditorMode(false);
    } else if (parsedExpression) {
      const newThresholdExpression = getThresholdExpression(parsedExpression);
      if (newThresholdExpression) {
        setAdvancedEditorMode(false);
        setThresholdExpression(newThresholdExpression);
      } else {
        setAdvancedEditorMode(true);
      }
    } else {
      enqueueSnackbar(
        "Error parsing alert expression. You can still edit this using the advanced editor, but you won't be able to use the UI expression editor.",
        {
          variant: 'error',
        },
      );
    }
  }, [enqueueSnackbar, parsedExpression]);

  const saveAlert = async () => {
    try {
      if (!rule) {
        throw new Error('Alert config empty');
      }

      const request = {
        networkId: match.params.networkId,
        rule: toAlertConfig(formState),
      };
      if (isNew) {
        await apiUtil.createAlertRule(request);
      } else {
        await apiUtil.editAlertRule(request);
      }
      enqueueSnackbar(`Successfully ${isNew ? 'added' : 'saved'} alert rule`, {
        variant: 'success',
      });
      onExit();
    } catch (error) {
      enqueueSnackbar(
        `Unable to create alert: ${
          error.response ? error.response.data.message : error.message
        }.`,
        {
          variant: 'error',
        },
      );
    }
  };

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
        ...formState,
        ...formUpdate(value),
      };
      setFormState(updatedForm);
      const updatedConfig = toAlertConfig(updatedForm);
      onRuleUpdated({
        ...rule,
        ...({rawRule: updatedConfig}: $Shape<GenericRule<AlertConfig>>),
      });
    },
    [formState, onRuleUpdated, rule],
  );

  // TODO: pull out common functionality between this and handleInputChange
  const updateThresholdExpression = React.useCallback(
    (expression: ThresholdExpression) => {
      setThresholdExpression(expression);
      const stringExpression = thresholdToPromQL(expression);
      const updatedForm = {
        ...formState,
        expression: stringExpression,
      };
      setFormState(updatedForm);
      const updatedConfig = toAlertConfig(updatedForm);
      onRuleUpdated({
        ...rule,
        ...({rawRule: updatedConfig}: $Shape<GenericRule<AlertConfig>>),
      });
    },
    [formState, onRuleUpdated, rule],
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
            ruleName={formState.ruleName}
            disabled={!isNew}
          />
        </Grid>
        {props.thresholdEditorEnabled ? (
          <ToggleableExpressionEditor
            apiUtil={props.apiUtil}
            onChange={handleInputChange}
            onThresholdExpressionChange={updateThresholdExpression}
            expression={thresholdExpression}
            stringExpression={formState.expression}
            toggleOn={advancedEditorMode}
            onToggleChange={val => setAdvancedEditorMode(val)}
          />
        ) : (
          <Grid item xs={4}>
            <AdvancedExpressionEditor
              expression={formState.expression}
              onChange={handleInputChange}
            />
          </Grid>
        )}

        <Grid item xs={12} sm={3}>
          <SeverityEditor
            onChange={handleInputChange}
            options={severityOptions}
            severity={formState.severity}
          />
        </Grid>
        <Grid item xs={12} sm={3}>
          <TimeEditor
            onChange={handleInputChange}
            timeNumber={formState.timeNumber}
            timeUnit={formState.timeUnit}
          />
        </Grid>
        <Grid item xs={12} sm={3}>
          <DescriptionEditor
            onChange={handleInputChange}
            description={formState.description}
          />
        </Grid>
      </Grid>

      <Grid item>
        <Button
          variant="outlined"
          onClick={() => onExit()}
          className={classes.button}>
          Close
        </Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => saveAlert()}
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
      description: '',
      timeNumber: 0,
      timeUnit: '',
    };
  }
  const timeString = rule.for ?? '';
  const {timeNumber, timeUnit} = getMostSignificantTime(
    parseTimeString(timeString),
  );
  return {
    ruleName: rule.alert || '',
    expression: rule.expr || '',
    severity: rule.labels?.severity || '',
    description: rule.annotations?.description || '',
    timeNumber: timeNumber,
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
    annotations: {
      description: form.description,
    },
  };
}

export type Duration = {
  hours: number,
  minutes: number,
  seconds: number,
};

export function parseTimeString(timeStamp: string): Duration {
  if (timeStamp === '') {
    return {hours: 0, minutes: 0, seconds: 0};
  }
  const durationRegex = /^((\d+)h)*((\d+)m)*((\d+)s)*$/;
  const duration = timeStamp.match(durationRegex);
  if (!duration) {
    return {hours: 0, minutes: 0, seconds: 0};
  }
  // Index is corresponding capture group from regex
  const hours = parseInt(duration[2], 10) || 0;
  const minutes = parseInt(duration[4], 10) || 0;
  const seconds = parseInt(duration[6], 10) || 0;
  return {hours, minutes, seconds};
}

function getMostSignificantTime(
  duration: Duration,
): {timeNumber: number, timeUnit: string} {
  if (duration.hours) {
    return {timeNumber: duration.hours, timeUnit: 'h'};
  } else if (duration.minutes) {
    return {timeNumber: duration.minutes, timeUnit: 'm'};
  }
  return {timeNumber: duration.seconds, timeUnit: 's'};
}

function getThresholdExpression(exp: PromQL.Expression): ?ThresholdExpression {
  if (
    !(
      exp instanceof PromQL.BinaryOperation &&
      exp.lh instanceof PromQL.InstantSelector &&
      exp.rh instanceof PromQL.Scalar
    )
  ) {
    return null;
  }

  const metricName = exp.lh.selectorName || '';
  const threshold = exp.rh.value;
  const filters = exp.lh.labels || new PromQL.Labels();
  filters.removeByName('networkID');
  const comparator = getBinaryComparator(exp.operator);
  if (!comparator) {
    return null;
  }
  return {
    metricName,
    filters,
    comparator,
    value: threshold,
  };
}

function getBinaryComparator(str: string): ?BinaryComparator {
  if (BINARY_COMPARATORS.includes(str)) {
    return BINARY_COMPARATORS[BINARY_COMPARATORS.indexOf(str)];
  }
  return null;
}
