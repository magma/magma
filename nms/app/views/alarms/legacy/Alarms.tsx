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
 */

import * as React from 'react';
import AppContext from '../../../components/context/AppContext';
import FBCAlarms from '../components/Alarms';
import type {ApiUtil} from '../components/AlarmsApi';
import type {Labels} from '../components/AlarmAPIType';

type Props = {
  apiUtil: ApiUtil;
};

export default function Alarms(props: Props) {
  const {apiUtil} = props;
  const {isFeatureEnabled} = React.useContext(AppContext);
  const disabledTabs = React.useMemo(
    () =>
      [
        isFeatureEnabled('alert_receivers') ? null : 'receivers',
        isFeatureEnabled('alert_routes') ? null : 'routes',
        isFeatureEnabled('alert_suppressions') ? null : 'suppressions',
      ].filter(Boolean) as Array<string>,
    [isFeatureEnabled],
  );
  return (
    <FBCAlarms
      apiUtil={apiUtil}
      makeTabLink={({networkId, keyName}) =>
        `/nms/${networkId || ''}/alerts/${keyName}`
      }
      disabledTabs={disabledTabs}
      thresholdEditorEnabled={true}
      filterLabels={filterSymphonyLabels}
    />
  );
}

/**
 * Filters out hidden system labels from the firing alerts table
 */
function filterSymphonyLabels(labels: Labels) {
  const labelsToFilter = ['monitor', 'networkID'];
  const filtered = {...labels};
  for (const label of labelsToFilter) {
    delete filtered[label];
  }
  return filtered;
}
