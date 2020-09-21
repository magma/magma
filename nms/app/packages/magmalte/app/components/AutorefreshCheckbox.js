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

import Checkbox from '@material-ui/core/Checkbox';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import React from 'react';
import Text from '../theme/design-system/Text';
import moment from 'moment';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

export type UseRefreshingDateRangeHook = (
  shouldAutorefreshByDefault: boolean,
  updateInterval: number,
  onDateRangeChange: () => void,
) => {|
  startDate: moment,
  endDate: moment,
  setStartDate: (date: moment) => void,
  setEndDate: (date: moment) => void,
  isAutorefreshing: boolean,
  setIsAutorefreshing: ((boolean => boolean) | boolean) => void,
|};

export const useRefreshingDateRange: UseRefreshingDateRangeHook = (
  shouldAutorefreshByDefault,
  updateInterval,
  onDateRangeChange,
) => {
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());
  const [isAutorefreshing, setIsAutorefreshing] = useState(
    shouldAutorefreshByDefault,
  );

  useEffect(() => {
    if (isAutorefreshing) {
      const interval = setInterval(() => {
        setEndDate(moment());
        onDateRangeChange();
      }, updateInterval);

      return () => clearInterval(interval);
    }
  }, [endDate, startDate, onDateRangeChange, isAutorefreshing, updateInterval]);

  const modifiedSetStartDate = useCallback(
    (date: moment) => {
      setIsAutorefreshing(false);
      setStartDate(date);
      onDateRangeChange();
    },
    [onDateRangeChange],
  );

  const modifiedSetEndDate = useCallback(
    (date: moment) => {
      setIsAutorefreshing(false);
      setEndDate(date);
      onDateRangeChange();
    },
    [onDateRangeChange],
  );

  return {
    startDate,
    endDate,
    setStartDate: modifiedSetStartDate,
    setEndDate: modifiedSetEndDate,
    isAutorefreshing,
    setIsAutorefreshing,
  };
};

const useStyles = makeStyles(() => ({
  autorefreshCheckbox: {
    color: colors.primary.comet,
  },
}));

type Props = {
  autorefreshEnabled: boolean,
  onToggle: (boolean | (boolean => boolean)) => void,
};

export default function AutorefreshCheckbox(props: Props) {
  const {autorefreshEnabled, onToggle} = props;
  const classes = useStyles();

  return (
    <FormControlLabel
      control={<Checkbox checked={autorefreshEnabled} onChange={onToggle} />}
      label={
        <Text variant="body3" className={classes.autorefreshCheckbox}>
          Autorefresh
        </Text>
      }
    />
  );
}
