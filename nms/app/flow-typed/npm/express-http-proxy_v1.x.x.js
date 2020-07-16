/* Types inferred from documentation:
 *  https://www.npmjs.com/package/express-http-proxy
 * 
 * The types are not entirely complete (e.g. most of the proxy options 
 * can also accept promises) but they cover all the use-cases we have
 */ 

import type {ExpressRequest, ExpressResponse,  Middleware} from 'express';

declare type Headers = {[string]: string}

declare type ProxyOptions = {
  proxyReqPathResolver: (req: ExpressRequest) => string,
  filter: (req: ExpressRequest, res: ExpressResponse) => boolean,
  userResDecorator: (
    proxyRes: ExpressResponse,
    proxyResData: Buffer,
    userReq: ExpressRequest,
    userRes: ExpressResponse
   ) => string,
  memoizeHost: boolean,
  userResHeaderDecorator: (
    headers: Headers,
    userReq: ExpressRequest,
    userRes: ExpressResponse,
    proxyReq: ExpressRequest,
    proxyRes: ExpressResponse
  ) => Headers,
  // This is experimental and may change in upcoming versions
  skipToNextHandlerFilter: (proxyRes: ExpressResponse) => boolean,
  proxyErrorHandler: (err: any, res: ExpressResponse, next: Middleware) => void,
  proxyReqOptDecorator: (proxyReqOpts: any, srcReq: ExpressRequest) => void,
  proxyReqBodyDecorator: (bodyContent: string, srcReq: ExpressRequest) => void,
  https: boolean,
  preserveHostHdr: boolean,
  parseReqBody: boolean,
  reqAsBuffer: boolean,
  reqBodyEncoding: string,
  timeout: number
}

declare module 'express-http-proxy' {
  declare export type ProxyOptions = ProxyOptions;

  declare module.exports: {
    // If you try to call like a function, it will use this signature
    (host: (string | () => string), options?: ProxyOptions): Middleware,
  };
}
