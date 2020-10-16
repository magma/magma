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

import type {flow_description} from '@fbcnms/magma-api';

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import DialogContent from '@material-ui/core/DialogContent';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import Typography from '@material-ui/core/Typography';

import {
  ACTION,
  DIRECTION,
  PROTOCOL,
} from '../../components/network/PolicyTypes';
import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import type {policy_rule} from '@fbcnms/magma-api';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  root: {
    '&$expanded': {
      minHeight: 'auto',
    },
    marginTop: '0px',
    marginBottom: '0px',
  },
  expanded: {marginTop: '-8px', marginBottom: '-8px'},
  block: {
    display: 'block',
  },
  flex: {display: 'flex'},
  panel: {flexGrow: 1},
  removeIcon: {alignSelf: 'baseline'},
  dialog: {height: '640px'},
  title: {textAlign: 'center', margin: 'auto', marginLeft: '0px'},
}));

type FieldProps = {
  index: number,
  flow: flow_description,
  handleDelete: number => void,
  onChange: (number, flow_description) => void,
};

function PolicyFlowFields2(props: FieldProps) {
  const classes = useStyles();
  const {flow} = props;

  const handleActionChange = action =>
    props.onChange(props.index, {
      ...props.flow,
      // $FlowIgnore: value guaranteed to match the string literals
      action,
    });

  const handleFieldChange = (field: string, value: number | string) =>
    props.onChange(props.index, {
      ...props.flow,
      match: {
        ...props.flow.match,
        [field]: value,
      },
    });

  return (
    <div className={classes.flex}>
      <Accordion className={classes.panel}>
        <AccordionSummary
          classes={{
            root: classes.root,
            expanded: classes.expanded,
          }}
          expandIcon={<ExpandMoreIcon />}>
          <Grid container justify="space-between">
            <Grid item className={classes.title}>
              <Text variant="body2">Flow {props.index + 1}</Text>
            </Grid>
            <Grid item>
              <IconButton
                className={classes.removeIcon}
                onClick={() => props.handleDelete(props.index)}>
                <DeleteOutline />
              </IconButton>
            </Grid>
          </Grid>
        </AccordionSummary>
        <AccordionDetails classes={{root: classes.block}}>
          <div className={classes.flex}>
            <Grid container spacing={2}>
              <Grid item xs={4}>
                <AltFormField disableGutters label={'Action'}>
                  <Select
                    fullWidth={true}
                    variant={'outlined'}
                    value={flow.action}
                    onChange={({target}) => {
                      handleActionChange(target.value);
                    }}
                    input={<OutlinedInput id="action" />}>
                    <MenuItem value={ACTION.PERMIT}>
                      <ListItemText primary={'Permit'} />
                    </MenuItem>
                    <MenuItem value={ACTION.DENY}>
                      <ListItemText primary={'Deny'} />
                    </MenuItem>
                  </Select>
                </AltFormField>
                {flow.match.ip_proto !== PROTOCOL.IPPROTO_ICMP && (
                  <AltFormField disableGutters label={'IPv4'}>
                    <Grid container spacing={0}>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Source'}>
                          <OutlinedInput
                            data-testid="ipv4Source"
                            placeholder="192.168.0.1/24"
                            fullWidth={true}
                            value={flow.match.ipv4_src}
                            onChange={({target}) =>
                              handleFieldChange('ipv4_src', target.value)
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Destination'}>
                          <OutlinedInput
                            data-testid="ipv4Destination"
                            placeholder="192.168.0.1/24"
                            fullWidth={true}
                            value={flow.match.ipv4_dst}
                            onChange={({target}) =>
                              handleFieldChange('ipv4_dst', target.value)
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                    </Grid>
                  </AltFormField>
                )}
              </Grid>
              <Grid item xs={4}>
                <AltFormField disableGutters label={'Direction'}>
                  <Select
                    fullWidth={true}
                    variant={'outlined'}
                    value={flow.match.direction}
                    onChange={({target}) => {
                      handleFieldChange('direction', target.value);
                    }}
                    input={<OutlinedInput id="direction" />}>
                    <MenuItem value={DIRECTION.UPLINK}>
                      <ListItemText primary={'Uplink'} />
                    </MenuItem>
                    <MenuItem value={DIRECTION.DOWNLINK}>
                      <ListItemText primary={'Downlink'} />
                    </MenuItem>
                  </Select>
                </AltFormField>
                {flow.match.ip_proto === PROTOCOL.IPPROTO_TCP && (
                  <AltFormField disableGutters label={'TCP'}>
                    <Grid container spacing={0}>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Source Port'}>
                          <OutlinedInput
                            data-testid="tcpSource"
                            placeholder="0"
                            fullWidth={true}
                            value={flow.match.tcp_src}
                            onChange={({target}) =>
                              handleFieldChange(
                                'tcp_src',
                                parseInt(target.value),
                              )
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Destination Port'}>
                          <OutlinedInput
                            data-testid="tcpDestination"
                            placeholder="0"
                            fullWidth={true}
                            value={flow.match.tcp_dst}
                            onChange={({target}) =>
                              handleFieldChange(
                                'tcp_dst',
                                parseInt(target.value),
                              )
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                    </Grid>
                  </AltFormField>
                )}
              </Grid>
              <Grid item xs={4}>
                <AltFormField disableGutters label={'Protocol'}>
                  <Select
                    fullWidth={true}
                    variant={'outlined'}
                    value={flow.match.ip_proto}
                    onChange={({target}) => {
                      handleFieldChange('ip_proto', target.value);
                    }}
                    input={<OutlinedInput id="protocol" />}>
                    <MenuItem value={PROTOCOL.IPPROTO_IP}>
                      <ListItemText primary={'IP'} />
                    </MenuItem>
                    <MenuItem value={PROTOCOL.IPPROTO_UDP}>
                      <ListItemText primary={'UDP'} />
                    </MenuItem>
                    <MenuItem value={PROTOCOL.IPPROTO_TCP}>
                      <ListItemText primary={'TCP'} />
                    </MenuItem>
                    <MenuItem value={PROTOCOL.IPPROTO_ICMP}>
                      <ListItemText primary={'ICMP'} />
                    </MenuItem>
                  </Select>
                </AltFormField>
                {flow.match.ip_proto === PROTOCOL.IPPROTO_UDP && (
                  <AltFormField disableGutters label={'UDP'}>
                    <Grid container spacing={0}>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Source Port'}>
                          <OutlinedInput
                            data-testid="tcpSource"
                            placeholder="0"
                            fullWidth={true}
                            value={flow.match.udp_src}
                            onChange={({target}) =>
                              handleFieldChange(
                                'udp_src',
                                parseInt(target.value),
                              )
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <AltFormFieldSubheading label={'Destination Port'}>
                          <OutlinedInput
                            data-testid="udpDestination"
                            placeholder="0"
                            fullWidth={true}
                            value={flow.match.udp_dst}
                            onChange={({target}) =>
                              handleFieldChange(
                                'udp_dst',
                                parseInt(target.value),
                              )
                            }
                          />
                        </AltFormFieldSubheading>
                      </Grid>
                    </Grid>
                  </AltFormField>
                )}
              </Grid>
            </Grid>
          </div>
        </AccordionDetails>
      </Accordion>
    </div>
  );
}

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyFlowsEdit(props: Props) {
  const handleAddFlow = () => {
    const flowList = [
      ...(props.policyRule.flow_list || []),
      {
        action: ACTION.DENY,
        match: {
          direction: DIRECTION.UPLINK,
          ip_proto: PROTOCOL.IPPROTO_IP,
        },
      },
    ];

    props.onChange({...props.policyRule, flow_list: flowList});
  };

  const onFlowChange = (index, flow) => {
    const flowList = [...(props.policyRule.flow_list || [])];
    flowList[index] = flow;
    props.onChange({...props.policyRule, flow_list: flowList});
  };

  const handleDeleteFlow = (index: number) => {
    const flowList = [...(props.policyRule.flow_list || [])];
    flowList.splice(index, 1);
    props.onChange({...props.policyRule, flow_list: flowList});
  };

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
            {"A policy's flows determines how it routes traffic"}
          </Typography>
          <ListItem dense disableGutters />
          {props.policyRule.flow_list.length > 0 && (
            <ListItem disableGutters>
              <Typography variant="h6">Flows</Typography>
            </ListItem>
          )}
          {(props.policyRule.flow_list || []).slice(0, 30).map((flow, i) => (
            <ListItem key={i} disableGutters>
              <PolicyFlowFields2
                index={i}
                flow={flow}
                handleDelete={handleDeleteFlow}
                onChange={onFlowChange}
              />
            </ListItem>
          ))}
          Add New Flow
          <IconButton onClick={handleAddFlow}>
            <AddCircleOutline />
          </IconButton>
        </List>
      </DialogContent>
    </>
  );
}
