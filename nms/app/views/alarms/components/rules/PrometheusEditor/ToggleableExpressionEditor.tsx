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

import * as PromQL from '../../prometheus/PromQL';
import * as React from 'react';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import RemoveCircleIcon from '@material-ui/icons/RemoveCircle';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import {LABEL_OPERATORS} from '../../prometheus/PromQLTypes';
import {Theme} from '@material-ui/core/styles';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../../AlarmContext';
import {useNetworkId} from '../../hooks';
import {useSnackbars} from '../../../../../hooks/useSnackbar';

import {InputChangeFunc} from './PrometheusEditor';
import {SelectProps} from '@material-ui/core/Select/Select';
import {getErrorMessage} from '../../../../../util/ErrorUtils';
import type {BinaryComparator} from '../../prometheus/PromQLTypes';

const useStyles = makeStyles<Theme>(theme => ({
  button: {
    marginLeft: -theme.spacing(0.5),
    margin: theme.spacing(1.5),
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
  metricName: string;
  // $FlowFixMe[value-as-type] migrated to TypeScript
  comparator: PromQL.BinaryComparator;
  // $FlowFixMe[value-as-type] migrated to TypeScript
  filters: PromQL.Labels;
  value: number;
};

type LabelValuesLookup = Map<string, Set<string>>;

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
  onChange: InputChangeFunc;
  onThresholdExpressionChange: (expresion: ThresholdExpression) => void;
  expression: ThresholdExpression;
  stringExpression: string;
  toggleOn: boolean;
  onToggleChange: () => void;
}) {
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const networkId = useNetworkId();
  const {response, error} = apiUtil.useAlarmsApi(apiUtil.getMetricNames, {
    networkId,
  });

  if (error) {
    snackbars.error(`Error retrieving metrics: ${getErrorMessage(error)}`);
  }

  return (
    <Grid container item xs={12}>
      <ThresholdExpressionEditor
        onChange={props.onThresholdExpressionChange}
        expression={props.expression}
        metricNames={response ?? []}
        onToggleChange={props.onToggleChange}
      />
    </Grid>
  );
}

export function AdvancedExpressionEditor(props: {
  onChange: InputChangeFunc;
  expression: string;
}) {
  return (
    <Grid item>
      <InputLabel htmlFor="metric-advanced-input">Expression</InputLabel>
      <TextField
        id="metric-advanced-input"
        required
        placeholder="SNR >= 0"
        value={props.expression}
        onChange={props.onChange(value => ({expression: value}))}
        fullWidth
      />
    </Grid>
  );
}

function ConditionSelector(props: {
  onChange: (expression: ThresholdExpression) => void;
  expression: ThresholdExpression;
}) {
  const conditions: Array<BinaryComparator> = [
    '>',
    '<',
    '==',
    '>=',
    '<=',
    '!=',
  ];
  return (
    <Grid>
      <InputLabel htmlFor="condition-input">Condition</InputLabel>
      <TextField
        id="condition-input"
        fullWidth
        required
        select
        value={props.expression.comparator.op}
        onChange={({target}) => {
          props.onChange({
            ...props.expression,
            comparator: new PromQL.BinaryComparator(
              // Cast to element type of conditions as it's item type
              target.value as typeof conditions[number],
            ),
          });
        }}>
        {conditions.map(item => (
          <MenuItem key={item} value={item}>
            {item}
          </MenuItem>
        ))}
      </TextField>
    </Grid>
  );
}

function ValueSelector(props: {
  onChange: (expression: ThresholdExpression) => void;
  expression: ThresholdExpression;
}) {
  return (
    <Grid item>
      <InputLabel htmlFor="value-input">Value</InputLabel>
      <TextField
        id="value-input"
        fullWidth
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
  );
}

function MetricSelector(props: {
  expression: ThresholdExpression;
  onChange: (expression: ThresholdExpression) => void;
  metricNames: Array<string>;
}) {
  const {metricNames} = props;
  return (
    <Grid item>
      <InputLabel htmlFor="metric-input">Metric</InputLabel>
      <Autocomplete
        id="metric-input"
        options={metricNames}
        groupBy={getMetricNamespace}
        value={props.expression.metricName}
        onChange={(_e, value) => {
          props.onChange({...props.expression, metricName: value!});
        }}
        renderInput={params => <TextField {...params} required />}
      />
    </Grid>
  );
}

function ThresholdExpressionEditor({
  expression,
  onChange,
  onToggleChange,
  metricNames,
}: {
  onChange: (expression: ThresholdExpression) => void;
  expression: ThresholdExpression;
  metricNames: Array<string>;
  onToggleChange: () => void;
}) {
  const networkId = useNetworkId();
  const {apiUtil} = useAlarmContext();
  const {metricName} = expression;
  // mapping from label name to all values in response
  const [labels, setLabels] = React.useState<LabelValuesLookup>(new Map());
  // cache all label names
  const labelNames = React.useMemo<Array<string>>(
    () => getFilteredListOfLabelNames(Array.from(labels.keys())),
    [labels],
  );
  React.useEffect(() => {
    async function getMetricLabels() {
      const response = (
        await apiUtil.getMetricSeries({
          name: metricName,
          networkId: networkId,
        })
      ).data;
      const labelValues = new Map<string, Set<string>>();
      for (const metric of response) {
        for (const labelName of Object.keys(metric)) {
          let set = labelValues.get(labelName);
          if (!set) {
            set = new Set<string>();
            labelValues.set(labelName, set);
          }
          const labelValue = metric[labelName];
          set.add(labelValue);
        }
      }
      setLabels(labelValues);
    }
    if (metricName != null && metricName !== '') {
      void getMetricLabels();
    }
  }, [metricName, networkId, setLabels, apiUtil]);

  return (
    <Grid item container spacing={1}>
      <Grid
        item
        container
        spacing={1}
        alignItems="flex-end"
        justifyContent="space-between">
        <Grid item xs={7}>
          <MetricSelector
            expression={expression}
            onChange={onChange}
            metricNames={metricNames}
          />
        </Grid>
        <Grid item xs={3}>
          <ConditionSelector expression={expression} onChange={onChange} />
        </Grid>
        <Grid item xs={2}>
          <ValueSelector expression={expression} onChange={onChange} />
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <MetricFilters
          labelNames={labelNames}
          labelValues={labels}
          expression={expression}
          onChange={onChange}
          onToggleChange={onToggleChange}
        />
      </Grid>
    </Grid>
  );
}

function MetricFilters(props: {
  labelNames: Array<string>;
  labelValues: LabelValuesLookup;
  expression: ThresholdExpression;
  onChange: (expression: ThresholdExpression) => void;
  onToggleChange: () => void;
}) {
  const classes = useStyles();
  const isMetricSelected =
    props.expression?.metricName != null && props.expression?.metricName !== '';
  return (
    <Grid item container direction="column">
      <Grid item>
        <Button
          className={classes.button}
          color="primary"
          size="small"
          disabled={!isMetricSelected}
          onClick={() => {
            const filtersCopy = props.expression.filters.copy();
            filtersCopy.addEqual('', '');
            props.onChange({
              ...props.expression,
              filters: filtersCopy,
            });
          }}>
          Add new filter
        </Button>
        <Button
          className={classes.button}
          color="primary"
          size="small"
          onClick={props.onToggleChange}>
          Write a custom expression
        </Button>
      </Grid>
      <Grid item container direction="column" spacing={3}>
        {props.expression.filters.labels.map((filter, idx) => (
          <Grid item key={idx}>
            <LabelFilter
              labelNames={props.labelNames}
              labelValues={props.labelValues}
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
      </Grid>
    </Grid>
  );
}

function LabelFilter(props: {
  labelNames: Array<string>;
  labelValues: LabelValuesLookup;
  onChange: (expression: ThresholdExpression) => void;
  onRemove: (filerIdx: number) => void;
  expression: ThresholdExpression;
  filterIdx: number;
  selectedLabel: string;
  selectedValue: string;
}) {
  const currentFilter = props.expression.filters.labels[props.filterIdx];
  const values = Array.from(props.labelValues.get(props.selectedLabel) ?? []);
  return (
    <Grid item container xs={12} spacing={1} alignItems="flex-start">
      <Grid item xs={6}>
        <InputLabel htmlFor={`metric-input-${props.filterIdx}`}>
          Label
        </InputLabel>
        <FilterSelector
          id={`metric-input-${props.filterIdx}`}
          values={props.labelNames}
          defaultVal=""
          onChange={({target}) => {
            const filtersCopy = props.expression.filters.copy();
            filtersCopy.setIndex(props.filterIdx, target.value, '');
            props.onChange({...props.expression, filters: filtersCopy});
          }}
          selectedValue={props.selectedLabel}
        />
      </Grid>
      <Grid item xs={2}>
        <Grid>
          <InputLabel htmlFor={`condition-input-${props.filterIdx}`}>
            Condition
          </InputLabel>
          <TextField
            id={`condition-input-${props.filterIdx}`}
            fullWidth
            required
            select
            value={currentFilter.operator}
            onChange={({target}) => {
              const filtersCopy = props.expression.filters.copy();
              const filterOperator = isRegexValue(target.value) ? '=~' : '=';
              filtersCopy.setIndex(
                props.filterIdx,
                currentFilter.name,
                currentFilter.value,
                filterOperator,
              );
              props.onChange({...props.expression, filters: filtersCopy});
            }}>
            {LABEL_OPERATORS.map(item => (
              <MenuItem key={item} value={item}>
                {item}
              </MenuItem>
            ))}
          </TextField>
        </Grid>
      </Grid>
      <Grid item xs={3}>
        <Grid item>
          <InputLabel htmlFor={`value-input-${props.filterIdx}`}>
            Value
          </InputLabel>
          <Autocomplete
            value={currentFilter.value}
            freeSolo
            options={values}
            onChange={(_e, value) => {
              const filtersCopy = props.expression.filters.copy();
              filtersCopy.setIndex(
                props.filterIdx,
                currentFilter.name,
                value!,
                currentFilter.operator,
              );
              props.onChange({
                ...props.expression,
                filters: filtersCopy,
              });
            }}
            renderInput={params => (
              <TextField
                {...params}
                required
                id={`value-input-${props.filterIdx}`}
              />
            )}
          />
        </Grid>
      </Grid>
      <Grid item xs={1} container alignItems="center" justifyContent="flex-end">
        <IconButton onClick={() => props.onRemove(props.filterIdx)} edge="end">
          <RemoveCircleIcon />
        </IconButton>
      </Grid>
    </Grid>
  );
}

function FilterSelector(props: {
  id: string;
  values: Array<string>;
  defaultVal: string;
  onChange: (event: React.ChangeEvent<{value: string}>) => void;
  selectedValue: string;
  disabled?: boolean;
}) {
  const classes = useStyles();
  const menuItems = props.values.map(val => (
    <MenuItem value={val} key={val}>
      {val}
    </MenuItem>
  ));

  return (
    <Select
      id={props.id}
      fullWidth
      disabled={props.disabled}
      displayEmpty
      className={classes.metricFilterItem}
      value={props.selectedValue}
      onChange={props.onChange as SelectProps['onChange']}>
      <MenuItem disabled value="">
        {props.defaultVal}
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

/**
 * Gets the first word application prefix of a prometheus metric name. This
 * is known by most client libraries as a namespace.
 */
function getMetricNamespace(option: string) {
  const index = option.indexOf('_');
  if (index > -1) {
    return option.slice(0, index);
  }
  return option;
}
