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

import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
import classNames from 'classnames';
import {Link, useLocation, useResolvedPath} from 'react-router-dom';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  icon: {
    color: colors.primary.gullGray,
    display: 'flex',
    justifyContent: 'center',
  },
  root: {
    display: 'flex',
    width: '100%',
    textDecoration: 'none',
    alignItems: 'center',
    padding: '15px 28px',
    outline: 'none',
    '&:hover $icon, &:hover $label, &:focus $icon, &:focus $label, &$selected $icon, &$selected $label': {
      color: colors.primary.white,
    },
  },
  selected: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
  label: {
    '&&': {
      color: colors.primary.gullGray,
      whiteSpace: 'nowrap',
      paddingLeft: '16px',
    },
  },
}));

type Props = {
  path: string,
  label: string,
  icon: any,
  expanded: boolean,
};

export default function SidebarItem(props: Props) {
  const classes = useStyles();
  const location = useLocation();
  const resolvedPath = useResolvedPath(props.path);
  const isSelected = location.pathname.startsWith(resolvedPath.pathname);

  return (
    <Link
      to={props.path}
      className={classNames({
        [classes.root]: true,
        [classes.selected]: isSelected,
      })}>
      <div className={classes.icon}>{props.icon}</div>
      {props.expanded && (
        <Text className={classes.label} variant="body3">
          {props.label}
        </Text>
      )}
    </Link>
  );
}
