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

import type {
  flow_description,
  policy_rule,
} from '../../../generated/MagmaAPIBindings';

import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';

import {
  ACTION,
  DIRECTION,
  PROTOCOL,
} from '../../components/network/PolicyTypes';
// $FlowFixMe migrated to typescript
import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';
import {useState} from 'react';

const useStyles = makeStyles(() => policyStyles);

type FieldProps = {
  index: number,
  flow: flow_description,
  handleDelete: number => void,
  onChange: (number, flow_description) => void,
};

function PolicyFlowFields2(props: FieldProps) {
  const classes = useStyles();
  const {flow} = props;
  const [ipAddrType, setIPAddrType] = useState<'IPv4' | 'IPv6'>('IPv4');
  const handleActionChange = action =>
    props.onChange(props.index, {
      ...props.flow,
      // $FlowIgnore: value guaranteed to match the string literals
      action,
    });

  const handleFieldChange = (field: string, value: number | string | {}) =>
    props.onChange(props.index, {
      ...props.flow,
      match: {
        ...props.flow.match,
        [field]: value,
      },
    });

  return (
    <div className={classes.flex}>
      <Accordion defaultExpanded className={classes.panel}>
        <AccordionSummary
          classes={{
            root: classes.root,
            expanded: classes.expanded,
          }}
          expandIcon={<ExpandMoreIcon />}>
          <Grid container justifyContent="space-between">
            <Grid item className={classes.title}>
              <Text weight="medium" variant="body2">
                Flow {props.index + 1}
              </Text>
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
                  <>
                    <AltFormField disableGutters label={'IP'}>
                      <Grid container spacing={0}>
                        <Grid item xs={12}>
                          <AltFormFieldSubheading label={'Source IP'}>
                            <OutlinedInput
                              data-testid="ipSrc"
                              placeholder="192.168.0.1/24"
                              fullWidth={true}
                              value={flow.match.ip_src?.address ?? ''}
                              onChange={({target}) =>
                                handleFieldChange('ip_src', {
                                  address: target.value,
                                  version: ipAddrType,
                                })
                              }
                            />
                          </AltFormFieldSubheading>
                        </Grid>
                        <Grid item xs={12}>
                          <AltFormFieldSubheading label={'Destination IP'}>
                            <OutlinedInput
                              data-testid="ipDest"
                              placeholder="192.168.0.1/24"
                              fullWidth={true}
                              value={flow.match.ip_dst?.address ?? ''}
                              onChange={({target}) =>
                                handleFieldChange('ip_dst', {
                                  address: target.value,
                                  version: ipAddrType,
                                })
                              }
                            />
                          </AltFormFieldSubheading>
                        </Grid>
                      </Grid>
                    </AltFormField>
                    <ListItem disableGutters={true}>
                      <ToggleButtonGroup
                        size="small"
                        value={ipAddrType}
                        exclusive
                        onChange={(_, nextAddrType) =>
                          setIPAddrType(nextAddrType)
                        }>
                        <ToggleButton value="IPv4">{'IPv4'}</ToggleButton>
                        <ToggleButton value="IPv6">{'IPv6'}</ToggleButton>
                      </ToggleButtonGroup>
                    </ListItem>
                  </>
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
                      <Grid item xs={12}>
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
                      <Grid item xs={12}>
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
                      <Grid item xs={12}>
                        <AltFormFieldSubheading label={'Source Port'}>
                          <OutlinedInput
                            data-testid="udpSource"
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
                      <Grid item xs={12}>
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
  inputClass: string,
};

export default function PolicyFlowsEdit(props: Props) {
  const classes = useStyles();
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

  const flowList = props.policyRule.flow_list || [];
  return (
    <div data-testid="flowEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {"A policy's flows determines how it routes traffic"}
      </Text>
      <ListItem dense disableGutters />
      {props.policyRule.flow_list && props.policyRule.flow_list.length > 0 && (
        <ListItem disableGutters>
          <Text weight="medium" variant="subtitle1">
            Flows
          </Text>
        </ListItem>
      )}
      {flowList.slice(0, 30).map((flow, i) => (
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
      <IconButton data-testid="addFlowButton" onClick={handleAddFlow}>
        <AddCircleOutline />
      </IconButton>
    </div>
  );
}
