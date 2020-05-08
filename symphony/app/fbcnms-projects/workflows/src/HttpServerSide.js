/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import request from 'superagent';

const HttpClient = {
  get: <T>(path: string, parentRequest: express$Request): Promise<T> =>
    new Promise<T>((resolve, reject) => {
      const req = request.get(path).accept('application/json');
      req.header['x-auth-organization'] =
        parentRequest.headers['x-auth-organization'];
      req.end((err, res) => {
        if (err) {
          if (res && res.error) {
            resolve(res.error.status);
            console.error(res.error.message);
          } else {
            reject(err);
          }
        } else {
          resolve(res.body);
        }
      });
    }),
  delete: <T>(path: string, data: T, parentRequest: express$Request) =>
    new Promise<T>((resolve, reject) => {
      // If data is empty object, convert it to null.
      // Otherwise the http library will send a request
      // with Content-Length: 2 :/
      let modifiedData = data;
      if (data && typeof data === 'object' && Object.keys(data).length === 0) {
        modifiedData = null;
      }
      const req = request.delete(path, modifiedData).accept('application/json');
      req.header['x-auth-organization'] =
        parentRequest.headers['x-auth-organization'];
      req.end((err, res) => {
        if (err) {
          resolve(err);
          reject(err);
        } else {
          if (res) {
            resolve(res);
          }
        }
      });
    }),

  post: <T>(path: string, data: T, parentRequest: express$Request) =>
    new Promise<T>((resolve, reject) => {
      const req = request
        .post(path, data)
        .set('Content-Type', 'application/json');
      req.header['x-auth-organization'] =
        parentRequest.headers['x-auth-organization'];
      req.end((err, res) => {
        if (err || !res.ok) {
          console.error('Error on post! ' + res);
          reject(err);
        } else {
          if (res) {
            resolve(res);
          }
        }
      });
    }),

  put: <T>(path: string, data: T, parentRequest: express$Request) =>
    new Promise<{}>((resolve, reject) => {
      const req = request.put(path, data).set('Accept', 'application/json');
      req.header['x-auth-organization'] =
        parentRequest.headers['x-auth-organization'];
      req
        .then(res => {
          resolve(res);
        })
        .catch(error => {
          reject(error);
        });
    }),
};

export default HttpClient;
