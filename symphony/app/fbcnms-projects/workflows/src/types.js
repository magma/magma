/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type Task = {
  name: string,
  type: string,
  defaultCase?: Array<Task>,
  decisionCases?: Map<string, Array<Task>>,
  forkTasks?: Array<Task>,
  subWorkflowParam?: {name: string},
};

export type ProxyRequest = any;
export type ProxyResponse = any;
export type ProxyNext = any => void;

export type ProxyCallback = (proxyOptions?: any) => void;

export type BeforeFun = (
  tenantId: string,
  req: ProxyRequest,
  res: ProxyResponse,
  proxyCallback: ProxyCallback,
) => void;

export type AfterFun = (
  tenantId: string,
  req: ProxyRequest,
  respObj: ?any,
  res: ProxyResponse,
) => void;

export type TransformerCtx = {
  proxyTarget: string,
};

export type HttpMethod = 'get' | 'post' | 'delete' | 'put';

type ExpressCallback = (
  req: ProxyRequest,
  res: ProxyResponse,
  next: ProxyNext,
) => any;
type ExpressMethodFun = (string, ExpressCallback) => void;
export type ExpressRouter = {[HttpMethod]: ExpressMethodFun};

export type TransformerEntry = {
  method: HttpMethod,
  url: string,
  before?: BeforeFun,
  after?: AfterFun,
};

export type TransformerRegistrationFun = (
  ctx: TransformerCtx,
) => Array<TransformerEntry>;

export type Event = {
  name: string,
  event: string,
};

export type Workflow = {
  name: string,
  tasks: Array<Task>,
};

export type StartWorkflowRequest = {
  name: string,
  workflowDef?: any,
  taskToDomain?: any,
};
