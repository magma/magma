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

import * as React from 'react';
import Grid from '@material-ui/core/Grid';
import InputLabel from '@material-ui/core/InputLabel';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import {Theme} from '@material-ui/core/styles';
import {makeStyles} from '@material-ui/styles';
import type {RuleInterfaceMap} from './RuleInterface';

const useRuleTypeStyles = makeStyles<Theme>(theme => ({
  buttonGroup: {
    paddingTop: theme.spacing(1),
  },
  button: {
    textTransform: 'capitalize',
  },
  label: {
    fontSize: theme.typography.pxToRem(14),
  },
}));

export default function SelectRuleType<TRuleUnion>({
  ruleMap,
  value,
  onChange,
}: {
  ruleMap: RuleInterfaceMap<TRuleUnion>;
  onChange: (ruleType: string) => void;
  value: string;
}) {
  const classes = useRuleTypeStyles();
  const ruleTypes = React.useMemo<Array<{type: string; friendlyName: string}>>(
    () =>
      Object.keys(ruleMap || {}).map(key => ({
        type: key,
        friendlyName: ruleMap[key].friendlyName || key,
      })),
    [ruleMap],
  );

  const handleChange = React.useCallback(
    (_e: React.MouseEvent<HTMLElement>, val: string) => {
      onChange(val);
    },
    [onChange],
  );

  // if there's < 2 rule types, just stick with the default rule type
  if (ruleTypes.length < 2) {
    return null;
  }

  /**
   * Grid structure is chosen here to match the selected editor's width
   * and padding.
   */
  return (
    <Grid item>
      <InputLabel className={classes.label}>Rule Type</InputLabel>
      <ToggleButtonGroup
        className={classes.buttonGroup}
        size="medium"
        color="primary"
        value={value}
        onChange={handleChange}
        exclusive>
        {ruleTypes.map(ruleType => (
          <ToggleButton
            className={classes.button}
            key={ruleType.type}
            value={ruleType.type}>
            {ruleType.friendlyName}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </Grid>
  );
}
