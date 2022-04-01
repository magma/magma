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

import * as React from 'react';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import InfoTinyIcon from '../Icons/Indications/InfoTinyIcon';
import Text from '../Text';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  tooltip: {
    borderRadius: '2px',
    backgroundColor: symphony.palette.secondary,
    padding: '8px 10px',
    width: '181px',
  },
  iconContainer: {
    display: 'inline-flex',
    '&:hover $icon': {
      fill: symphony.palette.primary,
    },
  },
  icon: {},
  arrow: {
    position: 'absolute',
    bottom: '0',
    left: '8px',
    '&:before': {
      borderBottom: '4px solid transparent',
      borderLeft: '4px solid transparent',
      borderRight: '4px solid transparent',
      borderTop: `4px solid ${symphony.palette.secondary}`,
      botom: '-8px',
      content: '""',
      position: 'absolute',
      zIndex: 10,
    },
  },
}));

type Props = {
  description: React.Node,
};

const InfoTooltip = ({description}: Props) => {
  const classes = useStyles();
  return (
    <BasePopoverTrigger
      position="above"
      popover={
        <div className={classes.tooltip}>
          <Text color="light" variant="caption">
            {description}
          </Text>
          <span className={classes.arrow} />
        </div>
      }>
      {(onShow, onHide, contextRef) => (
        <div
          className={classes.iconContainer}
          ref={contextRef}
          onMouseOver={onShow}
          onMouseOut={onHide}>
          <InfoTinyIcon color="gray" className={classes.icon} />
        </div>
      )}
    </BasePopoverTrigger>
  );
};

export default InfoTooltip;
