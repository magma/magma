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
import {Link} from 'react-router-dom';

import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  link: {
    textDecoration: 'none',
  },
}));

type Props = {
  children: React.Node,
  to: string,
  className?: string,
};

function NestedRouteLink(props: Props, ref: React.Ref<*>) {
  const classes = useStyles();
  const {children, to, className: childClassName, ...childProps} = props;
  return (
    <Link
      {...childProps}
      innerRef={ref}
      className={classNames(classes.link, childClassName)}
      to={to}>
      {children}
    </Link>
  );
}

export default React.forwardRef<Props, typeof NestedRouteLink>(NestedRouteLink);
