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

import * as React from 'react';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'inline-flex',
    alignItems: 'center',
    minHeight: '24px',
    boxSizing: 'initial',
  },
  followingText: {
    marginLeft: '8px',
    '&$wide': {
      marginLeft: '12px',
    },
  },
  wide: {},
}));

type Variant = 'body2' | 'subtitle2';
type Margin = 'regular' | 'wide';

export type TextPairingContainerProps = $ReadOnly<{|
  title: React.Node,
  variant?: ?Variant,
  margin?: ?Margin,
  disabled?: ?boolean,
  className?: ?string,
|}>;

type Props = $ReadOnly<{|
  children: React.Node,
  ...TextPairingContainerProps,
|}>;

const TextPairingContainer = (props: Props) => {
  const {
    children,
    title,
    variant,
    margin = 'regular',
    disabled: propDisabled = false,
    className,
  } = props;
  const classes = useStyles();
  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <div className={classNames(classes.root, className)}>
      {children}
      {title == null ? null : (
        <Text
          className={classNames(classes.followingText, {
            [classes.wide]: margin === 'wide',
          })}
          variant={variant || 'body2'}
          color={disabled ? 'gray' : undefined}>
          {title}
        </Text>
      )}
    </div>
  );
};

export default TextPairingContainer;
