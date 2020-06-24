/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import request from 'request';

import type {FrontendResponse} from './types';

function makeOptions(
  method: string,
  url: string,
  parentRequest: express$Request,
  body,
) {
  const options = {
    method,
    url,
    headers: {
      'x-auth-organization': parentRequest.headers['x-auth-organization'],
      'x-auth-user-email': parentRequest.headers['x-auth-user-email'],
      'x-auth-user-role': parentRequest.headers['x-auth-user-role'],
      cookie: parentRequest.headers['cookie'],
      'Content-Type': 'application/json',
    },
    body: undefined,
    // setting `json: true` would help by automatically converting, but it adds
    // `Accept: application/json` header that triggers conductor's bug #376,
    // see https://github.com/Netflix/conductor/issues/376
  };
  // If body is empty object, convert it to null.
  // Otherwise the http library will send a request
  // with Content-Length: 2 :/
  let modifiedBody = null;
  if (body && typeof body === 'object' && Object.keys(body).length > 0) {
    modifiedBody = JSON.stringify(body);
  }
  if (body != null) {
    options.body = modifiedBody;
  }
  return options;
}

function isSuccess(res) {
  return res.statusCode >= 200 && res.statusCode < 300;
}

function resolveSuccess(res, resolve) {
  const resCompatibileWithSuperagent = {
    text: res.body,
    statusCode: res.statusCode,
  };
  resolve(resCompatibileWithSuperagent);
}

function doHttpRequestWithOptions(options): Promise<FrontendResponse> {
  return new Promise<FrontendResponse>((resolve, reject) => {
    request(options, function(err, res) {
      if (res != null && isSuccess(res)) {
        resolveSuccess(res, resolve);
      } else if (res != null) {
        reject({message: 'Wrong status code', statusCode: res.statusCode});
      } else {
        reject(err);
      }
    });
  });
}

function doHttpRequest(
  method: string,
  url: string,
  parentRequest: express$Request,
  body,
): Promise<FrontendResponse> {
  return doHttpRequestWithOptions(
    makeOptions(method, url, parentRequest, body),
  );
}

const HttpClient = {
  // TODO: refactor usage so that get method can be simplified
  get: <T>(path: string, parentRequest: express$Request): Promise<T> =>
    new Promise<T>((resolve, reject) => {
      const options = makeOptions('GET', path, parentRequest);
      request(options, function(err, res) {
        if (res != null && isSuccess(res)) {
          // TODO get rid of explicit json parsing when
          // conductor's #376 is fixed
          resolve(JSON.parse(res.body)); // all GET methods return json
        } else if (res != null) {
          reject({message: 'Wrong status code', statusCode: res.statusCode});
        } else {
          reject(err);
        }
      });
    }),

  delete: <T>(
    path: string,
    data: T,
    parentRequest: ExpressRequest,
  ): Promise<FrontendResponse> => {
    return doHttpRequest('DELETE', path, parentRequest, data);
  },

  post: <T>(
    path: string,
    data: T,
    parentRequest: ExpressRequest,
  ): Promise<FrontendResponse> => {
    return doHttpRequest('POST', path, parentRequest, data);
  },

  put: <T>(
    path: string,
    data: T,
    parentRequest: ExpressRequest,
  ): Promise<FrontendResponse> => {
    return doHttpRequest('PUT', path, parentRequest, data);
  },
};

export default HttpClient;
