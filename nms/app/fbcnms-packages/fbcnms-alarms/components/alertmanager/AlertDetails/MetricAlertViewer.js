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
 *
 */

import * as React from 'react';
import DescriptionIcon from '@material-ui/icons/Description';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import {Detail, ObjectViewer, Section} from './AlertDetailsPane';
import {useAlarmContext} from '../../AlarmContext';

import type {AlertViewerProps} from '../../rules/RuleInterface';

export default function MetricAlertViewer({alert}: AlertViewerProps) {
  const {filterLabels} = useAlarmContext();
  const {labels, annotations} = alert || {};
  const {alertname: _a, severity: _s, ...extraLabels} = labels || {};
  const {description, ...extraAnnotations} = annotations || {};
  return (
    <Grid container data-testid="metric-alert-viewer" spacing={2}>
      <Section title={'Details'}>
        <Detail icon={DescriptionIcon} title="Description">
          <Typography color="textSecondary">{description}</Typography>
        </Detail>
      </Section>
      <Section title={'Labels'}>
        <ObjectViewer
          object={filterLabels ? filterLabels(extraLabels) : extraLabels}
        />
      </Section>
      <Section title={'Annotations'} divider={false}>
        <ObjectViewer object={extraAnnotations} />
      </Section>
    </Grid>
  );
}
