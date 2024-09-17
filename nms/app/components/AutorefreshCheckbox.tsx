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

import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import React from 'react';
import Text from '../theme/design-system/Text';
import makeStyles from '@mui/styles/makeStyles';
import {colors} from '../theme/default';
import {subHours} from 'date-fns';
import {useCallback, useEffect, useState} from 'react';

export type UseRefreshingDateRangeHook = (
  isAutoRefreshing: boolean,
  updateInterval: number,
  onDateRangeChange: () => void,
) => {
  startDate: Date;
  endDate: Date;
  setStartDate: (date: Date) => void;
  setEndDate: (date: Date) => void;
};

export const useRefreshingDateRange: UseRefreshingDateRangeHook = (
  isAutoRefreshing,
  updateInterval,
  onDateRangeChange,
) => {
  const [startDate, setStartDate] = useState(subHours(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());

  useEffect(() => {
    if (isAutoRefreshing) {
      const interval = setInterval(() => {
        setEndDate(new Date());
        onDateRangeChange();
      }, updateInterval);

      return () => clearInterval(interval);
    }
  }, [endDate, startDate, onDateRangeChange, isAutoRefreshing, updateInterval]);

  const modifiedSetStartDate = useCallback(
    (date: Date) => {
      setStartDate(date);
      onDateRangeChange();
    },
    [onDateRangeChange],
  );

  const modifiedSetEndDate = useCallback(
    (date: Date) => {
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
  };
};

const useStyles = makeStyles(() => ({
  autorefreshCheckbox: {
    color: colors.primary.comet,
  },
}));

type Props = {
  autorefreshEnabled: boolean;
  onToggle: () => void;
};

export default function AutorefreshCheckbox(props: Props) {
  const {autorefreshEnabled, onToggle} = props;
  const classes = useStyles();

  return (
    <FormControlLabel
      control={
        <Checkbox
          checked={autorefreshEnabled}
          onChange={onToggle}
          data-testid="autorefresh-checkbox"
        />
      }
      label={
        <Text variant="body3" className={classes.autorefreshCheckbox}>
          Autorefresh
        </Text>
      }
    />
  );
}
