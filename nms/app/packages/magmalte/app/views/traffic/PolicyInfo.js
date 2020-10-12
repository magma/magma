/*
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

import DialogContent from '@material-ui/core/DialogContent';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';

import {AltFormField} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import type {policy_rule} from '@fbcnms/magma-api';

const useStyles = makeStyles(() => ({
  title: {textAlign: 'center', margin: 'auto', marginLeft: '0px'},
  switch: {margin: 'auto 0px'},
}));

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  isNetworkWide: boolean,
  setIsNetworkWide: boolean => void,
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyInfoEdit(props: Props) {
  const classes = useStyles();
  return (
    <>
      <DialogContent
        data-testid="networkInfoEdit"
        className={props.dialogClass}>
        <List>
          <Typography
            variant="caption"
            display="block"
            className={props.descriptionClass}
            gutterBottom>
            {'Basic policy rule fields'}
          </Typography>
          <ListItem dense disableGutters />
          <AltFormField
            label={'Policy ID'}
            subLabel={'A unique identifier for the policy rule'}
            disableGutters>
            <OutlinedInput
              className={props.inputClass}
              data-testid="policyID"
              placeholder="Eg. policy_id"
              value={props.policyRule.id}
              onChange={({target}) => {
                props.onChange({...props.policyRule, id: target.value});
              }}
            />
          </AltFormField>
          <AltFormField
            label={'Priority Level'}
            subLabel={'Higher priority policies override lower priority ones'}
            disableGutters>
            <OutlinedInput
              className={props.inputClass}
              data-testid="policyPriority"
              placeholder="Value between 1 and 15"
              fullWidth={true}
              value={props.policyRule.priority}
              onChange={({target}) =>
                props.onChange({...props.policyRule, priority: target.value})
              }
            />
          </AltFormField>
          <Grid container justify="space-between" className={props.inputClass}>
            <Grid item className={classes.title}>
              <AltFormField disableGutters label={'Network Wide'} />
            </Grid>
            <Grid item className={classes.switch}>
              <FormControlLabel
                control={
                  <Switch
                    color="primary"
                    checked={props.isNetworkWide}
                    onChange={({target}) =>
                      props.setIsNetworkWide(target.checked)
                    }
                  />
                }
                label={props.isNetworkWide ? 'Enabled' : 'Disabled'}
                labelPlacement="start"
              />
            </Grid>
          </Grid>
        </List>
      </DialogContent>
    </>
  );
}
