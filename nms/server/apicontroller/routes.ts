/*
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

import HttpsProxyAgent, {HttpsProxyAgentOptions} from 'https-proxy-agent';
import auditLoggingDecorator from './auditLoggingDecorator';
import proxy from 'express-http-proxy';
import url from 'url';
import {API_HOST, apiCredentials} from '../../config/config';
import {AxiosError} from 'axios';
import {IncomingMessage, RequestOptions} from 'http';
import {NextFunction, Request, Response, Router} from 'express';
import {intersection} from 'lodash';

const router = Router();

const PROXY_TIMEOUT_MS = 30000;
const MUTATORS = ['POST', 'PUT', 'DELETE'];

let agent: HttpsProxyAgent | null = null;
if (process.env.HTTPS_PROXY) {
  const options = url.parse(process.env.HTTPS_PROXY) as HttpsProxyAgentOptions;
  agent = new HttpsProxyAgent(options);
}
const PROXY_OPTIONS = {
  https: true,
  memoizeHost: false,
  timeout: PROXY_TIMEOUT_MS,
  proxyReqOptDecorator: (proxyReqOpts: RequestOptions) => {
    return {
      ...proxyReqOpts,
      agent: agent!,
      cert: apiCredentials().cert,
      key: apiCredentials().key,
      rejectUnauthorized: false,
    };
  },
  proxyReqPathResolver: (req: Request) =>
    req.originalUrl.replace(/^\/nms\/apicontroller/, ''),
};

export async function apiFilter(req: Request): Promise<boolean> {
  if (req.user.isReadOnlyUser && MUTATORS.includes(req.method)) {
    return false;
  }

  if (req.organization) {
    const organization = await req.organization();

    // If the request isn't an organization network, block
    // the request
    const isOrganizationAllowed = containsNetworkID(
      organization.networkIDs,
      req.params.networkID,
    );
    if (!isOrganizationAllowed) {
      return false;
    }
  }

  // super users on standalone deployments
  // have access to all proxied API requests
  // for the organization
  if (req.user.isSuperUser) {
    return true;
  }
  return containsNetworkID(req.user.networkIDs, req.params.networkID);
}

export async function networksResponseDecorator(
  _proxyRes: IncomingMessage,
  proxyResData: Buffer,
  userReq: Request,
) {
  let result = JSON.parse(proxyResData.toString('utf8')) as Array<string>;
  if (userReq.organization) {
    const organization = await userReq.organization();
    result = intersection(result, organization.networkIDs);
  }
  if (!userReq.user.isSuperUser) {
    // the list of networks is further restricted to what the user
    // is allowed to see
    result = intersection(result, userReq.user.networkIDs);
  }
  return JSON.stringify(result);
}

const containsNetworkID = function (
  allowedNetworkIDs: Array<string>,
  networkID: string,
): boolean {
  return (
    allowedNetworkIDs.indexOf(networkID) !== -1 ||
    // Remove secondary condition after T34404422 is addressed. Reason:
    //   Request needs to be lower cased otherwise calling
    //   MagmaAPIUrls.gateways() potentially returns missing devices.
    allowedNetworkIDs
      .map(id => id.toString().toLowerCase())
      .indexOf(networkID.toString().toLowerCase()) !== -1
  );
};

const proxyErrorHandler = (err: unknown, res: Response, next: NextFunction) => {
  if ((err as AxiosError).code === 'ENOTFOUND') {
    res.status(503).send('Cannot reach Orchestrator server');
  } else {
    console.error('andreilee here: ', err, res);
    next();
  }
};

router.use(
  /^\/magma\/v1\/feg_lte$/,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    userResDecorator: networksResponseDecorator,
    proxyErrorHandler,
  }),
);

router.use(
  /^\/magma\/v1\/networks$/,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    userResDecorator: networksResponseDecorator,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/v1/networks/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: apiFilter,
    userResDecorator: auditLoggingDecorator,
    proxyErrorHandler,
  }),
);

const networkTypeRegex = '(cwf|feg|lte|feg_lte)';
router.use(
  `/magma/v1/:networkType(${networkTypeRegex})/:networkID`,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: apiFilter,
    userResDecorator: auditLoggingDecorator,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: req => req.method === 'GET',
  }),
);

router.use(
  '/magma/v1/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: req => req.method === 'GET',
  }),
);

router.use(
  '/magma/v1/tenants/targets_metadata',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: req => req.method === 'GET',
  }),
);

router.use(
  '/magma/v1/events/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: apiFilter,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/v1/events/:networkID/:streamName',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: apiFilter,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/v1/dp/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: apiFilter,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/v1/about/version',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: req => req.method === 'GET',
  }),
);

router.use('', (req, res) => {
  if (req.user.isReadOnlyUser && MUTATORS.includes(req.method)) {
    res.status(403).send('Mutation forbidden. Readonly access');
    return;
  }
  res.status(404).send('Not Found');
});

export default router;
