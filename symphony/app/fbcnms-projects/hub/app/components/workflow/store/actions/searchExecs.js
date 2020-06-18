import { sortBy } from "lodash";
import { HttpClient as http } from "../../common/HttpClient";
import { conductorApiUrlPrefix } from "../../constants";

export const RECEIVE_NEW_DATA = "RECEIVE_NEW_DATA";
export const HIERARCHY_NEW_DATA = "HIERARCHY_NEW_DATA";
export const UPDATE_LABEL = "UPDATE_LABEL";
export const UPDATE_QUERY = "UPDATE_QUERY";
export const DATA_SIZE = "DATA_SIZE";
export const CHECKED_WORKFLOWS = "CHECKED_WORKFLOWS";

export const updateSize = size => {
  return { type: DATA_SIZE, size };
};

export const updateLabel = label => {
  return { type: UPDATE_LABEL, label };
};

export const updateQuery = query => {
  return { type: UPDATE_QUERY, query };
};

const createQuery = ({ query, label }) => {
  let q = "",
    search = "";
  if (query) {
    for (let i = 0; i < query.length; i++) {
      search += "[" + query[i].toUpperCase() + query[i].toLowerCase() + "]";
    }
    q += "(workflowId:" + query + "+workflowType:/.*" + search + ".*/)";
  }
  if (label.length) {
    if (query) q += "AND";
    q += "(status:" + label + ")";
  }
  return q;
};

export const fetchNewData = (viewedPage, defaultPages, backendApiUrlPrefix) => {
  return (dispatch, getState) => {
    let q = createQuery(getState().searchReducer);
    let page = (viewedPage - 1) * defaultPages;
    http
      .get(
        backendApiUrlPrefix + "/executions/?q=&h=&freeText=" +
          q +
          "&start=" +
          page +
          "&size=" +
          defaultPages
      )
      .then(res => {
        if (res.result) {
          const data = res.result
            ? res.result.hits
              ? res.result.hits
              : []
            : [];
          dispatch(updateSize(res.result.totalHits));
          dispatch(receiveNewData(data));
        }
      });
  };
};

export const receiveNewData = data => {
  return { type: RECEIVE_NEW_DATA, data };
};

export const fetchParentWorkflows = (viewedPage, defaultPages, backendApiUrlPrefix) => {
  return (dispatch, getState) => {
    let page = viewedPage - 1;

    const { checkedWfs, size } = getState().searchReducer;
    let q = createQuery(getState().searchReducer);
    http
      .get(
        backendApiUrlPrefix + "/hierarchical/?freeText=" +
          q +
          "&start=" +
          checkedWfs[page] +
          "&size=" +
          defaultPages
      )
      .then(res => {
        if (res) {
          let parents = res.parents ? res.parents : [];
          let children = res.children ? res.children : [];
          if (
            res.count < res.hits &&
            (typeof checkedWfs[viewedPage] === "undefined" ||
              checkedWfs.length === 1)
          ) {
            checkedWfs.push(res.count);
            dispatch(updateSize(size + parents.length));
          }
          dispatch(checkedWorkflows(checkedWfs));
          parents = sortBy(parents, wf => new Date(wf.startTime)).reverse();
          dispatch(receiveHierarchicalData(parents, children));
        }
      });
  };
};

export const receiveHierarchicalData = (parents, children) => {
  return { type: HIERARCHY_NEW_DATA, parents, children };
};

export const checkedWorkflows = checkedWfs => {
  return { type: CHECKED_WORKFLOWS, checkedWfs };
};

export const updateParents = childInput => {
  return (dispatch, getState) => {
    const { parents, children } = getState().searchReducer;
    let dataset = parents;
    dataset.forEach((wfs, i) => {
      if (childInput.some(e => e.parentWorkflowId === wfs.workflowId)) {
        let unfoldChildren = childInput.filter(
          wf => wf.parentWorkflowId === wfs["workflowId"]
        );
        unfoldChildren = sortBy(unfoldChildren, wf => new Date(wf.startTime));
        unfoldChildren.forEach((wf, index) =>
          dataset.splice(index + 1 + i, 0, wf)
        );
      }
    });
    dispatch(receiveHierarchicalData(dataset, children));
  };
};

export const deleteParents = childInput => {
  return (dispatch, getState) => {
    const { parents, children } = getState().searchReducer;
    let dataset = parents;
    childInput.forEach(wfs => {
      dataset = dataset.filter(p => p.workflowId !== wfs.workflowId);
    });
    dispatch(receiveHierarchicalData(dataset, children));
  };
};
