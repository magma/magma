/* eslint-disable flowtype/require-valid-file-annotation */
/* Disabled flow and also disabled eslint flow annotation check TODO */
/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

import HttpClient from './HttpServerSide';
import Router from 'express';
import bodyParser from 'body-parser';
import filter from 'lodash/fp/filter';
import forEach from 'lodash/fp/forEach';
import identity from 'lodash/fp/identity';
import logging from '@fbcnms/logging';
import map from 'lodash/fp/map';
import moment from 'moment';
import transform from 'lodash/fp/transform';

import type {$Application, ExpressRequest, ExpressResponse} from 'express';
import type {TaskType} from './types';

const logger = logging.getLogger(module);
const http = HttpClient;

const findSchedule = (schedules, name) => {
  for (const schedule of schedules) {
    if (schedule.name === name) {
      return schedule;
    }
  }
  return null;
};

//TODO merge with proxy
export default async function(
  baseURL: string,
  addScheduleMetadata: boolean,
): Promise<$Application<ExpressRequest, ExpressResponse>> {
  const router = Router();
  const baseApiURL = baseURL + 'api/';
  const baseURLWorkflow = baseApiURL + 'workflow/';
  const baseURLMeta = baseApiURL + 'metadata/';
  const baseURLTask = baseApiURL + 'tasks/';
  const eventURL = baseApiURL + 'event';
  const baseURLSchedule = baseURL + 'schedule/';

  router.use(bodyParser.urlencoded({extended: false}));
  router.use('/', bodyParser.json());

  router.get('/metadata/taskdefs', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.get(baseURLMeta + 'taskdefs', req);
      res.status(200).send({result});
    } catch (err) {
      next(err);
    }
  });

  router.post('/metadata/taskdef', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.post(baseURLMeta + 'taskdefs', req.body, req);
      res.status(200).send(result);
    } catch (err) {
      res.status(400).send(err.response.body);
      next(err);
    }
  });

  router.get(
    '/metadata/taskdef/:name',
    async (req: ExpressRequest, res, next) => {
      try {
        const result = await http.get(
          baseURLMeta + 'taskdefs/' + req.params.name,
          req,
        );
        res.status(200).send({result});
      } catch (err) {
        next(err);
      }
    },
  );

  router.delete(
    '/metadata/taskdef/:name',
    async (req: ExpressRequest, res, next) => {
      try {
        const result = await http.delete(
          baseURLMeta + 'taskdefs/' + req.params.name,
          null,
          req,
        );
        res.status(200).send({result});
      } catch (err) {
        next(err);
      }
    },
  );

  router.get('/metadata/workflow', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.get(baseURLMeta + 'workflow', req);
      // combine with schedules
      // FIXME this should not be here, schedules must be isolated from
      // the rest of workflow API !
      if (addScheduleMetadata) {
        const schedules = await http.get(baseURLSchedule, req);
        for (const workflowDef of result) {
          const expectedScheduleName =
            workflowDef.name + ':' + workflowDef.version;
          const found = findSchedule(schedules, expectedScheduleName);
          workflowDef.hasSchedule = found != null;
          workflowDef.expectedScheduleName = expectedScheduleName;
        }
      }

      res.status(200).send({result});
    } catch (err) {
      next(err);
    }
  });

  router.delete(
    '/metadata/workflow/:name/:version',
    async (req: ExpressRequest, res, next) => {
      try {
        const result = await http.delete(
          baseURLMeta +
            'workflow/' +
            req.params.name +
            '/' +
            req.params.version,
          null,
          req,
        );
        res.status(200).send({result});
      } catch (err) {
        next(err);
      }
    },
  );

  router.get(
    '/metadata/workflow/:name/:version',
    async (req: ExpressRequest, res, next) => {
      try {
        const result = await http.get(
          baseURLMeta +
            'workflow/' +
            req.params.name +
            '?version=' +
            req.params.version,
          req,
        );
        res.status(200).send({result});
      } catch (err) {
        next(err);
      }
    },
  );

  router.put('/metadata', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.put(baseURLMeta + 'workflow/', req.body, req);
      res.status(200).send(result);
    } catch (err) {
      res.status(400).send(err?.response?.body);
      next(err);
    }
  });

  // Conductor only allows POST for event handler creation
  // and PUT for updating. This code works around the issue by
  // trying PUT if POST fails.
  router.post('/event', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.post(eventURL, req.body, req);
      res.status(200).send(result);
    } catch (err) {
      logger.info(`Got exception ${JSON.stringify(err)} while POSTing event`);
      try {
        const result = await http.put(eventURL, req.body, req);
        res.status(200).send(result);
      } catch (e) {
        logger.info(`Got exception ${JSON.stringify(e)} while PUTting event`);
        res.status(400).send('Post and Put failed');
        next(e);
      }
    }
  });

  router.post('/workflow', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.post(baseURLWorkflow, req.body, req);
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.get('/executions', async (req: ExpressRequest, res, next) => {
    try {
      const freeText = [];
      if (req.query.freeText !== '') {
        freeText.push(req.query.freeText);
      } else {
        freeText.push('*');
      }

      let h: string = '-1';
      if (req.query.h !== 'undefined' && req.query.h !== '') {
        /* FIXME req.query is user-controlled input, properties and values
         in this object are untrusted and should be validated before trusting */
        h = req.query.h;
      }
      if (h !== '-1') {
        freeText.push('startTime:[now-' + h + 'h TO now]');
      }
      let start: number = 0;
      if (!isNaN(req.query.start)) {
        // FIXME: isNaN is sketchy and accepts arrays
        start = req.query.start;
      }
      let size: number = 1000;
      if (req.query.size !== 'undefined' && req.query.size !== '') {
        /* FIXME req.query is user-controlled input, properties and values
         in this object are untrusted and should be validated before trusting */
        size = req.query.size;
      }

      const query = req.query.q;

      const url =
        baseURLWorkflow +
        'search?size=' +
        size +
        '&sort=startTime:DESC&freeText=' +
        encodeURIComponent(freeText.join(' AND ')) +
        '&start=' +
        start +
        '&query=' +
        /* FIXME: req.query is user-controlled input and could
         be an array. Needs to be checked */
        encodeURIComponent(query);
      const result = await http.get(url, req);
      const hits = result.results;
      res.status(200).send({result: {hits: hits, totalHits: result.totalHits}});
    } catch (err) {
      next(err);
    }
  });

  router.delete('/bulk/terminate', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.delete(
        baseURLWorkflow + 'bulk/terminate',
        req.body,
        req,
      );
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.put('/bulk/pause', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.put(
        baseURLWorkflow + 'bulk/pause',
        req.body,
        req,
      );
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.put('/bulk/resume', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.put(
        baseURLWorkflow + 'bulk/resume',
        req.body,
        req,
      );
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.post('/bulk/retry', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.post(
        baseURLWorkflow + 'bulk/retry',
        req.body,
        req,
      );
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.post('/bulk/restart', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.post(
        baseURLWorkflow + 'bulk/restart',
        req.body,
        req,
      );
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.delete(
    '/workflow/:workflowId',
    async (req: ExpressRequest, res, next) => {
      const archiveWorkflow = req.query.archiveWorkflow === 'true';
      try {
        const result = await http.delete(
          baseURLWorkflow +
            req.params.workflowId +
            '/remove?archiveWorkflow=' +
            archiveWorkflow.toString(),
          req.body,
          req,
        );
        res.status(200).send(result);
      } catch (err) {
        next(err);
      }
    },
  );

  router.get('/id/:workflowId', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.get(
        baseURLWorkflow + req.params.workflowId + '?includeTasks=true',
        req,
      );
      let meta = result.workflowDefinition;
      if (!meta) {
        meta = await http.get(
          baseURLMeta +
            'workflow/' +
            result.workflowType +
            '?version=' +
            result.version,
          req,
        );
      }

      const subs = filter(identity)(
        map((task: TaskType): ?TaskType => {
          if (task.taskType === 'SUB_WORKFLOW' && task.inputData) {
            const subWorkflowId = task.inputData.subWorkflowId;

            if (subWorkflowId != null) {
              return {
                name: task.inputData?.subWorkflowName,
                version: task.inputData?.subWorkflowVersion,
                referenceTaskName: task.referenceTaskName,
                subWorkflowId: subWorkflowId,
              };
            }
          }
        })(result.tasks || []),
      );

      const logs = map(task =>
        Promise.all([task, http.get(baseURLTask + task.taskId + '/log', req)]),
      )(result.tasks);
      const LOG_DATE_FORMAT = 'MM/DD/YY, HH:mm:ss:SSS';

      await Promise.all(logs).then(result => {
        forEach(([task, logs]) => {
          if (logs) {
            task.logs = map(
              ({createdTime, log}) =>
                `${moment(createdTime).format(LOG_DATE_FORMAT)} : ${log}`,
            )(logs);
          }
        })(result);
      });

      const fun: (
        Array<TaskType>,
      ) => Array<
        Promise<mixed>,
      > = map(({name, version, subWorkflowId, referenceTaskName}) =>
        Promise.all([
          referenceTaskName,
          http.get(
            baseURLMeta + 'workflow/' + name + '?version=' + version,
            req,
          ),
          http.get(baseURLWorkflow + subWorkflowId + '?includeTasks=true', req),
        ]),
      );
      const promises = fun(subs);

      const subworkflows = await Promise.all(promises).then(result => {
        return transform(
          result,
          (result, [key, meta, wfe]) => {
            result[key] = {meta, wfe};
          },
          {},
        );
      });

      res.status(200).send({result, meta, subworkflows: subworkflows});
    } catch (err) {
      console.log(err);
      next(err);
    }
  });

  router.get('/hierarchical', async (req: ExpressRequest, res, next) => {
    try {
      let size: number = 1000;
      if (req.query.size !== 'undefined' && req.query.size !== '') {
        /* FIXME req.query is user-controlled input, properties and values
         in this object are untrusted and should be validated before trusting */
        size = req.query.size;
      }

      let count = 0;
      let start: number = 0;
      if (!isNaN(req.query.start)) {
        // FIXME: isNaN is sketchy and accepts arrays
        start = req.query.start;
        count = Number(start);
      }

      const freeText = [];
      if (req.query.freeText !== '') {
        freeText.push(req.query.freeText);
      } else {
        freeText.push('*');
      }

      const parents = [];
      const children = [];

      let hits = 0;
      while (parents.length < size) {
        const url =
          baseURLWorkflow +
          'search?size=' +
          size * 10 +
          '&sort=startTime:DESC&freeText=' +
          encodeURIComponent(freeText.join(' AND ')) +
          '&start=' +
          start +
          '&query=';
        const result = await http.get(url, req);
        const allData = result.results ? result.results : [];
        hits = result.totalHits ? result.totalHits : 0;

        const separatedWfs = [];
        const chunk = 5;

        for (let i = 0, j = allData.length; i < j; i += chunk) {
          separatedWfs.push(allData.slice(i, i + chunk));
        }

        for (let i = 0; i < separatedWfs.length; i++) {
          const wfs = async function(sepWfs) {
            return await Promise.all(
              sepWfs.map(wf =>
                http.get(
                  baseURLWorkflow + wf.workflowId + '?includeTasks=false',
                  req,
                ),
              ),
            );
          };
          let checked = 0;
          const responses = await wfs(separatedWfs[i]);
          for (let j = 0; j < responses.length; j++) {
            if (responses[j].parentWorkflowId) {
              separatedWfs[i][j]['parentWorkflowId'] =
                responses[j].parentWorkflowId;
              children.push(separatedWfs[i][j]);
            } else {
              parents.push(separatedWfs[i][j]);
              if (parents.length === size) {
                checked = j + 1;
                break;
              }
            }
            checked = j + 1;
          }
          count += checked;
          if (parents.length >= size) break;
        }
        if (req.query.freeText !== '') {
          for (let i = 0; i < children.length; i++) {
            const parent = await http.get(
              baseURLWorkflow +
                children[i].parentWorkflowId +
                '?includeTasks=false',
              req,
            );
            parent.startTime = new Date(parent.startTime);
            parent.endTime = new Date(parent.endTime);
            if (
              parent.parentWorkflowId &&
              !children.find(wf => wf.workflowId === parent.workflowId)
            )
              children.push(parent);
            if (
              !parent.parentWorkflowId &&
              !parents.find(wf => wf.workflowId === parent.workflowId)
            )
              parents.push(parent);
          }
        }
        start = Number(start) + size * 10;
        if (Number(start) >= hits) break;
      }
      res.status(200).send({parents, children, count, hits});
    } catch (err) {
      next(err);
    }
  });

  router.get('/schedule/?', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.get(baseURLSchedule, req);
      res.status(200).send(result);
    } catch (err) {
      next(err);
    }
  });

  router.get('/schedule/:name', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.get(baseURLSchedule + req.params.name, req);
      res.status(200).send(result);
    } catch (err) {
      logger.warn('Failed to POST and PUT', {error: err});
      next(err);
    }
  });

  router.put('/schedule/:name', async (req: ExpressRequest, res, next) => {
    const urlWithName = baseURLSchedule + req.params.name;
    try {
      // create using POST
      const result = await http.post(baseURLSchedule, req.body, req);
      res.status(result.statusCode).send(result.text);
    } catch (e) {
      try {
        // update using PUT
        const result = await http.put(urlWithName, req.body, req);
        res.status(result.statusCode).send(result.text);
      } catch (err) {
        next(err);
      }
    }
  });

  router.delete('/schedule/:name', async (req: ExpressRequest, res, next) => {
    try {
      const result = await http.delete(
        baseURLSchedule + req.params.name,
        null,
        req,
      );
      res.status(result.statusCode).send(result.text);
    } catch (err) {
      next(err);
    }
  });

  return router;
}
