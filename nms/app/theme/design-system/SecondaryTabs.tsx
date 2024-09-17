/*
 * Copyright 2022 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import withStyles from '@mui/styles/withStyles';
import {colors, typography} from '../default';

export const MagmaTabs = withStyles({
  indicator: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
})(Tabs);

export const MagmaTab = withStyles({
  root: {
    fontFamily: typography.body1.fontFamily,
    fontWeight: typography.body1.fontWeight,
    fontSize: typography.body1.fontSize,
    lineHeight: typography.body1.lineHeight,
    letterSpacing: typography.body1.letterSpacing,
    color: colors.primary.brightGray,
    textTransform: 'none',
    opacity: 0.7,
    '&.Mui-selected': {
      opacity: 1,
      color: colors.primary.brightGray,
    },
  },
})(Tab);
