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
 */

import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import AddIcon from '@mui/icons-material/Add';
import Button from '@mui/material/Button';
import DeleteIcon from '@mui/icons-material/Delete';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FormControl from '@mui/material/FormControl';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemSecondaryAction from '@mui/material/ListItemSecondaryAction';
import ListItemText from '@mui/material/ListItemText';
import LteNetworkContext from '../../context/LteNetworkContext';
import MenuItem from '@mui/material/MenuItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import React from 'react';
import Select from '@mui/material/Select';
import Text from '../../theme/design-system/Text';
import {AltFormField} from '../../components/FormField';
import {EditSubscriberProps} from './SubscriberUtils';
import {makeStyles} from '@mui/styles';
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
                    }}
                    size="large">
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
