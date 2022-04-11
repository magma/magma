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

import React from 'react';

import type {GetHealthResponse} from './GrafanaAPIType';
import type {Task} from './handlers';

export default function GrafanaErrorMessage(props: {
  completedTasks: Array<Task>,
  errorTask: Task,
  grafanaHealth: GetHealthResponse,
}) {
  return (
    <div>
      <h2>Grafana Debug Information:</h2>
      <p>
        An unrecoverable error was encountered while trying to load Grafana.
        Please contact your NMS administrator or open an issue on
        <a href="https://github.com/facebookincubator/magma/issues"> github </a>
        and include this error log.
      </p>
      <h3>Completed Tasks:</h3>
      <ul>
        {props.completedTasks.map((task, index) => {
          return (
            <li key={index}>
              <TaskMessage task={task} />
            </li>
          );
        })}
      </ul>
      <h3>Task Responsible for Error:</h3>
      <TaskMessage task={props.errorTask} />
      <h3> Grafana Health Information: </h3>
      <p>Commit: {props.grafanaHealth.commit}</p>
      <p>Database: {props.grafanaHealth.database}</p>
      <p>Version: {props.grafanaHealth.version}</p>
    </div>
  );
}

function TaskMessage(props: {task: Task}) {
  return (
    <div>
      <p>
        <strong>{props.task.name}</strong>
      </p>
      <p>Status: {props.task.status}</p>
      <p>Message: {JSON.stringify(props.task.message)}</p>
    </div>
  );
}
