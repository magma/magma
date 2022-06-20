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

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import FormControl from '@material-ui/core/FormControl';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';

// $FlowFixMe migrated to typescript
import {AltFormField} from '../../components/FormField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {EditSubscriberProps} from './SubscriberUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
  accordionList: {
    width: '100%',
  },
  placeholder: {
    opacity: 0.5,
  },
  apnButton: {
    margin: '20px 0',
  },
}));

export default function EditSubscriberApnStaticIps(props: EditSubscriberProps) {
  const lteCtx = useContext(LteNetworkContext);
  const staticIpAssignments =
    lteCtx.state.cellular?.epc?.mobility?.enable_static_ip_assignments;
  const classes = useStyles();

  return (
    <div>
      <Button
        onClick={props.onAddApnStaticIP}
        disabled={!staticIpAssignments ?? false}
        className={classes.apnButton}>
        Add New APN Static IP
        <AddIcon />
      </Button>
      {props.subscriberStaticIPRows.map((apn, index) => (
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <List className={classes.accordionList}>
              <ListItem>
                <ListItemText
                  primary={
                    apn.apnName || (
                      <Text className={classes.placeholder}>{'APN'}</Text>
                    )
                  }
                />
                <ListItemSecondaryAction>
                  <IconButton
                    edge="end"
                    aria-label="delete"
                    onClick={event => {
                      event.stopPropagation();
                      props.onDeleteApn(apn);
                    }}>
                    <DeleteIcon />
                  </IconButton>
                </ListItemSecondaryAction>
              </ListItem>
            </List>
          </AccordionSummary>
          <AccordionDetails>
            <AltFormField label={'APN name'}>
              <FormControl className={classes.input}>
                <Select
                  value={apn.apnName}
                  onChange={({target}) => {
                    const staticIpApn = props.subscriberStaticIPRows.map(
                      apn => apn.apnName,
                    );
                    if (!staticIpApn.includes(target.value)) {
                      props.onTrafficPolicyChange(
                        'apnName',
                        target.value,
                        index,
                      );
                    }
                  }}
                  input={<OutlinedInput />}>
                  {(props.subscriberState.active_apns || []).map(apn => (
                    <MenuItem value={apn}>
                      <ListItemText primary={apn} />
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </AltFormField>
            <AltFormField label={'APN Static IP'}>
              <OutlinedInput
                className={classes.input}
                placeholder="Eg. 192.168.100.1"
                fullWidth={true}
                value={apn.staticIp}
                onChange={({target}) => {
                  props.onTrafficPolicyChange('staticIp', target.value, index);
                }}
              />
            </AltFormField>
          </AccordionDetails>
        </Accordion>
      ))}
    </div>
  );
}
