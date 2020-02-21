/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {flow_description} from '@fbcnms/magma-api';

import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControl from '@material-ui/core/FormControl';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';

import {ACTION, DIRECTION, PROTOCOL} from './PolicyTypes';
import {makeStyles} from '@material-ui/styles';

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
  },
  expanded: {},
  block: {
    display: 'block',
  },
  flex: {display: 'flex'},
  panel: {flexGrow: 1},
  removeIcon: {alignSelf: 'baseline'},
}));

type ActionType = $Keys<typeof ACTION>;
type Props = {
  index: number,
  flow: flow_description,
  handleActionChange: (number, ActionType) => void,
  handleFieldChange: (number, string, string | number) => void,
  handleDelete: number => void,
};

export default function PolicyFlowFields(props: Props) {
  const classes = useStyles();
  const {flow} = props;

  return (
    <div className={classes.flex}>
      <ExpansionPanel className={classes.panel}>
        <ExpansionPanelSummary
          classes={{root: classes.root, expanded: classes.expanded}}
          expandIcon={<ExpandMoreIcon />}>
          <Text variant="body2">Flow {props.index + 1}</Text>
        </ExpansionPanelSummary>
        <ExpansionPanelDetails classes={{root: classes.block}}>
          <div className={classes.flex}>
            <FormControl className={classes.input}>
              <InputLabel htmlFor="action">Action</InputLabel>
              <TypedSelect
                items={{
                  [ACTION.PERMIT]: 'Permit',
                  [ACTION.DENY]: 'Deny',
                }}
                value={flow.action}
                onChange={val => props.handleActionChange(props.index, val)}
                input={<Input id="action" />}
              />
            </FormControl>
            <FormControl className={classes.input}>
              <InputLabel htmlFor="direction">Direction</InputLabel>
              <TypedSelect
                items={{
                  [DIRECTION.UPLINK]: 'Uplink',
                  [DIRECTION.DOWNLINK]: 'Downllink',
                }}
                value={flow.match.direction}
                onChange={val =>
                  props.handleFieldChange(props.index, 'direction', val)
                }
                input={<Input id="direction" />}
              />
            </FormControl>
            <FormControl className={classes.input}>
              <InputLabel htmlFor="protocol">Protocol</InputLabel>
              <TypedSelect
                items={{
                  [PROTOCOL.IPPROTO_IP]: 'IP',
                  [PROTOCOL.IPPROTO_UDP]: 'UDP',
                  [PROTOCOL.IPPROTO_TCP]: 'TCP',
                  [PROTOCOL.IPPROTO_ICMP]: 'ICMP',
                }}
                value={flow.match.ip_proto}
                onChange={val =>
                  props.handleFieldChange(props.index, 'ip_proto', val)
                }
                input={<Input id="protocol" />}
              />
            </FormControl>
          </div>
          {flow.match.ip_proto !== PROTOCOL.IPPROTO_ICMP && (
            <div className={classes.flex}>
              <TextField
                className={classes.input}
                label="IPv4 Source"
                margin="normal"
                value={flow.match.ipv4_src}
                onChange={({target}) =>
                  props.handleFieldChange(props.index, 'ipv4_src', target.value)
                }
              />
              <TextField
                className={classes.input}
                label="IPv4 Destination"
                margin="normal"
                value={flow.match.ipv4_dst}
                onChange={({target}) =>
                  props.handleFieldChange(props.index, 'ipv4_dst', target.value)
                }
              />
            </div>
          )}
          {flow.match.ip_proto === PROTOCOL.IPPROTO_UDP && (
            <div className={classes.flex}>
              <TextField
                className={classes.input}
                label="UDP Source Port"
                margin="normal"
                value={flow.match.udp_src}
                onChange={({target}) =>
                  props.handleFieldChange(
                    props.index,
                    'udp_src',
                    parseInt(target.value),
                  )
                }
              />
              <TextField
                className={classes.input}
                label="UDP Destination Port"
                margin="normal"
                value={flow.match.udp_dst}
                onChange={({target}) =>
                  props.handleFieldChange(
                    props.index,
                    'udp_dst',
                    parseInt(target.value),
                  )
                }
              />
            </div>
          )}
          {flow.match.ip_proto === PROTOCOL.IPPROTO_TCP && (
            <div className={classes.flex}>
              <TextField
                className={classes.input}
                label="TCP Source Port"
                margin="normal"
                value={flow.match.tcp_src}
                onChange={({target}) =>
                  props.handleFieldChange(
                    props.index,
                    'tcp_src',
                    parseInt(target.value),
                  )
                }
              />
              <TextField
                className={classes.input}
                label="TCP Destination Port"
                margin="normal"
                value={flow.match.tcp_dst}
                onChange={({target}) =>
                  props.handleFieldChange(
                    props.index,
                    'tcp_dst',
                    parseInt(target.value),
                  )
                }
              />
            </div>
          )}
        </ExpansionPanelDetails>
      </ExpansionPanel>
      <IconButton
        className={classes.removeIcon}
        onClick={() => props.handleDelete(props.index)}>
        <RemoveCircleOutline />
      </IconButton>
    </div>
  );
}
