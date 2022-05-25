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
 * @flow strict-local
 * @format
 */

import type {TimeRange} from './AsyncMetric';

import * as React from 'react';
import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from './AsyncMetric';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TextField from '@material-ui/core/TextField';
import TimeRangeSelector from './TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';
import {useParams} from 'react-router-dom';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  appBar: {
    display: 'inline-block',
  },
  chartRow: {
    display: 'flex',
  },
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
  selectorAutocomplete: {
    width: '400px',
  },
}));

export type MetricGraphConfig = {
  basicQueryConfigs: BasicQueryConfig[],
  customQueryConfigs?: CustomQuery[],
  label: string,
  unit?: string,
  legendLabels?: string[],
};

export type CustomQuery = {
  resolveQuery: string => string,
};

export type BasicQueryConfig = {
  filters: MetricLabel[],
  metric: string,
};

export type MetricLabel = {
  name: string,
  value: string,
};

export function resolveQuery(
  config: MetricGraphConfig,
  filterName: string,
  filterValue: string,
): string[] {
  if (config.customQueryConfigs) {
    return resolveCustomQuery(config.customQueryConfigs, filterValue);
  }
  return resolveBasicQuery(config.basicQueryConfigs, filterName, filterValue);
}

function resolveBasicQuery(
  configs: BasicQueryConfig[],
  filterName: string,
  filterValue: string,
): string[] {
  return configs.map(config => {
    const filterString = resolveFilters(
      config.filters,
      filterName,
      filterValue,
    );
    return `${config.metric}{${filterString}}`;
  });
}

function resolveFilters(
  filters: MetricLabel[],
  filterName: string,
  filterValue: string,
): string {
  const dbFilters: string[] = filters.map(
    filter => filter.name + '="' + filter.value + '"',
  );
  dbFilters.push(`${filterName}="${filterValue}"`);
  return dbFilters.join(',');
}

function resolveCustomQuery(
  configs: CustomQuery[],
  filterValue: string,
): string[] {
  return configs.map(config => config.resolveQuery(filterValue));
}

export default function (props: {
  selectors: Array<string>,
  defaultSelector: string,
  onSelectorChange: (
    event: SyntheticInputEvent<EventTarget>,
    value: string,
  ) => void,
  configs: MetricGraphConfig[],
  selectorName: string,
  renderOptionOverride?: string => React.Node,
}) {
  const {selectedID} = useParams();
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');

  const selectedOrDefault = selectedID || props.defaultSelector;

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="devices">{props.selectorName}</InputLabel>
          <Autocomplete
            className={classes.selectorAutocomplete}
            defaultValue={props.defaultSelector}
            options={props.selectors}
            onChange={props.onSelectorChange}
            disableClearable
            renderOption={option =>
              props.renderOptionOverride
                ? props.renderOptionOverride(option)
                : option
            }
            renderInput={params => (
              <TextField
                ref={params.ref}
                InputLabelProps={params.InputLabelProps}
                InputProps={params.InputProps}
                inputProps={params.inputProps}
                variant="filled"
                fullWidth
              />
            )}
          />
        </FormControl>
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
      </AppBar>
      <GridList cols={2} cellHeight={300}>
        {props.configs.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Text variant="h6">{config.label}</Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(
                      config,
                      props.selectorName,
                      selectedOrDefault,
                    )}
                    timeRange={timeRange}
                  />
                </div>
              </CardContent>
            </Card>
          </GridListTile>
        ))}
      </GridList>
    </>
  );
}
