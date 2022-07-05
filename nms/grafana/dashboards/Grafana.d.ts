/**
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

declare module 'grafana-dash-gen' {
  type RowParam = {
    title: string;
  };
  export class Row {
    addPanel(pannel: any): void;

    constructor(param: RowParam);
  }

  type DashboardParam = {
    schemaVersion: number;
    title: string;
    templating: Array<TemplateConfig>;
    description: string;
    rows: Array<Row>;
  };

  type Option = {
    selected: boolean;
    text: string;
    value: string;
  };
  export class Dashboard {
    state: {
      editable: boolean;
      templating: {
        list: Array<{
          type: string;
          includeAll: boolean;
          options: Array<Option>;
          current: Option;
        }>;
      };
    };

    constructor(param: DashboardParam);

    generate(): Record<string, any>;
  }

  type GraphParams = {
    title: string;
    span: number;
    datasource: string;
    description: string;
  };
  export declare namespace Panels {
    export class Graph {
      state: {
        targets: Array<{expr: string; legendFormat?: string}>;
        grid: {leftMin: number | null; leftMax: number};
        legend: {max: boolean; avg: boolean};
        y_formats: Array<string>;
      };
      constructor(param: GraphParams);
    }
  }
}
