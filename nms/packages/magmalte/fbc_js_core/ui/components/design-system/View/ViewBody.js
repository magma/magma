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

import type {CombinedRefs} from '../hooks/useCombinedRefs';

import * as React from 'react';
import classNames from 'classnames';
import useCombinedRefs from '../hooks/useCombinedRefs';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useRef, useState} from 'react';

const paddingRight = 24;
const scrollWidth = 12;
const useStyles = makeStyles(() => ({
  viewWrapper: {
    flexGrow: 1,
    overflowX: 'hidden',
    overflowY: 'auto',
    display: 'flex',
    '&:not($plain)': {
      padding: `8px ${paddingRight}px 4px 24px`,
      '&$withScrollY': {
        paddingRight: `${paddingRight - scrollWidth}px`,
      },
    },
  },
  withScrollY: {},
  idented: {},
  plain: {},
}));

export const VARIANTS = {
  idented: 'idented',
  plain: 'plain',
};

export type Variant = $Keys<typeof VARIANTS>;

type Props = $ReadOnly<{|
  children: React.Node,
  variant?: ?Variant,
|}>;

const ViewBody = React.forwardRef<Props, HTMLElement>((props, ref) => {
  const {children, variant = VARIANTS.idented} = props;
  const classes = useStyles();
  const refs: CombinedRefs = [useRef(null), ref];
  const combinedRef = useCombinedRefs(refs);
  const [hasScrollY, setHasScrollY] = useState(false);

  useEffect(() => {
    window.requestAnimationFrame(() => {
      const element = combinedRef.current;
      setHasScrollY(element && element.scrollHeight > element.offsetHeight);
    });
  });

  return (
    <div
      ref={combinedRef}
      className={classNames(
        {
          [classes.withScrollY]: hasScrollY,
        },
        variant ? classes[variant] : null,
        classes.viewWrapper,
      )}>
      {children}
    </div>
  );
});

export default ViewBody;
