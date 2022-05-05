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
import Text from '../../../../../app/theme/design-system/Text';
import classNames from 'classnames';
import {colors} from '../../../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useInput} from './InputContext';
import {useMemo} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    color: colors.primary.comet,
    alignItems: 'center',
    marginRight: '8px',
    marginLeft: '4px',
  },
  hasValue: {},
  disabled: {
    '&:not($hasValue) $text': {
      color: colors.primary.gullGray,
    },
    pointerEvents: 'none',
    opacity: 0.5,
  },
  text: {
    color: colors.primary.comet,
  },
  clickable: {
    cursor: 'pointer',
  },
}));

type Props = {
  className?: string,
  children: React.Node,
  onClick?: () => void,
  inheritsDisability?: boolean,
};

const InputAffix = (props: Props) => {
  const {children, className, onClick, inheritsDisability = false} = props;
  const classes = useStyles();
  const {disabled: parentInputDisabled, value} = useInput();

  const {disabled: contextDisabled} = useFormElementContext();

  const disabled = useMemo(
    () => (parentInputDisabled && inheritsDisability) || contextDisabled,
    [parentInputDisabled, inheritsDisability, contextDisabled],
  );

  return (
    <div
      className={classNames(
        classes.root,
        {
          [classes.disabled]: disabled,
          [classes.hasValue]: Boolean(value),
          [classes.clickable]: onClick !== undefined && !disabled,
        },
        className,
      )}
      onClick={disabled ? null : onClick}>
      {typeof children === 'string' ? (
        <Text className={classes.text} variant="body2">
          {children}
        </Text>
      ) : (
        children
      )}
    </div>
  );
};

export default InputAffix;
