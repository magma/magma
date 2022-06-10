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

import type {AccessRoleLevel} from '../../shared/roles';

const path = require('path');

const {AccessRoles} = require('../../shared/roles');
const {ErrorCodes} = require('../../shared/errorCodes');
const addQueryParamsToUrl = require('./util').addQueryParamsToUrl;
const logger = require('../../shared/logging').getLogger(module);
const openRoutes = require('./openRoutes').default;

import type {ExpressResponse, NextFunction} from 'express';
import type {FBCNMSPassportRequest} from './passport';

type Options = {loginUrl: string};
// Final type, thus naming it as thus.
export type FBCNMSRequest = FBCNMSPassportRequest & { access: Options }; 
