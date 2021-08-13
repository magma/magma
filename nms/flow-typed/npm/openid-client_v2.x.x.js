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
declare module 'openid-client' {
  import type {KeyStore} from 'node-jose';

  declare type HTTPOptions = {
    proxy?: string,
  };

  declare type HTTPClient = {
    HTTPError: Error,
    get: (url, options) => Promise<any>,
    post: (url, options) => Promise<any>,
  };

  declare export type OpenidUserInfoClaims = {
    name: string,
    preferred_username: string,
    given_name: string,
    family_name: string,
    email: string,
    jti: string,
    exp: string,
    nbf: string,
    iat: string,
    iss: string,
    aud: string,
    sub: string,
    typ: string,
    azp: string,
    auth_time: string,
    session_state: string,
    acr: string,
    email_verified: string,
  };

  declare type GrantRequestOptions = {
    grant_type: string,
    username?: string,
    password?: string,
    scope?: string,
  };

  declare export class Client {
    constructor({
      client_id: string,
      client_secret: string,
    }): Client;
    grant: (options: GrantRequestOptions) => Promise<TokenSet>;
    decryptIdToken: (tokenSet: TokenSet) => Promise<TokenSet>;
    validateIdToken: (
      tokenSet: TokenSet,
      nonce: ?string,
      returnedBy: string,
      maxAge: ?string,
      state: ?string,
    ) => Promise<TokenSet>;
    refresh: (refreshToken: string) => Promise<TokenSet>;
    issuer: OpenidIssuer;
    userinfo: (token: string) => Promise<OpenidUserInfoClaims>;
  }

  declare export class TokenSet {
    constructor(tokenSet: TokenSet | string | void): TokenSet;

    /**
     * This is the user's api key - it should be sent to
     * the auth server whenever a protected resource is requested
     */
    access_token: string;
    expires_at: number;
    refresh_expires_in: number;
    //access token is short lived, we get a new one using the refresh token
    refresh_token: string;
    scope: string;
    /**
     * oidc jwt containing information about the user
     * not to be trusted until we validate the signature
     **/
    id_token: string;
    // the TokenSet class decodes and caches these from the id_token
    claims: OpenidUserInfoClaims;
    session_state: string;
    token_type: string;
    expired: () => boolean;
  }

  declare export class Issuer {
    static discover: (url: string) => Promise<Issuer>;
    Client: Class<Client>;
    keystore: (refresh?: boolean) => KeyStore;
    static defaultHttpOptions?: HTTPOptions;
    static httpClient?: HTTPClient;
  }

  declare export class Strategy<TUser> {
    constructor(
      {
        client: Client,
        passReqToCallback?: boolean,
        params: {
          redirect_uri: string,
        },
      },
      callback: HandleOidcResponse<TUser>,
    ): Strategy<TUser>;
  }

  declare export type HandleOidcResponse<TUser> = {
    (
      req: any,
      tokenSet: TokenSet,
      claims: OpenidUserInfoClaims,
      done: (error: Error | void, user: TUser) => any,
    ): any,
  };
}
