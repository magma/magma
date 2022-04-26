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

import React, {useCallback, useState} from 'react';
import Text from './design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import {Link} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../../../fbc_js_core/ui/hooks';
import {colors} from "../../../app/theme/default";

const useStyles = makeStyles(() => ({
  icon: {
    color: colors.primary.gullGray,
    display: 'flex',
  },
  link: {
    width: '100%',
  },
  root: {
    display: 'flex',
    justifyContent: 'center',
    width: '100%',
    padding: '15px 0px',
    '&:hover $icon, &$selected $icon': {
      color: colors.primary.white,
    },
  },
  selected: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
  tooltip: {
    position: 'relative',
    '&&': {
      padding: '8px 12px',
      backgroundColor: colors.primary.brightGray,
    },
  },
  arrow: {
    position: 'absolute',
    left: '-8px',
    '&:before': {
      borderBottom: '4px solid transparent',
      borderLeft: '4px solid transparent',
      borderRight: `4px solid ${colors.primary.brightGray}`,
      borderTop: '4px solid transparent',
      top: '-5px',
      content: '""',
      position: 'absolute',
      zIndex: 10,
    },
  },
  bootstrapPlacementLeft: {
    margin: '0 8px',
  },
  tooltipLabel: {
    '&&': {
      fontSize: '12px',
      lineHeight: '16px',
      color: colors.primary.white,
      fontWeight: 'bold',
    },
  },
}));

type Props = {
  path: string,
  label: string,
  icon: any,
  hidden: boolean,
  onClick?: ?() => void,
};

export default function NavListItem(props: Props) {
  const {hidden, onClick} = props;
  const classes = useStyles();
  const router = useRouter();
  const [arrowArrow, setArrowRef] = useState(null);
  const handleArrowRef = useCallback(node => {
    if (node !== null) {
      setArrowRef(node);
    }
  }, []);

  if (hidden) {
    return null;
  }

  const isSelected = router.location.pathname.includes(props.path);

  return (
    <Link
      to={props.path}
      className={classes.link}
      onClick={() => onClick && onClick()}>
      <Tooltip
        placement="right"
        title={
          <>
            <Text className={classes.tooltipLabel} variant="body2">
              {props.label}
            </Text>
            <span className={classes.arrow} ref={handleArrowRef} />
          </>
        }
        classes={{
          tooltip: classes.tooltip,
          tooltipPlacementLeft: classes.bootstrapPlacementLeft,
        }}
        PopperProps={{
          popperOptions: {
            modifiers: {
              arrow: {
                enabled: Boolean(arrowArrow),
                element: arrowArrow,
              },
            },
          },
        }}>
        <div
          className={classNames({
            [classes.root]: true,
            [classes.selected]: isSelected,
          })}>
          <div className={classes.icon}>{props.icon}</div>
        </div>
      </Tooltip>
    </Link>
  );
}

NavListItem.defaultProps = {
  hidden: false,
};
