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
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {typography} from '../default';

const styles = {
  h1: typography.h1,
  h2: typography.h2,
  h3: typography.h3,
  h4: typography.h4,
  h5: typography.h5,
  h6: typography.h6,
  subtitle1: typography.subtitle1,
  subtitle2: typography.subtitle2,
  subtitle3: typography.subtitle3,
  body1: typography.body1,
  body2: typography.body2,
  body3: typography.body3,
  caption: typography.caption,
  overline: typography.overline,
  lightWeight: {
    fontWeight: 300,
  },
  regularWeight: {
    fontWeight: 400,
  },
  mediumWeight: {
    fontWeight: 500,
  },
  boldWeight: {
    fontWeight: 600,
  },
  truncate: {
    display: 'block',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
};
export const typographyStyles = makeStyles<Props, typeof styles>(() => styles);

type Props = {
  children: ?React.Node,
  variant?:
    | 'h1'
    | 'h2'
    | 'h3'
    | 'h4'
    | 'h5'
    | 'h6'
    | 'subtitle1'
    | 'subtitle2'
    | 'subtitle3'
    | 'body1'
    | 'body2'
    | 'body3'
    | 'caption'
    | 'overline',
  className?: string,
  useEllipsis?: ?boolean,
  weight?: 'inherit' | 'light' | 'regular' | 'medium' | 'bold',
};

const Text = (props: Props) => {
  const {
    children,
    variant = 'body1',
    className,
    weight = 'inherit',
    useEllipsis = false,
    ...rest
  } = props;
  const classes = typographyStyles();
  return (
    <span
      {...rest}
      className={classNames(
        classes[variant],
        classes[`${weight ? weight : 'regular'}Weight`],
        {[classes.truncate]: useEllipsis},
        className,
      )}>
      {children}
    </span>
  );
};

export default Text;
