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

import type {flow_description} from '../../../generated/MagmaAPIBindings';

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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TextField from '@material-ui/core/TextField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import TypedSelect from '../TypedSelect';

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

type Props = {
  index: number,
  flow: flow_description,
  handleDelete: number => void,
  onChange: (number, flow_description) => void,
};

export default function PolicyFlowFields(props: Props) {
  const classes = useStyles();
  const {flow} = props;

  const handleActionChange = action =>
    props.onChange(props.index, {
      ...props.flow,
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
                onChange={handleActionChange}
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
                onChange={val => handleFieldChange('direction', val)}
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
                onChange={val => handleFieldChange('ip_proto', val)}
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
                value={flow.match.ip_src?.address ?? ''}
                onChange={({target}) => {
                  handleFieldChange('ip_src', {
                    address: target.value,
                    version: 'IPv4',
                  });
                }}
              />
              <TextField
                className={classes.input}
                label="IPv4 Destination"
                margin="normal"
                value={flow.match.ip_dst?.address ?? ''}
                onChange={({target}) => {
                  handleFieldChange('ip_dst', {
                    address: target.value,
                    version: 'IPv4',
                  });
                }}
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
                  handleFieldChange('udp_src', parseInt(target.value))
                }
              />
              <TextField
                className={classes.input}
                label="UDP Destination Port"
                margin="normal"
                value={flow.match.udp_dst}
                onChange={({target}) =>
                  handleFieldChange('udp_dst', parseInt(target.value))
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
                  handleFieldChange('tcp_src', parseInt(target.value))
                }
              />
              <TextField
                className={classes.input}
                label="TCP Destination Port"
                margin="normal"
                value={flow.match.tcp_dst}
                onChange={({target}) =>
                  handleFieldChange('tcp_dst', parseInt(target.value))
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
