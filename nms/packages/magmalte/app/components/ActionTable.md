```jsx
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

/**
 * This file is meant to be viewed using React Styleguidist
 * See the NMS readme.md
 */

import Chip from '@material-ui/core/Chip';
import defaultTheme from '../theme/default';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
import {makeStyles} from '@material-ui/styles';

const asdf = 'fdsa';

/*
type AlertRowType = {
  alertName: string,
  labels: {[string]: string},
  status: string,
  service: string,
  gatewayId: string,
  date: Date,
};
*/

const tableData = [ // Array<AlertRowType>
  {
    alertName: 'High Gateway CPU Load',
    labels: {},
    status: 'ACTIVE',
    service: 'magmad',
    gatewayId: 'nyc_soho_1',
    date: 'December 25, 2000'
  },
  {
    alertName: 'High Gateway CPU Load',
    labels: {},
    status: 'ACTIVE',
    service: 'magmad',
    gatewayId: 'nyc_soho_2',
    date: 'December 25, 2000'
  },
  {
    alertName: 'High Gateway CPU Load',
    labels: {},
    status: 'ACTIVE',
    service: 'magmad',
    gatewayId: 'yyz_1',
    date: 'December 25, 2000'
  }
];

<ThemeProvider theme={defaultTheme}>
  <ActionTable
    data={tableData}
    columns={[
      {title: 'Date', field: 'date', type: 'datetime', defaultSort: 'desc'},
      {title: 'Status', field: 'status', width: 200},
      {title: 'Alert Name', field: 'alertName', width: 200},
      {title: 'Service', field: 'service', width: 200},
      {title: 'Gateway', field: 'gatewayId', width: 200},
      {
        title: 'Labels',
        field: 'labels',
        render: (currRow/*: AlertRowType*/) => (
          <div>
            {Object.keys(currRow.labels)
              .filter(k => !ignoreLabelList.includes(k))
              .map(k => (
                <Chip
                  key={k}

                  label={
                    <span>
                      <em>{k}</em>={currRow.labels[k]}
                    </span>
                  }
                  size="small"
                />
              ))}
          </div>
        ),
      },
    ]}
    options={{
      actionsColumnIndex: -1,
      pageSizeOptions: [5, 10],
      toolbar: false,
    }}
    localization={{
      header: {actions: ''},
    }}
  />
</ThemeProvider>
```
