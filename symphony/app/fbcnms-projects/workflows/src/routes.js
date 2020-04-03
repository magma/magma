/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const moment = require('moment');
const filter = require('lodash/fp/filter');
const forEach = require('lodash/fp/forEach');
const map = require('lodash/fp/map');
const transform = require('lodash/transform');
const identity = require('lodash/identity');

const Router = require('express');
const http = require('./HttpServerSide').HttpClient;

const router = new Router();
//TODO make configurable
const baseURL = 'conductor-server:8080/api/';
const baseURLWorkflow = baseURL + 'workflow/';
const baseURLMeta = baseURL + 'metadata/';
const baseURLTask = baseURL + 'tasks/';

router.get('/metadata/taskdefs', async (req, res, next) => {
  try {
    const result = await http.get(baseURLMeta + 'taskdefs', req.token);
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.get('/metadata/taskdefs', async (req, res, next) => {
  try {
    const result = await http.get(baseURLMeta + 'taskdefs', req.token);
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.post('/metadata/taskdef', async (req, res, next) => {
  try {
    const result = await http.post(baseURLMeta + 'taskdefs', req.body);
    res.status(200).send(result);
  } catch (err) {
    res.status(400).send(err.response.body);
    next(err);
  }
});

router.get('/metadata/taskdef/:name', async (req, res, next) => {
  try {
    const result = await http.get(
      baseURLMeta + 'taskdefs/' + req.params.name,
      req.token,
    );
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.delete('/metadata/taskdef/:name', async (req, res, next) => {
  try {
    const result = await http.delete(
      baseURLMeta + 'taskdefs/' + req.params.name,
      req.token,
    );
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.get('/metadata/workflow', async (req, res, next) => {
  try {
    const result = await http.get(baseURLMeta + 'workflow', req.token);
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.delete('/metadata/workflow/:name/:version', async (req, res, next) => {
  try {
    const result = await http.delete(
      baseURLMeta + 'workflow/' + req.params.name + '/' + req.params.version,
      req.token,
    );
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.get('/metadata/workflow/:name/:version', async (req, res, next) => {
  try {
    const result = await http.get(
      baseURLMeta +
        'workflow/' +
        req.params.name +
        '?version=' +
        req.params.version,
      req.token,
    );
    res.status(200).send({result});
  } catch (err) {
    next(err);
  }
});

router.put('/metadata', async (req, res, next) => {
  try {
    const result = await http.put(baseURLMeta + 'workflow/', req.body);
    res.status(200).send(result);
  } catch (err) {
    res.status(400).send(err.response.body);
    next(err);
  }
});

router.post('/workflow', async (req, res, next) => {
  try {
    const result = await http.post(baseURLWorkflow, req.body);
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.get('/executions', async (req, res, next) => {
  try {
    const freeText = [];
    if (req.query.freeText !== '') {
      freeText.push(req.query.freeText);
    } else {
      freeText.push('*');
    }

    let h = '-1';
    if (req.query.h !== 'undefined' && req.query.h !== '') {
      h = req.query.h;
    }
    if (h !== '-1') {
      freeText.push('startTime:[now-' + h + 'h TO now]');
    }
    let start = 0;
    if (!isNaN(req.query.start)) {
      start = req.query.start;
    }
    let size = 1000;
    if (req.query.size !== 'undefined' && req.query.size !== '') {
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
      encodeURIComponent(query);
    const result = await http.get(url, req.token);
    const hits = result.results;
    res.status(200).send({result: {hits: hits, totalHits: result.totalHits}});
  } catch (err) {
    next(err);
  }
});

router.delete('/bulk/terminate', async (req, res, next) => {
  try {
    const result = await http.delete(
      baseURLWorkflow + 'bulk/terminate',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.put('/bulk/pause', async (req, res, next) => {
  try {
    const result = await http.put(
      baseURLWorkflow + 'bulk/pause',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.put('/bulk/resume', async (req, res, next) => {
  try {
    const result = await http.put(
      baseURLWorkflow + 'bulk/resume',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.post('/bulk/retry', async (req, res, next) => {
  try {
    const result = await http.post(
      baseURLWorkflow + 'bulk/retry',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.post('/bulk/restart', async (req, res, next) => {
  try {
    const result = await http.post(
      baseURLWorkflow + 'bulk/restart',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.delete('/workflow/:workflowId', async (req, res, next) => {
  try {
    const result = await http.delete(
      baseURLWorkflow + req.params.workflowId + '/remove',
      req.body,
      req.token,
    );
    res.status(200).send(result);
  } catch (err) {
    next(err);
  }
});

router.get('/id/:workflowId', async (req, res, next) => {
  try {
    const result = await http.get(
      baseURLWorkflow + req.params.workflowId + '?includeTasks=true',
      req.token,
    );
    let meta = result.workflowDefinition;
    if (!meta) {
      meta = await http.get(
        baseURLMeta +
          'workflow/' +
          result.workflowType +
          '?version=' +
          result.version,
        req.token,
      );
    }

    const subs = filter(identity)(
      map(task => {
        if (task.taskType === 'SUB_WORKFLOW') {
          const subWorkflowId = task.inputData && task.inputData.subWorkflowId;

          if (subWorkflowId != null) {
            return {
              name: task.inputData.subWorkflowName,
              version: task.inputData.subWorkflowVersion,
              referenceTaskName: task.referenceTaskName,
              subWorkflowId: subWorkflowId,
            };
          }
        }
      })(result.tasks || []),
    );

    (result.tasks || []).forEach(task => {
      if (task.taskType === 'SUB_WORKFLOW') {
        const subWorkflowId = task.inputData && task.inputData.subWorkflowId;

        if (subWorkflowId != null) {
          subs.push({
            name: task.inputData.subWorkflowName,
            version: task.inputData.subWorkflowVersion,
            referenceTaskName: task.referenceTaskName,
            subWorkflowId: subWorkflowId,
          });
        }
      }
    });

    const logs = map(task =>
      Promise.all([task, http.get(baseURLTask + task.taskId + '/log')]),
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

    const promises = map(({name, version, subWorkflowId, referenceTaskName}) =>
      Promise.all([
        referenceTaskName,
        http.get(baseURLMeta + 'workflow/' + name + '?version=' + version),
        http.get(baseURLWorkflow + subWorkflowId + '?includeTasks=true'),
      ]),
    )(subs);

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
    next(err);
  }
});

router.get('/hierarchical', async (req, res, next) => {
  try {
    let size = 1000;
    if (req.query.size !== 'undefined' && req.query.size !== '') {
      size = req.query.size;
    }

    let count = 0;
    let start = 0;
    if (!isNaN(req.query.start)) {
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
      const result = await http.get(url, req.token);
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
                req.token,
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
            req.token,
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

router.get('/queue/data', async (req, res, next) => {
  try {
    const sizes = await http.get(baseURLTask + 'queue/all', req.token);
    const polldata = await http.get(
      baseURLTask + 'queue/polldata/all',
      req.token,
    );
    polldata.forEach(pd => {
      let qname = pd.queueName;

      if (pd.domain != null) {
        qname = pd.domain + ':' + qname;
      }
      pd.qsize = sizes[qname];
    });
    res.status(200).send({polldata});
  } catch (err) {
    next(err);
  }
});

module.exports = router;
