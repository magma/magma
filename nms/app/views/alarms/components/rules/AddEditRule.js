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
 * @flow
 * @format
 */
import Grid from '@material-ui/core/Grid';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import RuleContext from './RuleContext';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useState} from 'react';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {GenericRule} from './RuleInterface';

type Props<TRuleUnion> = {
  onExit: () => void,
  //TODO rename?
  initialConfig: ?GenericRule<TRuleUnion>,
  isNew: boolean,
  defaultRuleType?: string,
};

const useStyles = makeStyles(_theme => ({
  gridContainer: {
    flexGrow: 1,
  },
}));

export default function AddEditRule<TRuleUnion>(props: Props<TRuleUnion>) {
  const {ruleMap} = useAlarmContext();
  const {isNew, onExit} = props;
  const classes = useStyles();
  const [rule, setRule] = useState<?GenericRule<TRuleUnion>>(
    props.initialConfig,
  );

  const [selectedRuleType, setSelectedRuleType] = React.useState<string>(
    rule?.ruleType || props.defaultRuleType || 'prometheus',
  );

  // null out in-progress rule so next editor doesnt see an incompatible schema
  const selectRuleType = React.useCallback(
    type => {
      setRule(null);
      return setSelectedRuleType(type);
    },
    [setRule, setSelectedRuleType],
  );

  const {RuleEditor} = ruleMap[selectedRuleType];
  return (
    <RuleContext.Provider
      value={{
        ruleMap: ruleMap,
        ruleType: selectedRuleType,
        selectRuleType: selectRuleType,
      }}>
      <Grid
        className={classes.gridContainer}
        container
        spacing={0}
        data-testid="add-edit-alert">
        <RuleEditor
          isNew={isNew}
          onExit={onExit}
          onRuleUpdated={setRule}
          rule={rule}
        />
      </Grid>
    </RuleContext.Provider>
  );
}
